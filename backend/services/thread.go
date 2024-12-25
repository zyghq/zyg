package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/getsentry/sentry-go"
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
		models.SetMarkdownBody(messageText),
	)
	newThread.SetNextInboundSeq(newMessage.PreviewText())

	threadMessage := models.ThreadMessage{
		Thread:  newThread,
		Message: newMessage,
	}
	insThread, insMessage, err := s.repo.InsertInboundThreadMessage(ctx, threadMessage)
	if err != nil {
		return models.Thread{}, models.Message{}, ErrThreadChat
	}
	return insThread, insMessage, nil
}

func (s *ThreadService) GetPostmarkInReplyThread(
	ctx context.Context, workspaceId string, mailMessageId string) (*models.Thread, error) {
	thread, err := s.repo.FindThreadByPostmarkReplyMessageId(ctx, workspaceId, mailMessageId)
	if errors.Is(err, repository.ErrEmpty) {
		return nil, ErrThreadNotFound
	}
	if err != nil {
		slog.Error("failed to fetch thread by postmark inbound reply message id", slog.Any("err", err))
		return nil, ErrThread
	}
	return &thread, nil
}

func (s *ThreadService) IsPostmarkInboundMessageProcessed(ctx context.Context, messageId string) (bool, error) {
	exists, err := s.repo.CheckPostmarkInboundMessageExists(ctx, messageId)
	if err != nil {
		return false, ErrPostmarkInboundNotFound
	}
	return exists, nil
}

// ProcessPostmarkInbound processes an inbound Postmark email, handling threading logic, message creation, and attachments.
// It creates or finds an existing email thread based on the provided inbound message and processes HTML/Markdown conversions.
// If attachments exist, they are uploaded and persisted. Returns the updated or new thread and message, or an error if any.
func (s *ThreadService) ProcessPostmarkInbound(
	ctx context.Context, workspaceId string,
	customer models.CustomerActor, createdBy models.MemberActor, inboundMessage *models.PostmarkInboundMessage,
) (models.Thread, models.Message, error) {
	hub := sentry.GetHubFromContext(ctx)

	// Check if an existing thread already exists for the inbound Postmark based on reply mail message ID
	// otherwise, creates a new thread for the channel.
	channel := models.ThreadChannel{}.Email()
	var thread, threadExists, err = func(channel string) (*models.Thread, bool, error) {
		// Check if this inboundMessage is a reply to existing inboundMessage.
		// It's possible that reply mail message ID might exist for the inbound without
		// the corresponding thread in our system.
		if inboundMessage.ReplyMailMessageId != nil {
			// Get existing thread for Postmark inbound in-reply inboundMessage if exists.
			// Otherwise, creates a new thread.
			existingThread, err := s.GetPostmarkInReplyThread(ctx, workspaceId, *inboundMessage.ReplyMailMessageId)
			if errors.Is(err, ErrThreadNotFound) {
				slog.Info("thread not found for postmark inbound reply mail message ID should start new thread")
				newThread := models.NewThread(
					workspaceId, customer, createdBy, channel,
					models.SetThreadTitle(inboundMessage.Subject),
					models.SetThreadDescription(inboundMessage.TextBody),
				)
				return newThread, false, nil
			}
			if err != nil {
				hub.CaptureException(err)
				slog.Error(
					"failed to get existing thread for postmark inbound in-reply", slog.Any("err", err))
				return nil, false, ErrThread
			}
			// Returns existing thread.
			return existingThread, true, nil
		}
		// If inboundMessage is not a reply, start a new thread.
		newThread := models.NewThread(
			workspaceId, customer, createdBy, channel,
			models.SetThreadTitle(inboundMessage.Subject),
			models.SetThreadDescription(inboundMessage.TextBody),
		)
		// Return new thread.
		return newThread, false, nil
	}(channel)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to get existing thread or create one", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrThread
	}

	// Clean HTML into markdown - in case of error use HTML as fallback
	var markdownBody string
	cleanedHTML, err := utils.CleanHTML(inboundMessage.HTMLBody, utils.DefaultHTMLMatchers())
	if err != nil {
		slog.Error("failed to clean up postmark inbound mail html", slog.Any("err", err))
		markdownBody = inboundMessage.HTMLBody // Set fallback to raw HTMLBody
	} else {
		// Proceed to attempt converting the cleaned HTML to Markdown
		markdownBody, err = utils.HTMLToMarkdown(cleanedHTML)
		if err != nil {
			slog.Error("failed to convert html to markdown", slog.Any("err", err))
			markdownBody = inboundMessage.HTMLBody // Set fallback to raw HTMLBody
		}
	}

	newMessage := models.NewMessage(
		thread.ThreadId, channel,
		models.SetMessageCustomer(customer),
		models.SetHTMLBody(inboundMessage.HTMLBody),
		models.SetMessageTextBody(inboundMessage.TextBody),
		models.SetMarkdownBody(markdownBody),
	)
	thread.SetNextInboundSeq(newMessage.PreviewText())
	postmarkMessageLog := inboundMessage.NewPostmarkInboundLog(newMessage.MessageId)

	// If thread exists, append to the existing thread.
	if threadExists {
		newMessage, err = s.repo.AppendPostmarkInboundThreadMessage(
			ctx, thread.ThreadId, thread.InboundMessage, &postmarkMessageLog, newMessage)
		if err != nil {
			hub.CaptureException(err)
			slog.Error("failed to append postmark inbound message to existing thread", slog.Any("err", err))
			return models.Thread{}, models.Message{}, ErrPostmarkInbound
		}
	} else {
		thread, newMessage, err = s.repo.InsertPostmarkInboundThreadMessage(
			ctx, thread, &postmarkMessageLog, newMessage)
		if err != nil {
			hub.CaptureException(err)
			slog.Error("failed to insert postmark inbound message to new thread", slog.Any("err", err))
			return models.Thread{}, models.Message{}, ErrPostmarkInbound
		}
	}

	accountId := zyg.CFAccountId()
	accessKeyId := zyg.R2AccessKeyId()
	accessKeySecret := zyg.R2AccessSecretKey()
	s3Bucket := zyg.S3Bucket()
	s3Client, err := store.NewS3(ctx, s3Bucket, accountId, accessKeyId, accessKeySecret)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to connect S3 to process inbound message attachments", slog.Any("err", err))
		return *thread, *newMessage, nil
	}

	// Process attachments if any.
	if len(inboundMessage.Attachments) > 0 {
		attachments := make([]models.MessageAttachment, 0, len(inboundMessage.Attachments))
		for _, a := range inboundMessage.Attachments {
			att, attErr := ProcessMessageAttachment(
				ctx, thread.WorkspaceId, thread.ThreadId, newMessage.MessageId,
				a.Content, a.ContentType, a.Name, s3Client,
			)
			if attErr != nil {
				hub.Scope().SetTag("messageId", att.MessageId)
				hub.Scope().SetTag("attachmentId", att.AttachmentId)
				hub.Scope().SetTag("attachmentName", att.Name)
				hub.Scope().SetTag("attachmentMD5Hash", att.MD5Hash)
				hub.CaptureException(attErr)
				slog.Error(
					"failed to process inbound message attachment",
					slog.Any("err", attErr),
					slog.Any("attachmentId", att.AttachmentId),
				)
			}
			attachments = append(attachments, att)
		}
		// Persists processed inbound message attachments
		// @sanchitrk: bulk inserts?
		for _, a := range attachments {
			_, err := s.repo.InsertMessageAttachment(ctx, a)
			if err != nil {
				slog.Error(
					"failed to insert inbound message attachment", slog.Any("err", err))
			}
		}
	}
	return *thread, *newMessage, nil
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
		return models.ThreadLabel{}, created, ErrLabel
	}

	return label, created, nil
}

