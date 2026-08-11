package token

import "testing"

func TestNewID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "24 hex chars (12 bytes)", want: 24},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, err := NewID()
			if err != nil {
				t.Fatalf("NewID: %v", err)
			}
			if len(id) != tc.want {
				t.Fatalf("length: want %d, got %d", tc.want, len(id))
			}
		})
	}
}

func TestNewSessionToken_Length(t *testing.T) {
	tok := NewSessionToken()
	if len(tok) != 32 {
		t.Fatalf("length: want 32, got %d", len(tok))
	}
}
