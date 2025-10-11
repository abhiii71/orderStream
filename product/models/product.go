package models

type Product struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Description string  `json:"description"`
	Price      float64 `json:"price"`
	AccountId  int     `json:"accountId"`
}

type ProductDocument struct {
	Name       string  `json:"name"`
	Description string  `json:"description"`
	Price      float64 `json:"price"`
}
