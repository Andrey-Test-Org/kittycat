# go-shipping

Carrier integration: fetch manifests, persist tracker state, notify hooks.

## Files

- `shipping.go` — HTTP client for carrier manifests
- `tracker.go` — local JSON state persistence
- `notifier.go` — outbound webhook notifications
