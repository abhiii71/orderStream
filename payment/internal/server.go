package internal

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)


func StartServers(service PaymentService, consumer sarama.Consumer, orderURL string, grpcPort, webhookPort int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 3)

	// Start kafka consumer if available 
	if consumer !=  nil {
		wg.Add(1)
		go func(){
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
go func(){
	defer wg.Done()
	if err := ListneGrpc
}()
}