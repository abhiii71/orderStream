package dto

import "github.com/dodopayments/dodopayments-go"

type WebhookMetadata struct {
	OrderId uint64 `json:"order_id"`
	UserId  uint64 `json:"user_id"`
}


// WebhookPayload represents the full webhook body.
type WebhookPayload struct {
	Type string       `json:"type"`
	Data WebhookData  `json:"data"`
}

// CustomerInfo holds basic customer details.
type CustomerInfo struct {
	CustomerID string `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}


// ProductItem represents one product in the cart.
type ProductItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}


type WebhookData struct {
	Customer      CustomerInfo          `json:"customer"`
	ProductCart   []ProductItem         `json:"product_cart"`
	PaymentID     string                `json:"payment_id"`
	Metadata      WebhookMetadata       `json:"metadata"`
	TotalAmount   int64                 `json:"total_amount"`
	SettledAmount int64                 `json:"settled_amount"`
	Currency      dodopayments.Currency `json:"currency"`
	Status        string                `json:"status"`
}