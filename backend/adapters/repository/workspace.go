package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (w *WorkspaceDB) InsertWorkspaceByAccountId(
	ctx context.Context, memberName string, workspace models.Workspace) (models.Workspace, error) {
	var member models.Member
	tx, err := w.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	defer tx.Rollback(ctx)

	workspaceId := workspace.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO workspace(workspace_id, account_id, name)
		VALUES ($1, $2, $3)
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`, workspaceId, workspace.AccountId, workspace.Name).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Workspace{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	memberId := member.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, name, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		workspaceId, workspace.AccountId, memberId, memberName, models.MemberRole{}.Primary()).Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Workspace{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	// commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit db tx", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) ModifyWorkspaceById(
	ctx context.Context, workspace models.Workspace) (models.Workspace, error) {
	err := w.db.QueryRow(ctx, `UPDATE workspace SET
		name = $1, updated_at = NOW()
		WHERE workspace_id = $2
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`,
		workspace.Name, workspace.WorkspaceId).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to update query", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) ModifyWorkspaceLabelById(
	ctx context.Context, workspaceId string, label models.Label,
) (models.Label, error) {
	err := w.db.QueryRow(ctx, `UPDATE label SET
		name = $1, icon = $2, updated_at = NOW()
		WHERE workspace_id = $3 AND label_id = $4
		RETURNING
		label_id, workspace_id, name, icon, created_at, updated_at`,
		label.Name, label.Icon, workspaceId, label.LabelId).Scan(
		&label.LabelId, &label.WorkspaceId,
		&label.Name, &label.Icon, &label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to update query", slog.Any("err", err))
		return models.Label{}, ErrQuery
	}

	return label, nil
}

func (w *WorkspaceDB) FetchWorkspaceById(ctx context.Context, workspaceId string) (models.Workspace, error) {
	var workspace models.Workspace
	err := w.db.QueryRow(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE workspace_id = $1`, workspaceId).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Workspace{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) LookupLinkedWorkspaceByAccountId(
	ctx context.Context, accountId string, workspaceId string) (models.Workspace, error) {
	var workspace models.Workspace
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
			AND ws.workspace_id = $2
	`

	err := w.db.QueryRow(ctx, stmt, accountId, workspaceId).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Workspace{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	return workspace, nil
}

func (w *WorkspaceDB) FetchLinkedWorkspacesByAccountId(
	ctx context.Context, accountId string) ([]models.Workspace, error) {
	var workspace models.Workspace
	workspaces := make([]models.Workspace, 0, 100)
	stmt := `
		SELECT
			ws.workspace_id AS workspace_id,
			ws.account_id AS account_id,
			ws.name AS name,
			ws.created_at AS created_at,
			ws.updated_at AS updated_at
		FROM
			workspace ws
			INNER JOIN member m ON ws.workspace_id = m.workspace_id
		WHERE
			m.account_id = $1
		ORDER BY
			ws.created_at DESC
		LIMIT 100
	`

	rows, _ := w.db.Query(ctx, stmt, accountId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	}, func() error {
		workspaces = append(workspaces, workspace)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Workspace{}, ErrQuery
	}

	return workspaces, nil
}

func (w *WorkspaceDB) InsertLabelByName(
	ctx context.Context, label models.Label) (models.Label, bool, error) {
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

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Label{}, isCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Label{}, isCreated, ErrQuery
	}

	return label, isCreated, nil
}

func (w *WorkspaceDB) LookupWorkspaceLabelById(
	ctx context.Context, workspaceId string, labelId string) (models.Label, error) {
	var label models.Label
	err := w.db.QueryRow(ctx, `SELECT
		label_id, workspace_id, name, icon, created_at, updated_at
		FROM label WHERE workspace_id = $1 AND label_id = $2`, workspaceId, labelId).Scan(
		&label.LabelId, &label.WorkspaceId, &label.Name,
		&label.Icon, &label.CreatedAt, &label.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Label{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Label{}, ErrQuery
	}

	return label, nil
}

func (w *WorkspaceDB) RetrieveLabelsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Label, error) {
	var label models.Label
	labels := make([]models.Label, 0, 100)
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
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Label{}, ErrQuery
	}

	return labels, nil
}

func (w *WorkspaceDB) InsertMemberIntoWorkspace(ctx context.Context, workspaceId string, member models.Member) (models.Member, error) {
	var m models.Member
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
		return models.Member{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.Member{}, ErrQuery
	}

	return m, nil
}

func (w *WorkspaceDB) InsertWidgetIntoWorkspace(
	ctx context.Context, workspaceId string, widget models.Widget) (models.Widget, error) {
	widgetId := widget.GenId()
	err := w.db.QueryRow(ctx, `INSERT INTO widget(workspace_id, widget_id, name, configuration)
		VALUES ($1, $2, $3, $4)
		RETURNING
		workspace_id, widget_id, name, configuration, created_at, updated_at`,
		workspaceId, widgetId, widget.Name, widget.Configuration).Scan(
		&widget.WorkspaceId, &widget.WidgetId,
		&widget.Name, &widget.Configuration,
		&widget.CreatedAt, &widget.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Widget{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Widget{}, ErrQuery
	}

	return widget, nil
}

func (w *WorkspaceDB) RetrieveWidgetsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Widget, error) {
	var widget models.Widget
	widgets := make([]models.Widget, 0, 100)
	stmt := `SELECT widget_id, workspace_id, name, configuration,
		created_at, updated_at
		FROM widget WHERE workspace_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, _ := w.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&widget.WidgetId, &widget.WorkspaceId, &widget.Name, &widget.Configuration,
		&widget.CreatedAt, &widget.UpdatedAt,
	}, func() error {
		widgets = append(widgets, widget)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Widget{}, ErrQuery
	}

	return widgets, nil
}

