package domain

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// taken from Supabase JWT encoding
type AuthJWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// custom jwt claims for customer
type CustomerJWTClaims struct {
	WorkspaceId string `json:"workspaceId"`
	ExternalId  string `json:"externalId"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	jwt.RegisteredClaims
}

// custom data type wrapper for SQL nullable string
func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

type ThreadStatus struct{}

func (s ThreadStatus) Todo() string {
	return "todo"
}

func (s ThreadStatus) Done() string {
	return "done"
}

func (s ThreadStatus) InProgress() string {
	return "inprogress"
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

func (a LabelAddedBy) AiAssist() string {
	return "aiassist"
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
	return "a_" + xid.New().String()
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
	Description sql.NullString
	UnMask      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ap AccountPAT) GenId() string {
	return "ap_" + xid.New().String()
}

func (ap AccountPAT) MarshalJSON() ([]byte, error) {
	var description *string
	var token string
	if ap.Description.Valid {
		description = &ap.Description.String
	}

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
		AccountId   string  `json:"accountId"`
		PatId       string  `json:"patId"`
		Token       string  `json:"token"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		AccountId:   ap.AccountId,
		PatId:       ap.PatId,
		Token:       token,
		Name:        ap.Name,
		Description: description,
		CreatedAt:   ap.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ap.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
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

type Member struct {
	WorkspaceId string
	AccountId   string
	MemberId    string
	Name        sql.NullString
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Member) GenId() string {
	return "m_" + xid.New().String()
}

func (m Member) MarshalJSON() ([]byte, error) {
	var name *string
	if m.Name.Valid {
		name = &m.Name.String
	}
	aux := &struct {
		WorkspaceId string  `json:"workspaceId"`
		AccountId   string  `json:"accountId"`
		MemberId    string  `json:"memberId"`
		Name        *string `json:"name"`
		Role        string  `json:"role"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: m.WorkspaceId,
		AccountId:   m.AccountId,
		MemberId:    m.MemberId,
		Name:        name,
		Role:        m.Role,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type Customer struct {
	WorkspaceId string
	CustomerId  string
	ExternalId  sql.NullString
	Email       sql.NullString
	Phone       sql.NullString
	Name        sql.NullString
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func (c Customer) GenId() string {
	return "c_" + xid.New().String()
}

func (c Customer) MarshalJSON() ([]byte, error) {
	var externalId, email, phone, name *string
	if c.ExternalId.Valid {
		externalId = &c.ExternalId.String
	}
	if c.Email.Valid {
		email = &c.Email.String
	}
	if c.Phone.Valid {
		phone = &c.Phone.String
	}
	if c.Name.Valid {
		name = &c.Name.String
	}

	aux := &struct {
		WorkspaceId string  `json:"workspaceId"`
		CustomerId  string  `json:"customerId"`
		ExternalId  *string `json:"externalId"`
		Email       *string `json:"email"`
		Phone       *string `json:"phone"`
		Name        *string `json:"name"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: c.WorkspaceId,
		CustomerId:  c.CustomerId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		Name:        name,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ThreadChat struct {
	WorkspaceId  string
	CustomerId   string
	CustomerName sql.NullString
	AssigneeId   sql.NullString
	AssigneeName sql.NullString
	ThreadChatId string
	Title        string
	Summary      string
	Sequence     int
	Status       string
	Read         bool
	Replied      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (th ThreadChat) GenId() string {
	return "th_" + xid.New().String()
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

type ThreadChatMessage struct {
	ThreadChatId        string
	ThreadChatMessageId string
	Body                string
	Sequence            int
	CustomerId          sql.NullString
	CustomerName        sql.NullString
	MemberId            sql.NullString
	MemberName          sql.NullString
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (thm ThreadChatMessage) GenId() string {
	return "thm_" + xid.New().String()
}

type ThreadChatWithMessage struct {
	ThreadChat ThreadChat
	Message    ThreadChatMessage
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
	return "l_" + xid.New().String()
}

type ThreadChatLabel struct {
	ThreadChatId      string
	LabelId           string
	ThreadChatLabelId string
	AddedBy           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (thl ThreadChatLabel) GenId() string {
	return xid.New().String()
}

type ThreadChatLabelled struct {
	ThreadChatLabelId string
	ThreadChatId      string
	LabelId           string
	Name              string
	Icon              string
	AddedBy           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
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
