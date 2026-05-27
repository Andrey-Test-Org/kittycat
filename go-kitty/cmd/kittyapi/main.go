package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-kitty/internal/httpapi"
	"github.com/Andrey-Test-Org/kittycat/go-kitty/internal/users"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svc := users.NewService(users.NewInMemoryRepository())
	srv := httpapi.NewServer(":8080", svc, logger)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server exited unexpectedly", "err", err)
			os.Exit(1)
		}
	}()
	logger.Info("kittyapi started", "addr", srv.Addr)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}
	logger.Info("kittyapi stopped")
}