func (s *ThreadService) ListThreadLabels(
	ctx context.Context, threadId string) ([]models.ThreadLabel, error) {
	labels, err := s.repo.FetchAttachedLabelsByThreadId(ctx, threadId)
	if err != nil {
		return labels, ErrLabel
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
		models.SetMarkdownBody(messageText),
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
		models.SetMarkdownBody(messageText),
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

func (s *ThreadService) SendThreadMailReply(
	ctx context.Context,
	workspace models.Workspace,
	setting models.PostmarkMailServerSetting, thread models.Thread,
	member models.Member, customer models.Customer, textBody, htmlBody string,
) (models.Message, error) {
	hub := sentry.GetHubFromContext(ctx)

	fromName := fmt.Sprintf("%s at %s", member.Name, workspace.Name)

	// extract from HTML if text is empty
	// fallback to specified text in any case
	var textBodyFmt string
	if textBody != "" {
		textBodyFmt = textBody
	} else {
		extractedText, err := utils.ExtractTextFromHTML(htmlBody)
		if err != nil {
			hub.CaptureException(err)
			textBodyFmt = textBody
		} else {
			textBodyFmt = extractedText
		}
	}

	markdownBody, err := utils.HTMLToMarkdown(htmlBody)
	if err != nil {
		hub.CaptureMessage("failed to convert HTML to markdown for send reply mail")
		hub.CaptureException(err)
		slog.Error("failed to convert HTML to markdown for send reply mail", slog.Any("err", err))
		markdownBody = htmlBody // fallback to HTML
	}

	// Get recent Postmark mail message ID
	// The mail message ID is used in header for `In-Reply-To` maintaining a mail thread.
	mailMsgId, err := s.GetRecentPostmarkLogMailMessageId(ctx, thread.ThreadId)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to get recent postmark message log mail message ID", slog.Any("err", err))
		return models.Message{}, ErrPostmarkOutbound
	}

	fmt.Println("fromName", fromName)
	fmt.Println("subject", thread.Title)
	fmt.Println("textBody", textBody)
	fmt.Println("htmlBody", htmlBody)
	fmt.Println("markdownBody", markdownBody)
	fmt.Println("In-Reply-To", mailMsgId)

	newMessage := models.NewMessage(
		thread.ThreadId, models.ThreadChannel{}.Email(),
		models.SetMessageMember(member.AsMemberActor()),
		models.SetHTMLBody(htmlBody),
		models.SetMessageTextBody(textBodyFmt),
		models.SetMarkdownBody(markdownBody),
	)
	thread.SetNextOutboundSeq(member.AsMemberActor(), newMessage.PreviewText())

	threadMessage := models.ThreadMessage{
		Thread:  &thread,
		Message: newMessage,
	}

	fmt.Println("threadMessage", threadMessage)

	return models.Message{}, nil
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
		return ErrLabel
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

func (s *ThreadService) GetMessageAttachment(
	ctx context.Context, messageId, attachmentId string) (models.MessageAttachment, error) {
	attachment, err := s.repo.FetchMessageAttachmentById(ctx, messageId, attachmentId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.MessageAttachment{}, ErrMessageAttachmentNotFound
	}
	if err != nil {
		return models.MessageAttachment{}, ErrMessageAttachment
	}
	return attachment, nil
}

func (s *ThreadService) GetRecentPostmarkLogMailMessageId(
	ctx context.Context, threadId string) (string, error) {
	msgId, err := s.repo.FindRecentPostmarkLogMailMessageIdByThreadId(ctx, threadId)
	if errors.Is(err, repository.ErrEmpty) {
		return "", nil
	}
	fmt.Println("msgId", msgId)
	return "", nil
}
