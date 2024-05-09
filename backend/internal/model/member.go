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

func (m Member) MarshalJSON() ([]byte, error) {
	var name *string
	if m.Name.Valid {
		name = &m.Name.String
	}
	aux := &struct {
		WorkspaceId string  `json:"workspaceId"`
		AccountId   string  `json:"accountId"`
		MemberId    string  `json:"memberId"`
		Name        *string `json:"name"`
		Role        string  `json:"role"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: m.WorkspaceId,
		AccountId:   m.AccountId,
		MemberId:    m.MemberId,
		Name:        name,
		Role:        m.Role,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (m Member) GenId() string {
	return "m_" + xid.New().String()
}

// done
func (m Member) GetWorkspaceMemberByAccountId(ctx context.Context, db *pgxpool.Pool) (Member, error) {
	err := db.QueryRow(ctx, `SELECT
		workspace_id, account_id, member_id, name, role, created_at, updated_at
		FROM member WHERE workspace_id = $1 AND account_id = $2`, m.WorkspaceId, m.AccountId).Scan(
		&m.WorkspaceId, &m.AccountId,
		&m.MemberId, &m.Name, &m.Role,
		&m.CreatedAt, &m.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return m, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return m, ErrQuery
	}

	return m, nil
}
