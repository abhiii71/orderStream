package config

import (
	"context"
	"log"

	"github.com/abhiii71/orderStream/payment/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.PaymentServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	C := pb.NewPaymentServiceClient(conn)
	return &Client{conn, C}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) CreateCustomerPortalSession(ctx context.Context, userId uint64, name, email string) (string, error) {
	res, err := c.service.CreateCustomerPortalSession(ctx, &pb.CustomerPortalRequest{
		UserId: userId,
		Name:   &name,
		Email:  &email,
	})
	if err != nil {
		return "", err
	}

	return res.Value, nil
}

func (c *Client) CreateCheckoutSession(ctx context.Context, orderId, userId int, name, email, redirectUrl string, products []*pb.CartItem) (string, error) {
	res, err := c.service.CreateCheckoutSession(ctx, &pb.CheckoutRequest{
		UserId:      uint64(userId),
		Name:        name,
		Email:       email,
		RedirectURL: redirectUrl,
		Products:    products,
		OrderId:     uint64(orderId),
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	return res.Value, nil
}
