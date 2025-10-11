package models

import "time"

type Customer struct {
	UserId       uint64
	CustomerId   string
	BillingEmail string
	CreatedAt    time.Time
	Transactions []Transaction
}

type CustomerInput struct {
	UserId string
}
