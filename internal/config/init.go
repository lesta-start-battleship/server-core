package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Port               string
	KafkaBrokers       []string
	TopicsToSend       []string
	MatchResults       string
	UsedItems          string
	GetAllItemsURL     string
	GetAllUserItemsURl string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Port = os.Getenv("GAME_CORE_PORT")
	if Port == "" {
		log.Fatal("Game core port not set")
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

	GetAllItemsURL = os.Getenv("INVENTORY_SERVICE_GET_ALL_ITEMS")
	if GetAllItemsURL == "" {
		log.Fatal("INVENTORY_SERVICE_URL not set")
	}

	GetAllUserItemsURl = os.Getenv("INVENTORY_SERVICE_GET_USER_ITEMS")
	if GetAllUserItemsURl == "" {
		log.Fatal("INVENTORY_SERVICE_URL not set")
	}

	TopicsToSend = []string{MatchResults, UsedItems}
}
