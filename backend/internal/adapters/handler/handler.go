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

func AuthenticateAccount(ctx context.Context, authz ports.AuthServicer, scheme string, cred string) (domain.Account, error) {
	if scheme == "token" {
		account, err := authz.CheckPatAccount(ctx, cred)
		if err != nil {
			return account, fmt.Errorf("failed to authenticate got error: %v", err)
		}
		slog.Info("authenticated account with PAT", slog.String("accountId", account.AccountId))
		return account, nil
	} else if scheme == "bearer" {
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return domain.Account{}, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
		}
		ac, err := services.ParseJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			return domain.Account{}, fmt.Errorf("failed to parse JWT token got error: %v", err)
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			return domain.Account{}, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}

		slog.Info("authenticated account with jwt", slog.String("authUserId", sub))
		account, err := authz.CheckAuthUser(ctx, sub)

		if errors.Is(err, services.ErrAccountNotFound) {
			slog.Warn(
				"account not found or does not exist",
				slog.String("authUserId", sub),
			)
			return account, fmt.Errorf("account not found or does not exist")
		}
		if errors.Is(err, services.ErrAccount) {
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
		return account, nil
	} else {
		return domain.Account{}, fmt.Errorf("unsupported scheme `%s` cannot authenticate", scheme)
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
	authService ports.AuthServicer,
	accountService ports.AccountServicer,
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadChatService ports.ThreadChatServicer,
) http.Handler {
	mux := http.NewServeMux()

	// creates service handlers
	ah := NewAccountHandler(accountService)
	wh := NewWorkspaceHandler(workspaceService, customerService)
	th := NewThreadChatHandler(workspaceService, threadChatService)

	mux.HandleFunc("GET /{$}", handleGetIndex)
	mux.HandleFunc("POST /accounts/auth/{$}", ah.handleGetOrCreateAccount)

	mux.Handle("POST /pats/{$}", NewEnsureAuth(ah.handleCreatePat, authService))
	mux.Handle("GET /pats/{$}", NewEnsureAuth(ah.handleGetPatList, authService))

	mux.Handle("POST /workspaces/{$}", NewEnsureAuth(wh.handleCreateWorkspace, authService))
	mux.Handle("GET /workspaces/{$}", NewEnsureAuth(wh.handleGetWorkspaces, authService))

	mux.Handle("GET /workspaces/{workspaceId}/{$}", NewEnsureAuth(wh.handleGetWorkspace, authService))
	mux.Handle("PATCH /workspaces/{workspaceId}/{$}", NewEnsureAuth(wh.handleUpdateWorkspace, authService))

	mux.Handle("GET /workspaces/{workspaceId}/members/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceMembers, authService))
	mux.Handle("GET /workspaces/{workspaceId}/members/me/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceMembership, authService))

	mux.Handle("GET /workspaces/{workspaceId}/customers/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceCustomers, authService))

	mux.Handle("POST /workspaces/{workspaceId}/labels/{$}",
		NewEnsureAuth(wh.handleGetOrCreateWorkspaceLabel, authService))
	mux.Handle("GET /workspaces/{workspaceId}/labels/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceLabels, authService))

	// returns all threads in a workspace
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/{$}",
		NewEnsureAuth(th.handleGetThreadChats, authService))

	// returns all threads in a workspace assigned to the authenticated member
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/with/me/{$}",
		NewEnsureAuth(th.handleGetMyThreadChats, authService))

	// returns all threads in a workspace unassigned
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/with/unassigned/{$}",
		NewEnsureAuth(th.handleGetUnassignedThreadChats, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/with/labels/{labelId}/{$}",
		NewEnsureAuth(th.handleGetLabelledThreadChats, authService))

	mux.Handle("POST /workspaces/{workspaceId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(th.handleCreateThChatMessage, authService))

	mux.Handle("PUT /workspaces/{workspaceId}/threads/chat/{threadId}/labels/{$}",
		NewEnsureAuth(th.handleSetThChatLabel, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/{threadId}/labels/{$}",
		NewEnsureAuth(th.handleGetThreadChatLabels, authService))

	mux.Handle("POST /workspaces/{workspaceId}/customers/tokens/{$}",
		NewEnsureAuth(wh.handleIssueCustomerToken, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/metrics/{$}",
		NewEnsureAuth(th.handleGetThreadChatMetrics, authService))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
