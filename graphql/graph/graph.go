package graph

import (
	"github.com/99designs/gqlgen/graphql"
	account "github.com/abhiii71/orderStream/account/client"
	"github.com/abhiii71/orderStream/graphql/generated"
	order "github.com/abhiii71/orderStream/order/client"
	payment "github.com/abhiii71/orderStream/payment/client"
	product "github.com/abhiii71/orderStream/product/client"
	recommender "github.com/abhiii71/orderStream/recommender/client"
)

type Server struct {
	accountClient     *account.Client
	productClient     *product.Client
	orderClient       *order.Client
	paymentClient     *payment.Client
	recommenderClient *recommender.Client
}

func NewGraphQLServer(accountUrl, productUrl, orderUrl, paymentUrl, recommenderUrl string) (*Server, error) {
	accClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	prodClient, err := product.NewClient(productUrl)
	if err != nil {
		accClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accClient.Close()
		prodClient.Close()
		return nil, err
	}

	paymentClient, err := payment.NewClient(paymentUrl)
	if err != nil {
		accClient.Close()
		prodClient.Close()
		orderClient.Close()
		return nil, err
	}

	recClient, err := recommender.NewClient(productUrl)
	if err != nil {
		accClient.Close()
		prodClient.Close()
		orderClient.Close()
		paymentClient.Close()
		return nil, err
	}

	return &Server{
		accountClient:     accClient,
		productClient:     prodClient,
		orderClient:       orderClient,
		paymentClient:     paymentClient,
		recommenderClient: recClient,
	}, nil
}

func (s *Server) Mutation() generated.MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() generated.QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() generated.AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}
