package types

type SyncReq struct {
	SignedPacket
	Packet *RequestPacket `json:"packet"`
}

type AsyncReq struct {
}

type SignedPacket struct {
	Signature string `json:"signature"`
}
