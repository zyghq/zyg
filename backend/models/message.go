package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Message struct {
	TextBody     string    `json:"textBody"`
	MarkdownBody string    `json:"markdownBody"`
	HTMLBody     string    `json:"htmlBody"`
	Channel      string    `json:"channel"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (m *Message) ToJSON() map[string]interface{} {
	result := make(map[string]interface{})
	data, err := json.Marshal(m)
	if err != nil {
		return result
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result
	}

	return result
}

type MessageOption func(message *Message)

func NewMessage(channel string, opts ...MessageOption) *Message {
	now := time.Now().UTC()
	message := &Message{
		Channel:   channel,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, opt := range opts {
		opt(message)
	}
	return message
}

func SetMessageTextBody(textBody string) MessageOption {
	return func(message *Message) {
		message.TextBody = textBody
	}
}

func SetMarkdownBody(body string) MessageOption {
	return func(message *Message) {
		message.MarkdownBody = body
	}
}

func SetHTMLBody(body string) MessageOption {
	return func(message *Message) {
		message.HTMLBody = body
	}
}

type ActivityAttachment struct {
	AttachmentId string    `json:"attachmentId"`
	ActivityID   string    `json:"activityId"`
	Name         string    `json:"name"`
	ContentType  string    `json:"contentType"`
	ContentKey   string    `json:"contentKey"`
	ContentUrl   string    `json:"contentUrl"`
	Spam         bool      `json:"spam"`
	HasError     bool      `json:"hasError"`
	Error        string    `json:"error"`
	MD5Hash      string    `json:"md5Hash"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (m *ActivityAttachment) GenId() string {
	return "at" + xid.New().String()
}

type ActivityWithAttachments struct {
	Activity
	Attachments []ActivityAttachment
}

type PostmarkMessageAttachment struct {
	Name        string
	ContentType string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PostmarkInboundMessage represents abstracted Postmark inbound webhook request.
// Attributes are specific to Postmark and as required by the application.
type PostmarkInboundMessage struct {
	// Raw inbound request payload
	Payload map[string]interface{}
	// `MessageID` from Postmark
	PostmarkMessageId string

	// From mail protocol `Message-ID` in headers
	MailMessageId string
	// From mail protocol `In-Reply-To` To `Message-ID` in headers
	ReplyMailMessageId *string

	FromEmail string
	FromName  string

	CreatedAt time.Time
	UpdatedAt time.Time

	Subject  string
	TextBody string
	HTMLBody string

	Attachments []PostmarkMessageAttachment
}

func (p *PostmarkInboundMessage) ToPostmarkMessageLog(activityID string) PostmarkMessageLog {
	now := time.Now().UTC()
	return PostmarkMessageLog{
		ActivityID:         activityID,
		Payload:            p.Payload,
		PostmarkMessageId:  p.PostmarkMessageId,
		MailMessageId:      p.MailMessageId,
		ReplyMailMessageId: p.ReplyMailMessageId,
		HasError:           false,
		SubmittedAt:        now,
		MessageType:        "inbound",
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

type PostmarkMessageLog struct {
	ActivityID         string
	Payload            map[string]interface{}
	PostmarkMessageId  string
	MailMessageId      string
	ReplyMailMessageId *string
	HasError           bool
	SubmittedAt        time.Time
	ErrorCode          int64
	PostmarkMessage    string
	MessageEvent       string
	Acknowledged       bool // has Postmark outbound message acknowledged with details API
	MessageType        string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SetOutboundMailMessageId sets outbound mail message ID with specified domain as configured.
// This should be only used for outbound mails as inbound mail already has it generated from the client.
func (m *PostmarkMessageLog) SetOutboundMailMessageId(d string) {
	m.MailMessageId = fmt.Sprintf("<%s@%s>", m.PostmarkMessageId, d)
}
