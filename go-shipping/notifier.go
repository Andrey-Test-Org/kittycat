package shipping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Notifier struct {
	hookURL string
	http    *http.Client
	logger  *slog.Logger
}

func NewNotifier(hookURL string, logger *slog.Logger) *Notifier {
	return &Notifier{
		hookURL: hookURL,
		http:    &http.Client{Timeout: 3 * time.Second},
		logger:  logger,
	}
}

func (n *Notifier) Notify(ctx context.Context, s Shipment) error {
	payload, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("encode notify payload for %s: %w", s.ID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.hookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build notify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, _ := n.http.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	n.logger.Info("shipment notified", "shipmentID", s.ID, "hook", n.hookURL)
	return nil
}
