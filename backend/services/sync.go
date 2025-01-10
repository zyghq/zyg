package services

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type SyncService struct {
	syncDB ports.Syncer
}

func NewSyncService(syncDB ports.Syncer) *SyncService {
	return &SyncService{syncDB: syncDB}
}

// SyncWorkspace synchronizes the provided workspace data with the database,
// returning the synchronized workspace details.
func (sy *SyncService) SyncWorkspace(
	ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error) {
	hub := sentry.GetHubFromContext(ctx)
	inSync, err := sy.syncDB.SaveWorkspace(ctx, workspace)
	if err != nil {
		hub.CaptureException(err)
		return models.WorkspaceInSync{}, err
	}
	return inSync, nil
}
