package types

type SyncResponsePacket struct {
	UID             string                 `json:"uid"`
	PacketType      string                 `json:"packetType"`
	EncryptionKeyId string                 `json:"encryptionKeyId"`
	SymmetricKey    string                 `json:"symmetricKey"`
	IV              string                 `json:"iv"`
	Data            map[string]interface{} `json:"data"`
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

type (
	AsyncResponse = ResponseList[AsyncResponsePacket]
	SyncResponse  = Response[SyncResponsePacket]
)
