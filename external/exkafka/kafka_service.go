package exkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"tax-management/external/exkafka/messages"
	"tax-management/pkg"
	terminal2 "tax-management/terminal"

	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaServiceImpl struct {
	Writer        *kafka.Writer
	Repository    pkg.ClientRepository
	Client        pkg.ClientLoggerExtension
	TaxClient     pkg.TaxClient
	Redis         pkg.RedisService
	Url           string
	TokenUrl      string
	ServerInfoUrl string
	Terminal      *terminal2.Terminal
}

func (kpi KafkaServiceImpl) Publish(msg messages.RawTransaction) error {
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

func (kpi KafkaServiceImpl) Consumer(message *messages.RawTransaction, consumerType string, err error) {

	fmt.Printf("receive message %s", consumerType)
	if err != nil {
		panic("receive message has error")
	}

	ctx := context.Background()

	id, err := kpi.Repository.InsertTaxData(ctx, consumerType, *message)
	if err != nil {
		fmt.Errorf("Insert raw taxData has error %s", err.Error())
	}
	fmt.Printf("id is %s", id)
	//token, err := kpi.TaxClient.GetToken()
	kpi.TaxClient.SendInvoice(*message)
	//	fmt.Printf("token is %s", token)

}
