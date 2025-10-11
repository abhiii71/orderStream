package internal

import (
	"context"
	"fmt"
	"log"
	"net"

	mapset "github.com/deckarep/golang-set/v2"

	account "github.com/abhiii71/orderStream/account/client"
	"github.com/abhiii71/orderStream/order/models"
	"github.com/abhiii71/orderStream/order/proto/pb"
	product "github.com/abhiii71/orderStream/product/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	productClient *product.Client
}

func ListenGRPC(service Service, accountURL string, productURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	productClient, err := product.NewClient(accountURL)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		productClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		pb.UnimplementedOrderServiceServer{},
		service,
		accountClient,
		productClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {

	_, err := s.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Println("error getting account", err)
		return nil, err
	}

	var productIDs []string
	for _, p := range request.Products {
		productIDs = append(productIDs, p.Id)
	}

	orderedProducts, err := s.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting ordered products", err)
		return nil, err
	}

	var products []*models.OrderedProduct
	totalPrice := 0.0

	for _, p := range orderedProducts {
		productObj := &models.OrderedProduct{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}

		for _, requestProduct := range request.Products {
			if requestProduct.Id == p.Id {
				productObj.Quantity = requestProduct.Quantity
				break
			}
		}
		if productObj.Quantity != 0 {
			products = append(products, productObj)
			totalPrice += productObj.Price * float64(productObj.Price)
		}
	}

	postOrder, err := s.service.PostOrder(ctx, request.AccountId, totalPrice, products)
	if err != nil {
		log.Println("error  posting postOrder", err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         uint64(postOrder.ID),
		TotalPrice: postOrder.TotalPrice,
		Products:   []*pb.ProductInfo{},
	}

	orderProto.CreatedAt, _ = postOrder.CreatedAt.MarshalBinary()
	for _, p := range postOrder.Products {
		orderProto.Products = append(orderProto.Products, &pb.ProductInfo{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &pb.PostOrderResponse{Order: orderProto}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, request *wrapperspb.UInt64Value) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, request.Value)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Taking unique products. We use set to avoid repeating
	productIdsSet := mapset.NewSet[string]()
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIdsSet.Add(p.ID)
		}
	}

	productIds := productIdsSet.ToSlice()

	products, err := s.productClient.GetProducts(ctx, 0, 0, productIds, "")
	if err != nil {
		log.Println("error getting account products: ", err)
		return nil, err
	}

	// collecting orders
	var orders []*pb.Order
	for _, order := range accountOrders {
		// encounter order
		encounterOrder := &pb.Order{
			AccountId:  order.AccountID,
			Id:         uint64(order.ID),
			TotalPrice: order.TotalPrice,
			Products:   []*pb.ProductInfo{},
		}

		encounterOrder.CreatedAt, _ = order.CreatedAt.MarshalBinary()

		// Decorate orders with products
		for _, orderedProduct := range order.Products {
			// Populate product fields
			for _, prod := range products {
				if prod.Id == orderedProduct.ID {
					orderedProduct.Name = prod.Name
					orderedProduct.Description = prod.Description
					orderedProduct.Price = prod.Price
					break
				}
			}

			encounterOrder.Products = append(encounterOrder.Products, &pb.ProductInfo{
				Id:          orderedProduct.ID,
				Name:        orderedProduct.Name,
				Description: orderedProduct.Description,
				Price:       orderedProduct.Price,
				Quantity:    orderedProduct.Quantity,
			})
		}
		orders = append(orders, encounterOrder)
	}

	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}

func (s *grpcServer) UpdateOrderStatus(ctx context.Context, request *pb.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateOrderPaymentStatus(ctx, request.OrderId, request.Status)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
