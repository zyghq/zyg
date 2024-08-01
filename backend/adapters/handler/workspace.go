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
	workspaces, err := h.ws.ListAccountWorkspaces(ctx, account.AccountId)
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

	workspace, err := h.ws.GetLinkedWorkspaceMember(ctx, account.AccountId, workspaceId)
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

	workspace, err := h.ws.GetLinkedWorkspaceMember(ctx, account.AccountId, workspaceId)
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

	workspace, err := h.ws.GetLinkedWorkspaceMember(ctx, account.AccountId, workspaceId)
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

	label, err := h.ws.GetWorkspaceLabel(ctx, workspaceId, labelId)
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

	label, err = h.ws.UpdateWorkspaceLabel(ctx, label.WorkspaceId, label)
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

	labels, err := h.ws.ListWorkspaceLabels(ctx, workspaceId)
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

// TODO: deprecate this.
// func (h *WorkspaceHandler) handleIssueCustomerToken(w http.ResponseWriter, r *http.Request, account *models.Account) {
// 	defer func(r io.ReadCloser) {
// 		_, _ = io.Copy(io.Discard, r)
// 		_ = r.Close()
// 	}(r.Body)

// 	ctx := r.Context()

// 	var rb CustomerTIReqPayload
// 	err := json.NewDecoder(r.Body).Decode(&rb)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	externalId := models.NullString(rb.Customer.ExternalId)
// 	email := models.NullString(rb.Customer.Email)
// 	phone := models.NullString(rb.Customer.Phone)
// 	name := models.NullString(rb.Customer.Name)
// 	if !externalId.Valid && !email.Valid && !phone.Valid {
// 		slog.Error("at least one of `externalId`, `email` or `phone` is required")
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	workspaceId := r.PathValue("workspaceId")

// 	workspace, err := h.ws.GetLinkedWorkspaceMember(ctx, account.AccountId, workspaceId)

// 	if errors.Is(err, services.ErrWorkspaceNotFound) {
// 		slog.Warn(
// 			"account workspace not found or does not exist",
// 			slog.String("workspaceId", workspaceId),
// 		)
// 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	if err != nil {
// 		slog.Error(
// 			"failed to get account workspace or does not exist "+
// 				"something went wrong",
// 			slog.String("workspaceId", workspaceId),
// 		)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}

// 	customer := models.Customer{
// 		WorkspaceId: workspace.WorkspaceId,
// 		ExternalId:  externalId,
// 		Email:       email,
// 		Phone:       phone,
// 		Name:        name,
// 		IsVerified:  true,
// 	}

// 	var isCreated bool
// 	var resp CustomerTIRespPayload

// 	if rb.Create {
// 		if rb.CreateBy == nil {
// 			slog.Error("requires `createBy` when `create` is enabled")
// 			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 			return
// 		}
// 		createBy := *rb.CreateBy
// 		slog.Info("create Customer if does not exists", slog.String("createBy", createBy))
// 		switch createBy {
// 		case "email":
// 			if !customer.Email.Valid {
// 				slog.Error("email is required for createBy email")
// 				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 				return
// 			}

// 			slog.Info("create Customer by email")
// 			customer, isCreated, err = h.ws.CreateCustomerWithEmail(ctx, customer)

// 			if err != nil {
// 				slog.Error(
// 					"failed to get or create Workspace Customer by email" +
// 						"something went wrong",
// 				)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 				return
// 			}
// 		case "phone":
// 			if !customer.Phone.Valid {
// 				slog.Error("phone is required for createBy phone")
// 				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 				return
// 			}

// 			slog.Info("create Customer by phone")
// 			customer, isCreated, err = h.ws.CreateCustomerWithPhone(ctx, customer)

// 			if err != nil {
// 				slog.Error(
// 					"failed to get or create Workspace Customer by phone " +
// 						"something went wrong",
// 				)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 				return
// 			}
// 		case "externalId":
// 			if !customer.ExternalId.Valid {
// 				slog.Error("externalId is required for createBy externalId")
// 				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 				return
// 			}

// 			slog.Info("create Customer by externalId")
// 			customer, isCreated, err = h.ws.CreateCustomerWithExternalId(ctx, customer)

