package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (c *CustomerDB) LookupWorkspaceCustomerById(
	ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error) {
	var customer models.Customer

	args := []any{workspaceId, customerId}
	stmt := `select
		workspace_id, customer_id, external_id, email, phone,
		name, avatar_url, anonymous_id,
		is_anonymous, role,
		created_at, updated_at
		from customer
		where
		workspace_id = $1 and customer_id = $2`

	if role != nil {
		stmt += " AND role = $3"
		args = append(args, *role)
	}

	err := c.db.QueryRow(ctx, stmt, args...).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone,
		&customer.Name, &customer.AvatarUrl, &customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Customer{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, ErrQuery
	}
	return customer, nil
}

func (c *CustomerDB) UpsertCustomerByExtId(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	avatarUrl := customer.GenerateAvatar(cId)
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, avatar_url, is_anonymous, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (workspace_id, external_id) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name, anonymous_id,
		is_anonymous, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name, avatar_url,
	anonymous_id, is_anonymous, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, external_id) = ($2, $3) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, avatarUrl, customer.IsAnonymous, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name, &customer.AvatarUrl, &customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt, &customer.UpdatedAt, &isCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrQuery
	}
	return customer, isCreated, nil
}

func (c *CustomerDB) UpsertCustomerByEmail(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	avatarUrl := customer.GenerateAvatar(cId)
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, avatar_url, is_anonymous, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (workspace_id, email) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name, avatar_url, anonymous_id,
		is_anonymous, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name, avatar_url,
	anonymous_id, is_anonymous, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, email) = ($2, $4) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, avatarUrl, customer.IsAnonymous, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name, &customer.AvatarUrl, &customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrQuery
	}
	return customer, isCreated, nil
}

func (c *CustomerDB) UpsertCustomerByPhone(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	avatarUrl := customer.GenerateAvatar(cId)
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, avatar_url, is_anonymous, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (workspace_id, phone) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name, avatar_url, anonymous_id,
		is_anonymous, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name, avatar_url,
	anonymous_id, is_anonymous, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, phone) = ($2, $5) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, avatarUrl, customer.IsAnonymous, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name, &customer.AvatarUrl, &customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrQuery
	}
	return customer, isCreated, nil
}

func (c *CustomerDB) FetchCustomersByWorkspaceId(
	ctx context.Context, workspaceId string, role *string) ([]models.Customer, error) {
	var customer models.Customer
	customers := make([]models.Customer, 0, 100)

	args := []any{workspaceId}

	stmt := `SELECT workspace_id, customer_id, external_id, email, phone,
		name, avatar_url, anonymous_id, is_anonymous, role,
		created_at, updated_at
		FROM customer
		WHERE
		workspace_id = $1`

	if role != nil {
		stmt += " AND role = $2"
		args = append(args, *role)
	}

	stmt += " ORDER BY created_at DESC LIMIT 100"

	rows, _ := c.db.Query(ctx, stmt, args...)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone,
		&customer.Name, &customer.AvatarUrl, &customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt, &customer.UpdatedAt,
	}, func() error {
		customers = append(customers, customer)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.Customer{}, ErrQuery
	}
	return customers, nil
}

func (c *CustomerDB) LookupSecretKeyByWidgetId(
	ctx context.Context, widgetId string) (models.WorkspaceSecret, error) {
	var sk models.WorkspaceSecret
	stmt := `select sk.workspace_id as workspace_id,
		sk.hmac as hmac, sk.created_at, sk.updated_at
		from widget w
		inner join workspace_secret sk on sk.workspace_id = w.workspace_id
		where w.widget_id = $1`

	err := c.db.QueryRow(ctx, stmt, widgetId).Scan(
		&sk.WorkspaceId, &sk.Hmac,
		&sk.CreatedAt, &sk.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.WorkspaceSecret{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.WorkspaceSecret{}, ErrQuery
	}
	return sk, nil
}

func (c *CustomerDB) UpsertCustomerByAnonId(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	avatarUrl := customer.GenerateAvatar(cId)
	stmt := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, anonymous_id, is_anonymous, name, avatar_url, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (anonymous_id) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name, avatar_url, anonymous_id,
		is_anonymous, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name, avatar_url,
	anonymous_id, is_anonymous, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE anonymous_id = $3 AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, stmt, cId, customer.WorkspaceId, customer.AnonId, customer.IsAnonymous,
		customer.Name, avatarUrl, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name, &customer.AvatarUrl, &customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, isCreated, ErrQuery
	}
	return customer, isCreated, nil
}

func (c *CustomerDB) ModifyCustomerById(
	ctx context.Context, customer models.Customer) (models.Customer, error) {
	stmt := `update customer set
		external_id = $2, email = $3, phone = $4, name = $5, avatar_url = $6, is_anonymous = $7, role = $8,
		updated_at = now()
		where
		customer_id = $1
		returning customer_id, workspace_id,
		external_id, email, phone,
		name, avatar_url,
		anonymous_id, is_anonymous, role,
		created_at, updated_at`
	err := c.db.QueryRow(ctx, stmt, customer.CustomerId,
		customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.AvatarUrl,
		customer.IsAnonymous, customer.Role).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name, &customer.AvatarUrl,
		&customer.AnonId,
		&customer.IsAnonymous, &customer.Role,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Customer{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, ErrQuery
	}
	return customer, nil
}

func (c *CustomerDB) CheckEmailExists(
	ctx context.Context, workspaceId string, email string) (bool, error) {
	var exists bool
	stmt := `select exists (
        select 1
        from customer
        where workspace_id = $1 and email = $2
    ) as exists`

	err := c.db.QueryRow(ctx, stmt, workspaceId, email).Scan(&exists)
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return exists, ErrQuery
	}
	return exists, nil
}

func (c *CustomerDB) InsertEmailIdentity(
	ctx context.Context, identity models.EmailIdentity) (models.EmailIdentity, error) {
	identityId := identity.GenId()

	stmt := `insert into email_identity (email_identity_id, customer_id, email, is_verified, has_conflict)
	values ($1, $2, $3, $4, $5)
		returning email_identity_id, customer_id, email, is_verified, has_conflict,
	created_at, updated_at
	`

	err := c.db.QueryRow(
		ctx, stmt, identityId, identity.CustomerId, identity.Email, identity.IsVerified, identity.HasConflict,
	).Scan(
		&identity.EmailIdentityId, &identity.CustomerId,
		&identity.Email, &identity.IsVerified, &identity.HasConflict,
		&identity.CreatedAt, &identity.UpdatedAt,
	)
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return identity, ErrQuery
	}
	return identity, nil
}
