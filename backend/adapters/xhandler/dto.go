package xhandler

import (
	"database/sql"
	"encoding/json"
	"time"
)

type CustomerTraits struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Name      *string `json:"name"`
}

type WidgetInitReq struct {
	AnonId             *string         `json:"anonId"`
	CustomerHash       *string         `json:"customerHash"`
	CustomerExternalId *string         `json:"customerExternalId"`
	CustomerEmail      *string         `json:"customerEmail"`
	CustomerPhone      *string         `json:"customerPhone"`
	Traits             *CustomerTraits `json:"traits"`
}

type WidgetInitResp struct {
	Jwt        string         `json:"jwt"`
	Create     bool           `json:"create"`
	IsVerified bool           `json:"isVerified"`
	Name       string         `json:"name"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	ExternalId sql.NullString `json:"externalId"`
}

func (w WidgetInitResp) MarshalJSON() ([]byte, error) {
	var email *string
	if w.Email.Valid {
		email = &w.Email.String
	}

	var phone *string
	if w.Phone.Valid {
		phone = &w.Phone.String
	}

	var externalId *string
	if w.ExternalId.Valid {
		externalId = &w.ExternalId.String
	}

	aux := &struct {
		Jwt        string  `json:"jwt"`
		Create     bool    `json:"create"`
		IsVerified bool    `json:"isVerified"`
		Name       string  `json:"name"`
		Email      *string `json:"email"`
		Phone      *string `json:"phone"`
		ExternalId *string `json:"externalId"`
	}{
		Jwt:        w.Jwt,
		Create:     w.Create,
		IsVerified: w.IsVerified,
		Name:       w.Name,
		Email:      email,
		Phone:      phone,
		ExternalId: externalId,
	}
	return json.Marshal(aux)
}

type ThChatReq struct {
	Message string `json:"message"`
}

type ThCustomerResp struct {
	CustomerId string `json:"customerId"`
	Name       string `json:"name"`
}

type ThMemberResp struct {
	MemberId string `json:"memberId"`
	Name     string `json:"name"`
}

type ThreadResp struct {
	ThreadId        string
	Customer        ThCustomerResp
	Title           string
	Description     string
	Sequence        int
	Status          string
	Read            bool
	Replied         bool
	Priority        string
	Spam            bool
	Channel         string
	PreviewText     string
	Assignee        *ThMemberResp
	IngressFirstSeq sql.NullInt64
	IngressLastSeq  sql.NullInt64
	IngressCustomer *ThCustomerResp
	EgressFirstSeq  sql.NullInt64
	EgressLastSeq   sql.NullInt64
	EgressMember    *ThMemberResp
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (t ThreadResp) MarshalJSON() ([]byte, error) {
	var customer ThCustomerResp
	var assignee *ThMemberResp
	var ingressCustomer *ThCustomerResp
	var egressMember *ThMemberResp

	var ingressFirstSeq, ingressLastSeq *int64
	var egressFirstSeq, egressLastSeq *int64

	if t.Assignee != nil {
		assignee = t.Assignee
	}

	if t.IngressCustomer != nil {
		ingressCustomer = t.IngressCustomer
	}

	if t.EgressMember != nil {
		egressMember = t.EgressMember
	}

	if t.IngressFirstSeq.Valid {
		ingressFirstSeq = &t.IngressFirstSeq.Int64
	}

	if t.IngressLastSeq.Valid {
		ingressLastSeq = &t.IngressLastSeq.Int64
	}

	if t.EgressFirstSeq.Valid {
		egressFirstSeq = &t.EgressFirstSeq.Int64
	}

	if t.EgressLastSeq.Valid {
		egressLastSeq = &t.EgressLastSeq.Int64
	}

	aux := &struct {
		ThreadId        string          `json:"threadId"`
		Customer        ThCustomerResp  `json:"customer"`
		Title           string          `json:"title"`
		Description     string          `json:"description"`
		Sequence        int             `json:"sequence"`
		Status          string          `json:"status"`
		Read            bool            `json:"read"`
		Replied         bool            `json:"replied"`
		Priority        string          `json:"priority"`
		Spam            bool            `json:"spam"`
		Channel         string          `json:"channel"`
		PreviewText     string          `json:"previewText"`
		Assignee        *ThMemberResp   `json:"assignee,omitempty"`
		IngressFirstSeq *int64          `json:"ingressFirstSeq,omitempty"`
		IngressLastSeq  *int64          `json:"ingressLastSeq,omitempty"`
		IngressCustomer *ThCustomerResp `json:"ingressCustomer,omitempty"`
		EgressFirstSeq  *int64          `json:"egressFirstSeq,omitempty"`
		EgressLastSeq   *int64          `json:"egressLastSeq,omitempty"`
		EgressMember    *ThMemberResp   `json:"egressMember,omitempty"`
		CreatedAt       string          `json:"createdAt"`
		UpdatedAt       string          `json:"updatedAt"`
	}{
		ThreadId:        t.ThreadId,
		Customer:        customer,
		Title:           t.Title,
		Description:     t.Description,
		Sequence:        t.Sequence,
		Status:          t.Status,
		Read:            t.Read,
		Replied:         t.Replied,
		Priority:        t.Priority,
		Spam:            t.Spam,
		Channel:         t.Channel,
		PreviewText:     t.PreviewText,
		Assignee:        assignee,
		IngressFirstSeq: ingressFirstSeq,
		IngressLastSeq:  ingressLastSeq,
		IngressCustomer: ingressCustomer,
		EgressFirstSeq:  egressFirstSeq,
		EgressLastSeq:   egressLastSeq,
		EgressMember:    egressMember,
		CreatedAt:       t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       t.UpdatedAt.Format(time.RFC3339),
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

// type ThreadChatResp struct {
// 	ThreadResp
// 	Chat ChatResp `json:"chat"`
// }

type ThreadChatResp struct {
	ThreadId        string
	Customer        ThCustomerResp
	Title           string
	Description     string
	Sequence        int
	Status          string
	Read            bool
	Replied         bool
	Priority        string
	Spam            bool
	Channel         string
	PreviewText     string
	Assignee        *ThMemberResp
	IngressFirstSeq sql.NullInt64
	IngressLastSeq  sql.NullInt64
	IngressCustomer *ThCustomerResp
	EgressFirstSeq  sql.NullInt64
	EgressLastSeq   sql.NullInt64
	EgressMember    *ThMemberResp
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Chat            ChatResp `json:"chat"`
}

func (t ThreadChatResp) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberResp
	var ingressCustomer *ThCustomerResp
	var egressMember *ThMemberResp

	var ingressFirstSeq, ingressLastSeq *int64
	var egressFirstSeq, egressLastSeq *int64

	if t.Assignee != nil {
		assignee = t.Assignee
	}

	if t.IngressCustomer != nil {
		ingressCustomer = t.IngressCustomer
	}

	if t.EgressMember != nil {
		egressMember = t.EgressMember
	}

	if t.IngressFirstSeq.Valid {
		ingressFirstSeq = &t.IngressFirstSeq.Int64
	}

	if t.IngressLastSeq.Valid {
		ingressLastSeq = &t.IngressLastSeq.Int64
	}

	if t.EgressFirstSeq.Valid {
		egressFirstSeq = &t.EgressFirstSeq.Int64
	}

	if t.EgressLastSeq.Valid {
		egressLastSeq = &t.EgressLastSeq.Int64
	}

	aux := &struct {
		ThreadId        string          `json:"threadId"`
		Customer        ThCustomerResp  `json:"customer"`
		Title           string          `json:"title"`
		Description     string          `json:"description"`
		Sequence        int             `json:"sequence"`
		Status          string          `json:"status"`
		Read            bool            `json:"read"`
		Replied         bool            `json:"replied"`
		Priority        string          `json:"priority"`
		Spam            bool            `json:"spam"`
		Channel         string          `json:"channel"`
		PreviewText     string          `json:"previewText"`
		Assignee        *ThMemberResp   `json:"assignee,omitempty"`
		IngressFirstSeq *int64          `json:"ingressFirstSeq,omitempty"`
		IngressLastSeq  *int64          `json:"ingressLastSeq,omitempty"`
		IngressCustomer *ThCustomerResp `json:"ingressCustomer,omitempty"`
		EgressFirstSeq  *int64          `json:"egressFirstSeq,omitempty"`
		EgressLastSeq   *int64          `json:"egressLastSeq,omitempty"`
		EgressMember    *ThMemberResp   `json:"egressMember,omitempty"`
		CreatedAt       string          `json:"createdAt"`
		UpdatedAt       string          `json:"updatedAt"`
		Chat            ChatResp        `json:"chat"`
	}{
		ThreadId:        t.ThreadId,
		Customer:        t.Customer,
		Title:           t.Title,
		Description:     t.Description,
		Sequence:        t.Sequence,
		Status:          t.Status,
		Read:            t.Read,
		Replied:         t.Replied,
		Priority:        t.Priority,
		Spam:            t.Spam,
		Channel:         t.Channel,
		PreviewText:     t.PreviewText,
		Assignee:        assignee,
		IngressFirstSeq: ingressFirstSeq,
		IngressLastSeq:  ingressLastSeq,
		IngressCustomer: ingressCustomer,
		EgressFirstSeq:  egressFirstSeq,
		EgressLastSeq:   egressLastSeq,
		EgressMember:    egressMember,
		CreatedAt:       t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       t.UpdatedAt.Format(time.RFC3339),
		Chat:            t.Chat,
	}
	return json.Marshal(aux)
}

type CustomerIdentitiesReq struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	External *string `json:"external"`
}

type AddCustomerIdentitiesResp struct {
	IsVerified bool           `json:"isVerified"`
	Name       string         `json:"name"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	ExternalId sql.NullString `json:"externalId"`
}

func (w AddCustomerIdentitiesResp) MarshalJSON() ([]byte, error) {
	var email *string
	if w.Email.Valid {
		email = &w.Email.String
	}

	var phone *string
	if w.Phone.Valid {
		phone = &w.Phone.String
	}

	var externalId *string
	if w.ExternalId.Valid {
		externalId = &w.ExternalId.String
	}

	aux := &struct {
		IsVerified bool    `json:"isVerified"`
		Name       string  `json:"name"`
		Email      *string `json:"email"`
		Phone      *string `json:"phone"`
		ExternalId *string `json:"externalId"`
	}{
		IsVerified: w.IsVerified,
		Name:       w.Name,
		Email:      email,
		Phone:      phone,
		ExternalId: externalId,
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
