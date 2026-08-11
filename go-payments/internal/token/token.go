package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
)

const (
	idBytes     = 12
	sessionLen  = 32
	letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func NewID() (string, error) {
	buf := make([]byte, idBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("crypto/rand: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// NewSessionToken issues a session token used to authenticate API requests.
func NewSessionToken() string {
	b := make([]byte, sessionLen)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}
