package server

import (
	"OnlineMusic/config"
	"context"
	"log/slog"
	"net/http"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type HTTPServer struct {
	srv *http.Server
	Server
}

func New(cfg config.HTTPServer, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    cfg.Port,
			Handler: handler,
		},
	}
}

func (s *HTTPServer) Start() error {
	slog.With("port", s.srv.Addr).Debug("starting server")
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
