package models

type ProductEventData struct {
	ProductID   *string
	Name        *string
	Description *string
	Price       *float64
	AccountID   *int
}

type ProductEvent struct {
	Type string
	Data ProductEventData
}
