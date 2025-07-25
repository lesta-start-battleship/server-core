package event

import (
	"github.com/lesta-battleship/server-core/internal/config"
	"github.com/lesta-battleship/server-core/internal/infra/kafka"
)

type MatchEventPublisher interface {
	PublishMatchResult(result MatchResult) error
	PublishUsedItem(item Item) error
}

type KafkaMatchEventPublisher struct {
	kafkaProducer kafka.KafkaProducer
}

func NewKafkaMatchEventPublisher(producer kafka.KafkaProducer) *KafkaMatchEventPublisher {
	return &KafkaMatchEventPublisher{kafkaProducer: producer}
}

func (p *KafkaMatchEventPublisher) PublishMatchResult(result MatchResult) error {
	return p.kafkaProducer.Send(config.MatchResults, result)
}

func (p *KafkaMatchEventPublisher) PublishUsedItem(item Item) error {
	return p.kafkaProducer.Send(config.UsedItems, item)
}
