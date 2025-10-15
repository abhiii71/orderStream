package graph

import (
	"context"
	"log"
	"time"

	"github.com/abhiii71/orderStream/graphql/generated"
	"github.com/abhiii71/orderStream/graphql/models"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) ID(ctx context.Context, obj *models.Account) (int, error) {
	return int(obj.ID), nil
}

func (r *accountResolver) Orders(ctx context.Context, obj *models.Account) ([]*generated.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetordersForAccount(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*generated.Order
	for _, order := range orderList {
		var products []*generated.OrderedProduct
		for _, orderProduct := range order.Products {
			products = append(products, &generated.OrderedProduct{
				ID:          orderProduct.ID,
				Name:        orderProduct.Name,
				Description: orderProduct.Description,
				Price:       order.TotalPrice,
				Quantity:    int(orderProduct.Quantity),
			})
		}

		orders = append(orders, &generated.Order{
			ID:         int(order.ID),
			CreatedAt:  order.CreatedAt,
			TotalPrice: order.TotalPrice,
			Products:   products,
		})
	}

	return orders, nil
}
