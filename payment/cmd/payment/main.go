package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/abhiii71/orderStream/payment/config"
	"github.com/abhiii71/orderStream/payment/internal"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"


	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

func main() {
	err := godotenv.Load(".env") // relative to project root
	if err != nil {
		log.Println(".env file not found!")
	}

	var repository internal.PaymentRepository

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	log.Println("DATABASE_URL:", dbURL)

	// Retry connecting to DB
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		db, err := sql.Open("pgx", dbURL)
		if err != nil {
			log.Println("DB connection error:", err)
			return err
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Println("DB ping error:", err)
			return err
		}

		// Wrap db in PostgresAccountRepository
		repository = internal.NewPostgresRepository(db)
		return nil
	})

	defer repository.Close()

	// Setup Kafka consumer
	var consumer sarama.Consumer
	if config.KafkaBrokers != "" {
		retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
			kafkaConfig := sarama.NewConfig()
			kafkaConfig.Consumer.Return.Errors = true

			consumer, err = sarama.NewConsumer([]string{config.KafkaBrokers}, kafkaConfig)
			if err != nil {
				log.Printf("Failed to create Kafka consumer: %v", err)
			}
			return
		})
	}
	dodoClient := internal.NewDodoClient(config.DodoAPIKEY, config.DodoTestMode)
	service := internal.NewPaymentService(dodoClient, repository)

	log.Fatal(internal.StartServers(service, consumer, config.OrderServiceURL, config.GrpcPort, config.WebhookPort))
}
