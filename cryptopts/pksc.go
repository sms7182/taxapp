package cryptops

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func SignPKCS1v15(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	digest := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
}

func VerifyPKCS1v15(signature []byte, plaintext []byte, publicKey *rsa.PublicKey) (bool, error) {
	hashed := sha256.Sum256(plaintext)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, err
	}
	return true, nil
}
