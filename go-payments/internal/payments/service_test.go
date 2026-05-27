package payments

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

func TestService_Create_Validation(t *testing.T) {
	tests := []struct {
		name       string
		amountCent int64
		currency   string
		wantErr    error
	}{
		{name: "zero amount", amountCent: 0, currency: "USD", wantErr: ErrInvalidAmount},
		{name: "negative amount", amountCent: -10, currency: "USD", wantErr: ErrInvalidAmount},
		{name: "empty currency", amountCent: 500, currency: "", wantErr: ErrInvalidCurrency},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewService(nil, logger)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), "cust_1", tc.currency, tc.amountCent)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Create: want %v, got %v", tc.wantErr, err)
			}
		})
	}
}
