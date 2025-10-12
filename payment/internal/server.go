package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/IBM/sarama"
	order "github.com/abhiii71/orderStream/order/client"
	"github.com/abhiii71/orderStream/payment/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServers(service PaymentService, consumer sarama.Consumer, orderURL string, grpcPort, webhookPort int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 3)

	// Start kafka consumer if available
	if consumer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eventConsumer := NewEventConsumer(consumer, service)
			ctx := context.Background()
			if err := eventConsumer.StartProductEventConsumer(ctx); err != nil {
				errCh <- fmt.Errorf("kafka consumer error: %w", err)
			}
		}()
	}

	// start gRPC Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ListenGRPC(service, orderURL, grpcPort); err != nil {
			errCh <- fmt.Errorf("grpc server error: %w", err)
		}
	}()

	// start webhook HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := listenWebhook(service, orderURL, webhookPort); err != nil {
			errCh <- fmt.Errorf("webhook server error: %w", err)
		}
	}()

	//wait for  first error or all server to complete
	go func() {
		wg.Wait()
		close(errCh)
	}()

	return <-errCh
}

func ListenGRPC(service PaymentService, orderURL string, port int) error {
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		orderClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterPaymentServiceServer(serv, &grpcServer{
		pb.UnimplementedPaymentServiceServer{},
		service,
		*orderClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}
func listenWebhook(service PaymentService, orderURL string, port int) error {
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		return err
	}

	defer orderClient.Close()

	webhookServer := &WebhookServer{
		service:     service,
		orderClient: orderClient,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/webhook/payment", webhookServer.HandlePaymentWebhook)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("webhook server listenin on port %d", port)
	return server.ListenAndServe()

}
