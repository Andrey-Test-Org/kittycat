# go-bookstore

A small bookstore service. Manages books, authors, inventory, carts, orders, and payments.

## Layout

```
go-bookstore/
├── cmd/bookstore/      # main entrypoint
└── internal/
    ├── book/           # book domain
    ├── author/         # author domain
    ├── order/          # order domain
    ├── inventory/      # stock tracking
    ├── cart/           # shopping carts
    ├── httpapi/        # HTTP transport
    ├── config/         # runtime configuration
    ├── token/          # crypto-secure ID issuance
    ├── concurrency/    # errgroup-backed helpers
    └── audit/          # audit log
```

## Run

```
cd go-bookstore
go run ./cmd/bookstore
```

## Test

```
go test ./...
```
