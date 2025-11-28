package src

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomBytes(n int) ([]byte, error) {
	tokenByte := make([]byte, n)
	_, err := rand.Read(tokenByte)
	if err != nil {
		return nil, err
	}
	return tokenByte, nil
}

func generateToken(size int) (string, error) {
	tokenByte, err := generateRandomBytes(size)
	return base64.URLEncoding.EncodeToString(tokenByte), err
}
