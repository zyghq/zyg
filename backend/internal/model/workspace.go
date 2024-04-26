package model

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

func (w Workspace) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		AccountId   string `json:"accountId"`
		Name        string `json:"name"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: w.WorkspaceId,
		AccountId:   w.AccountId,
		Name:        w.Name,
		CreatedAt:   w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   w.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (w Workspace) GenId() string {
	return "wrk" + xid.New().String()
}

func (w Workspace) Create(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	var member Member
	tx, err := db.Begin(ctx)

	// check if tx was started
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return w, ErrQuery
	}

	defer tx.Rollback(ctx)

	wId := w.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO workspace(workspace_id, account_id, name)
		VALUES ($1, $2, $3)
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`, wId, w.AccountId, w.Name).Scan(
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return w, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return w, ErrQuery
	}

	mId := member.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		wId, w.AccountId, mId, "primary").Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return w, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return w, ErrQuery
	}

	err = tx.Commit(ctx)
	// check if the tx was committed
	if err != nil {
		slog.Error("failed to commit db tx", "error", err)
		return w, ErrQuery
	}

	return w, nil
}

func (w Workspace) GetAccountWorkspace(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	err := db.QueryRow(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE account_id = $1 AND workspace_id = $2`, w.AccountId, w.WorkspaceId).Scan(
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return w, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return w, ErrQuery
	}

	return w, nil
}

func (w Workspace) GetById(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	err := db.QueryRow(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE workspace_id = $1`, w.WorkspaceId).Scan(
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return w, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return w, ErrQuery
	}

	return w, nil
}

func (w Workspace) GetListByAccountId(ctx context.Context, db *pgxpool.Pool) ([]Workspace, error) {
	ws := make([]Workspace, 0, 100)
	stmt := `SELECT workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE account_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, _ := db.Query(ctx, stmt, w.AccountId)

	_, err := pgx.ForEachRow(rows, []any{
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	}, func() error {
		ws = append(ws, w)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []Workspace{}, ErrQuery
	}

	defer rows.Close()

	return ws, nil
}
