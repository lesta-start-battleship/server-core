package main

import (
	"github.com/lesta-battleship/server-core/internal/api"
	"github.com/lesta-battleship/server-core/internal/config"
	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/infra/kafka"
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

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	api.SetupRoutes(router, dispatcher)

	log.Println("Listening and serving HTTP on:", config.Port)
	if err := router.Run(":" + config.Port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
