package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/store"
	"github.com/zyghq/zyg/integrations/email"
	"github.com/zyghq/zyg/utils"

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

func (s *ThreadService) IsPostmarkInboundProcessed(ctx context.Context, pmMessageId string) (bool, error) {
	exists, err := s.repo.CheckPostmarkInboundExists(ctx, pmMessageId)
	if err != nil {
		return false, ErrPostmarkInbound
	}
	return exists, nil
}

func (s *ThreadService) ProcessPostmarkInbound(
	ctx context.Context, workspaceId string,
	customer *models.Customer, createdBy *models.Member,
	inboundMessage *models.PostmarkInboundMessage) (models.Thread, models.Activity, error) {
	hub := sentry.GetHubFromContext(ctx)

	// Check if an existing thread already exists for the inbound Postmark based on reply mail message ID
	// otherwise, creates a new thread for the channel.
	channel := models.ThreadChannel{}.Email()
	var thread, threadExists, err = func() (*models.Thread, bool, error) {
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
					workspaceId,
					customer.AsCustomerActor(), createdBy.AsMemberActor(), channel,
					models.SetThreadTitle(inboundMessage.Subject),
					models.SetThreadDescription(inboundMessage.TextBody),
					models.SetThreadInboundTime(time.Now().UTC()), // set inbound time
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
			workspaceId, customer.AsCustomerActor(), createdBy.AsMemberActor(), channel,
			models.SetThreadTitle(inboundMessage.Subject),
			models.SetThreadDescription(inboundMessage.TextBody),
			models.SetThreadInboundTime(time.Now().UTC()), // set inbound time
		)
		// Return new thread.
		return newThread, false, nil
	}()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to get existing thread or create one", slog.Any("err", err))
		return models.Thread{}, models.Activity{}, ErrThread
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
	// create new message
	newMessage := models.NewMessage(
		channel,
		models.SetHTMLBody(inboundMessage.HTMLBody),
		models.SetMessageTextBody(inboundMessage.TextBody),
		models.SetMarkdownBody(markdownBody),
	)
	thread.SetPreviewText(inboundMessage.TextBody) // modify thread preview text

	// create new activity that embeds the message
	newActivity := models.NewActivity(
		thread.ThreadId,
		models.ActivityThreadMessage,
		models.SetActivityCustomer(customer.AsCustomerActor()),
		models.SetActivityBody(newMessage.ToJSON()),
	)

	// Convert postmark inbound message into Postmark message log.
	// The is persisted for both inbound and outbound messages.
	// Inbound as received from Postmark.
	postmarkMessageLog := inboundMessage.ToPostmarkMessageLog(newActivity.ActivityID)

	// If thread exists, append to the existing thread.
	if threadExists {
		thread.SetLatestInboundAt() // update latest inbound time for existing thread
		thread, newActivity, err = s.repo.SavePostmarkThreadActivity(
			ctx, thread, newActivity, &postmarkMessageLog)
		if err != nil {
			hub.CaptureException(err)
			slog.Error("failed to append postmark inbound to existing thread", slog.Any("err", err))
			return models.Thread{}, models.Activity{}, ErrPostmarkInbound
		}
	} else {
		// insert new thread with message
		thread, newActivity, err = s.repo.SavePostmarkThreadActivity(
			ctx, thread, newActivity, &postmarkMessageLog)
		if err != nil {
			hub.CaptureException(err)
			slog.Error("failed to insert postmark inbound to new thread", slog.Any("err", err))
			return models.Thread{}, models.Activity{}, ErrPostmarkInbound
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
		return *thread, *newActivity, nil
	}

	// Process attachments if any, each attachment is part of the activity.
	if len(inboundMessage.Attachments) > 0 {
		attachments := make([]models.ActivityAttachment, 0, len(inboundMessage.Attachments))
		for _, a := range inboundMessage.Attachments {
			att, attErr := ProcessMessageAttachment(
				ctx, thread.WorkspaceId, thread.ThreadId, newActivity.ActivityID,
				a.Content, a.ContentType, a.Name, s3Client,
			)
			if attErr != nil {
				hub.Scope().SetTag("activityId", att.ActivityID)
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
	return *thread, *newActivity, nil
}

func (s *ThreadService) SendThreadMailReply(
	ctx context.Context,
	workspace *models.Workspace,
	setting *models.PostmarkServerSetting, thread *models.Thread,
	member *models.Member, customer *models.Customer,
	textBody, htmlBody string) (models.Thread, models.Activity, error) {
	hub := sentry.GetHubFromContext(ctx)

	// make sure customer has email to send.
	if !customer.HasEmail() {
		slog.Error("customer has not valid email - cannot send reply mail",
			slog.Any("customerId", customer.CustomerId), slog.Any("email", customer.Email))
		hub.Scope().SetTag("customerId", customer.CustomerId)
		hub.CaptureMessage("customer has not valid email - cannot send reply mail")
		return models.Thread{}, models.Activity{}, ErrPostmarkOutbound
	}

	// set email params
	fromName := fmt.Sprintf("%s at %s", member.Name, workspace.Name)
	from := fmt.Sprintf("%s <%s>", fromName, setting.Email)
	replySubject := fmt.Sprintf("Re: %s", thread.Title)

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
		hub.CaptureMessage("failed to convert HTML to markdown for reply mail")
		hub.CaptureException(err)
		slog.Error("failed to convert HTML to markdown for reply mail", slog.Any("err", err))
		markdownBody = htmlBody // fallback to HTML
	}

	// Get most recent Postmark mail message ID
	// The mail message ID is used in header for `In-Reply-To` maintaining a mail thread.
	replyMailMessageId, err := s.GetRecentThreadMailMessageId(ctx, thread.ThreadId)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to get recent postmark message log mail message ID", slog.Any("err", err))
		return models.Thread{}, models.Activity{}, ErrPostmarkOutbound
	}
	thread.SetLatestOutboundAt() // set outbound time
	thread.SetPreviewText(textBodyFmt)

	newMessage := models.NewMessage(
		models.ThreadChannel{}.Email(),
		models.SetHTMLBody(htmlBody),
		models.SetMessageTextBody(textBodyFmt),
		models.SetMarkdownBody(markdownBody),
	)
	pmEmailReq := email.NewPostmarkEmailReq(
		replySubject, from, customer.Email.String,
		email.WithPostmarkHeader("In-Reply-To", replyMailMessageId),
		email.SetPostmarkTag(customer.CustomerId),
		email.SetPostmarkTextBody(newMessage.TextBody),
		email.SetPostmarkHTMLBody(newMessage.HTMLBody),
	)

	payload, err := utils.StructToMap(pmEmailReq)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to marshal postmark email request", slog.Any("err", err))
		return models.Thread{}, models.Activity{}, ErrPostmarkOutbound
	}

	resp, err := email.SendPostmarkMail(ctx, setting, pmEmailReq)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to send postmark email", slog.Any("err", err))
		return models.Thread{}, models.Activity{}, ErrPostmarkOutbound
	}

	// create new activity that embeds the message
	newActivity := models.NewActivity(
		thread.ThreadId,
		models.ActivityThreadMessage,
		models.SetActivityMember(member.AsMemberActor()),
		models.SetActivityBody(newMessage.ToJSON()),
	)

	messageLog := models.PostmarkMessageLog{
		ActivityID:         newActivity.ActivityID,
		Payload:            payload,
		PostmarkMessageId:  resp.MessageID,
		ReplyMailMessageId: &replyMailMessageId,
		HasError:           false,
		SubmittedAt:        resp.SubmittedAt,
		ErrorCode:          resp.ErrorCode,
		PostmarkMessage:    resp.Message,
		Acknowledged:       false,
		MessageType:        "outbound",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	// set outbound Mail Message ID for reply message with Postmark delivery domain
	messageLog.SetOutboundMailMessageId(zyg.PostmarkDeliveryDomain())

	thread, newActivity, err = s.repo.SavePostmarkThreadActivity(ctx, thread, newActivity, &messageLog)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to append postmark inbound thread message", slog.Any("err", err))
		return models.Thread{}, models.Activity{}, ErrPostmarkOutbound
	}
	return *thread, *newActivity, nil
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

