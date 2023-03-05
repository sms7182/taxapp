package exkafka

import (
	"tax-management/pkg"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader  *kafka.Reader
	Service pkg.Service
}
