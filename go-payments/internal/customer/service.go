package customer

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type Service struct {
	repo   *Repository
	logger *slog.Logger
}

func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Register(ctx context.Context, id, email string) (Customer, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return Customer{}, fmt.Errorf("register customer: invalid email %q", email)
	}
	c := Customer{ID: id, Email: email, CreatedAt: time.Now().UTC()}
	if err := s.repo.Upsert(ctx, c); err != nil {
		return Customer{}, fmt.Errorf("register customer: %w", err)
	}
	s.logger.Info("customer registered", "customerID", id)
	return c, nil
}
