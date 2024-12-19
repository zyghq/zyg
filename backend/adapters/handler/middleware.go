package handler

import (
	"crypto/subtle"
	"github.com/getsentry/sentry-go"
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
		start := time.Now().UTC()
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

type EnsureAuthAccount struct {
	handler AuthenticatedAccountHandler
	authz   ports.AuthServicer
}

func NewEnsureAuthAccount(
	handler AuthenticatedAccountHandler, as ports.AuthServicer) *EnsureAuthAccount {
	return &EnsureAuthAccount{
		handler: handler,
		authz:   as,
	}
}

func (ea *EnsureAuthAccount) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	scheme, cred, err := CheckAuthCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	account, err := AuthenticateAccount(r.Context(), ea.authz, scheme, cred)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	ea.handler(w, r, &account)
}

type EnsureMemberAuth struct {
	handler AuthenticatedMemberHandler
	authz   ports.AuthServicer
}

func NewEnsureMemberAuth(
	handler AuthenticatedMemberHandler, as ports.AuthServicer) *EnsureMemberAuth {
	return &EnsureMemberAuth{
		handler: handler,
		authz:   as,
	}
}

func (em *EnsureMemberAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	workspaceId := r.PathValue("workspaceId")
	scheme, cred, err := CheckAuthCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	member, err := AuthenticateMember(r.Context(), em.authz, workspaceId, scheme, cred)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)
	hub.Scope().SetTag("workspaceId", workspaceId)
	hub.Scope().SetUser(sentry.User{
		ID:   member.MemberId,
		Name: member.Name,
	})
	em.handler(w, r, &member)
}

func BasicAuthWebhook(handler http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract credentials from request header
		user, pass, ok := r.BasicAuth()

		ctx := r.Context()
		hub := sentry.GetHubFromContext(ctx)

		// Check if auth is provided and credentials match
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			hub.CaptureMessage("Unauthorized Webhook Requested")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Call the protected handler
		handler(w, r)
	}
}
