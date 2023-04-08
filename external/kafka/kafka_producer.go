package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Producer struct {
	Producer *kafka.Producer
}

func (k *Producer) Produce(topic string, msg []byte) error {
	return k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
}
