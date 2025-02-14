package handler

import (
	"encoding/json"
	"errors"
	"github.com/getsentry/sentry-go"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/zyghq/zyg"

	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

type WorkspaceHandler struct {
	ws ports.WorkspaceServicer
	as ports.AccountServicer
	cs ports.CustomerServicer
}

func NewWorkspaceHandler(
	ws ports.WorkspaceServicer, as ports.AccountServicer, cs ports.CustomerServicer) *WorkspaceHandler {
	return &WorkspaceHandler{
		ws: ws,
		as: as,
		cs: cs,
	}
}

func (h *WorkspaceHandler) handleCreateWorkspace(
	w http.ResponseWriter, r *http.Request, account *models.Account) {
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

	workspace, err := h.as.CreateWorkspace(ctx, *account, reqp.Name)
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
	workspaces, err := h.as.ListAccountLinkedWorkspaces(ctx, account.AccountId)
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

func (h *WorkspaceHandler) handleGetWorkspace(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

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

func (h *WorkspaceHandler) handleUpdateWorkspace(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
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

	var hasUpdates bool
	// apply updates if any
	if reqp.Name != "" {
		hasUpdates = true
		workspace.Name = reqp.Name
	}

	// if no updates then respond with same workspace as is.
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
		slog.Error("failed to update workspace", slog.Any("err", err))
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

func (h *WorkspaceHandler) handleCreateWorkspaceLabel(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp NewLabelReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

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

func (h *WorkspaceHandler) handleUpdateWorkspaceLabel(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	labelId := r.PathValue("labelId")
	var reqp NewLabelReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		slog.Error("failed to decode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	label, err := h.ws.GetLabel(ctx, member.WorkspaceId, labelId)
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
			LabelId:   label.LabelId,
			Name:      label.Name,
			Icon:      label.Icon,
			CreatedAt: label.CreatedAt,
			UpdatedAt: label.UpdatedAt,
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
		LabelId:   label.LabelId,
		Name:      label.Name,
		Icon:      label.Icon,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceLabels(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	labels, err := h.ws.ListLabels(ctx, member.WorkspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace labels", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]LabelResp, 0, len(labels))
	for _, l := range labels {
		item := LabelResp{
			LabelId:   l.LabelId,
			Name:      l.Name,
			Icon:      l.Icon,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
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

func (h *WorkspaceHandler) handleGetWorkspaceLabel(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	labelId := r.PathValue("labelId")
	label, err := h.ws.GetLabel(ctx, member.WorkspaceId, labelId)
	if errors.Is(err, services.ErrLabelNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch workspace label", slog.Any("err", err))
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleCreateWorkspaceCustomer(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
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
		slog.Error("at least one of `externalId`, `email` or `phone` is required")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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

	var customer models.Customer
	var isCreated bool

	if externalId.Valid {
		customer, isCreated, err = h.ws.CreateCustomerWithExternalId(
			ctx, workspace.WorkspaceId,
			externalId.String,
			reqp.Name,
		)
		if err != nil {
			slog.Error("failed to fetch or create customer by externalId", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if email.Valid {
		customer, isCreated, err = h.ws.CreateCustomerWithEmail(
			ctx, workspace.WorkspaceId,
			email.String,
			reqp.IsEmailVerified,
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
			reqp.Name,
		)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		slog.Error("at least one of `externalId`, `email` or `phone` is required")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return

	}

	resp := CustomerResp{
		CustomerId:      customer.CustomerId,
		Name:            customer.Name,
		AvatarUrl:       customer.AvatarUrl(),
		IsEmailVerified: customer.IsEmailVerified,
		Role:            customer.Role,
		ExternalId:      customer.ExternalId,
		Email:           customer.Email,
		Phone:           customer.Phone,
		CreatedAt:       customer.CreatedAt,
		UpdatedAt:       customer.UpdatedAt,
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

func (h *WorkspaceHandler) handleGetWorkspaceMembership(
	w http.ResponseWriter, _ *http.Request, member *models.Member) {
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

func (h *WorkspaceHandler) handleGetWorkspaceMember(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	memberId := r.PathValue("memberId")
	otherMember, err := h.ws.GetMember(ctx, member.WorkspaceId, memberId)
	if errors.Is(err, services.ErrMemberNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch workspace otherMember", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := MemberResp{
		MemberId:  otherMember.MemberId,
		Name:      otherMember.Name,
		Role:      otherMember.Role,
		CreatedAt: otherMember.CreatedAt,
		UpdatedAt: otherMember.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceCustomers(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	customers, err := h.ws.ListCustomers(ctx, member.WorkspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace customers", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]CustomerResp, 0, len(customers))
	for _, c := range customers {
		items = append(items, CustomerResp{
			CustomerId:      c.CustomerId,
			ExternalId:      c.ExternalId,
			Email:           c.Email,
			Phone:           c.Phone,
			Name:            c.Name,
			IsEmailVerified: c.IsEmailVerified,
			Role:            c.Role,
			CreatedAt:       c.CreatedAt,
			UpdatedAt:       c.UpdatedAt,
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

func (h *WorkspaceHandler) handleGetWorkspaceMembers(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	members, err := h.ws.ListMembers(ctx, member.WorkspaceId)
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

func (h *WorkspaceHandler) handleGenerateSecretKey(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

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

	// generate secret key for the workspace for the given length.
	// required for hashing secret or jwt.
	sk, err := h.ws.GenerateWorkspaceSecret(ctx, workspace.WorkspaceId, zyg.DefaultSecretKeyLength)
	if err != nil {
		slog.Error("failed to generate workspace secret key", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := WorkspaceSecretResp{
		Hmac:      sk.Hmac,
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

func (h *WorkspaceHandler) handleGetWorkspaceSecretKey(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	sk, err := h.ws.GetSecretKey(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrSecretKeyNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch workspace secret key", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := WorkspaceSecretResp{
		Hmac:      sk.Hmac,
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

func (h *WorkspaceHandler) handleGetWidgets(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	widgets, err := h.ws.ListWidgets(ctx, member.WorkspaceId)
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

func (h *WorkspaceHandler) handlePostmarkGetMailServer(
	w http.ResponseWriter, r *http.Request, member *models.Member) {

	ctx := r.Context()

	hub := sentry.GetHubFromContext(ctx)

	workspace, err := h.ws.GetWorkspace(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	setting, err := h.ws.GetPostmarkMailServerSetting(ctx, workspace.WorkspaceId)
	if errors.Is(err, services.ErrPostmarkSettingNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch postmark mail server setting", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
	}
}

func (h *WorkspaceHandler) handlePostmarkCreateMailServer(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	hub := sentry.GetHubFromContext(ctx)

	var reqp CreatePostmarkMailServer
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspace, err := h.ws.GetWorkspace(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// split the domain from email
	domain := strings.Split(reqp.Email, "@")[1]
	if domain == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	setting, err := h.ws.PostmarkCreateMailServer(ctx, workspace.WorkspaceId, reqp.Email, domain)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create postmark mail server", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
	}
}

func (h *WorkspaceHandler) handlePostmarkMailServerAddDNS(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	hub := sentry.GetHubFromContext(ctx)

	var reqp AddPostmarkMailServerDNS

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspace, err := h.ws.GetWorkspace(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setting, err := h.ws.GetPostmarkMailServerSetting(ctx, workspace.WorkspaceId)
	if errors.Is(err, services.ErrPostmarkSettingNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch postmark mail server setting", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Add domain in Postmark
	setting, created, err := h.ws.PostmarkMailServerAddDomain(ctx, setting, reqp.Domain)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to add postmark mail server domain", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if created {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(setting); err != nil {
			slog.Error("failed to encode json", slog.Any("err", err))
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
	}
}

func (h *WorkspaceHandler) handlePostmarkMailServerVerifyDNS(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	hub := sentry.GetHubFromContext(ctx)

	workspace, err := h.ws.GetWorkspace(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setting, err := h.ws.GetPostmarkMailServerSetting(ctx, workspace.WorkspaceId)
	if errors.Is(err, services.ErrPostmarkSettingNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch postmark mail server setting", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if dns domain ID exists
	// dns domain ID must exist before verifying
	if setting.DNSDomainId == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Verify domain in Postmark
	setting, err = h.ws.PostmarkMailServerVerifyDomain(ctx, setting)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to verify postmark mail server domain", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
	}
}

// handlePostmarkUpdateMailServer updates the Postmark mail server settings of a workspace based on the provided HTTP request.
func (h *WorkspaceHandler) handlePostmarkUpdateMailServer(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	hub := sentry.GetHubFromContext(ctx)

	var reqp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspace, err := h.ws.GetWorkspace(ctx, member.WorkspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	setting, err := h.ws.GetPostmarkMailServerSetting(ctx, workspace.WorkspaceId)
	if errors.Is(err, services.ErrPostmarkSettingNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch postmark mail server setting", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var isModified bool
	var blockEnabling bool

	// fields for updates
	fields := make([]string, 0, len(reqp))

	if hasForwardingEnabled, found := reqp["hasForwardingEnabled"]; found {
		isModified = true
		if hasForwardingEnabled == nil {
			setting.HasForwardingEnabled = false
			fields = append(fields, "hasForwardingEnabled")
		} else {
			f := hasForwardingEnabled.(bool)
			setting.HasForwardingEnabled = f
			fields = append(fields, "hasForwardingEnabled")
		}
	}

	if !setting.HasForwardingEnabled {
		blockEnabling = true
		setting.IsEnabled = false
		fields = append(fields, "enabled")
	}

	// If the forwarding is disabled, then enable flag cannot be updated and IsEnabled flag is set to false
	if enabled, found := reqp["enabled"]; found && !blockEnabling {
		isModified = true
		if enabled == nil {
			setting.IsEnabled = false
			fields = append(fields, "enabled")
		} else {
			f := enabled.(bool)
			if f && setting.HasForwardingEnabled {

			}
			setting.IsEnabled = f
			fields = append(fields, "enabled")
		}
	}

	// persist if modified
	if isModified {
		setting, err = h.ws.PostmarkMailServerUpdate(ctx, setting, fields)
		if err != nil {
			hub.CaptureException(err)
			slog.Error("failed to modify postmark mail server setting", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
	}
}
