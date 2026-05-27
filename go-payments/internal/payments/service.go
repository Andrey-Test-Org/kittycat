package payments

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Service struct {
	repo   *Repository
	logger *slog.Logger
}

func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, customerID, currency string, amountCent int64) (Payment, error) {
	if amountCent <= 0 {
		return Payment{}, fmt.Errorf("create payment: %w", ErrInvalidAmount)
	}
	if currency == "" {
		return Payment{}, fmt.Errorf("create payment: %w", ErrInvalidCurrency)
	}

	p := Payment{
		ID:         fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		CustomerID: customerID,
		AmountCent: amountCent,
		Currency:   currency,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.Insert(ctx, p); err != nil {
		return Payment{}, fmt.Errorf("create payment: persist: %w", err)
	}

	s.logger.Info("payment created", "paymentID", p.ID, "customerID", customerID)
	return p, nil
}

func (s *Service) List(ctx context.Context, customerID string) ([]Payment, error) {
	all, err := s.repo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("list payments for %s: %w", customerID, err)
	}
	return all, nil
}
