package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/store"

	"github.com/zyghq/zyg/integrations/email"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

type ThreadHandler struct {
	ws  ports.WorkspaceServicer
	ths ports.ThreadServicer
	ss  ports.SyncServicer
}

func NewThreadHandler(ws ports.WorkspaceServicer, ths ports.ThreadServicer, ss ports.SyncServicer) *ThreadHandler {
	return &ThreadHandler{ws: ws, ths: ths, ss: ss}
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
	hub := sentry.GetHubFromContext(ctx)

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
		hub.CaptureException(err)
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
		hub.CaptureException(err)
		slog.Error("failed to update thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// sync thread without setting labels, keep as is.
	// to clear the labels pass empty slice ref: &[]models.ThreadLabel
	inSync, err := h.ss.SyncThreadRPC(ctx, thread, nil)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("thread synced", slog.Any("versionID", inSync.VersionID))
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

func (h *ThreadHandler) handleReplyThreadMail(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	threadId := r.PathValue("threadId")

	var reqp ReplyThreadMailReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)

	// Get member workspace
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

	// Get workspace thread
	thread, err := h.ths.GetWorkspaceThread(ctx, workspace.WorkspaceId, threadId, nil)
	if errors.Is(err, services.ErrThreadNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get Postmark setting for the workspace
	// Postmark setting must be configured before sending a reply mail
	setting, err := h.ws.GetPostmarkMailServerSetting(ctx, workspace.WorkspaceId)
	if errors.Is(err, services.ErrPostmarkSettingNotFound) {
		hub.CaptureMessage("postmark setting required before sending reply")
		http.Error(w, http.StatusText(http.StatusPreconditionRequired), http.StatusPreconditionRequired)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch postmark mail server setting", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get thread customer to send the reply mail
	customer, err := h.ws.GetCustomer(ctx, workspace.WorkspaceId, thread.Customer.CustomerId, nil)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch customer", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	thread, activity, err := h.ths.SendThreadMailReply(
		ctx, &workspace, &setting, &thread, member, &customer, reqp.TextBody, reqp.HTMLBody)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to send thread mail reply", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var activityCustomer *CustomerActorResp
	var activityMember *MemberActorResp
	if activity.Customer != nil {
		activityCustomer = &CustomerActorResp{
			CustomerId: activity.Customer.CustomerId,
			Name:       activity.Customer.Name,
		}
	} else if activity.Member != nil {
		activityMember = &MemberActorResp{
			MemberId: activity.Member.MemberId,
			Name:     activity.Member.Name,
		}
	}

	resp := ActivityResp{
		ActivityID:   activity.ActivityID,
		ThreadID:     activity.ThreadID,
		ActivityType: activity.ActivityType,
		Customer:     activityCustomer,
		Member:       activityMember,
		Body:         activity.Body,
		CreatedAt:    activity.CreatedAt,
		UpdatedAt:    activity.UpdatedAt,
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

	activities, err := h.ths.ListThreadMessagesWithAttachments(ctx, thread.ThreadId)
	if err != nil {
		slog.Error("failed to fetch thread messages", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items := make([]ActivityWithAttachmentsResp, 0, 100)
	for _, activity := range activities {
		var activityCustomer *CustomerActorResp
		var activityMember *MemberActorResp
		if activity.Customer != nil {
			activityCustomer = &CustomerActorResp{
				CustomerId: activity.Customer.CustomerId,
				Name:       activity.Customer.Name,
			}
		} else if activity.Member != nil {
			activityMember = &MemberActorResp{
				MemberId: activity.Member.MemberId,
				Name:     activity.Member.Name,
			}
		}

		resp := ActivityWithAttachmentsResp{
			ActivityResp: ActivityResp{
				ActivityID:   activity.ActivityID,
				ThreadID:     activity.ThreadID,
				ActivityType: activity.ActivityType,
				Customer:     activityCustomer,
				Member:       activityMember,
				Body:         activity.Body,
				CreatedAt:    activity.CreatedAt,
				UpdatedAt:    activity.UpdatedAt,
			},
			Attachments: activity.Attachments,
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

func (h *ThreadHandler) handleGetMessageAttachment(
	w http.ResponseWriter, r *http.Request, _ *models.Member) {
	ctx := r.Context()

	messageId := r.PathValue("messageId")
	attachmentId := r.PathValue("attachmentId")
	attachment, err := h.ths.GetMessageAttachment(ctx, messageId, attachmentId)
	if errors.Is(err, services.ErrMessageAttachmentNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch message attachment", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Generate presigned URL
	accountId := zyg.CFAccountId()
	accessKeyId := zyg.R2AccessKeyId()
	accessKeySecret := zyg.R2AccessSecretKey()
	s3Bucket := zyg.S3Bucket()
	s3Client, err := store.NewS3(ctx, s3Bucket, accountId, accessKeyId, accessKeySecret)
	if err != nil {
		slog.Error("failed to create s3 client", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	expiresIn := time.Now().Add(time.Hour * 24) // Set expiry to 24hrs from now.
	signedUrl, err := store.PresignedUrl(ctx, s3Client, attachment.ContentKey, expiresIn)
	if err != nil {
		slog.Error(
			"failed to generate attachment signed url",
			slog.Any("attachmentId", attachment.AttachmentId),
			slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	attachment.ContentUrl = signedUrl
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(attachment); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
	}
}

func (h *ThreadHandler) handleSetThreadLabel(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()
	threadId := r.PathValue("threadId")
	hub := sentry.GetHubFromContext(ctx)

	var reqp ThChatLabelReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, member.WorkspaceId, threadId)
	if err != nil {
		hub.CaptureException(err)
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
		hub.CaptureException(err)
		slog.Error("failed to create label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threadLabel, isAdded, err := h.ths.SetLabel(ctx, threadId, label.LabelId, models.LabelAddedBy{}.User())
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to add label to thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	threadLabels := make([]models.ThreadLabel, 0, 1)
	threadLabels = append(threadLabels, threadLabel)
	inSync, err := h.ss.SyncThreadLabelsRPC(ctx, threadId, threadLabels)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync thread labels", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("thread labels synced", slog.Any("versionID", inSync.VersionID))
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
	hub := sentry.GetHubFromContext(ctx)

	threadId := r.PathValue("threadId")
	labelId := r.PathValue("labelId")
	thExist, err := h.ths.ThreadExistsInWorkspace(ctx, member.WorkspaceId, threadId)
	if err != nil {
		hub.CaptureException(err)
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
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace label", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.ths.RemoveThreadLabel(ctx, threadId, label.LabelId)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to delete label from thread", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	labelIds := make([]string, 0, 1)
	labelIds = append(labelIds, label.LabelId)
	inSync, err := h.ss.SyncDeleteThreadLabelsRPC(ctx, threadId, labelIds)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync thread labels", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Info("thread labels synced", slog.Any("versionID", inSync.VersionID))
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

	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)

	var reqp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to decode json payload", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workspaceId := r.PathValue("workspaceId")
	workspace, err := h.ws.GetWorkspace(ctx, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	inboundReq, err := email.FromPostmarkInboundRequest(reqp)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("error parsing postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Log inbound request payload for history and auditability.
	// logging makes sure that we capture Postmark event and don't lose, helps in re-running too.
	err = h.ths.LogPostmarkInboundRequest(ctx, workspaceId, inboundReq.MessageID, inboundReq.Payload)
	if err != nil {
		slog.Error("failed to log postmark inbound request", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Check if the Postmark inbound message request has already been processed.
	isProcessed, err := h.ths.IsPostmarkInboundProcessed(ctx, inboundReq.MessageID)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to check inbound message request processed", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if isProcessed {
		slog.Info("postmark inbound message is already processed")
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return
	}

	// Convert inbound request to Postmark Inbound Message for further processing.
	inboundMessage := inboundReq.ToPostmarkInboundMessage()
	// This also marks the email as verified for the customer as inbound is received directly from mail provider.
	customer, _, err := h.ws.CreateCustomerWithEmail(
		ctx, workspace.WorkspaceId, inboundMessage.FromEmail, true, inboundMessage.FromName)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create customer for postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the system member for the workspace which will process the inbound mail.
	member, err := h.ws.GetSystemMember(ctx, workspace.WorkspaceId)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to get system member for postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Process the Postmark inbound message.
	thread, activity, err := h.ths.ProcessPostmarkInbound(
		ctx, workspace.WorkspaceId, &customer, &member, &inboundMessage)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to process postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Info("processed postmark inbound message",
		slog.Any("threadId", thread.ThreadId), slog.Any("activityID", activity.ActivityID))

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		return
	}
}
