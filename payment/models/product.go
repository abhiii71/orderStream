package models

import "time"

type Product struct {
	ID            uint64
	ProductID     string
	DodoProductID string
	Price         int64
	Currency      string

	CreatedAt time.Time
	UpdatedAt time.Time
}
