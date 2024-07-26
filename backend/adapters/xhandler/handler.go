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
	anonName := models.Customer{}.AnonName()
	customerName := models.NullString(&anonName)

	if payload.Traits != nil {
		if payload.Traits.Name != nil {
			customerName = models.NullString(payload.Traits.Name)
		} else {
			if payload.Traits.FirstName != nil || payload.Traits.LastName != nil {
				n := ""
				if payload.Traits.FirstName != nil {
					n += *payload.Traits.FirstName
				}
				if payload.Traits.LastName != nil {
					n += " " + *payload.Traits.LastName
				}
				customerName = models.NullString(&n)
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

func (h *CustomerHandler) handleCreateCustomerThChat(
	w http.ResponseWriter, r *http.Request, customer *models.Customer,
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

	th := models.ThreadChat{
		WorkspaceId:  workspace.WorkspaceId,
		CustomerId:   customer.CustomerId,
		CustomerName: customer.Name,
	}
	th, thm, err := h.ths.CreateThreadWithMessage(ctx, th, message.Message)
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
		Priority:     th.Priority,
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
				"check the json encoding defn",
			slog.String("threadChatId", th.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleCreateThChatMessage(
	w http.ResponseWriter, r *http.Request, customer *models.Customer,
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

	thread, err := h.ths.GetThread(ctx, workspace.WorkspaceId, threadId)

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

	thm, err := h.ths.AddCustomerMessageToThread(ctx, thread, customer, message.Message)

	if err != nil {
		slog.Error(
			"failed to create thread chat message for customer "+
				"something went wrong",
			slog.String("threadChatId", thread.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var threadAssigneeRepr *ThMemberRespPayload
	var msgCustomerRepr *ThCustomerRespPayload
	var msgMemberRepr *ThMemberRespPayload

	// for thread
	threadCustomerRepr := ThCustomerRespPayload{
		CustomerId: thread.CustomerId,
		Name:       thread.CustomerName,
	}

	// for thread
	if thread.AssigneeId.Valid {
		threadAssigneeRepr = &ThMemberRespPayload{
			MemberId: thread.AssigneeId.String,
			Name:     thread.AssigneeName,
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
		ThreadChatId:        thread.ThreadChatId,
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
		ThreadChatId: thread.ThreadChatId,
		Sequence:     thread.Sequence,
		Status:       thread.Status,
		Read:         thread.Read,
		Replied:      thread.Replied,
		Priority:     thread.Priority,
		Customer:     threadCustomerRepr,
		Assignee:     threadAssigneeRepr,
		CreatedAt:    thread.CreatedAt,
		UpdatedAt:    thread.UpdatedAt,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat message to json "+
				"check the json encoding defn",
			slog.String("threadChatId", thread.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetThChatMesssages(
	w http.ResponseWriter, r *http.Request, customer *models.Customer,
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

	thread, err := h.ths.GetThread(ctx, workspace.WorkspaceId, threadId)

	if errors.Is(err, services.ErrThreadChatNotFound) {
		slog.Warn(
			"thread chat not found or does not exist for customer",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	results, err := h.ths.ListThreadMessages(ctx, thread.ThreadChatId)

	if err != nil {
		slog.Error(
			"failed to get list of thread chat messages for customer "+
				"something went wrong",
			slog.String("threadChatId", thread.ThreadChatId),
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
			ThreadChatId:        thread.ThreadChatId,
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
		CustomerId: thread.CustomerId,
		Name:       thread.CustomerName,
	}

	// for thread
	if thread.AssigneeId.Valid {
		threadAssigneeRepr = &ThMemberRespPayload{
			MemberId: thread.AssigneeId.String,
			Name:     thread.AssigneeName,
		}
	}

	resp := ThChatRespPayload{
		ThreadChatId: thread.ThreadChatId,
		Sequence:     thread.Sequence,
		Status:       thread.Status,
		Read:         thread.Read,
		Replied:      thread.Replied,
		Priority:     thread.Priority,
		Customer:     threadCustomerRepr,
		Assignee:     threadAssigneeRepr,
		CreatedAt:    thread.CreatedAt,
		UpdatedAt:    thread.UpdatedAt,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat messages to json "+
				"check the json encoding defn",
			slog.String("threadChatId", thread.ThreadChatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CustomerHandler) handleGetCustomerThChats(
	w http.ResponseWriter, r *http.Request, customer *models.Customer,
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

	results, err := h.ths.ListCustomerThreads(ctx, workspace.WorkspaceId, customer.CustomerId)

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
	for _, th := range results {
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
			Priority:     th.ThreadChat.Priority,
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
				"check the json encoding defn",
			slog.String("customerId", customer.CustomerId),
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
