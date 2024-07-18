package handler

import (
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg/ports"
)

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

	// initialize service handlers
	ah := NewAccountHandler(accountService, workspaceService)
	wh := NewWorkspaceHandler(workspaceService, customerService)
	th := NewThreadChatHandler(workspaceService, threadChatService)

	mux.HandleFunc("GET /{$}", handleGetIndex)
	mux.HandleFunc("POST /accounts/auth/{$}", ah.handleGetOrCreateAccount)

	mux.Handle("POST /pats/{$}", NewEnsureAuth(ah.handleCreatePat, authService))
	mux.Handle("GET /pats/{$}", NewEnsureAuth(ah.handleGetPatList, authService))
	mux.Handle("DELETE /pats/{patId}/{$}", NewEnsureAuth(ah.handleDeletePat, authService))

	mux.Handle("POST /workspaces/{$}", NewEnsureAuth(wh.handleCreateWorkspace, authService))
	mux.Handle("GET /workspaces/{$}", NewEnsureAuth(wh.handleGetWorkspaces, authService))

	mux.Handle("GET /workspaces/{workspaceId}/{$}", NewEnsureAuth(wh.handleGetWorkspace, authService))
	mux.Handle("PATCH /workspaces/{workspaceId}/{$}", NewEnsureAuth(wh.handleUpdateWorkspace, authService))

	// generate workspace secret key
	mux.Handle("POST /workspaces/{workspaceId}/sk/{$}",
		NewEnsureAuth(wh.handleGenerateSecretKey, authService))
	mux.Handle("GET /workspaces/{workspaceId}/sk/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceSecretKey, authService))

	mux.Handle("GET /workspaces/{workspaceId}/members/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceMembers, authService))
	mux.Handle("GET /workspaces/{workspaceId}/members/me/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceMembership, authService))

	mux.Handle("GET /workspaces/{workspaceId}/members/{memberId}/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceMember, authService))

	mux.Handle("GET /workspaces/{workspaceId}/customers/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceCustomers, authService))

	mux.Handle("POST /workspaces/{workspaceId}/labels/{$}",
		NewEnsureAuth(wh.handleGetOrCreateWorkspaceLabel, authService))
	mux.Handle("GET /workspaces/{workspaceId}/labels/{$}",
		NewEnsureAuth(wh.handleGetWorkspaceLabels, authService))
	mux.Handle("PATCH /workspaces/{workspaceId}/labels/{labelId}/{$}",
		NewEnsureAuth(wh.handleUpdateWorkspaceLabel, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/{$}",
		NewEnsureAuth(th.handleGetThreadChats, authService))

	mux.Handle("PATCH /workspaces/{workspaceId}/threads/chat/{threadId}/{$}",
		NewEnsureAuth(th.handleUpdateThreadChat, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/with/me/{$}",
		NewEnsureAuth(th.handleGetMyThreadChats, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/with/unassigned/{$}",
		NewEnsureAuth(th.handleGetUnassignedThreadChats, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/with/labels/{labelId}/{$}",
		NewEnsureAuth(th.handleGetLabelledThreadChats, authService))

	mux.Handle("POST /workspaces/{workspaceId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(th.handleCreateThChatMessage, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(th.handleGetThChatMesssages, authService))

	mux.Handle("PUT /workspaces/{workspaceId}/threads/chat/{threadId}/labels/{$}",
		NewEnsureAuth(th.handleSetThChatLabel, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/{threadId}/labels/{$}",
		NewEnsureAuth(th.handleGetThreadChatLabels, authService))

	mux.Handle("POST /workspaces/{workspaceId}/customers/tokens/{$}",
		NewEnsureAuth(wh.handleIssueCustomerToken, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/chat/metrics/{$}",
		NewEnsureAuth(th.handleGetThreadChatMetrics, authService))

	// creates a new workspace widget
	mux.Handle("POST /workspaces/{workspaceId}/widgets/{$}",
		NewEnsureAuth(wh.handleCreateWidget, authService))
	mux.Handle("GET /workspaces/{workspaceId}/widgets/{$}",
		NewEnsureAuth(wh.handleGetWidgets, authService))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
