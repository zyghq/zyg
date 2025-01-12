package models

import "time"

type WorkspaceShape struct {
	WorkspaceID string    `json:"workspaceId"`
	Name        string    `json:"name"`
	PublicName  string    `json:"publicName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	VersionID   string    `json:"versionId"`
	SyncedAt    time.Time `json:"syncedAt"`
}

type WorkspaceInSync struct {
	WorkspaceID string    `json:"workspaceId"`
	SyncedAt    time.Time `json:"syncedAt"`
	VersionID   string    `json:"versionId"`
}

type CustomerShape struct {
	CustomerID      string    `json:"customerId"`
	WorkspaceID     string    `json:"workspaceId"`
	ExternalID      *string   `json:"externalId"`
	Email           *string   `json:"email"`
	Phone           *string   `json:"phone"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
	AvatarURL       string    `json:"avatarUrl"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	VersionID       string    `json:"versionId"`
	SyncedAt        time.Time `json:"syncedAt"`
}

type CustomerInSync struct {
	CustomerID string    `json:"customerId"`
	SyncedAt   time.Time `json:"syncedAt"`
	VersionID  string    `json:"versionId"`
}

type MemberShape struct {
	MemberID    string                 `json:"memberId"`
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	PublicName  string                 `json:"publicName"`
	Role        string                 `json:"role"`
	Permissions map[string]interface{} `json:"permissions"`
	AvatarURL   string                 `json:"avatarUrl"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	VersionID   string                 `json:"versionId"`
	SyncedAt    time.Time              `json:"syncedAt"`
}

type MemberInSync struct {
	MemberID  string    `json:"memberId"`
	SyncedAt  time.Time `json:"syncedAt"`
	VersionID string    `json:"versionId"`
}

type ThreadShape struct {
	ThreadID          string                 `json:"threadId"`
	WorkspaceID       string                 `json:"workspaceId"`
	CustomerID        string                 `json:"customerId"`
	AssigneeID        *string                `json:"assigneeId"`
	AssignedAt        *time.Time             `json:"assignedAt"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	PreviewText       string                 `json:"previewText"`
	Status            string                 `json:"status"`
	StatusChangedAt   time.Time              `json:"statusChangedAt"`
	StatusChangedByID string                 `json:"statusChangedById"`
	Stage             string                 `json:"stage"`
	Replied           bool                   `json:"replied"`
	Priority          string                 `json:"priority"`
	Channel           string                 `json:"channel"`
	CreatedByID       string                 `json:"createdById"`
	UpdatedByID       string                 `json:"updatedById"`
	Labels            map[string]interface{} `json:"labels"`
	InboundSeqID      *string                `json:"inboundSeqId"`
	OutboundSeqID     *string                `json:"outboundSeqId"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	VersionID         string                 `json:"versionId"`
	SyncedAt          time.Time              `json:"syncedAt"`
}

type ThreadInSync struct {
	ThreadID  string    `json:"threadId"`
	SyncedAt  time.Time `json:"syncedAt"`
	VersionID string    `json:"versionId"`
}
