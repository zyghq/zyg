package esync

import (
	"context"
	"database/sql"
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

func (sy *SyncDB) SaveCustomer(
	ctx context.Context, customer models.CustomerShape) (models.CustomerInSync, error) {
	var inSync models.CustomerInSync
	hub := sentry.GetHubFromContext(ctx)

	var externalID sql.NullString
	var email sql.NullString
	var phone sql.NullString

	if customer.ExternalID != nil {
		externalID.String = *customer.ExternalID
		externalID.Valid = true
	}
	if customer.Email != nil {
		email.String = *customer.Email
		email.Valid = true
	}
	if customer.Phone != nil {
		phone.String = *customer.Phone
		phone.Valid = true
	}

	stmt := `
	INSERT INTO customer (
		customer_id, workspace_id,
		external_id, email, phone,
		name, role, avatar_url, is_email_verified,
		created_at, updated_at,
		version_id, synced_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT (customer_id) DO UPDATE SET
	workspace_id = EXCLUDED.workspace_id,
	external_id = EXCLUDED.external_id,
	email = EXCLUDED.email,
	phone = EXCLUDED.phone,
	name = EXCLUDED.name,
	role = EXCLUDED.role,
	avatar_url = EXCLUDED.avatar_url,
	is_email_verified = EXCLUDED.is_email_verified,
	created_at = EXCLUDED.created_at,
	updated_at = EXCLUDED.updated_at,
	version_id = EXCLUDED.version_id,
	synced_at = EXCLUDED.synced_at
	RETURNING customer_id, synced_at, version_id`

	err := sy.db.QueryRow(
		ctx, stmt, customer.CustomerID, customer.WorkspaceID,
		externalID, email, phone,
		customer.Name, customer.Role, customer.AvatarURL,
		customer.IsEmailVerified, customer.CreatedAt, customer.UpdatedAt,
		customer.VersionID, customer.SyncedAt,
	).Scan(&inSync.CustomerID, &inSync.SyncedAt, &inSync.VersionID)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.CustomerInSync{}, err
	}
	return inSync, nil
}

func (sy *SyncDB) SaveMember(
	ctx context.Context, member models.MemberShape) (models.MemberInSync, error) {
	var inSync models.MemberInSync
	hub := sentry.GetHubFromContext(ctx)

	stmt := `
    INSERT INTO member (
        member_id, workspace_id,
        name, public_name, role, permissions, avatar_url,
        created_at, updated_at,
        version_id, synced_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    ON CONFLICT (member_id) DO UPDATE SET
	workspace_id = EXCLUDED.workspace_id,
	name = EXCLUDED.name,
	public_name = EXCLUDED.public_name,
	role = EXCLUDED.role,
	permissions = EXCLUDED.permissions,
	avatar_url = EXCLUDED.avatar_url,
	created_at = EXCLUDED.created_at,
	updated_at = EXCLUDED.updated_at,
	version_id = EXCLUDED.version_id,
	synced_at = EXCLUDED.synced_at
	RETURNING member_id, synced_at, version_id`

	err := sy.db.QueryRow(
		ctx, stmt,
		member.MemberID, member.WorkspaceID,
		member.Name, member.PublicName, member.Role, member.Permissions, member.AvatarURL,
		member.CreatedAt, member.UpdatedAt,
		member.VersionID, member.SyncedAt,
	).Scan(&inSync.MemberID, &inSync.SyncedAt, &inSync.VersionID)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.MemberInSync{}, err
	}
	return inSync, nil
}
