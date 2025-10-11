package models

type EventData struct {
	AccountId int    `json:"user_id"`
	ProductId string `json:"product_id"`
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}
