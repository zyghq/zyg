package models

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

// AuthJWTClaims taken from Supabase JWT encoding
type AuthJWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// CustomerJWTClaims custom jwt claims for customer
type CustomerJWTClaims struct {
	WorkspaceId     string  `json:"workspaceId"`
	ExternalId      *string `json:"externalId"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	IsEmailVerified bool    `json:"isEmailVerified"`
	jwt.RegisteredClaims
}

type KycMailJWTClaims struct {
	WorkspaceId string `json:"workspaceId"`
	Email       string `json:"email"`
	RedirectUrl string `json:"redirectUrl"`
	jwt.RegisteredClaims
}

// NullString custom data type wrapper for SQL nullable string
func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// MemberActor identifies referenced Member.
type MemberActor struct {
	MemberId string
	Name     string
}

// AssignedMember represents the Member assigned with when it was assigned.
type AssignedMember struct {
	MemberId   string
	Name       string
	AssignedAt time.Time // The datetime the Member was assigned.
}

// CustomerActor identifies referenced Customer.
type CustomerActor struct {
	CustomerId string
	Name       string
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

func (w Workspace) NewWorkspace(accountId string, name string) Workspace {
	return Workspace{
		WorkspaceId: w.GenId(),
		AccountId:   accountId,
		Name:        name,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
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

// ThreadMetrics represents Thread count metrics for specific status and stage.
type ThreadMetrics struct {
	ActiveCount             int // sum of all threads in status Todo.
	NeedsFirstResponseCount int // sum of threads in status Todo and stage NeedsFirstResponse.
	WaitingOnCustomerCount  int // sum of threads in status Todo and stage WaitingOnCustomer.
	HoldCount               int // sum of threads in status Todo and stage Hold.
	NeedsNextResponseCount  int // sum of threads in status Todo and stage NeedsNextResponse.
}

// ThreadAssigneeMetrics represents Thread count metrics for member
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

type WorkspaceSecret struct {
	WorkspaceId string
	Hmac        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (sk WorkspaceSecret) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		Hmac        string `json:"hmac"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: sk.WorkspaceId,
		Hmac:        sk.Hmac,
		CreatedAt:   sk.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sk.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

type ClaimedMail struct {
	ClaimId      string
	WorkspaceId  string
	CustomerId   string
	Email        string
	HasConflict  bool
	ExpiresAt    time.Time
	Token        string
	IsMailSent   bool
	Platform     sql.NullString
	SenderId     sql.NullString
	SenderStatus sql.NullString
	SentAt       sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (cl ClaimedMail) GenId() string {
	return "cl" + xid.New().String()
}

func (cl ClaimedMail) NewVerification(
	workspaceId string, customerId string, email string,
	hasConflict bool, expiresAt time.Time, token string,
) ClaimedMail {
	return ClaimedMail{
		ClaimId:      cl.GenId(),
		WorkspaceId:  workspaceId,
		CustomerId:   customerId,
		Email:        email,
		HasConflict:  hasConflict,
		ExpiresAt:    expiresAt,
		Token:        token,
		IsMailSent:   false,
		Platform:     sql.NullString{},
		SenderId:     sql.NullString{},
		SenderStatus: sql.NullString{},
		SentAt:       sql.NullTime{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type WidgetSessionData struct {
	WorkspaceId  string `json:"workspaceId"`
	CustomerId   string `json:"customerId"`
	IdentityHash string `json:"identityHash"`
}

type WidgetSession struct {
	SessionId string
	WidgetId  string
	Data      string
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ws *WidgetSession) GenId() string {
	return "ws" + xid.New().String()
}

func (ws *WidgetSession) CreateSession(sessionId string, widgetId string) WidgetSession {
	return WidgetSession{
		SessionId: sessionId,
		WidgetId:  widgetId,
		ExpireAt:  time.Now().UTC().Add(time.Hour * 24), // 24 hours
	}
}

// CreateSessionData creates the widget session data
func (ws *WidgetSession) CreateSessionData(
	workspaceId string, customerId string, identityHash string) WidgetSessionData {
	return WidgetSessionData{
		WorkspaceId:  workspaceId,
		CustomerId:   customerId,
		IdentityHash: identityHash,
	}
}

// Encode encodes the session data and returns the encoded data
// The secret key is specific to each workspace and is used to verify the integrity of the session data.
func (ws *WidgetSession) Encode(sk string, session WidgetSessionData) (string, error) {
	// jsonify the session data
	sessJson, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	// stringify the JSON data
	sessionData := string(sessJson)
	encodedSessionData := base64.URLEncoding.EncodeToString([]byte(sessionData))

	// create HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(encodedSessionData))
	signature := h.Sum(nil)

	// Base64 encode the signature
	encodedSignature := base64.URLEncoding.EncodeToString(signature)

	// combine the session data and the signature
	signedSession := fmt.Sprintf("%s:%s", encodedSessionData, encodedSignature)
	return signedSession, nil
}

// SetEncodeData encodes the session data and sets the encoded data to the session
// The secret key is specific to each workspace and is used to verify the integrity of the session data.
func (ws *WidgetSession) SetEncodeData(sk string, session WidgetSessionData) error {
	encoded, err := ws.Encode(sk, session)
	if err != nil {
		return err
	}
	ws.Data = encoded
	return nil
}

// Decode splits the signed session string, verifies it, and decodes the session data
// The secret key is specific to each workspace and is used to verify the integrity of the session data.
func (ws *WidgetSession) Decode(sk string) (WidgetSessionData, error) {
	var session WidgetSessionData

	parts := strings.Split(ws.Data, ":")
	if len(parts) != 2 {
		return WidgetSessionData{}, errors.New("invalid signed session format")
	}

	encodedSessionData := parts[0]
	signature := parts[1]

	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(encodedSessionData))
	expectedSignature := h.Sum(nil)

	receivedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return WidgetSessionData{}, err
	}

	if !hmac.Equal(receivedSignature, expectedSignature) {
		return WidgetSessionData{}, errors.New(
			"invalid signature: perhaps the data is tampered or the secret keys updated")
	}

	sessionData, err := base64.URLEncoding.DecodeString(encodedSessionData)
	if err != nil {
		return WidgetSessionData{}, err
	}
	err = json.Unmarshal(sessionData, &session)
	if err != nil {
		return WidgetSessionData{}, err
	}
	return session, nil
}
