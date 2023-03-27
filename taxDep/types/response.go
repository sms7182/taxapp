package types

type SyncResponsePacket struct {
	UID             string                 `json:"uid"`
	PacketType      string                 `json:"packetType"`
	EncryptionKeyId string                 `json:"encryptionKeyId"`
	SymmetricKey    string                 `json:"symmetricKey"`
	IV              string                 `json:"iv"`
	Data            map[string]interface{} `json:"data"`
	Timestamp       *int64                 `json:"timestamp"`
}

type SyncResponsePacketInquiry struct {
	UID             string          `json:"uid"`
	PacketType      string          `json:"packetType"`
	EncryptionKeyId string          `json:"encryptionKeyId"`
	SymmetricKey    string          `json:"symmetricKey"`
	IV              string          `json:"iv"`
	Data            []InquiryResult `json:"data"`
	Timestamp       *int64          `json:"timestamp"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
	Code   string `json:"errorCode"`
}

type AsyncResponsePacket struct {
	UID             string `json:"uid"`
	ReferenceNumber string `json:"referenceNumber"`
	ErrorCode       string `json:"errorCode"`
	ErrorDetail     string `json:"errorDetail"`
}

type Response[T any] struct {
	Timestamp      int64           `json:"timestamp"`
	Result         T               `json:"result"`
	Errors         []ErrorResponse `json:"errors"`
	HttpStatusCode int             `json:"-"`
}

type ResponseList[T any] struct {
	Timestamp      int64           `json:"timestamp"`
	Result         []T             `json:"result"`
	Errors         []ErrorResponse `json:"errors"`
	HttpStatusCode int             `json:"-"`
}
type InquiryResult struct {
	ReferenceNumber string            `json:"referenceNumber"`
	UID             string            `json:"uid"`
	FiscalID        string            `json:"fiscalId"`
	Status          string            `json:"status"`
	PacketType      string            `json:"packetType"`
	Data            InquiryResultData `json:"data"`
}

type InquiryResultData struct {
	ConfirmationReferenceID string               `json:"confirmationReferenceId"`
	Error                   []any                `json:"error"`
	Success                 bool                 `json:"success"`
	Warning                 []InquiryDataWarning `json:"warning"`
}

type InquiryDataWarning struct {
	Code   string `json:"code"`
	Detail []any  `json:"detail"`
	Msg    string `json:"msg"`
}

type (
	AsyncResponse = ResponseList[AsyncResponsePacket]
	SyncResponse  = Response[SyncResponsePacket]
	SyncResponse2 = Response[SyncResponsePacketInquiry]
)
