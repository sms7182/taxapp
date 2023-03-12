package utility

import "tax-management/external/exkafka/messages"

type SignatureFirstTypeRequest struct {
	Authorization  string
	RequestTraceId string            `json:"requestTraceId"`
	TimeStamp      string            `json:"timestamp"`
	Packets        []PacketFirstType `json:"packets"`
}

type PacketFirstType struct {
	Uid             string             `json:"uid"`
	PacketType      string             `json:"packetType"`
	Retry           bool               `json:"retry"`
	Data            messages.AfterData `json:"data"`
	EncryptionKeyId string             `json:"encryptionKeyId"`
	SymmetricKey    string             `json:"symmetricKey"`
	IV              string             `json:"iv"`
	FiscalId        string             `json:"fiscalId"`
	DataSignature   string             `json:"dataSignature"`
}
