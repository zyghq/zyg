package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
	"io"
	"log/slog"
	"net/http"
)

type SyncService struct{}

func NewSyncService() *SyncService {
	return &SyncService{}
}

// RequestBody defines a custom type wrapping map[string]interface{}
type RequestBody map[string]interface{}

func (r RequestBody) SetField(key string, value interface{}) {
	r[key] = value
}

func (r RequestBody) GetField(key string) interface{} {
	return r[key]
}

func (r RequestBody) RemoveField(key string) {
	delete(r, key)
}

func (r RequestBody) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

func (sy *SyncService) SyncWorkspaceRPC(
	ctx context.Context, workspace models.Workspace) (models.WorkspaceInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync := models.WorkspaceInSync{}
	shape := models.WorkspaceShape{
		WorkspaceID: workspace.WorkspaceId,
		Name:        workspace.Name,
		PublicName:  workspace.Name,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
	restateBaseUrl := zyg.RestateRPCURL()
	client := &http.Client{}

	jsonData, err := json.Marshal(shape)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal workspace shape", slog.Any("err", err))
		return models.WorkspaceInSync{}, err
	}
	url := fmt.Sprintf("%s/sync/%s/upsertWorkspace", restateBaseUrl, workspace.WorkspaceId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create sync request", slog.Any("err", err))
		return models.WorkspaceInSync{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync workspace", slog.Any("err", err))
		return models.WorkspaceInSync{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
		return models.WorkspaceInSync{}, err
	}
	err = json.Unmarshal(body, &inSync)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to unmarshal sync response", slog.Any("err", err))
		return models.WorkspaceInSync{}, err
	}
	return inSync, nil
}

func (sy *SyncService) SyncCustomerRPC(
	ctx context.Context, customer models.Customer) (models.CustomerInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync := models.CustomerInSync{}

	restateBaseUrl := zyg.RestateRPCURL()
	client := &http.Client{}

	requestBody := RequestBody{
		"customerId":      customer.CustomerId,
		"workspaceId":     customer.WorkspaceId,
		"name":            customer.Name,
		"publicName":      customer.Name,
		"role":            customer.Role,
		"avatarUrl":       customer.AvatarUrl(),
		"isEmailVerified": customer.IsEmailVerified,
		"createdAt":       customer.CreatedAt,
		"updatedAt":       customer.UpdatedAt,
	}

	if customer.ExternalId.Valid {
		requestBody.SetField("externalId", customer.ExternalId.String)
	} else {
		requestBody.SetField("externalId", nil)
	}

	if customer.Email.Valid {
		requestBody.SetField("email", customer.Email.String)
	} else {
		requestBody.SetField("email", nil)
	}

	if customer.Phone.Valid {
		requestBody.SetField("phone", customer.Phone.String)
	} else {
		requestBody.SetField("phone", nil)
	}

	jsonData, err := requestBody.ToJSON()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal json", slog.Any("err", err))
		return models.CustomerInSync{}, err
	}
	url := fmt.Sprintf("%s/sync/%s/upsertCustomer", restateBaseUrl, customer.CustomerId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create sync request", slog.Any("err", err))
		return models.CustomerInSync{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync customer", slog.Any("err", err))
		return models.CustomerInSync{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
		return models.CustomerInSync{}, err
	}
	err = json.Unmarshal(responseBody, &inSync)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to unmarshal sync response", slog.Any("err", err))
		return models.CustomerInSync{}, err
	}
	return inSync, nil
}

func (sy *SyncService) SyncMemberRPC(
	ctx context.Context, member models.Member) (models.MemberInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync := models.MemberInSync{}

	restateBaseUrl := zyg.RestateRPCURL()
	client := &http.Client{}

	// Todo: add permission, for now dummy.
	permissions := make(map[string]interface{})
	requestBody := RequestBody{
		"memberId":    member.MemberId,
		"workspaceId": member.WorkspaceId,
		"name":        member.Name,
		"publicName":  member.Name,
		"role":        member.Role,
		"permissions": permissions,
		"avatarUrl":   member.AvatarUrl(),
		"createdAt":   member.CreatedAt,
		"updatedAt":   member.UpdatedAt,
	}

	jsonData, err := requestBody.ToJSON()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal json", slog.Any("err", err))
		return models.MemberInSync{}, err
	}
	url := fmt.Sprintf("%s/sync/%s/upsertMember", restateBaseUrl, member.MemberId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create sync request", slog.Any("err", err))
		return models.MemberInSync{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync member", slog.Any("err", err))
		return models.MemberInSync{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
		return models.MemberInSync{}, err
	}
	err = json.Unmarshal(responseBody, &inSync)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to unmarshal sync response", slog.Any("err", err))
		return models.MemberInSync{}, err
	}
	return inSync, nil
}

func (sy *SyncService) SyncThreadRPC(
	ctx context.Context, thread models.Thread, labels *[]models.ThreadLabel) (models.ThreadInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync := models.ThreadInSync{}

	restateBaseUrl := zyg.RestateRPCURL()
	client := &http.Client{}

	var previewText string
	if thread.InboundMessage != nil {
		previewText = thread.InboundMessage.PreviewText
	} else if thread.OutboundMessage != nil {
		previewText = thread.OutboundMessage.PreviewText
	}

	requestBody := RequestBody{
		"threadId":          thread.ThreadId,
		"workspaceId":       thread.WorkspaceId,
		"customerId":        thread.Customer.CustomerId,
		"title":             thread.Title,
		"description":       thread.Description,
		"previewText":       previewText,
		"status":            thread.ThreadStatus.Status,
		"statusChangedAt":   thread.ThreadStatus.StatusChangedAt,
		"statusChangedById": thread.ThreadStatus.StatusChangedBy.MemberId,
		"stage":             thread.ThreadStatus.Stage,
		"replied":           thread.Replied,
		"priority":          thread.Priority,
		"channel":           thread.Channel,
		"createdById":       thread.CreatedBy.MemberId,
		"updatedById":       thread.UpdatedBy.MemberId,
		"createdAt":         thread.CreatedAt,
		"updatedAt":         thread.UpdatedAt,
	}

	// set or remove assigned member
	if thread.AssignedMember != nil {
		requestBody.SetField("assigneeId", thread.AssignedMember.MemberId)
		requestBody.SetField("assignedAt", thread.AssignedMember.AssignedAt)
	} else {
		requestBody.SetField("assigneeId", nil)
		requestBody.SetField("assignedAt", nil)
	}

	// set or remove inbound sequence ID
	if thread.InboundMessage != nil {
		requestBody.SetField("inboundSeqId", thread.InboundMessage.LastSeqId)
	} else {
		requestBody.SetField("inboundSeqId", nil)
	}

	// set or remove outbound sequence ID
	if thread.OutboundMessage != nil {
		requestBody.SetField("outboundSeqId", thread.OutboundMessage.LastSeqId)
	} else {
		requestBody.SetField("outboundSeqId", nil)
	}

	// only set labels if the labels is not nil, otherwise ignore setting the value in request body.
	if labels != nil {
		labelsMap := make(map[string]interface{})
		for _, label := range *labels {
			labelsMap[label.LabelId] = map[string]interface{}{
				"labelId":   label.LabelId,
				"name":      label.Name,
				"createdAt": label.CreatedAt,
				"updatedAt": label.UpdatedAt,
			}
		}
		requestBody.SetField("labels", labelsMap)
	}

	jsonData, err := requestBody.ToJSON()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal json", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	url := fmt.Sprintf("%s/sync/%s/upsertThread", restateBaseUrl, thread.ThreadId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create sync request", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync thread", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
		return models.ThreadInSync{}, err
	}
	err = json.Unmarshal(responseBody, &inSync)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to unmarshal sync response", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	return inSync, nil
}

