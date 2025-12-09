package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HttpServerConfig struct {
	Host         string `env:"HTTP_HOST" env-default:"localhost" `
	Port         string `env:"HTTP_PORT" env-default:"8080" `
	ReadTimeout  string `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout string `env:"HTTP_WRITE_TIMEOUT" env-default:"10s" `
}

type Server struct {
	server *http.Server
}

func NewHttpServer(handler http.Handler, c HttpServerConfig) (*Server, error) {
	readTimeout, err := time.ParseDuration(c.ReadTimeout)
	if err != nil {
		return nil, fmt.Errorf("incorrect format HTTP_READ_TIMEOUT or default: %v", err)
	}

	writeTimeout, err := time.ParseDuration(c.WriteTimeout)
	if err != nil {
		return nil, fmt.Errorf("incorrect format HTTP_WRITE_TIMEOUT or default: %v", err)
	}

	srv := &http.Server{
		Addr:         c.Host + ":" + c.Port,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return &Server{
		server: srv,
	}, nil
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error during graceful shutdown: %w", err)
	}
	return nil
}

func (s *Server) Close() error {
	if err := s.server.Close(); err != nil {
		return fmt.Errorf("error when closing the server: %w", err)
	}
	return nil
}
