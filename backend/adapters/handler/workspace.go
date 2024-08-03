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

type WorkspaceHandler struct {
	ws ports.WorkspaceServicer
	cs ports.CustomerServicer
}

func NewWorkspaceHandler(ws ports.WorkspaceServicer, cs ports.CustomerServicer) *WorkspaceHandler {
	return &WorkspaceHandler{ws: ws, cs: cs}
}

func (h *WorkspaceHandler) handleCreateWorkspace(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	var reqp WorkspaceReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspace, err := h.ws.CreateWorkspace(ctx, account.AccountId, account.Name, reqp.Name)
	if err != nil {
		slog.Error("failed to create workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := WorkspaceResp{
		WorkspaceId: workspace.WorkspaceId,
		Name:        workspace.Name,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaces(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()
	workspaces, err := h.ws.ListAccountLinkedWorkspaces(ctx, account.AccountId)
	if err != nil {
		slog.Error("failed fetch workspaces", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]WorkspaceResp, 0, len(workspaces))
	for _, w := range workspaces {
		item := WorkspaceResp{
			WorkspaceId: w.WorkspaceId,
			Name:        w.Name,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspace(w http.ResponseWriter, r *http.Request, account *models.Account) {

	ctx := r.Context()
	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetAccountLinkedWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := WorkspaceResp{
		WorkspaceId: workspace.WorkspaceId,
		Name:        workspace.Name,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// keeping it simple for now
// future shall handle more workspace updates.
func (h *WorkspaceHandler) handleUpdateWorkspace(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp WorkspaceReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetAccountLinkedWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var hasUpdates bool

	// apply updates if any
	if reqp.Name != "" {
		hasUpdates = true
		workspace.Name = reqp.Name
	}

	// short circuit if no updates
	if !hasUpdates {
		resp := WorkspaceResp{
			WorkspaceId: workspace.WorkspaceId,
			Name:        workspace.Name,
			CreatedAt:   workspace.CreatedAt,
			UpdatedAt:   workspace.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	workspace, err = h.ws.UpdateWorkspace(ctx, workspace)
	if err != nil {
		slog.Error("failed to update workspce", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workspace); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetOrCreateWorkspaceLabel(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")

	var reqp NewLabelReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspace, err := h.ws.GetAccountLinkedWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	label, isCreated, err := h.ws.CreateLabel(ctx, workspace.WorkspaceId, reqp.Name, reqp.Icon)
	if err != nil {
		slog.Error("failed to create label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := LabelResp{
		LabelId:   label.LabelId,
		Name:      label.Name,
		Icon:      label.Icon,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}

	if isCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (h *WorkspaceHandler) handleUpdateWorkspaceLabel(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")
	labelId := r.PathValue("labelId")
	var reqp NewLabelReq

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	label, err := h.ws.GetLabel(ctx, workspaceId, labelId)
	if errors.Is(err, services.ErrLabelNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var hasUpdates bool
	if reqp.Name != "" {
		hasUpdates = true
		label.Name = reqp.Name
	}
	if reqp.Icon != "" {
		hasUpdates = true
		label.Icon = reqp.Icon
	}

	if !hasUpdates {
		resp := LabelResp{
			LabelId:     label.LabelId,
			WorkspaceId: label.WorkspaceId,
			Name:        label.Name,
			Icon:        label.Icon,
			CreatedAt:   label.CreatedAt,
			UpdatedAt:   label.UpdatedAt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	label, err = h.ws.UpdateLabel(ctx, label)
	if err != nil {
		slog.Error("failed to update workspace label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := LabelResp{
		LabelId:     label.LabelId,
		WorkspaceId: label.WorkspaceId,
		Name:        label.Name,
		Icon:        label.Icon,
		CreatedAt:   label.CreatedAt,
		UpdatedAt:   label.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceLabels(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	labels, err := h.ws.ListLabels(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace labels", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]LabelResp, 0, len(labels))

	for _, l := range labels {
		item := LabelResp{
			LabelId:     l.LabelId,
			WorkspaceId: l.WorkspaceId,
			Name:        l.Name,
			Icon:        l.Icon,
			CreatedAt:   l.CreatedAt,
			UpdatedAt:   l.UpdatedAt,
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleCreateWorkspaceCustomer(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	var reqp CreateCustomerReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	externalId := models.NullString(reqp.ExternalId)
	email := models.NullString(reqp.Email)
	phone := models.NullString(reqp.Phone)
	if !externalId.Valid && !email.Valid && !phone.Valid {
		slog.Error("atleast one of `externalId`, `email` or `phone` is required")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetAccountLinkedWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var customer models.Customer
	var isCreated bool

	if externalId.Valid {
		customer, isCreated, err = h.ws.CreateCustomerWithExternalId(
			ctx, workspace.WorkspaceId,
			externalId.String,
			true,
			reqp.Name,
		)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if email.Valid {
		customer, isCreated, err = h.ws.CreateCustomerWithEmail(
			ctx, workspace.WorkspaceId,
			email.String,
			true,
			reqp.Name,
		)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if phone.Valid {
		customer, isCreated, err = h.ws.CreateCustomerWithPhone(
			ctx, workspace.WorkspaceId,
			phone.String,
			true,
			reqp.Name,
		)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		slog.Error("atleast one of `externalId`, `email` or `phone` is required")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return

	}

	resp := CustomerResp{
		CustomerId: customer.CustomerId,
		Name:       customer.Name,
		IsVerified: customer.IsVerified,
		Role:       customer.Role,
		ExternalId: customer.ExternalId,
		Email:      customer.Email,
		Phone:      customer.Phone,
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
	}

	if isCreated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceMembership(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	member, err := h.ws.GetAccountLinkedMember(ctx, workspaceId, account.AccountId)
	if err != nil {
		slog.Error("failed to fetch workspace membership", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := MemberResp{
		MemberId:  member.MemberId,
		Name:      member.Name,
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceMember(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	memberId := r.PathValue("memberId")

	member, err := h.ws.GetMember(ctx, workspaceId, memberId)
	if errors.Is(err, services.ErrMemberNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace member", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := MemberResp{
		MemberId:  member.MemberId,
		Name:      member.Name,
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceCustomers(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	customers, err := h.ws.ListCustomers(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace customers", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]CustomerResp, 0, len(customers))
	for _, c := range customers {
		items = append(items, CustomerResp{
			CustomerId: c.CustomerId,
			ExternalId: c.ExternalId,
			Email:      c.Email,
			Phone:      c.Phone,
			Name:       c.Name,
			IsVerified: c.IsVerified,
			Role:       c.Role,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceMembers(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	members, err := h.ws.ListMembers(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace members", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]MemberResp, 0, len(members))
	for _, m := range members {
		item := MemberResp{
			MemberId:  m.MemberId,
			Name:      m.Name,
			Role:      m.Role,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGenerateSecretKey(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetAccountLinkedWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sk, err := h.ws.GenerateSecretKey(ctx, workspace.WorkspaceId, 64)
	if err != nil {
		slog.Error("failed to generate workspace secret key", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := SKResp{
		SecretKey: sk.SecretKey,
		CreatedAt: sk.CreatedAt,
		UpdatedAt: sk.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceSecretKey(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	sk, err := h.ws.GetSecretKey(ctx, workspaceId)
	if errors.Is(err, services.ErrSecretKeyNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace secret key", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := SKResp{
		SecretKey: sk.SecretKey,
		CreatedAt: sk.CreatedAt,
		UpdatedAt: sk.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWidgets(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	ctx := r.Context()

	widgets, err := h.ws.ListWidgets(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch widgets", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make([]WidgetResp, 0, 100)
	for _, widget := range widgets {
		resp := WidgetResp{
			WidgetId:      widget.WidgetId,
			Name:          widget.Name,
			Configuration: widget.Configuration,
			CreatedAt:     widget.CreatedAt,
			UpdatedAt:     widget.UpdatedAt,
		}
		response = append(response, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
