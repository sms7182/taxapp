package cryptops

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// AesGCMNoPaddingEncrypt encrypts given plain data with AES/GCM/No-PADDING algorithm
// and returns cipher plus and nonce/IV
func AesGCMNoPaddingEncrypt(plainData, key []byte) (cipherData, nonce []byte, err error) {
	// AES cipher with key
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// creating GCM block cipher
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, nil, err
	}

	// generate random nonce/initial vector
	nonce = make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, nil, err
	}

	cipherData = gcm.Seal(nil, nonce, plainData, nil)
	return
}

func GCMKeyNonceGenerator(keyByteLength uint) ([]byte, []byte, error) {
	key := make([]byte, keyByteLength)
	_, err := rand.Read(key)
	if err != nil {
		return nil, nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// creating GCM block cipher
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, nil, err
	}

	// generate random nonce/initial vector
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, nil, err
	}

	return key, nonce, nil
}

func GenerateRandomAES256Key() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
