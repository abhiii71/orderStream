package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/abhiii71/orderStream/account"
	"github.com/abhiii71/orderStream/account/internal"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

func main() {
	err := godotenv.Load(".env") // relative to project root
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	log.Println("DATABASE_URL:", dbURL)

	var repository internal.AccountRepository

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
		repository = internal.NewAccountRepository(db)
		return nil
	})

	defer repository.Close()
	
	port := account.Port
	log.Printf("Listening on port %d...", port)

	service := internal.Newservice(repository)
	log.Fatal(internal.ListenGRPC(service, port))
}
