package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/abhiii71/orderStream/account"
	"github.com/abhiii71/orderStream/order/config"
	"github.com/abhiii71/orderStream/order/internal"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

func main() {
	err := godotenv.Load(".env") // relative to project root
	if err != nil {
		log.Println(".env file not found!")
	}

	producer, err := sarama.NewAsyncProducer([]string{config.BootStrapServers}, nil)
	if err != nil {
		log.Println(err)
	}

	defer func(producer sarama.AsyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Println(err)
		}
	}(producer)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	log.Println("DATABASE_URL:", dbURL)

	var repository internal.OrderRepository

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

	port := account.Port
	log.Printf("Listening on port %d...", port)

	service := internal.NewOrderService(repository, producer)
	log.Fatal(internal.ListenGRPC(service, config.AccountURL, config.ProductURL, port))
}
