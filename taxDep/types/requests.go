package types

type SyncReq struct {
	SignedPacket
	Packet *RequestPacket `json:"packet"`
}

type AsyncReq struct {
	SignedPacket
	Packets []RequestPacket `json:"packets"`
}

type SignedPacket struct {
	Signature string `json:"signature"`
}
