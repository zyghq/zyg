package xhandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/zyghq/zyg"
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

// handleGetWidgetConfig returns the widget configuration.
// TODO: #71 get widget configuration from db/redis.
func (h *CustomerHandler) handleGetWidgetConfig(w http.ResponseWriter, _ *http.Request) {
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

// handleInitWidget initializes the configured widget with customer details.
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

	// from request body
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
	// (XXX):simplify parsing of traits?
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
		skipIdentityCheck = true // if the customer hash is provided, skip the identity check.
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

	// Shall use the latest claimed email.
	if !customer.IsVerified && !skipIdentityCheck {
		claimedEmail, err := h.cs.GetLatestValidClaimedEmail(ctx, customer.WorkspaceId, customer.CustomerId)
		// customer haven't provided claim able identity yet or is cleared.
		if errors.Is(err, services.ErrClaimedEmailNotFound) {
			customer.Email = models.NullString(nil)
		}
		// If there was an error checking claimed email by customer, ask again hence set nil
		// This might create duplicates, but that is fine.
		if errors.Is(err, services.ErrClaimedEmail) {
			customer.Email = models.NullString(nil)
		}
		// If claimed email is not empty string, set it to the customer.
		if claimedEmail != "" {
			customer.Email = models.NullString(&claimedEmail)
		}
	}

	resp := WidgetInitResp{
		Jwt:    jwt,
		Create: isCreated,
		CustomerResp: CustomerResp{
			CustomerId:      customer.CustomerId,
			Name:            customer.Name,
			AvatarUrl:       customer.AvatarUrl(),
			Email:           customer.Email,
			IsEmailVerified: false,
			IsEmailPrimary:  false,
			Phone:           customer.Phone,
			ExternalId:      customer.ExternalId,
			Role:            customer.Role,
			CreatedAt:       customer.CreatedAt,
			UpdatedAt:       customer.UpdatedAt,
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
	var emailOrClaimed sql.NullString
	var isEmailVerified, IsEmailPrimary bool

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

	if widget.WorkspaceId != customer.WorkspaceId {
		slog.Error("invalid workspace customer or widget configured")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if customer.IsVerified {
		emailOrClaimed = customer.Email
		isEmailVerified = true
		IsEmailPrimary = true
	} else {
		claimed, err := h.cs.GetLatestValidClaimedEmail(ctx, customer.WorkspaceId, customer.CustomerId)
		if err != nil {
			emailOrClaimed = models.NullString(nil)
		} else {
			emailOrClaimed = models.NullString(&claimed)
		}
	}
	resp := CustomerResp{
		CustomerId:      customer.CustomerId,
		Name:            customer.Name,
		AvatarUrl:       customer.AvatarUrl(),
		Email:           emailOrClaimed,
		IsEmailVerified: isEmailVerified,
		IsEmailPrimary:  IsEmailPrimary,
		Phone:           customer.Phone,
		ExternalId:      customer.ExternalId,
		Role:            customer.Role,
		CreatedAt:       customer.CreatedAt,
		UpdatedAt:       customer.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCreateThreadChat(
	w http.ResponseWriter, r *http.Request, customer *models.Customer) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp CreateThreadChatReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Add claimed email for verification, only if the customer's email is not verified yet.
	if !customer.IsVerified && reqp.Email != nil {
		redirectTo := zyg.ZygUrl() + "/?utm_source=zyg&utm_medium=kyc"
		if reqp.RedirectHost != nil {
			redirectTo = *reqp.RedirectHost + "/?utm_source=zyg&utm_medium=kyc"
		}

		// workspace secret key must exists.
		sk, err := h.ws.GetOrGenerateSecretKey(ctx, customer.WorkspaceId)
		if err != nil {
			slog.Error("failed to fetch workspace secret key", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		hasConflict, err := h.ws.DoesEmailConflict(ctx, customer.WorkspaceId, *reqp.Email)
		if err != nil {
			slog.Error("failed to check email conflict", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = h.cs.ClaimEmailForVerification(
			ctx, *customer, sk.Hmac, *reqp.Email, reqp.Name,
			hasConflict, reqp.Message, redirectTo,
		)
		if err != nil {
			slog.Error("failed to claim email for verification", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

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
		reqp.Message,
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

	var message MessageThreadReq
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

// (XXX) not a API endpoint, will be used for redirecting from mail verification URL.
// In all the cases we redirect to either the default target URL or the URL provided in the JWT token.
func (h *CustomerHandler) handleMailRedirectKyc(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	redirectTo := zyg.ZygUrl() + "/?utm_source=zyg&utm_medium=redirect"
	if t == "" {
		http.Redirect(w, r, redirectTo, http.StatusFound)
		return
	}

	ctx := r.Context()

	// Get valid claimed email by token.
	// Makes sure that the token exists in DB and not expired yet.
	// This is to have backend control on token handling.
	claimed, err := h.cs.GetValidClaimedEmailByToken(ctx, t)
	if errors.Is(err, services.ErrClaimedEmailExpired) {
		slog.Error("claimed email token is expired", slog.Any("err", err))
		http.Redirect(w, r, redirectTo, http.StatusFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch claimed email or does not exists", slog.Any("err", err))
		http.Redirect(w, r, redirectTo, http.StatusFound)
		return
	}

	// Get the secret key associated with the workspace.
	sk, err := h.ws.GetSecretKey(ctx, claimed.WorkspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace secret key", slog.Any("err", err))
		http.Redirect(w, r, redirectTo, http.StatusFound)
		return
	}

	// Verify the jwt token against the secret key.
	// NEVER trust the token before verifying it.
	j, err := h.cs.VerifyEmailVerificationToken([]byte(sk.Hmac), t)
	if err != nil {
		slog.Error("failed to verify claimed email token", slog.Any("err", err))
		return
	}

	// update the redirect URL to the URL provided in the JWT token.
	redirectTo = j.RedirectUrl + "/?utm_source=zyg&utm_medium=redirect"

	go func(claim models.KycMailJWTClaims) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// Fetch the workspace customer linked with the claimed email.
		// If the customer does not exists or failed, then return do nothing.
		role := models.Customer{}.Lead()
		claimedCustomer, err := h.ws.GetCustomer(ctx, claim.WorkspaceId, claim.Subject, &role)
		if err != nil {
			slog.Error("failed to fetch subject lead customer", slog.Any("err", err))
			return
		}

		// Check for actual customer associated with this email as primary.
		_, err = h.ws.GetCustomerByEmail(ctx, claim.WorkspaceId, claim.Email)
		// If no customer is associated with claimed email as primary, then trust the claimed email and customer
		// linking them together.
		if errors.Is(err, services.ErrCustomerNotFound) {
			claimedCustomer.Email = models.NullString(&claim.Email)
			claimedCustomer.IsVerified = true
			claimedCustomer.Role = models.Customer{}.Engaged()
			claimedCustomer, err = h.cs.UpdateCustomer(ctx, claimedCustomer)
			if err != nil {
				slog.Error("failed to update lead customer", slog.Any("err", err))
				return
			}
		}
		if err != nil {
			slog.Error("failed to fetch existing customer", slog.Any("err", err))
			return
		}
		// clear the claimed email identity linked with the lead customer.
		// err = h.cs.RemoveCustomerClaimedEmail(ctx, leadCustomer.WorkspaceId, leadCustomer.CustomerId, claim.Email)
		// if err != nil {
		// 	slog.Error("failed to remove email identity", slog.Any("err", err))
		// 	return
		// }
	}(j)

	http.Redirect(w, r, redirectTo, http.StatusFound)
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

	mux.HandleFunc("GET /mail/kyc/{$}", ch.handleMailRedirectKyc)

	mux.HandleFunc("GET /widgets/{widgetId}/config/{$}", ch.handleGetWidgetConfig)
	mux.HandleFunc("POST /widgets/{widgetId}/init/{$}", ch.handleInitWidget)

	mux.Handle("GET /widgets/{widgetId}/me/{$}",
		NewEnsureAuth(ch.handleGetCustomer, authService))

	mux.Handle("POST /widgets/{widgetId}/threads/chat/{$}",
		NewEnsureAuth(ch.handleCreateThreadChat, authService))
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
