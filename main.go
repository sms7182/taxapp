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
	"time"

	"github.com/gofrs/uuid"
)

func main() {

}

func send() {
	//	url := fmt.Sprintf("https://tp.tax.gov.ir/req/")
	//	trace_id, _ := uuid.NewV4()

	//	request, err := http.NewRequest("POST", get_tokenUrl, nil)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// }
	// request.Header.Set("Authorization", "")
	// request.Header.Set("requestTraceId", trace_id.String())
	// request.Header.Set("timestamp", time.Now().String())
}
func get_token() (*string, error) {
	url := fmt.Sprintf("https://tp.tax.gov.ir/req/")

	tokenUrl := fmt.Sprintf(url, "api/self-tsp/sync/GET_TOKEN")
	token_req := TokenBody{
		UserName: "",
	}
	jsonBytes, err := json.Marshal(token_req)
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
	request.Header.Set("timestamp", time.Now().String())
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

func normalize(obj interface{}) {
	t := reflect.TypeOf(obj)
	if kind := t.Kind(); kind != reflect.Struct {
		log.Fatalf("This program expects to work on a struct; we got a %v instead.", kind)
	}

	fields := traverseStruct(t, obj)
	sort.Strings(fields)

}
func (v Value) Interface() interface {
}
func traverseStruct(rType reflect.Type, value interface{}) []string {
	if kind := rType.Kind(); kind != reflect.Struct {
		log.Fatalf("expects to a struct type %v", kind)
	}
	p := reflect.ValueOf(&value).Elem()
	var fields []string
	for i := 0; i < rType.NumField(); i++ {
		f := rType.Field(i)

		if f.Type.Kind() == reflect.Struct {
			nested_fields := traverseStruct(f.Type, f.Interface())
			for j := range nested_fields {
				nested := nested_fields[j]
				fields = append(fields, fmt.Sprint(f.Name, ".", nested))
			}

		} else {
			fields = append(fields, f.Name)
		}
	}
	return fields
}
