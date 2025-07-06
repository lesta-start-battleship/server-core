package kafka

type KafkaProducer interface {
	Send(topic string, message any) error
	Close() error
	CreateTopic(topic string, numPartitions int32, replicationFactor int16) error
}
