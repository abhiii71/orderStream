package models

type EventData struct {
	Id          *string  `json:"product_id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	AccountID   *int     `json:"accountID,omitempty"`
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}

// ProductsListedEventData for when multiple products are listed
type ProductsListedEventData struct {
	ProductIDs []string `json:"product_ids"`
	Count      int      `json:"count"`
}

type ProductsListedEvent struct {
	Type string                  `json:"type"`
	Data ProductsListedEventData `json:"data"`
}
