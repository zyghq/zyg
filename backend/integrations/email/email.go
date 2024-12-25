package email

import (
	"context"
	"encoding/json"
	"github.com/zyghq/postmark"
	"github.com/zyghq/zyg/integrations"
	"github.com/zyghq/zyg/models"
	"log/slog"
	"time"
)

// PostmarkInboundMessageReq represents inbound webhook request from Postmark
// with the raw JSON payload.
type PostmarkInboundMessageReq struct {
	postmark.InboundWebhook
	Payload map[string]interface{}
}

// ToPostmarkInboundMessage converts a PostmarkInboundMessageReq to a PostmarkInboundMessage model instance.
func (p *PostmarkInboundMessageReq) ToPostmarkInboundMessage() models.PostmarkInboundMessage {
	now := time.Now().UTC()
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
		Attachments:       p.ToPostmarkMessageAttachments(),
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

func (p *PostmarkInboundMessageReq) ToPostmarkMessageAttachments() []models.PostmarkMessageAttachment {
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
func FromPostmarkInboundRequest(reqp map[string]interface{}) (PostmarkInboundMessageReq, error) {
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

type PostmarkEmailReqOption func(req *postmark.Email)

func NewPostmarkEmailReq(subject, from, to string, opts ...PostmarkEmailReqOption) *postmark.Email {
	req := &postmark.Email{
		Subject: subject,
		From:    from,
		To:      to,
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

func SetPostmarkTextBody(textBody string) PostmarkEmailReqOption {
	return func(req *postmark.Email) {
		req.TextBody = textBody
	}
}

func SetPostmarkHTMLBody(htmlBody string) PostmarkEmailReqOption {
	return func(req *postmark.Email) {
		req.HTMLBody = htmlBody
	}
}

func WithPostmarkHeader(name, value string) PostmarkEmailReqOption {
	return func(req *postmark.Email) {
		req.Headers = append(req.Headers, postmark.Header{
			Name:  name,
			Value: value,
		})
	}
}

func SetPostmarkTag(tag string) PostmarkEmailReqOption {
	return func(req *postmark.Email) {
		req.Tag = tag
	}
}

func SendPostmarkMail(
	ctx context.Context, setting models.PostmarkMailServerSetting, email *postmark.Email,
) (postmark.EmailResponse, error) {
	client := postmark.NewClient(setting.ServerToken, "")
	r, err := client.SendEmail(ctx, *email)
	if err != nil {
		slog.Error("failed to send email", slog.Any("error", err), slog.Any("email", email))
		return postmark.EmailResponse{}, integrations.ErrPostmarkSendMail
	}
	return r, nil
}
