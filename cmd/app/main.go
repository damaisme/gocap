package main

import (
	"time"

	"github.com/damaisme/gocap/internal/config"
	"github.com/damaisme/gocap/internal/database"
	"github.com/damaisme/gocap/internal/models"
	"github.com/damaisme/gocap/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var (
	// Define a cookie store for session management
	store = sessions.NewCookieStore([]byte("your-secret-key"))
)

func main() {

	database.InitDB()

	config.InitSession()

	// Seed a voucher (for testing)
	database.DB.Create(&model.Voucher{Code: "VOUCHER123", Expiry: time.Now().Add(24 * time.Hour), MaxUses: 2})
	database.DB.Create(&model.User{Username: "test", Password: "test", Expiry: time.Now().Add(24 * 30 * time.Hour)})

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	routes.RegisterRoutes(router)

	router.Run(":8080")
}
