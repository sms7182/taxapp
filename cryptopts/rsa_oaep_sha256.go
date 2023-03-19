package cryptops

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func EncryptRsaOAEPSHA256(plainData []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plainData, nil)
}
