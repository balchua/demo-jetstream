package model

import "google.golang.org/genproto/googleapis/type/decimal"

type UserTransaction struct {
	TransactionID int
	UserID        int
	Status        string
	Amount        decimal.Decimal
}

type TransactionPublisher interface {
	Publish(message string, subject string) error
}
