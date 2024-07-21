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

	var rb WorkspaceReqPayload
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspace := models.Workspace{AccountId: account.AccountId, Name: rb.Name}
	workspace, err = h.ws.CreateWorkspace(ctx, *account, workspace)

	if err != nil {
		slog.Error(
			"failed to create workspace "+
				"something went wrong",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(workspace); err != nil {
		slog.Error(
			"failed to encode workspace to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaces(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaces, err := h.ws.ListMemberWorkspaces(ctx, account.AccountId)

	if err != nil {
		slog.Error(
			"failed to get list of workspaces "+
				"something went wrong",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workspaces); err != nil {
		slog.Error(
			"failed to encode workspaces to json "+
				"check the json encoding defn",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspace(w http.ResponseWriter, r *http.Request, account *models.Account) {

	ctx := r.Context()
	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workspace); err != nil {
		slog.Error(
			"failed to encode workspace to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspace.WorkspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleUpdateWorkspace(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var rb WorkspaceReqPayload
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// apply updates
	workspace.Name = rb.Name

	workspace, err = h.ws.UpdateWorkspace(ctx, workspace)
	if err != nil {
		slog.Error(
			"failed to update workspace "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workspace); err != nil {
		slog.Error(
			"failed to encode workspace to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspace.WorkspaceId),
		)
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
	var reqp CrLabelReqPayload

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	label := models.Label{
		WorkspaceId: workspace.WorkspaceId,
		Name:        reqp.Name,
		Icon:        reqp.Icon,
	}

	label, isCreated, err := h.ws.CreateLabel(ctx, label)

	if err != nil {
		slog.Error(
			"failed to get or create label something went wrong",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := CrLabelRespPayload{
		LabelId:   label.LabelId,
		Name:      label.Name,
		Icon:      label.Icon,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}

	if isCreated {
		slog.Info("created label", "labelId", label.LabelId)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode label to json "+
					"check the json encoding defn",
				"labelId", label.LabelId,
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		slog.Info("label already exists", "labelId", label.LabelId)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode label to json "+
					"check the json encoding defn",
				"labelId", label.LabelId,
			)
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
	var reqp CrLabelReqPayload

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	label, err := h.ws.GetWorkspaceLabel(ctx, workspaceId, labelId)
	if errors.Is(err, services.ErrLabelNotFound) {
		slog.Warn(
			"label not found or does not exist",
			slog.String("labelId", labelId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get label by id "+
				"something went wrong",
			slog.String("labelId", labelId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// apply updates
	label.Name = reqp.Name
	label.Icon = reqp.Icon

	label, err = h.ws.SetWorkspaceLabel(ctx, label.WorkspaceId, label)
	if err != nil {
		slog.Error(
			"failed to update label "+
				"something went wrong",
			slog.String("labelId", labelId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := CrLabelRespPayload{
		LabelId:   label.LabelId,
		Name:      label.Name,
		Icon:      label.Icon,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			"failed to encode label to json "+
				"check the json encoding defn",
			slog.String("labelId", label.LabelId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceLabels(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	results, err := h.ws.ListWorkspaceLabels(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get workspace labels "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	labels := make([]CrLabelRespPayload, 0, len(results))
	for _, l := range results {
		labels = append(labels, CrLabelRespPayload{
			LabelId:   l.LabelId,
			Name:      l.Name,
			Icon:      l.Icon,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(labels); err != nil {
		slog.Error(
			"failed to encode workspace labels to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspaceId),
		)
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

// 	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

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
// 			customer, isCreated, err = h.ws.CreateWorkspaceCustomerWithEmail(ctx, customer)

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
// 			customer, isCreated, err = h.ws.CreateWorkspaceCustomerWithPhone(ctx, customer)

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
// 			customer, isCreated, err = h.ws.CreateCustomerByExternalId(ctx, customer)

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

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	member, err := h.ws.GetWorkspaceMember(ctx, account.AccountId, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get workspace membership "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(member); err != nil {
		slog.Error(
			"failed to encode workspace membership to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceMember(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	memberId := r.PathValue("memberId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	member, err := h.ws.GetWorkspaceMemberById(ctx, workspace.WorkspaceId, memberId)
	if errors.Is(err, services.ErrMemberNotFound) {
		slog.Warn(
			"workspace member not found or does not exist",
			slog.String("memberId", memberId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace member or does not exist "+
				"something went wrong",
			slog.String("memberId", memberId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(member); err != nil {
		slog.Error(
			"failed to encode workspace member to json "+
				"check the json encoding defn",
			slog.String("memberId", memberId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceCustomers(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	customers, err := h.ws.ListWorkspaceCustomers(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get workspace customers "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customers); err != nil {
		slog.Error(
			"failed to encode workspace customers to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceMembers(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	members, err := h.ws.ListWorkspaceMembers(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get workspace members "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(members); err != nil {
		slog.Error(
			"failed to encode workspace members to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspaceId),
		)
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

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sk, err := h.ws.GenerateSecretKey(ctx, workspace.WorkspaceId, 64)
	if err != nil {
		slog.Error(
			"failed to generate workspace secret key "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := SKRespPayload{
		SecretKey: sk.SecretKey,
		CreatedAt: sk.CreatedAt,
		UpdatedAt: sk.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode workspace secret key to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWorkspaceSecretKey(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get account workspace or does not exist "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sk, err := h.ws.GetWorkspaceSecretKey(ctx, workspace.WorkspaceId)

	if errors.Is(err, services.ErrSecretKeyNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get workspace secret key "+
				"something went wrong",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := SKRespPayload{
		SecretKey: sk.SecretKey,
		CreatedAt: sk.CreatedAt,
		UpdatedAt: sk.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode workspace secret key to json "+
				"check the json encoding defn",
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleGetWidgets(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to get workspace by id "+
			"something went wrong", "accountId", account.AccountId, "workspaceId", workspaceId)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	widgets, err := h.ws.ListWidgets(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get list of widgets for workspace "+
				"something went wrong",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make([]WidgetRespPayload, 0, 100)
	for _, widget := range widgets {
		resp := WidgetRespPayload{
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
		slog.Error(
			"failed to encode widgets to json "+
				"check the json encoding defn",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
