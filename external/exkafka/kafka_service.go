package exkafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	Redis         pkg.RedisService
	Url           string
	TokenUrl      string
	ServerInfoUrl string
}

func getTokenKey() string {
	return "tax-token"
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

func (kpi KafkaServiceImpl) Consumer(message *messages.RawTransaction, err error) {
	fmt.Printf("receive message")
	if err != nil {
		panic("receive message has error")
	}

	ctx := context.Background()

	id, err := kpi.Repository.InsertTaxData(ctx, *message)
	if err != nil {
		fmt.Errorf("Insert raw taxData has error %s", err.Error())
	}
	fmt.Println("id is %s", id)

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

func (ksi KafkaServiceImpl) get_token() (*string, error) {

	url := fmt.Sprintf(ksi.Url, ksi.TokenUrl)

	rqId, _ := uuid.NewV4()
	timeNow := time.Now().Unix()
	tstr := strconv.FormatInt(timeNow, 10)
	sPacketReq := utility.SignaturePacketRequest{

		RequestTraceId: rqId.String(),
		TimeStamp:      tstr,
		ContentType:    "application/json",
		Packet: utility.Packet{
			Uid:        rqId.String(),
			PacketType: "GET_TOKEN",
			Retry:      false,
			Data: utility.TokenBody{
				UserName: "A118GE",
			},
		},
	}

	normalized, err := utility.Normalize(sPacketReq)
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
	postRequest := utility.PostDataRequest{
		RequestTraceId: rqId.String(),
		TimeStamp:      tstr,
		ContentType:    "application/json",
		Packet: utility.Packet{
			Uid:        rqId.String(),
			PacketType: "GET_TOKEN",
			Retry:      false,
			Data: utility.TokenBody{
				UserName: "A118GE",
			},
		},
		Signature: signature,
	}
	jsonBytes, err := json.Marshal(postRequest)
	if err != nil {
		fmt.Printf("json marshal has error %s", err.Error())
		return nil, err
	}
	reader := bytes.NewReader(jsonBytes)

	request, err := http.NewRequest("POST", url, reader)

	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return nil, err
	}
	traceId, _ := uuid.NewV4()
	request.Header.Set("requestTraceId", traceId.String())
	request.Header.Set("timestamp", tstr)
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	//client := ksi.Client
	//resp, err := client.Do(postRequest.RequestTraceId, *signature, postRequest.Packet.PacketType, request, url)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response has error %s", err.Error())
		return nil, err
	}
	var tokenResponse utility.TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Printf("responseJson has error %s", err.Error())
		return nil, err
	}
	return &tokenResponse.Token, nil
}

func (ksi KafkaServiceImpl) getServerList() string {
	url := fmt.Sprintf(ksi.Url, ksi.ServerInfoUrl)
	id, _ := uuid.NewV4()

	bodyReq := utility.BodyReq{
		Time: 2,
		Packet: utility.Packet{
			Uid:             id.String(),
			PacketType:      "GET_SERVER_INFORMATION",
			Retry:           false,
			Data:            utility.TokenBody{},
			EncryptionKeyId: "",
			SymmetricKey:    "",
			IV:              "",
			FiscalId:        "",
			DataSignature:   "",
		},
	}
	marshaled, err := json.Marshal(bodyReq)
	if err != nil {
		fmt.Printf("has error create request")
	}

	jsonBytes := bytes.NewReader(marshaled)
	request, err := http.NewRequest("POST", url, jsonBytes)
	if err != nil {
		fmt.Printf("Create Post request has error %s", err.Error())
	}
	request.Header.Set("requestTraceId", id.String())
	request.Header.Set("timestampt", time.Now().String())
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("response has error %s", err.Error())

	}
	body, err := ioutil.ReadAll(resp.Body)
	var bodyObj utility.BodyResponse

	if err != nil {
		fmt.Printf("read response has error %s", err.Error())

	}
	json.Unmarshal(body, &bodyObj)
	return bodyObj.Result.Data.PublicKeys[0].Key
}
