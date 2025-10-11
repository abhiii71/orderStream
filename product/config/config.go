package config

import "os"

var (
	ElasticsearchURL string
	BootstrapServers string
)

func init() {
	ElasticsearchURL = os.Getenv("ELASTICSEARCH_URL")
	BootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
}
