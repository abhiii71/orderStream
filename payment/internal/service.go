package internal

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/abhiii71/orderStream/payment/models"
	"github.com/abhiii71/orderStream/payment/proto/pb"
	"github.com/dodopayments/dodopayments-go"
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

func NewPaymentService(client PaymentClient, paymentRepository PaymentRepository) PaymentService {
	return &paymentService{client: client, paymentRepository: paymentRepository}
}

func (ds *paymentService) RegisterProduct(ctx context.Context, name string, price int64, customerId, productId string) error {

	// we will use USD  as currency and Digital products as tax category for now to keep it simple
	product, err := ds.client.CreateProduct(ctx, name, price, dodopayments.CurrencyUsd, dodopayments.TaxCategoryDigitalProducts, customerId, productId)
	if err != nil {
		return err
	}
	return ds.paymentRepository.SaveProduct(ctx, &models.Product{
		ProductID:     productId,
		DodoProductID: product.ProductID,
		Price:         product.Price.FixedPrice,
		Currency:      string(product.Price.Currency),
	})
}

func (ds *paymentService) UpdateProduct(ctx context.Context, productId string, name string, price int64) error {

	err := ds.client.UpdateProduct(ctx, productId, name, price)
	if err != nil {
		return err
	}

	product, err := ds.paymentRepository.GetProductByProductId(ctx, productId)
	if err != nil {
		return err
	}

	if product.Price != price {
		product.Price = price
		err = ds.paymentRepository.UpdateProduct(ctx, product)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *paymentService) DeleteProduct(ctx context.Context, productId string) error {
	err := ds.client.ArchiveProduct(ctx, productId)
	if err != nil {
		return err
	}

	return ds.paymentRepository.DeleteProduct(ctx, productId)
}

func (ds *paymentService) CreateCustomerPortalSession(ctx context.Context, customer *models.Customer) (string, error) {
	customerPortalLink, err := ds.client.CreateCustomerSession(ctx, customer.CustomerId)
	if err != nil {
		return "", err
	}

	return customerPortalLink, nil
}

func (ds *paymentService) FindOrCreateCustomer(ctx context.Context, userId uint64, name, email string) (*models.Customer, error) {
	existingCustomer, err := ds.paymentRepository.GetCustomerByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	if err == nil {
		return existingCustomer, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	customer, err := ds.client.CreateCustomer(ctx, int64(userId), name, email)
	if err != nil {
		return nil, err
	}

	err = ds.paymentRepository.SaveCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

// createcheckoutsession - returns url to check out page  and error
func (ds *paymentService) CreateCheckoutSession(ctx context.Context, userId uint64, customerId, redirect string, products []*pb.CartItem, orderId uint64) (checkoutURL string, err error) {
	productIds := make([]string, len(products))
	productQuantities := make(map[string]uint64, len(products))

	for i, product := range products {
		productIds[i] = product.ProductId
		productQuantities[product.ProductId] = product.Quantity
	}

	modelsProducts, err := ds.paymentRepository.GetProductsByIds(ctx, productIds)
	if err != nil {
		return "", err
	}

	var dodoProducts []dodopayments.CheckoutSessionRequestProductCartParam

	for _, product := range modelsProducts {
		dodoProducts = append(dodoProducts, dodopayments.CheckoutSessionRequestProductCartParam{
			ProductID: dodopayments.F(product.DodoProductID),
			Quantity:  dodopayments.F(int64(productQuantities[product.ProductID])),
		})

	}

	return ds.client.CreateCheckoutSession(ctx, int64(userId), customerId, redirect, dodoProducts, orderId)
}

func (ds *paymentService) HandlePaymentWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.Transaction, error) {
	updatedTransaction, err := ds.client.HandleWebhook(w, r)
	if err != nil {
		return nil, err
	}

	err = ds.paymentRepository.UpdatedTransaction(ctx, updatedTransaction)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, err
}
