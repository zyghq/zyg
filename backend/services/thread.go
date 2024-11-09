package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type ThreadService struct {
	repo ports.ThreadRepositorer
}

func NewThreadChatService(
	repo ports.ThreadRepositorer) *ThreadService {
	return &ThreadService{
		repo: repo,
	}
}

// CreateNewInboundThreadChat CreateInboundThreadChat creates a new inbound thread chat for the customer.
// This is usually triggered when a customer sends a message.
// Inbound is always assumed as a customer message.
func (s *ThreadService) CreateNewInboundThreadChat(
	ctx context.Context, workspaceId string, customer models.Customer,
	createdBy models.MemberActor, message string,
) (models.Thread, models.Chat, error) {
	// Freeze the datetime for this transaction.
	// This is to ensure that Thread and Chat are happening in the same time space.
	now := time.Now().UTC()
	channel := models.ThreadChannel{}.InAppChat() // source channel the thread belongs to
	customerActor := customer.AsCustomerActor()

	newThread := models.NewThread(
		workspaceId, customer.AsCustomerActor(), createdBy, channel,
	)

	// Create chat for the thread in the same time space.
	// Mark the chat as head.
	chat := (&models.Chat{}).CreateNewCustomerChat(
		newThread.ThreadId, customerActor.CustomerId, true, message,
		now, now,
	)

	newThread.SetNewInboundMessage(customerActor, chat.PreviewText())

	thread, chat, err := s.repo.InsertInboundThreadChat(ctx, *newThread, chat)
	if err != nil {
		return models.Thread{}, models.Chat{}, ErrThreadChat
	}
	return thread, chat, nil
}

func (s *ThreadService) GetPostmarkInboundInReplyThread(
	ctx context.Context, workspaceId string, inboundMessage *models.PostmarkInboundMessage) (*models.Thread, error) {
	if inboundMessage == nil || inboundMessage.ReplyMailMessageId == nil {
		return nil, ErrPostmarkInbound
	}
	thread, err := s.repo.FetchThreadByPostmarkInboundInReplyMessageId(ctx, workspaceId, *inboundMessage.ReplyMailMessageId)
	if errors.Is(err, repository.ErrEmpty) {
		return nil, ErrThreadNotFound
	}
	if err != nil {
		slog.Error("failed to fetch thread by postmark inbound reply message id", slog.Any("err", err))
		return nil, ErrThread
	}
	return &thread, nil
}

