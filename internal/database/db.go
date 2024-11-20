package database

import (
	"github.com/damaisme/go-captive-portal/internal/models"
	"gorm.io/driver/sqlite" // Example with SQLite, use your DB driver here
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	DB.AutoMigrate(&model.User{}, &model.Voucher{}, &model.Transaction{})

}
