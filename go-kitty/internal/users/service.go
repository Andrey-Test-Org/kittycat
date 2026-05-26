package users

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-kitty/internal/token"
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Register(ctx context.Context, email string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return User{}, fmt.Errorf("register %q: %w", email, ErrInvalidEmail)
	}

	apiKey, err := token.NewAPIKey()
	if err != nil {
		return User{}, fmt.Errorf("register: generate API key: %w", err)
	}

	u := User{
		ID:        token.NewID(),
		Email:     email,
		APIKey:    apiKey,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return User{}, fmt.Errorf("register: persist user: %w", err)
	}

	s.logger.Info("user registered", "userID", u.ID, "email", u.Email)
	return u, nil
}

func (s *Service) Get(ctx context.Context, id string) (User, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}
