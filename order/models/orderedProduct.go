package models

type ProductInfo struct {
	Id        uint64
	OrderID   uint64
	ProductId string
	Quantity  int
}

// func (ProductInfo) TableName() string {
// 	return "order_products"
// }
