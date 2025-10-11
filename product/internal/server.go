package internal

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/abhiii71/orderStream/product/models"
	"github.com/abhiii71/orderStream/product/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type grpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()

	pb.RegisterProductServiceServer(serv, &grpcServer{
		UnimplementedProductServiceServer: pb.UnimplementedProductServiceServer{},
		service:                           s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) GetProduct(ctx context.Context, request *wrapperspb.StringValue) (*pb.ProductResponse, error) {
	product, err := s.service.GetProduct(ctx, request.Value)
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{Product: &pb.Product{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, request *pb.GetProductsRequest) (*pb.ProductsResponse, error) {
	var res []models.Product
	var err error

	if request.Query != "" {
		res, err = s.service.SearchProducts(ctx, request.Query, request.Skip, request.Take)
	} else if len(request.Ids) != 0 {
		res, err = s.service.GetProductsWithIds(ctx, request.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, request.Skip, request.Take)
	}
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})

	}

	return &pb.ProductsResponse{Products: products}, nil
}

func (s *grpcServer) PostProduct(ctx context.Context, request *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product, err := s.service.PostProduct(ctx, request.Name, request.Description, request.Price, int(request.GetAccountId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ProductResponse{Product: &pb.Product{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, nil
}

func (s *grpcServer) UpdateProduct(ctx context.Context, request *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product, err := s.service.UpdateProduct(ctx, request.Id, request.Name, request.Description, request.Price, int(request.GetAccountId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ProductResponse{Product: &pb.Product{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, nil
}

func (s *grpcServer) DeleteProduct(ctx context.Context, request *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteProduct(ctx, request.GetProductId(), int(request.GetAccountId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
