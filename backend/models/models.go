package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

func GenToken(length int, prefix string) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return prefix + base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

// AuthJWTClaims taken from Supabase JWT encoding
type AuthJWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// CustomerJWTClaims custom jwt claims for customer
type CustomerJWTClaims struct {
	WorkspaceId string `json:"workspaceId"`
	ExternalId  string `json:"externalId"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	jwt.RegisteredClaims
}

// NullString custom data type wrapper for SQL nullable string
func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// IsValidUUID validates if a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type Workspace struct {
	WorkspaceId string
	AccountId   string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (w Workspace) GenId() string {
	return "wrk" + xid.New().String()
}

func (w Workspace) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		AccountId   string `json:"accountId"`
		Name        string `json:"name"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: w.WorkspaceId,
		AccountId:   w.AccountId,
		Name:        w.Name,
		CreatedAt:   w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   w.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ThreadStatus struct{}

func (s ThreadStatus) Todo() string {
	return "todo"
}

func (s ThreadStatus) Done() string {
	return "done"
}

func (s ThreadStatus) Snoozed() string {
	return "snoozed"
}

func (s ThreadStatus) UnSnoozed() string {
	return "unsnoozed"
}

func (s ThreadStatus) DefaultStatus() string {
	return s.Todo()
}

func (s ThreadStatus) IsValid(status string) bool {
	switch status {
	case s.Done():
		return true
	case s.Todo():
		return true
	default:
		return false
	}
}

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

func (p ThreadPriority) DefaultPriority() string {
	return p.Normal()
}

func (p ThreadPriority) IsValid(s string) bool {
	switch s {
	case p.Urgent(), p.High(), p.Normal(), p.Low():
		return true
	default:
		return false
	}
}

type ThreadChannel struct{}

func (c ThreadChannel) Chat() string {
	return "chat"
}

type LabelAddedBy struct{}

func (a LabelAddedBy) User() string {
	return "user"
}

func (a LabelAddedBy) System() string {
	return "system"
}

func (a LabelAddedBy) Ai() string {
	return "ai"
}

type Account struct {
	AccountId  string
	Email      string
	Provider   string
	AuthUserId string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (a Account) GenId() string {
	return "ac" + xid.New().String()
}

func (a Account) MarshalJSON() ([]byte, error) {
	aux := &struct {
		AccountId string `json:"accountId"`
		Email     string `json:"email"`
		Provider  string `json:"provider"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		AccountId: a.AccountId,
		Email:     a.Email,
		Provider:  a.Provider,
		Name:      a.Name,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type AccountPAT struct {
	AccountId   string
	PatId       string
	Token       string
	Name        string
	Description string
	UnMask      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ap AccountPAT) GenId() string {
	return "ap_" + xid.New().String()
}

func (ap AccountPAT) MarshalJSON() ([]byte, error) {
	var token string
	maskLeft := func(s string) string {
		rs := []rune(s)
		for i := range rs[:len(rs)-4] {
			rs[i] = '*'
		}
		return string(rs)
	}

	if !ap.UnMask {
		token = maskLeft(ap.Token)
	} else {
		token = ap.Token
	}

	aux := &struct {
		AccountId   string `json:"accountId"`
		PatId       string `json:"patId"`
		Token       string `json:"token"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		AccountId:   ap.AccountId,
		PatId:       ap.PatId,
		Token:       token,
		Name:        ap.Name,
		Description: ap.Description,
		CreatedAt:   ap.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ap.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type Member struct {
	WorkspaceId string
	AccountId   string
	MemberId    string
	Name        string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Member) GenId() string {
	return "mm" + xid.New().String()
}

func (m Member) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		AccountId   string `json:"accountId"`
		MemberId    string `json:"memberId"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: m.WorkspaceId,
		AccountId:   m.AccountId,
		MemberId:    m.MemberId,
		Name:        m.Name,
		Role:        m.Role,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type MemberRole struct{}

func (mr MemberRole) Primary() string {
	return "primary"
}

func (mr MemberRole) Owner() string {
	return "owner"
}

func (mr MemberRole) Admin() string {
	return "administrator"
}

func (mr MemberRole) Member() string {
	return "member"
}

func (mr MemberRole) DefaultRole() string {
	return mr.Member()
}

func (mr MemberRole) IsValid(s string) bool {
	switch s {
	case mr.Primary(), mr.Owner(), mr.Admin(), mr.Member():
		return true
	default:
		return false
	}
}

type Customer struct {
	WorkspaceId string
	CustomerId  string
	ExternalId  sql.NullString
	Email       sql.NullString
	Phone       sql.NullString
	Name        string
	AnonId      string
	IsVerified  bool
	Role        string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func (c Customer) GenId() string {
	return "cs" + xid.New().String()
}

func (c Customer) Visitor() string {
	return "visitor"
}

func (c Customer) Lead() string {
	return "lead"
}

func (c Customer) Engaged() string {
	return "engaged"
}

func (c Customer) MarshalJSON() ([]byte, error) {
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
		WorkspaceId string  `json:"workspaceId"`
		CustomerId  string  `json:"customerId"`
		ExternalId  *string `json:"externalId"`
		Email       *string `json:"email"`
		Phone       *string `json:"phone"`
		Name        string  `json:"name"`
		AnonId      string  `json:"anonId"`
		IsVerified  bool    `json:"isVerified"`
		Role        string  `json:"role"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: c.WorkspaceId,
		CustomerId:  c.CustomerId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		Name:        c.Name,
		AnonId:      c.AnonId,
		IsVerified:  c.IsVerified,
		Role:        c.Role,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (c Customer) AnonName() string {
	return "Anon User"
}

func (c Customer) AddAnonimizedEmail(email string) string {
	return c.AnonId + "~" + email
}

func (c Customer) AddAnonimizedPhone(phone string) string {
	return c.AnonId + "~" + phone
}

func (c Customer) AddAnonimizedExternalId(externalId string) string {
	return c.AnonId + "~" + externalId
}

func (c Customer) DeAnonEmail() string {
	splits := strings.Split(c.Email.String, "~")
	if len(splits) == 2 {
		return splits[1]
	}
	return c.Email.String
}

func (c Customer) DeAnonPhone() string {
	splits := strings.Split(c.Phone.String, "~")
	if len(splits) == 2 {
		return splits[1]
	}
	return c.Phone.String
}

func (c Customer) DeAnonExternalId() string {
	splits := strings.Split(c.ExternalId.String, "~")
	if len(splits) == 2 {
		return splits[1]
	}
	return c.ExternalId.String
}

type IngressMessage struct {
	MessageId  string
	CustomerId string
	FirstSeq   int
	LastSeq    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (im IngressMessage) GenId() string {
	return "im" + xid.New().String()
}

type EgressMessage struct {
	MessageId string
	MemberId  string
	FirstSeq  int
	LastSeq   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (em EgressMessage) GenId() string {
	return "em" + xid.New().String()
}

type Thread struct {
	WorkspaceId         string
	ThreadId            string
	CustomerId          string
	CustomerName        string
	AssigneeId          sql.NullString
	AssigneeName        sql.NullString
	Title               string
	Description         string
	Sequence            int
	Status              string
	Read                bool
	Replied             bool
	Priority            string
	Spam                bool
	Channel             string
	PreviewText         string
	IngressMessageId    sql.NullString
	IngressFirstSeq     sql.NullInt64
	IngressLastSeq      sql.NullInt64
	IngressCustomerId   sql.NullString
	IngressCustomerName sql.NullString
	EgressMessageId     sql.NullString
	EgressFirstSeq      sql.NullInt64
	EgressLastSeq       sql.NullInt64
	EgressMemberId      sql.NullString
	EgressMemberName    sql.NullString
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (t Thread) GenId() string {
	return "th" + xid.New().String()
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

func (l Label) MarshalJSON() ([]byte, error) {
	aux := &struct {
		LabelId     string `json:"labelId"`
		WorkspaceId string `json:"workspaceId"`
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		LabelId:     l.LabelId,
		WorkspaceId: l.WorkspaceId,
		Name:        l.Name,
		Icon:        l.Icon,
		CreatedAt:   l.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   l.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ThreadQAA struct {
	WorkspaceId string
	ThreadQAId  string
	AnswerId    string
	Answer      string
	Sequence    int
	Eval        sql.NullInt32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ThreadQA struct {
	WorkspaceId    string
	CustomerId     string
	ThreadId       string
	ParentThreadId sql.NullString
	Query          string
	Title          string
	Summary        string
	Sequence       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Label struct {
	WorkspaceId string
	LabelId     string
	Name        string
	Icon        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (l Label) GenId() string {
	return "lb" + xid.New().String()
}

type ThreadLabel struct {
	ThreadLabelId string
	ThreadId      string
	LabelId       string
	Name          string
	Icon          string
	AddedBy       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (thl ThreadLabel) GenId() string {
	return xid.New().String()
}

type ThreadLabelMetric struct {
	LabelId string
	Name    string
	Icon    string
	Count   int
}

type ThreadMetrics struct {
	ActiveCount  int // sum of threads in Todo and Snoozed
	DoneCount    int
	TodoCount    int
	SnoozedCount int
}

type ThreadAssigneeMetrics struct {
	MeCount            int
	UnAssignedCount    int
	OtherAssignedCount int
}

type ThreadMemberMetrics struct {
	ThreadMetrics
	ThreadAssigneeMetrics
	ThreadLabelMetrics []ThreadLabelMetric
}

type Widget struct {
	WorkspaceId   string
	WidgetId      string
	Name          string
	Configuration map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (w Widget) GenId() string {
	return "wd" + xid.New().String()
}

type SecretKey struct {
	WorkspaceId string
	SecretKey   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (sk SecretKey) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		SecretKey   string `json:"secretKey"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: sk.WorkspaceId,
		SecretKey:   sk.SecretKey,
		CreatedAt:   sk.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sk.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}
