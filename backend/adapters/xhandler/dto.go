package xhandler

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/zyghq/zyg/models"
)

type CustomerResp struct {
	CustomerId        string
	ExternalId        sql.NullString
	Email             sql.NullString
	Phone             sql.NullString
	Name              string
	AvatarUrl         string
	IsVerified        bool
	Role              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RequireIdentities []string
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
		CustomerId        string   `json:"customerId"`
		ExternalId        *string  `json:"externalId"`
		Email             *string  `json:"email"`
		Phone             *string  `json:"phone"`
		Name              string   `json:"name"`
		AvatarUrl         string   `json:"avatarUrl"`
		IsVerified        bool     `json:"isVerified"`
		Role              string   `json:"role"`
		CreatedAt         string   `json:"createdAt"`
		UpdatedAt         string   `json:"updatedAt"`
		RequireIdentities []string `json:"requireIdentities"`
	}{
		CustomerId:        c.CustomerId,
		ExternalId:        externalId,
		Email:             email,
		Phone:             phone,
		Name:              c.Name,
		AvatarUrl:         c.AvatarUrl,
		IsVerified:        c.IsVerified,
		Role:              c.Role,
		CreatedAt:         c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         c.UpdatedAt.Format(time.RFC3339),
		RequireIdentities: c.RequireIdentities,
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
	IsVerified         *bool           `json:"isVerified"`
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
	ThreadId           string
	Customer           ThCustomerResp
	Title              string
	Description        string
	Sequence           int
	Status             string
	Read               bool
	Replied            bool
	Priority           string
	Spam               bool
	Channel            string
	PreviewText        string
	Assignee           *ThMemberResp
	InboundFirstSeqId  *string
	InboundLastSeqId   *string
	InboundCustomer    *ThCustomerResp
	OutboundFirstSeqId *string
	OutboundLastSeqId  *string
	OutboundMember     *ThMemberResp
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (th ThreadResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId           string          `json:"threadId"`
		Customer           ThCustomerResp  `json:"customer"`
		Title              string          `json:"title"`
		Description        string          `json:"description"`
		Sequence           int             `json:"sequence"`
		Status             string          `json:"status"`
		Read               bool            `json:"read"`
		Replied            bool            `json:"replied"`
		Priority           string          `json:"priority"`
		Spam               bool            `json:"spam"`
		Channel            string          `json:"channel"`
		PreviewText        string          `json:"previewText"`
		Assignee           *ThMemberResp   `json:"assignee,omitempty"`
		InboundFirstSeqId  *string         `json:"inboundFirstSeqId,omitempty"`
		InboundLastSeqId   *string         `json:"inboundLastSeqId,omitempty"`
		InboundCustomer    *ThCustomerResp `json:"inboundCustomer,omitempty"`
		OutboundFirstSeqId *string         `json:"outboundFirstSeqId,omitempty"`
		OutboundLastSeqId  *string         `json:"outboundLastSeqId,omitempty"`
		OutboundMember     *ThMemberResp   `json:"outboundMember,omitempty"`
		CreatedAt          string          `json:"createdAt"`
		UpdatedAt          string          `json:"updatedAt"`
	}{
		ThreadId:           th.ThreadId,
		Customer:           th.Customer,
		Title:              th.Title,
		Description:        th.Description,
		Sequence:           th.Sequence,
		Status:             th.Status,
		Read:               th.Read,
		Replied:            th.Replied,
		Priority:           th.Priority,
		Spam:               th.Spam,
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
	var threadAssignee, outboundMember *ThMemberResp
	var inboundCustomer *ThCustomerResp
	var inboundFirstSeqId, inboundLastSeqId, outboundFirstSeqId, outboundLastSeqId *string

	customer := ThCustomerResp{
		CustomerId: thread.CustomerId,
		Name:       thread.CustomerName,
	}

	if thread.AssigneeId.Valid {
		threadAssignee = &ThMemberResp{
			MemberId: thread.AssigneeId.String,
			Name:     thread.AssigneeName.String,
		}
	}

	if thread.InboundMessage != nil {
		inboundCustomer = &ThCustomerResp{
			CustomerId: thread.InboundMessage.CustomerId,
			Name:       thread.InboundMessage.CustomerName,
		}
		inboundFirstSeqId = &thread.InboundMessage.FirstSeqId
		inboundLastSeqId = &thread.InboundMessage.LastSeqId
	}

	if thread.OutboundMessage != nil {
		outboundMember = &ThMemberResp{
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
		Sequence:           thread.Sequence,
		Status:             thread.Status,
		Read:               thread.Read,
		Replied:            thread.Replied,
		Priority:           thread.Priority,
		Spam:               thread.Spam,
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

type ThreadChatResp struct {
	ThreadId           string
	Customer           ThCustomerResp
	Title              string
	Description        string
	Sequence           int
	Status             string
	Read               bool
	Replied            bool
	Priority           string
	Spam               bool
	Channel            string
	PreviewText        string
	Assignee           *ThMemberResp
	InboundFirstSeqId  *string
	InboundLastSeqId   *string
	InboundCustomer    *ThCustomerResp
	OutboundFirstSeqId *string
	OutboundLastSeqId  *string
	OutboundMember     *ThMemberResp
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Chat               ChatResp `json:"chat"`
}

func (t ThreadChatResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId           string          `json:"threadId"`
		Customer           ThCustomerResp  `json:"customer"`
		Title              string          `json:"title"`
		Description        string          `json:"description"`
		Sequence           int             `json:"sequence"`
		Status             string          `json:"status"`
		Read               bool            `json:"read"`
		Replied            bool            `json:"replied"`
		Priority           string          `json:"priority"`
		Spam               bool            `json:"spam"`
		Channel            string          `json:"channel"`
		PreviewText        string          `json:"previewText"`
		Assignee           *ThMemberResp   `json:"assignee,omitempty"`
		InboundFirstSeqId  *string         `json:"inboundFirstSeqId,omitempty"`
		InboundLastSeqId   *string         `json:"inboundLastSeqId,omitempty"`
		InboundCustomer    *ThCustomerResp `json:"inboundCustomer,omitempty"`
		OutboundFirstSeqId *string         `json:"outboundFirstSeqId,omitempty"`
		OutboundLastSeqId  *string         `json:"outboundLastSeqId,omitempty"`
		OutboundMember     *ThMemberResp   `json:"outboundMember,omitempty"`
		CreatedAt          string          `json:"createdAt"`
		UpdatedAt          string          `json:"updatedAt"`
		Chat               ChatResp        `json:"chat"`
	}{
		ThreadId:           t.ThreadId,
		Customer:           t.Customer,
		Title:              t.Title,
		Description:        t.Description,
		Sequence:           t.Sequence,
		Status:             t.Status,
		Read:               t.Read,
		Replied:            t.Replied,
		Priority:           t.Priority,
		Spam:               t.Spam,
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

type CustomerIdentitiesReq struct {
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	ExternalId *string `json:"externalId"`
}
