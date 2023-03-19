package transfer

import "crypto/rsa"

type Transfer struct {
	cfg          *ApiConfig
	serverPubKey *rsa.PublicKey
	pubKeyID     string
}

func NewApiTransfer(cfg *ApiConfig) (*Transfer, error) {
	transfer := &Transfer{
		cfg: cfg,
	}

	return transfer, transfer.setServerInfos()
}

func (t *Transfer) setServerInfos() error {
	pubkey, id, err := t.GetServerPublicKey()
	if err != nil {
		return err
	}

	t.pubKeyID = id
	t.serverPubKey = pubkey
	return nil
}
