package transfer

import (
	"crypto/rsa"
	cryptops "tax-management/cryptopts"
)

type PriorityLevel uint8

const (
	PriorityLevelNormal PriorityLevel = iota + 1
	PriorityLevelHigh
)

type ApiConfig struct {
	baseURL string
	prvKey  *rsa.PrivateKey
	pubKey  *rsa.PublicKey

	normalizer func(map[string]interface{}) (string, error)

	signer func([]byte, *rsa.PrivateKey) ([]byte, error)

	encrypter func(plainData, key []byte) (cipherData, nonce []byte, err error)
}

func DefaultAPIConfig(prvKey *rsa.PrivateKey, pubKey *rsa.PublicKey) *ApiConfig {
	return &ApiConfig{
		baseURL: "https://tp.tax.gov.ir/req/api/self-tsp",
		//baseURL:    "https://sandboxrc.tax.gov.ir/req/api/self-tsp",
		prvKey:     prvKey,
		pubKey:     pubKey,
		normalizer: cryptops.NormalizeJsonObj,
		signer:     cryptops.SignPKCS1v15,
		encrypter:  cryptops.AesGCMNoPaddingEncrypt,
	}
}
