package api

import (
	"database/sql"
	"encoding/json"
	"time"
)

type PATReqPayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type WorkspaceReqPayload struct {
	Name string `json:"name"`
}

type CustomerTIReqPayload struct {
	Create   bool    `json:"create"`
	CreateBy *string `json:"createBy"` // optional
	Customer struct {
		ExternalId *string `json:"externalId"` // optional
		Email      *string `json:"email"`      // optional
		Phone      *string `json:"phone"`      // optional
	} `json:"customer"`
}

type CustomerTIRespPayload struct {
	Create     bool   `json:"create"`
	CustomerId string `json:"customerId"`
	Jwt        string `json:"jwt"`
}

type ThChatReqPayload struct {
	Message string `json:"message"`
}

type ThCustomerRespPayload struct {
	CustomerId string
	Name       sql.NullString
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
	ThreadId  string
	Sequence  int
	Status    string
	Customer  ThCustomerRespPayload
	Assignee  *ThMemberRespPayload
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []ThChatMessageRespPayload
}

func (thresp ThChatRespPayload) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberRespPayload

	if thresp.Assignee != nil {
		assignee = thresp.Assignee
	}

	aux := &struct {
		ThreadId  string                     `json:"threadId"`
		Sequence  int                        `json:"sequence"`
		Status    string                     `json:"status"`
		Customer  ThCustomerRespPayload      `json:"customer"`
		Assignee  *ThMemberRespPayload       `json:"assignee"`
		CreatedAt string                     `json:"createdAt"`
		UpdatedAt string                     `json:"updatedAt"`
		Messages  []ThChatMessageRespPayload `json:"messages"`
	}{
		ThreadId:  thresp.ThreadId,
		Sequence:  thresp.Sequence,
		Status:    thresp.Status,
		Customer:  thresp.Customer,
		Assignee:  assignee,
		CreatedAt: thresp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: thresp.UpdatedAt.Format(time.RFC3339),
		Messages:  thresp.Messages,
	}
	return json.Marshal(aux)
}

type ThreadQAReqPayload struct {
	Query string `json:"query"`
}

type ThreadQAARespPayload struct {
	AnswerId string
	Answer   string
	Eval     sql.NullInt32
	Sequence int
}

func (tqar ThreadQAARespPayload) MarshalJSON() ([]byte, error) {
	var eval *int32
	if tqar.Eval.Valid {
		eval = &tqar.Eval.Int32
	}
	aux := &struct {
		AnswerId string `json:"answerId"`
		Answer   string `json:"answer"`
		Eval     *int32 `json:"eval"`
		Sequence int    `json:"sequence"`
	}{
		AnswerId: tqar.AnswerId,
		Answer:   tqar.Answer,
		Eval:     eval,
		Sequence: tqar.Sequence,
	}
	return json.Marshal(aux)
}

type ThreadQARespPayload struct {
	ThreadId  string
	Query     string
	Sequence  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Answers   []ThreadQAARespPayload
}

func (thr ThreadQARespPayload) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId  string                 `json:"threadId"`
		Query     string                 `json:"query"`
		Sequence  int                    `json:"sequence"`
		CreatedAt string                 `json:"createdAt"`
		UpdatedAt string                 `json:"updatedAt"`
		Answers   []ThreadQAARespPayload `json:"answers"`
	}{
		ThreadId:  thr.ThreadId,
		Query:     thr.Query,
		Sequence:  thr.Sequence,
		CreatedAt: thr.CreatedAt.Format(time.RFC3339),
		UpdatedAt: thr.UpdatedAt.Format(time.RFC3339),
		Answers:   thr.Answers,
	}
	return json.Marshal(aux)
}
