package models

import (
	"github.com/rs/xid"
	"time"
)

const (
	ActivityThreadMessage = "thread.message"
)

type Activity struct {
	ActivityID   string                 `json:"activityId"`
	ThreadID     string                 `json:"threadId"`
	ActivityType string                 `json:"activityType"`
	Customer     *CustomerActor         `json:"customer,omitempty"`
	Member       *MemberActor           `json:"member,omitempty"`
	Body         map[string]interface{} `json:"body"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

func (a *Activity) GenID() string {
	return "act" + xid.New().String()
}

type ActivityOption func(activity *Activity)

func NewActivity(threadID string, activityType string, opts ...ActivityOption) *Activity {
	activityID := (&Activity{}).GenID()
	now := time.Now()
	activity := &Activity{
		ThreadID:     threadID,
		ActivityID:   activityID,
		ActivityType: activityType,
		Body:         make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	for _, opt := range opts {
		opt(activity)
	}
	return activity
}

func SetActivityCustomer(customer CustomerActor) ActivityOption {
	return func(activity *Activity) {
		activity.Customer = &customer
		activity.Member = nil
	}
}

func SetActivityMember(member MemberActor) ActivityOption {
	return func(activity *Activity) {
		activity.Member = &member
		activity.Customer = nil
	}
}

func SetActivityBody(body map[string]interface{}) ActivityOption {
	return func(activity *Activity) {
		activity.Body = body
	}
}
