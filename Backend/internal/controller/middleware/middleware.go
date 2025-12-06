package middleware

import (
	"bytes"
	"net/http"
	"runtime/debug"

	"github.com/PlopyBlopy/notebot/config"
	"github.com/rs/zerolog/log"
)

// Type for middleware
type Middleware func(h http.Handler) http.Handler

// A struct for implementing Fluent-Builder, middleware storage
type MiddlewareBuilder struct {
	middlewares []Middleware
}

// Create builder
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

// Adds error handling for each request
func (m *MiddlewareBuilder) GlobalExceptionHandler(c config.Config) *MiddlewareBuilder {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseInterceptor{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
			}

			next.ServeHTTP(rw, r)

			if rw.statusCode >= 400 {
				if c.Isdev {
					log.Error().
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Int("status", rw.statusCode).
						Str("response_body", rw.body.String()).
						Msg("HTTP error")
				} else {
					log.Error().
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Int("status", rw.statusCode).
						Msg("HTTP error")
				}
			}
		})
	}
	m.middlewares = append(m.middlewares, middleware)
	return m
}

// Processing a panic to convert to an error
func (m *MiddlewareBuilder) PanicMiddleware() *MiddlewareBuilder {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error().Msgf("PANIC recovered: %v", err)
					stack := debug.Stack()
					log.Error().Msgf("stack trace: %s", stack)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
	m.middlewares = append(m.middlewares, middleware)
	return m
}

// For custom middleware
func (m *MiddlewareBuilder) Use(middleware Middleware) *MiddlewareBuilder {
	m.middlewares = append(m.middlewares, middleware)
	return m
}

// The call middleware: Mn->M1->M0 : M0->M1->Mn
func (m *MiddlewareBuilder) Build(h http.Handler) http.Handler {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		h = m.middlewares[i](h)
	}
	return h
}
