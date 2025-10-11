package models

import "time"

type Order struct {
	ID            uint
	CreatedAt     time.Time
	TotalPrice    float64
	AccountID     uint64
	Status        string
	PaymentStatus string
	ProductInfos  []ProductInfo
	Products      []*OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}
