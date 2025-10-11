package config

import "os"

var (
	DatabaseURL     string
	AccountURL      string
	ProductURL      string
	BootStrapServer string
)

func init() {
	DatabaseURL = os.Getenv("DATABASE_URL")
	AccountURL = os.Getenv("ACCOUNT_URL")
	ProductURL = os.Getenv("PRODUCT_URL")
	BootStrapServer = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
}
