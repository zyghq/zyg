package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/zyghq/zyg/integrations/email"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

type ThreadHandler struct {
	ws  ports.WorkspaceServicer
	ths ports.ThreadServicer
}

func NewThreadHandler(ws ports.WorkspaceServicer, ths ports.ThreadServicer) *ThreadHandler {
	return &ThreadHandler{ws: ws, ths: ths}
}

// handleGetThreads returns a list of threads associated with the given member's workspace.
func (h *ThreadHandler) handleGetThreads(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	threads, err := h.ths.ListWorkspaceThreads(ctx, member.WorkspaceId)
	if err != nil {
		slog.Error("failed to fetch workspace threads", slog.Any("err", err))
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
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleUpdateThread(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()
	threadId := r.PathValue("threadId")

	var reqp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	thread, err := h.ths.GetWorkspaceThread(ctx, member.WorkspaceId, threadId, nil)
	if errors.Is(err, services.ErrThreadNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fields := make([]string, 0, len(reqp))

	// Modify priority if present, otherwise set default priority.
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

	// Modify stage which indirectly modifies status, otherwise set default stage and status.
	if stage, found := reqp["stage"]; found {
		if stage == nil {
			// set the default stage
			thread.SetDefaultStatus(member.AsMemberActor())
			fields = append(fields, "stage")
		} else {
			stage := stage.(string)
			isValid := thread.ThreadStatus.IsValidStage(stage)
			if !isValid {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			thread.SetStatusStage(stage, member.AsMemberActor())
			fields = append(fields, "stage")
		}
	}

	// Modify assignee if present, otherwise set default assignee.
	// Assignee is a workspace member as member actor.
	if assignee, found := reqp["assignee"]; found {
		if assignee == nil {
			thread.ClearAssignedMember()
			fields = append(fields, "assignee")
		} else {
			assigneeId := assignee.(string)
			member, err := h.ws.GetMember(ctx, member.WorkspaceId, assigneeId)
			if errors.Is(err, services.ErrMemberNotFound) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if err != nil {
				slog.Error("failed to fetch assignee", slog.Any("err", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			thread.AssignMember(member.AsMemberActor(), time.Now().UTC())
			fields = append(fields, "assignee")
		}
	}

	thread, err = h.ths.UpdateThread(ctx, thread, fields)
	if err != nil {
		slog.Error("failed to update thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := ThreadResp{}.NewResponse(&thread)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleGetMyThreads(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	threads, err := h.ths.ListMemberThreads(ctx, member.MemberId)
	if err != nil {
		slog.Error("failed to fetch assigned threads", slog.Any("err", err))
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
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleGetUnassignedThreads(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	threads, err := h.ths.ListUnassignedThreads(ctx, member.WorkspaceId)
	if err != nil {
		slog.Error("failed to fetch unassigned threads", slog.Any("err", err))
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
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleGetLabelledThreads(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	labelId := r.PathValue("labelId")
	label, err := h.ws.GetLabel(ctx, member.WorkspaceId, labelId)
	if errors.Is(err, services.ErrLabelNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return

	}
	if err != nil {
		slog.Error("failed to fetch labelled threads", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threads, err := h.ths.ListLabelledThreads(ctx, label.LabelId)
	if err != nil {
		slog.Error("failed to fetch labelled threads", slog.Any("err", err))
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
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleCreateThreadChatMessage(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	threadId := r.PathValue("threadId")

	var reqp ThChatReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	thread, err := h.ths.GetWorkspaceThread(ctx, member.WorkspaceId, threadId, nil)
	if errors.Is(err, services.ErrThreadNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	message, err := h.ths.AppendOutboundThreadChat(ctx, thread, *member, reqp.Message)
	if err != nil {
		slog.Error("failed to append thread chat message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var messageCustomer *CustomerActorResp
	var messageMember *MemberActorResp
	if message.Customer != nil {
		messageCustomer = &CustomerActorResp{
			CustomerId: message.Customer.CustomerId,
			Name:       message.Customer.Name,
		}
	} else if message.Member != nil {
		messageMember = &MemberActorResp{
			MemberId: message.Member.MemberId,
			Name:     message.Member.Name,
		}
	}

	resp := MessageResp{
		ThreadId:  message.ThreadId,
		MessageId: message.MessageId,
		TextBody:  message.TextBody,
		Body:      message.Body,
		Customer:  messageCustomer,
		Member:    messageMember,
		Channel:   message.Channel,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleGetThreadMessages(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	threadId := r.PathValue("threadId")
	thread, err := h.ths.GetWorkspaceThread(ctx, member.WorkspaceId, threadId, nil)
	if errors.Is(err, services.ErrThreadNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	messages, err := h.ths.ListThreadChatMessages(ctx, thread.ThreadId)
	if err != nil {
		slog.Error("failed to fetch thread messages", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]MessageResp, 0, 100)
	for _, message := range messages {
		var messageCustomer *CustomerActorResp
		var messageMember *MemberActorResp
		if message.Customer != nil {
			messageCustomer = &CustomerActorResp{
				CustomerId: message.Customer.CustomerId,
				Name:       message.Customer.Name,
			}
		} else if message.Member != nil {
			messageMember = &MemberActorResp{
				MemberId: message.Member.MemberId,
				Name:     message.Member.Name,
			}
		}
		resp := MessageResp{
			ThreadId:  message.ThreadId,
			MessageId: message.MessageId,
			TextBody:  message.TextBody,
			Body:      message.Body,
			Customer:  messageCustomer,
			Member:    messageMember,
			Channel:   message.Channel,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		}
		items = append(items, resp)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("failed to encode thread chat messages to json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ThreadHandler) handleSetThreadLabel(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	threadId := r.PathValue("threadId")

	var reqp ThChatLabelReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, member.WorkspaceId, threadId)
	if err != nil {
		slog.Error("failed checking thread existence in workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !thExist {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	label, isCreated, err := h.ws.CreateLabel(ctx, member.WorkspaceId, reqp.Name, reqp.Icon)
	if err != nil {
		slog.Error("failed to create label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threadLabel, isAdded, err := h.ths.SetLabel(ctx, threadId, label.LabelId, models.LabelAddedBy{}.User())
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

func (h *ThreadHandler) handleGetThreadLabels(
	w http.ResponseWriter, r *http.Request, member *models.Member) {

	ctx := r.Context()

	threadId := r.PathValue("threadId")
	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, member.WorkspaceId, threadId)
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
		slog.Error("failed to fetch labels for thread", slog.Any("err", err))
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

func (h *ThreadHandler) handleDeleteThreadLabel(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()

	threadId := r.PathValue("threadId")
	labelId := r.PathValue("labelId")
	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, member.WorkspaceId, threadId)
	if err != nil {
		slog.Error("failed checking thread existence in workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !thExist {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

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
	err = h.ths.RemoveThreadLabel(ctx, threadId, label.LabelId)
	if err != nil {
		slog.Error("failed to delete label from thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ThreadHandler) handleGetThreadMetrics(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()
	metrics, err := h.ths.GenerateMemberThreadMetrics(ctx, member.WorkspaceId, member.MemberId)
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
		Active:             metrics.ActiveCount,
		NeedsFirstResponse: metrics.NeedsFirstResponseCount,
		WaitingOnCustomer:  metrics.WaitingOnCustomerCount,
		HoldCount:          metrics.HoldCount,
		NeedsNextResponse:  metrics.NeedsNextResponseCount,
		AssignedToMe:       metrics.MeCount,
		Unassigned:         metrics.UnAssignedCount,
		OtherAssigned:      metrics.OtherAssignedCount,
		Labels:             labels,
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

func (h *WorkspaceHandler) handleCreateWidget(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp CreateWidgetReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	configuration := map[string]interface{}{}
	if reqp.Configuration != nil {
		cf := *reqp.Configuration
		for k, v := range cf {
			configuration[k] = v
		}
	}

	widget, err := h.ws.CreateWidget(ctx, member.WorkspaceId, reqp.Name, configuration)
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

func (h *ThreadHandler) handlePostmarkInboundWebhook(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		slog.Error("error decoding json payload", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	workspace, err := h.ws.GetWorkspace(ctx, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("got postmark inbound message for workspace",
		slog.Any("workspaceId", workspace.WorkspaceId),
	)

	inboundReq, err := email.FromPostmarkInboundRequest(reqp)
	if err != nil {
		slog.Error("error parsing postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Log inbound request for history and auditability.
	err = h.ths.LogPostmarkInboundRequest(ctx, workspaceId, inboundReq.MessageID, inboundReq.Payload)
	if err != nil {
		slog.Error("failed to log postmark inbound request", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Check if the Postmark inbound message request has already been processed.
	isProcessed, err := h.ths.IsPostmarkInboundProcessed(ctx, inboundReq.MessageID)
	if err != nil {
		slog.Error("failed to check inbound message request processed", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if isProcessed {
		slog.Info("postmark inbound message is already processed")
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return
	}

	inboundMessage := inboundReq.ToPostmarkInboundMessage()

	customer, _, err := h.ws.CreateCustomerWithEmail(
		ctx, workspace.WorkspaceId, inboundMessage.FromEmail, true, inboundMessage.FromName)
	if err != nil {
		slog.Error("failed to create customer for postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the system member for the workspace which will process the inbound mail.
	member, err := h.ws.GetSystemMember(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error("failed to get system member for postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Process the Postmark inbound message.
	thread, message, err := h.ths.ProcessPostmarkInbound(
		ctx, workspace.WorkspaceId, customer.AsCustomerActor(),
		member.AsMemberActor(), &inboundMessage,
	)
	if err != nil {
		slog.Error("failed to process postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Info("processed postmark inbound message with threadId %s and messageId %s",
		thread.ThreadId, message.MessageId)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		return
	}
}
