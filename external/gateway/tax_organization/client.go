package taxorganization

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	cryptops "tax-management/cryptopts"
	"tax-management/external/exkafka/messages"
	"tax-management/pkg"
	"tax-management/terminal"
	"tax-management/types"
	"tax-management/utility"

	"time"

	"github.com/gofrs/uuid"
)

type TaxAPIType int

const (
	RequestTraceIDHeader = "requestTraceId"
	TimestampHeader      = "timestamp"
)

const (
	GetServerInformation TaxAPIType = iota
	GetToken
	GetFiscalInformation
	SendInvoice
	InquiryByUId
)

func (ts TaxAPIType) String() string {
	return []string{"GET_SERVER_INFORMATION", "GET_TOKEN", "GET_FISCAL_INFORMATION", "INVOICE.V01", "INQUIRY_BY_UID"}[ts]
}

func DefaultClientImpl() *ClientImpl {
	return &ClientImpl{

		normalizer: cryptops.NormalizeJsonObj,
		signer:     cryptops.SignPKCS1v15,
		encrypter:  cryptops.AesGCMNoPaddingEncrypt,
	}
}

type ClientImpl struct {
	HttpClient           pkg.ClientLoggerExtension
	Url                  string
	ServerInformationUrl string
	TokenUrl             string
	FiscalInformationUrl string
	InquiryByIdUrl       string
	SendInvoicUrl        string
	Repository           pkg.ClientRepository
	UserName             string
	Terminal             *terminal.Terminal
	normalizer           func(map[string]interface{}) (string, error)
	PrvKey               *rsa.PrivateKey
	PubKey               *rsa.PublicKey

	signer func([]byte, *rsa.PrivateKey) ([]byte, error)

	encrypter func(plainData, key []byte) (cipherData, nonce []byte, err error)
}

