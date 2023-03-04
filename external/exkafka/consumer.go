package exkafka

import (
	"tax-app/pkg"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader  *kafka.Reader
	Service pkg.Service
}
