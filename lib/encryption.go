package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// CipherEncrypt for encrypt data with AES algorithm
func CipherEncrypt(plaintext, key string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err == nil {
		gcm, err := cipher.NewGCM(c)
		if err == nil {
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
				return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
			}
		}
	}

	return nil, err
}

// CipherDecrypt for decrypt data with AES algorithm
func CipherDecrypt(ciphertext []byte, key string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err == nil {
		gcm, err := cipher.NewGCM(c)
		if err == nil {
			nonceSize := gcm.NonceSize()
			if len(ciphertext) < nonceSize {
				return nil, errors.New("ciphertext too short")
			}
			nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
			return gcm.Open(nil, nonce, ciphertext, nil)
		}
	}
	return nil, err
}
