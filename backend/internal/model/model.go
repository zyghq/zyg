package model

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type mErr string

func (err mErr) Error() string {
	return string(err)
}

const (
	ErrEmpty = mErr("nothing found")
	ErrQuery = mErr("failed to query")
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

type CustomerJWTClaims struct {
	WorkspaceId string `json:"workspaceId"`
	ExternalId  string `json:"externalId"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	jwt.RegisteredClaims
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

type Workspace struct {
	WorkspaceId string
	AccountId   string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

type ThreadChatLabel struct {
	ThreadChatId      string
	LabelId           string
	ThreadChatLabelId string
	AddedBy           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
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
