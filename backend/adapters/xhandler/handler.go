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

func (h *CustomerHandler) handleGetWidgetConfig(w http.ResponseWriter, _ *http.Request) {
	// TODO: probably get widget configuration from db/redis cache.
	resp := WidgetConfig{
		DomainsOnly:    false,
		Domains:        []string{},
		BubblePosition: "right",
		HeaderColor:    "#9370DB",
		ProfilePicture: "",
		IconColor:      "#ffff",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleInitWidget(w http.ResponseWriter, r *http.Request) {
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

	// Get or generate a new secret key for the workspace.
	sk, err := h.ws.GetOrGenerateSecretKey(ctx, widget.WorkspaceId)
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

	var isVerified bool
	if reqp.IsVerified != nil {
		isVerified = *reqp.IsVerified
	}

	sessionId := models.NullString(reqp.SessionId)
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

	var skipIdentityCheck bool // don't want redundant identity check if a valid identity is provided.

	// The client provides customer hash, to verify the end customer.
	//
	// Creates the customer if not found.
	// Priority is externalId, then email, then phone.
	// Other values are ignored.
	if customerHash.Valid {
		skipIdentityCheck = true // if the customer has been provided, skip the identity check.
		if customerExternalId.Valid {
			if h.cs.VerifyExternalId(sk.Hmac, customerHash.String, customerExternalId.String) {
				customer, isCreated, err = h.ws.CreateCustomerWithExternalId(
					ctx, widget.WorkspaceId,
					customerExternalId.String,
					isVerified,
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
					isVerified,
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
					isVerified,
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
	} else if sessionId.Valid {
		// Check if the session with the session ID is already created and verify the session.
		// Otherwise, create a new session with an anonymous customer.
		customer, err = h.ws.ValidateWidgetSession(ctx, sk.Hmac, widget.WidgetId, sessionId.String)
		if errors.Is(err, services.ErrWidgetSessionInvalid) {
			customer, isCreated, err = h.ws.CreateWidgetSession(
				ctx, sk.Hmac, widget.WorkspaceId, widget.WidgetId, sessionId.String, customerName)
		}
		if errors.Is(err, services.ErrWidgetSession) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		// force the client to provide the anonymousId.
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Generate JWT token for the customer and secret key.
	jwt, err := h.cs.GenerateCustomerJwt(customer, sk.Hmac)
	if err != nil {
		slog.Error("failed to make jwt token with error", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if the customer needs to provide additional identities.
	// Ask for identities if the customer doesn't have a natural identity and isn't verified yet.
	//
	// Note: We can configure the workspace settings to have this enabled or disabled?
	RequireIdentities := make([]string, 0, 3) // email, phone, externalId
	if !customer.HasNaturalIdentity() && !customer.IsVerified && !skipIdentityCheck {
		hasEmailIdentity, err := h.cs.HasProvidedEmailIdentity(ctx, customer.CustomerId)
		if err != nil {
			slog.Error("Failed to check customer email identity", slog.Any("err", err))
			// If there's an error, we'll ask for email to be safe
			RequireIdentities = append(RequireIdentities, "email")
		} else if !hasEmailIdentity {
			RequireIdentities = append(RequireIdentities, "email")
		}
		// Here you could add checks for other identity types if needed,
		// For example, phone, external ID, etc.
	}

	resp := WidgetInitResp{
		Jwt:    jwt,
		Create: isCreated,
		CustomerResp: CustomerResp{
			CustomerId:        customer.CustomerId,
			Name:              customer.Name,
			AvatarUrl:         customer.AvatarUrl(),
			Email:             customer.Email,
			Phone:             customer.Phone,
			ExternalId:        customer.ExternalId,
			IsVerified:        customer.IsVerified,
			Role:              customer.Role,
			CreatedAt:         customer.CreatedAt,
			UpdatedAt:         customer.UpdatedAt,
			RequireIdentities: RequireIdentities,
		},
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

func (h *CustomerHandler) handleGetCustomer(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	var externalId, email, phone string

	if customer.Email.Valid {
		email = customer.Email.String
	}
	if customer.Phone.Valid {
		phone = customer.Phone.String
	}
	if customer.ExternalId.Valid {
		externalId = customer.ExternalId.String
	}

	ctx := r.Context()

	// Check if the customer needs to provide additional identities.
	// Ask for identities if the customer doesn't have a natural identity and isn't verified yet.
	//
	// Note: We can configure the workspace settings to have this enabled or disabled?
	RequireIdentities := make([]string, 0, 3) // email, phone, externalId
	if !customer.HasNaturalIdentity() && !customer.IsVerified {
		hasEmailIdentity, err := h.cs.HasProvidedEmailIdentity(ctx, customer.CustomerId)
		if err != nil {
			slog.Error("Failed to check customer email identity", slog.Any("err", err))
			// If there's an error, we'll ask for email to be safe
			RequireIdentities = append(RequireIdentities, "email")
		} else if !hasEmailIdentity {
			RequireIdentities = append(RequireIdentities, "email")
		}
		// Here you could add checks for other identity types if needed,
		// For example, phone, external ID, etc.
	}

	resp := CustomerResp{
		CustomerId:        customer.CustomerId,
		Name:              customer.Name,
		AvatarUrl:         customer.AvatarUrl(),
		Email:             models.NullString(&email),
		Phone:             models.NullString(&phone),
		ExternalId:        models.NullString(&externalId),
		IsVerified:        customer.IsVerified,
		Role:              customer.Role,
		CreatedAt:         customer.CreatedAt,
		UpdatedAt:         customer.UpdatedAt,
		RequireIdentities: RequireIdentities,
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

	// Allow modification only if the customer is not verified
	// from the external widget.
	// Don't want to allow modification of verified customer externally.
	if customer.IsVerified {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Add Email identity if provided for now, as it is more common.
	// Add other identities if needed.
	//
	// Intentionally don't want to modify the primary email identifier from external widget api.
	// Use other api that do direct primary email modification.
	if reqp.Email != nil {
		// Check if email already exists for a customer, if then there is a conflict.
		// Having a conflict doesn't mean we can't add multiple email identities.
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

		// Convert if the Customer is `visitor` to a `lead`.
		if customer.IsVisitor() {
			customer.Role = customer.Lead()
			updatedCustomer, err := h.cs.UpdateCustomer(ctx, *customer)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			customer = &updatedCustomer // updated customer
		}

		// Respond with customer details, with empty requireIdentities.
		resp := CustomerResp{
			CustomerId:        emailIdentity.CustomerId,
			ExternalId:        customer.ExternalId,
			Email:             customer.Email,
			Phone:             customer.Phone,
			Name:              customer.Name,
			AvatarUrl:         customer.AvatarUrl(),
			IsVerified:        customer.IsVerified,
			Role:              customer.Role,
			CreatedAt:         customer.CreatedAt,
			UpdatedAt:         customer.UpdatedAt,
			RequireIdentities: []string{},
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

func (h *CustomerHandler) handleCreateCustomerThreadChat(
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

	// return the system member for the workspace
	member, err := h.ws.GetSystemMember(ctx, customer.WorkspaceId)
	if errors.Is(err, services.ErrMemberNotFound) {
		// system member isn't found, create a new one.
		member, err = h.ws.CreateNewSystemMember(ctx, customer.WorkspaceId)
		// error creating system member, return server error.
		if err != nil {
			slog.Error("failed to create system member", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	if err != nil {
		slog.Error("failed to fetch system member", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	thread, chat, err := h.ths.CreateNewInboundThreadChat(
		ctx, customer.WorkspaceId, *customer, member.AsMemberActor(),
		message.Message,
	)
	if err != nil {
		slog.Error("failed to create thread chat", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := ThreadChatResp{}.NewResponse(&thread, &chat)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetCustomerThreadChats(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	ctx := r.Context()

	threads, err := h.ths.ListCustomerThreadChats(ctx, customer.CustomerId)
	if err != nil {
		slog.Error("failed to fetch thread chats", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadResp, 0, 100)
	for _, thread := range threads {
		resp := ThreadResp{}.NewResponse(&thread)
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

func (h *CustomerHandler) handleCreateThreadChatMessage(
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

	chat, err := h.ths.CreateInboundChatMessage(ctx, thread, message.Message)
	if err != nil {
		slog.Error("failed to create thread chat message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var chatCustomer *CustomerActorResp
	var chatMember *MemberActorResp

	// for chat - either of them
	if chat.CustomerId.Valid {
		chatCustomer = &CustomerActorResp{
			CustomerId: chat.CustomerId.String,
			Name:       chat.CustomerName.String,
		}
	} else if chat.MemberId.Valid {
		chatMember = &MemberActorResp{
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

func (h *CustomerHandler) handleGetThreadChatMessages(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	threadId := r.PathValue("threadId")
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

	chats, err := h.ths.ListThreadChatMessages(ctx, thread.ThreadId)
	if err != nil {
		slog.Error("failed to fetch thread chat messages", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	messages := make([]ChatResp, 0, 100)
	for _, chat := range chats {
		var chatCustomer *CustomerActorResp
		var chatMember *MemberActorResp
		if chat.CustomerId.Valid {
			chatCustomer = &CustomerActorResp{
				CustomerId: chat.CustomerId.String,
				Name:       chat.CustomerName.String,
			}
		} else if chat.MemberId.Valid {
			chatMember = &MemberActorResp{
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

	mux.HandleFunc("GET /widgets/{widgetId}/config/{$}", ch.handleGetWidgetConfig)
	mux.HandleFunc("POST /widgets/{widgetId}/init/{$}", ch.handleInitWidget)

	mux.Handle("GET /widgets/{widgetId}/me/{$}",
		NewEnsureAuth(ch.handleGetCustomer, authService))
	mux.Handle("POST /widgets/{widgetId}/me/identities/{$}",
		NewEnsureAuth(ch.handleCustomerIdentities, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleCreateCustomerThreadChat, authService))
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleGetCustomerThreadChats, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleCreateThreadChatMessage, authService))
	mux.Handle("GET /widgets/{widgetId}/threads/chat/{threadId}/messages/{$}",
		NewEnsureAuth(ch.handleGetThreadChatMessages, authService))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	handler := LoggingMiddleware(c.Handler(mux))

	return handler
}
