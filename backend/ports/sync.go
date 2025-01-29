package ports

import (
	"context"
	"github.com/zyghq/zyg/models"
)

type SyncServicer interface {
	SyncWorkspaceRPC(ctx context.Context, workspace models.Workspace) (models.WorkspaceInSync, error)
	SyncCustomerRPC(ctx context.Context, customer models.Customer) (models.CustomerInSync, error)
	SyncMemberRPC(ctx context.Context, member models.Member) (models.MemberInSync, error)
	SyncThread(ctx context.Context, thread models.Thread, labels []models.ThreadLabel) (models.ThreadInSync, error)
}
