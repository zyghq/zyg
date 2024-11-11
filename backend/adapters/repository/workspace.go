package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/cristalhq/builq"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
)

func workspaceCols() builq.Columns {
	return builq.Columns{
		"workspace_id",
		"account_id",
		"name",
		"created_at",
		"updated_at",
	}
}

func (wrk *WorkspaceDB) InsertWorkspaceWithMembers(
	ctx context.Context, workspace models.Workspace, members []models.Member) (models.Workspace, error) {
	// Start the DB transaction
	tx, err := wrk.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	// Persist workspace.
	q := builq.New()
	workspaceCols := workspaceCols()

	q("INSERT INTO workspace (%s)", workspaceCols)
	q("VALUES (%$, %$, %$, %$, %$)",
		workspace.WorkspaceId, workspace.AccountId, workspace.Name, workspace.CreatedAt, workspace.UpdatedAt)
	q("RETURNING %s", workspaceCols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(
		ctx, stmt, workspace.WorkspaceId, workspace.AccountId, workspace.Name,
		workspace.CreatedAt, workspace.UpdatedAt).Scan(
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

	// Persist workspace members.
	insertCols := builq.Columns{
		"member_id",
		"workspace_id",
		"account_id",
		"name",
		"role",
		"created_at",
		"updated_at",
	}

	// Add insert-able queries to batch.
	// Handle the null case for account_id.
	batch := &pgx.Batch{}
	for _, m := range members {
		if m.IsMemberSystem() {
			var accountId sql.NullString // defaults to nullable string.
			// build insert query for system member.
			q = builq.New()
			q("INSERT INTO member (%s)", insertCols)
			q("VALUES (%$, %$, %$, %$, %$, %$, %$)", m.MemberId, m.WorkspaceId, accountId, m.Name, m.Role, m.CreatedAt, m.UpdatedAt)
			stmt, _, err = q.Build()
			if err != nil {
				slog.Error("failed to build query", slog.Any("err", err))
				return models.Workspace{}, ErrQuery
			}
			if zyg.DBQueryDebug() {
				debug := q.DebugBuild()
				debugQuery(debug)
			}
			batch.Queue(stmt, m.MemberId, m.WorkspaceId, accountId, m.Name, m.Role, m.CreatedAt, m.UpdatedAt)
		} else {
			accountId := sql.NullString{String: workspace.AccountId, Valid: true} // reference to account
			// build insert query for non-system member.
			q = builq.New()
			q("INSERT INTO member (%s)", insertCols)
			q("VALUES (%$, %$, %$, %$, %$, %$, %$)", m.MemberId, m.WorkspaceId, accountId, m.Name, m.Role, m.CreatedAt, m.UpdatedAt)
			stmt, _, err = q.Build()
			if err != nil {
				slog.Error("failed to build query", slog.Any("err", err))
				return models.Workspace{}, ErrQuery
			}
			if zyg.DBQueryDebug() {
				debug := q.DebugBuild()
				debugQuery(debug)
			}
			batch.Queue(stmt, m.MemberId, m.WorkspaceId, accountId, m.Name, m.Role, m.CreatedAt, m.UpdatedAt)
		}
	}

	results := tx.SendBatch(ctx, batch)
	defer func(results pgx.BatchResults) {
		err := results.Close()
		if err != nil {
			return
		}
	}(results)

	for _, m := range members {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				slog.Error("member already exists", slog.Any("memberId", m.MemberId))
				continue
			}
			slog.Error("failed to insert members in batch", slog.Any("err", err))
			return models.Workspace{}, ErrQuery // some other error, return should roll back.
		}
	}

	err = results.Close()
	if err != nil {
		slog.Error("failed to close batch results", slog.Any("err", err))
		return models.Workspace{}, ErrQuery
	}

	// All good, commit the transaction
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

func (wrk *WorkspaceDB) LookupWidgetSessionById(
	ctx context.Context, widgetId string, sessionId string) (models.WidgetSession, error) {
	var session models.WidgetSession

	stmt := `SELECT
		session_id, widget_id, data, expire_at,
		created_at, updated_at
		FROM widget_session WHERE widget_id = $1 AND session_id = $2`

	err := wrk.db.QueryRow(ctx, stmt, widgetId, sessionId).Scan(
		&session.SessionId, &session.WidgetId,
		&session.Data, &session.ExpireAt,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.WidgetSession{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.WidgetSession{}, ErrQuery
	}

	return session, nil
}

func (wrk *WorkspaceDB) UpsertWidgetSessionById(
	ctx context.Context, session models.WidgetSession) (models.WidgetSession, bool, error) {
	stmt := `WITH ins AS (
		INSERT INTO widget_session (session_id, widget_id, data, expire_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (session_id) DO UPDATE SET
			data = $3,
			expire_at = $4,
			updated_at = now()
		RETURNING session_id, widget_id, data, expire_at, created_at, updated_at, TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT session_id, widget_id, data, expire_at, created_at, updated_at, FALSE AS is_created FROM widget_session
	WHERE session_id = $1 AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := wrk.db.QueryRow(ctx, stmt, session.SessionId, session.WidgetId, session.Data, session.ExpireAt).Scan(
		&session.SessionId, &session.WidgetId, &session.Data, &session.ExpireAt,
		&session.CreatedAt, &session.UpdatedAt, &isCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.WidgetSession{}, isCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.WidgetSession{}, isCreated, ErrQuery
	}

	return session, isCreated, nil
}

// InsertSystemMember inserts a system member without the account ID.
func (wrk *WorkspaceDB) InsertSystemMember(ctx context.Context, member models.Member) (models.Member, error) {
	q := builq.New()
	cols := builq.Columns{
		"member_id",
		"workspace_id",
		"name",
		"role",
		"created_at",
		"updated_at",
	}

	q("INSERT INTO member (%s)", cols)
	q("VALUES (%$, %$, %$, %$, %$, %$)",
		member.MemberId, member.WorkspaceId, member.Name, member.Role, member.CreatedAt, member.UpdatedAt)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Member{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = wrk.db.QueryRow(
		ctx, stmt, member.MemberId, member.WorkspaceId,
		member.Name, member.Role, member.CreatedAt, member.UpdatedAt).Scan(
		&member.MemberId, &member.WorkspaceId, &member.Name, &member.Role, &member.CreatedAt, &member.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Member{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Member{}, ErrQuery
	}
	return member, nil
}

// LookupSystemMemberByOldest looks up system member with the oldest created in the workspace.
func (wrk *WorkspaceDB) LookupSystemMemberByOldest(
	ctx context.Context, workspaceId string) (models.Member, error) {
	var member models.Member

	q := builq.New()
	cols := builq.Columns{
		"member_id",
		"workspace_id",
		"name",
		"role",
		"created_at",
		"updated_at",
	}

	q("SELECT %s FROM member", cols)
	q("WHERE workspace_id = %$ AND role = %$", workspaceId, models.MemberRole{}.System())
	q("ORDER BY created_at ASC LIMIT 1")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Member{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = wrk.db.QueryRow(ctx, stmt, workspaceId, models.MemberRole{}.System()).Scan(
		&member.MemberId, &member.WorkspaceId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Member{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Member{}, ErrQuery
	}
	return member, nil
}
