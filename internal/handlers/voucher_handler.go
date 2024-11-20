package handlers

import (
	"github.com/damaisme/gocap/internal/database"
	"github.com/damaisme/gocap/internal/models"
	"time"
)

// Validate and use a voucher
func ValidateVoucher(code string) (bool, model.Voucher) {
	var voucher model.Voucher
	result := database.DB.Where("code = ?", code).First(&voucher)
	if result.Error != nil || time.Now().After(voucher.Expiry) || voucher.Uses >= voucher.MaxUses {
		return false, voucher
	}

	// Increment usage count
	database.DB.Model(&voucher).Update("Uses", voucher.Uses+1)
	return true, voucher
}
