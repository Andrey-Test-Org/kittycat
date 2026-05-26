package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	apiKeyBytes = 32
	idBytes     = 12
)

func NewAPIKey() (string, error) {
	buf := make([]byte, apiKeyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read crypto/rand: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func NewID() string {
	buf := make([]byte, idBytes)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("crypto/rand unavailable: %w", err))
	}
	return hex.EncodeToString(buf)
}
