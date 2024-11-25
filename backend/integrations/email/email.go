package email

import (
	"encoding/json"
	"github.com/zyghq/postmark"
	"github.com/zyghq/zyg/models"
	"time"
)

type PostmarkInboundMessageReq struct {
	postmark.InboundMessageDetail
	Payload map[string]interface{}
}

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
