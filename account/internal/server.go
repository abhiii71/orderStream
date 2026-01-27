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

func (s *grpcServer) Login(ctx context.Context, request *pb.LoginRequest) (*wrapperspb.StringValue, error) {
	token, err := s.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{
		Value: token,
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, request *wrapperspb.UInt64Value) (*pb.AccountResponse, error) {
	id := request.GetValue()

	account, err := s.service.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.AccountResponse{Account: &pb.Account{
		Id:    account.ID,
		Name:  account.Name,
		Email: account.Email,
	}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	getAccounts, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var accounts []*pb.Account
	for _, getAccount := range getAccounts {
		accounts = append(accounts, &pb.Account{
			Id:    uint64(int(getAccount.ID)),
			Name:  getAccount.Name,
			Email: getAccount.Email,
		},
		)
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
