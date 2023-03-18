package taxorganization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"tax-management/pkg"
	"tax-management/utility"
	"time"

	"github.com/gofrs/uuid"
)

type TaxAPIType int

const (
	GetServerInformation TaxAPIType = iota
	GetToken
	GetFiscalInformation
	InquiryByUId
)

func (ts TaxAPIType) String() string {
	return []string{"GET_SERVER_INFORMATION", "GET_TOKEN", "GET_FISCAL_INFORMATION", "INQUIRY_BY_UID"}[ts]
}

type ClientImpl struct {
	HttpClient           pkg.ClientLoggerExtension
	Url                  string
	ServerInformationUrl string
	TokenUrl             string
	FiscalInformationUrl string
	InquiryByIdUrl       string
	Repository           pkg.ClientRepository
	UserName             string
}

func (client ClientImpl) GetServerInformation() (*string, error) {
	url := client.Url + client.ServerInformationUrl
	id, _ := uuid.NewV4()
	var stui string
	stui = id.String()
	bodyReq := utility.BodyReq{
		Time: 2,
		Packet: utility.Packet{
			Uid:             &stui,
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

	resp, err := client.HttpClient.Do(nil, nil, id.String(), request, GetServerInformation.String())
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

func (client ClientImpl) GetToken() (*utility.TokenResponse, error) {

	url := client.Url + client.TokenUrl

	rqId, _ := uuid.NewV4()
	timeNow := time.Now().Unix()
	tstr := strconv.FormatInt(timeNow, 10)
	var stui string
	stui = rqId.String()
	sPacketReq := utility.SignaturePacketRequest{

		RequestTraceId: rqId.String(),
		TimeStamp:      tstr,
		ContentType:    "application/json",
		Packet: utility.Packet{
			Uid:        &stui,
			PacketType: GetToken.String(),
			Retry:      false,
			Data: utility.TokenBody{
				UserName: client.UserName,
			},
		},
	}

	normalized, err := utility.Normalize(sPacketReq)
	if err != nil {
		// update for retry has error in normalize
		// notif to developer
		fmt.Printf("normalize has error,%s", err.Error())
		return nil, err
	}
	signature, err := utility.Sign(*normalized) //utility.SignAndVerify(normalized)
	if err != nil {
		fmt.Printf("sign has error %s", err.Error())
		// update for retry has error in normalize
		// notif to developer
		return nil, err
	}
	postRequest := utility.PostDataRequest{

		Time: 1,
		Packet: utility.Packet{
			Uid:        nil,
			PacketType: GetToken.String(),
			Retry:      false,
			Data: utility.TokenBody{
				UserName: client.UserName,
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
	//traceId, _ := uuid.NewV4()
	request.Header.Set("requestTraceId", rqId.String())
	request.Header.Set("timestamp", tstr)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.HttpClient.Do(nil, nil, rqId.String(), request, GetToken.String())
	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return nil, err

	}

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
	return &tokenResponse, nil
}

func (client ClientImpl) GetFiscalInformation(token string) {
	url := client.Url + client.FiscalInformationUrl
	rqId, _ := uuid.NewV4()
	timeNow := time.Now().Unix()
	tstr := strconv.FormatInt(timeNow, 10)
	sPacketReq := utility.FiscalInformationRequest{
		Authorization:  token,
		RequestTraceId: rqId.String(),
		TimeStamp:      tstr,
		ContentType:    "application/json",
		Packet: utility.FiscalInformationPacket{
			Uid:        rqId.String(),
			PacketType: GetFiscalInformation.String(),
			Retry:      false,
			Data:       nil,
		},
	}

	normalized, err := utility.Normalize(sPacketReq)
	if err != nil {
		// update for retry has error in normalize
		// notif to developer
		fmt.Printf("normalize has error,%s", err.Error())
		return
	}
	signature, err := utility.SignAndVerify(normalized)
	if err != nil {
		fmt.Printf("sign has error %s", err.Error())
		// update for retry has error in normalize
		// notif to developer
		return
	}
	var stui string
	stui = rqId.String()
	postRequest := utility.PostDataRequest{
		// RequestTraceId: rqId.String(),
		// TimeStamp:      tstr,
		// ContentType:    "application/json",
		Packet: utility.Packet{
			Uid:        &stui,
			PacketType: GetFiscalInformation.String(),
			Retry:      false,
			Data: utility.TokenBody{
				UserName: client.UserName,
			},
		},
		Signature: signature,
	}
	jsonBytes, err := json.Marshal(postRequest)
	if err != nil {
		fmt.Printf("json marshal has error %s", err.Error())
		return
	}
	reader := bytes.NewReader(jsonBytes)

	request, err := http.NewRequest("POST", url, reader)

	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return
	}

	request.Header.Set("requestTraceId", rqId.String())
	request.Header.Set("timestamp", tstr)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.HttpClient.Do(nil, nil, rqId.String(), request, GetFiscalInformation.String())
	if err != nil {
		fmt.Printf("response has error %s", err.Error())

	}

	if err != nil {
		fmt.Printf("response has error %s", err.Error())

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response has error %s", err.Error())

	}
	var fiscalInfoResponse utility.FiscalInformationResponse
	err = json.Unmarshal(body, &fiscalInfoResponse)
	if err != nil {
		fmt.Printf("responseJson has error %s", err.Error())

	}

}

func (client ClientImpl) InquiryById(token string) {
	url := client.Url + client.InquiryByIdUrl
	rqId, _ := uuid.NewV4()
	timeNow := time.Now().Unix()
	tstr := strconv.FormatInt(timeNow, 10)

	sPacketReq := utility.InquiryByIdRequest{
		Authorization:  token,
		RequestTraceId: rqId.String(),
		TimeStamp:      tstr,
		ContentType:    "application/json",
		Packet: utility.InquiryByIdPacket{
			Uid:        rqId.String(),
			PacketType: InquiryByUId.String(),
			Retry:      false,
		},
	}
	//todo complete this code
	sPacketReq.Packet.Data = append(sPacketReq.Packet.Data, utility.InquiryByIdBody{
		UId:      "",
		FiscalId: "",
	})

	normalized, err := utility.Normalize(sPacketReq)
	if err != nil {
		// update for retry has error in normalize
		// notif to developer
		fmt.Printf("normalize has error,%s", err.Error())
		return
	}
	signature, err := utility.SignAndVerify(normalized)
	if err != nil {
		fmt.Printf("sign has error %s", err.Error())
		// update for retry has error in normalize
		// notif to developer
		return
	}
	var stui string
	stui = rqId.String()
	postRequest := utility.PostDataRequest{
		// RequestTraceId: rqId.String(),
		// TimeStamp:      tstr,
		// ContentType:    "application/json",
		Packet: utility.Packet{
			Uid:        &stui,
			PacketType: InquiryByUId.String(),
			Retry:      false,
			Data: utility.TokenBody{
				UserName: client.UserName,
			},
		},
		Signature: signature,
	}
	jsonBytes, err := json.Marshal(postRequest)
	if err != nil {
		fmt.Printf("json marshal has error %s", err.Error())
		return
	}
	reader := bytes.NewReader(jsonBytes)

	request, err := http.NewRequest("POST", url, reader)

	if err != nil {
		fmt.Printf("response has error %s", err.Error())
		return
	}

	request.Header.Set("requestTraceId", rqId.String())
	request.Header.Set("timestamp", tstr)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.HttpClient.Do(nil, nil, rqId.String(), request, InquiryByUId.String())
	if err != nil {
		fmt.Printf("response has error %s", err.Error())

	}

	if err != nil {
		fmt.Printf("response has error %s", err.Error())

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response has error %s", err.Error())

	}
	var inquiryResponse utility.InquiryByIdResponse
	err = json.Unmarshal(body, &inquiryResponse)
	if err != nil {
		fmt.Printf("responseJson has error %s", err.Error())

	}

}
