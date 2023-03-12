package utility

type FiscalInformationResponse struct {
	Signature      interface{} `json:"signature"`
	SignatureKeyID interface{} `json:"signatureKeyId"`
	Timestamp      int64       `json:"timestamp"`
	Result         struct {
		UID        interface{} `json:"uid"`
		PacketType string      `json:"packetType"`
		Data       struct {
			NameTrade     string  `json:"nameTrade"`
			FiscalStatus  string  `json:"fiscalStatus"`
			SaleThreshold float64 `json:"saleThreshold"`
			EconomicCode  string  `json:"economicCode"`
		} `json:"data"`
		EncryptionKeyID interface{} `json:"encryptionKeyId"`
		SymmetricKey    interface{} `json:"symmetricKey"`
		Iv              interface{} `json:"iv"`
	} `json:"result"`
}

type FiscalInformationRequest struct {
	Authorization  string
	ContentType    string `json:"Content-Type"`
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      string `json:"timestamp"`
	Packet         Packet `json:"packet"`
}
