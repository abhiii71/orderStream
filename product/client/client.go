package client

import (
	"context"
	"log"

	"github.com/abhiii71/orderStream/product/models"
	"github.com/abhiii71/orderStream/product/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewProductServiceClient(conn)
	return &Client{conn, client}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	res, err := c.service.GetProduct(ctx, &wrapperspb.StringValue{Value: id})
	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         res.Product.Id,
		Name:       res.Product.Name,
		Decription: res.Product.Description,
		Price:      res.Product.Price,
		AccountId:  int(res.Product.GetAccountId()),
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, ids []string, query string) ([]models.Product, error) {
	res, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	var products []models.Product
	for _, p := range res.Products {
		products = append(products, models.Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			AccountID:   int(p.AccountId),
		})
	}
	return products, nil
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64, acccountId int64) (*models.Product, error) {
	res, err := c.service.PostProduct(ctx, &pb.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
		AccountId:   acccountId,
	})
	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         res.Product.Id,
		Name:       res.Product.Name,
		Decription: res.Product.Description,
		Price:      res.Product.Price,
		AccountId:  int(res.Product.GetAccountId()),
	}, nil
}

func (c *Client) UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId int64) (*models.Product, error) {
	res, err := c.service.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		AccountId:   accountId,
	})
	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         res.Product.Id,
		Name:       res.Product.Name,
		Decription: res.Product.Description,
		Price:      res.Product.Price,
		AccountId:  int(res.Product.GetAccountId()),
	}, nil
}

func (c *Client) DeleteProduct(ctx context.Context, productId string, accountId int64) error {
	_, err := c.service.DeleteProduct(ctx, &pb.DeleteProductRequest{ProductId: productId, AccountId: accountId})
	return err
}
