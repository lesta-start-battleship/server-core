package kafka

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.AsyncProducer
	admin    sarama.ClusterAdmin
	brokers  []string

	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewProducer(brokers []string, topics []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		asyncProducer.Close()
		return nil, err
	}

	producer := &Producer{
		producer: asyncProducer,
		admin:    admin,
		brokers:  brokers,
		stopChan: make(chan struct{}),
	}

	for _, topic := range topics {
		if err := producer.CreateTopic(topic, 3, 1); err != nil {
			log.Printf("[KAFKA] Create topic failed: %v", err)
		}
	}

	producer.wg.Add(2)
	go producer.listenForSuccess()
	go producer.listenForErrors()

	return producer, nil
}

func (p *Producer) listenForSuccess() {
	defer p.wg.Done()
	for {
		select {
		case success, ok := <-p.producer.Successes():
			if !ok {
				return
			}
			log.Printf("[KAFKA] Message sent to topic %s (partition %d, offset %d)",
				success.Topic, success.Partition, success.Offset)
		case <-p.stopChan:
			return
		}
	}
}

func (p *Producer) listenForErrors() {
	defer p.wg.Done()
	for {
		select {
		case err, ok := <-p.producer.Errors():
			if !ok {
				return
			}
			log.Printf("[KAFKA] Failed to send message: %v", err)
		case <-p.stopChan:
			return
		}
	}
}

func (p *Producer) Send(topic string, message any) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("[KAFKA] Failed to marshal message: %v", err)
		return err
	}

	log.Printf("[KAFKA] Sending message to topic %s: %s", topic, string(msgBytes))

	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msgBytes),
	}

	return nil
}

func (p *Producer) Close() error {
	close(p.stopChan)

	if err := p.producer.Close(); err != nil {
		return err
	}

	p.wg.Wait()

	return p.admin.Close()
}

func (p *Producer) CreateTopic(topic string, numPartitions int32, replicationFactor int16) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err := p.admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		if terr, ok := err.(*sarama.TopicError); ok && terr.Err == sarama.ErrTopicAlreadyExists {
			log.Printf("[KAFKA] Topic %s already exists", topic)
			return nil
		}
		log.Printf("[KAFKA] Failed to create topic %s: %v", topic, err)
		return err
	}

	log.Printf("[KAFKA] Successfully created topic %s with %d partitions and replication factor %d",
		topic, numPartitions, replicationFactor)
	return nil
}
