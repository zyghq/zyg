package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (c *CustomerDB) LookupByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string,
) (models.Customer, error) {
	var customer models.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND customer_id = $2`, workspaceId, customerId).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) FetchWorkspaceCustomerByExtId(ctx context.Context, workspaceId string, externalId string,
) (models.Customer, error) {
	var customer models.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND external_id = $2`, workspaceId, externalId).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) RetrieveWorkspaceCustomerByEmail(ctx context.Context, workspaceId string, email string,
) (models.Customer, error) {
	var customer models.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND email = $2`, workspaceId, email).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) LookupWorkspaceCustomerByPhone(ctx context.Context, workspaceId string, phone string,
) (models.Customer, error) {
	var customer models.Customer
	err := c.db.QueryRow(ctx, `SELECT
		workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND phone = $2`, workspaceId, phone).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, ErrQuery
	}

	return customer, nil
}

func (c *CustomerDB) UpsertCustomerByExtId(ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (workspace_id, external_id) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, external_id) = ($2, $3) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone, customer.Name,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, isCreated, ErrQuery
	}

	return customer, isCreated, nil
}

func (c *CustomerDB) UpsertCustomerByEmail(ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (workspace_id, email) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, email) = ($2, $4) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone, customer.Name,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, isCreated, ErrQuery
	}

	return customer, isCreated, nil
}

func (c *CustomerDB) UpsertCustomerByPhone(ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (workspace_id, phone) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, phone) = ($2, $5) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone, customer.Name,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, isCreated, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.Customer{}, isCreated, ErrQuery
	}

	return customer, isCreated, nil
}

func (c *CustomerDB) FetchCustomersByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Customer, error) {
	var customer models.Customer
	customers := make([]models.Customer, 0, 100)
	stmt := `SELECT workspace_id, customer_id, external_id, email, phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, _ := c.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	}, func() error {
		customers = append(customers, customer)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.Customer{}, ErrQuery
	}

	return customers, nil
}
