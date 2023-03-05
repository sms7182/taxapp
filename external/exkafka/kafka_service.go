package exkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"tax-management/pkg"

	"github.com/segmentio/kafka-go"
)

type KafkaServiceImpl struct {
	Reader *kafka.Reader

	Service pkg.Service
}

func (kpi KafkaServiceImpl) Consumer(id string, err error) {
	fmt.Printf("id from message %s", id)
}
func (kpi KafkaServiceImpl) Read(id string, callback func(string, error)) {
	for {

		ctx := context.Background()
		message, err := kpi.Reader.ReadMessage(ctx)

		if err != nil {
			callback(id, err)
			return
		}

		err = json.Unmarshal(message.Value, &id)
		if err != nil {
			callback(id, err)
			continue
		}
		err = kpi.Reader.CommitMessages(ctx, message)
		if err != nil {
			callback(id, err)
			continue
		}
		callback(id, nil)
	}
}
