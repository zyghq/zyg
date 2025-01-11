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
