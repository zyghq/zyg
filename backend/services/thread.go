package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type ThreadChatService struct {
	repo ports.ThreadRepositorer
}

func NewThreadChatService(repo ports.ThreadRepositorer) *ThreadChatService {
	return &ThreadChatService{
		repo: repo,
	}
}

func (s *ThreadChatService) CreateThreadInAppChat(
	ctx context.Context, workspaceId string, customerId string, message string) (models.Thread, models.Chat, error) {
	chat := models.Chat{
		Body:       message,
		CustomerId: models.NullString(&customerId),
		IsHead:     true,
	}
	thread := models.Thread{
		WorkspaceId:       workspaceId,
		CustomerId:        customerId,
		Status:            models.ThreadStatus{}.Todo(),
		Channel:           models.ThreadChannel{}.Chat(),
		MessageBody:       chat.PreviewBody(),
		MessageCustomerId: chat.CustomerId,
		MessageMemberId:   chat.MemberId,
	}
	thread, chat, err := s.repo.InsertInAppThreadChat(ctx, thread, chat)
	if err != nil {
		return models.Thread{}, models.Chat{}, ErrThreadChat
	}
	return thread, chat, nil
}

func (s *ThreadChatService) UpdateThread(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {
	thread, err := s.repo.ModifyThreadById(ctx, thread, fields)

	if err != nil {
		return models.Thread{}, ErrThreadChat
	}

	return thread, nil
}

func (s *ThreadChatService) GetWorkspaceThread(
	ctx context.Context, workspaceId string, threadId string) (models.Thread, error) {
	thread, err := s.repo.LookupByWorkspaceThreadId(ctx, workspaceId, threadId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Thread{}, ErrThreadChatNotFound
	}

	if err != nil {
		return models.Thread{}, ErrThreadChat
	}

	return thread, nil
}

func (s *ThreadChatService) ListCustomerThreadChats(
	ctx context.Context, workspaceId string, customerId string) ([]models.Thread, error) {
	threads, err := s.repo.RetrieveWorkspaceThChatsByCustomerId(ctx, workspaceId, customerId)
	if err != nil {
		return threads, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) AssignMemberToThread(ctx context.Context, threadId string, assigneeId string) (models.Thread, error) {
	thread, err := s.repo.UpdateAssignee(ctx, threadId, assigneeId)
	if err != nil {
		return thread, ErrThreadChatAssign
	}
	return thread, nil
}

func (s *ThreadChatService) SetThreadReplyStatus(ctx context.Context, threadChatId string, replied bool) (models.Thread, error) {
	thread, err := s.repo.UpdateRepliedStatus(ctx, threadChatId, replied)
	if err != nil {
		return thread, ErrThreadChatReplied
	}
	return thread, nil
}

func (s *ThreadChatService) ListWorkspaceThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	threads, err := s.repo.FetchThChatsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListMemberAssignedThreadChats(
	ctx context.Context, memberId string) ([]models.Thread, error) {
	threads, err := s.repo.FetchAssignedThChatsByMemberId(ctx, memberId)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListUnassignedThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	threads, err := s.repo.FetchUnassignedThChatsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListLabelledThreadChats(
	ctx context.Context, labelId string) ([]models.Thread, error) {
	threads, err := s.repo.FetchThChatsByLabelId(ctx, labelId)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ThreadExistsInWorkspace(
	ctx context.Context, workspaceId string, threadId string) (bool, error) {
	exist, err := s.repo.CheckWorkspaceExistenceByThreadId(ctx, workspaceId, threadId)
	if err != nil {
		return false, ErrThreadChat
	}
	return exist, nil
}

func (s *ThreadChatService) AttachLabelToThread(
	ctx context.Context, threadId string, labelId string, addedBy string) (models.ThreadLabel, bool, error) {
	label := models.ThreadLabel{
		ThreadId: threadId,
		LabelId:  labelId,
		AddedBy:  addedBy,
	}
	label, created, err := s.repo.SetLabelToThread(ctx, label)
	if err != nil {
		return models.ThreadLabel{}, created, ErrThChatLabel
	}

	return label, created, nil
}

func (s *ThreadChatService) ListThreadLabels(
	ctx context.Context, threadId string) ([]models.ThreadLabel, error) {
	labels, err := s.repo.RetrieveLabelsByThreadId(ctx, threadId)
	if err != nil {
		return labels, ErrThChatLabel
	}
	return labels, nil
}

func (s *ThreadChatService) AddCustomerMessageToThread(
	ctx context.Context, threadId string, customerId string, message string) (models.Chat, error) {
	chat := models.Chat{
		ThreadId:   threadId,
		Body:       message,
		CustomerId: models.NullString(&customerId),
		IsHead:     false,
	}
	chat, err := s.repo.InsertCustomerMessage(ctx, chat)
	if err != nil {
		return models.Chat{}, ErrThChatMessage
	}

	return chat, nil
}

func (s *ThreadChatService) AddMemberMessageToThreadChat(
	ctx context.Context, threadId string, memberId string, message string) (models.Chat, error) {
	chat := models.Chat{
		ThreadId: threadId,
		Body:     message,
		MemberId: models.NullString(&memberId),
		IsHead:   false,
	}
	chat, err := s.repo.InsertThChatMemberMessage(ctx, chat)
	if err != nil {
		return models.Chat{}, ErrThChatMessage
	}

	return chat, nil
}

func (s *ThreadChatService) ListThreadChatMessages(
	ctx context.Context, threadId string) ([]models.Chat, error) {
	messages, err := s.repo.FetchThChatMessagesByThreadId(ctx, threadId)
	if err != nil {
		return []models.Chat{}, ErrThChatMessage
	}
	return messages, nil
}

func (s *ThreadChatService) GenerateMemberThreadMetrics(
	ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error) {
	statusMetrics, err := s.repo.ComputeStatusMetricsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadChatMetrics
	}

	assignmentMetrics, err := s.repo.ComputeAssigneeMetricsByMember(ctx, workspaceId, memberId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadChatMetrics
	}

	labelMetrics, err := s.repo.ComputeLabelMetricsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadChatMetrics
	}

	metrics := models.ThreadMemberMetrics{
		ThreadMetrics:         statusMetrics,
		ThreadAssigneeMetrics: assignmentMetrics,
		ThreadLabelMetrics:    labelMetrics,
	}

	return metrics, nil
}
