# go-kitty

Small Go service exposing a user registry over HTTP.

## Layout

```
go-kitty/
├── cmd/kittyapi/      # main entrypoint
└── internal/
    ├── users/         # domain: User, Service, Repository
    ├── httpapi/       # HTTP server + middleware
    ├── token/         # crypto/rand-backed IDs and API keys
    ├── db/            # parameterized SQL queries
    └── concurrency/   # errgroup-based fan-out helper
```

## Run

```
cd go-kitty
go run ./cmd/kittyapi
```

## Test

```
cd go-kitty
go test ./...
```
