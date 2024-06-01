package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	AccountID uint
	Amount    float64 `gorm:"not null"`
	Type      string  `gorm:"not null"`
}
