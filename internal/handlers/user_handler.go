package handlers

import (
	"github.com/damaisme/go-captive-portal/internal/config"
	"github.com/damaisme/go-captive-portal/internal/database"
	"github.com/damaisme/go-captive-portal/internal/models"
	"github.com/damaisme/go-captive-portal/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func GetIndex(c *gin.Context) {
	session, err := config.Store.Get(c.Request, "session")
	if err != nil {
		log.Panic("error get session")
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		c.HTML(http.StatusOK, "index.html", gin.H{"voucherCode": session.Values["voucherCode"], "voucherExpiry": session.Values["voucherExpiry"]})
	} else {
		c.Redirect(http.StatusSeeOther, "/login")
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	voucherCode := c.PostForm("voucher")

	if username == "" && voucherCode == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{"Error": "Username and password or voucher code are required."})
		return
	}

	if username != "" && AuthenticateUser(username, password) {
		// Start session
		session, _ := config.Store.Get(c.Request, "session")
		clientIP := c.ClientIP()

		session.Values["authenticated"] = true
		session.Values["loginType"] = "user"
		session.Values["ip"] = clientIP // Store client IP in session
		session.Values["username"] = username
		err := session.Save(c.Request, c.Writer)
		if err != nil {
			log.Printf("Failed to save session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		utils.AddIptablesRule(clientIP)

		c.Redirect(http.StatusSeeOther, "/")
	} else if voucherCode != "" {

		valid, voucher := ValidateVoucher(voucherCode)

		log.Println(valid)

		if valid {
			// Start session
			session, _ := config.Store.Get(c.Request, "session")
			clientIP := c.ClientIP()

			session.Values["authenticated"] = true
			session.Values["loginType"] = "voucher"
			session.Values["ip"] = clientIP // Store client IP in session
			session.Values["voucherCode"] = voucher.Code
			session.Values["voucherExpiry"] = voucher.Expiry
			session.Save(c.Request, c.Writer)

			utils.AddIptablesRule(c.ClientIP())
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Access denied"})
		return

	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Access denied"})
		return
	}
}

func Logout(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	clientIP, _ := session.Values["ip"].(string)
	voucherCode, _ := session.Values["voucherCode"].(string)

	if voucherCode != "" {
		var voucher model.Voucher
		result := database.DB.Where("code = ?", voucherCode).First(&voucher)
		if result.Error != nil {
			database.DB.Model(&voucher).Update("Uses", voucher.Uses-1)
		}
	}

	utils.DeleteIptablesRule(clientIP)

	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusSeeOther, "/login")
}

func BuyVoucher(c *gin.Context) {
	plan := c.PostForm("plan")
	name := c.PostForm("name")

	// if err := c.Shouldatabase.DBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	log.Println(&req)
	// 	return
	// }
	//
	var GrossAmt int64
	var Expiry time.Time

	if plan == "1day" {
		GrossAmt = 5000
		Expiry = time.Now().Add(24 * time.Hour)
	}

	// Save transaction details to the database
	transaction := model.Transaction{
		ID:            uuid.New(),
		GrossAmt:      GrossAmt,
		PaymentStatus: "pending",
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	voucher := model.Voucher{
		ID:            uuid.New(),
		Code:          uuid.New().String(),
		Expiry:        Expiry,
		MaxUses:       1,
		Uses:          0,
		Name:          name,
		TransactionID: transaction.ID,
		Transaction:   transaction,
	}

	if err := database.DB.Create(&voucher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	snapResp := CreateTransaction(name, GrossAmt, plan, &voucher.TransactionID)
	c.Redirect(http.StatusFound, snapResp.RedirectUrl)

}

func Finish(c *gin.Context) {
	// Retrieve query parameters from the URL
	orderID, _ := uuid.Parse(c.Query("order_id"))

	statusCode := c.Query("status_code")
	transactionStatus := c.Query("transaction_status")

	log.Printf("Order ID: %s, Status Code: %s, Transaction Status: %s", orderID, statusCode, transactionStatus)
	log.Println(statusCode)
	log.Println(transactionStatus)
	// Check the transaction status and respond accordingly
	switch transactionStatus {
	case "settlement":
		if statusCode == "200" {

			if CheckTransactionStatus(orderID.String()) {
				var voucher model.Voucher
				result := database.DB.Where(&model.Voucher{TransactionID: orderID}).First(&voucher)
				if result.Error != nil {
					// Handle error, e.g., record not found
					log.Panic(result.Error)
				}

				// set voucher to IsActive and save
				voucher.IsActive = true
				if err := database.DB.Save(&voucher).Error; err != nil {
					log.Printf("Error updating transaction status: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
					return
				}

				session, _ := config.Store.Get(c.Request, "session")
				clientIP := c.ClientIP()

				session.Values["authenticated"] = true
				session.Values["loginType"] = "voucher"
				session.Values["ip"] = clientIP
				session.Values["voucherCode"] = voucher.Code
				session.Values["voucherExpiry"] = voucher.Expiry
				err := session.Save(c.Request, c.Writer)
				if err != nil {
					log.Printf("Failed to save session: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
					return
				}

				utils.AddIptablesRule(c.ClientIP())
				c.Redirect(http.StatusSeeOther, "/")
				return
			}
		} else {
			// Handle potential issues with capture
			c.JSON(http.StatusOK, gin.H{"status": "Payment capture issue", "order_id": orderID})
		}
	default:
		c.HTML(http.StatusOK, "voc.html", gin.H{"Error": "Payment Error, order_id: " + orderID.String()})
		// c.JSON(http.StatusOK, gin.H{"status": "Unknown payment status", "order_id": orderID})
	}

}

// Authenticate user with SQLite database
func AuthenticateUser(username, password string) bool {
	var user model.User
	result := database.DB.Where("username = ? AND password = ?", username, password).First(&user)
	return result.Error == nil
}
