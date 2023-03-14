package utility

type InquiryByIdRequest struct {
	Authorization  string
	ContentType    string            `json:"Content-Type"`
	RequestTraceId string            `json:"requestTraceId"`
	TimeStamp      string            `json:"timestamp"`
	Packet         InquiryByIdPacket `json:"packet"`
}
type InquiryByIdPacket struct {
	Uid             string            `json:"uid"`
	PacketType      string            `json:"packetType"`
	Retry           bool              `json:"retry"`
	Data            []InquiryByIdBody `json:"data"`
	EncryptionKeyId string            `json:"encryptionKeyId"`
	SymmetricKey    string            `json:"symmetricKey"`
	IV              string            `json:"iv"`
	FiscalId        string            `json:"fiscalId"`
	DataSignature   string            `json:"dataSignature"`
}

type InquiryByIdBody struct {
	UId      string `json:"uid"`
	FiscalId string `json:"fiscalId"`
}
type InquiryByIdResponse struct {
	UID        interface{} `json:"uid"`
	PacketType string      `json:"packetType"`
	Data       []struct {
		ReferenceNumber string `json:"referenceNumber"`
		UID             string `json:"uid"`
		Status          string `json:"status"`
		Data            struct {
			ConfirmationReferenceID interface{} `json:"confirmationReferenceId"`
			TaxResult               string      `json:"taxResult"`
		} `json:"data"`
		PacketType string      `json:"packetType"`
		FiscalID   interface{} `json:"fiscalId"`
	} `json:"data"`
	EncryptionKeyID interface{} `json:"encryptionKeyId"`
	SymmetricKey    interface{} `json:"symmetricKey"`
	Iv              interface{} `json:"iv"`
}
