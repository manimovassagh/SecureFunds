package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	AccountID uint    `gorm:"not null" json:"account_id"`
	Amount    float64 `gorm:"not null" json:"amount"`
	Type      string  `gorm:"not null" json:"type"` // "deposit" or "withdrawal"
}
