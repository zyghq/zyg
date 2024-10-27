package xhandler

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/zyghq/zyg/models"
)

// CustomerResp represents the API response for a customer.
type CustomerResp struct {
	CustomerId      string
	ExternalId      sql.NullString
	Email           sql.NullString
	IsEmailVerified bool
	IsEmailPrimary  bool
	Phone           sql.NullString
	Name            string
	AvatarUrl       string
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
		IsEmailVerified bool    `json:"isEmailVerified"`
		IsEmailPrimary  bool    `json:"isEmailPrimary"`
		Phone           *string `json:"phone"`
		Name            string  `json:"name"`
		AvatarUrl       string  `json:"avatarUrl"`
		Role            string  `json:"role"`
		CreatedAt       string  `json:"createdAt"`
		UpdatedAt       string  `json:"updatedAt"`
	}{
		CustomerId:      c.CustomerId,
		ExternalId:      externalId,
		Email:           email,
		IsEmailVerified: c.IsEmailVerified,
		IsEmailPrimary:  c.IsEmailPrimary,
		Phone:           phone,
		Name:            c.Name,
		AvatarUrl:       c.AvatarUrl,
		Role:            c.Role,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type CustomerTraits struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Name      *string `json:"name"`
}

type WidgetConfig struct {
	DomainsOnly    bool     `json:"domainsOnly"`
	Domains        []string `json:"domains"`
	BubblePosition string   `json:"bubblePosition"`
	HeaderColor    string   `json:"headerColor"`
	ProfilePicture string   `json:"profilePicture"`
	IconColor      string   `json:"iconColor"`
}

type WidgetInitReq struct {
	SessionId          *string         `json:"sessionId"`
	IsEmailVerified    *bool           `json:"IsEmailVerified"`
	CustomerHash       *string         `json:"customerHash"`
	CustomerExternalId *string         `json:"customerExternalId"`
	CustomerEmail      *string         `json:"customerEmail"`
	CustomerPhone      *string         `json:"customerPhone"`
	Traits             *CustomerTraits `json:"traits"`
}

type WidgetInitResp struct {
	Jwt    string `json:"jwt"`
	Create bool   `json:"create"`
	CustomerResp
}

func (w WidgetInitResp) MarshalJSON() ([]byte, error) {
	customerJson, err := json.Marshal(w.CustomerResp)
	if err != nil {
		return nil, err
	}

	var mergedMap map[string]interface{}
	err = json.Unmarshal(customerJson, &mergedMap)
	if err != nil {
		return nil, err
	}

	mergedMap["jwt"] = w.Jwt
	mergedMap["create"] = w.Create
	return json.Marshal(mergedMap)
}

type CreateThreadChatReq struct {
	Message      string  `json:"message"`
	Email        *string `json:"email"`
	Name         *string `json:"name"`
	RedirectHost *string `json:"host"`
}

type MessageThreadReq struct {
	Message string `json:"message"`
}

type CustomerActorResp struct {
	CustomerId string `json:"customerId"`
	Name       string `json:"name"`
}

type MemberActorResp struct {
	MemberId string `json:"memberId"`
	Name     string `json:"name"`
}

