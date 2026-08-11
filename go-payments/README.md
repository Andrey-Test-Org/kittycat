# go-payments

Small payments service: customers, payments, audit log, token issuance.

## Layout

```
go-payments/
├── cmd/payments/      # main entrypoint
└── internal/
    ├── payments/      # domain
    ├── customer/      # customer registry
    ├── token/         # token generation
    ├── httpapi/       # HTTP server
    ├── audit/         # audit log
    ├── config/        # runtime config
    └── concurrency/   # batch helpers
```

## Run

```
cd go-payments
go run ./cmd/payments
```
