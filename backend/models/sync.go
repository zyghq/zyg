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
