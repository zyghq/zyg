package esync

import (
	"context"
	"fmt"
	"github.com/zyghq/zyg/models"
)

func (sy *SyncDB) SaveWorkspace(
	ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error) {
	return models.WorkspaceInSync{}, fmt.Errorf("not implemented")
}
