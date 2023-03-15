package exkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"tax-management/external/exkafka/messages"
	"tax-management/pkg"
	"tax-management/utility"
	"time"

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

	serverInformation, err := kpi.TaxClient.GetServerInformation()
	if err != nil {
		fmt.Printf("GetServer information has error %s", err.Error())
		return
	}
	fmt.Printf("Server information is %s", *serverInformation)

	// normalize
	requestToNormalize := utility.SignatureFirstTypeRequest{
		Authorization:  "1234",
		RequestTraceId: "",
		TimeStamp:      time.Now().String(),
	}
	pft := utility.PacketFirstType{
		Uid:        "",
		PacketType: "invoice",
		Retry:      false,
		Data:       message.After,
	}
	requestToNormalize.Packets = append(requestToNormalize.Packets, pft)
	normalized, err := utility.Normalize(requestToNormalize)
	if err != nil {
		// update for retry has error in normalize
		// notif to developer
		fmt.Printf("normalize has error,%s", err.Error())
	}

	signature, err := utility.SignAndVerify(normalized)
	if err != nil {
		fmt.Printf("sign has error %s", err.Error())
		// update for retry has error in normalize
		// notif to developer
	}
	fmt.Printf("signature: %s", *signature)
	// token, err := kpi.Redis.Get(ctx, getTokenKey())
	// if err != nil {
	// 	tokenResp, err := kpi.get_token()
	// 	if err != nil {
	// 		fmt.Printf("Get token has error %s", err.Error())
	// 	} else {
	// 		if tokenResp != nil {
	// 			kpi.Redis.Set(ctx, getTokenKey(), *tokenResp, 3000)
	// 		}
	// 	}
	// }
	// fmt.Printf("just for using token %s", token)
}
