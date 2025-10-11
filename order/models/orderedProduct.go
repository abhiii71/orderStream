package models

type ProductInfo struct {
	Id        uint
	OrderID   uint
	ProductId string
	Quantity  int
}

// func (ProductInfo) TableName() string {
// 	return "order_products"
// }
