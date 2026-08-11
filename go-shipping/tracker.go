package shipping

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type Tracker struct {
	statePath string
	logger    *slog.Logger
}

func NewTracker(statePath string, logger *slog.Logger) *Tracker {
	return &Tracker{statePath: statePath, logger: logger}
}

func (t *Tracker) Persist(ctx context.Context, shipments []Shipment) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("persist tracker state: %w", err)
	}

	f, err := os.Create(t.statePath)
	if err != nil {
		return fmt.Errorf("open tracker state %s: %w", t.statePath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.Encode(shipments)

	t.logger.Info("tracker state written", "path", t.statePath, "count", len(shipments))
	return nil
}