func (client ClientImpl) GetServerInformation() (*string, error) {
	url := client.Url + client.ServerInformationUrl
	id, _ := uuid.NewV4()
	var stui string
	stui = id.String()
	bodyReq := utility.BodyReq{
		Time: 2,
		Packet: utility.Packet{
			Uid:             stui,
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

func (client ClientImpl) SendPacket(packet *types.RequestPacket, version string, headers map[string]string, encrypt, sign bool, url string) (*types.SyncResponse, error) {

	if packet == nil {
		return nil, nil
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	client.fillEssentialHeader(headers)

	if sign {
		client.signPacket(packet)
	}

	if encrypt {
		client.encryptPacket(packet)
	}

	normalizedForm, err := client.normalizer(client.mergePacketAndHeaders(packet, headers))
	if err != nil {
		return nil, err
	}

	requestSign, err := client.signer([]byte(normalizedForm), client.PrvKey)
	if err != nil {
		return nil, err
	}

	reqJsonBody, err := json.Marshal(&types.SyncReq{
		SignedPacket: types.SignedPacket{
			Signature: base64.StdEncoding.EncodeToString(requestSign),
		},
		Packet: packet,
	})

	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqJsonBody))
	if err != nil {
		return nil, err
	}

	headers["Content-Type"] = "application/json"
	for k, v := range headers {
		httpReq.Header[k] = []string{v}
	}

	resp, err := client.HttpClient.Do(nil, nil, packet.UID, httpReq, version) // http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	sr := new(types.SyncResponse)

	return sr, json.NewDecoder(resp.Body).Decode(sr)
}
func (client ClientImpl) GetToken() (string, error) {
	t := client.Terminal
	packet := t.BuildRequestPacket(struct {
		Username string `json:"username"`
	}{
		Username: client.UserName,
	}, "GET_TOKEN")
	url := client.Url + client.TokenUrl
	resp, err := client.SendPacket(packet, "GET_TOKEN", nil, false, false, url)
	if err != nil {
		return "", err
	}
	fmt.Println(resp)
	token := (*resp).Result.Data["token"].(string)
	exp := time.UnixMilli(int64(resp.Result.Data["expiresIn"].(float64)))
	fmt.Printf("fix for redis exp %v", exp)
	return token, nil
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
			Uid:        stui,
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
			Uid:        stui,
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

func (client ClientImpl) SendInvoice(rawdata messages.RawTransaction) {
	t := client.Terminal
	packet := t.BuildRequestPacket(rawdata.After, SendInvoice.String())
	url := client.Url + client.SendInvoicUrl
	token, err := client.GetToken()
	if err != nil {
		log.Fatal("token has error %s", err)
	}
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	resp, err := client.SendPacket(packet, SendInvoice.String(), headers, true, true, url)
	if err != nil {

	}

	fmt.Println(resp)
}

func (t *ClientImpl) fillEssentialHeader(headers map[string]string) {
	unixMilli := fmt.Sprint(time.Now().UnixMilli())
	if _, ok := headers[RequestTraceIDHeader]; !ok {
		headers[RequestTraceIDHeader] = unixMilli
	}

	if _, ok := headers[TimestampHeader]; !ok {
		headers[TimestampHeader] = unixMilli
	}
}

func (t *ClientImpl) signPacket(packet *types.RequestPacket) error {
	normalizedForm, err := t.normalizer(packet.GetDataJSONMap())
	if err != nil {
		return err
	}

	sig, err := t.signer([]byte(normalizedForm), t.PrvKey)
	if err != nil {
		return err
	}

	packet.DataSignature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func (t *ClientImpl) encryptPacket(packet *types.RequestPacket) error {
	key := make([]byte, 32)
	rand.Read(key)

	packet.SymmetricKey = hex.EncodeToString(key)
	jsonBytes, err := json.Marshal(packet.Data)
	if err != nil {
		return err
	}

	cipherData, nonce, err := t.encrypter(t.xorBytes(jsonBytes, key), key)
	if err != nil {
		return err
	}

	packet.IV = hex.EncodeToString(nonce)
	packet.Data = base64.StdEncoding.EncodeToString(cipherData)

	return nil
}

func (t *ClientImpl) xorBytes(a, b []byte) []byte {
	if len(b) > len(a) {
		a, b = b, a
	}
	c := make([]byte, len(a))
	for i := range c {
		c[i] = a[i%len(a)] ^ b[i%len(b)]
	}
	return c
}

func (t *ClientImpl) mergePacketAndHeaders(packet *types.RequestPacket, headers map[string]string) map[string]interface{} {
	result := packet.GetJSONMap()

	for k, v := range headers {
		result[k] = v
	}

	return result
}

func (client ClientImpl) FirstGetToken() (*utility.TokenResponse, error) {

	url := client.Url + client.TokenUrl

	rqId, _ := uuid.NewV4()

	t := time.Now().UnixNano() / int64(time.Millisecond)
	tstr := strconv.FormatInt(t, 10)
	var stui string
	stui = rqId.String()
	sPacketReq := utility.SignaturePacketRequest{

		RequestTraceId: tstr,
		TimeStamp:      tstr,
		Packet: utility.Packet{
			Uid:        stui,
			PacketType: GetToken.String(),
			Retry:      false,
			Data: utility.TokenBody{
				UserName: client.UserName,
			},
			FiscalId: client.UserName,
		},
	}

	normalized, err := utility.Normalize(sPacketReq)
	// if err != nil {

	// 	fmt.Printf("normalize has error,%s", err.Error())
	// 	return nil, err
	// }
	//normalized := fmt.Sprintf("A11T1F#####A11T1F###GET_TOKEN#%s#false###%s#%s", tstr, tstr, stui)
	signature, err := utility.Sign(*normalized)

	if err != nil {
		fmt.Printf("sign has error %s", err.Error())

		return nil, err
	}
	postRequest := utility.PostDataRequest{

		Packet: utility.Packet{
			Uid:        stui,
			PacketType: GetToken.String(),
			Retry:      false,
			Data: utility.TokenBody{
				UserName: client.UserName,
			},
			FiscalId:      client.UserName,
			IV:            "",
			DataSignature: "",
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

	request.Header.Set("requestTraceId", tstr)
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
