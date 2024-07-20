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

type WidgetInitReqPayload struct {
	AnonymousId        *string         `json:"anonymousId"`
	CustomerHash       *string         `json:"customerHash"`
	CustomerExternalId *string         `json:"customerExternalId"`
	CustomerEmail      *string         `json:"customerEmail"`
	CustomerPhone      *string         `json:"customerPhone"`
	Traits             *CustomerTraits `json:"traits"`
}

type WidgetInitRespPayload struct {
	Jwt        string         `json:"jwt"`
	Create     bool           `json:"create"`
	IsVerified bool           `json:"isVerified"`
	Name       sql.NullString `json:"name"`
}

func (w WidgetInitRespPayload) MarshalJSON() ([]byte, error) {
	var name *string
	if w.Name.Valid {
		name = &w.Name.String
	}
	aux := &struct {
		Jwt        string  `json:"jwt"`
		Create     bool    `json:"create"`
		IsVerified bool    `json:"isVerified"`
		Name       *string `json:"name"`
	}{
		Jwt:        w.Jwt,
		Create:     w.Create,
		IsVerified: w.IsVerified,
		Name:       name,
	}
	return json.Marshal(aux)
}

type ThChatReqPayload struct {
	Message string `json:"message"`
}

type ThCustomerRespPayload struct {
	CustomerId string
	Name       sql.NullString
}

func (c ThCustomerRespPayload) MarshalJSON() ([]byte, error) {
	var name *string
	if c.Name.Valid {
		name = &c.Name.String
	}
	aux := &struct {
		CustomerId string  `json:"customerId"`
		Name       *string `json:"name"`
	}{
		CustomerId: c.CustomerId,
		Name:       name,
	}
	return json.Marshal(aux)
}

type ThMemberRespPayload struct {
	MemberId string
	Name     sql.NullString
}

func (m ThMemberRespPayload) MarshalJSON() ([]byte, error) {
	var name *string
	if m.Name.Valid {
		name = &m.Name.String
	}
	aux := &struct {
		MemberId string  `json:"memberId"`
		Name     *string `json:"name"`
	}{
		MemberId: m.MemberId,
		Name:     name,
	}
	return json.Marshal(aux)
}

type ThChatMessageRespPayload struct {
	ThreadChatId        string
	ThreadChatMessageId string
	Body                string
	Sequence            int
	Customer            *ThCustomerRespPayload
	Member              *ThMemberRespPayload
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (thcresp ThChatMessageRespPayload) MarshalJSON() ([]byte, error) {
	var customer *ThCustomerRespPayload
	var member *ThMemberRespPayload

	if thcresp.Customer != nil {
		customer = thcresp.Customer
	}

	if thcresp.Member != nil {
		member = thcresp.Member
	}

	aux := &struct {
		ThreadChatId        string                 `json:"threadChatId"`
		ThreadChatMessageId string                 `json:"threadChatMessageId"`
		Body                string                 `json:"body"`
		Sequence            int                    `json:"sequence"`
		Customer            *ThCustomerRespPayload `json:"customer,omitempty"`
		Member              *ThMemberRespPayload   `json:"member,omitempty"`
		CreatedAt           string                 `json:"createdAt"`
		UpdatedAt           string                 `json:"updatedAt"`
	}{
		ThreadChatId:        thcresp.ThreadChatId,
		ThreadChatMessageId: thcresp.ThreadChatMessageId,
		Body:                thcresp.Body,
		Sequence:            thcresp.Sequence,
		Customer:            customer,
		Member:              member,
		CreatedAt:           thcresp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           thcresp.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ThChatRespPayload struct {
	ThreadChatId string
	Sequence     int
	Status       string
	Read         bool
	Replied      bool
	Priority     string
	Customer     ThCustomerRespPayload
	Assignee     *ThMemberRespPayload
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Messages     []ThChatMessageRespPayload
}

func (thresp ThChatRespPayload) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberRespPayload

	if thresp.Assignee != nil {
		assignee = thresp.Assignee
	}

	aux := &struct {
		ThreadChatId string                     `json:"threadChatId"`
		Sequence     int                        `json:"sequence"`
		Status       string                     `json:"status"`
		Read         bool                       `json:"read"`
		Replied      bool                       `json:"replied"`
		Priority     string                     `json:"priority"`
		Customer     ThCustomerRespPayload      `json:"customer"`
		Assignee     *ThMemberRespPayload       `json:"assignee"`
		CreatedAt    string                     `json:"createdAt"`
		UpdatedAt    string                     `json:"updatedAt"`
		Messages     []ThChatMessageRespPayload `json:"messages"`
	}{
		ThreadChatId: thresp.ThreadChatId,
		Sequence:     thresp.Sequence,
		Status:       thresp.Status,
		Read:         thresp.Read,
		Replied:      thresp.Replied,
		Priority:     thresp.Priority,
		Customer:     thresp.Customer,
		Assignee:     assignee,
		CreatedAt:    thresp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    thresp.UpdatedAt.Format(time.RFC3339),
		Messages:     thresp.Messages,
	}
	return json.Marshal(aux)
}
