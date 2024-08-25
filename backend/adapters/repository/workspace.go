package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (wrk *WorkspaceDB) InsertWorkspaceWithMember(
	ctx context.Context, workspace models.Workspace, member models.Member) (models.Workspace, error) {
	tx, err := wrk.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	workspaceId := workspace.GenId()
	err = tx.QueryRow(ctx, `insert into workspace(workspace_id, account_id, name)
		values ($1, $2, $3)
		returning
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
	err = tx.QueryRow(ctx, `insert into member(workspace_id, account_id, member_id, name, role)
		values ($1, $2, $3, $4, $5)
		returning
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		workspaceId, workspace.AccountId, memberId, member.Name, member.Role).Scan(
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

func (wrk *WorkspaceDB) ModifyWorkspaceById(
	ctx context.Context, workspace models.Workspace) (models.Workspace, error) {
	err := wrk.db.QueryRow(ctx, `update workspace set
		name = $1, updated_at = now()
		where workspace_id = $2
		returning
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

func (wrk *WorkspaceDB) ModifyLabelById(
	ctx context.Context, label models.Label,
) (models.Label, error) {
	err := wrk.db.QueryRow(ctx, `update label set
		name = $1, icon = $2, updated_at = now()
		where workspace_id = $3 and label_id = $4
		returning
		label_id, workspace_id, name, icon, created_at, updated_at`,
		label.Name, label.Icon, label.WorkspaceId, label.LabelId).Scan(
		&label.LabelId, &label.WorkspaceId,
		&label.Name, &label.Icon, &label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to update query", slog.Any("err", err))
		return models.Label{}, ErrQuery
	}

	return label, nil
}

func (wrk *WorkspaceDB) FetchByWorkspaceId(
	ctx context.Context, workspaceId string) (models.Workspace, error) {
	var workspace models.Workspace
	err := wrk.db.QueryRow(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE workspace_id = $1`, workspaceId).Scan(
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

func (wrk *WorkspaceDB) LookupWorkspaceByAccountId(
	ctx context.Context, workspaceId string, accountId string) (models.Workspace, error) {
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

	err := wrk.db.QueryRow(ctx, stmt, accountId, workspaceId).Scan(
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

func (wrk *WorkspaceDB) FetchWorkspacesByAccountId(
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

	rows, _ := wrk.db.Query(ctx, stmt, accountId)

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

func (wrk *WorkspaceDB) InsertLabelByName(
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

	err := wrk.db.QueryRow(ctx, stmt, lId, label.WorkspaceId, label.Name, label.Icon).Scan(
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

func (wrk *WorkspaceDB) LookupWorkspaceLabelById(
	ctx context.Context, workspaceId string, labelId string) (models.Label, error) {
	var label models.Label
	err := wrk.db.QueryRow(ctx, `SELECT
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

func (wrk *WorkspaceDB) FetchLabelsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Label, error) {
	var label models.Label
	labels := make([]models.Label, 0, 100)
	stmt := `SELECT label_id, workspace_id, name, icon, created_at, updated_at
		FROM label WHERE workspace_id = $1`

	rows, _ := wrk.db.Query(ctx, stmt, workspaceId)

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

func (wrk *WorkspaceDB) InsertWidget(
	ctx context.Context, widget models.Widget) (models.Widget, error) {
	widgetId := widget.GenId()
	err := wrk.db.QueryRow(ctx, `insert into widget(workspace_id, widget_id, name, configuration)
		values ($1, $2, $3, $4)
		returning
		workspace_id, widget_id, name, configuration, created_at, updated_at`,
		widget.WorkspaceId, widgetId, widget.Name, widget.Configuration).Scan(
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

func (wrk *WorkspaceDB) FetchWidgetsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Widget, error) {
	var widget models.Widget
	widgets := make([]models.Widget, 0, 100)
	stmt := `SELECT widget_id, workspace_id, name, configuration,
		created_at, updated_at
		FROM widget WHERE workspace_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, _ := wrk.db.Query(ctx, stmt, workspaceId)

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

func (wrk *WorkspaceDB) InsertWorkspaceSecret(
	ctx context.Context, workspaceId string, sk string) (models.WorkspaceSecret, error) {
	var secretKey models.WorkspaceSecret
	stmt := `
		insert into workspace_secret(workspace_id, hmac)
		values ($1, $2)
		on conflict (workspace_id) do update set hmac = $2
		returning
		workspace_id, hmac, created_at, updated_at
	`
	err := wrk.db.QueryRow(ctx, stmt, workspaceId, sk).Scan(
		&secretKey.WorkspaceId, &secretKey.Hmac,
		&secretKey.CreatedAt, &secretKey.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.WorkspaceSecret{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.WorkspaceSecret{}, ErrQuery
	}

	return secretKey, nil
}

func (wrk *WorkspaceDB) FetchSecretKeyByWorkspaceId(
	ctx context.Context, workspaceId string) (models.WorkspaceSecret, error) {
	var secretKey models.WorkspaceSecret
	err := wrk.db.QueryRow(ctx, `select
		workspace_id, hmac, created_at, updated_at
		from workspace_secret where workspace_id = $1`, workspaceId).Scan(
		&secretKey.WorkspaceId, &secretKey.Hmac,
		&secretKey.CreatedAt, &secretKey.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.WorkspaceSecret{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.WorkspaceSecret{}, ErrQuery
	}
	return secretKey, nil
}

func (wrk *WorkspaceDB) LookupWidgetById(
	ctx context.Context, widgetId string) (models.Widget, error) {
	var widget models.Widget

	stmt := `SELECT
		workspace_id, widget_id, name, configuration,
		created_at, updated_at
		FROM widget WHERE widget_id = $1`

	err := wrk.db.QueryRow(ctx, stmt, widgetId).Scan(
		&widget.WorkspaceId,
		&widget.WidgetId,
		&widget.Name,
		&widget.Configuration,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Widget{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Widget{}, ErrQuery
	}
	return widget, nil
}
