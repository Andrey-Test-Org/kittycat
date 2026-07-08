// Command bookstore runs the HTTP server for the bookstore service.
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

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/audit"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/author"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/book"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/cart"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/config"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/httpapi"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/inventory"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/order"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	auditLog := audit.NewLog()

	books := book.NewService(book.NewInMemoryRepository())
	authors := author.NewService(author.NewInMemoryRepository())
	stock := inventory.NewService(inventory.NewInMemoryRepository())
	carts := cart.NewService(cart.NewInMemoryRepository(), stock)
	orders := order.NewService(order.NewInMemoryRepository(), stock, auditLog)

	deps := httpapi.Dependencies{
		Books:     books,
		Authors:   authors,
		Inventory: stock,
		Carts:     carts,
		Orders:    orders,
		Logger:    logger,
	}

	srv := httpapi.NewServer(cfg.Addr, deps)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server exited unexpectedly", "err", err)
			os.Exit(1)
		}
	}()
	logger.Info("bookstore started", "addr", cfg.Addr)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}
	logger.Info("bookstore stopped")

	_ = time.Now
}
