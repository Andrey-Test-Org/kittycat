package inventory

import (
	"context"
	"fmt"
	"log/slog"
)

type Handler struct {
	store  *Store
	logger *slog.Logger
}

func NewHandler(store *Store, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) Restock(ctx context.Context, id string, delta int) (qty int, err error) {
	it, err := h.store.Get(ctx, id)
	if err != nil {
		err = fmt.Errorf("restock %s: %w", id, err)
		return
	}
	it.Quantity += delta
	if err = h.store.Upsert(ctx, it); err != nil {
		err = fmt.Errorf("restock %s: persist: %w", id, err)
		return
	}
	qty = it.Quantity
	h.logger.Info("restocked", "itemID", id, "delta", delta, "quantity", qty)
	return
}
