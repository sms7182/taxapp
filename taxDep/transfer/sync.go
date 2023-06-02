package transfer

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"tax-management/taxDep/types"
)

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if priv != nil {
		return priv.(*rsa.PrivateKey), nil
	}
	return nil, nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}
func (t *Transfer) SendPacket(taxRawId *uint, taxProcessId *uint, requestUniqueId string, packet *types.RequestPacket, version string, token string, encrypt, sign bool, prvKey string, customerId string) (*types.SyncResponse, error) {
	if packet == nil {
		return nil, nil
	}
	serverPubKey, id, err := t.GetServerPublicKeyWithCustomerId(customerId)
	if err != nil && serverPubKey == nil {
		return nil, nil
	}
	packet.EncryptionKeyId = id
	print(id)
	headers := make(map[string]string)
	t.fillEssentialHeader(headers)
	if len(token) > 0 {
		headers["Authorization"] = token
	}

	if sign {
		t.signPacket(packet, prvKey)
	}

	if encrypt {
		t.encryptPacket(packet, serverPubKey)
	}

	normalizedForm, err := t.cfg.normalizer(t.mergePacketAndHeaders(packet, headers))
	if err != nil {
		return nil, err
	}
	privateKey, err := ParseRsaPrivateKeyFromPemStr(prvKey)
	if err != nil {
		return nil, err
	}
	requestSign, err := t.cfg.signer([]byte(normalizedForm), privateKey)
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

	url := t.cfg.baseURL + "sync/" + version

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqJsonBody))
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

func (t *Transfer) SendPacketInquiry(taxRawId *uint, taxProcessId *uint, requestUniqueId string, packet *types.RequestPacket, version string, token string, encrypt, sign bool, privateKey string, customerId string) (*types.SyncResponse2, error) {
	rsaPrv, err := ParseRsaPrivateKeyFromPemStr(privateKey)
	if packet == nil {
		return nil, nil
	}
	if rsaPrv == nil {
		return nil, nil
	}
	serverPubKey, id, err := t.GetServerPublicKeyWithCustomerId(customerId)
	if err != nil && serverPubKey == nil {
		return nil, nil
	}
	packet.EncryptionKeyId = id
	headers := make(map[string]string)
	t.fillEssentialHeader(headers)
	if len(token) > 0 {
		headers["Authorization"] = token
	}

	if sign {
		t.signPacket(packet, privateKey)
	}

	if encrypt {
		t.encryptPacket(packet, serverPubKey)
	}

	normalizedForm, err := t.cfg.normalizer(t.mergePacketAndHeaders(packet, headers))
	if err != nil {
		return nil, err
	}

	requestSign, err := t.cfg.signer([]byte(normalizedForm), rsaPrv)
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
