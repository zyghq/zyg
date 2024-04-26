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

func (a Account) MarshalJSON() ([]byte, error) {
	aux := &struct {
		AccountId string `json:"accountId"`
		Email     string `json:"email"`
		Provider  string `json:"provider"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		AccountId: a.AccountId,
		Email:     a.Email,
		Provider:  a.Provider,
		Name:      a.Name,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (a Account) GenId() string {
	return "a_" + xid.New().String()
}

func (a Account) GetOrCreateByAuthUserId(ctx context.Context, db *pgxpool.Pool) (Account, bool, error) {
	var isCreated bool
	aId := a.GenId()
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

	err := db.QueryRow(ctx, stmt, aId, a.AuthUserId, a.Email, a.Provider, a.Name).Scan(
		&a.AccountId, &a.AuthUserId,
		&a.Email, &a.Provider, &a.Name,
		&a.CreatedAt, &a.UpdatedAt,
		&isCreated,
	)

	// no rows return error
	if errors.Is(err, pgx.ErrNoRows) {
		return a, isCreated, ErrEmpty
	}

	// query error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return a, isCreated, ErrQuery
	}
	return a, isCreated, nil
}

func (a Account) GetByAuthUserId(ctx context.Context, db *pgxpool.Pool) (Account, error) {
	err := db.QueryRow(ctx, `SELECT 
		account_id, auth_user_id, email,
		provider, name, created_at, updated_at
		FROM account WHERE auth_user_id = $1`, a.AuthUserId).Scan(
		&a.AccountId, &a.AuthUserId,
		&a.Email, &a.Provider, &a.Name,
		&a.CreatedAt, &a.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return a, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return a, ErrQuery
	}
	return a, nil
}

func (ap AccountPAT) MarshalJSON() ([]byte, error) {
	var description *string
	var token string
	if ap.Description.Valid {
		description = &ap.Description.String
	}

	maskLeft := func(s string) string {
		rs := []rune(s)
		for i := range rs[:len(rs)-4] {
			rs[i] = '*'
		}
		return string(rs)
	}

	if !ap.UnMask {
		token = maskLeft(ap.Token)
	} else {
		token = ap.Token
	}

	aux := &struct {
		AccountId   string  `json:"accountId"`
		PatId       string  `json:"patId"`
		Token       string  `json:"token"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		AccountId:   ap.AccountId,
		PatId:       ap.PatId,
		Token:       token,
		Name:        ap.Name,
		Description: description,
		CreatedAt:   ap.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ap.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (ap AccountPAT) GenId() string {
	return "ap_" + xid.New().String()
}

func (ap AccountPAT) Create(ctx context.Context, db *pgxpool.Pool) (AccountPAT, error) {
	apId := ap.GenId()

	token, err := GenToken(32, "pt_")
	if err != nil {
		slog.Error("failed to generate token got error", "error", err)
		return ap, err
	}

	stmt := `INSERT INTO account_pat(account_id, pat_id, token, name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING account_id, pat_id, token, name, description, created_at, updated_at`

	err = db.QueryRow(ctx, stmt, ap.AccountId, apId, token, ap.Name, ap.Description).Scan(
		&ap.AccountId, &ap.PatId, &ap.Token,
		&ap.Name, &ap.Description, &ap.CreatedAt, &ap.UpdatedAt,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return ap, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query got error", "error", err)
		return ap, ErrQuery
	}

	return ap, nil
}

func (ap AccountPAT) GetListByAccountId(ctx context.Context, db *pgxpool.Pool) ([]AccountPAT, error) {
	var pat AccountPAT
	aps := make([]AccountPAT, 0, 100)

	stmt := `SELECT account_id, pat_id, token, name, description,
		created_at, updated_at
		FROM account_pat WHERE account_id = $1
		ORDER BY created_at DESC LIMIT 100`

	// ignore the error - handled by the caller
	rows, _ := db.Query(ctx, stmt, ap.AccountId)

	_, err := pgx.ForEachRow(rows, []any{
		&pat.AccountId, &pat.PatId, &pat.Token,
		&pat.Name, &pat.Description, &pat.CreatedAt, &pat.UpdatedAt,
	}, func() error {
		aps = append(aps, pat)
		return nil
	})

	if err != nil {
		slog.Error("failed to query got error", "error", err)
		return []AccountPAT{}, ErrQuery
	}

	defer rows.Close()

	return aps, nil
}

func (ap AccountPAT) GetAccountByToken(ctx context.Context, db *pgxpool.Pool) (Account, error) {
	var account Account

	stmt := `SELECT
		a.account_id, a.email,
		a.provider, a.auth_user_id, a.name,
		a.created_at, a.updated_at
		FROM account a
		INNER JOIN account_pat ap ON a.account_id = ap.account_id
		WHERE ap.token = $1`

	err := db.QueryRow(ctx, stmt, ap.Token).Scan(
		&account.AccountId, &account.Email,
		&account.Provider, &account.AuthUserId, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return account, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return account, ErrQuery
	}

	return account, nil
}