func (s *ThreadService) ProcessPostmarkInbound(
	ctx context.Context, workspaceId string,
	customer models.CustomerActor, createdBy models.MemberActor, inboundMessage *models.PostmarkInboundMessage,
) (models.Thread, models.Message, error) {
	if inboundMessage == nil {
		return models.Thread{}, models.Message{}, ErrPostmarkInbound
	}

	channel := models.ThreadChannel{}.Email()
	var thread, exists, err = func(channel string) (*models.Thread, bool, error) {
		// Check for in-reply mail message ID for existing reply message.
		if inboundMessage.ReplyMailMessageId != nil {
			// Get existing thread for postmark inbound in-reply mail message if exists.
			thread, err := s.GetPostmarkInboundInReplyThread(ctx, workspaceId, inboundMessage)
			if errors.Is(err, ErrThreadNotFound) {
				// Create a new thread
				thread := models.NewThread(
					workspaceId, customer, createdBy, channel,
					models.SetThreadTitle(inboundMessage.Subject()),
					models.SetThreadDescription(inboundMessage.PlainText()),
				)
				return thread, false, nil
			}
			// Something went wrong
			if err != nil {
				return nil, false, ErrThread
			}
			// Return existing thread, exists flag and error.
			return thread, true, nil
		}
		// Create new thread
		thread := models.NewThread(
			workspaceId, customer, createdBy, channel,
			models.SetThreadTitle(inboundMessage.Subject()),
			models.SetThreadDescription(inboundMessage.PlainText()),
		)
		return thread, false, nil
	}(channel)
	if err != nil {
		slog.Error("failed to get or create thread", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrThread
	}

	message := models.NewMessage(
		thread.ThreadId, channel,
		models.SetMessageTextBody(inboundMessage.PlainText()),
		models.SetMessageBody(inboundMessage.Html()),
		models.SetMessageCustomer(customer),
	)

	// Check if the thread has inbound message sequence.
	// Set next new inbound sequence else create new inbound sequence.
	if thread.InboundMessage != nil {
		thread.SetNextInboundSeq(message.TextBody)
	} else {
		thread.SetNewInboundMessage(customer, message.TextBody)
	}

	inbound := models.ThreadMessageWithPostmarkInbound{
		Thread:                 thread,
		Message:                message,
		PostmarkInboundMessage: inboundMessage,
	}

	// If thread exists, append postmark inbound message to existing thread.
	if exists {
		insMessage, err := s.repo.AppendPostmarkInboundThreadMessage(ctx, inbound)
		if err != nil {
			return models.Thread{}, models.Message{}, ErrPostmarkInbound
		}
		return *thread, insMessage, nil
	}

	insThread, insMessage, err := s.repo.InsertPostmarkInboundThreadMessage(ctx, inbound)
	if err != nil {
		return models.Thread{}, models.Message{}, ErrPostmarkInbound
	}
	return insThread, insMessage, nil
}

func (s *ThreadService) UpdateThread(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {
	thread, err := s.repo.ModifyThreadById(ctx, thread, fields)

	if err != nil {
		return models.Thread{}, ErrThreadChat
	}

	return thread, nil
}

func (s *ThreadService) GetWorkspaceThread(
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

func (s *ThreadService) ListCustomerThreadChats(
	ctx context.Context, customerId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.InAppChat()
	threads, err := s.repo.FetchThreadsByCustomerId(ctx, customerId, &channel)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadService) ListWorkspaceThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.InAppChat()
	threads, err := s.repo.FetchThreadsByWorkspaceId(ctx, workspaceId, &channel, nil)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadService) ListMemberThreadChats(
	ctx context.Context, memberId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.InAppChat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByAssignedMemberId(ctx, memberId, &channel, &role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadService) ListUnassignedThreadChats(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.InAppChat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByMemberUnassigned(ctx, workspaceId, &channel, &role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadService) ListLabelledThreadChats(
	ctx context.Context, labelId string) ([]models.Thread, error) {
	channel := models.ThreadChannel{}.InAppChat()
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByLabelId(ctx, labelId, &channel, &role)
	if err != nil {
		return []models.Thread{}, ErrThreadChat
	}
	return threads, nil
}

func (s *ThreadService) ThreadExistsInWorkspace(
	ctx context.Context, workspaceId string, threadId string) (bool, error) {
	exist, err := s.repo.CheckThreadInWorkspaceExists(ctx, workspaceId, threadId)
	if err != nil {
		return false, ErrThreadChat
	}
	return exist, nil
}

func (s *ThreadService) SetLabel(
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

func (s *ThreadService) ListThreadLabels(
	ctx context.Context, threadId string) ([]models.ThreadLabel, error) {
	labels, err := s.repo.FetchAttachedLabelsByThreadId(ctx, threadId)
	if err != nil {
		return labels, ErrThChatLabel
	}
	return labels, nil
}

// CreateInboundChatMessage adds an inbound message to the existing thread.
// Checks if the thread already has an inbound message reference otherwise creates a new one.
func (s *ThreadService) CreateInboundChatMessage(
	ctx context.Context, thread models.Thread, message string) (models.Chat, error) {
	chat := models.Chat{
		ThreadId:   thread.ThreadId,
		Body:       message,
		CustomerId: models.NullString(&thread.Customer.CustomerId),
		IsHead:     false,
	}

	if thread.InboundMessage != nil {
		// Set next thread inbound message sequence.
		thread.SetNextInboundSeq(chat.PreviewText())
	} else {
		thread.SetNewInboundMessage(thread.Customer, chat.PreviewText())
	}

	chat, err := s.repo.InsertCustomerChat(ctx, thread, chat)
	if err != nil {
		return models.Chat{}, ErrThChatMessage
	}
	return chat, nil
}

// CreateOutboundChatMessage creates an outbound message to the existing thread chat.
// Checks if the thread already has an outbound message reference otherwise creates a new one.
func (s *ThreadService) CreateOutboundChatMessage(
	ctx context.Context, thread models.Thread, member models.Member, message string) (models.Chat, error) {
	var outboundMessage models.OutboundMessage
	chat := models.Chat{
		ThreadId: thread.ThreadId,
		Body:     message,
		MemberId: models.NullString(&member.MemberId),
		IsHead:   false,
	}
	// If an existing outbound message already exists, then update
	// the existing with the latest value of last sequence ID,
	// else create a new outbound message for the thread chat.
	if thread.OutboundMessage != nil {
		outboundMessage = *thread.OutboundMessage
		// Deprecate usage use xid.
		lastSeqId := ksuid.New().String()
		outboundMessage.LastSeqId = lastSeqId
		outboundMessage.PreviewText = chat.PreviewText()
	} else {
		// Deprecate usage use xid.
		seqId := ksuid.New().String()
		outboundMessage = models.OutboundMessage{
			MessageId:   outboundMessage.GenId(),
			Member:      member.AsMemberActor(),
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

func (s *ThreadService) ListThreadChatMessages(
	ctx context.Context, threadId string) ([]models.Chat, error) {
	messages, err := s.repo.FetchThChatMessagesByThreadId(ctx, threadId)
	if err != nil {
		return []models.Chat{}, ErrThChatMessage
	}
	return messages, nil
}

func (s *ThreadService) GenerateMemberThreadMetrics(
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

func (s *ThreadService) RemoveThreadLabel(
	ctx context.Context, threadId string, labelId string) error {
	err := s.repo.DeleteThreadLabelById(ctx, threadId, labelId)
	if err != nil {
		return ErrThChatLabel
	}
	return nil
}
