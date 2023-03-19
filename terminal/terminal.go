package terminal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"tax-management/transfer"
	"tax-management/types"
	"time"

	"github.com/google/uuid"
)

type Terminal struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	TransferAPI *transfer.Transfer

	token    string
	clientID string
	exp      time.Time
}

func New(opt types.TerminalOptions) (*Terminal, error) {
	prv, pub, err := getPrivateKey(opt.PrivatePemPath)
	if err != nil {
		return nil, err
	}

	return &Terminal{
		PrivateKey:  prv,
		PublicKey:   pub,
		clientID:    opt.ClientID,
		TransferAPI: transfer.NewApiTransfer(transfer.DefaultAPIConfig(prv, pub)),
	}, nil
}

func getPrivateKey(pvPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	prvPemBytes, err := os.ReadFile(pvPath)
	if err != nil {
		return nil, nil, err
	}

	prvBlock, _ := pem.Decode(prvPemBytes)
	if prvBlock == nil {
		return nil, nil, errors.New("invalid kitchen private key")
	}

	prv, err := x509.ParsePKCS8PrivateKey(prvBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey := prv.(*rsa.PrivateKey)

	return privateKey, &privateKey.PublicKey, nil
}

func (t *Terminal) BuildRequestPacket(data any, packetType string) *types.RequestPacket {
	uid := uuid.NewString()
	return &types.RequestPacket{
		UID:        uid,
		PacketType: packetType,
		Retry:      false,
		Data:       data,
		FiscalId:   t.clientID,
	}
}
