package client

import (
	"context"
	"log"

	"github.com/abhiii71/orderStream/account/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	C := pb.NewAccountServiceClient(conn)
	return &Client{conn, C}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) Register(ctx context.Context, name, email, password string) (string, error) {
	response, err := c.service.Register(ctx, &pb.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return response.Value, nil

}
