package ports

import (
	"context"
	"github.com/zyghq/zyg/models"
)

type SyncServicer interface {
	SyncWorkspace(ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error)
	SyncCustomer(ctx context.Context, customer models.CustomerShape) (models.CustomerInSync, error)
	SyncMember(ctx context.Context, member models.MemberShape) (models.MemberInSync, error)
	SyncThread(ctx context.Context, thread models.ThreadShape) (models.ThreadInSync, error)
}

// Syncer defines methods for synchronizing entities and values between application and sync systems.
// Think it as sync repository
type Syncer interface {
	SaveWorkspace(ctx context.Context, workspace models.WorkspaceShape) (models.WorkspaceInSync, error)
	SaveCustomer(ctx context.Context, customer models.CustomerShape) (models.CustomerInSync, error)
	SaveMember(ctx context.Context, member models.MemberShape) (models.MemberInSync, error)
	SaveThread(ctx context.Context, member models.ThreadShape) (models.ThreadInSync, error)
}
