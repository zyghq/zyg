package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

type CustomerHandler struct {
	ws ports.WorkspaceServicer
	cs ports.CustomerServicer
}

func NewCustomerHandler(
	ws ports.WorkspaceServicer, cs ports.CustomerServicer) *CustomerHandler {
	return &CustomerHandler{ws, cs}
}

type customerIdentifiers struct {
	CustomerId *string
	ExternalId *string
	Email      *string
}

// validateCustomerIdentifiers ensures at least one customer identifier is provided (CustomerId, ExternalId, Email).
func validateCustomerIdentifiers(customer customerIdentifiers) error {
	if (customer.CustomerId != nil && *customer.CustomerId != "") ||
		(customer.ExternalId != nil && *customer.ExternalId != "") ||
		(customer.Email != nil && *customer.Email != "") {
		return nil
	}
	return errors.New("at least one customer identifier must be provided")
}

// handleCreateCustomerEvent processes the HTTP request to create a customer event.
// It validates the input, checks for customer existence, and appends the event to the customer's record.
func (h *CustomerHandler) handleCreateCustomerEvent(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	var reqp CustomerEventReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check for customer identifier
	ci := reqp.Customer
	identifiers := customerIdentifiers{
		CustomerId: ci.CustomerId,
		ExternalId: ci.ExternalId,
		Email:      ci.Email,
	}
	err = validateCustomerIdentifiers(identifiers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspace, err := h.ws.GetWorkspace(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Add event for customer identified by customer ID.
	if ci.CustomerId != nil {
		customer, err := h.ws.GetCustomer(ctx, workspace.WorkspaceId, *ci.CustomerId, nil)
		if errors.Is(err, services.ErrCustomerNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if err != nil {
			slog.Error("failed to fetch customer", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newEvent, err := models.NewEvent(
			reqp.Title,
			models.SetEventCustomer(customer.AsCustomerActor()),
			models.SetEventSeverity(reqp.Severity),
			models.SetEventTimestampFromStr(reqp.Timestamp),
			models.WithEventComponents(reqp.Components),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, err := h.cs.AddEvent(ctx, *newEvent)
		if err != nil {
			slog.Error("failed to add customer event", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp := CustomerEventAddedResp{
			EventID:   event.EventID,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if ci.ExternalId != nil {
		customer, _, err := h.ws.CreateCustomerWithExternalId(ctx, workspace.WorkspaceId, *ci.ExternalId, *ci.Name)
		if err != nil {
			slog.Error("failed to fetch or create customer by externalId", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newEvent, err := models.NewEvent(
			reqp.Title,
			models.SetEventCustomer(customer.AsCustomerActor()),
			models.SetEventSeverity(reqp.Severity),
			models.SetEventTimestampFromStr(reqp.Timestamp),
			models.WithEventComponents(reqp.Components),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, err := h.cs.AddEvent(ctx, *newEvent)
		if err != nil {
			slog.Error("failed to add customer event", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp := CustomerEventAddedResp{
			EventID:   event.EventID,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if ci.Email != nil {
		// check email verification flag
		var isEmailVerified bool
		if ci.IsEmailVerified != nil {
			isEmailVerified = *ci.IsEmailVerified
		}
		customer, _, err := h.ws.CreateCustomerWithEmail(
			ctx, workspace.WorkspaceId, *ci.Email, isEmailVerified, *ci.Name)
		if err != nil {
			slog.Error("failed to fetch or create customer by email", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newEvent, err := models.NewEvent(
			reqp.Title,
			models.SetEventCustomer(customer.AsCustomerActor()),
			models.SetEventSeverity(reqp.Severity),
			models.SetEventTimestampFromStr(reqp.Timestamp),
			models.WithEventComponents(reqp.Components),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, err := h.cs.AddEvent(ctx, *newEvent)
		if err != nil {
			slog.Error("failed to add customer event", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp := CustomerEventAddedResp{
			EventID:   event.EventID,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		slog.Error("at least one of `customerId`, `externalId` or `email` is required")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (h *CustomerHandler) handleGetCustomerEvents(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	customerId := r.PathValue("customerId")
	ctx := r.Context()

	items := make([]CustomerEventResp, 0, 11)

	// check if the customer does exist in workspace.
	customer, err := h.ws.GetCustomer(ctx, member.WorkspaceId, customerId, nil)
	if errors.Is(err, services.ErrCustomerNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(items); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		slog.Error("failed to fetch customer", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetch list of customer events.
	events, err := h.cs.ListEvents(ctx, customer.CustomerId)
	if err != nil {
		slog.Error("failed to list customer events", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, e := range events {
		resp := CustomerEventResp{}.NewResponse(&e)
		items = append(items, resp)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
