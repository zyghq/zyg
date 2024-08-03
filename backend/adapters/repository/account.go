package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (a *AccountDB) UpsertByAuthUserId(
	ctx context.Context, account models.Account) (models.Account, bool, error) {
	var isCreated bool
	accountId := account.GenId()
	stmt := `
		WITH ins AS (
		INSERT INTO account(account_id, auth_user_id, email, provider, name)
			VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (auth_user_id)
			DO NOTHING
		RETURNING
		account_id, auth_user_id, email, provider, name,
		created_at, updated_at, 
		TRUE AS is_created)
		SELECT *
		FROM ins
		UNION ALL
		SELECT account_id, auth_user_id, email, provider,
			name, created_at, updated_at, FALSE AS is_created
		FROM account
		WHERE auth_user_id = $2
			AND NOT EXISTS (SELECT 1 FROM ins)
	`

	err := a.db.QueryRow(ctx, stmt, accountId, account.AuthUserId, account.Email,
		account.Provider, account.Name).Scan(
		&account.AccountId, &account.AuthUserId,
		&account.Email, &account.Provider, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
		&isCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Account{}, isCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Account{}, isCreated, ErrQuery
	}
	return account, isCreated, nil
}

func (a *AccountDB) FetchByAuthUserId(
	ctx context.Context, authUserId string) (models.Account, error) {
	var account models.Account

	err := a.db.QueryRow(ctx, `SELECT 
		account_id, auth_user_id, email,
		provider, name, created_at, updated_at
		FROM account WHERE auth_user_id = $1`, authUserId).Scan(
		&account.AccountId, &account.AuthUserId,
		&account.Email, &account.Provider, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Account{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Account{}, ErrQuery
	}
	return account, nil
}

func (a *AccountDB) InsertPersonalAccessToken(
	ctx context.Context, pat models.AccountPAT) (models.AccountPAT, error) {
	patId := pat.GenId()
	token, err := models.GenToken(32, "pt")
	if err != nil {
		slog.Error("failed to generate pat token", slog.Any("err", err))
		return models.AccountPAT{}, err
	}

	stmt := `
		INSERT INTO account_pat(account_id, pat_id, token, name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
		account_id, pat_id, token, name, description, created_at, updated_at
	`

	err = a.db.QueryRow(ctx, stmt, pat.AccountId, patId, token, pat.Name, pat.Description).Scan(
		&pat.AccountId, &pat.PatId, &pat.Token,
		&pat.Name, &pat.Description, &pat.CreatedAt, &pat.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.AccountPAT{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.AccountPAT{}, ErrQuery
	}

	return pat, nil
}

func (a *AccountDB) FetchPatsByAccountId(
	ctx context.Context, accountId string) ([]models.AccountPAT, error) {
	var pat models.AccountPAT
	aps := make([]models.AccountPAT, 0, 100)

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
		slog.Error("failed to query", slog.Any("err", err))
		return []models.AccountPAT{}, ErrQuery
	}

	return aps, nil
}

func (a *AccountDB) FetchPatById(
	ctx context.Context, patId string) (models.AccountPAT, error) {
	var pat models.AccountPAT
	stmt := `
		SELECT account_id, pat_id, token, name,
		description, created_at, updated_at
		FROM account_pat
		WHERE pat_id = $1
	`

	err := a.db.QueryRow(ctx, stmt, patId).Scan(
		&pat.AccountId, &pat.PatId, &pat.Token,
		&pat.Name, &pat.Description, &pat.CreatedAt, &pat.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.AccountPAT{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.AccountPAT{}, ErrQuery
	}

	return pat, nil
}

func (a *AccountDB) DeletePatById(
	ctx context.Context, patId string) error {
	stmt := `DELETE FROM account_pat WHERE pat_id = $1`
	_, err := a.db.Exec(ctx, stmt, patId)
	if err != nil {
		slog.Error("failed to delete query", slog.Any("err", err))
		return ErrQuery
	}
	return nil
}

func (a *AccountDB) LookupByToken(
	ctx context.Context, token string) (models.Account, error) {
	var account models.Account

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
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Account{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Account{}, ErrQuery
	}

	return account, nil
}
