package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"-"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	PIN         string    `json:"pin"`
}

// Account Model
type Account struct {
	ID               uuid.UUID        `json:"id"`
	UserID           uuid.UUID        `json:"user_id"`
	Balance          int              `json:"balance"`
	Deposits         []Deposit        `json:"deposits"`
	Withdrawals      []Withdrawal     `json:"withdrawals"`
	BalanceInquiries []BalanceInquiry `json:"balance_inquiries"`
}

// Deposit Model
type Deposit struct {
	ID            uuid.UUID `json:"id"`
	AccountID     uuid.UUID `json:"account_id"`
	DepositAmount int       `json:"deposit_amount"`
	DepositDate   time.Time `json:"deposit_date"`
}

// Withdrawal Model
type Withdrawal struct {
	ID               uuid.UUID `json:"id"`
	AccountID        uuid.UUID `json:"account_id"`
	WithdrawalAmount int       `json:"withdrawal_amount"`
	WithdrawalDate   time.Time `json:"withdrawal_date"`
}

// BalanceInquiry Model
type BalanceInquiry struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	CurrentBalance int       `json:"current_balance"`
	InquiryDate    time.Time `json:"inquiry_date"`
}

// PinChange Model
type PinChange struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	OldPIN     string    `json:"old_pin"`
	NewPIN     string    `json:"new_pin"`
	ChangeDate time.Time `json:"change_date"`
}
