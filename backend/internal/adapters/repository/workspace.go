package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/internal/domain"
)

func (w *WorkspaceDB) CreateWorkspaceByAccount(ctx context.Context, account domain.Account, workspace domain.Workspace,
) (domain.Workspace, error) {
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
	err = tx.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, name, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		wId, workspace.AccountId, mId, account.Name, domain.MemberRole{}.Primary()).Scan(
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

func (w *WorkspaceDB) UpdateWorkspaceById(
	ctx context.Context, workspace domain.Workspace,
) (domain.Workspace, error) {
	err := w.db.QueryRow(ctx, `UPDATE workspace SET
		name = $1, updated_at = NOW()
		WHERE workspace_id = $2
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`, workspace.Name, workspace.WorkspaceId).Scan(
		&workspace.WorkspaceId, &workspace.AccountId, &workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to update query", "error", err)
		return domain.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) UpdateWorkspaceLabelById(
	ctx context.Context, workspaceId string, label domain.Label,
) (domain.Label, error) {
	err := w.db.QueryRow(ctx, `UPDATE label SET
		name = $1, icon = $2, updated_at = NOW()
		WHERE workspace_id = $3 AND label_id = $4
		RETURNING
		label_id, workspace_id, name, icon, created_at, updated_at`, label.Name, label.Icon, workspaceId, label.LabelId).Scan(
		&label.LabelId, &label.WorkspaceId, &label.Name, &label.Icon, &label.CreatedAt, &label.UpdatedAt,
	)

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to update query", "error", err)
		return domain.Label{}, ErrQuery
	}

	return label, nil
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

func (w *WorkspaceDB) GetByAccountWorkspaceId(
	ctx context.Context, accountId string, workspaceId string,
) (domain.Workspace, error) {
	var workspace domain.Workspace

	stmt := `
		SELECT
			ws.workspace_id as workspace_id,
			ws.account_id as account_id,
			ws.name as name,
			ws.created_at as created_at,
			ws.updated_at as updated_at
		FROM
			workspace ws
		INNER JOIN
			member m ON ws.workspace_id = m.workspace_id
		WHERE
			m.account_id = $1
			AND ws.workspace_id = $2`

	err := w.db.QueryRow(ctx, stmt, accountId, workspaceId).Scan(
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

func (w *WorkspaceDB) GetListByMemberAccountId(ctx context.Context, accountId string) ([]domain.Workspace, error) {
	var workspace domain.Workspace
	ws := make([]domain.Workspace, 0, 100)
	// stmt := `SELECT workspace_id, account_id, name, created_at, updated_at
	// 	FROM workspace WHERE account_id = $1
	// 	ORDER BY created_at DESC LIMIT 100`

	stmt := `SELECT
		ws.workspace_id as workspace_id,
		ws.account_id as account_id,
		ws.name as name,
		ws.created_at as created_at,
		ws.updated_at as updated_at
	FROM
		workspace ws
	INNER JOIN
		member m ON ws.workspace_id = m.workspace_id
	WHERE
		m.account_id = $1
	ORDER BY ws.created_at DESC LIMIT 100`

	rows, _ := w.db.Query(ctx, stmt, accountId)

	defer rows.Close()

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

func (w *WorkspaceDB) GetWorkspaceLabelById(
	ctx context.Context, workspaceId string, labelId string,
) (domain.Label, error) {
	var label domain.Label
	err := w.db.QueryRow(ctx, `SELECT
		label_id, workspace_id, name, icon, created_at, updated_at
		FROM label WHERE workspace_id = $1 AND label_id = $2`, workspaceId, labelId).Scan(
		&label.LabelId, &label.WorkspaceId, &label.Name,
		&label.Icon, &label.CreatedAt, &label.UpdatedAt,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Label{}, ErrEmpty
	}

	// query error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Label{}, ErrQuery
	}

	return label, nil
}

func (w *WorkspaceDB) GetLabelListByWorkspaceId(ctx context.Context, workspaceId string) ([]domain.Label, error) {
	var label domain.Label
	labels := make([]domain.Label, 0, 100)
	stmt := `SELECT label_id, workspace_id, name, icon, created_at, updated_at
		FROM label WHERE workspace_id = $1`

	rows, _ := w.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&label.LabelId, &label.WorkspaceId, &label.Name, &label.Icon,
		&label.CreatedAt, &label.UpdatedAt,
	}, func() error {
		labels = append(labels, label)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.Label{}, ErrQuery
	}

	return labels, nil
}

func (w *WorkspaceDB) AddMemberByWorkspaceId(ctx context.Context, workspaceId string, member domain.Member) (domain.Member, error) {
	var m domain.Member
	memberId := member.GenId()
	err := w.db.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, name, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		workspaceId, member.AccountId, memberId, member.Name, member.Role).Scan(
		&m.WorkspaceId, &m.AccountId,
		&m.MemberId, &m.Name, &m.Role,
		&m.CreatedAt, &m.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Member{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.Member{}, ErrQuery
	}

	return m, nil
}
