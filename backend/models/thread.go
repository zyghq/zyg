package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/xid"
)

// Represents Thread status.
const (
	todo = "todo"
	done = "done"
)

// Represents Thread stages.
const (
	spam               = "spam"
	needsFirstResponse = "needs_first_response"
	waitingOnCustomer  = "waiting_on_customer"
	hold               = "hold"
	needsNextResponse  = "needs_next_response"
	resolved           = "resolved"
)

// Represents Thread communication sources.
const (
	inAppChat = "in_app_chat"
	email     = "email"
)

// ThreadStatus represents the high level status of the Thread.
type ThreadStatus struct {
	Status          string
	StatusChangedAt time.Time
	StatusChangedBy MemberActor
	Stage           string
}

// InitialStatus creates the new ThreadStatus with initial defaults.
func (ts *ThreadStatus) InitialStatus(member MemberActor) {
	ts.NeedsFirstResponse(member)
}

// Spam when the conversations seem like spam or suspicious.
// Reduces the noise in the support Thread requests.
func (ts *ThreadStatus) Spam(member MemberActor) {
	ts.Stage = spam
	ts.Status = ts.Todo()
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// NeedsFirstResponse is when the Customer made a request on an issue,
// and is waiting for a response from the support team.
func (ts *ThreadStatus) NeedsFirstResponse(member MemberActor) {
	ts.Stage = needsFirstResponse // update the stage
	ts.Status = ts.Todo()         // status for this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// WaitingOnCustomer is when Member has responded and is waiting for a
// response from the Customer.
func (ts *ThreadStatus) WaitingOnCustomer(member MemberActor) {
	ts.Stage = waitingOnCustomer // update the stage
	ts.Status = ts.Todo()        // status for this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// Hold is when resolution is waiting on some dependency.
// Dependencies may include review, external dependencies, or finding the cause.
func (ts *ThreadStatus) Hold(member MemberActor) {
	ts.Stage = hold       // update the stage
	ts.Status = ts.Todo() // status for this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// NeedsNextResponse is when the Customer has responded and a Member is yet to respond
// or take next steps.
// There might be some back-and-forth with Customer before moving to next steps.
func (ts *ThreadStatus) NeedsNextResponse(member MemberActor) {
	ts.Stage = needsNextResponse // update the stage
	ts.Status = ts.Todo()        // status for this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// Resolved is the final stage for the Thread. This also marks the ThreadStatus as done.
// It also means customer's concerns were addressed.
// If the Customer again responds, the stage transitions into needs next response.
func (ts *ThreadStatus) Resolved(member MemberActor) {
	ts.Stage = resolved   // update the stage
	ts.Status = ts.Done() // status of this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

func (ts *ThreadStatus) Todo() string {
	return todo
}

func (ts *ThreadStatus) Done() string {
	return done
}

func (ts *ThreadStatus) MarkDone(member MemberActor) {
	ts.Status = ts.Done()
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// IsValidStage checks if the given stage is valid.
// Returns true if valid otherwise false.
func (ts *ThreadStatus) IsValidStage(stage string) bool {
	switch stage {
	case spam, needsFirstResponse, waitingOnCustomer, hold, needsNextResponse, resolved:
		return true
	default:
		return false
	}
}

// ThreadPriority represents the priority of the Thread.
// Invoke methods for specific priority.
type ThreadPriority struct{}

func (p ThreadPriority) Urgent() string {
	return "urgent"
}

func (p ThreadPriority) High() string {
	return "high"
}

func (p ThreadPriority) Normal() string {
	return "normal"
}

func (p ThreadPriority) Low() string {
	return "low"
}

// DefaultPriority returns the default priority that can be set for the Thread.
func (p ThreadPriority) DefaultPriority() string {
	return p.Normal()
}

// IsValid checks if the given priority is valid
// Returns true if valid otherwise false.
func (p ThreadPriority) IsValid(s string) bool {
	switch s {
	case p.Urgent(), p.High(), p.Normal(), p.Low():
		return true
	default:
		return false
	}
}

// ThreadChannel represents the source channel of the Thread.
type ThreadChannel struct{}

func (c ThreadChannel) InAppChat() string {
	return inAppChat
}

func (c ThreadChannel) Email() string {
	return email
}

// InboundMessage tracks the inbound message received from the Customer.
// Common across channels.
type InboundMessage struct {
	MessageId   string
	Customer    CustomerActor
	PreviewText string
	FirstSeqId  string
	LastSeqId   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (im InboundMessage) GenId() string {
	return "im" + xid.New().String()
}

// OutboundMessage tracks the outbound message sent by the Member.
// Common across channels.
type OutboundMessage struct {
	MessageId   string
	Member      MemberActor
	PreviewText string
	FirstSeqId  string
	LastSeqId   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (om OutboundMessage) GenId() string {
	return "om" + xid.New().String()
}

// A Thread represents conversation with a Customer on a specific issue or topic.
type Thread struct {
	ThreadId        string           // ThreadId represents the unique ID of the Thread.
	WorkspaceId     string           // WorkspaceId is the ID of the Workspace this Thread belongs to.
	Customer        CustomerActor    // The attached Customer.
	AssignedMember  *AssignedMember  // The Member assigned to the Thread.
	Title           string           // The Title of the Thread, which allows to quickly identify what it is about.
	Description     string           // The Description of the Thread could be descriptive.
	ThreadStatus    ThreadStatus     // The status of the Thread.
	Replied         bool             // If the Member has anytime replied to the Thread.
	Priority        string           // The Priority of the Thread as per ThreadPriority.
	Channel         string           // The source channel this Thread belongs to as per ThreadChannel.
	InboundMessage  *InboundMessage  // InboundMessage tracks the inbound message from Customer
	OutboundMessage *OutboundMessage // OutboundMessage tracks the outbound message from Member
	CreatedBy       MemberActor      // The Member who created this Thread.
	UpdatedBy       MemberActor      // The Member who updated this Thread.
	CreatedAt       time.Time        // When the Thread was created
	UpdatedAt       time.Time        // When the Thread was last updated.
}

type ThreadOption func(*Thread)

func (th *Thread) GenId() string {
	return "th" + xid.New().String()
}

func NewThread(
	workspaceId string, customer CustomerActor, createdBy MemberActor, channel string,
	opts ...ThreadOption,
) *Thread {
	threadId := (&Thread{}).GenId()
	now := time.Now().UTC()

	status := ThreadStatus{}
	status.InitialStatus(createdBy)

	thread := &Thread{
		ThreadId:     threadId,
		WorkspaceId:  workspaceId,
		Customer:     customer,
		Channel:      channel,
		ThreadStatus: status,
		Replied:      false,
		Priority:     ThreadPriority{}.DefaultPriority(),
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
	}
	thread.SetDefaultTitle()

	for _, opt := range opts {
		opt(thread)
	}
	return thread
}

func SetThreadTitle(title string) ThreadOption {
	maxLength := 255
	if len(title) > maxLength {
		title = title[:maxLength]
	}
	return func(thread *Thread) {
		thread.Title = title
	}
}

func SetThreadDescription(description string) ThreadOption {
	maxLength := 512
	if len(description) > maxLength {
		description = description[:maxLength]
	}
	return func(thread *Thread) {
		thread.Description = description
	}
}

// AssignMember assigns the member to the thread and when the assignment was made.
func (th *Thread) AssignMember(member MemberActor, assignedAt time.Time) {
	th.AssignedMember = &AssignedMember{
		MemberId:   member.MemberId,
		Name:       member.Name,
		AssignedAt: assignedAt,
	}
}

func (th *Thread) ClearAssignedMember() {
	th.AssignedMember = nil
}

// SetNewInboundMessage adds the inbound message info to the Thread.
// Inbound messages are messages from the Customer.
func (th *Thread) SetNewInboundMessage(customer CustomerActor, previewText string) {
	messageId := InboundMessage{}.GenId()
	seqId := xid.New().String()
	now := time.Now().UTC()
	th.InboundMessage = &InboundMessage{
		MessageId:   messageId,
		Customer:    customer,
		PreviewText: previewText,
		FirstSeqId:  seqId,
		LastSeqId:   seqId, // starts with first seq.
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (th *Thread) SetNextInboundSeq(previewText string) {
	seqId := xid.New().String()
	now := time.Now().UTC()
	if th.InboundMessage != nil {
		th.InboundMessage.PreviewText = previewText
		th.InboundMessage.LastSeqId = seqId
		th.InboundMessage.UpdatedAt = now
	} else {
		th.SetNewInboundMessage(th.Customer, previewText)
	}
}

func (th *Thread) ClearInboundMessage() {
	th.InboundMessage = nil
}

// SetNewOutboundMessage adds the outbound message info to the Thread.
// Outbound messages are messages from the Member.
func (th *Thread) SetNewOutboundMessage(
	messageId string, member MemberActor, previewText string,
) {
	seqId := xid.New().String()
	now := time.Now().UTC()
	th.OutboundMessage = &OutboundMessage{
		MessageId:   messageId,
		Member:      member,
		PreviewText: previewText,
		FirstSeqId:  seqId,
		LastSeqId:   seqId, // starts with first seq.
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (th *Thread) SetNextOutboundSeq(member MemberActor, previewText string) {
	seqId := xid.New().String()
	now := time.Now().UTC()
	if th.OutboundMessage != nil {
		th.OutboundMessage.PreviewText = previewText
		th.OutboundMessage.LastSeqId = seqId
		th.OutboundMessage.UpdatedAt = now
	} else {
		messageId := OutboundMessage{}.GenId()
		th.SetNewOutboundMessage(messageId, member, previewText)
	}
}

func (th *Thread) ClearOutboundMessage() {
	th.OutboundMessage = nil
}

func (th *Thread) SetDefaultTitle() {
	th.Title = "Support Request"
}

// PreviewText
// TODO:
//   - possibly show the latest or
//   - have some kind logic based on upcoming thread stages.
func (th *Thread) PreviewText() string {
	if th.InboundMessage != nil {
		return th.InboundMessage.PreviewText
	}
	if th.OutboundMessage != nil {
		return th.OutboundMessage.PreviewText
	}
	return ""
}

// CustomerPreviewText
// TODO:
//   - update based on PreviewText.
func (th *Thread) CustomerPreviewText() string {
	if th.OutboundMessage != nil {
		return th.OutboundMessage.PreviewText
	}
	if th.InboundMessage != nil {
		return th.InboundMessage.PreviewText
	}
	return ""
}

// SetDefaultStatus checks if the Thread has been already been replied,
// If not then it sets the default status as NeedsFirstResponse.
// else it sets the default status as NeedsNextResponse.
func (th *Thread) SetDefaultStatus(member MemberActor) {
	if th.Replied {
		th.ThreadStatus.NeedsNextResponse(member)
	} else {
		th.ThreadStatus.NeedsFirstResponse(member)
	}
}

func (th *Thread) SetStatusStage(stage string, member MemberActor) {
	switch stage {
	case spam:
		th.ThreadStatus.Spam(member)
	case needsFirstResponse:
		th.ThreadStatus.NeedsFirstResponse(member)
	case waitingOnCustomer:
		th.ThreadStatus.WaitingOnCustomer(member)
	case hold:
		th.ThreadStatus.Hold(member)
	case needsNextResponse:
		th.ThreadStatus.NeedsNextResponse(member)
	case resolved:
		th.ThreadStatus.Resolved(member)
	default:
		th.SetDefaultStatus(member)
	}
}

// Chat represents a chat message in a thread.
// Deprecated: This struct is deprecated migrate to Message instead.
type Chat struct {
	ThreadId     string
	ChatId       string
	Body         string
	Sequence     int
	CustomerId   sql.NullString
	CustomerName sql.NullString
	MemberId     sql.NullString
	MemberName   sql.NullString
	IsHead       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c Chat) GenId() string {
	return "ch" + xid.New().String()
}

func (c Chat) CreateNewCustomerChat(
	threadId string, customerId string, isHead bool, message string,
	createdAt time.Time, updatedAt time.Time) Chat {
	chat := Chat{
		ThreadId:   threadId,
		CustomerId: NullString(&customerId),
		IsHead:     isHead,
		Body:       message,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
	// Defaults
	chat.ChatId = chat.GenId()
	return chat
}

func (c Chat) PreviewText() string {
	if len(c.Body) > 255 {
		return c.Body[:255]
	}
	return c.Body
}

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

// PostmarkInboundMessage represents 1:1 mapping with Message
// Attributes are specific to Postmark.
type PostmarkInboundMessage struct {
	// raw http request payload.
	Payload map[string]interface{}
	// MessageID generated by Postmark
	PMMessageId string
	// from mail protocol mail message-ID from headers
	MailMessageId string
	// from mail protocol In Reply To mail message-ID
	ReplyMailMessageId *string
	// Datetime from Postmark/Email `Date` or system.
	CreatedAt time.Time
	UpdatedAt time.Time

	subject  string
	textBody string
	htmlBody string

	fromEmail string
	fromName  string
}

func (p *PostmarkInboundMessage) GenId() string {
	return "msg" + xid.New().String()
}

// FromPayload initializes a PostmarkInboundMessage instance from the given payload map.
func (p *PostmarkInboundMessage) FromPayload(
	payload map[string]interface{}) (*PostmarkInboundMessage, error) {
	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	// current time space
	now := time.Now().UTC()

	// store the raw payload as is.
	p.Payload = payload
	p.CreatedAt = now
	p.UpdatedAt = now // initially

	// Message ID from Postmark, not to be confused with Message-ID in Headers
	// yet could be mixed or formatted version of the same but treat it as different.
	if msgID, ok := payload["MessageID"].(string); ok {
		p.PMMessageId = msgID
	} else {
		return nil, fmt.Errorf("postmark MessageID not found in payload refer postmark docs")
	}

	// Extract mail message IDs from headers array
	if headers, ok := payload["Headers"].([]interface{}); ok {
		for _, header := range headers {
			if headerMap, ok := header.(map[string]interface{}); ok {
				name, hasName := headerMap["Name"].(string)
				value, hasValue := headerMap["Value"].(string)

				if !hasName || !hasValue {
					continue
				}

				switch name {
				case "Message-ID":
					p.MailMessageId = value
				case "In-Reply-To":
					p.ReplyMailMessageId = &value
				}
			}
		}
	}

	// If mail Message ID is empty use the system message ID as mail message ID prefixed by `zyg:`
	if p.MailMessageId == "" {
		p.MailMessageId = fmt.Sprintf("zyg:%s", p.PMMessageId)
	}

	if subject, ok := payload["Subject"].(string); ok {
		p.subject = subject
	}

	if textBody, ok := payload["TextBody"].(string); ok {
		p.textBody = textBody
	}

	if htmlBody, ok := payload["HtmlBody"].(string); ok {
		p.htmlBody = htmlBody
	}

	// get From details.
	if fromFull, ok := payload["FromFull"].(map[string]interface{}); ok {
		fromEmail, hasEmail := fromFull["Email"].(string)
		if hasEmail {
			p.fromEmail = fromEmail
		}

		fromName, hasName := fromFull["Name"].(string)
		if hasName {
			p.fromName = fromName
		}
	}
	return p, nil
}

func (p *PostmarkInboundMessage) Subject() string {
	return p.subject
}

func (p *PostmarkInboundMessage) PlainText() string {
	return p.textBody
}

func (p *PostmarkInboundMessage) Html() string {
	return p.htmlBody
}

func (p *PostmarkInboundMessage) FromEmail() string {
	return p.fromEmail
}

func (p *PostmarkInboundMessage) FromName() string {
	return p.fromName
}

// ThreadMessage combines a Thread and its associated Message.
type ThreadMessage struct {
	Thread  *Thread
	Message *Message
}

// ThreadMessageWithPostmarkInbound combines a ThreadMessage and a PostmarkInboundMessage.
type ThreadMessageWithPostmarkInbound struct {
	ThreadMessage
	PostmarkInboundMessage *PostmarkInboundMessage
}
