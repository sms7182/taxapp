package utility

type SignaturePacketRequest struct {
	//Authorization  string
	ContentType    string `json:"Content-Type"`
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      string `json:"timestamp"`
	Packet         Packet `json:"packet"`
}
type PostDataRequest struct {
	// ContentType    string  `json:"Content-Type"`
	// RequestTraceId string  `json:"requestTraceId"`
	// TimeStamp      string  `json:"timestamp"`
	Packet    Packet  `json:"packet"`
	Signature *string `json:"signature"`
	Time      int     `json:"time"`
}
type TestSignature struct {
	Packet PacketTest `json:"packet"`
}
type PacketTest struct {
	Data interface{} `json:"data" {asghar}`
}
type SignaturePacketsRequest struct {
	Authorization  string
	RequestTraceId string   `json:"requestTraceId"`
	TimeStamp      string   `json:"timestamp"`
	Packets        []Packet `json:"packets"`
}

type BodyReq struct {
	Time   int    `json:"time"`
	Packet Packet `json:"packet"`
}
type Packet struct {
	Uid             *string   `json:"uid"`
	PacketType      string    `json:"packetType"`
	Retry           bool      `json:"retry"`
	Data            TokenBody `json:"data"`
	EncryptionKeyId string    `json:"encryptionKeyId"`
	SymmetricKey    string    `json:"symmetricKey"`
	IV              string    `json:"iv"`
	FiscalId        string    `json:"fiscalId"`
	DataSignature   string    `json:"dataSignature"`
}
type Base struct {
	SecondLevelOne int
	ThirdLevelOne  LevelTwo
	FirstLevelOne  string
}
type LevelTwo struct {
	FirstLevelTwo int
}

type BodyResponse struct {
	Signature      *string `json:"signature"`
	SignatureKeyId *string `json:"signatureKeyId"`
	TimeStamp      int64   `json:"timestamp"`
	Result         struct {
		UId        string `json:"uid"`
		PacketType string `json:"packetType"`
		Data       struct {
			ServerTime int64 `json:"serverTime"`
			PublicKeys []struct {
				Key       string `json:"key"`
				Id        string `json:"id"`
				Algorithm string `json:"RSA"`
				Purpose   int    `json:"purpose"`
			} `json:"publicKeys"`
		} `json:"data"`
		EncryptionKeyId *string `json:"encryptionKeyId"`
		SymmetricKey    *string `json:"symmetricKey"`
		Iv              *string `json:"iv"`
	} `json:"result"`
}
type TestData struct {
	Name string `json:"name"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}
type TokenBody struct {
	UserName string `json:"username"`
}

type ToNormalizeData struct {
	Authorization  string
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      int64  `json:"timestamp"`
	Body           BodyRequest
}
type BodyRequest struct {
	Packets        []interface{} `json:"packets"`
	Signature      string        `json:"signature"`
	SignatureKeyId string        `json:"signatureKeyId"`
}
