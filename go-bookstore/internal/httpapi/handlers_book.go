package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/book"
)

const maxBookBodyBytes = 1 << 20

func (s *Server) registerBookRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /books", s.handleBookCreate)
	mux.HandleFunc("GET /books/{id}", s.handleBookGet)
	mux.HandleFunc("PATCH /books/{id}", s.handleBookUpdate)
	mux.HandleFunc("DELETE /books/{id}", s.handleBookDelete)
	mux.HandleFunc("GET /books", s.handleBookList)
	mux.HandleFunc("GET /books/search", s.handleBookSearch)
}

type bookCreateRequest struct {
	ISBN        string `json:"isbn"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	AuthorID    string `json:"authorId"`
	PriceCents  int64  `json:"priceCents"`
	Currency    string `json:"currency"`
	Genre       string `json:"genre"`
	PageCount   int    `json:"pageCount"`
	Description string `json:"description"`
}

func (s *Server) handleBookCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBookBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req bookCreateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode book create: %w", err), s.deps.Logger)
		return
	}
	if extra := dec.Decode(&struct{}{}); !errors.Is(extra, io.EOF) {
		writeError(w, http.StatusBadRequest, errors.New("body must contain a single JSON object"), s.deps.Logger)
		return
	}

	b, err := s.deps.Books.Create(r.Context(), book.CreateInput{
		ISBN:        req.ISBN,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		AuthorID:    req.AuthorID,
		PriceCents:  req.PriceCents,
		Currency:    req.Currency,
		Genre:       req.Genre,
		PageCount:   req.PageCount,
		Description: req.Description,
	})
	if err != nil {
		writeBookError(w, err, s.deps.Logger)
		return
	}
	s.deps.Logger.Info("book created", "bookID", b.ID)
	writeJSON(w, http.StatusCreated, b, s.deps.Logger)
}

func (s *Server) handleBookGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.deps.Books.Get(r.Context(), id)
	if err != nil {
		writeBookError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, b, s.deps.Logger)
}

type bookUpdateRequest struct {
	Title       *string `json:"title,omitempty"`
	Subtitle    *string `json:"subtitle,omitempty"`
	PriceCents  *int64  `json:"priceCents,omitempty"`
	Currency    *string `json:"currency,omitempty"`
	Genre       *string `json:"genre,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *Server) handleBookUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	r.Body = http.MaxBytesReader(w, r.Body, maxBookBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req bookUpdateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode book update: %w", err), s.deps.Logger)
		return
	}
	b, err := s.deps.Books.Update(r.Context(), id, book.UpdateInput{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		PriceCents:  req.PriceCents,
		Currency:    req.Currency,
		Genre:       req.Genre,
		Description: req.Description,
	})
	if err != nil {
		writeBookError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, b, s.deps.Logger)
}

func (s *Server) handleBookDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.deps.Books.Delete(r.Context(), id); err != nil {
		writeBookError(w, err, s.deps.Logger)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleBookList(w http.ResponseWriter, r *http.Request) {
	offset, limit := pagination(r, 25)
	books, err := s.deps.Books.List(r.Context(), offset, limit)
	if err != nil {
		writeBookError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, books, s.deps.Logger)
}

func (s *Server) handleBookSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	_, limit := pagination(r, 25)
	books, err := s.deps.Books.Search(r.Context(), q, limit)
	if err != nil {
		writeBookError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, books, s.deps.Logger)
}

func writeBookError(w http.ResponseWriter, err error, logger *slog.Logger) {
	switch {
	case errors.Is(err, book.ErrNotFound):
		writeError(w, http.StatusNotFound, err, logger)
	case errors.Is(err, book.ErrInvalidTitle),
		errors.Is(err, book.ErrInvalidISBN),
		errors.Is(err, book.ErrInvalidPrice),
		errors.Is(err, book.ErrInvalidCurrency),
		errors.Is(err, book.ErrInvalidPageCount),
		errors.Is(err, book.ErrInvalidGenre),
		errors.Is(err, book.ErrInvalidQuery):
		writeError(w, http.StatusBadRequest, err, logger)
	case errors.Is(err, book.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err, logger)
	default:
		writeError(w, http.StatusInternalServerError, err, logger)
	}
}

func pagination(r *http.Request, defaultLimit int) (offset, limit int) {
	limit = defaultLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return offset, limit
}
