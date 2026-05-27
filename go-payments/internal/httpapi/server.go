package httpapi

import (
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
	logger *slog.Logger
}

func NewServer(addr string, logger *slog.Logger) *Server {
	s := &Server{logger: logger}
	mux := http.NewServeMux()
	registerRoutes(mux, logger)
	s.Server = &http.Server{
		Addr:              addr,
		Handler:           WithRequestLogging(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}
