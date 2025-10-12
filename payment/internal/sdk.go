package internal

import (
	"context"
	"net/http"

	"github.com/abhiii71/orderStream/payment/models"
	"github.com/dodopayments/dodopayments-go"
	"github.com/dodopayments/dodopayments-go/option"
)

type PaymentClient interface {
	CreateProduct(ctx context.Context, name string, price int64, currency dodopayments.Currency, taxCategory dodopayments.TaxCategory, customerId, productId string) (*dodopayments.Product, error)
	UpdateProduct(ctx context.Context, productId, name string, price int64) error
	ArchiveProduct(ctx context.Context, productId string) error
	CreateCutomer(ctx context.Context, userId int64, name, email string) (*models.Customer, error)
	CreateCustomerSession(ctx context.Context, customerId string) (string, error)
	CreateCheckoutSession(ctx context.Context, userId int64, customerId string, redirect string, dodoProducts []dodopayments.CheckoutSessionRequestProductCartParam, orderId uint64) (checkoutURL string, err error)
	HandleWebhook(w http.ResponseWriter, r *http.Request) (*models.Transaction, error)
}

func NewDodoClient(apiKey string, testMode bool) PaymentClient {
	if testMode {
		return &dodoClient{
			client: dodopayments.NewClient(
				option.WithBearerToken(apiKey),
				option.WithEnvironmentTestMode(),
			),
		}
	}
	return &dodoClient{
		client: dodopayments.NewClient(
			option.WithBearerToken(apiKey),
		),
	}
}

type dodoClient struct {
	client *dodopayments.Client
}


func(d *dodoClient)	CreateProduct(ctx context.Context, name string, price int64, currency dodopayments.Currency, taxCategory dodopayments.TaxCategory, customerId, productId string) (*dodopayments.Product, error){
	product, err := d.client.Products.New(ctx , dodopayments.ProductNewParams{
		Name: dodopayments.F(name),
		Price: dodopayments.F[dodopayments.PriceUnionParam](
			dodopayments.PriceOneTimePriceParam{
			Price: dodopayments.F(price),
			Currency: dodopayments.F(currency),
			Discount: dodopayments.F[int64](0),
		},
	),
	TaxCategory: dodopayments.F(taxCategory),
})
if err != nil {
	return nil,err 
}
return prodct, nil
}

func(d *dodoClient)UpdateProduct(ctx context.Context, productId, name string, price int64) error{

}
func(d *dodoClient)	ArchiveProduct(ctx context.Context, productId string) error {

}

func(d *dodoClient)	CreateCutomer(ctx context.Context, userId int64, name, email string) (*models.Customer, error){

}
	func(d *dodoClient)CreateCustomerSession(ctx context.Context, customerId string) (string, error)
	{
	}
	func(d *dodoClient)CreateCheckoutSession(ctx context.Context, userId int64, customerId string, redirect string, dodoProducts []dodopayments.CheckoutSessionRequestProductCartParam, orderId uint64) (checkoutURL string, err error)
	
	{}
	func(d *dodoClient)HandleWebhook(w http.ResponseWriter, r *http.Request) (*models.Transaction, error){}