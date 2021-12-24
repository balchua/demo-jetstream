package model

import "github.com/shopspring/decimal"

type UserTransaction struct {
	TransactionID int             `json:"TransactionID"`
	UserID        int             `json:"UserID"`
	Status        string          `json:"Status"`
	Amount        decimal.Decimal `json:"Amount"`
}

type TransactionPublisher interface {
	Publish(message string, subject string) error
}
