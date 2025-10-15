package models

import "time"

type Order struct {
	ID         uint64             `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	TotalPrice float64            `json:"total_price"`
	Products   []*OrderedProducts `json:"products"`
}

type OrderedProducts struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}
