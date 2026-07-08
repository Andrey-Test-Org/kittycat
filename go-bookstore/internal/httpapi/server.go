// Package httpapi exposes the bookstore services over HTTP.
package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/author"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/book"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/cart"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/inventory"
	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/order"
)

// Dependencies aggregates the services the HTTP layer wires up.
type Dependencies struct {
	Books     *book.Service
	Authors   *author.Service
	Inventory *inventory.Service
	Carts     *cart.Service
	Orders    *order.Service
	Logger    *slog.Logger
}

// Server is the HTTP transport.
type Server struct {
	*http.Server
	deps Dependencies
}

// NewServer constructs a Server bound to addr.
func NewServer(addr string, deps Dependencies) *Server {
	s := &Server{deps: deps}
	mux := http.NewServeMux()
	s.registerBookRoutes(mux)
	s.registerAuthorRoutes(mux)
	s.registerOrderRoutes(mux)
	s.registerCartRoutes(mux)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true}, deps.Logger)
	})

	s.Server = &http.Server{
		Addr:              addr,
		Handler:           WithRequestLogging(deps.Logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

func writeJSON(w http.ResponseWriter, code int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Error("encode response", "err", err)
	}
}

func writeError(w http.ResponseWriter, code int, err error, logger *slog.Logger) {
	logger.Warn("request failed", "status", code, "err", err)
	writeJSON(w, code, map[string]string{"error": err.Error()}, logger)
}
