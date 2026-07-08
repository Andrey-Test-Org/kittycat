package token

import "testing"

func TestNewID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{name: "24 hex chars", want: 24},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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

func TestNewAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{name: "64 hex chars", want: 64},
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
				t.Fatalf("length: want %d, got %d", tc.want, len(key))
			}
		})
	}
}

func TestNewID_Unique(t *testing.T) {
	t.Parallel()
	seen := make(map[string]struct{}, 200)
	for i := 0; i < 200; i++ {
		id, err := NewID()
		if err != nil {
			t.Fatalf("NewID: %v", err)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate ID after %d iterations: %s", i, id)
		}
		seen[id] = struct{}{}
	}
}
