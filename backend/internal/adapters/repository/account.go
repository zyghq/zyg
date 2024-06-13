package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/internal/domain"
)

func (a *AccountDB) GetOrCreateByAuthUserId(ctx context.Context, account domain.Account,
) (domain.Account, bool, error) {
	var isCreated bool

	aId := account.GenId()
	stmt := `WITH ins AS (
		INSERT INTO account(account_id, auth_user_id, email, provider, name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (auth_user_id) DO NOTHING
		RETURNING
		account_id, auth_user_id, email,
		provider, name,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT account_id, auth_user_id, email, provider, name,
	created_at, updated_at, FALSE AS is_created FROM account
	WHERE auth_user_id = $2 AND NOT EXISTS (SELECT 1 FROM ins)`

	err := a.db.QueryRow(ctx, stmt, aId, account.AuthUserId, account.Email,
		account.Provider, account.Name).Scan(
		&account.AccountId, &account.AuthUserId,
		&account.Email, &account.Provider, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
		&isCreated,
	)

	// no rows returned error
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Account{}, isCreated, ErrEmpty
	}

	// query error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.Account{}, isCreated, ErrQuery
	}
	return account, isCreated, nil
}

func (a *AccountDB) GetByAuthUserId(ctx context.Context, authUserId string,
) (domain.Account, error) {
	var account domain.Account

	err := a.db.QueryRow(ctx, `SELECT 
		account_id, auth_user_id, email,
		provider, name, created_at, updated_at
		FROM account WHERE auth_user_id = $1`, authUserId).Scan(
		&account.AccountId, &account.AuthUserId,
		&account.Email, &account.Provider, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Account{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Account{}, ErrQuery
	}
	return account, nil
}

func (a *AccountDB) CreatePersonalAccessToken(ctx context.Context, ap domain.AccountPAT,
) (domain.AccountPAT, error) {
	apId := ap.GenId()
	token, err := domain.GenToken(32, "pt_")
	if err != nil {
		slog.Error("failed to generate token got error", "error", err)
		return domain.AccountPAT{}, err
	}

	stmt := `INSERT INTO account_pat(account_id, pat_id, token, name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING account_id, pat_id, token, name, description, created_at, updated_at`

	err = a.db.QueryRow(ctx, stmt, ap.AccountId, apId, token, ap.Name, ap.Description).Scan(
		&ap.AccountId, &ap.PatId, &ap.Token,
		&ap.Name, &ap.Description, &ap.CreatedAt, &ap.UpdatedAt,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.AccountPAT{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query got error", "error", err)
		return domain.AccountPAT{}, ErrQuery
	}

	return ap, nil
}

func (a *AccountDB) GetPatListByAccountId(ctx context.Context, accountId string) ([]domain.AccountPAT, error) {
	var pat domain.AccountPAT
	aps := make([]domain.AccountPAT, 0, 100)

	stmt := `SELECT account_id, pat_id, token, name, description,
		created_at, updated_at
		FROM account_pat WHERE account_id = $1
		ORDER BY created_at DESC LIMIT 100`

	// ignore the error - handled by the caller
	rows, _ := a.db.Query(ctx, stmt, accountId)

	defer rows.Close()

	// iterate over the each row
	// specific to pgx
	_, err := pgx.ForEachRow(rows, []any{
		&pat.AccountId, &pat.PatId, &pat.Token,
		&pat.Name, &pat.Description, &pat.CreatedAt, &pat.UpdatedAt,
	}, func() error {
		aps = append(aps, pat)
		return nil
	})

	if err != nil {
		slog.Error("failed to query got error", "error", err)
		return []domain.AccountPAT{}, ErrQuery
	}

	return aps, nil
}

func (a *AccountDB) GetPatByPatId(ctx context.Context, patId string) (domain.AccountPAT, error) {
	var pat domain.AccountPAT

	stmt := `SELECT
		account_id, pat_id, token, name, description,
		created_at, updated_at
		FROM account_pat
		WHERE pat_id = $1`

	err := a.db.QueryRow(ctx, stmt, patId).Scan(
		&pat.AccountId, &pat.PatId, &pat.Token,
		&pat.Name, &pat.Description, &pat.CreatedAt, &pat.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.AccountPAT{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.AccountPAT{}, ErrQuery
	}

	return pat, nil
}

func (a *AccountDB) HardDeletePatByPatId(ctx context.Context, patId string) error {
	stmt := `DELETE FROM account_pat WHERE pat_id = $1`
	_, err := a.db.Exec(ctx, stmt, patId)
	if err != nil {
		slog.Error("failed to delete pat", "error", err)
		return ErrQuery
	}
	return nil
}

func (a *AccountDB) GetAccountByToken(ctx context.Context, token string) (domain.Account, error) {
	var account domain.Account

	stmt := `SELECT
		a.account_id, a.email,
		a.provider, a.auth_user_id, a.name,
		a.created_at, a.updated_at
		FROM account a
		INNER JOIN account_pat ap ON a.account_id = ap.account_id
		WHERE ap.token = $1`

	err := a.db.QueryRow(ctx, stmt, token).Scan(
		&account.AccountId, &account.Email,
		&account.Provider, &account.AuthUserId, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Account{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Account{}, ErrQuery
	}

	return account, nil
}
