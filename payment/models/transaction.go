package models

type Transaction struct {
	OrderId      uint64
	UserId       uint64
	CustomerId   string
	PaymentId    string
	TotalPrice   int64
	SettledPrice int64
	Currency     string
	Status       string
}
