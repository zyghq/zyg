package handler

import (
	"github.com/zyghq/zyg"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg/ports"
)

func handleGetIndex(w http.ResponseWriter, _ *http.Request) {
	tm := time.Now().UTC().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK\n"))
	if err != nil {
		return
	}
}

// NewServer initializes and returns a http.Handler with all route handlers set up.
// It takes multiple service interfaces as parameters for dependency injection.
func NewServer(
	authService ports.AuthServicer,
	accountService ports.AccountServicer,
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadService ports.ThreadServicer,
	syncService ports.SyncServicer,
) http.Handler {
	mux := http.NewServeMux()

	// initialize service handlers
	ah := NewAccountHandler(accountService, workspaceService)
	wh := NewWorkspaceHandler(workspaceService, accountService, customerService, syncService)
	th := NewThreadHandler(workspaceService, threadService, syncService)
	ch := NewCustomerHandler(workspaceService, customerService)
	ss := NewSyncHandler(workspaceService, threadService)

	webhookUsername := zyg.WebhookUsername()
	webhookPassword := zyg.WebhookPassword()

	mux.HandleFunc("GET /{$}", handleGetIndex)

	//mux.HandleFunc("GET /auth/login/{$}", handleWorkOSAuthLogin)
	//mux.HandleFunc("GET /auth/callback/{$}", handleWorkOSAuthCallback)

	mux.HandleFunc("POST /accounts/auth/{$}", ah.handleGetOrCreateAccount)

	// Todo: deprecate PAT usage, instead create workspace member tokens, with permissions.
	mux.Handle("POST /pats/{$}", NewEnsureAuthAccount(ah.handleCreatePat, authService))
	mux.Handle("GET /pats/{$}", NewEnsureAuthAccount(ah.handleGetPatList, authService))
	mux.Handle("DELETE /pats/{patId}/{$}", NewEnsureAuthAccount(ah.handleDeletePat, authService))

	// synced
	mux.Handle("POST /workspaces/{$}", NewEnsureAuthAccount(wh.handleCreateWorkspace, authService))
	mux.Handle("GET /workspaces/{$}", NewEnsureAuthAccount(wh.handleGetWorkspaces, authService))

	mux.Handle("GET /workspaces/{workspaceId}/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspace, authService))
	// synced
	mux.Handle("PATCH /workspaces/{workspaceId}/{$}",
		NewEnsureMemberAuth(wh.handleUpdateWorkspace, authService))

	mux.Handle("POST /workspaces/{workspaceId}/sk/{$}",
		NewEnsureMemberAuth(wh.handleGenerateSecretKey, authService))
	mux.Handle("GET /workspaces/{workspaceId}/sk/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceSecretKey, authService))

	mux.Handle("GET /workspaces/{workspaceId}/members/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceMembers, authService))
	mux.Handle("GET /workspaces/{workspaceId}/members/me/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceMembership, authService))
	mux.Handle("GET /workspaces/{workspaceId}/members/{memberId}/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceMember, authService))

	// synced
	mux.Handle("POST /workspaces/{workspaceId}/customers/{$}",
		NewEnsureMemberAuth(wh.handleCreateWorkspaceCustomer, authService))
	mux.Handle("GET /workspaces/{workspaceId}/customers/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceCustomers, authService))

	mux.Handle("POST /workspaces/{workspaceId}/customers/events/{$}",
		NewEnsureMemberAuth(ch.handleCreateCustomerEvent, authService))
	mux.Handle("GET /workspaces/{workspaceId}/customers/events/{customerId}/{$}",
		NewEnsureMemberAuth(ch.handleGetCustomerEvents, authService))

	mux.Handle("POST /workspaces/{workspaceId}/labels/{$}",
		NewEnsureMemberAuth(wh.handleCreateWorkspaceLabel, authService))
	mux.Handle("GET /workspaces/{workspaceId}/labels/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceLabels, authService))

	mux.Handle("PATCH /workspaces/{workspaceId}/labels/{labelId}/{$}",
		NewEnsureMemberAuth(wh.handleUpdateWorkspaceLabel, authService))
	mux.Handle("GET /workspaces/{workspaceId}/labels/{labelId}/{$}",
		NewEnsureMemberAuth(wh.handleGetWorkspaceLabel, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/{$}",
		NewEnsureMemberAuth(th.handleGetThreads, authService))
	// synced
	mux.Handle("PATCH /workspaces/{workspaceId}/threads/{threadId}/{$}",
		NewEnsureMemberAuth(th.handleUpdateThread, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/parts/me/{$}",
		NewEnsureMemberAuth(th.handleGetMyThreads, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/parts/unassigned/{$}",
		NewEnsureMemberAuth(th.handleGetUnassignedThreads, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/parts/labels/{labelId}/{$}",
		NewEnsureMemberAuth(th.handleGetLabelledThreads, authService))

	mux.Handle("POST /workspaces/{workspaceId}/threads/email/{threadId}/messages/{$}",
		NewEnsureMemberAuth(th.handleReplyThreadMail, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/{threadId}/messages/{$}",
		NewEnsureMemberAuth(th.handleGetThreadMessages, authService))

	mux.Handle("GET /workspaces/{workspaceId}/messages/{messageId}/attachments/{attachmentId}/{$}",
		NewEnsureMemberAuth(th.handleGetMessageAttachment, authService))

	// synced
	mux.Handle("PUT /workspaces/{workspaceId}/threads/{threadId}/labels/{$}",
		NewEnsureMemberAuth(th.handleSetThreadLabel, authService))
	mux.Handle("GET /workspaces/{workspaceId}/threads/{threadId}/labels/{$}",
		NewEnsureMemberAuth(th.handleGetThreadLabels, authService))
	// synced
	mux.Handle("DELETE /workspaces/{workspaceId}/threads/{threadId}/labels/{labelId}/{$}",
		NewEnsureMemberAuth(th.handleDeleteThreadLabel, authService))

	mux.Handle("GET /workspaces/{workspaceId}/threads/metrics/{$}",
		NewEnsureMemberAuth(th.handleGetThreadMetrics, authService))

	mux.Handle("POST /workspaces/{workspaceId}/widgets/{$}",
		NewEnsureMemberAuth(wh.handleCreateWidget, authService))
	mux.Handle("GET /workspaces/{workspaceId}/widgets/{$}",
		NewEnsureMemberAuth(wh.handleGetWidgets, authService))

	mux.Handle("GET /workspaces/{workspaceId}/postmark/servers/{$}",
		NewEnsureMemberAuth(wh.handlePostmarkGetMailServer, authService))

	mux.Handle("POST /workspaces/{workspaceId}/postmark/servers/{$}",
		NewEnsureMemberAuth(wh.handlePostmarkCreateMailServer, authService))

	mux.Handle("POST /workspaces/{workspaceId}/postmark/servers/parts/dns/add/{$}",
		NewEnsureMemberAuth(wh.handlePostmarkMailServerAddDNS, authService))

	mux.Handle("PUT /workspaces/{workspaceId}/postmark/servers/parts/dns/verify/{$}",
		NewEnsureMemberAuth(wh.handlePostmarkMailServerVerifyDNS, authService))

	mux.Handle("PATCH /workspaces/{workspaceId}/postmark/servers/{$}",
		NewEnsureMemberAuth(wh.handlePostmarkUpdateMailServer, authService))

	// v1 sync handlers
	mux.Handle("GET /v1/sync/workspaces/{workspaceId}/shapes/parts/members/{$}",
		NewEnsureMemberAuth(ss.syncWorkspaceMemberShapesV1, authService))

	mux.Handle("GET /v1/sync/workspaces/{workspaceId}/shapes/parts/customers/{$}",
		NewEnsureMemberAuth(ss.syncWorkspaceCustomerShapesV1, authService))

	mux.Handle("GET /v1/sync/workspaces/{workspaceId}/shapes/parts/threads/{$}",
		NewEnsureMemberAuth(ss.syncWorkspaceCustomerShapesV1, authService))

	// Webhooks
	// handles postmark inbound message webhook for workspace.
	// This URL path must also be configured in the postmark inbound settings.
	mux.HandleFunc("POST /webhooks/{workspaceId}/postmark/inbound/{$}",
		BasicAuthWebhook(th.handlePostmarkInboundWebhook, webhookUsername, webhookPassword))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
