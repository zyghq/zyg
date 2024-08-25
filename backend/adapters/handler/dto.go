package handler

import (
	"database/sql"
	"encoding/json"
	"time"
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
	CustomerId string
	ExternalId sql.NullString
	Email      sql.NullString
	Phone      sql.NullString
	Name       string
	IsVerified bool
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
		CustomerId string  `json:"customerId"`
		ExternalId *string `json:"externalId"`
		Email      *string `json:"email"`
		Phone      *string `json:"phone"`
		Name       string  `json:"name"`
		IsVerified bool    `json:"isVerified"`
		Role       string  `json:"role"`
		CreatedAt  string  `json:"createdAt"`
		UpdatedAt  string  `json:"updatedAt"`
	}{
		CustomerId: c.CustomerId,
		ExternalId: externalId,
		Email:      email,
		Phone:      phone,
		Name:       c.Name,
		IsVerified: c.IsVerified,
		Role:       c.Role,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  c.UpdatedAt.Format(time.RFC3339),
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
	Name       string  `json:"name"`
	ExternalId *string `json:"externalId"` // optional
	Email      *string `json:"email"`      // optional
	Phone      *string `json:"phone"`      // optional
}

type ThreadLabelCountResp struct {
	LabelId string `json:"labelId"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Count   int    `json:"count"`
}

type ThreadCountResp struct {
	ActiveCount   int                    `json:"active"`
	DoneCount     int                    `json:"done"`
	TodoCount     int                    `json:"todo"`
	SnoozedCount  int                    `json:"snoozed"`
	AssignedToMe  int                    `json:"assignedToMe"`
	Unassigned    int                    `json:"unassigned"`
	OtherAssigned int                    `json:"otherAssigned"`
	Labels        []ThreadLabelCountResp `json:"labels"`
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

type ThCustomerResp struct {
	CustomerId string
	Name       string
}

func (c ThCustomerResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		CustomerId string `json:"customerId"`
		Name       string `json:"name"`
	}{
		CustomerId: c.CustomerId,
		Name:       c.Name,
	}
	return json.Marshal(aux)
}

type ThMemberResp struct {
	MemberId string
	Name     string
}

func (m ThMemberResp) MarshalJSON() ([]byte, error) {
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
	ThreadId          string
	Customer          ThCustomerResp
	Title             string
	Description       string
	Sequence          int
	Status            string
	Read              bool
	Replied           bool
	Priority          string
	Spam              bool
	Channel           string
	PreviewText       string
	Assignee          *ThMemberResp
	InboundFirstSeqId *string
	InboundLastSeqId  *string
	InboundCustomer   *ThCustomerResp
	EgressFirstSeq    sql.NullInt64
	EgressLastSeq     sql.NullInt64
	EgressMember      *ThMemberResp
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (t ThreadResp) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberResp
	var inboundCustomer *ThCustomerResp
	var egressMember *ThMemberResp
	var egressFirstSeq, egressLastSeq *int64

	if t.Assignee != nil {
		assignee = t.Assignee
	}

	if t.InboundCustomer != nil {
		inboundCustomer = t.InboundCustomer
	}

	if t.EgressMember != nil {
		egressMember = t.EgressMember
	}

	if t.EgressFirstSeq.Valid {
		egressFirstSeq = &t.EgressFirstSeq.Int64
	}

	if t.EgressLastSeq.Valid {
		egressLastSeq = &t.EgressLastSeq.Int64
	}

	aux := &struct {
		ThreadId          string          `json:"threadId"`
		Customer          ThCustomerResp  `json:"customer"`
		Title             string          `json:"title"`
		Description       string          `json:"description"`
		Sequence          int             `json:"sequence"`
		Status            string          `json:"status"`
		Read              bool            `json:"read"`
		Replied           bool            `json:"replied"`
		Priority          string          `json:"priority"`
		Spam              bool            `json:"spam"`
		Channel           string          `json:"channel"`
		PreviewText       string          `json:"previewText"`
		Assignee          *ThMemberResp   `json:"assignee,omitempty"`
		InboundFirstSeqId *string         `json:"inboundFirstSeqId,omitempty"`
		InboundLastSeqId  *string         `json:"inboundLastSeqId,omitempty"`
		InboundCustomer   *ThCustomerResp `json:"inboundCustomer,omitempty"`
		EgressFirstSeq    *int64          `json:"egressFirstSeq,omitempty"`
		EgressLastSeq     *int64          `json:"egressLastSeq,omitempty"`
		EgressMember      *ThMemberResp   `json:"egressMember,omitempty"`
		CreatedAt         string          `json:"createdAt"`
		UpdatedAt         string          `json:"updatedAt"`
	}{
		ThreadId:          t.ThreadId,
		Customer:          t.Customer,
		Title:             t.Title,
		Description:       t.Description,
		Sequence:          t.Sequence,
		Status:            t.Status,
		Read:              t.Read,
		Replied:           t.Replied,
		Priority:          t.Priority,
		Spam:              t.Spam,
		Channel:           t.Channel,
		PreviewText:       t.PreviewText,
		Assignee:          assignee,
		InboundFirstSeqId: t.InboundFirstSeqId,
		InboundLastSeqId:  t.InboundLastSeqId,
		InboundCustomer:   inboundCustomer,
		EgressFirstSeq:    egressFirstSeq,
		EgressLastSeq:     egressLastSeq,
		EgressMember:      egressMember,
		CreatedAt:         t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         t.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ChatResp struct {
	ThreadId  string
	ChatId    string
	Body      string
	Sequence  int
	Customer  *ThCustomerResp
	Member    *ThMemberResp
	IsHead    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ch ChatResp) MarshalJSON() ([]byte, error) {
	var customer *ThCustomerResp
	var member *ThMemberResp

	if ch.Customer != nil {
		customer = ch.Customer
	}

	if ch.Member != nil {
		member = ch.Member
	}

	aux := &struct {
		ThreadId  string          `json:"threadId"`
		ChatId    string          `json:"chatId"`
		Body      string          `json:"body"`
		Sequence  int             `json:"sequence"`
		IsHead    bool            `json:"isHead"`
		Customer  *ThCustomerResp `json:"customer,omitempty"`
		Member    *ThMemberResp   `json:"member,omitempty"`
		CreatedAt string          `json:"createdAt"`
		UpdatedAt string          `json:"updatedAt"`
	}{
		ThreadId:  ch.ThreadId,
		ChatId:    ch.ChatId,
		Body:      ch.Body,
		Sequence:  ch.Sequence,
		IsHead:    ch.IsHead,
		Customer:  customer,
		Member:    member,
		CreatedAt: ch.CreatedAt.Format(time.RFC3339),
		UpdatedAt: ch.UpdatedAt.Format(time.RFC3339),
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
