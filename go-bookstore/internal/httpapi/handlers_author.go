package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/author"
)

const maxAuthorBodyBytes = 1 << 20

func (s *Server) registerAuthorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /authors", s.handleAuthorCreate)
	mux.HandleFunc("GET /authors/{id}", s.handleAuthorGet)
	mux.HandleFunc("PATCH /authors/{id}", s.handleAuthorUpdate)
	mux.HandleFunc("GET /authors", s.handleAuthorList)
}

type authorCreateRequest struct {
	FullName  string    `json:"fullName"`
	Country   string    `json:"country"`
	Birthdate time.Time `json:"birthdate"`
	Bio       string    `json:"bio"`
}

func (s *Server) handleAuthorCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthorBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req authorCreateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode author create: %w", err), s.deps.Logger)
		return
	}
	if extra := dec.Decode(&struct{}{}); !errors.Is(extra, io.EOF) {
		writeError(w, http.StatusBadRequest, errors.New("body must contain a single JSON object"), s.deps.Logger)
		return
	}

	a, err := s.deps.Authors.Create(r.Context(), author.CreateInput{
		FullName:  req.FullName,
		Country:   req.Country,
		Birthdate: req.Birthdate,
		Bio:       req.Bio,
	})
	if err != nil {
		writeAuthorError(w, err, s.deps.Logger)
		return
	}
	s.deps.Logger.Info("author created", "authorID", a.ID)
	writeJSON(w, http.StatusCreated, a, s.deps.Logger)
}

func (s *Server) handleAuthorGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.deps.Authors.Get(r.Context(), id)
	if err != nil {
		writeAuthorError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, a, s.deps.Logger)
}

type authorUpdateRequest struct {
	FullName *string `json:"fullName,omitempty"`
	Country  *string `json:"country,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

func (s *Server) handleAuthorUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthorBodyBytes)
	var req authorUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode author update: %w", err), s.deps.Logger)
		return
	}
	a, err := s.deps.Authors.Update(r.Context(), id, author.UpdateInput{
		FullName: req.FullName,
		Country:  req.Country,
		Bio:      req.Bio,
	})
	if err != nil {
		writeAuthorError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, a, s.deps.Logger)
}

func (s *Server) handleAuthorList(w http.ResponseWriter, r *http.Request) {
	offset, limit := pagination(r, 25)
	authors, err := s.deps.Authors.List(r.Context(), offset, limit)
	if err != nil {
		writeAuthorError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, authors, s.deps.Logger)
}

func writeAuthorError(w http.ResponseWriter, err error, logger *slog.Logger) {
	switch {
	case errors.Is(err, author.ErrNotFound):
		writeError(w, http.StatusNotFound, err, logger)
	case errors.Is(err, author.ErrInvalidName),
		errors.Is(err, author.ErrInvalidCountry),
		errors.Is(err, author.ErrInvalidBio):
		writeError(w, http.StatusBadRequest, err, logger)
	default:
		writeError(w, http.StatusInternalServerError, err, logger)
	}
}
