package models

import (
	"database/sql"
	"github.com/rs/xid"
	"time"
)

const (
	notificationStatusIgnored = "ignored"
	notificationStatusSent    = "sent"
	notificationStatusAborted = "aborted"
	notificationStatusSending = "sending"
	eventSeverityInfo         = "info"
	eventSeverityWarning      = "warning"
	eventSeverityError        = "error"
	eventSeverityCritical     = "critical"
)

// CustomerEvent represents an event related to a customer, including details like severity, timestamps, and status.
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
	now := time.Now()
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
