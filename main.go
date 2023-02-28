package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func main() {
	getServerList()
	lt := LevelTwo{
		FirstLevelTwo: 75,
	}
	base := Base{

		SecondLevelOne: 70,
		ThirdLevelOne:  lt,
		FirstLevelOne:  "Test",
	}
	normalize(base)
}

type SignaturePacketRequest struct {
	Authorization  string
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      string `json:"timestamp"`
	Packet         Packet `json:"packet"`
}
type SignaturePacketsRequest struct {
	Authorization  string
	RequestTraceId string   `json:"requestTraceId"`
	TimeStamp      string   `json:"timestamp"`
	Packets        []Packet `json:"packets"`
}

type BodyReq struct {
	Time   int    `json:"time"`
	Packet Packet `json:"packet"`
}
type Packet struct {
	Uid             string      `json:"uid"`
	PacketType      string      `json:"packetType"`
	Retry           bool        `json:"retry"`
	Data            interface{} `json:"data"`
	EncryptionKeyId string      `json:"encryptionKeyId"`
	SymmetricKey    string      `json:"symmetricKey"`
	IV              string      `json:"iv"`
	FiscalId        string      `json:"fiscalId"`
	DataSignature   string      `json:"dataSignature"`
}
type Base struct {
	SecondLevelOne int
	ThirdLevelOne  LevelTwo
	FirstLevelOne  string
}
type LevelTwo struct {
	FirstLevelTwo int
}

func getServerList() {
	url := fmt.Sprintf("https://tp.tax.gov.ir/req/api/tsp/sync/GET_SERVER_INFORMATION")
	id, _ := uuid.NewV4()
	bodyReq := BodyReq{
		Time: 2,
		Packet: Packet{
			Uid:             id.String(),
			PacketType:      "GET_SERVER_INFORMATION",
			Retry:           false,
			Data:            nil,
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
	var bodyObj BodyResponse

	if err != nil {
		fmt.Printf("read response has error %s", err.Error())

	}
	json.Unmarshal(body, &bodyObj)
}

type BodyResponse struct {
	Signature      *string `json:"signature"`
	SignatureKeyId *string `json:"signatureKeyId"`
	TimeStamp      int64   `json:"timestamp"`
	Result         struct {
		UId        string `json:"uid"`
		PacketType string `json:"packetType"`
		Data       struct {
			ServerTime int64 `json:"serverTime"`
			PublicKeys []struct {
				Key       string `json:"key"`
				Id        string `json:"id"`
				Algorithm string `json:"RSA"`
				Purpose   int    `json:"purpose"`
			} `json:"publicKeys"`
		} `json:"data"`
		EncryptionKeyId *string `json:"encryptionKeyId"`
		SymmetricKey    *string `json:"symmetricKey"`
		Iv              *string `json:"iv"`
	} `json:"result"`
}

func get_token() (*string, error) {

	url := fmt.Sprintf("https://tp.tax.gov.ir/req/")

	tokenUrl := fmt.Sprintf(url, "api/self-tsp/sync/GET_TOKEN")

	rqId, _ := uuid.NewV4()
	timeNow := time.Now()

	sPacketReq := SignaturePacketRequest{
		Authorization:  "",
		RequestTraceId: rqId.String(),
		TimeStamp:      timeNow.String(),
		Packet: Packet{
			Uid:        rqId.String(),
			PacketType: "GET_TOKEN",
			Retry:      false,
			Data: TokenBody{
				UserName: "A118GE",
			},
		},
	}
	jsonBytes, err := json.Marshal(sPacketReq)
	if err != nil {
		fmt.Printf("json marshal has error %s", err.Error())
		return nil, err
	}
	reader := bytes.NewReader(jsonBytes)

	request, err := http.NewRequest("POST", tokenUrl, reader)

	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return nil, err
	}
	traceId, _ := uuid.NewV4()
	request.Header.Set("requestTraceId", traceId.String())
	request.Header.Set("timestamp", timeNow.String())
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
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Printf("responseJson has error %s", err.Error())
		return nil, err
	}
	return &tokenResponse.Token, nil
}

func normalize(obj interface{}) {
	t := reflect.TypeOf(obj)
	if kind := t.Kind(); kind != reflect.Struct {
		log.Fatalf("This program expects to work on a struct; we got a %v instead.", kind)
	}
	maps := make(map[string]interface{})
	fields := traverseObject(t)

	sort.Strings(fields)
	value := reflect.ValueOf(&obj).Elem().Elem()

	for i := range fields {
		field := fields[i]
		splited := strings.Split(field, ".")
		if len(splited) > 1 {
			temp := value
			for j := 0; j < len(splited); j++ {
				tempValue := temp.FieldByName(splited[j])
				tempInterface := tempValue.Interface()
				if j+1 < len(splited) {
					temp = reflect.ValueOf(&tempInterface).Elem().Elem()
				} else {
					obj := reflect.ValueOf(&tempInterface).Elem()

					maps[field] = obj
				}
			}

		} else {
			fieldValue := value.FieldByName(field)

			maps[field] = fieldValue.Interface()
		}
	}
	str := ""
	for i := range fields {
		vl := maps[fields[i]]
		if vl != nil {
			if i == 0 {
				str = fmt.Sprintf("%v", vl)
			} else {
				str = fmt.Sprintf("%s#%v", str, vl)
			}
		} else {
			str = fmt.Sprintf("%s###", str)
		}
	}
	fmt.Print(str)

}

func traverseObject(rType reflect.Type) []string {
	if kind := rType.Kind(); kind != reflect.Struct {
		log.Fatalf("expects to a struct type %v", kind)
	}
	var fields []string
	for i := 0; i < rType.NumField(); i++ {
		f := rType.Field(i)
		field := rType.Field(i)
		if f.Type.Kind() == reflect.Struct {

			nested_fields := traverseObject(f.Type)
			for j := range nested_fields {
				nested := nested_fields[j]
				fields = append(fields, fmt.Sprint(field.Name, ".", nested))
			}

		} else {

			fields = append(fields, field.Name)
		}
	}
	return fields
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}
type TokenBody struct {
	UserName string `json:"username"`
}

type ToNormalizeData struct {
	Authorization  string
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      int64  `json:"timestamp"`
	Body           BodyRequest
}
type BodyRequest struct {
	Packets        []interface{} `json:"packets"`
	Signature      string        `json:"signature"`
	SignatureKeyId string        `json:"signatureKeyId"`
}
