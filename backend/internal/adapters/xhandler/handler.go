package xhandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func AuthenticateCustomer(
	ctx context.Context, authz ports.CustomerAuthServicer,
	scheme string, cred string,
) (domain.Customer, error) {
	var customer domain.Customer
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

		customer, err = authz.GetWorkspaceCustomer(ctx, cc.WorkspaceId, sub)

		if errors.Is(err, services.ErrCustomerNotFound) {
			slog.Warn(
				"customer not found or does not exist",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("customer not found or does not exist")
		}

		if err != nil {
			slog.Error(
				"failed to get customer by customer id"+
					"something went wrong",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("failed to get customer by customer id: %s got error: %v", sub, err)
		}

		return customer, nil
	} else {
		return customer, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *domain.Customer)

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type CustomerHandler struct {
	ws  ports.WorkspaceServicer
	cs  ports.CustomerServicer
	ths ports.ThreadChatServicer
}

func NewCustomerHandler(
	ws ports.WorkspaceServicer,
	cs ports.CustomerServicer,
	ths ports.ThreadChatServicer,
) *CustomerHandler {
	return &CustomerHandler{
		ws:  ws,
		cs:  cs,
		ths: ths,
	}
}

func (h *CustomerHandler) handleGetCustomer(w http.ResponseWriter, r *http.Request, customer *domain.Customer) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		slog.Error(
			"failed to encode customer to json "+
				"might need to check the json encoding defn",
			slog.String("customerId", customer.CustomerId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCreateCustomerThChat(
	w http.ResponseWriter, r *http.Request, customer *domain.Customer,
) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var message ThChatReqPayload

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspace, err := h.ws.GetWorkspace(ctx, customer.WorkspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for customer",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	th := domain.ThreadChat{
		WorkspaceId:  workspace.WorkspaceId,
		CustomerId:   customer.CustomerId,
		CustomerName: customer.Name,
	}
	th, thm, err := h.ths.CreateCustomerThread(ctx, th, message.Message)
	if err != nil {
		slog.Error(
			"failed to create thread chat for customer "+
				"something went wrong",
			slog.String("customerId", customer.CustomerId),
			slog.String("workspaceid", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	messages := make([]ThChatMessageRespPayload, 0, 1)

	var msgCustomerRepr *ThCustomerRespPayload
	var msgMemberRepr *ThMemberRespPayload

	// for thread message - either of them
	if thm.CustomerId.Valid {
		msgCustomerRepr = &ThCustomerRespPayload{
			CustomerId: thm.CustomerId.String,
			Name:       thm.CustomerName,
		}
	} else if thm.MemberId.Valid {
		msgMemberRepr = &ThMemberRespPayload{
			MemberId: thm.MemberId.String,
			Name:     thm.MemberName,
		}
	}

	threadMessage := ThChatMessageRespPayload{
		ThreadChatId:        th.ThreadChatId,
		ThreadChatMessageId: thm.ThreadChatMessageId,
		Body:                thm.Body,
		Sequence:            thm.Sequence,
		Customer:            msgCustomerRepr,
		Member:              msgMemberRepr,
		CreatedAt:           thm.CreatedAt,
		UpdatedAt:           thm.UpdatedAt,
	}

	messages = append(messages, threadMessage)

	var threadAssigneeRepr *ThMemberRespPayload

	// for thread
	threadCustomerRepr := ThCustomerRespPayload{
		CustomerId: th.CustomerId,
		Name:       th.CustomerName,
	}

	// for thread
	if th.AssigneeId.Valid {
		threadAssigneeRepr = &ThMemberRespPayload{
			MemberId: th.AssigneeId.String,
			Name:     th.AssigneeName,
		}
	}

	resp := ThChatRespPayload{
		ThreadChatId: th.ThreadChatId,
		Sequence:     th.Sequence,
		Status:       th.Status,
		Read:         th.Read,
		Replied:      th.Replied,
		Customer:     threadCustomerRepr,
		Assignee:     threadAssigneeRepr,
		CreatedAt:    th.CreatedAt,
		UpdatedAt:    th.UpdatedAt,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat to json "+
				"might need to check the json encoding defn",
			slog.String("threadChatId", th.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCreateThChatMessage(
	w http.ResponseWriter, r *http.Request, customer *domain.Customer,
) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	threadId := r.PathValue("threadId")

	var message ThChatReqPayload

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspace, err := h.ws.GetWorkspace(ctx, customer.WorkspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for customer",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	th, err := h.ths.GetWorkspaceThread(ctx, workspace.WorkspaceId, threadId)

	if errors.Is(err, services.ErrThreadChatNotFound) {
		slog.Warn(
			"thread chat not found or does not exist for customer",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get thread chat by id "+
				"something went wrong",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	thm, err := h.ths.CreateCustomerMessage(ctx, th, customer, message.Message)

	if err != nil {
		slog.Error(
			"failed to create thread chat message for customer "+
				"something went wrong",
			slog.String("threadChatId", th.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var threadAssigneeRepr *ThMemberRespPayload
	var msgCustomerRepr *ThCustomerRespPayload
	var msgMemberRepr *ThMemberRespPayload

	// for thread
	threadCustomerRepr := ThCustomerRespPayload{
		CustomerId: th.CustomerId,
		Name:       th.CustomerName,
	}

	// for thread
	if th.AssigneeId.Valid {
		threadAssigneeRepr = &ThMemberRespPayload{
			MemberId: th.AssigneeId.String,
			Name:     th.AssigneeName,
		}
	}

	// for thread message - either of them
	if thm.CustomerId.Valid {
		msgCustomerRepr = &ThCustomerRespPayload{
			CustomerId: thm.CustomerId.String,
			Name:       thm.CustomerName,
		}
	} else if thm.MemberId.Valid {
		msgMemberRepr = &ThMemberRespPayload{
			MemberId: thm.MemberId.String,
			Name:     thm.MemberName,
		}
	}

	threadMessage := ThChatMessageRespPayload{
		ThreadChatId:        th.ThreadChatId,
		ThreadChatMessageId: thm.ThreadChatMessageId,
		Body:                thm.Body,
		Sequence:            thm.Sequence,
		Customer:            msgCustomerRepr,
		Member:              msgMemberRepr,
		CreatedAt:           thm.CreatedAt,
		UpdatedAt:           thm.UpdatedAt,
	}

	messages := make([]ThChatMessageRespPayload, 0, 1)
	messages = append(messages, threadMessage)
	resp := ThChatRespPayload{
		ThreadChatId: th.ThreadChatId,
		Sequence:     th.Sequence,
		Status:       th.Status,
		Read:         th.Read,
		Replied:      th.Replied,
		Customer:     threadCustomerRepr,
		Assignee:     threadAssigneeRepr,
		CreatedAt:    th.CreatedAt,
		UpdatedAt:    th.UpdatedAt,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat message to json "+
				"might need to check the json encoding defn",
			slog.String("threadChatId", th.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetThChatMesssages(
	w http.ResponseWriter, r *http.Request, customer *domain.Customer,
) {
	threadId := r.PathValue("threadId")

	ctx := r.Context()

	workspace, err := h.ws.GetWorkspace(ctx, customer.WorkspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for customer",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	th, err := h.ths.GetWorkspaceThread(ctx, workspace.WorkspaceId, threadId)

	if errors.Is(err, services.ErrThreadChatNotFound) {
		slog.Warn(
			"thread chat not found or does not exist for customer",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	results, err := h.ths.GetMessageList(ctx, th.ThreadChatId)

	if err != nil {
		slog.Error(
			"failed to get list of thread chat messages for customer "+
				"something went wrong",
			slog.String("threadChatId", th.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	messages := make([]ThChatMessageRespPayload, 0, 100)
	for _, thm := range results {
		var msgCustomerRepr *ThCustomerRespPayload
		var msgMemberRepr *ThMemberRespPayload

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerRespPayload{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberRespPayload{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThChatMessageRespPayload{
			ThreadChatId:        th.ThreadChatId,
			ThreadChatMessageId: thm.ThreadChatMessageId,
			Body:                thm.Body,
			Sequence:            thm.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           thm.CreatedAt,
			UpdatedAt:           thm.UpdatedAt,
		}

		messages = append(messages, threadMessage)
	}

	var threadAssigneeRepr *ThMemberRespPayload

	// for thread
	threadCustomerRepr := ThCustomerRespPayload{
		CustomerId: th.CustomerId,
		Name:       th.CustomerName,
	}

	// for thread
	if th.AssigneeId.Valid {
		threadAssigneeRepr = &ThMemberRespPayload{
			MemberId: th.AssigneeId.String,
			Name:     th.AssigneeName,
		}
	}

	resp := ThChatRespPayload{
		ThreadChatId: th.ThreadChatId,
		Sequence:     th.Sequence,
		Status:       th.Status,
		Read:         th.Read,
		Replied:      th.Replied,
		Customer:     threadCustomerRepr,
		Assignee:     threadAssigneeRepr,
		CreatedAt:    th.CreatedAt,
		UpdatedAt:    th.UpdatedAt,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat messages to json "+
				"might need to check the json encoding defn",
			slog.String("threadChatId", th.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetCustomerThChats(
	w http.ResponseWriter, r *http.Request, customer *domain.Customer,
) {
	ctx := r.Context()
	workspace, err := h.ws.GetWorkspace(ctx, customer.WorkspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for customer",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			slog.String("workspaceId", customer.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ths, err := h.ths.GetWorkspaceCustomerList(ctx, workspace.WorkspaceId, customer.CustomerId)

	if errors.Is(err, services.ErrThreadChat) {
		slog.Error(
			"failed to get list of thread chats for customer "+
				"perhaps a failed query or mapping",
			slog.String("customerId", customer.CustomerId),
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get list of thread chats for customer "+
				"something went wrong",
			slog.String("customerId", customer.CustomerId),
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threads := make([]ThChatRespPayload, 0, 100)
	for _, th := range ths {
		messages := make([]ThChatMessageRespPayload, 0, 1)

		var threadAssigneeRepr *ThMemberRespPayload
		var msgCustomerRepr *ThCustomerRespPayload
		var msgMemberRepr *ThMemberRespPayload

		// for thread
		threadCustomerRepr := ThCustomerRespPayload{
			CustomerId: th.ThreadChat.CustomerId,
			Name:       th.ThreadChat.CustomerName,
		}

		// for thread
		if th.ThreadChat.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberRespPayload{
				MemberId: th.ThreadChat.AssigneeId.String,
				Name:     th.ThreadChat.AssigneeName,
			}
		}

		// for thread message - either of them
		if th.Message.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerRespPayload{
				CustomerId: th.Message.CustomerId.String,
				Name:       th.Message.CustomerName,
			}
		} else if th.Message.MemberId.Valid {
			msgMemberRepr = &ThMemberRespPayload{
				MemberId: th.Message.MemberId.String,
				Name:     th.Message.MemberName,
			}
		}

		message := ThChatMessageRespPayload{
			ThreadChatId:        th.ThreadChat.ThreadChatId,
			ThreadChatMessageId: th.Message.ThreadChatMessageId,
			Body:                th.Message.Body,
			Sequence:            th.Message.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           th.Message.CreatedAt,
			UpdatedAt:           th.Message.UpdatedAt,
		}
		messages = append(messages, message)
		threads = append(threads, ThChatRespPayload{
			ThreadChatId: th.ThreadChat.ThreadChatId,
			Sequence:     th.ThreadChat.Sequence,
			Status:       th.ThreadChat.Status,
			Read:         th.ThreadChat.Read,
			Replied:      th.ThreadChat.Replied,
			Customer:     threadCustomerRepr,
			Assignee:     threadAssigneeRepr,
			CreatedAt:    th.ThreadChat.CreatedAt,
			UpdatedAt:    th.ThreadChat.UpdatedAt,
			Messages:     messages,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(threads); err != nil {
		slog.Error(
			"failed to encode thread chats to json "+
				"might need to check the json encoding defn",
			slog.String("customerId", customer.CustomerId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func NewServer(
	ctx context.Context, // deprecate context passing, shall we use req.Context() instead?
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadChatService ports.ThreadChatServicer,
) http.Handler {
	// initialize new server mux
	mux := http.NewServeMux()

	ch := NewCustomerHandler(workspaceService, customerService, threadChatService)

	mux.HandleFunc("GET /{$}", handleGetIndex)

	mux.Handle("GET /me/{$}", NewEnsureAuth(ch.handleGetCustomer, customerService))

	mux.Handle("POST /threads/chat/{$}", NewEnsureAuth(ch.handleCreateCustomerThChat, customerService))

	mux.Handle("GET /threads/chat/{$}", NewEnsureAuth(ch.handleGetCustomerThChats, customerService))

	mux.Handle("POST /threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleCreateThChatMessage, customerService))

	mux.Handle("GET /threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleGetThChatMesssages, customerService))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
