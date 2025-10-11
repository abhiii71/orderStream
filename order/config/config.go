package config

import "os"

var (
	DatabaseURL     string
	AccountURL      string
	ProductURL      string
	BootStrapServers string
)

func init() {
	DatabaseURL = os.Getenv("DATABASE_URL")
	AccountURL = os.Getenv("ACCOUNT_URL")
	ProductURL = os.Getenv("PRODUCT_URL")
	BootStrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
}
