package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type ConsumerService interface {
	GetConsumer() sarama.Consumer
}

// StartEventConsumer starts a simple Kafka consumer that listen to given topic
func StartEventConsumer(ctx context.Context, service ConsumerService, topic string, OnEvent func(p int32, pc sarama.PartitionConsumer)) error {
	partitions, err := service.GetConsumer().Partitions(topic)
	if err != nil {
		return err
	}

	log.Printf("payment kafka consumer starting; topic=%s partitions=%v", topic, partitions)

	done := make(chan struct{})
	for _, partition := range partitions {
		pc, err := service.GetConsumer().ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Printf("error starting partition consumer p=%d: %v", partition, err)
			continue
		}
		go OnEvent(partition, pc)
	}
	<-ctx.Done()
	close(done)

	return nil
}

func CloseConsumer(service ConsumerService) {
	if err := service.GetConsumer().Close(); err != nil {
		log.Printf("failed to close consumer: %v\n", err)
	} else {
		done <- true
	}
}
