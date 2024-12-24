package models

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

// Message represents multi-channel Thread message
type Message struct {
	MessageId    string
	ThreadId     string
	TextBody     string
	MarkdownBody string
	HTMLBody     string
	Customer     *CustomerActor
	Member       *MemberActor
	// Deprecated
	Channel   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MessageOption func(message *Message)

func (m *Message) GenId() string {
	return "msg" + xid.New().String()
}

func (m *Message) PreviewText() string {
	if len(m.TextBody) > 255 {
		return m.TextBody[:255]
	}
	return m.TextBody
}

func NewMessage(threadId string, channel string, opts ...MessageOption) *Message {
	messageId := (&Message{}).GenId()
	now := time.Now().UTC()
	message := &Message{
		MessageId: messageId,
		ThreadId:  threadId,
		Channel:   channel,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, opt := range opts {
		opt(message)
	}
	return message
}

func SetMessageCustomer(customer CustomerActor) MessageOption {
	return func(message *Message) {
		message.Customer = &customer
		message.Member = nil // it's either the Customer or the Member
	}
}

func SetMessageMember(member MemberActor) MessageOption {
	return func(message *Message) {
		message.Member = &member
		message.Customer = nil // it's either the Member or the Customer
	}
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

// MessageAttachment represents metadata and identification details for a file attachment linked to a message.
type MessageAttachment struct {
	AttachmentId string    `json:"attachmentId"`
	MessageId    string    `json:"messageId"`
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

func (m *MessageAttachment) GenId() string {
	return "at" + xid.New().String()
}

type MessageWithAttachments struct {
	Message
	Attachments []MessageAttachment
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

func (p *PostmarkInboundMessage) NewPostmarkInboundLog(messageId string) PostmarkMessageLog {
	now := time.Now().UTC()
	return PostmarkMessageLog{
		MessageId:          messageId,
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

// ThreadMessage combines a Thread and its associated Message.
type ThreadMessage struct {
	Thread  *Thread
	Message *Message
}

// PostmarkInboundThreadMessage combines a ThreadMessage and a PostmarkInboundMessage.
//type PostmarkInboundThreadMessage struct {
//	ThreadMessage
//	PostmarkInboundMessage *PostmarkInboundMessage
//}

type PostmarkMessageLog struct {
	MessageId          string
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
	m.MailMessageId = fmt.Sprintf("<%s@%s>", m.MessageId, d)
}
