package exkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"tax-management/pkg"

	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaServiceImpl struct {
	Reader     *kafka.Reader
	Writer     *kafka.Writer
	Repository pkg.ClientRepository
}

func (kpi KafkaServiceImpl) Publish(msg string) error {
	key, _ := uuid.NewV4()

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = kpi.Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte((key).String()),
		Value: bytes,
	})
	return err
}

func (kpi KafkaServiceImpl) Consumer(id string, err error) {
	fmt.Printf("id from message %s", id)
}
func (kpi KafkaServiceImpl) Read(id string, callback func(string, error)) {
	for {

		ctx := context.Background()
		message, err := kpi.Reader.FetchMessage(ctx)

		if err != nil {
			callback(id, err)
			return
		}

		err = json.Unmarshal(message.Value, &id)
		if err != nil {
			callback(id, err)
			continue
		}
		if err := kpi.Reader.CommitMessages(ctx, message); err != nil {
			callback(id, err)
			continue
		}
		callback(id, nil)
	}
}
