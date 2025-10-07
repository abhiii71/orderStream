package internal

import (
	"context"
	"fmt"
	"net"

	"github.com/abhiii71/orderStream/account/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service AccountService
}

func ListenGRPC(service AccountService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()

	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		service:                           service})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) Register(ctx context.Context, request *pb.RegisterRequest) (*wrapperspb.StringValue, error) {
	token, err := s.service.Register(ctx, request.Name, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{
		Value: token,
	}, nil
}
