package esync

import (
	"context"
	"github.com/cristalhq/builq"
	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg/models"
	"log/slog"
)

func (sy *SyncDB) SaveWorkspace(
	ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error) {
	var inSync models.WorkspaceInSync
	hub := sentry.GetHubFromContext(ctx)

	q := builq.New()
	q(`INSERT INTO workspace (workspace_id, name, public_name, created_at, updated_at, version_id, synced_at)`)
	q(`VALUES (%$, %$, %$, %$, %$, %$, %$)`,
		workspace.WorkspaceID, workspace.Name, workspace.PublicName,
		workspace.CreatedAt, workspace.UpdatedAt, workspace.VersionID, workspace.SyncedAt)
	q(`ON CONFLICT (workspace_id) DO UPDATE SET`)
	q(`name = EXCLUDED.name,`)
	q(`public_name = EXCLUDED.public_name,`)
	q(`created_at = EXCLUDED.created_at,`)
	q(`updated_at = EXCLUDED.updated_at,`)
	q(`version_id = EXCLUDED.version_id,`)
	q(`synced_at = EXCLUDED.synced_at`)
	q(`RETURNING workspace_id, synced_at, version_id`)

	stmt, _, err := q.Build()
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to build query", slog.Any("err", err))
		return models.WorkspaceInSync{}, err
	}

	err = sy.db.QueryRow(
		ctx, stmt, workspace.WorkspaceID, workspace.Name, workspace.PublicName,
		workspace.CreatedAt, workspace.UpdatedAt, workspace.VersionID, workspace.SyncedAt,
	).Scan(&inSync.WorkspaceID, &inSync.SyncedAt, &inSync.VersionID)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.WorkspaceInSync{}, err
	}
	return inSync, nil
}
