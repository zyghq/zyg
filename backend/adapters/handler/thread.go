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

type ThreadChatHandler struct {
	ws  ports.WorkspaceServicer
	ths ports.ThreadChatServicer
}

func NewThreadChatHandler(
	ws ports.WorkspaceServicer,
	ths ports.ThreadChatServicer,
) *ThreadChatHandler {
	return &ThreadChatHandler{ws: ws, ths: ths}
}

func (h *ThreadChatHandler) handleGetThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")

	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

	// workspace not found
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for account",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// error workspace
	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threads, err := h.ths.ListWorkspaceThreads(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get list of thread chats for workspace "+
				"something went wrong",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make([]ThChatRespPayload, 0, 100)
	for _, th := range threads {
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
		response = append(response, ThChatRespPayload{
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			"failed to encode thread chats to json "+
				"check the json encoding defn",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleUpdateThreadChat(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	var reqp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	thread, err := h.ths.GetThread(ctx, workspaceId, threadId)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		slog.Warn(
			"no thread chat found",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get thread chat "+
				"something went wrong",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// make a slice of the fields
	fields := make([]string, 0, len(reqp))

	if priority, found := reqp["priority"]; found {
		if priority == nil {
			slog.Error(
				"invalid priority",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		ps := priority.(string)
		isValid := models.ThreadPriority{}.IsValid(ps)
		if !isValid {
			slog.Error(
				"invalid priority",
				slog.String("threadChatId", threadId),
				slog.String("priority", ps),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		thread.Priority = ps
		fields = append(fields, "priority")
	}

	if status, found := reqp["status"]; found {
		if status != nil {
			status := status.(string)
			isValid := models.ThreadStatus{}.IsValid(status)
			if !isValid {
				slog.Error(
					"invalid status",
					slog.String("threadChatId", threadId),
					slog.String("status", status),
				)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			thread.Status = status
			fields = append(fields, "status")
		}
	}

	if assignee, found := reqp["assignee"]; found {
		if assignee != nil {
			assigneeId := assignee.(string)
			member, err := h.ws.GetWorkspaceMemberById(ctx, workspaceId, assigneeId)
			if errors.Is(err, services.ErrMemberNotFound) {
				slog.Warn(
					"no member found in workspace",
					slog.String("accountId", account.AccountId),
					slog.String("workspaceId", workspaceId),
				)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if err != nil {
				slog.Error(
					"failed to get member "+
						"something went wrong",
					slog.String("accountId", account.AccountId),
					slog.String("workspaceId", workspaceId),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			thread.AssigneeId = models.NullString(&member.MemberId)
		}
		if assignee == nil {
			thread.AssigneeId = models.NullString(nil)
		}
		fields = append(fields, "assignee")
	}

	thread, err = h.ths.UpdateThread(ctx, thread, fields)
	if err != nil {
		slog.Error(
			"failed to update thread chat "+
				"something went wrong",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var threadAssigneeRepr *ThMemberRespPayload

	threadCustomerRepr := ThCustomerRespPayload{
		CustomerId: thread.CustomerId,
		Name:       thread.CustomerName,
	}

	if thread.AssigneeId.Valid {
		threadAssigneeRepr = &ThMemberRespPayload{
			MemberId: thread.AssigneeId.String,
			Name:     thread.AssigneeName,
		}
	}

	resp := ThChatUpdateRespPayload{
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
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread chat to json "+
				"check the json encoding defn",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetMyThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")

	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

	// not found workspace
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for account",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// error workspace
	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	member, err := h.ws.GetWorkspaceMember(ctx, account.AccountId, workspace.WorkspaceId)
	// error workspace member
	if err != nil {
		slog.Error(
			"failed to get member "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	threads, err := h.ths.ListMemberAssignedThreads(ctx, workspace.WorkspaceId, member.MemberId)
	if err != nil {
		slog.Error(
			"failed to get list of thread chats for workspace "+
				"something went wrong",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make([]ThChatRespPayload, 0, 100)
	for _, th := range threads {
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
		response = append(response, ThChatRespPayload{
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			"failed to encode thread chats to json "+
				"check the json encoding defn",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetUnassignedThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")

	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

	// not found workspace
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for account",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// error workspace
	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threads, err := h.ths.ListUnassignedThreads(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get list of thread chats for workspace "+
				"something went wrong",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make([]ThChatRespPayload, 0, 100)
	for _, th := range threads {
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
		response = append(response, ThChatRespPayload{
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			"failed to encode thread chats to json "+
				"check the json encoding defn",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetLabelledThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	labelId := r.PathValue("labelId")

	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for account",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// error workspace
	if err != nil {
		slog.Error(
			"failed to get workspace by id "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	label, err := h.ws.GetWorkspaceLabel(ctx, workspace.WorkspaceId, labelId)
	if errors.Is(err, services.ErrLabelNotFound) {
		slog.Warn(
			"label not found or does not exist",
			"labelId", labelId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return

	}

	if err != nil {
		slog.Error(
			"failed to get label "+
				"something went wrong",
			"labelId", labelId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threads, err := h.ths.ListLabelledThreads(ctx, workspace.WorkspaceId, label.LabelId)
	if err != nil {
		slog.Error(
			"failed to get list of thread chats for workspace "+
				"something went wrong",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make([]ThChatRespPayload, 0, 100)
	for _, th := range threads {
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
		response = append(response, ThChatRespPayload{
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			"failed to encode thread chats to json "+
				"check the json encoding defn",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleCreateThChatMessage(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	var message ThChatReqPayload

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// check member against workspace
	member, err := h.ws.GetWorkspaceMember(ctx, account.AccountId, workspaceId)

	if errors.Is(err, services.ErrMemberNotFound) {
		slog.Warn(
			"no member found in workspace",
			slog.String("accountId", account.AccountId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get member "+
				"something went wrong",
			slog.String("accountId", account.AccountId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check thread chat against workspace
	thread, err := h.ths.GetThread(ctx, workspaceId, threadId)

	if errors.Is(err, services.ErrThreadChatNotFound) {
		slog.Warn(
			"no thread chat found",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get thread chat "+
				"something went wrong",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// create thread chat message
	thm, err := h.ths.AddMemberMessageToThread(ctx, thread, &member, message.Message)

	if err != nil {
		slog.Error(
			"failed to create thread chat message "+
				"something went wrong",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check for assignment and replied mark
	// if not assigned try assigning the member who sent is sending the message
	// mark replied if not already marked.
	//
	// when assigning or marking as replied, error is ignored and silently log it.
	//
	// (sanchitrk):
	// probably have some way to set the auto marking and assignment on member reply
	// as configurable in settings.
	if !thread.AssigneeId.Valid {
		slog.Info("thread chat not yet assigned", "threadChatId", thread.ThreadChatId, "memberId", member.MemberId)
		t := thread // make a temp copy before assigning
		thread, err = h.ths.AssignMemberToThread(ctx, thread.ThreadChatId, member.MemberId)
		// if error when assigning - revert back
		if err != nil {
			slog.Error("(silent) failed to assign member to Thread Chat", slog.Any("error", err))
			thread = t
		}
	}

	if !thread.Replied {
		slog.Info("thread chat not yet replied", "threadChatId", thread.ThreadChatId, "memberId", member.MemberId)
		t := thread // make a temp copy before marking replied
		thread, err = h.ths.SetThreadReplyStatus(ctx, thread.ThreadChatId, true)
		if err != nil {
			slog.Error("(silent) failed to mark thread chat as replied", slog.Any("error", err))
			thread = t
		}
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
			"failed to encode response",
			slog.Any("error", err),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetThChatMesssages(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	ctx := r.Context()

	thread, err := h.ths.GetThread(ctx, workspaceId, threadId)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		slog.Warn(
			"no thread chat found",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error(
			"failed to get thread chat "+
				"something went wrong",
			slog.String("threadChatId", threadId),
			slog.String("workspaceId", workspaceId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	messages, err := h.ths.ListThreadMessages(ctx, thread.ThreadChatId)
	if err != nil {
		slog.Error(
			"failed to get list of thread chat messages for thread chat "+
				"something went wrong",
			slog.String("threadChatId", threadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	results := make([]ThChatMessageRespPayload, 0, 100)
	for _, message := range messages {
		var msgCustomerRepr *ThCustomerRespPayload
		var msgMemberRepr *ThMemberRespPayload

		// for thread message - either of them
		if message.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerRespPayload{
				CustomerId: message.CustomerId.String,
				Name:       message.CustomerName,
			}
		} else if message.MemberId.Valid {
			msgMemberRepr = &ThMemberRespPayload{
				MemberId: message.MemberId.String,
				Name:     message.MemberName,
			}
		}

		threadMessage := ThChatMessageRespPayload{
			ThreadChatId:        thread.ThreadChatId,
			ThreadChatMessageId: message.ThreadChatMessageId,
			Body:                message.Body,
			Sequence:            message.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           message.CreatedAt,
			UpdatedAt:           message.UpdatedAt,
		}
		results = append(results, threadMessage)
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
		Messages:     results,
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

func (h *ThreadChatHandler) handleSetThChatLabel(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")
	var reqp SetThChatLabelReqPayload

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	isThExist, err := h.ths.ThreadExistsInWorkspace(ctx, workspaceId, threadId)

	if err != nil {
		slog.Error(
			"failed to check if thread chat exist in workspace",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isThExist {
		slog.Warn(
			"thread chat not found or does not exist in workspace",
			"threadChatId", threadId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	label := models.Label{
		WorkspaceId: workspaceId,
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

	thChatLabel := models.ThreadChatLabel{
		ThreadChatId: threadId,
		LabelId:      label.LabelId,
		AddedBy:      models.LabelAddedBy{}.User(),
	}

	thChatLabel, isAdded, err := h.ths.AttachLabelToThread(ctx, thChatLabel)
	if err != nil {
		slog.Error(
			"failed to add label to thread chat something went wrong",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cr := CrLabelRespPayload{
		LabelId:   label.LabelId,
		Name:      label.Name,
		Icon:      label.Icon,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}
	resp := SetThChatLabelRespPayload{
		ThreadChatLabelId:  thChatLabel.ThreadChatLabelId,
		ThreadChatId:       thChatLabel.ThreadChatId,
		AddedBy:            thChatLabel.AddedBy,
		CreatedAt:          thChatLabel.CreatedAt,
		UpdatedAt:          thChatLabel.UpdatedAt,
		CrLabelRespPayload: cr,
	}

	// if any of the workspace label or thread chat label is created
	if isCreated || isAdded {
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

func (h *ThreadChatHandler) handleGetThreadChatLabels(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	ctx := r.Context()

	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, workspaceId, threadId)
	if err != nil {
		slog.Error(
			"failed to check if thread chat exists in workspace "+
				"perhaps a failed query or returned nothing",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !thExist {
		slog.Warn(
			"thread chat not found or does not exist in workspace",
			"threadChatId", threadId,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	resp := make([]SetThChatLabelRespPayload, 0, 100)

	labels, err := h.ths.ListThreadLabels(ctx, threadId)
	if err != nil {
		slog.Error(
			"failed to get list of labels for thread chat "+
				"something went wrong",
			"threadChatId", threadId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, label := range labels {
		cr := CrLabelRespPayload{
			LabelId:   label.LabelId,
			Name:      label.Name,
			Icon:      label.Icon,
			CreatedAt: label.CreatedAt,
			UpdatedAt: label.UpdatedAt,
		}
		resp = append(resp, SetThChatLabelRespPayload{
			ThreadChatLabelId:  label.ThreadChatLabelId,
			ThreadChatId:       label.ThreadChatId,
			AddedBy:            label.AddedBy,
			CreatedAt:          label.CreatedAt,
			UpdatedAt:          label.UpdatedAt,
			CrLabelRespPayload: cr,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode labels to json "+
				"check the json encoding defn",
			"threadChatId", threadId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetThreadChatMetrics(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	ctx := r.Context()

	workspace, err := h.ws.GetMemberWorkspace(ctx, account.AccountId, workspaceId)

	if errors.Is(err, services.ErrWorkspaceNotFound) {
		slog.Warn(
			"workspace not found or does not exist for account",
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

	member, err := h.ws.GetWorkspaceMember(ctx, account.AccountId, workspace.WorkspaceId)
	if err != nil {
		slog.Error(
			"failed to get member "+
				"something went wrong",
			"accountId", account.AccountId, "workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	metrics, err := h.ths.GenerateMemberThreadMetrics(ctx, workspace.WorkspaceId, member.MemberId)
	if err != nil {
		slog.Error(
			"failed to generate thread chat metrics "+
				"something went wrong",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var label ThreadLabelCountRespPayload
	labels := make([]ThreadLabelCountRespPayload, 0, 100)

	for _, l := range metrics.ThreadLabelMetrics {
		label = ThreadLabelCountRespPayload{
			LabelId: l.LabelId,
			Name:    l.Name,
			Icon:    l.Icon,
			Count:   l.Count,
		}
		labels = append(labels, label)
	}

	count := ThreadCountRespPayload{
		ActiveCount:   metrics.ActiveCount,
		DoneCount:     metrics.DoneCount,
		TodoCount:     metrics.TodoCount,
		SnoozedCount:  metrics.SnoozedCount,
		AssignedToMe:  metrics.MeCount,
		Unassigned:    metrics.UnAssignedCount,
		OtherAssigned: metrics.OtherAssignedCount,
		Labels:        labels,
	}

	resp := ThreadMetricsRespPayload{
		Count: count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(
			"failed to encode thread metrics to json "+
				"check the json encoding defn",
			"workspaceId", workspace.WorkspaceId,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
