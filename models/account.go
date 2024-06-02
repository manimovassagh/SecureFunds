package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID       uint          `gorm:"not null" json:"user_id"`
	AccountType  string        `gorm:"not null" json:"account_type"`
	Balance      float64       `gorm:"not null" json:"balance"`
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"transactions"`
	User         User          `gorm:"foreignKey:UserID" json:"-"`
}
