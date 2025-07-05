package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	KafkaBrokers []string
	TopicsToSend []string
	MatchResults string
	UsedItems    string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS not set")
	}
	KafkaBrokers = strings.Split(kafkaBrokers, ",")

	MatchResults = os.Getenv("MATCH_RESULTS")
	if MatchResults == "" {
		log.Fatal("MATCH_RESULTS not set")
	}

	UsedItems = os.Getenv("USED_ITEMS")
	if UsedItems == "" {
		log.Fatal("USED_ITEMS not set")
	}

	TopicsToSend = []string{MatchResults, UsedItems}
}
