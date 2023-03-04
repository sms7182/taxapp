package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"tax-app/external/exkafka"
	"tax-app/external/gateway"
	"tax-app/external/pg"
	"tax-app/utility"

	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	setUpViper()
}

func setUpViper() {
	viper.SetConfigName(getEnv("CONFIG_NAME", "dev-conf"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	db := getGormDb()
	repository := pg.RepositoryImpl{
		DB: db,
	}

	client := gateway.ClientLoggerExtensionImpl{
		GatewayRepository: repository,
	}
	service := pg.ServiceImp{
		Repository: repository,
		Client:     client,
	}

	kafka := exkafka.KafkaConsumer{
		Service: service,
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %+v \n", err)
	}
}

func getGormDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(viper.GetString("postgresSource")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to initial gorm DB")
	}

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func get_token() (*string, error) {

	url := fmt.Sprintf("https://tp.tax.gov.ir/req/api/self-tsp/sync/GET_TOKEN")

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
	signature, err := utility.SignAndVerify(normalized)
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

func getServerList() string {
	url := fmt.Sprintf("https://tp.tax.gov.ir/req/api/tsp/sync/GET_SERVER_INFORMATION")
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
