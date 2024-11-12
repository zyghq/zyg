package xhandler

import (
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg/ports"
)

type CustomerHandler struct {
	ws  ports.WorkspaceServicer
	cs  ports.CustomerServicer
	ths ports.ThreadServicer
}

func NewCustomerHandler(
	ws ports.WorkspaceServicer,
	cs ports.CustomerServicer,
	ths ports.ThreadServicer,
) *CustomerHandler {
	return &CustomerHandler{
		ws:  ws,
		cs:  cs,
		ths: ths,
	}
}

type ThreadHandler struct {
	ws  ports.WorkspaceServicer
	ths ports.ThreadServicer
}

func NewThreadHandler(ws ports.WorkspaceServicer, ths ports.ThreadServicer) *ThreadHandler {
	return &ThreadHandler{
		ws:  ws,
		ths: ths,
	}
}

// handleGetIndex returns the API index.
func handleGetIndex(w http.ResponseWriter, _ *http.Request) {
	tm := time.Now().UTC().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

func NewServer(
	authService ports.CustomerAuthServicer,
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadChatService ports.ThreadServicer,
) http.Handler {
	// init new server mux
	mux := http.NewServeMux()
	// init handlers
	ch := NewCustomerHandler(workspaceService, customerService, threadChatService)
	th := NewThreadHandler(workspaceService, threadChatService)

	mux.HandleFunc("GET /{$}", handleGetIndex)

	mux.HandleFunc("GET /mail/kyc/{$}", ch.handleMailRedirectKyc)

	mux.HandleFunc("GET /widgets/{widgetId}/config/{$}", ch.handleGetWidgetConfig)
	mux.HandleFunc("POST /widgets/{widgetId}/init/{$}", ch.handleInitWidget)

	mux.Handle("GET /widgets/{widgetId}/me/{$}",
		NewEnsureAuth(ch.handleGetCustomer, authService))

	// Creates a new chat thread.
	mux.Handle("POST /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleCreateThreadChat, authService))
	// Returns a list of chat threads.
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleGetCustomerThreadChats, authService))
	// Creates a new thread chat message.
	mux.Handle("POST /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleCreateThreadChatMessage, authService))
	// Returns a list of thread chat messages.
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleGetThreadChatMessages, authService))

	// handles postmark inbound message webhook for workspace.
	// This URL path must also be configured in the postmark inbound settings.
	mux.HandleFunc("POST /webhooks/{workspaceId}/postmark/inbound/{$}",
		th.handlePostmarkInboundMessage)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
