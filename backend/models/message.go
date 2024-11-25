package models

import (
	"time"

	"github.com/rs/xid"
)

// Message represents multi-channel Thread message
type Message struct {
	MessageId string
	ThreadId  string
	TextBody  string
	Body      string
	Customer  *CustomerActor
	Member    *MemberActor
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

func SetMessageBody(body string) MessageOption {
	return func(message *Message) {
		message.Body = body
	}
}

type MessageAttachment struct {
	AttachmentId string    `json:"attachmentId"`
	Name         string    `json:"name"`
	ContentType  string    `json:"contentType"`
	ContentKey   string    `json:"contentKey"`
	ContentUrl   string    `json:"contentUrl"`
	Spam         bool      `json:"spam"`
	HasError     bool      `json:"hasError"`
	Error        string    `json:"error"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// PostmarkInboundMessage represents 1:1 mapping with Message
// Attributes are specific to Postmark.
type PostmarkInboundMessage struct {
	// Raw inbound request payload
	Payload map[string]interface{}
	// `MessageID` from Postmark
	PostmarkMessageId string

	// From mail protocol `Message-ID` in headers
	MailMessageId string
	// From mail protocol `In-Reply-To` To `Message-ID` in headers
	ReplyMailMessageId *string

	Subject  string
	TextBody string
	HTMLBody string

	FromEmail string
	FromName  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

//func (p *PostmarkInboundMessage) GenId() string {
//	return "msg" + xid.New().String()
//}

// FromPayload initializes a PostmarkInboundMessage instance from the given payload map.
//func (p *PostmarkInboundMessage) FromPayload(
//	payload map[string]interface{}) (*PostmarkInboundMessage, error) {
//	if payload == nil {
//		return nil, fmt.Errorf("payload cannot be nil")
//	}
//
//	// current time space
//	now := time.Now().UTC()
//
//	// store the raw payload as is.
//	p.Payload = payload
//	p.CreatedAt = now
//	p.UpdatedAt = now // initially
//
//	// Message ID from Postmark, not to be confused with Message-ID in Headers
//	// yet could be mixed or formatted version of the same but treat it as different.
//	if msgID, ok := payload["MessageID"].(string); ok {
//		p.PostmarkMessageId = msgID
//	} else {
//		return nil, fmt.Errorf("postmark MessageID not found in payload refer postmark docs")
//	}
//
//	// Extract mail message IDs from headers array
//	if headers, ok := payload["Headers"].([]interface{}); ok {
//		for _, header := range headers {
//			if headerMap, ok := header.(map[string]interface{}); ok {
//				name, hasName := headerMap["Name"].(string)
//				value, hasValue := headerMap["Value"].(string)
//
//				if !hasName || !hasValue {
//					continue
//				}
//
//				switch name {
//				case "Message-ID":
//					p.MailMessageId = value
//				case "In-Reply-To":
//					p.ReplyMailMessageId = &value
//				}
//			}
//		}
//	}
//
//	// If mail Message ID is empty use the system message ID as mail message ID prefixed by `zyg:`
//	if p.MailMessageId == "" {
//		p.MailMessageId = fmt.Sprintf("zyg:%s", p.PostmarkMessageId)
//	}
//
//	if Subject, ok := payload["Subject"].(string); ok {
//		p.Subject = Subject
//	}
//
//	if TextBody, ok := payload["TextBody"].(string); ok {
//		p.TextBody = TextBody
//	}
//
//	if HTMLBody, ok := payload["HtmlBody"].(string); ok {
//		p.HTMLBody = HTMLBody
//	}
//
//	// get From details.
//	if fromFull, ok := payload["FromFull"].(map[string]interface{}); ok {
//		FromEmail, hasEmail := fromFull["Email"].(string)
//		if hasEmail {
//			p.FromEmail = FromEmail
//		}
//
//		FromName, hasName := fromFull["Name"].(string)
//		if hasName {
//			p.FromName = FromName
//		}
//	}
//	return p, nil
//}

//func (p *PostmarkInboundMessage) Subject() string {
//	return p.Subject
//}

//func (p *PostmarkInboundMessage) PlainText() string {
//	return p.TextBody
//}

//func (p *PostmarkInboundMessage) HTML() string {
//	return p.HTMLBody
//}

//func (p *PostmarkInboundMessage) FromEmail() string {
//	return p.FromEmail
//}
//
//func (p *PostmarkInboundMessage) FromName() string {
//	return p.FromName
//}

// ThreadMessage combines a Thread and its associated Message.
type ThreadMessage struct {
	Thread  *Thread
	Message *Message
}

// PostmarkInboundThreadMessage combines a ThreadMessage and a PostmarkInboundMessage.
type PostmarkInboundThreadMessage struct {
	ThreadMessage
	PostmarkInboundMessage *PostmarkInboundMessage
}
