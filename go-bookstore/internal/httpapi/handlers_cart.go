package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const maxCartBodyBytes = 64 << 10

func (s *Server) registerCartRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /carts", s.handleCartCreate)
	mux.HandleFunc("GET /carts/{id}", s.handleCartGet)
	mux.HandleFunc("POST /carts/{id}/items", s.handleCartAddItem)
	mux.HandleFunc("DELETE /carts/{id}/items/{bookId}", s.handleCartRemoveItem)
	mux.HandleFunc("DELETE /carts/{id}", s.handleCartClear)
}

type cartCreateRequest struct {
	CustomerID string `json:"customerId"`
}

func (s *Server) handleCartCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxCartBodyBytes)
	var req cartCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode cart create: %w", err), s.deps.Logger)
		return
	}
	c, err := s.deps.Carts.Create(r.Context(), req.CustomerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusCreated, c, s.deps.Logger)
}

func (s *Server) handleCartGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.deps.Carts.Get(r.Context(), id)
	if err != nil {
		writeCartError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, c, s.deps.Logger)
}

type cartAddItemRequest struct {
	BookID   string `json:"bookId"`
	Quantity int    `json:"quantity"`
}

func (s *Server) handleCartAddItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	r.Body = http.MaxBytesReader(w, r.Body, maxCartBodyBytes)
	var req cartAddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode add item: %w", err), s.deps.Logger)
		return
	}
	c, err := s.deps.Carts.AddItem(r.Context(), id, req.BookID, req.Quantity)
	if err != nil {
		writeCartError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, c, s.deps.Logger)
}

func (s *Server) handleCartRemoveItem(w http.ResponseWriter, r *http.Request) {
	cid := r.PathValue("id")
	bid := r.PathValue("bookId")
	c, err := s.deps.Carts.RemoveItem(r.Context(), cid, bid)
	if err != nil {
		writeCartError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, c, s.deps.Logger)
}

func (s *Server) handleCartClear(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.deps.Carts.Clear(r.Context(), id); err != nil {
		writeCartError(w, err, s.deps.Logger)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeCartError(w http.ResponseWriter, err error, logger *slog.Logger) {
	if errors.Is(err, errCartNotFound) {
		writeError(w, http.StatusNotFound, err, logger)
		return
	}
	writeError(w, http.StatusInternalServerError, err, logger)
}

// errCartNotFound is a placeholder; cart package exports its own ErrNotFound,
// but we shadow it locally here just to keep import surface narrow.
var errCartNotFound = errors.New("cart not found")
