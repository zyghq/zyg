package models

import (
	"database/sql"
	"time"

	"github.com/rs/xid"
)

// ThreadStage represents the lifecycle stage of the Thread.
type ThreadStage struct {
	Stage string
}

func (ts *ThreadStage) String() string {
	return ts.Stage
}

// NeedsFirstResponse is when the Customer made a request on an issue,
// and is waiting for a response from the support team.
func (ts *ThreadStage) NeedsFirstResponse() string {
	return "needs_first_response"
}

// Spam when the conversations seem like spam or suspicious.
// Reduces the noise in the support Thread requests.
func (ts *ThreadStage) Spam() string {
	return "spam"
}

// WaitingOnCustomer when Member has responded and is waiting on a
// response from the Customer.
func (ts *ThreadStage) WaitingOnCustomer() string {
	return "waiting_on_customer"
}

// NeedsNextResponse when the Customer has responded, a Member is yet to respond
// or take next steps.
// There might be some back-and-forth with Customer before moving to next steps.
func (ts *ThreadStage) NeedsNextResponse() string {
	return "needs_next_response"
}

// Hold when resolution is waiting on some dependency.
// Dependencies may include review, external dependencies, or finding the cause.
func (ts *ThreadStage) Hold() string {
	return "hold"
}

// Resolved is the final stage for the Thread. This also marks the ThreadStatus as done.
// It also means customer's concerns were addressed.
// If the Customer again responds, the stage transitions into needs next response.
func (ts *ThreadStage) Resolved() string {
	return "resolved"
}

// ThreadStatus represents the high level status of the Thread.
type ThreadStatus struct {
	Status          string
	StatusChangedAt time.Time
	StatusChangedBy MemberActor
	Stage           ThreadStage
}

// NewThreadStatus creates a new ThreadStatus.
//func (s *ThreadStatus) NewThreadStatus(
//	status string, statusChangedAt time.Time, statusChangedBy MemberActor) ThreadStatus {
//	if !s.IsValid(status) {
//		status = s.DefaultStatus() // fallback to default status.
//	}
//	return ThreadStatus{
//		Status:          status,
//		StatusChangedAt: statusChangedAt,
//		StatusChangedBy: statusChangedBy,
//	}
//}

func (s *ThreadStatus) Todo() string {
	return "todo"
}

func (s *ThreadStatus) Done() string {
	return "done"
}

// Deprecated: DefaultStatus returns default status that can be set to the Thread.
func (s *ThreadStatus) DefaultStatus() string {
	return s.Todo()
}

func (s *ThreadStatus) MarkDone(member MemberActor) {
	s.Status = s.Done()
	s.StatusChangedAt = time.Now()
	s.StatusChangedBy = member
}

// IsValid checks if the given status is valid.
// Returns true if valid otherwise false.
func (s *ThreadStatus) IsValid(status string) bool {
	switch status {
	case s.Done():
		return true
	case s.Todo():
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

// ThreadCustomer represents attached Thread Customer.
type ThreadCustomer struct {
	CustomerId string
	Name       string
}

// A Thread represents conversation with a Customer on a specific issue or topic.
type Thread struct {
	WorkspaceId     string           // WorkspaceId is the ID of the Workspace this Thread belongs to.
	ThreadId        string           // ThreadId represents the unique ID of the Thread.
	Customer        ThreadCustomer   // The attached Customer.
	CustomerId      string           // Deprecated: CustomerId represents the ID of the Customer this Thread is part of.
	CustomerName    string           // Deprecated: CustomerName represents the name of the Customer as per CustomerId.
	AssigneeId      sql.NullString   // Deprecated: AssignedId represents the Member this Thread is assigned to.
	AssigneeName    sql.NullString   // Deprecated: AssigneeName represents the name of the Member as per AssigneeId.
	AssignedMember  *AssignedMember  // The Member assigned to the Thread.
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
	CreatedAt       time.Time        // When the Thread was created
	UpdatedAt       time.Time        // When the Thread was last updated.
}

func (th *Thread) GenId() string {
	return "th" + xid.New().String()
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

// AddInboundMessage adds the inbound message info to the Thread.
// Inbound messages are messages from the Customer.
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

func (c Chat) PreviewText() string {
	if len(c.Body) > 255 {
		return c.Body[:255]
	}
	return c.Body
}
