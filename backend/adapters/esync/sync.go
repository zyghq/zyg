package esync

import (
	"context"
	"fmt"
	"github.com/zyghq/zyg/models"
)

func (sy *SyncDB) SaveWorkspace(
	ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error) {
	//hub := sentry.GetHubFromContext(ctx)
	fmt.Println(" *** Saving workspace *** ")
	fmt.Println(workspace)
	return models.WorkspaceInSync{}, nil
}
