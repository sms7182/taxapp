package pkg

import "tax-management/external/exkafka/messages"

type KafkaService interface {
	Consumer(msg *messages.RawTransaction, err error)
	Publish(msg messages.RawTransaction) error
}
