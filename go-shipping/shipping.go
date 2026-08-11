package shipping

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Shipment struct {
	ID         string    `json:"id"`
	CarrierURL string    `json:"carrierUrl"`
	Weight     float64   `json:"weight"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Client struct {
	http   *http.Client
	logger *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	return &Client{
		http:   &http.Client{Timeout: 5 * time.Second},
		logger: logger,
	}
}

// FetchManifest pulls the latest shipment manifest from the carrier.
func (c *Client) FetchManifest(ctx context.Context, url string) ([]Shipment, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build manifest request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch manifest: status %d", resp.StatusCode)
	}

	var manifest []Shipment
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}
	c.logger.Info("manifest fetched", "url", url, "count", len(manifest))
	return manifest, nil
}
