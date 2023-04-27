package terminal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"sync"
	"tax-management/pkg"
	"time"

	"tax-management/taxDep/transfer"
	"tax-management/taxDep/types"
)

type Terminal struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	transferAPI *transfer.Transfer

	token    string
	clientID string
	exp      time.Time
	mtx      sync.Mutex
}

func New(opt types.TerminalOptions, httpClientExtension pkg.ClientLoggerExtension) (*Terminal, error) {
	tr, err := transfer.NewApiTransfer(transfer.DefaultAPIConfig(opt.ClientID, opt.TerminalBaseURl), httpClientExtension)
	if err != nil {
		return nil, err
	}

	return &Terminal{
		// PrivateKey:  prv,
		// PublicKey:   pub,
		clientID:    opt.ClientID,
		transferAPI: tr,
	}, nil
}

func getPrivateKey(pvPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	prvPemBytes, err := os.ReadFile(pvPath)
	if err != nil {
		return nil, nil, err
	}

	prvBlock, _ := pem.Decode(prvPemBytes)
	if prvBlock == nil {
		return nil, nil, errors.New("invalid trip private key")
	}

	prv, err := x509.ParsePKCS8PrivateKey(prvBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey := prv.(*rsa.PrivateKey)

	return privateKey, &privateKey.PublicKey, nil
}

func (t *Terminal) buildRequestPacket(data any, packetType string, uuid string) *types.RequestPacket {
	return &types.RequestPacket{
		UID:        uuid,
		PacketType: packetType,
		Retry:      false,
		Data:       data,
		FiscalId:   t.clientID,
	}
}
