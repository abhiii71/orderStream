package internal

import (
	"context"
	"net/http"

	"github.com/abhiii71/orderStream/payment/models"
	"github.com/abhiii71/orderStream/payment/proto/pb"
)

type PaymentService interface {
	RegisterProduct(ctx context.Context, name string, price int64, customerId, productId string) error
	UpdateProduct(ctx context.Context, productId string, name string, price int64) error
	DeleteProduct(ctx context.Context, productId string) error
	CreateCustomerPortalSession(ctx context.Context, customer *models.Customer) (string, error)
	FindOrCreateCustomer(ctx context.Context, userId uint64, name, email string) (*models.Customer, error)
	CreateCheckoutSession(ctx context.Context, userId uint64, customerId, redirect string, products []*pb.CartItem, orderId uint64) (checkoutURL string, err error)
	HandlePaymentWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.Transaction, error)
}

type paymentService struct {
	client            PaymentClient
	paymentRepository PaymentRepository
}
