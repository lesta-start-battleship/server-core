package main

import (
	"lesta-battleship/server-core/internal/game-core/api"
	"lesta-battleship/server-core/internal/game-core/config"
	"lesta-battleship/server-core/internal/game-core/event"
	"lesta-battleship/server-core/internal/game-core/infra/kafka"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	producer, err := kafka.NewProducer(config.KafkaBrokers, config.TopicsToSend)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	publisher := event.NewKafkaMatchEventPublisher(producer)
	dispatcher := event.NewMatchEventDispatcher(publisher)

	router := gin.Default()

	api.SetupRoutes(router, dispatcher)

	router.Run(":8080")
}
