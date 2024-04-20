package routes

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(ctx context.Context, db *pgxpool.Pool) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handleGetIndex)
	// TODO: rename route
	// web
	mux.Handle("POST /accounts/auth/{$}",
		handleGetOrCreateAuthAccount(ctx, db))
	// sdk+web
	mux.Handle("POST /pats/{$}", handleCreatePAT(ctx, db))
	// sdk+web
	mux.Handle("GET /pats/{$}", handleGetPATs(ctx, db))
	// sdk+web
	mux.Handle("POST /workspaces/{$}", handleCreateWorkspace(ctx, db))
	// sdk+web
	mux.Handle("GET /workspaces/{$}", handleGetWorkspaces(ctx, db))
	// sdk+web
	mux.Handle("POST /workspaces/{workspaceId}/threads/chat/{threadId}/messages/{$}",
		handleCreateMemberThChatMessage(ctx, db))
	// sdk+web
	mux.Handle("POST /workspaces/{workspaceId}/tokens/{$}",
		handleCustomerTokenIssue(ctx, db))
	// customer
	mux.Handle("GET /-/me/{$}", handleGetCustomer(ctx, db))
	// customer
	mux.Handle("POST /-/threads/chat/{$}", handleInitCustomerThreadChat(ctx, db))
	// customer
	mux.Handle("GET /-/threads/chat/{$}", handleGetCustomerThreadChats(ctx, db))
	// customer
	mux.Handle("POST /-/threads/chat/{threadId}/messages/{$}",
		handleCreateCustomerThChatMessage(ctx, db))
	// customer
	mux.Handle("GET /-/threads/chat/{threadId}/messages/{$}",
		handleGetCustomerThChatMessages(ctx, db))
	// customer
	mux.Handle("POST /-/threads/qa/{$}", handleInitCustomerThreadQA(ctx, db))

	return mux
}
