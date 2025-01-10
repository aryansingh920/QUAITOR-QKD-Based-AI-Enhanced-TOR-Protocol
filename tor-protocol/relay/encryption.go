package relay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

var privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

func Encrypt(data []byte) ([]byte, error) {
    return rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateKey.PublicKey, data, nil)
}

func Decrypt(data []byte) ([]byte, error) {
    return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
}
