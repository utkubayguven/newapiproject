// models/models.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// User Model
type User struct {
	gorm.Model
	Username    string `gorm:"unique;not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	PhoneNumber string `gorm:"unique;not null"`
	PIN         string `gorm:"not null"`
}

// Account Model
type Account struct {
	gorm.Model
	UserID           uint             `gorm:"not null"` // Kullanıcı ID'sine referans
	Balance          int              `gorm:"default:1000"`
	Deposits         []Deposit        `gorm:"foreignKey:AccountID"` // Bir hesaba birden fazla para yatırma işlemi yapılabilir
	Withdrawals      []Withdrawal     `gorm:"foreignKey:AccountID"` // Bir hesaptan birden fazla para çekme işlemi olabilir
	BalanceInquiries []BalanceInquiry `gorm:"foreignKey:AccountID"` // Bir hesaba ait birden fazla bakiye sorgulama kaydı olabilir
}

// Deposit Model
type Deposit struct {
	gorm.Model
	AccountID     uint    `gorm:"not null"` // Hesap ID'sine referans
	Account       Account `gorm:"foreignKey:AccountID;references:ID"`
	DepositAmount int
	DepositDate   time.Time `gorm:"default:current_timestamp"`
}

// Withdrawal Model
type Withdrawal struct {
	gorm.Model
	AccountID        uint    `gorm:"not null"` // Hesap ID'sine referans
	Account          Account `gorm:"foreignKey:AccountID;references:ID"`
	WithdrawalAmount int
	WithdrawalDate   time.Time `gorm:"default:current_timestamp"`
}

// BalanceInquiry Model
type BalanceInquiry struct {
	gorm.Model
	AccountID      uint    `gorm:"not null"` // Hesap ID'sine referans
	Account        Account `gorm:"foreignKey:AccountID;references:ID"`
	CurrentBalance int
	InquiryDate    time.Time `gorm:"default:current_timestamp"`
}

// PinChange Model
type PinChange struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`          // Kullanıcı ID'sine referans
	User       User      `gorm:"foreignKey:UserID"` // Belongs-to ilişki tanımı
	OldPIN     string    `gorm:"not null"`
	NewPIN     string    `gorm:"not null"`
	ChangeDate time.Time `gorm:"default:current_timestamp"`
}
