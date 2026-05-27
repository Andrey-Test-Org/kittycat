package token

import "testing"

func TestNewAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{name: "length is 64 hex chars (32 bytes)", want: 64},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			key, err := NewAPIKey()
			if err != nil {
				t.Fatalf("NewAPIKey: %v", err)
			}
			if len(key) != tc.want {
				t.Fatalf("NewAPIKey length: want %d, got %d", tc.want, len(key))
			}
		})
	}
}

func TestNewAPIKey_Unique(t *testing.T) {
	t.Parallel()

	seen := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		key, err := NewAPIKey()
		if err != nil {
			t.Fatalf("NewAPIKey: %v", err)
		}
		if _, dup := seen[key]; dup {
			t.Fatalf("duplicate API key after %d iterations: %s", i, key)
		}
		seen[key] = struct{}{}
	}
}

func TestNewID(t *testing.T) {
	t.Parallel()

	id, err := NewID()
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	if len(id) != 24 {
		t.Fatalf("NewID length: want 24, got %d", len(id))
	}
}
