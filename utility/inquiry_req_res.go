package utility

type InquiryByIdRequest struct {
	Authorization  string
	ContentType    string `json:"Content-Type"`
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      string `json:"timestamp"`
	Packet         Packet `json:"packet"`
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
