package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/internal/domain"
)

func (c *CustomerDB) GetByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string,
) (domain.Customer, error) {
	var customer domain.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND customer_id = $2`, workspaceId, customerId).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) GetWorkspaceCustomerByExtId(ctx context.Context, workspaceId string, externalId string,
) (domain.Customer, error) {
	var customer domain.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND external_id = $2`, workspaceId, externalId).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) GetWorkspaceCustomerByEmail(ctx context.Context, workspaceId string, email string,
) (domain.Customer, error) {
	var customer domain.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND email = $2`, workspaceId, email).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) GetWorkspaceCustomerByPhone(ctx context.Context, workspaceId string, phone string,
) (domain.Customer, error) {
	var customer domain.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND phone = $2`, workspaceId, phone).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) GetOrCreateCustomerByExtId(ctx context.Context, customer domain.Customer) (domain.Customer, bool, error) {
	cId := customer.GenId()
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
	err := c.db.QueryRow(ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, isCreated, ErrQuery
	}

	return customer, isCreated, nil
}

func (c *CustomerDB) GetOrCreateCustomerByEmail(ctx context.Context, customer domain.Customer) (domain.Customer, bool, error) {
	cId := customer.GenId()
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
	err := c.db.QueryRow(ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, isCreated, ErrQuery
	}

	return customer, isCreated, nil
}

func (c *CustomerDB) GetOrCreateCustomerByPhone(ctx context.Context, customer domain.Customer) (domain.Customer, bool, error) {
	cId := customer.GenId()
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
	err := c.db.QueryRow(ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.Customer{}, isCreated, ErrQuery
	}

	return customer, isCreated, nil
}
