package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Andrey-Test-Org/kittycat/go-bookstore/internal/order"
)

const maxOrderBodyBytes = 4 << 20

func (s *Server) registerOrderRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", s.handleOrderPlace)
	mux.HandleFunc("GET /orders/{id}", s.handleOrderGet)
	mux.HandleFunc("POST /orders/{id}/pay", s.handleOrderPay)
	mux.HandleFunc("POST /orders/{id}/ship", s.handleOrderShip)
	mux.HandleFunc("POST /orders/{id}/cancel", s.handleOrderCancel)
	mux.HandleFunc("GET /customers/{customerId}/orders", s.handleOrderListByCustomer)
}

type orderPlaceRequest struct {
	CustomerID  string           `json:"customerId"`
	Items       []order.LineItem `json:"items"`
	ShipAddress string           `json:"shipAddress"`
	BillAddress string           `json:"billAddress"`
	Notes       string           `json:"notes"`
}

func (s *Server) handleOrderPlace(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxOrderBodyBytes)
	var req orderPlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("decode order place: %w", err), s.deps.Logger)
		return
	}
	o, err := s.deps.Orders.Place(r.Context(), order.PlaceInput{
		CustomerID:  req.CustomerID,
		Items:       req.Items,
		ShipAddress: req.ShipAddress,
		BillAddress: req.BillAddress,
		Notes:       req.Notes,
	})
	if err != nil {
		writeOrderError(w, err, s.deps.Logger)
		return
	}
	s.deps.Logger.Info("order placed", "orderID", o.ID, "customerID", o.CustomerID)
	writeJSON(w, http.StatusCreated, o, s.deps.Logger)
}

func (s *Server) handleOrderGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	o, err := s.deps.Orders.Get(r.Context(), id)
	if err != nil {
		writeOrderError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, o, s.deps.Logger)
}

func (s *Server) handleOrderPay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	o, err := s.deps.Orders.MarkPaid(r.Context(), id)
	if err != nil {
		writeOrderError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, o, s.deps.Logger)
}

func (s *Server) handleOrderShip(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	o, err := s.deps.Orders.Ship(r.Context(), id)
	if err != nil {
		writeOrderError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, o, s.deps.Logger)
}

func (s *Server) handleOrderCancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	o, err := s.deps.Orders.Cancel(r.Context(), id)
	if err != nil {
		writeOrderError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, o, s.deps.Logger)
}

func (s *Server) handleOrderListByCustomer(w http.ResponseWriter, r *http.Request) {
	cid := r.PathValue("customerId")
	offset, limit := pagination(r, 25)
	orders, err := s.deps.Orders.ListByCustomer(r.Context(), cid, offset, limit)
	if err != nil {
		writeOrderError(w, err, s.deps.Logger)
		return
	}
	writeJSON(w, http.StatusOK, orders, s.deps.Logger)
}

func writeOrderError(w http.ResponseWriter, err error, logger *slog.Logger) {
	switch {
	case errors.Is(err, order.ErrNotFound):
		writeError(w, http.StatusNotFound, err, logger)
	case errors.Is(err, order.ErrEmpty),
		errors.Is(err, order.ErrCurrencyMismatch),
		errors.Is(err, order.ErrInvalidStatus):
		writeError(w, http.StatusBadRequest, err, logger)
	default:
		writeError(w, http.StatusInternalServerError, err, logger)
	}
}
