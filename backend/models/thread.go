package models

import (
	"database/sql"
	"time"

	"github.com/rs/xid"
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

// ThreadStatus represents the high level status of the Thread.
type ThreadStatus struct {
	Status          string
	StatusChangedAt time.Time
	StatusChangedBy MemberActor
	Stage           string
}

// NewStatus creates the new ThreadStatus with initial defaults.
func (ts *ThreadStatus) NewStatus(member MemberActor) ThreadStatus {
	t := ThreadStatus{}
	t.NeedsFirstResponse(member)
	return t
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

// WaitingOnCustomer when Member has responded and is waiting for a
// response from the Customer.
func (ts *ThreadStatus) WaitingOnCustomer(member MemberActor) {
	ts.Stage = waitingOnCustomer // update the stage
	ts.Status = ts.Todo()        // status for this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// Hold when resolution is waiting on some dependency.
// Dependencies may include review, external dependencies, or finding the cause.
func (ts *ThreadStatus) Hold(member MemberActor) {
	ts.Stage = hold       // update the stage
	ts.Status = ts.Todo() // status for this stage
	ts.StatusChangedAt = time.Now().UTC()
	ts.StatusChangedBy = member
}

// NeedsNextResponse when the Customer has responded, a Member is yet to respond
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
	return "todo"
}

func (ts *ThreadStatus) Done() string {
	return "done"
}

// Deprecated: DefaultStatus returns default status that can be set to the Thread.
func (ts *ThreadStatus) DefaultStatus() string {
	return ts.Todo()
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
// This also represents the communication channel for the Customer.
type ThreadChannel struct{}

func (c ThreadChannel) Chat() string {
	return "chat"
}

// InboundMessage tracks the inbound message received from the Customer.
// Holds required attributes to track the inbound message.
// Common across channels.
type InboundMessage struct {
	MessageId    string
	CustomerId   string
	CustomerName string
	PreviewText  string
	FirstSeqId   string
	LastSeqId    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (im InboundMessage) GenId() string {
	return "im" + xid.New().String()
}

// OutboundMessage tracks the outbound message sent by the Member.
// Holds required attributes to track an outbound message.
// Common across channels.
type OutboundMessage struct {
	MessageId   string
	MemberId    string
	MemberName  string
	PreviewText string
	FirstSeqId  string
	LastSeqId   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (em OutboundMessage) GenId() string {
	return "em" + xid.New().String()
}

// A Thread represents conversation with a Customer on a specific issue or topic.
type Thread struct {
	WorkspaceId     string           // WorkspaceId is the ID of the Workspace this Thread belongs to.
	ThreadId        string           // ThreadId represents the unique ID of the Thread.
	Customer        CustomerActor    // The attached Customer.
	CustomerId      string           // Deprecated: CustomerId represents the ID of the Customer this Thread is part of.
	CustomerName    string           // Deprecated: CustomerName represents the name of the Customer as per CustomerId.
	AssignedMember  *AssignedMember  // The Member assigned to the Thread.
	AssigneeId      sql.NullString   // Deprecated: AssignedId represents the Member this Thread is assigned to.
	AssigneeName    sql.NullString   // Deprecated: AssigneeName represents the name of the Member as per AssigneeId.
	Title           string           // The Title of the Thread, which allows to quickly identify what it is about.
	Description     string           // The Description of the Thread could be descriptive.
	Sequence        int              // Deprecated: will be removed in the next release.
	Status          string           // Deprecated: use ThreadStatus instead.
	ThreadStatus    ThreadStatus     // The status of the Thread. TODO: rename to `Status` post removal.
	Read            bool             // Deprecated: use upcoming stages instead.
	Replied         bool             // If the Member has anytime replied to the Thread.
	Priority        string           // The Priority of the Thread as per ThreadPriority.
	Spam            bool             // Deprecated: will be removed in the next release.
	Channel         string           // The channel this Thread belongs to as per ThreadChannel.
	InboundMessage  *InboundMessage  // InboundMessage tracks the inbound message from Customer
	OutboundMessage *OutboundMessage // OutboundMessage tracks the outbound message from Member
	CreatedBy       MemberActor      // The Member who created this Thread.
	UpdatedBy       MemberActor      // The Member who updated this Thread.
	CreatedAt       time.Time        // When the Thread was created
	UpdatedAt       time.Time        // When the Thread was last updated.
}

func (th *Thread) GenId() string {
	return "th" + xid.New().String()
}

// CreateNewThread creates a new thread with the provided attributes and defaults.
func (th *Thread) CreateNewThread(
	workspaceId string, customer CustomerActor,
	createdBy MemberActor, updatedBy MemberActor, channel string) Thread {
	thread := Thread{}
	status := ThreadStatus{}
	thread.WorkspaceId = workspaceId
	thread.Customer = customer
	thread.CreatedBy = createdBy
	thread.UpdatedBy = updatedBy
	thread.Channel = channel
	// Defaults
	thread.ThreadId = thread.GenId()
	thread.SetDefaultTitle()
	thread.ThreadStatus = status.NewStatus(createdBy)
	thread.Replied = false
	thread.Priority = ThreadPriority{}.DefaultPriority()
	return thread
}

// AssignMember assigns the member to the thread and when the assignment was made.
// TODO: pass MemberActor instead of memberId and name.
func (th *Thread) AssignMember(memberId string, name string, assignedAt time.Time) {
	th.AssignedMember = &AssignedMember{
		MemberId:   memberId,
		Name:       name,
		AssignedAt: assignedAt,
	}
}

func (th *Thread) ClearAssignedMember() {
	th.AssignedMember = nil
}

// AddInboundMessage adds the inbound message info to the Thread.
// Inbound messages are messages from the Customer.
// TODO: use CustomerActor instead of passing customerId and customerName.
func (th *Thread) AddInboundMessage(messageId string, customerId string, customerName string,
	previewText string,
	firstSeqId string, lastSeqId string,
	createdAt time.Time, updatedAt time.Time,
) {
	th.InboundMessage = &InboundMessage{
		MessageId: messageId, CustomerId: customerId, CustomerName: customerName,
		PreviewText: previewText,
		FirstSeqId:  firstSeqId, LastSeqId: lastSeqId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (th *Thread) ClearInboundMessage() {
	th.InboundMessage = nil
}

// AddOutboundMessage adds the outbound message info to the Thread.
// Outbound messages are messages from the Member.
// TODO: use MemberActor instead of passing memberId and memberName.
func (th *Thread) AddOutboundMessage(messageId string, memberId string, memberName string,
	previewText string,
	firstSeqId string, lastSeqId string,
	createdAt time.Time, updatedAt time.Time,
) {
	th.OutboundMessage = &OutboundMessage{
		MessageId: messageId, MemberId: memberId, MemberName: memberName,
		PreviewText: previewText,
		FirstSeqId:  firstSeqId, LastSeqId: lastSeqId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
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

// Chat
// TODO:
//   - use CustomerActor instead of passing customerId and customerName.
//   - use MemberActor instead of passing memberId and memberName.
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
	threadId string, customerId string, isHead bool, message string) Chat {
	chat := Chat{
		ThreadId:   threadId,
		CustomerId: NullString(&customerId),
		IsHead:     isHead,
		Body:       message,
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
