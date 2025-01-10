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
	q(`INSERT INTO workspace (
		"workspaceId", "name", "publicName", "createdAt", "updatedAt",
		"versionId", "syncedAt"
	) VALUES (%$, %$, %$, %$, %$, %$, %$)`,
		workspace.WorkspaceID, workspace.Name, workspace.PublicName,
		workspace.CreatedAt, workspace.UpdatedAt, workspace.VersionID, workspace.SyncedAt,
	)
	q(`ON CONFLICT ("workspaceId") DO UPDATE SET`)
	q(`"name" = EXCLUDED."name",`)
	q(`"publicName" = EXCLUDED."publicName",`)
	q(`"createdAt" = EXCLUDED."createdAt",`)
	q(`"updatedAt" = EXCLUDED."updatedAt",`)
	q(`"versionId" = EXCLUDED."versionId",`)
	q(`"syncedAt" = EXCLUDED."syncedAt"`)
	q(`RETURNING "workspaceId", "syncedAt", "versionId"`)

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
