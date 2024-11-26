package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/store"
	"github.com/zyghq/zyg/utils"
	"log/slog"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type ThreadService struct {
	repo ports.ThreadRepositorer
}

func NewThreadService(
	repo ports.ThreadRepositorer) *ThreadService {
	return &ThreadService{
		repo: repo,
	}
}

// CreateInboundThreadChat creates a new inbound thread chat for the customer.
// This is usually triggered when a customer sends a message.
// Inbound is always assumed as a customer message.
func (s *ThreadService) CreateInboundThreadChat(
	ctx context.Context, workspaceId string, customer models.Customer,
	createdBy models.MemberActor, messageText string,
) (models.Thread, models.Message, error) {
	// Creates new thread for the in-app chat channel.
	channel := models.ThreadChannel{}.InAppChat() // source channel the thread belongs to
	newThread := models.NewThread(
		workspaceId, customer.AsCustomerActor(), createdBy, channel,
	)
	// Create new Customer inbound message.
	newMessage := models.NewMessage(
		newThread.ThreadId, channel,
		models.SetMessageCustomer(customer.AsCustomerActor()),
		models.SetMessageTextBody(messageText),
		models.SetMessageBody(messageText),
	)
	// Set new inbound message sequence.
	newThread.SetNewInboundMessage(customer.AsCustomerActor(), newMessage.PreviewText())

	inbound := models.ThreadMessage{
		Thread:  newThread,
		Message: newMessage,
	}
	insThread, insMessage, err := s.repo.InsertInboundThreadMessage(ctx, inbound)
	if err != nil {
		return models.Thread{}, models.Message{}, ErrThreadChat
	}
	return insThread, insMessage, nil
}

func (s *ThreadService) GetPostmarkInboundInReplyThread(
	ctx context.Context, workspaceId string, inboundMessage *models.PostmarkInboundMessage) (*models.Thread, error) {
	if inboundMessage == nil || inboundMessage.ReplyMailMessageId == nil {
		return nil, ErrPostmarkInbound
	}
	thread, err := s.repo.FetchThreadByPostmarkInboundInReplyMessageId(
		ctx, workspaceId, *inboundMessage.ReplyMailMessageId)
	if errors.Is(err, repository.ErrEmpty) {
		return nil, ErrThreadNotFound
	}
	if err != nil {
		slog.Error("failed to fetch thread by postmark inbound reply message id", slog.Any("err", err))
		return nil, ErrThread
	}
	return &thread, nil
}

func (s *ThreadService) IsPostmarkInboundProcessed(ctx context.Context, messageId string) (bool, error) {
	exists, err := s.repo.CheckPostmarkInboundMessageExists(ctx, messageId)
	if err != nil {
		return false, ErrPostmarkInbound
	}
	return exists, nil
}

// htmlToMarkdown converts an HTML string to a Markdown string.
// Parameters:
// - html: A string containing the HTML content to convert.
// Returns:
// - A string containing the converted Markdown content.
// - An error if the conversion fails.
func htmlToMarkdown(html string) (string, error) {
	cleaned, err := utils.CleanHTML(html, utils.DefaultHTMLMatchers())
	if err != nil {
		return "", err
	}
	m, err := utils.HTMLToMarkdown(cleaned)
	if err != nil {
		return "", err
	}
	return m, nil
}

