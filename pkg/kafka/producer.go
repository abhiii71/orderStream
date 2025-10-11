package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

var done = make(chan bool)

type ProducerService interface {
	GetProducer() sarama.AsyncProducer
}

func SendMessageToRecommender(service ProducerService, event any, topic string) error {
	jsonMessage, err := json.Marshal(event)
	if err != nil {
		log.Println("failed to marshal event: ", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	// send the message asynchronously
	service.GetProducer().Input() <- msg

	return nil
}

func CloseProducer(service ProducerService) {
	if err := service.GetProducer().Close(); err != nil {
		log.Printf("failed to close producer: %v\n", err)
	} else {
		done <- true
	}
}
