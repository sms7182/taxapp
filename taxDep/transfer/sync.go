package transfer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"tax-management/taxDep/types"
)

func (t *Transfer) SendPacket(taxRawId *uint, taxProcessId *uint, requestUniqueId string, packet *types.RequestPacket, version string, token string, encrypt, sign bool) (*types.SyncResponse, error) {
	if packet == nil {
		return nil, nil
	}

	headers := make(map[string]string)
	t.fillEssentialHeader(headers)
	if len(token) > 0 {
		headers["Authorization"] = token
	}

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

	u.Path = path.Join(u.Path, filepath.Join("sync", version))

	httpReq, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(reqJsonBody))
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

	resp, err := t.HttpClientLogger.Do(taxRawId, taxProcessId, requestUniqueId, httpReq, version)
	if err != nil {
		return nil, err
	}

	sr := &types.SyncResponse{
		HttpStatusCode: resp.StatusCode,
	}

	return sr, json.NewDecoder(resp.Body).Decode(sr)
}

func (t *Transfer) SendPacketInquiry(taxRawId *uint, taxProcessId *uint, requestUniqueId string, packet *types.RequestPacket, version string, token string, encrypt, sign bool) (*types.SyncResponse2, error) {
	if packet == nil {
		return nil, nil
	}

	headers := make(map[string]string)
	t.fillEssentialHeader(headers)
	if len(token) > 0 {
		headers["Authorization"] = token
	}

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

	u.Path = path.Join(u.Path, filepath.Join("sync", version))

	httpReq, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(reqJsonBody))
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

	resp, err := t.HttpClientLogger.Do(taxRawId, taxProcessId, requestUniqueId, httpReq, version)
	if err != nil {
		return nil, err
	}

	sr := &types.SyncResponse2{
		HttpStatusCode: resp.StatusCode,
	}

	return sr, json.NewDecoder(resp.Body).Decode(sr)
}
