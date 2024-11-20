package handlers

import (
	"github.com/damaisme/gocap/internal/database"
	"github.com/damaisme/gocap/internal/models"
	"github.com/damaisme/gocap/internal/utils"
	"log"
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

func StartVoucherExpiryCheck() {
	// Set the interval to check every minute (adjust as needed)
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				checkAndRemoveExpiredVouchers()
			}
		}
	}()
}

func checkAndRemoveExpiredVouchers() {
	now := time.Now()

	// Query for all vouchers that are expired and active
	var vouchers []model.Voucher
	if err := database.DB.Where("expiry < ? AND is_active = ?", now, true).Find(&vouchers).Error; err != nil {
		log.Printf("Error fetching expired vouchers: %v", err)
		return
	}

	for _, voucher := range vouchers {
		// Remove the iptables rule for the expired voucher's IP
		utils.DeleteIptablesRule(voucher.Ip)

		// Mark the voucher as inactive (or delete it based on your logic)
		if err := database.DB.Model(&voucher).Update("is_active", false).Error; err != nil {
			log.Printf("Error updating voucher status: %v", err)
		}

		// Optionally log the action
		log.Printf("Voucher %s expired and iptables rule removed for IP %s", voucher.Code, voucher.Ip)
	}
}
