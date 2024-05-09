package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/auth"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
	"github.com/zyghq/zyg/internal/services"
)

func CheckAuthCredentials(r *http.Request) (string, string, error) {
	ath := r.Header.Get("Authorization")
	if ath == "" {
		return "", "", fmt.Errorf("no authorization header provided")
	}
	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])
	return scheme, cred[1], nil
}

func AuthenticateAccount(
	ctx context.Context, authz ports.AuthServicer,
	scheme string, cred string,
) (domain.Account, error) {
	var account domain.Account
	if scheme == "token" {
		slog.Info("authenticate account with PAT...")
		account, err := authz.GetPatAccount(ctx, cred)
		if err != nil {
			return domain.Account{}, fmt.Errorf("failed to authenticate got error: %v", err)
		}
		slog.Info("authenticated account with PAT", slog.String("accountId", account.AccountId))
		return account, nil
	} else if scheme == "bearer" {
		slog.Info("authenticate account with JWT...")
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return domain.Account{}, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
		}
		ac, err := auth.ParseJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			return domain.Account{}, fmt.Errorf("failed to parse JWT token got error: %v", err)
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			return account, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		slog.Info("authenticated account", slog.String("authUserId", sub))

		account, err = authz.GetAuthUser(ctx, sub)

		if errors.Is(err, services.ErrAccountNotFound) {
			slog.Warn(
				"account not found or does not exist",
				slog.String("authUserId", sub),
			)
			return domain.Account{}, fmt.Errorf("account not found or does not exist")
		}
		if errors.Is(err, services.ErrAccount) {
			slog.Error(
				"failed to get account by auth user id "+
					"perhaps a failed query or mapping",
				slog.String("authUserId", sub),
			)
			return domain.Account{}, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		if err != nil {
			slog.Error(
				"failed to get account by auth user id "+
					"something went wrong",
				slog.String("authUserId", sub),
			)
			return domain.Account{}, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		return account, nil
	} else {
		return account, fmt.Errorf("unsupported scheme `%s` cannot authenticate", scheme)
	}
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *domain.Account)

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func NewServer(
	ctx context.Context, // deprecate context passing, shall we use req.Context() instead?
	authService ports.AuthServicer,
	accountService ports.AccountServicer,
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadChatService ports.ThreadChatServicer,
) http.Handler {
	// initialize new server mux
	mux := http.NewServeMux()

	ah := NewAccountHandler(accountService)

	mux.HandleFunc("GET /{$}", handleGetIndex)

	// mux.Handle("POST /accounts/auth/{$}", ah.handleGetOrCreateAuthAccount(ctx))
	// mux.Handle("POST /pats/{$}", ah.handleCreatePAT(ctx))

	// mux.Handle("GET /pats/{$}", ah.handleGetPATs(ctx))

	mux.Handle("GET /pats/{$}", NewEnsureAuth(ah.handleGetPatList, authService))

	// wh := NewWorkspaceHandler(accountService, workspaceService)

	// mux.Handle("POST /workspaces/{$}", wh.handleCreateWorkspace(ctx))
	// mux.Handle("GET /workspaces/{$}", wh.handleGetWorkspaces(ctx))
	// mux.Handle("GET /workspaces/{workspaceId}/{$}", wh.handleGetWorkspace(ctx))
	// mux.Handle("POST /workspaces/{workspaceId}/labels/{$}", wh.handleGetOrCreateWorkspaceLabel(ctx))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
