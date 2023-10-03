package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Player struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}

type Wallet struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

type Account struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Username      string `json:"username"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
}

type BankAccount struct {
	ID            int    `gorm:"primaryKey" json:"id"`
	Username      string `json:"username"`
	AccountName   string `json:"accountName"`   // Ensure this matches the JSON field
	AccountNumber string `json:"AccountNumber"` // Ensure this matches the JSON field
	BankName      string `json:"bankName"`
}

type PlayerResponse struct {
	ID            uint    `gorm:"primary_key" json:"id"`
	Username      string  `json:"username"`
	Password      string  `json:"-"`
	Email         string  `json:"email"`
	WalletID      uint    `json:"wallet_id"`
	WalletBalance float64 `json:"wallet_balance"`
	AccountID     uint    `json:"account_id"`
	AccountName   string  `json:"account_name"`
	AccountNumber string  `json:"account_number"`
	BankName      string  `json:"bank_name"`
}

type StoredPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
