package xhandler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
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

func handleGetIndex(w http.ResponseWriter, _ *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

func (h *CustomerHandler) handleGetOrCreateCustomer(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp WidgetInitReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	widgetId := r.PathValue("widgetId")
	widget, err := h.ws.GetWidget(ctx, widgetId)
	if errors.Is(err, services.ErrWidgetNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace widget", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sk, err := h.ws.GetSecretKey(ctx, widget.WorkspaceId)
	if errors.Is(err, services.ErrSecretKeyNotFound) {
		// if the secret key is not found, then the widget cannot be authorized.
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace secret key", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var isCreated bool
	var customer models.Customer

	customerHash := models.NullString(reqp.CustomerHash)
	customerExternalId := models.NullString(reqp.CustomerExternalId)
	customerEmail := models.NullString(reqp.CustomerEmail)
	customerPhone := models.NullString(reqp.CustomerPhone)

	anonId := models.NullString(reqp.AnonId)
	customerName := models.Customer{}.AnonName()

	// if the customer traits are provided, then check for name in traits.
	if reqp.Traits != nil {
		if reqp.Traits.Name != nil {
			customerName = *reqp.Traits.Name
		} else {
			if reqp.Traits.FirstName != nil || reqp.Traits.LastName != nil {
				n := ""
				if reqp.Traits.FirstName != nil {
					n += *reqp.Traits.FirstName
				}
				if reqp.Traits.LastName != nil {
					n += " " + *reqp.Traits.LastName
				}
				customerName = strings.Trim(n, " ")
			}
		}
	}

	if customerHash.Valid {
		if customerExternalId.Valid {
			if h.cs.VerifyExternalId(sk.Hmac, customerHash.String, customerExternalId.String) {
				customer, isCreated, err = h.ws.CreateCustomerWithExternalId(
					ctx, widget.WorkspaceId,
					customerExternalId.String,
					true,
					customerName,
				)
				if err != nil {
					slog.Error("failed to create customer by externalId", slog.Any("err", err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		} else if customerEmail.Valid {
			if h.cs.VerifyEmail(sk.Hmac, customerHash.String, customerEmail.String) {
				customer, isCreated, err = h.ws.CreateCustomerWithEmail(
					ctx, widget.WorkspaceId,
					customerEmail.String,
					true,
					customerName,
				)
				if err != nil {
					slog.Error("failed to create customer by email", slog.Any("err", err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		} else if customerPhone.Valid {
			if h.cs.VerifyPhone(sk.Hmac, customerHash.String, customerPhone.String) {
				customer, isCreated, err = h.ws.CreateCustomerWithPhone(
					ctx, widget.WorkspaceId,
					customerPhone.String,
					true,
					customerName,
				)
				if err != nil {
					slog.Error("failed to create customer by phone", slog.Any("err", err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
	} else if anonId.Valid {
		isValid := models.IsValidUUID(anonId.String)
		if !isValid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		customer, isCreated, err = h.ws.CreateAnonymousCustomer(
			ctx, widget.WorkspaceId,
			anonId.String,
			customerName,
		)
		if err != nil {
			slog.Error("failed to create anonymous customer", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		// force the client to provide the anonymousId.
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	jwt, err := h.cs.GenerateCustomerJwt(customer, sk.Hmac)
	if err != nil {
		slog.Error("failed to make jwt token with error", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := WidgetInitResp{
		Jwt:         jwt,
		Create:      isCreated,
		IsAnonymous: customer.IsAnonymous,
		Name:        customer.Name,
		AvatarUrl:   customer.AvatarUrl,
		Email:       customer.Email,
		Phone:       customer.Phone,
		ExternalId:  customer.ExternalId,
	}
	if isCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode response", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode response", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (h *CustomerHandler) handleGetCustomer(w http.ResponseWriter, _ *http.Request, customer *models.Customer) {
	var email string
	var phone string
	var externalId string

	if customer.Email.Valid {
		email = customer.Email.String
	}

	if customer.Phone.Valid {
		phone = customer.Phone.String
	}

	if customer.ExternalId.Valid {
		externalId = customer.ExternalId.String
	}

	resp := CustomerResp{
		CustomerId:  customer.CustomerId,
		Name:        customer.Name,
		AvatarUrl:   customer.AvatarUrl,
		Email:       models.NullString(&email),
		Phone:       models.NullString(&phone),
		ExternalId:  models.NullString(&externalId),
		IsAnonymous: customer.IsAnonymous,
		Role:        customer.Role,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCustomerIdentities(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {

	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp CustomerIdentitiesReq

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	widgetId := r.PathValue("widgetId")
	widget, err := h.ws.GetWidget(ctx, widgetId)
	if errors.Is(err, services.ErrWidgetNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get workspace widget", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// get the customer at the widget.
	widgetCustomer, err := h.ws.GetCustomer(ctx, widget.WorkspaceId, customer.CustomerId, nil)
	if errors.Is(err, services.ErrCustomerNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get workspace customer", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// modify only if the customer is anonymous.
	// from the external widget.
	if !widgetCustomer.IsAnonymous {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// create email identity.
	if reqp.Email != nil {
		// check if email already exists for a customer, if then there is a conflict.
		hasConflict, err := h.ws.DoesEmailConflict(ctx, customer.WorkspaceId, *reqp.Email)
		if err != nil {
			slog.Error("failed to check email conflict", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		emailIdentity := models.EmailIdentity{
			CustomerId:  customer.CustomerId,
			Email:       *reqp.Email,
			IsVerified:  false,
			HasConflict: hasConflict,
		}
		emailIdentity, err = h.cs.AddCustomerEmailIdentity(ctx, emailIdentity)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp := AddCustomerIdentitiesResp{
			CustomerId:       emailIdentity.CustomerId,
			Email:            &emailIdentity.Email,
			HasEmailConflict: &emailIdentity.HasConflict,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode response", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (h *CustomerHandler) handleCreateCustomerThChat(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var message ThChatReq
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
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	thread, chat, err := h.ths.CreateInboundThreadChat(
		ctx, workspace.WorkspaceId, customer.CustomerId, message.Message)
	if err != nil {
		slog.Error("failed to create thread chat", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var threadAssignee, egressMember, chatMember *ThMemberResp
	var inboundCustomer, chatCustomer *ThCustomerResp
	var inboundFirstSeqId, inboundLastSeqId *string

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

	threadCustomer := ThCustomerResp{
		CustomerId: thread.CustomerId,
		Name:       thread.CustomerName,
	}

	if thread.AssigneeId.Valid {
		threadAssignee = &ThMemberResp{
			MemberId: thread.AssigneeId.String,
			Name:     thread.AssigneeName.String,
		}
	}

	if thread.InboundMessage != nil {
		inboundCustomer = &ThCustomerResp{
			CustomerId: thread.InboundMessage.CustomerId,
			Name:       thread.InboundMessage.CustomerName,
		}
		inboundFirstSeqId = &thread.InboundMessage.FirstSeqId
		inboundLastSeqId = &thread.InboundMessage.LastSeqId
	}

	if thread.EgressMessageId.Valid {
		egressMember = &ThMemberResp{
			MemberId: thread.EgressMemberId.String,
			Name:     thread.EgressMemberName.String,
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

	resp := ThreadChatResp{
		ThreadId:          thread.ThreadId,
		Customer:          threadCustomer,
		Title:             thread.Title,
		Description:       thread.Description,
		Sequence:          thread.Sequence,
		Status:            thread.Status,
		Read:              thread.Read,
		Replied:           thread.Replied,
		Priority:          thread.Priority,
		Spam:              thread.Spam,
		Channel:           thread.Channel,
		PreviewText:       thread.PreviewText,
		Assignee:          threadAssignee,
		InboundFirstSeqId: inboundFirstSeqId,
		InboundLastSeqId:  inboundLastSeqId,
		InboundCustomer:   inboundCustomer,
		EgressFirstSeq:    thread.EgressFirstSeq,
		EgressLastSeq:     thread.EgressLastSeq,
		EgressMember:      egressMember,
		CreatedAt:         thread.CreatedAt,
		UpdatedAt:         thread.UpdatedAt,
		Chat:              chatResp,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetCustomerThChats(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	ctx := r.Context()

	threads, err := h.ths.ListCustomerThreadChats(ctx, customer.CustomerId, nil)
	if err != nil {
		slog.Error("failed to fetch thread chats", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadResp, 0, 100)
	for _, thread := range threads {
		var threadAssignee, egressMember *ThMemberResp
		var inboundCustomer *ThCustomerResp
		var inboundFirstSeqId, inboundLastSeqId *string

		threadCustomer := ThCustomerResp{
			CustomerId: thread.CustomerId,
			Name:       thread.CustomerName,
		}

		if thread.AssigneeId.Valid {
			threadAssignee = &ThMemberResp{
				MemberId: thread.AssigneeId.String,
				Name:     thread.AssigneeName.String,
			}
		}

		if thread.InboundMessage != nil {
			inboundCustomer = &ThCustomerResp{
				CustomerId: thread.InboundMessage.CustomerId,
				Name:       thread.InboundMessage.CustomerName,
			}
			inboundFirstSeqId = &thread.InboundMessage.FirstSeqId
			inboundLastSeqId = &thread.InboundMessage.LastSeqId
		}

		if thread.EgressMessageId.Valid {
			egressMember = &ThMemberResp{
				MemberId: thread.EgressMemberId.String,
				Name:     thread.EgressMemberName.String,
			}
		}

		resp := ThreadResp{
			ThreadId:          thread.ThreadId,
			Customer:          threadCustomer,
			Title:             thread.Title,
			Description:       thread.Description,
			Sequence:          thread.Sequence,
			Status:            thread.Status,
			Read:              thread.Read,
			Replied:           thread.Replied,
			Priority:          thread.Priority,
			Spam:              thread.Spam,
			Channel:           thread.Channel,
			PreviewText:       thread.PreviewText,
			Assignee:          threadAssignee,
			InboundFirstSeqId: inboundFirstSeqId,
			InboundLastSeqId:  inboundLastSeqId,
			InboundCustomer:   inboundCustomer,
			EgressFirstSeq:    thread.EgressFirstSeq,
			EgressLastSeq:     thread.EgressLastSeq,
			EgressMember:      egressMember,
			CreatedAt:         thread.CreatedAt,
			UpdatedAt:         thread.UpdatedAt,
		}
		items = append(items, resp)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode json", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCreateThChatMessage(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	threadId := r.PathValue("threadId")

	var message ThChatReq
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	channel := models.ThreadChannel{}.Chat()
	thread, err := h.ths.GetWorkspaceThread(ctx, customer.WorkspaceId, threadId, &channel)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch thread chat", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	chat, err := h.ths.AddInboundMessage(ctx, thread, message.Message)
	if err != nil {
		slog.Error("failed to create thread chat message", slog.Any("err", err))
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
		slog.Error("failed to encode json", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetThChatMessages(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	threadId := r.PathValue("threadId")
	ctx := r.Context()

	channel := models.ThreadChannel{}.Chat()
	thread, err := h.ths.GetWorkspaceThread(ctx, customer.WorkspaceId, threadId, &channel)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	chats, err := h.ths.ListThreadChatMessages(ctx, thread.ThreadId)
	if err != nil {
		slog.Error("failed to fetch thread chat messages", slog.Any("err", err))
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
		slog.Error("failed to encode json", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func NewServer(
	authService ports.CustomerAuthServicer,
	workspaceService ports.WorkspaceServicer,
	customerService ports.CustomerServicer,
	threadChatService ports.ThreadServicer,
) http.Handler {
	mux := http.NewServeMux()
	ch := NewCustomerHandler(workspaceService, customerService, threadChatService)

	mux.HandleFunc("GET /{$}", handleGetIndex)

	mux.HandleFunc("POST /widgets/{widgetId}/init/{$}", ch.handleGetOrCreateCustomer)
	mux.Handle("GET /widgets/{widgetId}/me/{$}", NewEnsureAuth(ch.handleGetCustomer, authService))
	mux.Handle("POST /widgets/{widgetId}/me/identities/{$}",
		NewEnsureAuth(ch.handleCustomerIdentities, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleCreateCustomerThChat, authService))
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleGetCustomerThChats, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleCreateThChatMessage, authService))
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleGetThChatMessages, authService))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
