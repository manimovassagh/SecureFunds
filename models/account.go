package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID       uint
	AccountType  string  `gorm:"not null"`
	Balance      float64 `gorm:"not null"`
	Transactions []Transaction
}
