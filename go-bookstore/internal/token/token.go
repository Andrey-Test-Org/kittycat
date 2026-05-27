// Package token issues cryptographically random identifiers.
package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	idBytes     = 12
	apiKeyBytes = 32
)

// NewID returns a hex-encoded 12-byte random identifier sourced from crypto/rand.
func NewID() (string, error) {
	buf := make([]byte, idBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read crypto/rand: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// NewAPIKey returns a hex-encoded 32-byte random API key sourced from crypto/rand.
func NewAPIKey() (string, error) {
	buf := make([]byte, apiKeyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read crypto/rand: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
