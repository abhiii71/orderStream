package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/abhiii71/orderStream/product/config"
	"github.com/abhiii71/orderStream/product/internal"
	"github.com/tinrab/retry"
)

func main() {
	var repo internal.Repository

	producer, err := sarama.NewAsyncProducer([]string{config.BootstrapServers}, nil)
	if err != nil {
		log.Println(err)
	}
	defer func(producer sarama.AsyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Println(err)
		}
	}(producer)

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repo, err = internal.NewElasticRepository(config.ElasticsearchURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repo.Close()

	log.Println("listening on port 8080...")
	service := internal.NewProductService(repo, producer)
	log.Fatal(internal.ListenGRPC(service, 8080))
}