// 			if err != nil {
// 				slog.Error(
// 					"failed to get or create Workspace Customer by externalId" +
// 						"something went wrong",
// 				)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 				return
// 			}
// 		default:
// 			slog.Warn("unsupported createBy value")
// 			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 			return
// 		}
// 	} else {
// 		slog.Info("based on identifiers check for Customer in Workspace", slog.String("workspaceId", workspaceId))
// 		if customer.ExternalId.Valid {
// 			slog.Info("get customer by externalId")
// 			customer, err = h.cs.GetCustomerByExternalId(ctx, workspace.WorkspaceId, customer.ExternalId.String)
// 			if errors.Is(err, services.ErrCustomerNotFound) {
// 				slog.Warn(
// 					"Customer not found by externalId" +
// 						"perhaps the customer is not created or is not returned",
// 				)
// 				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 				return
// 			}

// 			if err != nil {
// 				slog.Error(
// 					"failed to get Workspace Customer by externalId" +
// 						"something went wrong",
// 				)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 				return
// 			}
// 		} else if customer.Email.Valid {
// 			slog.Info("get customer by email")
// 			customer, err = h.cs.GetCustomerByEmail(ctx, workspace.WorkspaceId, customer.Email.String)

// 			if errors.Is(err, services.ErrCustomerNotFound) {
// 				slog.Warn(
// 					"Customer not found by email" +
// 						"perhaps the customer is not created or is not returned",
// 				)
// 				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 				return
// 			}

// 			if err != nil {
// 				slog.Error(
// 					"failed to get Workspace Customer by email" +
// 						"something went wrong",
// 				)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 				return
// 			}
// 		} else if customer.Phone.Valid {
// 			slog.Info("get customer by phone")
// 			customer, err = h.cs.GetCustomerByPhone(ctx, workspace.WorkspaceId, customer.Phone.String)

// 			if errors.Is(err, services.ErrCustomerNotFound) {
// 				slog.Warn(
// 					"Customer not found by phone" +
// 						"perhaps the customer is not created or is not returned",
// 				)
// 				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 				return
// 			}

// 			if err != nil {
// 				slog.Error(
// 					"failed to get Workspace Customer by phone" +
// 						"something went wrong",
// 				)
// 				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 				return
// 			}
// 		} else {
// 			fmt.Println("unsupported customer identifier")
// 			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	slog.Info("got Workspace Customer",
// 		slog.String("customerId", customer.CustomerId),
// 		slog.Bool("isCreated", isCreated),
// 	)
// 	slog.Info("issue Customer JWT token")
// 	jwt, err := h.cs.GenerateCustomerToken(customer)
// 	if err != nil {
// 		slog.Error("failed to make jwt token with error", slog.Any("error", err))
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}

// 	resp = CustomerTIRespPayload{
// 		Create:     isCreated,
// 		CustomerId: customer.CustomerId,
// 		Jwt:        jwt,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		slog.Error(
// 			"failed to encode response to json "+
// 				"check the json encoding defn",
// 			slog.String("customerId", customer.CustomerId),
// 		)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// }

func (h *WorkspaceHandler) handleGetWorkspaceMembership(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	member, err := h.ws.GetWorkspaceAccountMember(ctx, account.AccountId, workspaceId)
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

	member, err := h.ws.GetWorkspaceMemberById(ctx, workspaceId, memberId)
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

	customers, err := h.ws.ListWorkspaceCustomers(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace customers", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]CustomerResp, 0, len(customers))
	for _, c := range customers {
		items = append(items, CustomerResp{
			WorkspaceId: workspaceId,
			CustomerId:  c.CustomerId,
			ExternalId:  c.ExternalId,
			Email:       c.Email,
			Phone:       c.Phone,
			Name:        c.Name,
			IsVerified:  c.IsVerified,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
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
	members, err := h.ws.ListWorkspaceMembers(ctx, workspaceId)
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

	workspace, err := h.ws.GetLinkedWorkspaceMember(ctx, account.AccountId, workspaceId)
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceSecretKey(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	sk, err := h.ws.GetWorkspaceSecretKey(ctx, workspaceId)
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
