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

func (s *ThreadChatService) CreateInboundThreadChat(
	ctx context.Context, workspaceId string, customerId string, message string) (models.Thread, models.Chat, error) {
	inbound := models.IngressMessage{
		CustomerId: customerId,
	}
	chat := models.Chat{
		Body:       message,
		CustomerId: models.NullString(&customerId),
		IsHead:     true,
	}
	thread := models.Thread{
		WorkspaceId: workspaceId,
		CustomerId:  customerId,
		Status:      models.ThreadStatus{}.Todo(),
		Priority:    models.ThreadPriority{}.Normal(),
		Channel:     models.ThreadChannel{}.Chat(),
		PreviewText: chat.PreviewText(),
	}
	thread, chat, err := s.repo.InsertInboundThreadChat(ctx, inbound, thread, chat)
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
	ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error) {
	thread, err := s.repo.LookupByWorkspaceThreadId(ctx, workspaceId, threadId, channel)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Thread{}, ErrThreadChatNotFound
	}

	if err != nil {
		return models.Thread{}, ErrThreadChat
	}

	return thread, nil
}

func (s *ThreadChatService) ListCustomerThreadChats(
	ctx context.Context, customerId string, role *string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	threads, err := s.repo.FetchThreadsByCustomerId(ctx, customerId, &channel, role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) AssignMember(
	ctx context.Context, threadId string, assigneeId string) (models.Thread, error) {
	thread, err := s.repo.UpdateAssignee(ctx, threadId, assigneeId)
	if err != nil {
		return models.Thread{}, ErrThreadChatAssign
	}
	return thread, nil
}

func (s *ThreadChatService) SetReplyStatus(
	ctx context.Context, threadChatId string, replied bool) (models.Thread, error) {
	thread, err := s.repo.UpdateRepliedState(ctx, threadChatId, replied)
	if err != nil {
		return models.Thread{}, ErrThreadChatReplied
	}
	return thread, nil
}

func (s *ThreadChatService) ListWorkspaceThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByWorkspaceId(ctx, workspaceId, &channel, &role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListMemberThreadChats(
	ctx context.Context, memberId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByAssignedMemberId(ctx, memberId, &channel, &role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListUnassignedThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByMemberUnassigned(ctx, workspaceId, &channel, &role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListLabelledThreadChats(
	ctx context.Context, labelId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByLabelId(ctx, labelId, &channel, &role)
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

func (s *ThreadChatService) SetLabel(
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

func (s *ThreadChatService) AddInboundMessage(
	ctx context.Context, thread models.Thread, customerId string, message string) (models.Chat, error) {
	var ingressMessageId *string
	if thread.IngressMessageId.Valid {
		ingressMessageId = &thread.IngressMessageId.String
	}
	chat := models.Chat{
		ThreadId:   thread.ThreadId,
		Body:       message,
		CustomerId: models.NullString(&customerId),
		IsHead:     false,
	}
	chat, err := s.repo.InsertCustomerChat(ctx, ingressMessageId, chat)
	if err != nil {
		return models.Chat{}, ErrThChatMessage
	}

	return chat, nil
}

func (s *ThreadChatService) AddOutboundMessage(
	ctx context.Context, thread models.Thread, memberId string, message string) (models.Chat, error) {
	var outboundMessageId *string
	if thread.IngressMessageId.Valid {
		outboundMessageId = &thread.EgressMessageId.String
	}
	chat := models.Chat{
		ThreadId: thread.ThreadId,
		Body:     message,
		MemberId: models.NullString(&memberId),
		IsHead:   false,
	}
	chat, err := s.repo.InsertMemberChat(ctx, outboundMessageId, chat)
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

func (s *ThreadChatService) RemoveThreadLabel(
	ctx context.Context, threadId string, labelId string) error {
	err := s.repo.DeleteThreadLabelByCompId(ctx, threadId, labelId)
	if err != nil {
		return ErrThChatLabel
	}
	return nil
}
