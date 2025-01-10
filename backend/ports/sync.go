package ports

import (
	"context"
	"github.com/zyghq/zyg/models"
)

type Syncer interface {
	SaveWorkspace(ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error)
}

type SyncServicer interface {
	SyncWorkspace(ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error)
}
