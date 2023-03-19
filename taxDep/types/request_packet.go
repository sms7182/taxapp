package types

import "encoding/json"

type RequestPacket struct {
	UID             string      `json:"uid"`
	PacketType      string      `json:"packetType"`
	Retry           bool        `json:"retry"`
	Data            interface{} `json:"data"`
	EncryptionKeyId string      `json:"encryptionKeyId"`
	SymmetricKey    string      `json:"symmetricKey"`
	IV              string      `json:"iv"`
	FiscalId        string      `json:"fiscalId"`
	DataSignature   string      `json:"dataSignature"`
}

func (p *RequestPacket) GetDataJSONMap() map[string]interface{} {
	m := make(map[string]interface{})
	b, _ := json.Marshal(p.Data)
	json.Unmarshal(b, &m)
	return m
}

func (p *RequestPacket) GetJSONMap() map[string]interface{} {
	m := make(map[string]interface{})
	b, _ := json.Marshal(*p)
	json.Unmarshal(b, &m)
	return m
}
