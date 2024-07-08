package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (m *MemberDB) GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (models.Member, error) {
	var member models.Member
	err := m.db.QueryRow(ctx, `SELECT
		workspace_id, account_id, member_id, name, role, created_at, updated_at
		FROM member WHERE account_id = $1 AND workspace_id = $2`, accountId, workspaceId).Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Member{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Member{}, ErrQuery
	}

	return member, nil
}

func (m *MemberDB) GetListByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Member, error) {
	var member models.Member
	members := make([]models.Member, 0, 100)
	stmt := `SELECT workspace_id, account_id, member_id, name, role,
	created_at, updated_at
	FROM member WHERE workspace_id = $1 ORDER BY created_at DESC LIMIT 100`

	rows, _ := m.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&member.WorkspaceId, &member.AccountId, &member.MemberId, &member.Name,
		&member.Role, &member.CreatedAt, &member.UpdatedAt,
	}, func() error {
		members = append(members, member)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.Member{}, ErrQuery
	}

	return members, nil
}

func (m *MemberDB) GetByWorkspaceMemberId(ctx context.Context, workspaceId string, memberId string) (models.Member, error) {
	var member models.Member
	err := m.db.QueryRow(ctx, `SELECT
		workspace_id, account_id, member_id, name, role, created_at, updated_at
		FROM member WHERE workspace_id = $1 AND member_id = $2`, workspaceId, memberId).Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Member{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Member{}, ErrQuery
	}

	return member, nil
}
