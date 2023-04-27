package transfer

import (
	"crypto/rsa"

	"tax-management/taxDep/cryptoutils"
)

type PriorityLevel uint8

const (
	PriorityLevelNormal PriorityLevel = iota + 1
	PriorityLevelHigh
)

type ApiConfig struct {
	baseURL  string
	clientID string
	prvKey   *rsa.PrivateKey
	pubKey   *rsa.PublicKey
	public   string
	private  string

	normalizer func(map[string]interface{}) (string, error)

	signer func([]byte, *rsa.PrivateKey) ([]byte, error)

	encrypter func(plainData, key []byte, nonceSize int) (cipherData, nonce []byte, err error)

	pubKeyEncrypter func([]byte, *rsa.PublicKey) ([]byte, error)
}

func DefaultAPIConfig(prvKey *rsa.PrivateKey, pubKey *rsa.PublicKey, clientID, baseURL string) *ApiConfig {
	return &ApiConfig{
		baseURL:         baseURL,
		prvKey:          prvKey,
		pubKey:          pubKey,
		normalizer:      cryptoutils.NormalizeJsonObj,
		signer:          cryptoutils.SignPKCS1v15,
		encrypter:       cryptoutils.AesGCMNoPaddingEncrypt,
		pubKeyEncrypter: cryptoutils.EncryptRsaOAEPSHA256,
		clientID:        clientID,
	}
}
