package ports

import (
	"context"
	"github.com/zyghq/zyg/models"
)

type SyncServicer interface {
	SyncWorkspaceRPC(ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error)
	SyncCustomer(ctx context.Context, customer models.CustomerShape) (models.CustomerInSync, error)
	SyncMember(ctx context.Context, member models.MemberShape) (models.MemberInSync, error)
	SyncThread(ctx context.Context, thread models.ThreadShape) (models.ThreadInSync, error)
}
