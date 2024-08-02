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

type WidgetInitRespPayload struct {
	Jwt        string         `json:"jwt"`
	Create     bool           `json:"create"`
	IsVerified bool           `json:"isVerified"`
	Name       string         `json:"name"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	ExternalId sql.NullString `json:"externalId"`
}

func (w WidgetInitRespPayload) MarshalJSON() ([]byte, error) {
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

type ThChatReqPayload struct {
	Message string `json:"message"`
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
	ThreadId        string
	Sequence        int
	Status          string
	Read            bool
	Replied         bool
	Priority        string
	Assignee        *ThMemberResp
	Title           string
	Summary         string
	Spam            bool
	Channel         string
	Body            string
	MessageSequence int
	Customer        *ThCustomerResp
	Member          *ThMemberResp
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (t ThreadResp) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberResp
	var customer *ThCustomerResp
	var member *ThMemberResp

	if t.Assignee != nil {
		assignee = t.Assignee
	}

	if t.Customer != nil {
		customer = t.Customer
	}

	if t.Member != nil {
		member = t.Member
	}

	aux := &struct {
		ThreadId        string          `json:"threadId"`
		Sequence        int             `json:"sequence"`
		Status          string          `json:"status"`
		Read            bool            `json:"read"`
		Replied         bool            `json:"replied"`
		Priority        string          `json:"priority"`
		Assignee        *ThMemberResp   `json:"assignee"`
		Title           string          `json:"title"`
		Summary         string          `json:"summary"`
		Spam            bool            `json:"spam"`
		Channel         string          `json:"channel"`
		Body            string          `json:"body"`
		MessageSequence int             `json:"messageSequence"`
		Customer        *ThCustomerResp `json:"customer,omitempty"`
		Member          *ThMemberResp   `json:"member,omitempty"`
		CreatedAt       string          `json:"createdAt"`
		UpdatedAt       string          `json:"updatedAt"`
	}{
		ThreadId:        t.ThreadId,
		Sequence:        t.Sequence,
		Status:          t.Status,
		Read:            t.Read,
		Replied:         t.Replied,
		Priority:        t.Priority,
		Assignee:        assignee,
		Title:           t.Title,
		Summary:         t.Summary,
		Spam:            t.Spam,
		Channel:         t.Channel,
		Body:            t.Body,
		MessageSequence: t.MessageSequence,
		Customer:        customer,
		Member:          member,
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

type ThreadChatResp struct {
	ThreadId  string
	Sequence  int
	Status    string
	Read      bool
	Replied   bool
	Priority  string
	Customer  ThCustomerResp
	Assignee  *ThMemberResp
	Title     string
	Summary   string
	Spam      bool
	Channel   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Chat      ChatResp
}

func (thc ThreadChatResp) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberResp

	if thc.Assignee != nil {
		assignee = thc.Assignee
	}

	aux := &struct {
		ThreadId  string         `json:"threadId"`
		Sequence  int            `json:"sequence"`
		Status    string         `json:"status"`
		Read      bool           `json:"read"`
		Replied   bool           `json:"replied"`
		Priority  string         `json:"priority"`
		Customer  ThCustomerResp `json:"customer"`
		Assignee  *ThMemberResp  `json:"assignee"`
		Title     string         `json:"title"`
		Summary   string         `json:"summary"`
		Spam      bool           `json:"spam"`
		Channel   string         `json:"channel"`
		CreatedAt string         `json:"createdAt"`
		UpdatedAt string         `json:"updatedAt"`
		Chat      ChatResp       `json:"chat"`
	}{
		ThreadId:  thc.ThreadId,
		Sequence:  thc.Sequence,
		Status:    thc.Status,
		Read:      thc.Read,
		Replied:   thc.Replied,
		Priority:  thc.Priority,
		Customer:  thc.Customer,
		Assignee:  assignee,
		Title:     thc.Title,
		Summary:   thc.Summary,
		Spam:      thc.Spam,
		Channel:   thc.Channel,
		CreatedAt: thc.CreatedAt.Format(time.RFC3339),
		UpdatedAt: thc.UpdatedAt.Format(time.RFC3339),
		Chat:      thc.Chat,
	}
	return json.Marshal(aux)
}

// type ThChatRespPayload struct {
// 	ThreadChatId string
// 	Sequence     int
// 	Status       string
// 	Read         bool
// 	Replied      bool
// 	Priority     string
// 	Customer     ThCustomerRespPayload
// 	Assignee     *ThMemberRespPayload
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// 	Messages     []ThChatMessageRespPayload
// }

// func (thresp ThChatRespPayload) MarshalJSON() ([]byte, error) {
// 	var assignee *ThMemberRespPayload

// 	if thresp.Assignee != nil {
// 		assignee = thresp.Assignee
// 	}

// 	aux := &struct {
// 		ThreadChatId string                     `json:"threadChatId"`
// 		Sequence     int                        `json:"sequence"`
// 		Status       string                     `json:"status"`
// 		Read         bool                       `json:"read"`
// 		Replied      bool                       `json:"replied"`
// 		Priority     string                     `json:"priority"`
// 		Customer     ThCustomerRespPayload      `json:"customer"`
// 		Assignee     *ThMemberRespPayload       `json:"assignee"`
// 		CreatedAt    string                     `json:"createdAt"`
// 		UpdatedAt    string                     `json:"updatedAt"`
// 		Messages     []ThChatMessageRespPayload `json:"messages"`
// 	}{
// 		ThreadChatId: thresp.ThreadChatId,
// 		Sequence:     thresp.Sequence,
// 		Status:       thresp.Status,
// 		Read:         thresp.Read,
// 		Replied:      thresp.Replied,
// 		Priority:     thresp.Priority,
// 		Customer:     thresp.Customer,
// 		Assignee:     assignee,
// 		CreatedAt:    thresp.CreatedAt.Format(time.RFC3339),
// 		UpdatedAt:    thresp.UpdatedAt.Format(time.RFC3339),
// 		Messages:     thresp.Messages,
// 	}
// 	return json.Marshal(aux)
// }

type AddCustomerIdentitiesReqPayload struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	External *string `json:"external"`
}

type AddCustomerIdentitiesRespPayload struct {
	IsVerified bool           `json:"isVerified"`
	Name       string         `json:"name"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	ExternalId sql.NullString `json:"externalId"`
}

func (w AddCustomerIdentitiesRespPayload) MarshalJSON() ([]byte, error) {
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