func (s *ThreadService) ProcessPostmarkInbound(
	ctx context.Context, workspaceId string,
	customer models.CustomerActor, createdBy models.MemberActor, message *models.PostmarkInboundMessage,
) (models.Thread, models.Message, error) {
	var insMessage models.Message
	if message == nil {
		return models.Thread{}, models.Message{}, ErrPostmarkInbound
	}

	channel := models.ThreadChannel{}.Email()
	var thread, exists, err = func(channel string) (*models.Thread, bool, error) {
		// Check if this message is a reply to existing message.
		if message.ReplyMailMessageId != nil {
			// Get existing thread for Postmark inbound in-reply message if exists.
			// Otherwise, creates a new thread.
			existingThread, err := s.GetPostmarkInboundInReplyThread(ctx, workspaceId, message)
			if errors.Is(err, ErrThreadNotFound) {
				newThread := models.NewThread(
					workspaceId, customer, createdBy, channel,
					models.SetThreadTitle(message.Subject),
					models.SetThreadDescription(message.TextBody),
				)
				return newThread, false, nil
			}
			if err != nil {
				return nil, false, ErrThread
			}
			// Returns existing thread.
			return existingThread, true, nil
		}
		// If message is not a reply, start a new thread.
		newThread := models.NewThread(
			workspaceId, customer, createdBy, channel,
			models.SetThreadTitle(message.Subject),
			models.SetThreadDescription(message.TextBody),
		)
		// Return new thread.
		return newThread, false, nil
	}(channel)
	if err != nil {
		slog.Error("failed to get or create thread", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrThread
	}

	body := message.HTMLBody
	markdown, err := htmlToMarkdown(message.HTMLBody)
	if err != nil {
		slog.Error("failed to convert html to markdown", slog.Any("err", err))
	} else {
		body = markdown
	}

	newMessage := models.NewMessage(
		thread.ThreadId, channel,
		models.SetMessageCustomer(customer),
		models.SetMessageTextBody(message.TextBody),
		models.SetMessageBody(body),
	)

	// Check if the thread has inbound message sequence.
	// Set next new inbound sequence else create new inbound sequence.
	if thread.InboundMessage != nil {
		thread.SetNextInboundSeq(newMessage.PreviewText())
	} else {
		thread.SetNewInboundMessage(customer, newMessage.PreviewText())
	}

	inbound := models.PostmarkInboundThreadMessage{
		ThreadMessage: models.ThreadMessage{
			Thread:  thread,
			Message: newMessage,
		},
		PostmarkInboundMessage: message,
	}

	// If thread exists, append message to the existing thread.
	if exists {
		insMessage, err = s.repo.AppendPostmarkInboundThreadMessage(ctx, inbound)
		if err != nil {
			return models.Thread{}, models.Message{}, ErrPostmarkInbound
		}
	} else {
		// Insert a new thread.
		var insThread models.Thread
		insThread, insMessage, err = s.repo.InsertPostmarkInboundThreadMessage(ctx, inbound)
		if err != nil {
			return models.Thread{}, models.Message{}, ErrPostmarkInbound
		}
		// Set the new inserted thread.
		thread = &insThread
	}

	// Process attachments if any.
	if len(inbound.PostmarkInboundMessage.Attachments) > 0 {
		attachments := make([]models.MessageAttachment, 0, len(inbound.PostmarkInboundMessage.Attachments))
		accountId := zyg.CFAccountId()
		accessKeyId := zyg.R2AccessKeyId()
		accessKeySecret := zyg.R2AccessSecretKey()
		s3Bucket := zyg.S3Bucket()
		s3Client, err := store.NewS3(ctx, s3Bucket, accountId, accessKeyId, accessKeySecret)
		if err != nil {
			slog.Error("failed to connect s3 to process message attachments", slog.Any("err", err))
		} else {
			for _, a := range inbound.PostmarkInboundMessage.Attachments {
				att, aErr := ProcessMessageAttachment(
					ctx, thread.WorkspaceId, thread.ThreadId, insMessage.MessageId,
					a.Content, a.ContentType, a.Name, s3Client,
				)
				if aErr != nil {
					slog.Error("failed to process message attachment", slog.Any("err", aErr))
				}
				attachments = append(attachments, att)
			}
		}
		// Persists processed message attachments
		for _, a := range attachments {
			_, err := s.repo.UpsertMessageAttachment(ctx, a)
			if err != nil {
				slog.Error("failed to insert message attachment", slog.Any("err", err))
			}
		}
	}
	return *thread, insMessage, nil
}

func (s *ThreadService) UpdateThread(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {
	thread, err := s.repo.ModifyThreadById(ctx, thread, fields)

	if err != nil {
		return models.Thread{}, ErrThread
	}

	return thread, nil
}

func (s *ThreadService) GetWorkspaceThread(
	ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error) {
	thread, err := s.repo.LookupByWorkspaceThreadId(ctx, workspaceId, threadId, channel)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Thread{}, ErrThreadNotFound
	}
	if err != nil {
		return models.Thread{}, ErrThread
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

func (s *ThreadService) ListWorkspaceThreads(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByWorkspaceId(ctx, workspaceId, nil, &role)
	if err != nil {
		return []models.Thread{}, ErrThread
	}
	return threads, nil
}

func (s *ThreadService) ListMemberThreads(
	ctx context.Context, memberId string) ([]models.Thread, error) {
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByAssignedMemberId(ctx, memberId, nil, &role)
	if err != nil {
		return []models.Thread{}, ErrThread
	}
	return threads, nil
}

func (s *ThreadService) ListUnassignedThreads(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByMemberUnassigned(ctx, workspaceId, nil, &role)
	if err != nil {
		return []models.Thread{}, ErrThread
	}
	return threads, nil
}

func (s *ThreadService) ListLabelledThreads(
	ctx context.Context, labelId string) ([]models.Thread, error) {
	role := models.Customer{}.Engaged()
	threads, err := s.repo.FetchThreadsByLabelId(ctx, labelId, nil, &role)
	if err != nil {
		return []models.Thread{}, ErrThread
	}
	return threads, nil
}

func (s *ThreadService) ThreadExistsInWorkspace(
	ctx context.Context, workspaceId string, threadId string) (bool, error) {
	exist, err := s.repo.CheckThreadInWorkspaceExists(ctx, workspaceId, threadId)
	if err != nil {
		return false, ErrThread
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
		return models.ThreadLabel{}, created, ErrThreadLabel
	}

	return label, created, nil
}

func (s *ThreadService) ListThreadLabels(
	ctx context.Context, threadId string) ([]models.ThreadLabel, error) {
	labels, err := s.repo.FetchAttachedLabelsByThreadId(ctx, threadId)
	if err != nil {
		return labels, ErrThreadLabel
	}
	return labels, nil
}

// AppendInboundThreadChat adds inbound message to an existing thread.
func (s *ThreadService) AppendInboundThreadChat(
	ctx context.Context, thread models.Thread, messageText string) (models.Message, error) {

	channel := models.ThreadChannel{}.InAppChat()
	newMessage := models.NewMessage(
		thread.ThreadId, channel,
		models.SetMessageCustomer(thread.Customer),
		models.SetMessageTextBody(messageText),
		models.SetMessageBody(messageText),
	)
	thread.SetNextInboundSeq(newMessage.PreviewText())

	threadMessage := models.ThreadMessage{
		Thread:  &thread,
		Message: newMessage,
	}
	message, err := s.repo.AppendInboundThreadMessage(ctx, threadMessage)
	if err != nil {
		return models.Message{}, ErrThreadMessage
	}
	return message, nil
}

// AppendOutboundThreadChat appends a new chat message to an existing thread.
func (s *ThreadService) AppendOutboundThreadChat(
	ctx context.Context, thread models.Thread, member models.Member, messageText string) (models.Message, error) {

	channel := models.ThreadChannel{}.InAppChat()
	newMessage := models.NewMessage(
		thread.ThreadId, channel,
		models.SetMessageMember(member.AsMemberActor()),
		models.SetMessageTextBody(messageText),
		models.SetMessageBody(messageText),
	)
	thread.SetNextInboundSeq(newMessage.PreviewText())

	threadMessage := models.ThreadMessage{
		Thread:  &thread,
		Message: newMessage,
	}
	message, err := s.repo.AppendOutboundThreadMessage(ctx, threadMessage)
	if err != nil {
		return models.Message{}, ErrThreadMessage
	}
	return message, nil
}

func (s *ThreadService) ListThreadMessages(
	ctx context.Context, threadId string) ([]models.Message, error) {
	messages, err := s.repo.FetchMessagesByThreadId(ctx, threadId)
	if err != nil {
		return []models.Message{}, ErrThreadMessage
	}
	return messages, nil
}

func (s *ThreadService) ListThreadMessagesWithAttachments(
	ctx context.Context, threadId string) ([]models.MessageWithAttachments, error) {
	messages, err := s.repo.FetchMessagesWithAttachmentsByThreadId(ctx, threadId)
	if err != nil {
		return []models.MessageWithAttachments{}, ErrThreadMessage
	}
	return messages, nil
}

func (s *ThreadService) GenerateMemberThreadMetrics(
	ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error) {
	statusMetrics, err := s.repo.ComputeStatusMetricsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadMetrics
	}

	assignmentMetrics, err := s.repo.ComputeAssigneeMetricsByMember(ctx, workspaceId, memberId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadMetrics
	}

	labelMetrics, err := s.repo.ComputeLabelMetricsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadMetrics
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
		return ErrThreadLabel
	}
	return nil
}

func (s *ThreadService) LogPostmarkInboundRequest(
	ctx context.Context, workspaceId string, messageId string, payload map[string]interface{}) error {
	accountId := zyg.CFAccountId()
	accessKeyId := zyg.R2AccessKeyId()
	accessSecretKey := zyg.R2AccessSecretKey()

	s3Client, err := store.NewS3(ctx, "zygdev", accountId, accessKeyId, accessSecretKey)
	if err != nil {
		return fmt.Errorf("failed to create S3: %v", err)
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// In format: <workspaceId>/logs/<integration>/<event>/<id>
	bucketKey := fmt.Sprintf("%s/logs/postmark/inbound/%s.json", workspaceId, messageId)

	input := &s3.PutObjectInput{
		Bucket:      &s3Client.BucketName,
		Key:         &bucketKey,
		Body:        bytes.NewReader(jsonData),
		ContentType: aws.String("application/json"),
	}

	_, err = s3Client.Client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put object: %v", err)
	}
	return nil
}
