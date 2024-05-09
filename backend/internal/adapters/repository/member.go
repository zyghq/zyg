package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/internal/domain"
)

func (m *MemberDB) GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (domain.Member, error) {
	var member domain.Member
	err := m.db.QueryRow(ctx, `SELECT
		workspace_id, account_id, member_id, name, role, created_at, updated_at
		FROM member WHERE AND account_id = $1 AND workspace_id = $2`, accountId, workspaceId).Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Member{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Member{}, ErrQuery
	}

	return member, nil
}
