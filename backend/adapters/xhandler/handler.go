package xhandler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

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

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *CustomerHandler) handleGetOrCreateCustomer(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var payload WidgetInitReqPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	widgetId := r.PathValue("widgetId")

	widget, err := h.ws.GetWorkspaceWidget(ctx, widgetId)
	if errors.Is(err, services.ErrWidgetNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace widget "+
				"something went wrong",
			slog.String("widgetId", widgetId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sk, err := h.ws.GetWorkspaceSecretKey(ctx, widget.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get workspace secret key "+
				"something went wrong",
			slog.String("workspaceId", widget.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var isCreated bool
	var isVerified bool
	var customer models.Customer

	customerHash := models.NullString(payload.CustomerHash)
	customerExternalId := models.NullString(payload.CustomerExternalId)
	customerEmail := models.NullString(payload.CustomerEmail)
	customerPhone := models.NullString(payload.CustomerPhone)

	anonId := models.NullString(payload.AnonId)
	customerName := models.Customer{}.AnonName()

	if payload.Traits != nil {
		if payload.Traits.Name != nil {
			customerName = *payload.Traits.Name
		} else {
			if payload.Traits.FirstName != nil || payload.Traits.LastName != nil {
				n := ""
				if payload.Traits.FirstName != nil {
					n += *payload.Traits.FirstName
				}
				if payload.Traits.LastName != nil {
					n += " " + *payload.Traits.LastName
				}
				customerName = n
			}
		}
	}

	if customerHash.Valid {
		sk, err := h.ws.GetWorkspaceSecretKey(ctx, widget.WorkspaceId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if customerExternalId.Valid {
			if h.cs.VerifyExternalId(sk.SecretKey, customerHash.String, customerExternalId.String) {
				isVerified = true
				customer = models.Customer{
					WorkspaceId: widget.WorkspaceId,
					ExternalId:  customerExternalId,
					IsVerified:  true,
					Role:        models.Customer{}.Engaged(),
					Name:        customerName,
				}
				customer, isCreated, err = h.ws.CreateCustomerWithExternalId(ctx, customer)
				if err != nil {
					slog.Error(
						"failed to create customer by externalId " +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		} else if customerEmail.Valid {
			if h.cs.VerifyEmail(sk.SecretKey, customerHash.String, customerEmail.String) {
				isVerified = true
				customer = models.Customer{
					WorkspaceId: widget.WorkspaceId,
					Email:       customerEmail,
					IsVerified:  true,
					Role:        models.Customer{}.Engaged(),
					Name:        customerName,
				}
				customer, isCreated, err = h.ws.CreateCustomerWithEmail(ctx, customer)
				if err != nil {
					slog.Error(
						"failed to create customer by email " +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		} else if customerPhone.Valid {
			if h.cs.VerifyPhone(sk.SecretKey, customerHash.String, customerPhone.String) {
				isVerified = true
				customer = models.Customer{
					WorkspaceId: widget.WorkspaceId,
					Phone:       customerPhone,
					IsVerified:  true,
					Role:        models.Customer{}.Engaged(),
					Name:        customerName,
				}
				customer, isCreated, err = h.ws.CreateCustomerWithPhone(ctx, customer)
				if err != nil {
					slog.Error(
						"failed to create customer by phone " +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
	} else if anonId.Valid {
		// make sure the anonymousId is a valid UUID
		isValid := models.IsValidUUID(anonId.String)
		if !isValid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		customer = models.Customer{
			WorkspaceId: widget.WorkspaceId,
			AnonId:      anonId.String,
			IsVerified:  false,
			Role:        models.Customer{}.Visitor(),
			Name:        customerName,
		}
		customer, isCreated, err = h.ws.CreateAnonymousCustomer(ctx, customer)
		if err != nil {
			slog.Error(
				"failed to create anonymous customer "+
					"something went wrong", slog.Any("error", err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		// force client to provide the anonymousId.
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	jwt, err := h.cs.GenerateCustomerToken(customer, sk.SecretKey)
	if err != nil {
		slog.Error("failed to make jwt token with error", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	denonEmail := customer.DeAnonEmail()
	denonPhone := customer.DeAnonPhone()
	denonExternalId := customer.DeAnonExternalId()

	resp := WidgetInitRespPayload{
		Jwt:        jwt,
		Create:     isCreated,
		IsVerified: isVerified,
		Name:       customer.Name,
		Email:      models.NullString(&denonEmail),
		Phone:      models.NullString(&denonPhone),
		ExternalId: models.NullString(&denonExternalId),
	}

	if isCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode customer to json " +
					"check the json encoding defn",
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode customer to json " +
					"check the json encoding defn",
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (h *CustomerHandler) handleGetCustomer(w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		slog.Error(
			"failed to encode customer to json "+
				"check the json encoding defn",
			slog.String("customerId", customer.CustomerId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleAddCustomerIdentities(w http.ResponseWriter, r *http.Request, customer *models.Customer) {

	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var payload AddCustomerIdentitiesReqPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	widgetId := r.PathValue("widgetId")

	widget, err := h.ws.GetWorkspaceWidget(ctx, widgetId)

	if errors.Is(err, services.ErrWidgetNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace widget "+
				"something went wrong",
			slog.String("widgetId", widgetId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	widgetCustomer, err := h.cs.GetWorkspaceCustomerById(ctx, widget.WorkspaceId, customer.CustomerId)
	if errors.Is(err, services.ErrCustomerNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace customer by id "+
				"something went wrong",
			slog.String("customerId", customer.CustomerId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// don't want to modify a verified customer
	// from the external widget.
	if widgetCustomer.IsVerified {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hasModified := false
	if payload.Email != nil {
		email := widgetCustomer.AddAnonimizedEmail(*payload.Email)
		widgetCustomer.Email = models.NullString(&email)
		hasModified = true
	}

	if payload.Phone != nil {
		phone := widgetCustomer.AddAnonimizedPhone(*payload.Phone)
		widgetCustomer.Phone = models.NullString(&phone)
		hasModified = true
	}

	if payload.External != nil {
		externalId := widgetCustomer.AddAnonimizedExternalId(*payload.External)
		widgetCustomer.ExternalId = models.NullString(&externalId)
		hasModified = true
	}

	if hasModified {
		widgetCustomer, err = h.cs.UpdateCustomer(ctx, widgetCustomer)
		if err != nil {
			slog.Error(
				"failed to update customer "+
					"something went wrong",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		denonEmail := widgetCustomer.DeAnonEmail()
		denonPhone := widgetCustomer.DeAnonPhone()
		denonExternalId := widgetCustomer.DeAnonExternalId()

		resp := AddCustomerIdentitiesRespPayload{
			IsVerified: widgetCustomer.IsVerified,
			Name:       widgetCustomer.Name,
			Email:      models.NullString(&denonEmail),
			Phone:      models.NullString(&denonPhone),
			ExternalId: models.NullString(&denonExternalId),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode customer to json "+
					"check the json encoding defn",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (h *CustomerHandler) handleCreateCustomerThChat(w http.ResponseWriter, r *http.Request, customer *models.Customer) {
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

	thread, chat, err := h.ths.CreateThreadInAppChat(ctx, workspace.WorkspaceId, customer.CustomerId, message.Message)
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

	var chatCustomer *ThCustomerResp
	var chatMember *ThMemberResp
	var threadAssignee *ThMemberResp

	// for chat - either of them
	if chat.CustomerId.Valid {
		chatCustomer = &ThCustomerResp{
			CustomerId: chat.CustomerId.String,
			Name:       chat.CustomerName.String,
		}
	} else if chat.MemberId.Valid {
		chatMember = &ThMemberResp{
			MemberId: chat.MemberId.String,
			Name:     chat.MemberName.String,
		}
	}

	chatResp := ChatResp{
		ThreadId:  thread.ThreadId,
		ChatId:    chat.ChatId,
		Body:      chat.Body,
		Sequence:  chat.Sequence,
		IsHead:    chat.IsHead,
		Customer:  chatCustomer,
		Member:    chatMember,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}

	// for thread
	threadCustomer := ThCustomerResp{
		CustomerId: thread.CustomerId,
		Name:       thread.CustomerName,
	}

	// for thread
	if thread.AssigneeId.Valid {
		threadAssignee = &ThMemberResp{
			MemberId: thread.AssigneeId.String,
			Name:     thread.AssigneeName.String,
		}
	}

	resp := ThreadChatResp{
		ThreadId:  thread.ThreadId,
		Sequence:  thread.Sequence,
		Status:    thread.Status,
		Read:      thread.Read,
		Replied:   thread.Replied,
		Priority:  thread.Priority,
		Customer:  threadCustomer,
		Assignee:  threadAssignee,
		Title:     thread.Title,
		Summary:   thread.Summary,
		Spam:      thread.Spam,
		Channel:   thread.Channel,
		CreatedAt: thread.CreatedAt,
		UpdatedAt: thread.UpdatedAt,
		Chat:      chatResp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat to json "+
				"check the json encoding defn",
			slog.String("threadId", thread.ThreadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetCustomerThChats(w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	ctx := r.Context()
	workspace, err := h.ws.GetWorkspace(ctx, customer.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
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

	threads, err := h.ths.ListCustomerThreadChats(ctx, workspace.WorkspaceId, customer.CustomerId)
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

	items := make([]ThreadResp, 0, 100)

	for _, thread := range threads {
		var threadAssignee *ThMemberResp
		var threadCustomer *ThCustomerResp
		var threadMember *ThMemberResp

		if thread.AssigneeId.Valid {
			threadAssignee = &ThMemberResp{
				MemberId: thread.AssigneeId.String,
				Name:     thread.AssigneeName.String,
			}
		}

		if thread.MessageCustomerId.Valid {
			threadCustomer = &ThCustomerResp{
				CustomerId: thread.MessageCustomerId.String,
				Name:       thread.MessageCustomerName.String,
			}
		} else if thread.MessageMemberId.Valid {
			threadMember = &ThMemberResp{
				MemberId: thread.MessageMemberId.String,
				Name:     thread.MessageMemberName.String,
			}
		}
		items = append(items, ThreadResp{
			ThreadId:        thread.ThreadId,
			Sequence:        thread.Sequence,
			Status:          thread.Status,
			Read:            thread.Read,
			Replied:         thread.Replied,
			Priority:        thread.Priority,
			Assignee:        threadAssignee,
			Title:           thread.Title,
			Summary:         thread.Summary,
			Spam:            thread.Spam,
			Channel:         thread.Channel,
			Body:            thread.MessageBody,
			MessageSequence: thread.MessageSequence,
			Customer:        threadCustomer,
			Member:          threadMember,
			CreatedAt:       thread.CreatedAt,
			UpdatedAt:       thread.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error(
			"failed to encode thread chats to json "+
				"check the json encoding defn",
			slog.String("customerId", customer.CustomerId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCreateThChatMessage(w http.ResponseWriter, r *http.Request, customer *models.Customer) {
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
	thread, err := h.ths.GetWorkspaceThread(ctx, workspace.WorkspaceId, threadId)

	if errors.Is(err, services.ErrThreadChatNotFound) {
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

	chat, err := h.ths.AddCustomerMessageToThread(ctx, thread.ThreadId, customer.CustomerId, message.Message)
	if err != nil {
		slog.Error(
			"failed to create thread chat message for customer "+
				"something went wrong",
			slog.String("threadId", thread.ThreadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var chatCustomer *ThCustomerResp
	var chatMember *ThMemberResp

	// for chat - either of them
	if chat.CustomerId.Valid {
		chatCustomer = &ThCustomerResp{
			CustomerId: chat.CustomerId.String,
			Name:       chat.CustomerName.String,
		}
	} else if chat.MemberId.Valid {
		chatMember = &ThMemberResp{
			MemberId: chat.MemberId.String,
			Name:     chat.MemberName.String,
		}
	}

	resp := ChatResp{
		ThreadId:  thread.ThreadId,
		ChatId:    chat.ChatId,
		Body:      chat.Body,
		Sequence:  chat.Sequence,
		IsHead:    chat.IsHead,
		Customer:  chatCustomer,
		Member:    chatMember,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat message to json "+
				"check the json encoding defn",
			slog.String("threadId", thread.ThreadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetThChatMesssages(w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	threadId := r.PathValue("threadId")

	ctx := r.Context()

	workspace, err := h.ws.GetWorkspace(ctx, customer.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
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
	thread, err := h.ths.GetWorkspaceThread(ctx, workspace.WorkspaceId, threadId)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	chats, err := h.ths.ListThreadChatMessages(ctx, thread.ThreadId)
	if err != nil {
		slog.Error(
			"failed to get list of thread chat messages for customer "+
				"something went wrong",
			slog.String("threadId", thread.ThreadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	messages := make([]ChatResp, 0, 100)
	for _, chat := range chats {
		var chatCustomer *ThCustomerResp
		var chatMember *ThMemberResp
		if chat.CustomerId.Valid {
			chatCustomer = &ThCustomerResp{
				CustomerId: chat.CustomerId.String,
				Name:       chat.CustomerName.String,
			}
		} else if chat.MemberId.Valid {
			chatMember = &ThMemberResp{
				MemberId: chat.MemberId.String,
				Name:     chat.MemberName.String,
			}
		}
		chatResp := ChatResp{
			ThreadId:  thread.ThreadId,
			ChatId:    chat.ChatId,
			Body:      chat.Body,
			Sequence:  chat.Sequence,
			IsHead:    chat.IsHead,
			Customer:  chatCustomer,
			Member:    chatMember,
			CreatedAt: chat.CreatedAt,
			UpdatedAt: chat.UpdatedAt,
		}
		messages = append(messages, chatResp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		slog.Error(
			"failed to encode thread chat messages to json "+
				"check the json encoding defn",
			slog.String("threadId", thread.ThreadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func NewServer(
	authService ports.CustomerAuthServicer,
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadChatService ports.ThreadChatServicer,
) http.Handler {
	mux := http.NewServeMux()

	ch := NewCustomerHandler(workspaceService, customerService, threadChatService)

	mux.HandleFunc("GET /{$}", handleGetIndex)

	mux.HandleFunc("POST /widgets/{widgetId}/init/{$}", ch.handleGetOrCreateCustomer)
	mux.Handle("GET /widgets/{widgetId}/me/{$}", NewEnsureAuth(ch.handleGetCustomer, authService))
	mux.Handle("POST /widgets/{widgetId}/me/identities/{$}", NewEnsureAuth(ch.handleAddCustomerIdentities, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{$}", NewEnsureAuth(ch.handleCreateCustomerThChat, authService))
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{$}", NewEnsureAuth(ch.handleGetCustomerThChats, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleCreateThChatMessage, authService))
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleGetThChatMesssages, authService))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
