package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/rs/xid"
	"github.com/sanchitrk/namingo"
	"github.com/zyghq/zyg"
	"strings"
	"time"
)

const (
	// notificationStatusIgnored indicates that a notification has been ignored
	// and will not be sent or processed further.
	notificationStatusIgnored = "ignored"

	// notificationStatusSent indicates that the notification has been successfully sent.
	notificationStatusSent = "sent"

	// notificationStatusAborted indicates that the notification process has been
	// terminated and will not be completed.
	notificationStatusAborted = "aborted"

	// notificationStatusSending indicates that the notification is currently in the process of being sent.
	notificationStatusSending = "sending"

	// eventSeverityInfo represents an informational event severity level.
	eventSeverityInfo = "info"

	// eventSeverityWarning indicates a warning severity level for an event.
	eventSeverityWarning = "warning"

	// eventSeverityError represents an error severity level for customer events.
	eventSeverityError = "error"

	// eventSeverityCritical indicates a critical severity level for an event.
	eventSeverityCritical = "critical"
)

type Customer struct {
	WorkspaceId     string
	CustomerId      string
	ExternalId      sql.NullString
	Email           sql.NullString
	Phone           sql.NullString
	Name            string
	IsEmailVerified bool
	Role            string
	UpdatedAt       time.Time
	CreatedAt       time.Time
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

func (c Customer) IsVisitor() bool {
	return c.Role == c.Visitor()
}

func (c Customer) AnonName() string {
	return namingo.Generate(2, " ", namingo.TitleCase())
}

func (c Customer) AvatarUrl() string {
	url := zyg.GetAvatarBaseURL()
	// url may or may not have a trailing slash
	// add a trailing slash if it doesn't have one
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url + c.CustomerId
}

func (c Customer) AsCustomerActor() CustomerActor {
	return CustomerActor{
		CustomerId: c.CustomerId,
		Name:       c.Name,
	}
}

// IdentityHash is a hash of the customer's identity
// Combined these fields create a unique hash for the customer
// (XXX): You might have to update this if you plan to add more identity fields
func (c Customer) IdentityHash() string {
	h := sha256.New()
	// Combine all fields into a single string
	identityString := fmt.Sprintf("%s:%s:%s:%s:%s:%t",
		c.WorkspaceId,
		c.CustomerId,
		c.ExternalId.String,
		c.Email.String,
		c.Phone.String,
		c.IsEmailVerified,
	)

	// Write the combined string to the hash
	h.Write([]byte(identityString))

	// Return the hash as a base64 encoded string
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c Customer) MakeCopy() Customer {
	return Customer{
		WorkspaceId:     c.WorkspaceId,
		CustomerId:      c.CustomerId,
		ExternalId:      c.ExternalId,
		Email:           c.Email,
		Phone:           c.Phone,
		Name:            c.Name,
		IsEmailVerified: c.IsEmailVerified,
		Role:            c.Role,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

// CustomerEvent represents an event related to a customer,
// including details like severity, timestamps, and status.
type CustomerEvent struct {
	EventId            string
	CustomerId         string
	ThreadId           sql.NullString
	Event              string
	EventBody          string
	Severity           string
	EventTimestamp     time.Time
	IdempotencyKey     sql.NullString
	NotificationStatus string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (cv *CustomerEvent) GenID() string {
	return "ev" + xid.New().String()
}

// NewEvent creates and returns a new CustomerEvent instance with specified details and default timestamps.
func (cv *CustomerEvent) NewEvent(
	customerId string, event string, body string, severity string) CustomerEvent {
	now := time.Now().UTC()
	c := CustomerEvent{
		EventId:        cv.GenID(),
		CustomerId:     customerId,
		Event:          event,
		EventBody:      body,
		EventTimestamp: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	c.SetSeverity(severity)
	c.NotificationIgnored()
	return c
}

func (cv *CustomerEvent) SetTimestampFromStr(ts string) error {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}
	cv.EventTimestamp = t
	return nil
}

func (cv *CustomerEvent) SetTimestampFromTime(t time.Time) {
	cv.EventTimestamp = t
}

func (cv *CustomerEvent) NotificationIgnored() {
	cv.NotificationStatus = notificationStatusIgnored
}

func (cv *CustomerEvent) NotificationAborted() {
	cv.NotificationStatus = notificationStatusAborted
}

func (cv *CustomerEvent) NotificationSending() {
	cv.NotificationStatus = notificationStatusSending
}

func (cv *CustomerEvent) NotificationSent() {
	cv.NotificationStatus = notificationStatusSent
}

func (cv *CustomerEvent) SetThreadId(threadId string) {
	cv.ThreadId = NullString(&threadId)
}

func (cv *CustomerEvent) SetIdempotencyKey(idempotencyKey string) {
	cv.IdempotencyKey = NullString(&idempotencyKey)
}

func (cv *CustomerEvent) SeverityInfo() {
	cv.Severity = eventSeverityInfo
}

func (cv *CustomerEvent) SeverityWarning() {
	cv.Severity = eventSeverityWarning
}

func (cv *CustomerEvent) SeverityError() {
	cv.Severity = eventSeverityError
}

func (cv *CustomerEvent) SeverityCritical() {
	cv.Severity = eventSeverityCritical
}

// IsSeverityValid checks if the given severity string is one of the
// predefined valid severity levels for a CustomerEvent.
func (cv *CustomerEvent) IsSeverityValid(s string) bool {
	return s == eventSeverityInfo ||
		s == eventSeverityWarning ||
		s == eventSeverityError ||
		s == eventSeverityCritical
}

// SetSeverity sets the severity level of the CustomerEvent,
// defaulting to "info" if the provided level is invalid.
func (cv *CustomerEvent) SetSeverity(s string) {
	if cv.IsSeverityValid(s) {
		cv.Severity = s
	} else {
		cv.SeverityInfo()
	}
}