//func (s *ThreadService) ListCustomerThreadChats(
//	ctx context.Context, customerId string) ([]models.Thread, error) {
//	channel := models.ThreadChannel{}.InAppChat()
//	threads, err := s.repo.FetchThreadsByCustomerId(ctx, customerId, &channel)
//	if err != nil {
//		return []models.Thread{}, ErrThreadChat
//	}
//	return threads, nil
//}

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

func (s *ThreadService) ListThreadMessageActivities(
	ctx context.Context, threadId string) ([]models.Activity, error) {
	activities, err := s.repo.FetchMessagesByThreadId(ctx, threadId)
	if err != nil {
		return []models.Activity{}, ErrThreadActivity
	}
	return activities, nil
}

func (s *ThreadService) ListThreadMessagesWithAttachments(
	ctx context.Context, threadId string) ([]models.ActivityWithAttachments, error) {
	messages, err := s.repo.FetchMessagesWithAttachmentsByThreadId(ctx, threadId)
	if err != nil {
		return []models.ActivityWithAttachments{}, ErrThreadActivity
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
	ctx context.Context, messageId, attachmentId string) (models.ActivityAttachment, error) {
	attachment, err := s.repo.FetchMessageAttachmentById(ctx, messageId, attachmentId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.ActivityAttachment{}, ErrMessageAttachmentNotFound
	}
	if err != nil {
		return models.ActivityAttachment{}, ErrMessageAttachment
	}
	return attachment, nil
}

func (s *ThreadService) GetRecentThreadMailMessageId(ctx context.Context, threadId string) (string, error) {
	mailMessageId, err := s.repo.GetRecentMailMessageIdByThreadId(ctx, threadId)
	if errors.Is(err, repository.ErrEmpty) {
		return "", ErrPostmarkLogNotFound
	}
	if err != nil {
		return "", ErrPostmarkLog
	}
	return mailMessageId, nil
}
