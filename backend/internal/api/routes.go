package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/auth"
	"github.com/zyghq/zyg/internal/model"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info(
			"request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("duration", time.Since(start).String()),
		)
	})
}

func HttpAuthCredentials(r *http.Request) (string, string, error) {
	ath := r.Header.Get("Authorization")
	if ath == "" {
		return "", "", fmt.Errorf("no authorization header provided")
	}
	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])
	return scheme, cred[1], nil
}

func AuthenticateAccount(ctx context.Context, db *pgxpool.Pool, r *http.Request) (model.Account, error) {
	var account model.Account

	scheme, cred, err := HttpAuthCredentials(r)
	if err != nil {
		return account, fmt.Errorf("failed to get auth credentials got error: %v", err)
	}

	if scheme == "token" {
		slog.Info("authenticate account with PAT")
		token := model.AccountPAT{Token: cred}
		account, err := token.GetAccountByToken(ctx, db)
		if err != nil {
			return account, fmt.Errorf("failed to authenticate got error: %v", err)
		}
		slog.Info("authenticated account with PAT", slog.String("accountId", account.AccountId))
		return account, nil
	} else if scheme == "bearer" {
		slog.Info("authenticate account with JWT")
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return account, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
		}
		ac, err := auth.ParseJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			return account, fmt.Errorf("failed to parse JWT token got error: %v", err)
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			return account, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		slog.Info("authenticated account with auth user id", slog.String("authUserId", sub))

		// fetch the authenticated account
		account = model.Account{AuthUserId: sub}
		account, err = account.GetByAuthUserId(ctx, db)
		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"account not found or does not exist",
				slog.String("authUserId", sub),
			)
			return account, fmt.Errorf("account not found or does not exist")
		}
		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get account by auth user id "+
					"perhaps a failed query or mapping",
				slog.String("authUserId", sub),
			)
			return account, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		if err != nil {
			slog.Error(
				"failed to get account by auth user id "+
					"something went wrong",
				slog.String("authUserId", sub),
			)
			return account, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		// return the authenticated account
		return account, nil
	} else {
		return account, fmt.Errorf("unsupported scheme `%s` cannot authenticate", scheme)
	}
}

func AuthenticateCustomer(ctx context.Context, db *pgxpool.Pool, r *http.Request) (model.Customer, error) {
	var customer model.Customer

	scheme, cred, err := HttpAuthCredentials(r)
	if err != nil {
		return customer, fmt.Errorf("failed to get auth credentials with error: %v", err)
	}

	if scheme == "bearer" {
		slog.Info("authenticate with customer JWT")
		// TODO: update to specific secret key for customer jwt.
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return customer, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
		}
		cc, err := auth.ParseCustomerJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			return customer, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}
		sub, err := cc.RegisteredClaims.GetSubject()
		if err != nil {
			return customer, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		slog.Info("authenticated customer with customer id", slog.String("customerId", sub))

		// fetch the authenticated customer
		customer = model.Customer{WorkspaceId: cc.WorkspaceId, CustomerId: sub}
		customer, err = customer.GetWrkCustomerById(ctx, db)
		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"customer not found or does not exist",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("customer not found or does not exist")
		}
		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get customer by customer id"+
					"perhaps a failed query or mapping",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("failed to get customer by customer id: %s got error: %v", sub, err)
		}
		if err != nil {
			slog.Error(
				"failed to get customer by customer id"+
					"something went wrong",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("failed to get customer by customer id: %s got error: %v", sub, err)
		}
		// return the authenticated customer
		return customer, nil
	} else {
		return customer, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

func NewHandler(ctx context.Context, db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handleGetIndex)

	mux.Handle("POST /accounts/auth/{$}",
		handleGetOrCreateAuthAccount(ctx, db))
	mux.Handle("POST /pats/{$}",
		handleCreatePAT(ctx, db))
	mux.Handle("GET /pats/{$}",
		handleGetPATs(ctx, db))

	mux.Handle("POST /workspaces/{$}",
		handleCreateWorkspace(ctx, db))

	mux.Handle("GET /workspaces/{$}",
		handleGetWorkspaces(ctx, db))
	mux.Handle("GET /workspaces/{workspaceId}/{$}",
		handleGetWorkspace(ctx, db))

	mux.Handle("POST /workspaces/{workspaceId}/threads/chat/{threadId}/messages/{$}",
		handleCreateMemberThChatMessage(ctx, db))
	mux.Handle("POST /workspaces/{workspaceId}/tokens/{$}",
		handleCustomerTokenIssue(ctx, db))

	mux.Handle("GET /-/me/{$}",
		handleGetCustomer(ctx, db))

	mux.Handle("POST /-/threads/chat/{$}",
		handleInitCustomerThreadChat(ctx, db))
	mux.Handle("GET /-/threads/chat/{$}",
		handleGetCustomerThreadChats(ctx, db))

	mux.Handle("POST /-/threads/chat/{threadId}/messages/{$}",
		handleCreateCustomerThChatMessage(ctx, db))
	mux.Handle("GET /-/threads/chat/{threadId}/messages/{$}",
		handleGetCustomerThChatMessages(ctx, db))

	// ignored
	mux.Handle("POST /-/threads/qa/{$}", handleInitCustomerThreadQA(ctx, db))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler

}