type ThreadResp struct {
	ThreadId           string
	Customer           CustomerActorResp
	Title              string
	Description        string
	Status             string
	Replied            bool
	Priority           string
	Channel            string
	PreviewText        string
	Assignee           *MemberActorResp
	InboundFirstSeqId  *string
	InboundLastSeqId   *string
	InboundCustomer    *CustomerActorResp
	OutboundFirstSeqId *string
	OutboundLastSeqId  *string
	OutboundMember     *MemberActorResp
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (th ThreadResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId           string             `json:"threadId"`
		Customer           CustomerActorResp  `json:"customer"`
		Title              string             `json:"title"`
		Description        string             `json:"description"`
		Status             string             `json:"status"`
		Replied            bool               `json:"replied"`
		Priority           string             `json:"priority"`
		Channel            string             `json:"channel"`
		PreviewText        string             `json:"previewText"`
		Assignee           *MemberActorResp   `json:"assignee,omitempty"`
		InboundFirstSeqId  *string            `json:"inboundFirstSeqId,omitempty"`
		InboundLastSeqId   *string            `json:"inboundLastSeqId,omitempty"`
		InboundCustomer    *CustomerActorResp `json:"inboundCustomer,omitempty"`
		OutboundFirstSeqId *string            `json:"outboundFirstSeqId,omitempty"`
		OutboundLastSeqId  *string            `json:"outboundLastSeqId,omitempty"`
		OutboundMember     *MemberActorResp   `json:"outboundMember,omitempty"`
		CreatedAt          string             `json:"createdAt"`
		UpdatedAt          string             `json:"updatedAt"`
	}{
		ThreadId:           th.ThreadId,
		Customer:           th.Customer,
		Title:              th.Title,
		Description:        th.Description,
		Status:             th.Status,
		Replied:            th.Replied,
		Priority:           th.Priority,
		Channel:            th.Channel,
		PreviewText:        th.PreviewText,
		Assignee:           th.Assignee,
		InboundFirstSeqId:  th.InboundFirstSeqId,
		InboundLastSeqId:   th.InboundLastSeqId,
		InboundCustomer:    th.InboundCustomer,
		OutboundFirstSeqId: th.OutboundFirstSeqId,
		OutboundLastSeqId:  th.OutboundLastSeqId,
		OutboundMember:     th.OutboundMember,
		CreatedAt:          th.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          th.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (th ThreadResp) NewResponse(thread *models.Thread) ThreadResp {
	var threadAssignee, outboundMember *MemberActorResp
	var inboundCustomer *CustomerActorResp
	var inboundFirstSeqId, inboundLastSeqId, outboundFirstSeqId, outboundLastSeqId *string

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

	if thread.InboundMessage != nil {
		inboundCustomer = &CustomerActorResp{
			CustomerId: thread.InboundMessage.CustomerId,
			Name:       thread.InboundMessage.CustomerName,
		}
		inboundFirstSeqId = &thread.InboundMessage.FirstSeqId
		inboundLastSeqId = &thread.InboundMessage.LastSeqId
	}

	if thread.OutboundMessage != nil {
		outboundMember = &MemberActorResp{
			MemberId: thread.OutboundMessage.MemberId,
			Name:     thread.OutboundMessage.MemberName,
		}
		outboundFirstSeqId = &thread.OutboundMessage.FirstSeqId
		outboundLastSeqId = &thread.OutboundMessage.LastSeqId

	}

	return ThreadResp{
		ThreadId:           thread.ThreadId,
		Customer:           customer,
		Title:              thread.Title,
		Description:        thread.Description,
		Status:             thread.ThreadStatus.Status,
		Replied:            thread.Replied,
		Priority:           thread.Priority,
		Channel:            thread.Channel,
		PreviewText:        thread.CustomerPreviewText(),
		Assignee:           threadAssignee,
		InboundFirstSeqId:  inboundFirstSeqId,
		InboundLastSeqId:   inboundLastSeqId,
		InboundCustomer:    inboundCustomer,
		OutboundFirstSeqId: outboundFirstSeqId,
		OutboundLastSeqId:  outboundLastSeqId,
		OutboundMember:     outboundMember,
		CreatedAt:          thread.CreatedAt,
		UpdatedAt:          thread.UpdatedAt,
	}
}

type ChatResp struct {
	ThreadId  string
	ChatId    string
	Body      string
	Sequence  int
	Customer  *CustomerActorResp
	Member    *MemberActorResp
	IsHead    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ch ChatResp) MarshalJSON() ([]byte, error) {
	var customer *CustomerActorResp
	var member *MemberActorResp

	if ch.Customer != nil {
		customer = ch.Customer
	}

	if ch.Member != nil {
		member = ch.Member
	}

	aux := &struct {
		ThreadId  string             `json:"threadId"`
		ChatId    string             `json:"chatId"`
		Body      string             `json:"body"`
		Sequence  int                `json:"sequence"`
		IsHead    bool               `json:"isHead"`
		Customer  *CustomerActorResp `json:"customer,omitempty"`
		Member    *MemberActorResp   `json:"member,omitempty"`
		CreatedAt string             `json:"createdAt"`
		UpdatedAt string             `json:"updatedAt"`
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

type ThreadChatResp struct {
	ThreadId           string
	Customer           CustomerActorResp
	Title              string
	Description        string
	Status             string
	Replied            bool
	Priority           string
	Channel            string
	PreviewText        string
	Assignee           *MemberActorResp
	InboundFirstSeqId  *string
	InboundLastSeqId   *string
	InboundCustomer    *CustomerActorResp
	OutboundFirstSeqId *string
	OutboundLastSeqId  *string
	OutboundMember     *MemberActorResp
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Chat               ChatResp `json:"chat"`
}

func (t ThreadChatResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId           string             `json:"threadId"`
		Customer           CustomerActorResp  `json:"customer"`
		Title              string             `json:"title"`
		Description        string             `json:"description"`
		Status             string             `json:"status"`
		Replied            bool               `json:"replied"`
		Priority           string             `json:"priority"`
		Channel            string             `json:"channel"`
		PreviewText        string             `json:"previewText"`
		Assignee           *MemberActorResp   `json:"assignee,omitempty"`
		InboundFirstSeqId  *string            `json:"inboundFirstSeqId,omitempty"`
		InboundLastSeqId   *string            `json:"inboundLastSeqId,omitempty"`
		InboundCustomer    *CustomerActorResp `json:"inboundCustomer,omitempty"`
		OutboundFirstSeqId *string            `json:"outboundFirstSeqId,omitempty"`
		OutboundLastSeqId  *string            `json:"outboundLastSeqId,omitempty"`
		OutboundMember     *MemberActorResp   `json:"outboundMember,omitempty"`
		CreatedAt          string             `json:"createdAt"`
		UpdatedAt          string             `json:"updatedAt"`
		Chat               ChatResp           `json:"chat"`
	}{
		ThreadId:           t.ThreadId,
		Customer:           t.Customer,
		Title:              t.Title,
		Description:        t.Description,
		Status:             t.Status,
		Replied:            t.Replied,
		Priority:           t.Priority,
		Channel:            t.Channel,
		PreviewText:        t.PreviewText,
		Assignee:           t.Assignee,
		InboundFirstSeqId:  t.InboundFirstSeqId,
		InboundLastSeqId:   t.InboundLastSeqId,
		InboundCustomer:    t.InboundCustomer,
		OutboundFirstSeqId: t.OutboundFirstSeqId,
		OutboundLastSeqId:  t.OutboundLastSeqId,
		OutboundMember:     t.OutboundMember,
		CreatedAt:          t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          t.UpdatedAt.Format(time.RFC3339),
		Chat:               t.Chat,
	}
	return json.Marshal(aux)
}

func (t ThreadChatResp) NewResponse(thread *models.Thread, chat *models.Chat) ThreadChatResp {
	var threadAssignee, outboundMember, chatMember *MemberActorResp
	var inboundCustomer, chatCustomer *CustomerActorResp
	var inboundFirstSeqId, inboundLastSeqId, outboundFirstSeqId, outboundLastSeqId *string

	// for chat - either of them
	if chat.CustomerId.Valid {
		chatCustomer = &CustomerActorResp{
			CustomerId: chat.CustomerId.String,
			Name:       chat.CustomerName.String,
		}
	} else if chat.MemberId.Valid {
		chatMember = &MemberActorResp{
			MemberId: chat.MemberId.String,
			Name:     chat.MemberName.String,
		}
	}

	chatResp := ChatResp{
		ThreadId:  thread.ThreadId,
		ChatId:    chat.ChatId,
		Body:      chat.Body,
		Sequence:  chat.Sequence,
		IsHead:    chat.IsHead,
		Customer:  chatCustomer,
		Member:    chatMember,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}

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

	if thread.InboundMessage != nil {
		inboundCustomer = &CustomerActorResp{
			CustomerId: thread.InboundMessage.CustomerId,
			Name:       thread.InboundMessage.CustomerName,
		}
		inboundFirstSeqId = &thread.InboundMessage.FirstSeqId
		inboundLastSeqId = &thread.InboundMessage.LastSeqId
	}

	if thread.OutboundMessage != nil {
		outboundMember = &MemberActorResp{
			MemberId: thread.OutboundMessage.MemberId,
			Name:     thread.OutboundMessage.MemberName,
		}
		outboundFirstSeqId = &thread.OutboundMessage.FirstSeqId
		outboundLastSeqId = &thread.OutboundMessage.LastSeqId

	}

	return ThreadChatResp{
		ThreadId:           thread.ThreadId,
		Customer:           customer,
		Title:              thread.Title,
		Description:        thread.Description,
		Status:             thread.ThreadStatus.Status,
		Replied:            thread.Replied,
		Priority:           thread.Priority,
		Channel:            thread.Channel,
		PreviewText:        thread.CustomerPreviewText(),
		Assignee:           threadAssignee,
		InboundFirstSeqId:  inboundFirstSeqId,
		InboundLastSeqId:   inboundLastSeqId,
		InboundCustomer:    inboundCustomer,
		OutboundFirstSeqId: outboundFirstSeqId,
		OutboundLastSeqId:  outboundLastSeqId,
		OutboundMember:     outboundMember,
		CreatedAt:          thread.CreatedAt,
		UpdatedAt:          thread.UpdatedAt,
		Chat:               chatResp,
	}
}

type EmailProfileReq struct {
	Email           string  `json:"email"`           // claimed email identity.
	Name            string  `json:"name"`            // claimed name identity.
	RedirectHost    *string `json:"host"`            // redirect host for redirects.
	ContextThreadId *string `json:"contextThreadId"` // customer identity context thread ID.
}
