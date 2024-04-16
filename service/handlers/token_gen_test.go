package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGenerateRandomToken(t *testing.T) {
	length := 32
	token, err := GenerateRandomToken(length)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(token) != base64.RawURLEncoding.EncodedLen(sha256.Size) {
		t.Errorf("Generated token has incorrect length")
	}
}
