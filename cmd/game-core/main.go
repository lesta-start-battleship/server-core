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

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	api.SetupRoutes(router, dispatcher)

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
	log.Println("Listening and serving HTTP on: ", config.Port)
}
