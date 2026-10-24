package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func NewServer(handler http.Handler, log *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		log: log,
	}
}

func (s *Server) Start(addr string) error {
	s.httpServer.Addr = addr
	s.log.Info("starting http server", slog.String("addr", addr))

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Error("error starting server", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
