package transfer

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"tax-management/types"
	"time"
)

type Transfer struct {
	cfg *ApiConfig
}

func NewApiTransfer(cfg *ApiConfig) *Transfer {
	return &Transfer{
		cfg: cfg,
	}
}

func (t *Transfer) SendPacket(packet *types.RequestPacket, version string, headers map[string]string, encrypt, sign bool) (*types.SyncResponse, error) {
	if packet == nil {
		return nil, nil
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	t.fillEssentialHeader(headers)

	if sign {
		t.signPacket(packet)
	}

	if encrypt {
		t.encryptPacket(packet)
	}

	normalizedForm, err := t.cfg.normalizer(t.mergePacketAndHeaders(packet, headers))
	if err != nil {
		return nil, err
	}

	requestSign, err := t.cfg.signer([]byte(normalizedForm), t.cfg.prvKey)
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

	u, err := url.Parse(t.cfg.baseURL)
	if err != nil {
		return nil, err
	}

	u.Path = u.Path + "/sync/" + version //path.Join(u.Path, filepath.Join("sync", version))

	httpReq, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(reqJsonBody))
	if err != nil {
		return nil, err
	}

	headers["Content-Type"] = "application/json"
	for k, v := range headers {
		httpReq.Header[k] = []string{v}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	sr := new(types.SyncResponse)

	return sr, json.NewDecoder(resp.Body).Decode(sr)
}

func (t *Transfer) fillEssentialHeader(headers map[string]string) {
	unixMilli := fmt.Sprint(time.Now().UnixMilli())
	if _, ok := headers[RequestTraceIDHeader]; !ok {
		headers[RequestTraceIDHeader] = unixMilli
	}

	if _, ok := headers[TimestampHeader]; !ok {
		headers[TimestampHeader] = unixMilli
	}
}

func (t *Transfer) signPacket(packet *types.RequestPacket) error {
	normalizedForm, err := t.cfg.normalizer(packet.GetDataJSONMap())
	if err != nil {
		return err
	}

	sig, err := t.cfg.signer([]byte(normalizedForm), t.cfg.prvKey)
	if err != nil {
		return err
	}

	packet.DataSignature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func (t *Transfer) encryptPacket(packet *types.RequestPacket) error {
	key := make([]byte, 32)
	rand.Read(key)

	packet.SymmetricKey = hex.EncodeToString(key)
	jsonBytes, err := json.Marshal(packet.Data)
	if err != nil {
		return err
	}

	cipherData, nonce, err := t.cfg.encrypter(t.xorBytes(jsonBytes, key), key)
	if err != nil {
		return err
	}

	packet.IV = hex.EncodeToString(nonce)
	packet.Data = base64.StdEncoding.EncodeToString(cipherData)

	return nil
}

func (t *Transfer) xorBytes(a, b []byte) []byte {
	if len(b) > len(a) {
		a, b = b, a
	}
	c := make([]byte, len(a))
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
