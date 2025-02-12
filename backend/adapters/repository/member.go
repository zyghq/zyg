package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cristalhq/builq"
	"github.com/zyghq/zyg"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

// Returns the required columns for the member table.
// The order of the columns matters when returning the results.
func memberCols() builq.Columns {
	return builq.Columns{
		"member_id",    // PK
		"workspace_id", // FK to workspace
		"name",
		"role",
		"created_at",
		"updated_at",
	}
}

// LookupByWorkspaceAccountId returns the member by workspace ID and account ID.
// The member is uniquely identified by the combination of `workspace_id` and `account_id`
// Human Member can authenticate to the workspace, hence the link to account ID.
func (m *MemberDB) LookupByWorkspaceAccountId(
	ctx context.Context, workspaceId string, accountId string) (models.Member, error) {
	var member models.Member

	q := builq.New()
	q("SELECT %s FROM member", memberCols())
	q("WHERE workspace_id = %$ AND account_id = %$", workspaceId, accountId)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Member{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = m.db.QueryRow(ctx, stmt, workspaceId, accountId).Scan(
		&member.MemberId, &member.WorkspaceId,
		&member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Member{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Member{}, ErrQuery
	}
	return member, nil
}

func (m *MemberDB) FetchMembersByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Member, error) {
	var member models.Member
	limit := 100
	members := make([]models.Member, 0, limit)

	q := builq.New()
	q("SELECT %s FROM member", memberCols())
	q("WHERE workspace_id = %$", workspaceId)
	q("ORDER BY created_at DESC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.Member{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	rows, _ := m.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&member.MemberId, &member.WorkspaceId,
		&member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	}, func() error {
		members = append(members, member)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.Member{}, ErrQuery
	}
	return members, nil
}

func (m *MemberDB) FetchByWorkspaceMemberId(
	ctx context.Context, workspaceId string, memberId string) (models.Member, error) {
	var member models.Member

	q := builq.New()
	q("SELECT %s FROM member", memberCols())
	q("WHERE workspace_id = %$ AND member_id = %$", workspaceId, memberId)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Member{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = m.db.QueryRow(ctx, stmt, workspaceId, memberId).Scan(
		&member.MemberId, &member.WorkspaceId,
		&member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Member{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Member{}, ErrQuery
	}
	return member, nil
}
