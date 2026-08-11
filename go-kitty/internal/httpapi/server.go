// Package httpapi exposes the user Service over HTTP.
package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-kitty/internal/users"
)

// maxBodyBytes caps incoming JSON request bodies.
const maxBodyBytes = 1 << 20 // 1 MiB

// Server is the HTTP transport for the user Service.
type Server struct {
	*http.Server
	svc    *users.Service
	logger *slog.Logger
}

// NewServer builds an HTTP server bound to addr. The logger is used by the
// transport layer (request log, response encode errors) only; the user
// Service itself never logs.
func NewServer(addr string, svc *users.Service, logger *slog.Logger) *Server {
	s := &Server{svc: svc, logger: logger}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", s.handleRegister)
	mux.HandleFunc("GET /users/{id}", s.handleGet)
	mux.HandleFunc("GET /users", s.handleList)

	s.Server = &http.Server{
		Addr:              addr,
		Handler:           WithRequestLogging(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

type registerRequest struct {
	Email string `json:"email"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req registerRequest
	if err := dec.Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("decode register: %w", err))
		return
	}
	if extra := dec.Decode(&struct{}{}); !errors.Is(extra, io.EOF) {
		s.writeError(w, http.StatusBadRequest, errors.New("request body must contain a single JSON object"))
		return
	}

	u, err := s.svc.Register(r.Context(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrInvalidEmail):
			s.writeError(w, http.StatusBadRequest, err)
		case errors.Is(err, users.ErrAlreadyExists):
			s.writeError(w, http.StatusConflict, err)
		default:
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	// Transport-layer logging (not library logging): log non-PII identifiers only.
	s.logger.Info("user registered", "userID", u.ID)
	s.writeJSON(w, http.StatusCreated, u)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u, err := s.svc.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrInvalidID):
			s.writeError(w, http.StatusBadRequest, err)
		case errors.Is(err, users.ErrNotFound):
			s.writeError(w, http.StatusNotFound, err)
		default:
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	s.writeJSON(w, http.StatusOK, u)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	all, err := s.svc.List(r.Context())
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}
	s.writeJSON(w, http.StatusOK, all)
}

func (s *Server) writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.logger.Error("encode response", "err", err)
	}
}

func (s *Server) writeError(w http.ResponseWriter, code int, err error) {
	s.logger.Warn("request failed", "status", code, "err", err)
	s.writeJSON(w, code, map[string]string{"error": err.Error()})
}
