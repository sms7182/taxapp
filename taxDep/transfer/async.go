package transfer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"tax-management/taxDep/types"
)

func (t *Transfer) SendPackets(
	taxRawId *uint,
	taxProcessId *uint,
	requestUniqueId string,
	packets []types.RequestPacket,
	version string,
	token string,
	encrypt,
	sign bool,
	privateKey string) (*types.AsyncResponse, error) {
	if len(packets) == 0 {
		return nil, nil
	}
	rsaPrv, err := ParseRsaPrivateKeyFromPemStr(privateKey)
	headers := make(map[string]string)
	t.fillEssentialHeader(headers)
	if len(token) > 0 {
		headers["Authorization"] = token
	}

	for i := range packets {
		packets[i].EncryptionKeyId = t.pubKeyID
	}

	if sign {
		for i := range packets {
			t.signPacket(&packets[i], privateKey)
		}
	}

	if encrypt {
		for i := range packets {
			t.encryptPacket(&packets[i])
		}
	}

	m, err := t.getPacketsMap(packets, headers)
	if err != nil {
		return nil, err
	}

	normalizedFrom, err := t.cfg.normalizer(m)
	if err != nil {
		return nil, err
	}

	requestSign, err := t.cfg.signer([]byte(normalizedFrom), rsaPrv)
	if err != nil {
		return nil, err
	}

	requestBody, err := json.MarshalIndent(&types.AsyncReq{
		SignedPacket: types.SignedPacket{
			Signature: base64.StdEncoding.EncodeToString(requestSign),
		},
		Packets: packets,
	}, "", " ")

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	url := t.cfg.baseURL + "async/" + version

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	headers["Content-Type"] = "application/json"
	if len(token) > 0 {
		headers["Authorization"] = "Bearer " + token
	}

	for k, v := range headers {
		httpReq.Header[k] = []string{v}
	}

	resp, err := t.HttpClientLogger.Do(taxRawId, taxProcessId, requestUniqueId, httpReq, "SendInvoice")
	if err != nil {
		return nil, err
	}

	sr := &types.AsyncResponse{
		HttpStatusCode: resp.StatusCode,
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return sr, json.Unmarshal(all, sr)
}

func (t *Transfer) getPacketsMap(packets []types.RequestPacket, headers map[string]string) (map[string]any, error) {
	obj := &struct {
		Packets []types.RequestPacket `json:"packets"`
	}{
		Packets: packets,
	}

	m := make(map[string]any)

	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	for k, v := range headers {
		m[k] = v
	}

	return m, nil
}
