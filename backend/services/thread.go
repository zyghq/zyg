package services

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type ThreadChatService struct {
	repo ports.ThreadRepositorer
}

func NewThreadChatService(
	repo ports.ThreadRepositorer) *ThreadChatService {
	return &ThreadChatService{
		repo: repo,
	}
}

// CreateNewInboundThreadChat CreateInboundThreadChat creates a new inbound thread chat for the customer.
// This is usually triggered when a customer sends a message.
// Inbound is always assumed as a customer message.
func (s *ThreadChatService) CreateNewInboundThreadChat(
	ctx context.Context, workspaceId string, customer models.Customer,
	createdBy models.MemberActor, message string,
) (models.Thread, models.Chat, error) {
	// Freeze the datetime for this transaction.
	// This is to ensure that all this is happening in the same time space.
	now := time.Now()
	channel := models.ThreadChannel{}.Chat() // channel the thread belongs to
	customerActor := customer.AsCustomerActor()

	// Create thread in the same time space.
	thread := (&models.Thread{}).CreateNewThread(workspaceId, customerActor, createdBy, createdBy, channel)
	thread.CreatedAt = now
	thread.UpdatedAt = now

	// Create chat for the thread in the same time space.
	// Mark the chat as head.
	chat := (&models.Chat{}).CreateNewCustomerChat(
		thread.ThreadId, customerActor.CustomerId, true, message)
	chat.CreatedAt = now
	chat.UpdatedAt = now

	// TODO: directly pass the customer actor, once method is updated.
	messageId := models.InboundMessage{}.GenId()
	seqId := ksuid.New().String()

	thread.AddInboundMessage(
		messageId, customerActor.CustomerId, customerActor.Name,
		chat.PreviewText(), seqId, seqId,
		now, now,
	)

	thread, chat, err := s.repo.InsertInboundThreadChat(ctx, thread, chat)
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
	ctx context.Context, customerId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	threads, err := s.repo.FetchThreadsByCustomerId(ctx, customerId, &channel)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadChatService) ListWorkspaceThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.Chat()
	threads, err := s.repo.FetchThreadsByWorkspaceId(ctx, workspaceId, &channel, nil)
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
	exist, err := s.repo.CheckThreadInWorkspaceExists(ctx, workspaceId, threadId)
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
	label, created, err := s.repo.SetThreadLabel(ctx, label)
	if err != nil {
		return models.ThreadLabel{}, created, ErrThChatLabel
	}

	return label, created, nil
}

func (s *ThreadChatService) ListThreadLabels(
	ctx context.Context, threadId string) ([]models.ThreadLabel, error) {
	labels, err := s.repo.FetchAttachedLabelsByThreadId(ctx, threadId)
	if err != nil {
		return labels, ErrThChatLabel
	}
	return labels, nil
}

// CreateInboundChatMessage adds an inbound message to the existing thread.
// Checks if the thread already has an inbound message reference otherwise creates a new one.
func (s *ThreadChatService) CreateInboundChatMessage(
	ctx context.Context, thread models.Thread, message string) (models.Chat, error) {
	var inboundMessage models.InboundMessage
	chat := models.Chat{
		ThreadId:   thread.ThreadId,
		Body:       message,
		CustomerId: models.NullString(&thread.Customer.CustomerId),
		IsHead:     false,
	}
	// If an existing inbound message already exists, then update
	// the existing inbound message with the latest value of last sequence ID.
	// Else create a new inbound message for the thread.
	if thread.InboundMessage != nil {
		inboundMessage = *thread.InboundMessage
		lastSeqId := ksuid.New().String()
		inboundMessage.LastSeqId = lastSeqId
		inboundMessage.PreviewText = chat.PreviewText()
	} else {
		seqId := ksuid.New().String()
		inboundMessage = models.InboundMessage{
			MessageId:   inboundMessage.GenId(),
			CustomerId:  thread.Customer.CustomerId,
			FirstSeqId:  seqId,
			LastSeqId:   seqId,
			PreviewText: chat.PreviewText(),
		}
	}
	chat, err := s.repo.InsertCustomerChat(ctx, thread, inboundMessage, chat)
	if err != nil {
		return models.Chat{}, ErrThChatMessage
	}
	return chat, nil
}

// CreateOutboundChatMessage creates an outbound message to the existing thread chat.
// Checks if the thread already has an outbound message reference otherwise creates a new one.
func (s *ThreadChatService) CreateOutboundChatMessage(
	ctx context.Context, thread models.Thread, memberId string, message string) (models.Chat, error) {
	var outboundMessage models.OutboundMessage
	chat := models.Chat{
		ThreadId: thread.ThreadId,
		Body:     message,
		MemberId: models.NullString(&memberId),
		IsHead:   false,
	}
	// If an existing outbound message already exists, then update
	// the existing with the latest value of last sequence ID,
	// else create a new outbound message for the thread chat.
	if thread.OutboundMessage != nil {
		outboundMessage = *thread.OutboundMessage
		lastSeqId := ksuid.New().String()
		outboundMessage.LastSeqId = lastSeqId
		outboundMessage.PreviewText = chat.PreviewText()
	} else {
		seqId := ksuid.New().String()
		outboundMessage = models.OutboundMessage{
			MessageId:   outboundMessage.GenId(),
			MemberId:    memberId,
			FirstSeqId:  seqId,
			LastSeqId:   seqId,
			PreviewText: chat.PreviewText(),
		}
	}
	chat, err := s.repo.InsertMemberChat(ctx, thread, outboundMessage, chat)
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
	err := s.repo.DeleteThreadLabelById(ctx, threadId, labelId)
	if err != nil {
		return ErrThChatLabel
	}
	return nil
}
