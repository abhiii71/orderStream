package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/abhiii71/orderStream/account/config"
	"github.com/abhiii71/orderStream/account/internal"
	"github.com/tinrab/retry"
	"gorm.io/driver/postgres"
)

func main() {
	var repo internal.AccountRepository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		db, err := sql.Open(postgres.Open(config.DatabaseURL), &sql.Config{})
		if err != nil {
			log.Fatal(err)
		}

		repo := internal.Newrepo(db)

		return
	})

	defer repo.Close()
	log.Println("Listening on port 8080...")
	service := internal.Newservice(repo)
	log.Fatal(internal.ListenGRPC(service, 8080))
}
