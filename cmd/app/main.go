package main

import (
	"time"

	"github.com/damaisme/gocap/internal/config"
	"github.com/damaisme/gocap/internal/database"
	"github.com/damaisme/gocap/internal/models"
	"github.com/damaisme/gocap/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"log"
	"os"
)

var (
	// Define a cookie store for session management
	store = sessions.NewCookieStore([]byte("your-secret-key"))
)

func main() {

	config.LoadConfig()

	database.InitDB()

	config.InitSession()

	// Seed a voucher (for testing)
	database.DB.Create(&model.Voucher{Code: "VOUCHER123", Expiry: time.Now().Add(24 * time.Hour), MaxUses: 2})
	database.DB.Create(&model.User{Username: "test", Password: "test", Expiry: time.Now().Add(24 * 30 * time.Hour)})

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/public", "public/")

	routes.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // Default port
	}
	log.Printf("Starting server on port: %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
