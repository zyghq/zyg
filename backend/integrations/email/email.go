package email

import (
	"encoding/json"
	"github.com/zyghq/postmark"
	"github.com/zyghq/zyg/models"
	"time"
)

type PostmarkInboundMessageReq struct {
	postmark.InboundWebhook
	Payload map[string]interface{}
}

// ToPostmarkInboundMessage converts a PostmarkInboundMessageReq to a PostmarkInboundMessage model instance.
func (p *PostmarkInboundMessageReq) ToPostmarkInboundMessage() models.PostmarkInboundMessage {
	now := time.Now().UTC()
	//var textBody string
	//if p.StrippedTextReply != "" {
	//	textBody = p.StrippedTextReply
	//} else {
	//	textBody = p.TextBody
	//}
	message := models.PostmarkInboundMessage{
		Payload:           p.Payload,
		PostmarkMessageId: p.MessageID, // Postmark MessageID
		Subject:           p.Subject,
		TextBody:          p.TextBody,
		HTMLBody:          p.HTMLBody,
		FromEmail:         p.FromFull.Email,
		FromName:          p.FromFull.Name,
		CreatedAt:         now,
		UpdatedAt:         now,
		Attachments:       p.ToMessageAttachments(),
	}
	for _, h := range p.Headers {
		if h.Name == "Message-ID" {
			message.MailMessageId = h.Value // From mail protocol headers
		}
		if h.Name == "In-Reply-To" {
			message.ReplyMailMessageId = &h.Value // From mail protocol headers
		}
	}
	// If this message is a reply to an existing mail message ID
	// check if the stripped text is provided - take that as the text body instead.
	if message.ReplyMailMessageId != nil && p.StrippedTextReply != "" {
		message.TextBody = p.StrippedTextReply
	}
	return message
}

func (p *PostmarkInboundMessageReq) ToMessageAttachments() []models.PostmarkMessageAttachment {
	var attachments []models.PostmarkMessageAttachment
	now := time.Now().UTC()
	for _, m := range p.Attachments {
		attachments = append(attachments, models.PostmarkMessageAttachment{
			CreatedAt:   now,
			UpdatedAt:   now,
			Name:        m.Name,
			ContentType: m.ContentType,
			Content:     m.Content,
		})
	}
	return attachments
}

// FromPostmarkInboundRequest parses an inbound webhook payload from Postmark into a
// PostmarkInboundMessageReq structure.
// It takes a map[string]interface{} request payload and returns the parsed
// PostmarkInboundMessageReq or an error if parsing fails.
func FromPostmarkInboundRequest(
	reqp map[string]interface{}) (PostmarkInboundMessageReq, error) {
	// Convert to JSON bytes
	jsonBytes, err := json.Marshal(reqp)
	if err != nil {
		return PostmarkInboundMessageReq{}, err
	}

	// Parse JSON bytes into struct
	var payload PostmarkInboundMessageReq
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return PostmarkInboundMessageReq{}, err
	}

	payload.Payload = reqp
	return payload, nil
}

// PostmarkOutboxQueue represents a structure for managing message details in a Postmark-based outbox queue system.
// It contains fields for tracking message identifiers, recipient details, error statuses, and timestamps.
type PostmarkOutboxQueue struct {
	MessageId          string    `json:"messageId"`
	PostmarkMessageId  string    `json:"postmarkMessageId"`
	ReplyMailMessageId *string   `json:"replyMailMessageId"`
	HasError           bool      `json:"hasError"`
	MailTo             string    `json:"mailTo"`
	SubmittedAt        time.Time `json:"submittedAt"`
	ErrorCode          int64     `json:"errorCode"`
	Message            string    `json:"message"`
	Acknowledged       bool      `json:"acknowledged"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
