package transfer

import (
	"crypto/rsa"
	"tax-management/pkg"
)

type Transfer struct {
	cfg          *ApiConfig
	serverPubKey *rsa.PublicKey
	//pubKeyID         string
	HttpClientLogger pkg.ClientLoggerExtension
}

func NewApiTransfer(cfg *ApiConfig, httpClientLogger pkg.ClientLoggerExtension) (*Transfer, error) {
	transfer := &Transfer{
		cfg:              cfg,
		HttpClientLogger: httpClientLogger,
	}

	return transfer, nil
}

func (t *Transfer) setServerInfos() error {
	// pubkey, id, err := t.GetServerPublicKey()
	// if err != nil {
	// 	return err
	// }

	// //t.pubKeyID = id
	// t.serverPubKey = pubkey
	return nil
}
