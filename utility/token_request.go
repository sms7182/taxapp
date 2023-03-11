package utility

type SignaturePacketRequest struct {
	//Authorization  string
	ContentType    string `json:"Content-Type"`
	RequestTraceId string `json:"requestTraceId"`
	TimeStamp      string `json:"timestamp"`
	Packet         Packet `json:"packet"`
}
type PostDataRequest struct {
	ContentType    string  `json:"Content-Type"`
	RequestTraceId string  `json:"requestTraceId"`
	TimeStamp      string  `json:"timestamp"`
	Packet         Packet  `json:"packet"`
	Signature      *string `json:"signature"`
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
	Uid             string    `json:"uid"`
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

// func Encrypt() string {
// 	data := DataToEncrypt{
// 		Description: "test encryption",
// 	}
// 	data.Detail = append(data.Detail, Detail{
// 		Quantity: 2,
// 		Price:    250,
// 	})
// 	key := make([]byte, 32)
// 	if _, err := rand.Read(key); err != nil {
// 		fmt.Printf("has error %s", err.Error())
// 	}
// 	bytes, err := json.Marshal(key)
// 	if err != nil {
// 		fmt.Printf("has error %s", err.Error())
// 	}

// 	block, err := aes.NewCipher(key)

// 	aesGCM, err := cipher.NewGCM(block)

// 	nonce := make([]byte, aesGCM.NonceSize())
// 	tag := make([]byte, 16)
// 	if _, err := rand.Read(tag); err != nil {
// 		fmt.Printf("has error %s", err.Error())
// 	}
// 	result := aesGCM.Seal(nonce, nonce, bytes, tag)
// 	return fmt.Sprintf("%x", result)
// }
