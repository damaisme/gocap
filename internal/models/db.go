package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string
	Expiry   time.Time
	Ip       string
}

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;"`
	GrossAmt      int64
	PaymentStatus string
}

type Voucher struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code          string    `gorm:"uniqueIndex"`
	Expiry        time.Time
	Name          string
	MaxUses       int
	Uses          int
	IsActive      bool
	TransactionID uuid.UUID
	Transaction   Transaction
	Ip            string
}
