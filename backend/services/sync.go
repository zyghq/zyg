package services

import (
	"context"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type SyncService struct {
	syncDB ports.Syncer
}

func NewSyncService(syncDB ports.Syncer) *SyncService {
	return &SyncService{syncDB: syncDB}
}

// SyncWorkspace - Todo: implement this
func (sy *SyncService) SyncWorkspace(ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error) {
	inSync, err := sy.syncDB.SaveWorkspace(ctx, workspace)
	if err != nil {
		return models.WorkspaceInSync{}, err
	}
	return inSync, nil
}
