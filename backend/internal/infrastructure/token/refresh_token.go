package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func HashRefreshToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))

	return hex.EncodeToString(hash[:])
}
