package handler

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/zyghq/zyg/models"
)

type PATReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WorkspaceReq struct {
	Name string `json:"name"`
}

type WorkspaceResp struct {
	WorkspaceId string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (w WorkspaceResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		Name        string `json:"name"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: w.WorkspaceId,
		Name:        w.Name,
		CreatedAt:   w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   w.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type CustomerResp struct {
	CustomerId      string
	ExternalId      sql.NullString
	Email           sql.NullString
	Phone           sql.NullString
	Name            string
	AvatarUrl       string
	IsEmailVerified bool
	Role            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c CustomerResp) MarshalJSON() ([]byte, error) {
	var externalId, email, phone *string
	if c.ExternalId.Valid {
		externalId = &c.ExternalId.String
	}
	if c.Email.Valid {
		email = &c.Email.String
	}
	if c.Phone.Valid {
		phone = &c.Phone.String
	}

	aux := &struct {
		CustomerId      string  `json:"customerId"`
		ExternalId      *string `json:"externalId"`
		Email           *string `json:"email"`
		Phone           *string `json:"phone"`
		Name            string  `json:"name"`
		IsEmailVerified bool    `json:"isEmailVerified"`
		Role            string  `json:"role"`
		CreatedAt       string  `json:"createdAt"`
		UpdatedAt       string  `json:"updatedAt"`
	}{
		CustomerId:      c.CustomerId,
		ExternalId:      externalId,
		Email:           email,
		Phone:           phone,
		Name:            c.Name,
		IsEmailVerified: c.IsEmailVerified,
		Role:            c.Role,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type MemberResp struct {
	MemberId  string
	Name      string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m MemberResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		MemberId  string `json:"memberId"`
		Name      string `json:"name"`
		Role      string `json:"role"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		MemberId:  m.MemberId,
		Name:      m.Name,
		Role:      m.Role,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type NewLabelReq struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ThChatReq struct {
	Message string `json:"message"`
}

type ThChatLabelReq struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type CreateCustomerReq struct {
	Name            string  `json:"name"`
	IsEmailVerified bool    `json:"isEmailVerified"` // defaults to false
	ExternalId      *string `json:"externalId"`      // optional
	Email           *string `json:"email"`           // optional
	Phone           *string `json:"phone"`           // optional
}

type ThreadLabelCountResp struct {
	LabelId string `json:"labelId"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Count   int    `json:"count"`
}

type ThreadCountResp struct {
	Active             int                    `json:"active"`
	NeedsFirstResponse int                    `json:"needsFirstResponse"`
	WaitingOnCustomer  int                    `json:"waitingOnCustomer"`
	HoldCount          int                    `json:"hold"`
	NeedsNextResponse  int                    `json:"needsNextResponse"`
	AssignedToMe       int                    `json:"assignedToMe"`
	Unassigned         int                    `json:"unassigned"`
	OtherAssigned      int                    `json:"otherAssigned"`
	Labels             []ThreadLabelCountResp `json:"labels"`
}

type ThreadMetricsResp struct {
	Count ThreadCountResp `json:"count"`
}

type CreateWidgetReq struct {
	Name          string                  `json:"name"`
	Configuration *map[string]interface{} `json:"configuration"`
}

type WidgetResp struct {
	WidgetId      string                 `json:"widgetId"`
	Name          string                 `json:"name"`
	Configuration map[string]interface{} `json:"configuration"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

func (w WidgetResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WidgetId      string                 `json:"widgetId"`
		Name          string                 `json:"name"`
		Configuration map[string]interface{} `json:"configuration"`
		CreatedAt     string                 `json:"createdAt"`
		UpdatedAt     string                 `json:"updatedAt"`
	}{
		WidgetId:      w.WidgetId,
		Name:          w.Name,
		Configuration: w.Configuration,
		CreatedAt:     w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     w.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type WorkspaceSecretResp struct {
	Hmac      string    `json:"hmac"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (sk WorkspaceSecretResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		Hmac      string `json:"hmac"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Hmac:      sk.Hmac,
		CreatedAt: sk.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sk.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type CustomerActorResp struct {
	CustomerId string
	Name       string
}

func (c CustomerActorResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		CustomerId string `json:"customerId"`
		Name       string `json:"name"`
	}{
		CustomerId: c.CustomerId,
		Name:       c.Name,
	}
	return json.Marshal(aux)
}

type MemberActorResp struct {
	MemberId string
	Name     string
}

func (m MemberActorResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		MemberId string `json:"memberId"`
		Name     string `json:"name"`
	}{
		MemberId: m.MemberId,
		Name:     m.Name,
	}
	return json.Marshal(aux)
}

type ThreadResp struct {
	ThreadId        string
	Customer        CustomerActorResp
	Title           string
	Description     string
	Status          string
	StatusChangedAt time.Time
	Stage           string
	Replied         bool
	Priority        string
	Channel         string
	PreviewText     string
	Assignee        *MemberActorResp
	LastInboundAt   *time.Time
	LastOutboundAt  *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (th ThreadResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId        string            `json:"threadId"`
		Customer        CustomerActorResp `json:"customer"`
		Title           string            `json:"title"`
		Description     string            `json:"description"`
		Status          string            `json:"status"`
		StatusChangedAt string            `json:"statusChangedAt"`
		Stage           string            `json:"stage"`
		Replied         bool              `json:"replied"`
		Priority        string            `json:"priority"`
		Channel         string            `json:"channel"`
		PreviewText     string            `json:"previewText"`
		Assignee        *MemberActorResp  `json:"assignee,omitempty"`
		LastInboundAt   *time.Time        `json:"lastInboundAt,omitempty"`
		LastOutboundAt  *time.Time        `json:"lastOutboundAt,omitempty"`
		CreatedAt       string            `json:"createdAt"`
		UpdatedAt       string            `json:"updatedAt"`
	}{
		ThreadId:        th.ThreadId,
		Customer:        th.Customer,
		Title:           th.Title,
		Description:     th.Description,
		Status:          th.Status,
		StatusChangedAt: th.StatusChangedAt.Format(time.RFC3339),
		Stage:           th.Stage,
		Replied:         th.Replied,
		Priority:        th.Priority,
		Channel:         th.Channel,
		PreviewText:     th.PreviewText,
		Assignee:        th.Assignee,
		LastInboundAt:   th.LastInboundAt,
		LastOutboundAt:  th.LastOutboundAt,
		CreatedAt:       th.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       th.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (th ThreadResp) NewResponse(thread *models.Thread) ThreadResp {
	var threadAssignee *MemberActorResp
	var lastInboundAt, lastOutboundAt *time.Time

	customer := CustomerActorResp{
		CustomerId: thread.Customer.CustomerId,
		Name:       thread.Customer.Name,
	}
	if thread.AssignedMember != nil {
		threadAssignee = &MemberActorResp{
			MemberId: thread.AssignedMember.MemberId,
			Name:     thread.AssignedMember.Name,
		}
	}

	if thread.LastInboundAt != nil {
		lastInboundAt = thread.LastInboundAt
	}
	if thread.LastOutboundAt != nil {
		lastOutboundAt = thread.LastOutboundAt
	}

	return ThreadResp{
		ThreadId:        thread.ThreadId,
		Customer:        customer,
		Title:           thread.Title,
		Description:     thread.Description,
		Status:          thread.ThreadStatus.Status,
		StatusChangedAt: thread.ThreadStatus.StatusChangedAt,
		Stage:           thread.ThreadStatus.Stage,
		Replied:         thread.Replied,
		Priority:        thread.Priority,
		Channel:         thread.Channel,
		PreviewText:     thread.PreviewText,
		Assignee:        threadAssignee,
		LastInboundAt:   lastInboundAt,
		LastOutboundAt:  lastOutboundAt,
		CreatedAt:       thread.CreatedAt,
		UpdatedAt:       thread.UpdatedAt,
	}
}

type ActivityResp struct {
	ActivityID   string
	ThreadID     string
	ActivityType string
	Customer     *CustomerActorResp
	Member       *MemberActorResp
	Body         map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m ActivityResp) MarshalJSON() ([]byte, error) {
	var customer *CustomerActorResp
	var member *MemberActorResp

	if m.Customer != nil {
		customer = m.Customer
	}
	if m.Member != nil {
		member = m.Member
	}

	aux := &struct {
		ActivityID   string                 `json:"activityId"`
		ThreadID     string                 `json:"threadId"`
		ActivityType string                 `json:"activityType"`
		Customer     *CustomerActorResp     `json:"customer,omitempty"`
		Member       *MemberActorResp       `json:"member,omitempty"`
		Body         map[string]interface{} `json:"body"`
		CreatedAt    string                 `json:"createdAt"`
		UpdatedAt    string                 `json:"updatedAt"`
	}{
		ActivityID:   m.ActivityID,
		ThreadID:     m.ThreadID,
		ActivityType: m.ActivityType,
		Customer:     customer,
		Member:       member,
		Body:         m.Body,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ActivityWithAttachmentsResp struct {
	ActivityResp
	Attachments []models.ActivityAttachment `json:"attachments"`
}

func (m ActivityWithAttachmentsResp) MarshalJSON() ([]byte, error) {
	var customer *CustomerActorResp
	var member *MemberActorResp

	if m.Customer != nil {
		customer = m.Customer
	}
	if m.Member != nil {
		member = m.Member
	}

	type attachment struct {
		AttachmentId string `json:"attachmentId"`
		ActivityID   string `json:"activityId"`
		Name         string `json:"name"`
		ContentType  string `json:"contentType"`
		ContentKey   string `json:"contentKey"`
		Spam         bool   `json:"spam"`
		HasError     bool   `json:"hasError"`
		Error        string `json:"error"`
		MD5Hash      string `json:"md5Hash"`
		CreatedAt    string `json:"createdAt"`
		UpdatedAt    string `json:"updatedAt"`
	}

	formattedAttachments := make([]attachment, len(m.Attachments))
	hasError := false

	for i, att := range m.Attachments {
		formattedAttachments[i] = attachment{
			AttachmentId: att.AttachmentId,
			ActivityID:   att.ActivityID,
			Name:         att.Name,
			ContentType:  att.ContentType,
			ContentKey:   att.ContentKey,
			Spam:         att.Spam,
			HasError:     att.HasError,
			Error:        att.Error,
			MD5Hash:      att.MD5Hash,
			CreatedAt:    att.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    att.UpdatedAt.Format(time.RFC3339),
		}
		if att.HasError {
			hasError = true
		}
	}

	aux := &struct {
		ActivityID          string                 `json:"activityId"`
		ThreadID            string                 `json:"threadId"`
		ActivityType        string                 `json:"activityType"`
		Customer            *CustomerActorResp     `json:"customer,omitempty"`
		Member              *MemberActorResp       `json:"member,omitempty"`
		Body                map[string]interface{} `json:"body"`
		CreatedAt           string                 `json:"createdAt"`
		UpdatedAt           string                 `json:"updatedAt"`
		Attachments         []attachment           `json:"attachments"`
		AttachmentsHasError bool                   `json:"attachmentsHasError"`
	}{
		ActivityID:          m.ActivityID,
		ThreadID:            m.ThreadID,
		ActivityType:        m.ActivityType,
		Customer:            customer,
		Member:              member,
		Body:                m.Body,
		CreatedAt:           m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           m.UpdatedAt.Format(time.RFC3339),
		Attachments:         formattedAttachments,
		AttachmentsHasError: hasError,
	}
	return json.Marshal(aux)
}

type LabelResp struct {
	LabelId   string `json:"labelId"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (l LabelResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		LabelId   string `json:"labelId"`
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		LabelId:   l.LabelId,
		Name:      l.Name,
		Icon:      l.Icon,
		CreatedAt: l.CreatedAt.Format(time.RFC3339),
		UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ThreadLabelResp struct {
	ThreadLabelId string    `json:"threadLabelId"`
	ThreadId      string    `json:"threadId"`
	LabelId       string    `json:"labelId"`
	Name          string    `json:"name"`
	Icon          string    `json:"icon"`
	AddedBy       string    `json:"addedBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (tl ThreadLabelResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadLabelId string `json:"threadLabelId"`
		ThreadId      string `json:"threadId"`
		LabelId       string `json:"labelId"`
		Name          string `json:"name"`
		Icon          string `json:"icon"`
		AddedBy       string `json:"addedBy"`
		CreatedAt     string `json:"createdAt"`
		UpdatedAt     string `json:"updatedAt"`
	}{
		ThreadLabelId: tl.ThreadLabelId,
		ThreadId:      tl.ThreadId,
		LabelId:       tl.LabelId,
		Name:          tl.Name,
		Icon:          tl.Icon,
		AddedBy:       tl.AddedBy,
		CreatedAt:     tl.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     tl.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type CustomerEventReq struct {
	Customer struct {
		CustomerId      *string `json:"customerId"`
		ExternalId      *string `json:"externalId"`
		Email           *string `json:"email"`
		Name            *string `json:"name"`
		IsEmailVerified *bool   `json:"isEmailVerified"`
	} `json:"customer"`
	Title      string                  `json:"title"`
	Severity   string                  `json:"severity"`
	Timestamp  string                  `json:"timestamp"`
	Components []models.EventComponent `json:"components"`
}

type CustomerEventAddedResp struct {
	EventID   string    `json:"eventId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CustomerEventResp struct {
	EventID    string                  `json:"eventId"`
	Title      string                  `json:"title"`
	Severity   string                  `json:"severity"`
	Timestamp  time.Time               `json:"timestamp"`
	Components []models.EventComponent `json:"components"`
	Customer   CustomerActorResp       `json:"customer"`
	CreatedAt  time.Time               `json:"createdAt"`
	UpdatedAt  time.Time               `json:"updatedAt"`
}

func (cv CustomerEventResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		EventID    string                  `json:"eventId"`
		Title      string                  `json:"title"`
		Severity   string                  `json:"severity"`
		Timestamp  string                  `json:"timestamp"`
		Components []models.EventComponent `json:"components"`
		Customer   CustomerActorResp       `json:"customer"`
		CreatedAt  string                  `json:"createdAt"`
		UpdatedAt  string                  `json:"updatedAt"`
	}{
		EventID:    cv.EventID,
		Title:      cv.Title,
		Severity:   cv.Severity,
		Timestamp:  cv.Timestamp.Format(time.RFC3339),
		Components: cv.Components,
		Customer:   cv.Customer,
		CreatedAt:  cv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  cv.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (cv CustomerEventResp) NewResponse(event *models.Event) CustomerEventResp {
	return CustomerEventResp{
		EventID:    event.EventID,
		Title:      event.Title,
		Severity:   event.Severity.String(),
		Timestamp:  event.Timestamp,
		Components: event.Components,
		Customer: CustomerActorResp{
			CustomerId: event.Customer.CustomerId,
			Name:       event.Customer.Name,
		},
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

type CreatePostmarkMailServer struct {
	Email string `json:"email"`
}

type AddPostmarkMailServerDNS struct {
	Domain string `json:"domain"`
}

// ReplyThreadMailReq represents the reply thread mail request body
type ReplyThreadMailReq struct {
	HTMLBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`
}
