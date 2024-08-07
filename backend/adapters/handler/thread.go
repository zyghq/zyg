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
	ths ports.ThreadServicer
}

func NewThreadChatHandler(
	ws ports.WorkspaceServicer,
	ths ports.ThreadServicer,
) *ThreadChatHandler {
	return &ThreadChatHandler{ws: ws, ths: ths}
}

func (h *ThreadChatHandler) handleGetThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")

	ctx := r.Context()

	threads, err := h.ths.ListWorkspaceThreadChats(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace threads", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadResp, 0, 100)

	for _, thread := range threads {
		var threadAssignee *ThMemberResp
		var ingressCustomer *ThCustomerResp
		var egressMember *ThMemberResp

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

		if thread.IngressMessageId.Valid {
			ingressCustomer = &ThCustomerResp{
				CustomerId: thread.IngressCustomerId.String,
				Name:       thread.IngressCustomerName.String,
			}
		}

		if thread.EgressMessageId.Valid {
			egressMember = &ThMemberResp{
				MemberId: thread.EgressMemberId.String,
				Name:     thread.EgressMemberName.String,
			}
		}

		resp := ThreadResp{
			ThreadId:        thread.ThreadId,
			Customer:        threadCustomer,
			Title:           thread.Title,
			Description:     thread.Description,
			Sequence:        thread.Sequence,
			Status:          thread.Status,
			Read:            thread.Read,
			Replied:         thread.Replied,
			Priority:        thread.Priority,
			Spam:            thread.Spam,
			Channel:         thread.Channel,
			PreviewText:     thread.PreviewText,
			Assignee:        threadAssignee,
			IngressFirstSeq: thread.IngressFirstSeq,
			IngressLastSeq:  thread.IngressLastSeq,
			IngressCustomer: ingressCustomer,
			EgressFirstSeq:  thread.EgressFirstSeq,
			EgressLastSeq:   thread.EgressLastSeq,
			EgressMember:    egressMember,
			CreatedAt:       thread.CreatedAt,
			UpdatedAt:       thread.UpdatedAt,
		}
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
	channel := models.ThreadChannel{}.Chat()
	thread, err := h.ths.GetWorkspaceThread(ctx, workspaceId, threadId, &channel)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fields := make([]string, 0, len(reqp))

	if priority, found := reqp["priority"]; found {
		if priority == nil {
			// set default priority
			thread.Priority = models.ThreadPriority{}.DefaultPriority()
			fields = append(fields, "priority")
		} else {
			ps := priority.(string)
			isValid := models.ThreadPriority{}.IsValid(ps)
			if !isValid {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			thread.Priority = ps
			fields = append(fields, "priority")
		}
	}

	if status, found := reqp["status"]; found {
		if status == nil {
			// set default status
			thread.Status = models.ThreadStatus{}.DefaultStatus()
			fields = append(fields, "status")
		} else {
			status := status.(string)
			isValid := models.ThreadStatus{}.IsValid(status)
			if !isValid {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			thread.Status = status
			fields = append(fields, "status")
		}
	}

	if assignee, found := reqp["assignee"]; found {
		if assignee == nil {
			thread.AssigneeId = models.NullString(nil)
			fields = append(fields, "assignee")
		} else {
			assigneeId := assignee.(string)
			member, err := h.ws.GetMember(ctx, workspaceId, assigneeId)
			if errors.Is(err, services.ErrMemberNotFound) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if err != nil {
				slog.Error("failed to fetch assignee", slog.Any("err", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			thread.AssigneeId = models.NullString(&member.MemberId)
			fields = append(fields, "assignee")
		}
	}

	if read, found := reqp["read"]; found {
		if read == nil {
			// set default read
			thread.Read = false
			fields = append(fields, "read")
		} else {
			read := read.(bool)
			thread.Read = read
			fields = append(fields, "read")
		}
	}

	if replied, found := reqp["replied"]; found {
		if replied == nil {
			// set default replied
			thread.Replied = false
			fields = append(fields, "replied")
		} else {
			replied := replied.(bool)
			thread.Replied = replied
			fields = append(fields, "replied")
		}
	}

	if spam, found := reqp["spam"]; found {
		if spam == nil {
			// set default spam
			thread.Spam = false
			fields = append(fields, "spam")
		} else {
			spam := spam.(bool)
			thread.Spam = spam
			fields = append(fields, "spam")
		}
	}

	thread, err = h.ths.UpdateThread(ctx, thread, fields)
	if err != nil {
		slog.Error("failed to update thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var threadAssignee *ThMemberResp
	var ingressCustomer *ThCustomerResp
	var egressMember *ThMemberResp

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

	if thread.IngressMessageId.Valid {
		ingressCustomer = &ThCustomerResp{
			CustomerId: thread.IngressCustomerId.String,
			Name:       thread.IngressCustomerName.String,
		}
	}

	if thread.EgressMessageId.Valid {
		egressMember = &ThMemberResp{
			MemberId: thread.EgressMemberId.String,
			Name:     thread.EgressMemberName.String,
		}
	}

	resp := ThreadResp{
		ThreadId:        thread.ThreadId,
		Customer:        threadCustomer,
		Title:           thread.Title,
		Description:     thread.Description,
		Sequence:        thread.Sequence,
		Status:          thread.Status,
		Read:            thread.Read,
		Replied:         thread.Replied,
		Priority:        thread.Priority,
		Spam:            thread.Spam,
		Channel:         thread.Channel,
		PreviewText:     thread.PreviewText,
		Assignee:        threadAssignee,
		IngressFirstSeq: thread.IngressFirstSeq,
		IngressLastSeq:  thread.IngressLastSeq,
		IngressCustomer: ingressCustomer,
		EgressFirstSeq:  thread.EgressFirstSeq,
		EgressLastSeq:   thread.EgressLastSeq,
		EgressMember:    egressMember,
		CreatedAt:       thread.CreatedAt,
		UpdatedAt:       thread.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetMyThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")

	ctx := r.Context()

	member, err := h.ws.GetAccountLinkedMember(ctx, workspaceId, account.AccountId)
	if errors.Is(err, services.ErrMemberNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch member", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	threads, err := h.ths.ListMemberThreadChats(ctx, member.MemberId)
	if err != nil {
		slog.Error("failed to fetch assigned threads", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadResp, 0, 100)

	for _, thread := range threads {
		var threadCustomer ThCustomerResp
		var threadAssignee *ThMemberResp
		var ingressCustomer *ThCustomerResp
		var egressMember *ThMemberResp

		threadCustomer = ThCustomerResp{
			CustomerId: thread.CustomerId,
			Name:       thread.CustomerName,
		}

		if thread.AssigneeId.Valid {
			threadAssignee = &ThMemberResp{
				MemberId: thread.AssigneeId.String,
				Name:     thread.AssigneeName.String,
			}
		}

		if thread.IngressMessageId.Valid {
			ingressCustomer = &ThCustomerResp{
				CustomerId: thread.IngressCustomerId.String,
				Name:       thread.IngressCustomerName.String,
			}
		}

		if thread.EgressMessageId.Valid {
			egressMember = &ThMemberResp{
				MemberId: thread.EgressMemberId.String,
				Name:     thread.EgressMemberName.String,
			}
		}

		resp := ThreadResp{
			ThreadId:        thread.ThreadId,
			Customer:        threadCustomer,
			Title:           thread.Title,
			Description:     thread.Description,
			Sequence:        thread.Sequence,
			Status:          thread.Status,
			Read:            thread.Read,
			Replied:         thread.Replied,
			Priority:        thread.Priority,
			Spam:            thread.Spam,
			Channel:         thread.Channel,
			PreviewText:     thread.PreviewText,
			Assignee:        threadAssignee,
			IngressFirstSeq: thread.IngressFirstSeq,
			IngressLastSeq:  thread.IngressLastSeq,
			IngressCustomer: ingressCustomer,
			EgressFirstSeq:  thread.EgressFirstSeq,
			EgressLastSeq:   thread.EgressLastSeq,
			EgressMember:    egressMember,
			CreatedAt:       thread.CreatedAt,
			UpdatedAt:       thread.UpdatedAt,
		}
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

func (h *ThreadChatHandler) handleGetUnassignedThChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")

	ctx := r.Context()

	threads, err := h.ths.ListUnassignedThreadChats(ctx, workspaceId)
	if err != nil {
		slog.Error("failed to fetch unassigned threads", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadResp, 0, 100)

	for _, thread := range threads {
		var threadCustomer ThCustomerResp
		var threadAssignee *ThMemberResp
		var ingressCustomer *ThCustomerResp
		var egressMember *ThMemberResp

		threadCustomer = ThCustomerResp{
			CustomerId: thread.CustomerId,
			Name:       thread.CustomerName,
		}

		if thread.AssigneeId.Valid {
			threadAssignee = &ThMemberResp{
				MemberId: thread.AssigneeId.String,
				Name:     thread.AssigneeName.String,
			}
		}

		if thread.IngressMessageId.Valid {
			ingressCustomer = &ThCustomerResp{
				CustomerId: thread.IngressCustomerId.String,
				Name:       thread.IngressCustomerName.String,
			}
		}

		if thread.EgressMessageId.Valid {
			egressMember = &ThMemberResp{
				MemberId: thread.EgressMemberId.String,
				Name:     thread.EgressMemberName.String,
			}
		}

		resp := ThreadResp{
			ThreadId:        thread.ThreadId,
			Customer:        threadCustomer,
			Title:           thread.Title,
			Description:     thread.Description,
			Sequence:        thread.Sequence,
			Status:          thread.Status,
			Read:            thread.Read,
			Replied:         thread.Replied,
			Priority:        thread.Priority,
			Spam:            thread.Spam,
			Channel:         thread.Channel,
			PreviewText:     thread.PreviewText,
			Assignee:        threadAssignee,
			IngressFirstSeq: thread.IngressFirstSeq,
			IngressLastSeq:  thread.IngressLastSeq,
			IngressCustomer: ingressCustomer,
			EgressFirstSeq:  thread.EgressFirstSeq,
			EgressLastSeq:   thread.EgressLastSeq,
			EgressMember:    egressMember,
			CreatedAt:       thread.CreatedAt,
			UpdatedAt:       thread.UpdatedAt,
		}
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

func (h *ThreadChatHandler) handleGetLabelledThreadChats(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	labelId := r.PathValue("labelId")

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

	threads, err := h.ths.ListLabelledThreadChats(ctx, label.LabelId)
	if err != nil {
		slog.Error("failed to fetch labelled threads", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadResp, 0, 100)

	for _, thread := range threads {
		var threadCustomer ThCustomerResp
		var threadAssignee *ThMemberResp
		var ingressCustomer *ThCustomerResp
		var egressMember *ThMemberResp

		threadCustomer = ThCustomerResp{
			CustomerId: thread.CustomerId,
			Name:       thread.CustomerName,
		}

		if thread.AssigneeId.Valid {
			threadAssignee = &ThMemberResp{
				MemberId: thread.AssigneeId.String,
				Name:     thread.AssigneeName.String,
			}
		}

		if thread.IngressMessageId.Valid {
			ingressCustomer = &ThCustomerResp{
				CustomerId: thread.IngressCustomerId.String,
				Name:       thread.IngressCustomerName.String,
			}
		}

		if thread.EgressMessageId.Valid {
			egressMember = &ThMemberResp{
				MemberId: thread.EgressMemberId.String,
				Name:     thread.EgressMemberName.String,
			}
		}

		resp := ThreadResp{
			ThreadId:        thread.ThreadId,
			Customer:        threadCustomer,
			Title:           thread.Title,
			Description:     thread.Description,
			Sequence:        thread.Sequence,
			Status:          thread.Status,
			Read:            thread.Read,
			Replied:         thread.Replied,
			Priority:        thread.Priority,
			Spam:            thread.Spam,
			Channel:         thread.Channel,
			PreviewText:     thread.PreviewText,
			Assignee:        threadAssignee,
			IngressCustomer: ingressCustomer,
			EgressMember:    egressMember,
			CreatedAt:       thread.CreatedAt,
			UpdatedAt:       thread.UpdatedAt,
		}
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

func (h *ThreadChatHandler) handleCreateThChatMessage(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	var reqp ThChatReq

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	member, err := h.ws.GetAccountLinkedMember(ctx, workspaceId, account.AccountId)
	if errors.Is(err, services.ErrMemberNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch member", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	channel := models.ThreadChannel{}.Chat()
	thread, err := h.ths.GetWorkspaceThread(ctx, workspaceId, threadId, &channel)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	chat, err := h.ths.AddOutboundMessage(ctx, thread, member.MemberId, reqp.Message)
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

	// improvements:
	// shall we use go routines for async assignment and replied marking?
	// also lets check for workspace setttings for auto assignment and replied marking
	// for now keep it as is.
	// if !thread.AssigneeId.Valid {
	// 	slog.Info("thread chat not yet assigned", "threadId", thread.ThreadId, "memberId", member.MemberId)
	// 	t := thread // make a temp copy before assigning
	// 	thread, err = h.ths.AssignMember(ctx, thread.ThreadId, member.MemberId)
	// 	// if error when assigning - revert back
	// 	if err != nil {
	// 		slog.Error("(silent) failed to assign member to Thread Chat", slog.Any("error", err))
	// 		thread = t
	// 	}
	// }

	// if !thread.Replied {
	// 	slog.Info("thread chat not yet replied", "threadId", thread.ThreadId, "memberId", member.MemberId)
	// 	t := thread // make a temp copy before marking replied
	// 	thread, err = h.ths.SetReplyStatus(ctx, thread.ThreadId, true)
	// 	if err != nil {
	// 		slog.Error("(silent) failed to mark thread chat as replied", slog.Any("error", err))
	// 		thread = t
	// 	}
	// }

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
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleGetThChatMesssages(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	ctx := r.Context()

	channel := models.ThreadChannel{}.Chat()
	thread, err := h.ths.GetWorkspaceThread(ctx, workspaceId, threadId, &channel)
	if errors.Is(err, services.ErrThreadChatNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch thread", slog.Any("err", err))
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
		slog.Error(
			"failed to encode thread chat messages to json "+
				"check the json encoding defn",
			slog.String("threadId", thread.ThreadId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadChatHandler) handleSetThreadChatLabel(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")
	var reqp ThChatLabelReq

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, workspaceId, threadId)
	if err != nil {
		slog.Error("failed checking thread existence in workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !thExist {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	label, isCreated, err := h.ws.CreateLabel(ctx, workspaceId, reqp.Name, reqp.Icon)
	if err != nil {
		slog.Error("failed to create label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threadLabel, isAdded, err := h.ths.SetLabel(
		ctx, threadId, label.LabelId, models.LabelAddedBy{}.User())
	if err != nil {
		slog.Error("failed to add label to thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := ThreadLabelResp{
		ThreadLabelId: threadLabel.ThreadLabelId,
		ThreadId:      threadLabel.ThreadId,
		LabelId:       threadLabel.LabelId,
		Name:          threadLabel.Name,
		Icon:          threadLabel.Icon,
		AddedBy:       threadLabel.AddedBy,
		CreatedAt:     threadLabel.CreatedAt,
		UpdatedAt:     threadLabel.UpdatedAt,
	}

	if isCreated || isAdded {
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

func (h *ThreadChatHandler) handleGetThChatLabels(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	threadId := r.PathValue("threadId")

	ctx := r.Context()

	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, workspaceId, threadId)
	if err != nil {
		slog.Error("failed checking thread existence in workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !thExist {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	labels, err := h.ths.ListThreadLabels(ctx, threadId)
	if err != nil {
		slog.Error("failed to fetch list of labels for thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ThreadLabelResp, 0, 100)
	for _, label := range labels {
		item := ThreadLabelResp{
			ThreadLabelId: label.ThreadLabelId,
			ThreadId:      label.ThreadId,
			LabelId:       label.LabelId,
			Name:          label.Name,
			Icon:          label.Icon,
			AddedBy:       label.AddedBy,
			CreatedAt:     label.CreatedAt,
			UpdatedAt:     label.UpdatedAt,
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

func (h *ThreadChatHandler) handleGetThreadChatMetrics(w http.ResponseWriter, r *http.Request, account *models.Account) {
	workspaceId := r.PathValue("workspaceId")
	ctx := r.Context()

	member, err := h.ws.GetAccountLinkedMember(ctx, workspaceId, account.AccountId)
	if errors.Is(err, services.ErrMemberNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch member", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	metrics, err := h.ths.GenerateMemberThreadMetrics(ctx, workspaceId, member.MemberId)
	if err != nil {
		slog.Error("failed to generate thread metrics", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var label ThreadLabelCountResp
	labels := make([]ThreadLabelCountResp, 0, 100)

	for _, l := range metrics.ThreadLabelMetrics {
		label = ThreadLabelCountResp{
			LabelId: l.LabelId,
			Name:    l.Name,
			Icon:    l.Icon,
			Count:   l.Count,
		}
		labels = append(labels, label)
	}

	count := ThreadCountResp{
		ActiveCount:   metrics.ActiveCount,
		DoneCount:     metrics.DoneCount,
		TodoCount:     metrics.TodoCount,
		SnoozedCount:  metrics.SnoozedCount,
		AssignedToMe:  metrics.MeCount,
		Unassigned:    metrics.UnAssignedCount,
		OtherAssigned: metrics.OtherAssignedCount,
		Labels:        labels,
	}

	resp := ThreadMetricsResp{
		Count: count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *WorkspaceHandler) handleCreateWidget(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	workspaceId := r.PathValue("workspaceId")

	var reqp CreateWidgetReq

	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspace, err := h.as.GetAccountLinkedWorkspace(ctx, account.AccountId, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	configuration := map[string]interface{}{}
	if reqp.Configuration != nil {
		cf := *reqp.Configuration
		for k, v := range cf {
			configuration[k] = v
		}
	}

	widget, err := h.ws.CreateWidget(ctx, workspace.WorkspaceId, reqp.Name, configuration)
	if err != nil {
		slog.Error("failed to create widget", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := WidgetResp{
		WidgetId:      widget.WidgetId,
		Name:          widget.Name,
		Configuration: widget.Configuration,
		CreatedAt:     widget.CreatedAt,
		UpdatedAt:     widget.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
