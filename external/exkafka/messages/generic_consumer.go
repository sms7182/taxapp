package messages

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Consumer[T comparable] struct {
	Type    string
	reader  *kafka.Reader
	Dialer  *kafka.Dialer
	Topic   string
	Brokers []string
}

func (c *Consumer[T]) CreateConnection() {

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.Brokers,
		Topic:   c.Topic,

		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		GroupID:  "tax-management",
	})
	c.reader.SetOffset(0)

}

func (c *Consumer[T]) Read(model T, callback func(T, error)) {
	for {
		ctx := context.Background()
		message, err := c.reader.ReadMessage(ctx)

		if err != nil {
			callback(model, err)
			return
		}
		err = json.Unmarshal(message.Value, &model)
		if err != nil {
			callback(model, err)
			continue
		}
		callback(model, nil)
	}
}
