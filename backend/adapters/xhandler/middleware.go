package xhandler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/zyghq/zyg/ports"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default
		}
		next.ServeHTTP(wrapper, r)
		slog.Info(
			"http",
			slog.Int("status", wrapper.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("duration", time.Since(start).String()),
		)
	})
}

type EnsureAuth struct {
	handler AuthenticatedHandler
	authz   ports.CustomerAuthServicer
}

func NewEnsureAuth(handler AuthenticatedHandler, as ports.CustomerAuthServicer) *EnsureAuth {
	return &EnsureAuth{
		handler: handler,
		authz:   as,
	}
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	scheme, cred, err := CheckAuthCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	customer, err := AuthenticateCustomer(r.Context(), ea.authz, scheme, cred)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	ea.handler(w, r, &customer)
}
