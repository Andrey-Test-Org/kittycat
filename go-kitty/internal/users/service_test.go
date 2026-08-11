package users

import (
	"context"
	"errors"
	"testing"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewInMemoryRepository())
}

func TestService_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{name: "valid", email: "luna@example.com", wantErr: nil},
		{name: "missing at", email: "lunaexample.com", wantErr: ErrInvalidEmail},
		{name: "blank", email: "   ", wantErr: ErrInvalidEmail},
		{name: "too long", email: longEmail(), wantErr: ErrInvalidEmail},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := newTestService(t)
			_, err := svc.Register(context.Background(), tc.email)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("Register: unexpected error: %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Register: want %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	created, err := svc.Register(context.Background(), "mochi@example.com")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{name: "found", id: created.ID, wantErr: nil},
		{name: "missing", id: "deadbeefdeadbeefdeadbeef", wantErr: ErrNotFound},
		{name: "invalid id empty", id: "", wantErr: ErrInvalidID},
		{name: "invalid id non-hex", id: "not-a-hex-id", wantErr: ErrInvalidID},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.Get(context.Background(), tc.id)
			if tc.wantErr == nil && err != nil {
				t.Fatalf("Get: unexpected error: %v", err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("Get: want %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func longEmail() string {
	prefix := make([]byte, maxEmailLength)
	for i := range prefix {
		prefix[i] = 'a'
	}
	return string(prefix) + "@example.com"
}
