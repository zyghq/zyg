package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/zyghq/zyg"
)

func (c Customer) MarshalJSON() ([]byte, error) {
	var externalId, email, phone, name *string
	if c.ExternalId.Valid {
		externalId = &c.ExternalId.String
	}
	if c.Email.Valid {
		email = &c.Email.String
	}
	if c.Phone.Valid {
		phone = &c.Phone.String
	}
	if c.Name.Valid {
		name = &c.Name.String
	}

	aux := &struct {
		WorkspaceId string  `json:"workspaceId"`
		CustomerId  string  `json:"customerId"`
		ExternalId  *string `json:"externalId"`
		Email       *string `json:"email"`
		Phone       *string `json:"phone"`
		Name        *string `json:"name"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: c.WorkspaceId,
		CustomerId:  c.CustomerId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		Name:        name,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (c Customer) GenId() string {
	return "c_" + xid.New().String()
}

func (c Customer) GetWrkCustomerById(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	err := db.QueryRow(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND customer_id = $2`, c.WorkspaceId, c.CustomerId).Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, ErrQuery
	}

	return c, nil
}

func (c Customer) GetWrkCustomerByExtId(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	err := db.QueryRow(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND external_id = $2`, c.WorkspaceId, c.ExternalId).Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, ErrQuery
	}

	return c, nil
}

func (c Customer) GetWrkCustomerByEmail(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	err := db.QueryRow(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND email = $2`, c.WorkspaceId, c.Email).Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrEmpty

	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, ErrQuery
	}

	return c, nil
}

func (c Customer) GetWrkCustomerByPhone(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	err := db.QueryRow(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND phone = $2`, c.WorkspaceId, c.Phone).Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, ErrQuery
	}

	return c, nil
}

func (c Customer) GetOrCreateWrkCustomerByExtId(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {
	cId := c.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, external_id) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, external_id) = ($2, $3) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := db.QueryRow(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone).Scan(
		&c.CustomerId, &c.WorkspaceId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.CreatedAt,
		&c.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, isCreated, ErrQuery
	}

	return c, isCreated, nil
}

func (c Customer) GetOrCreateWrkCustomerByEmail(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {
	cId := c.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, email) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, email) = ($2, $4) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := db.QueryRow(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone).Scan(
		&c.CustomerId, &c.WorkspaceId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.CreatedAt,
		&c.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, isCreated, ErrQuery
	}

	return c, isCreated, nil
}

func (c Customer) GetOrCreateWrkCustomerByPhone(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {
	cId := c.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, phone) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, phone) = ($2, $5) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := db.QueryRow(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone).Scan(
		&c.CustomerId, &c.WorkspaceId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.CreatedAt,
		&c.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return c, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return c, isCreated, ErrQuery
	}

	return c, isCreated, nil
}

func (c Customer) MakeJWT() (string, error) {

	var externalId string
	var email string
	var phone string

	audience := []string{"customer"}

	sk, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
	if err != nil {
		return "", fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
	}

	if !c.ExternalId.Valid {
		externalId = ""
	} else {
		externalId = c.ExternalId.String
	}

	if !c.Email.Valid {
		email = ""
	} else {
		email = c.Email.String
	}

	if !c.Phone.Valid {
		phone = ""
	} else {
		phone = c.Phone.String
	}

	claims := CustomerJWTClaims{
		WorkspaceId: c.WorkspaceId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.zyg.ai",
			Subject:   c.CustomerId,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)), // Set ExpiresAt to 1 year from now
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        c.WorkspaceId + ":" + c.CustomerId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(sk))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token got error: %v", err)
	}
	return jwt, nil
}
