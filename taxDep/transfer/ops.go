package transfer

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tax-management/taxDep/types"

	"github.com/google/uuid"
)

func (t *Transfer) fillEssentialHeader(headers map[string]string) {
	unixMilli := fmt.Sprint(time.Now().UnixMilli())
	if _, ok := headers[RequestTraceIDHeader]; !ok {
		headers[RequestTraceIDHeader] = unixMilli
	}

	if _, ok := headers[TimestampHeader]; !ok {
		headers[TimestampHeader] = unixMilli
	}
}

func (t *Transfer) signPacket(packet *types.RequestPacket, privateKey string) error {
	rsaPrv, err := ParseRsaPrivateKeyFromPemStr(privateKey)
	normalizedForm, err := t.cfg.normalizer(packet.GetDataJSONMap())
	if err != nil {
		return err
	}

	sig, err := t.cfg.signer([]byte(normalizedForm), rsaPrv)
	if err != nil {
		return err
	}

	packet.DataSignature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func (t *Transfer) encryptPacket(packet *types.RequestPacket) error {
	aesKey := make([]byte, 32)
	rand.Read(aesKey)

	hexKeyBuff := new(bytes.Buffer)
	_, err := hex.NewEncoder(hexKeyBuff).Write(aesKey)
	if err != nil {
		return err
	}

	symmetricKey, err := t.cfg.pubKeyEncrypter(hexKeyBuff.Bytes(), t.serverPubKey)
	if err != nil {
		return err
	}

	packet.SymmetricKey = base64.StdEncoding.EncodeToString(symmetricKey)

	jsonBytes, err := json.Marshal(packet.Data)
	if err != nil {
		return err
	}

	cipherData, nonce, err := t.cfg.encrypter(t.xorBytes(jsonBytes, aesKey), aesKey, 16)
	if err != nil {
		return err
	}

	packet.IV = hex.EncodeToString(nonce)
	packet.Data = base64.StdEncoding.EncodeToString(cipherData)

	return nil
}

func (t *Transfer) xorBytes(a, b []byte) []byte {
	maxLen := len(a)
	if len(b) > len(a) {
		maxLen = len(b)
	}

	c := make([]byte, maxLen)
	for i := range c {
		c[i] = a[i%len(a)] ^ b[i%len(b)]
	}

	return c
}

func (t *Transfer) mergePacketAndHeaders(packet *types.RequestPacket, headers map[string]string) map[string]interface{} {
	result := packet.GetJSONMap()

	for k, v := range headers {
		result[k] = v
	}

	return result
}
func (t *Transfer) GetServerPublicKeyWithCustomerId(customerId string) (*rsa.PublicKey, string, error) {
	headers := make(map[string]string)
	t.fillEssentialHeader(headers)

	req := &types.SyncReq{
		Packet: &types.RequestPacket{
			UID:        uuid.NewString(),
			PacketType: "GET_SERVER_INFORMATION",
			FiscalId:   customerId,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, "", err
	}

	// u, err := url.Parse(t.cfg.baseURL)
	// if err != nil {
	// 	return nil, "", err
	// }
	url := t.cfg.baseURL + "sync/GET_SERVER_INFORMATION"
	// u.Path = path.Join(u.Path, filepath.Join("sync", "GET_SERVER_INFORMATION"))

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, "", err
	}

	headers["Content-Type"] = "application/json"
	for k, v := range headers {
		httpReq.Header[k] = []string{v}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, "", err
	}

	info := new(serverInfoResponse)

	if err := json.NewDecoder(resp.Body).Decode(info); err != nil {
		return nil, "", err
	}

	pubKey, err := t.parsePubKey(info.Result.Data.PublicKeys[0].Key)
	if err != nil {
		return nil, "", err
	}

	return pubKey, info.Result.Data.PublicKeys[0].ID, nil
}

func (t *Transfer) parsePubKey(keyBase64 string) (*rsa.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}
func (t *Transfer) GetServerPublicKey() (*rsa.PublicKey, string, error) {
	return t.GetServerPublicKeyWithCustomerId(t.cfg.clientID)
}

type serverInfoResponse struct {
	Result struct {
		Data struct {
			PublicKeys []struct {
				Key string `json:"key"`
				ID  string `json:"id"`
			} `json:"publicKeys"`
		} `json:"data"`
	} `json:"result"`
}