func (sy *SyncService) SyncThreadLabelsRPC(
	ctx context.Context, threadId string, labels []models.ThreadLabel) (models.ThreadInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync := models.ThreadInSync{}

	restateBaseUrl := zyg.RestateRPCURL()
	client := &http.Client{}

	requestBody := RequestBody{
		"threadId": threadId,
	}

	labelsMap := make(map[string]interface{})
	for _, label := range labels {
		labelsMap[label.LabelId] = map[string]interface{}{
			"labelId":   label.LabelId,
			"name":      label.Name,
			"createdAt": label.CreatedAt,
			"updatedAt": label.UpdatedAt,
		}
	}
	requestBody.SetField("labels", labelsMap)

	jsonData, err := requestBody.ToJSON()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal json", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	url := fmt.Sprintf("%s/sync/%s/addThreadLabels", restateBaseUrl, threadId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create sync request", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync thread", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
		return models.ThreadInSync{}, err
	}
	err = json.Unmarshal(responseBody, &inSync)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to unmarshal sync response", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	return inSync, nil
}

func (sy *SyncService) SyncDeleteThreadLabelsRPC(
	ctx context.Context, threadId string, labelIds []string) (models.ThreadInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync := models.ThreadInSync{}

	restateBaseUrl := zyg.RestateRPCURL()
	client := &http.Client{}

	requestBody := RequestBody{
		"threadId": threadId,
		"labelIds": labelIds,
	}

	jsonData, err := requestBody.ToJSON()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal json", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	url := fmt.Sprintf("%s/sync/%s/deleteThreadLabels", restateBaseUrl, threadId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create sync request", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to sync thread", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
		return models.ThreadInSync{}, err
	}
	err = json.Unmarshal(responseBody, &inSync)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to unmarshal sync response", slog.Any("err", err))
		return models.ThreadInSync{}, err
	}
	return inSync, nil
}
