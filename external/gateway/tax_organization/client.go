package taxorganization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tax-management/pkg"
	"tax-management/utility"
	"time"

	"github.com/gofrs/uuid"
)

type TaxAPIType int

const (
	GetServerInformation TaxAPIType = iota
	GetToken
)

func (ts TaxAPIType) String() string {
	return []string{"GET_SERVER_INFORMATION", "GET_TOKEN"}[ts]
}

type ClientImpl struct {
	HttpClient           pkg.ClientLoggerExtension
	Url                  string
	ServerInformationUrl string
	Repository           pkg.ClientRepository
}

func (client ClientImpl) GetServerInformation() (*string, error) {
	url := client.Url + client.ServerInformationUrl
	id, _ := uuid.NewV4()

	bodyReq := utility.BodyReq{
		Time: 2,
		Packet: utility.Packet{
			Uid:             id.String(),
			PacketType:      GetServerInformation.String(),
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
		return nil, err
	}

	jsonBytes := bytes.NewReader(marshaled)
	request, err := http.NewRequest("POST", url, jsonBytes)
	if err != nil {
		fmt.Printf("Create Post request has error %s", err.Error())
		return nil, err
	}
	request.Header.Set("requestTraceId", id.String())
	request.Header.Set("timestampt", time.Now().String())
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.HttpClient.Do(nil, nil, id.String(), request, "GET_SERVER_INFORMATION")
	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return nil, err

	}
	body, err := ioutil.ReadAll(resp.Body)
	var bodyObj utility.BodyResponse

	if err != nil {
		fmt.Printf("read response has error %s", err.Error())
		return nil, err

	}
	json.Unmarshal(body, &bodyObj)

	rs := bodyObj.Result.Data.PublicKeys[0].Key
	return &rs, nil
}
