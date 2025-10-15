package graph

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/abhiii71/orderStream/graphql/generated"
	"github.com/abhiii71/orderStream/graphql/models"
	payment "github.com/abhiii71/orderStream/payment/proto/pb"
	"github.com/abhiii71/orderStream/pkg/auth"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) Register(ctx context.Context, in generated.RegisterInput) (*generated.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := r.server.accountClient.Register(ctx, in.Name, in.Email, in.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ginContext, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, errors.New("could not retrieve gin context")
	}
	ginContext.SetCookie("token", token, 3600, "/", "localhost", false, true)

	return &generated.AuthResponse{Token: token}, nil
}

func (r *mutationResolver) Login(ctx context.Context, in generated.LoginInput) (*generated.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := r.server.accountClient.Login(ctx, in.Email, in.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ginContext, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, errors.New("could not retrieve gin context")
	}
	ginContext.SetCookie("token", token, 3600, "/", "localhost", false, true)

	return &generated.AuthResponse{Token: token}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in generated.CreateProductInput) (*generated.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log.Println("CreateProduct called with input: ", in)

	accountId, err := auth.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, err
	}
	log.Println("CreateProduct called with accountId: ", accountId)
	postProduct, err := r.server.productClient.PostProduct(ctx, in.Name, in.Description, in.Price, int64(accountId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Created product: ", postProduct)
	log.Println("Product Id: ", postProduct.Id)

	return &generated.Product{
		ID:          postProduct.Id,
		Name:        postProduct.Name,
		Description: postProduct.Description,
		Price:       postProduct.Price,
		AccountID:   accountId,
	}, nil
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, in generated.UpdateProductInput) (*generated.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	accountId, err := auth.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := r.server.productClient.UpdateProduct(ctx, in.ID, in.Name, in.Description, in.Price, int64(accountId))
	if err != nil {
		return nil, err
	}

	return &generated.Product{
		ID:          updatedProduct.Id,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		AccountID:   accountId,
	}, nil
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (*bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	accountId, err := auth.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, err
	}

	err = r.server.productClient.DeleteProduct(ctx, id, int64(accountId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	success := true
	return &success, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in generated.OrderInput) (*generated.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []*models.OrderedProducts
	for _, product := range in.Products {
		if product.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}

		products = append(products, &models.OrderedProducts{
			ID:       product.ID,
			Quantity: product.Quantity,
		})
	}

	accountId, err := auth.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	postOrder, err := r.server.orderClient.PostOrder(ctx, uint64(accountId), products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orderedProducts []*generated.OrderedProduct
	for _, orderProduct := range postOrder.Products {
		orderedProducts = append(orderedProducts, &generated.OrderedProduct{
			ID:          orderProduct.ID,
			Name:        orderProduct.Name,
			Description: orderProduct.Description,
			Price:       orderProduct.Price,
			Quantity:    int(orderProduct.Quantity),
		})
	}

	return &generated.Order{
		ID:         int(postOrder.ID),
		CreatedAt:  postOrder.CreatedAt,
		TotalPrice: postOrder.TotalPrice,
		Products:   orderedProducts,
	}, nil
}

func (r *mutationResolver) CreateCustomerPortalSession(ctx context.Context, credentials *generated.CustomerPortalSessionInput) (*generated.RedirectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	UrlWithSession, err := r.server.paymentClient.CreateCustomerPortalSession(ctx, uint64(credentials.AccounntID), credentials.Name, credentials.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &generated.RedirectResponse{URL: UrlWithSession}, nil
}

func (r *mutationResolver) CreateCheckoutSession(ctx context.Context, details *generated.CheckoutInput) (*generated.RedirectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var products []*payment.CartItem
	for _, product := range details.Products {
		products = append(products, &payment.CartItem{
			ProductId: product.ID,
			Quantity:  uint64(product.Quantity),
		})
	}

	UrlWithCheckoutSession, err := r.server.paymentClient.CreateCheckoutSession(ctx, details.OrderID, details.AccounID, details.Name, details.Email,
		details.RedirectURL, products)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &generated.RedirectResponse{URL: UrlWithCheckoutSession}, nil
}