func (w *WorkspaceDB) InsertSecretKeyIntoWorkspace(
	ctx context.Context, workspaceId string, sk string) (models.SecretKey, error) {
	var secretKey models.SecretKey
	err := w.db.QueryRow(ctx, `INSERT INTO secret_key(workspace_id, secret_key)
		VALUES ($1, $2)
		ON CONFLICT (workspace_id) DO UPDATE SET secret_key = $2
		RETURNING
		workspace_id, secret_key, created_at, updated_at`, workspaceId, sk).Scan(
		&secretKey.WorkspaceId, &secretKey.SecretKey,
		&secretKey.CreatedAt, &secretKey.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.SecretKey{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.SecretKey{}, ErrQuery
	}

	return secretKey, nil
}

func (r *WorkspaceDB) FetchSecretKeyByWorkspaceId(
	ctx context.Context, workspaceId string) (models.SecretKey, error) {
	var secretKey models.SecretKey
	err := r.db.QueryRow(ctx, `SELECT
		workspace_id,secret_key, created_at, updated_at
		FROM secret_key WHERE workspace_id = $1`, workspaceId).Scan(
		&secretKey.WorkspaceId, &secretKey.SecretKey,
		&secretKey.CreatedAt, &secretKey.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.SecretKey{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.SecretKey{}, ErrQuery
	}

	return secretKey, nil
}

func (r *WorkspaceDB) LookupWorkspaceWidget(ctx context.Context, widgetId string) (models.Widget, error) {
	var widget models.Widget

	stmt := `SELECT
		workspace_id, widget_id, name, configuration,
		created_at, updated_at
		FROM widget WHERE widget_id = $1`

	err := r.db.QueryRow(ctx, stmt, widgetId).Scan(
		&widget.WorkspaceId,
		&widget.WidgetId,
		&widget.Name,
		&widget.Configuration,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, sql.ErrNoRows) {
		return models.Widget{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Widget{}, ErrQuery
	}

	return widget, nil
}
