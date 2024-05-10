package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/internal/domain"
)

func (w *WorkspaceDB) CreateWorkspace(ctx context.Context, workspace domain.Workspace) (domain.Workspace, error) {
	var member domain.Member
	tx, err := w.db.Begin(ctx)

	// check if tx was started
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	defer tx.Rollback(ctx)

	wId := workspace.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO workspace(workspace_id, account_id, name)
		VALUES ($1, $2, $3)
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`, wId, workspace.AccountId, workspace.Name).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Workspace{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	mId := member.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		wId, workspace.AccountId, mId, "primary").Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Workspace{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	err = tx.Commit(ctx)
	// check if the tx was committed
	if err != nil {
		slog.Error("failed to commit db tx", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) GetWorkspaceById(ctx context.Context, workspaceId string) (domain.Workspace, error) {
	var workspace domain.Workspace
	err := w.db.QueryRow(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE workspace_id = $1`, workspaceId).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Workspace{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (domain.Workspace, error) {
	var workspace domain.Workspace
	err := w.db.QueryRow(ctx, `SELECT
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE account_id = $1 AND workspace_id = $2`, accountId, workspaceId).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Workspace{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) GetListByAccountId(ctx context.Context, accountId string) ([]domain.Workspace, error) {
	var workspace domain.Workspace
	ws := make([]domain.Workspace, 0, 100)
	stmt := `SELECT workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE account_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, _ := w.db.Query(ctx, stmt, accountId)

	_, err := pgx.ForEachRow(rows, []any{
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	}, func() error {
		ws = append(ws, workspace)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.Workspace{}, ErrQuery
	}

	defer rows.Close()

	return ws, nil
}

func (w *WorkspaceDB) GetOrCreateLabel(ctx context.Context, label domain.Label) (domain.Label, bool, error) {
	var isCreated bool
	lId := label.GenId()
	stmt := `WITH ins AS (
		INSERT INTO label (label_id, workspace_id, name, icon)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (workspace_id, name) DO NOTHING
		RETURNING label_id, workspace_id, name, icon, created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT label_id, workspace_id, name, icon,
	created_at, updated_at, FALSE AS is_created FROM label
	WHERE workspace_id = $2 AND name = $3 AND NOT EXISTS (SELECT 1 FROM ins)`

	err := w.db.QueryRow(ctx, stmt, lId, label.WorkspaceId, label.Name, label.Icon).Scan(
		&label.LabelId, &label.WorkspaceId, &label.Name, &label.Icon,
		&label.CreatedAt, &label.UpdatedAt, &isCreated,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Label{}, isCreated, ErrEmpty
	}

	// query error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.Label{}, isCreated, ErrQuery
	}

	return label, isCreated, nil
}
