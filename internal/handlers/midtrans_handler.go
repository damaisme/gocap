package handlers

import (
	"log"
	"os"

	"github.com/damaisme/gocap/internal/models"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func CheckTransactionStatus(orderID string) bool {

	var client = coreapi.Client{}
	if os.Getenv("ENVIRONMENT") == "production" {
		client.New(os.Getenv("SERVERKEY"), midtrans.Production)
	} else {
		client.New(os.Getenv("SERVERKEY"), midtrans.Sandbox)
	}

	resp, err := client.CheckTransaction(orderID)
	if err != nil {
		log.Fatalf("Failed to check transaction status: %v", err)
	}

	// Process transaction status based on the response
	switch resp.TransactionStatus {
	case "capture":
		if resp.FraudStatus == "accept" {
			log.Printf("Transaction %s is successfully captured\n", orderID)
			return true
		}
	case "settlement":
		if resp.FraudStatus == "accept" {
			log.Printf("Transaction %s is successfully captured\n", orderID)
			return true
		}
	case "deny", "cancel", "expire":
		log.Printf("Transaction %s is %s\n", orderID, resp.TransactionStatus)
		return false
	case "pending":
		log.Printf("Transaction %s is pending\n", orderID)
		return false
	}

	return false
}

func CreateTransaction(name string, GrossAmt int64, plan string, OrderID *uuid.UUID) model.MidtransResponse {

	var snapClient = snap.Client{}

	if os.Getenv("ENVIRONMENT") == "production" {
		snapClient.New(os.Getenv("SERVERKEY"), midtrans.Production)
	} else {
		snapClient.New(os.Getenv("SERVERKEY"), midtrans.Sandbox)
	}

	// customer
	custAddress := &midtrans.CustomerAddress{
		FName:       "John",
		LName:       "Doe",
		Phone:       "081234567890",
		Address:     "Baker Street 97th",
		City:        "Jakarta",
		Postcode:    "16000",
		CountryCode: "IDN",
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  OrderID.String(),
			GrossAmt: GrossAmt,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    name,
			LName:    "Filo",
			Email:    "john@doe.com",
			Phone:    "081234567890",
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "Voc-" + plan,
				Qty:   1,
				Price: int64(GrossAmt),
				Name:  plan,
			},
		},
	}

	response, errSnap := snapClient.CreateTransaction(req)
	if errSnap != nil {
		log.Panic(errSnap.GetRawError())
	}

	midtransReponse := model.MidtransResponse{
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransReponse
}
