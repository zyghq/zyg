package xhandler

import (
	"database/sql"
	"encoding/json"
	"time"
)

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
