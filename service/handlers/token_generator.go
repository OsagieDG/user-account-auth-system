package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func GenerateRandomToken(length int) (string, error) {
	if length < 32 {
		return "", errors.New("token length is too short")
	}

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Hashing the random bytes with SHA-256
	hashedBytes := sha256.Sum256(randomBytes)

	// Encoding the hashed bytes as a base64 string
	token := base64.RawURLEncoding.EncodeToString(hashedBytes[:])

	return token, nil
}
