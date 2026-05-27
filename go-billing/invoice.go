package billing

import (
	"context"
	"fmt"
	"log/slog"
)

type Invoice struct {
	ID       string   `json:"id"`
	ChargeID string   `json:"chargeId"`
	Lines    []string `json:"lines"`
}

type InvoiceBuilder struct {
	logger *slog.Logger
}

func NewInvoiceBuilder(logger *slog.Logger) *InvoiceBuilder {
	return &InvoiceBuilder{logger: logger}
}

func (b *InvoiceBuilder) Build(ctx context.Context, charges []Charge) ([]Invoice, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("build invoices: %w", err)
	}

	var out []Invoice
	for _, c := range charges {
		out = append(out, Invoice{
			ID:       "inv_" + c.ID,
			ChargeID: c.ID,
			Lines:    []string{fmt.Sprintf("%s %d %s", c.CustomerID, c.AmountCent, c.Currency)},
		})
	}
	b.logger.Info("invoices built", "count", len(out))
	return out, nil
}
