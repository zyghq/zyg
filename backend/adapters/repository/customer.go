package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cristalhq/builq"
	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
)

// customerCols returns the required columns for the `customer` table.
func customerCols() builq.Columns {
	return builq.Columns{
		"customer_id",
		"workspace_id",
		"external_id",
		"email",
		"phone",
		"name",
		"is_verified",
		"role",
		"created_at",
		"updated_at",
	}
}

func claimedMailCols() builq.Columns {
	return builq.Columns{
		"claim_id", "workspace_id", "customer_id", "email",
		"has_conflict", "expires_at", "token",
		"is_mail_sent", "platform", "sender_id",
		"sender_status", "sent_at",
		"created_at", "updated_at",
	}
}

func customerEventCols() builq.Columns {
	return builq.Columns{
		"event_id", "customer_id", "thread_id", "event", "event_body",
		"severity", "event_timestamp", "notification_status", "idempotency_key",
		"created_at", "updated_at",
	}
}

// LookupWorkspaceCustomerById returns the workspace customer by ID.
func (c *CustomerDB) LookupWorkspaceCustomerById(
	ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error) {
	var customer models.Customer

	params := []any{workspaceId, customerId}
	stmt := `select
		workspace_id, customer_id, external_id, email, phone,
		name, is_verified, role,
		created_at, updated_at
		from customer
		where
		workspace_id = $1 and customer_id = $2`

	if role != nil {
		stmt += " AND role = $3"
		params = append(params, *role)
	}

	err := c.db.QueryRow(ctx, stmt, params...).Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone,
		&customer.Name, &customer.IsVerified,
		&customer.Role, &customer.CreatedAt, &customer.UpdatedAt,
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

// LookupWorkspaceCustomerByEmail returns the workspace customer by email with optional role
func (c *CustomerDB) LookupWorkspaceCustomerByEmail(
	ctx context.Context, workspaceId string, email string, role *string) (models.Customer, error) {
	var customer models.Customer

	cols := customerCols()
	q := builq.New()
	params := []any{workspaceId, email}

	q("SELECT %s FROM %s", cols, "customer")
	q("WHERE workspace_id = %$ AND email = %$", workspaceId, email)
	if role != nil {
		q("AND role = %$", *role)
		params = append(params, *role)
	}

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Customer{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = c.db.QueryRow(ctx, stmt, params...).Scan(
		&customer.CustomerId, &customer.WorkspaceId, &customer.ExternalId,
		&customer.Email, &customer.Phone, &customer.Name, &customer.IsVerified, &customer.Role,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Customer{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Customer{}, ErrQuery
	}
	return customer, nil
}

// UpsertCustomerByExtId upsert(insert or update) the customer by external ID.
func (c *CustomerDB) UpsertCustomerByExtId(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, is_verified, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (workspace_id, external_id) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name,
		is_verified, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name,
	is_verified, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, external_id) = ($2, $3) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsVerified, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsVerified, &customer.Role,
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

// UpsertCustomerByEmail upsert(insert or update) the customer by email.
func (c *CustomerDB) UpsertCustomerByEmail(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, is_verified, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (workspace_id, email) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name,
		is_verified, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name,
	is_verified, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, email) = ($2, $4) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsVerified, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsVerified, &customer.Role,
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

// UpsertCustomerByPhone upsert(insert or update) the customer by phone.
func (c *CustomerDB) UpsertCustomerByPhone(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	cId := customer.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, is_verified, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (workspace_id, phone) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone, name,
		is_verified, role,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name,
	is_verified, role, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, phone) = ($2, $5) AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, st, cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsVerified, customer.Role,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsVerified, &customer.Role,
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

// FetchCustomersByWorkspaceId fetches the customers by workspace ID with optional role.
// Also, excludes the customer with role `visitor`, as we don't consider them as part of the system.
// visitors are more like tentative customers, who don't have any identity.
func (c *CustomerDB) FetchCustomersByWorkspaceId(
	ctx context.Context, workspaceId string, role *string) ([]models.Customer, error) {
	var customer models.Customer
	customers := make([]models.Customer, 0, 100)

	params := []any{workspaceId}

	stmt := `SELECT workspace_id, customer_id, external_id, email, phone,
		name, is_verified, role,
		created_at, updated_at
		FROM customer
		WHERE
		workspace_id = $1`

	if role != nil {
		stmt += " AND role = $2"
		params = append(params, *role)
	}

	// exclude visitors
	stmt += " AND role <> 'visitor'"

	stmt += " ORDER BY created_at DESC LIMIT 100"

	rows, _ := c.db.Query(ctx, stmt, params...)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email, &customer.Phone,
		&customer.Name,
		&customer.IsVerified, &customer.Role,
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

// LookupSecretKeyByWidgetId returns the secret key linked to workspace by widget ID.
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

// UpsertCustomerById upsert(insert or update) the customer by ID.
func (c *CustomerDB) UpsertCustomerById(
	ctx context.Context, customer models.Customer) (models.Customer, bool, error) {
	stmt := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, role, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (customer_id) DO UPDATE SET
			external_id = $3,
			email = $4,
			phone = $5,
			name = $6,
			role = $7,
			is_verified = $8,
			updated_at = now()
		RETURNING customer_id, workspace_id,
		external_id, email, phone, name, role,
		is_verified,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name, role,
	is_verified, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE customer_id = $1 AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, stmt, customer.CustomerId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.Role, customer.IsVerified,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.Role, &customer.IsVerified,
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

// ModifyCustomerById updates the customer by ID.
func (c *CustomerDB) ModifyCustomerById(
	ctx context.Context, customer models.Customer) (models.Customer, error) {
	stmt := `update customer set
		external_id = $2,
		email = $3,
		phone = $4,
		name = $5,
		is_verified = $6,
		role = $7,
		updated_at = now()
		where
		customer_id = $1
		returning customer_id, workspace_id,
		external_id, email, phone,
		name, is_verified, role,
		created_at, updated_at`
	err := c.db.QueryRow(ctx, stmt, customer.CustomerId,
		customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsVerified, customer.Role).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsVerified, &customer.Role,
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

// CheckEmailExists checks if the email exists in the customer table.
func (c *CustomerDB) CheckEmailExists(
	ctx context.Context, workspaceId string, email string) (bool, error) {
	var exists bool
	stmt := `select exists (
        select 1
        from customer
        where workspace_id = $1 and email = $2
    ) as exists`

	err := c.db.QueryRow(ctx, stmt, workspaceId, email).Scan(&exists)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return exists, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return exists, ErrQuery
	}
	return exists, nil
}

// InsertClaimedMail inserts claimed mail for verification
func (c *CustomerDB) InsertClaimedMail(
	ctx context.Context, claimed models.ClaimedMail) (models.ClaimedMail, error) {
	claimId := claimed.GenId()

	q := builq.New()
	cols := claimedMailCols()

	q("INSERT INTO %s (%s)", "claimed_mail", cols)
	q("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$)",
		claimId, claimed.WorkspaceId, claimed.CustomerId, claimed.Email,
		claimed.HasConflict, claimed.ExpiresAt, claimed.Token,
		claimed.IsMailSent, claimed.Platform, claimed.SenderId,
		claimed.SenderStatus, claimed.SentAt,
		claimed.CreatedAt, claimed.UpdatedAt,
	)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.ClaimedMail{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = c.db.QueryRow(ctx, stmt,
		claimId, claimed.WorkspaceId, claimed.CustomerId, claimed.Email,
		claimed.HasConflict, claimed.ExpiresAt, claimed.Token,
		claimed.IsMailSent, claimed.Platform, claimed.SenderId,
		claimed.SenderStatus, claimed.SentAt,
		claimed.CreatedAt, claimed.UpdatedAt,
	).Scan(
		&claimed.ClaimId, &claimed.WorkspaceId, &claimed.CustomerId, &claimed.Email,
		&claimed.HasConflict, &claimed.ExpiresAt, &claimed.Token,
		&claimed.IsMailSent, &claimed.Platform, &claimed.SenderId,
		&claimed.SenderStatus, &claimed.SentAt,
		&claimed.CreatedAt, &claimed.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.ClaimedMail{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.ClaimedMail{}, ErrQuery
	}
	return claimed, nil
}

// LookupClaimedMailByToken returns the claimed email verification by token.
// Always make sure the signed token is verified before usage.
func (c *CustomerDB) LookupClaimedMailByToken(
	ctx context.Context, token string) (models.ClaimedMail, error) {
	var claimed models.ClaimedMail
	q := builq.New()
	cols := claimedMailCols()

	q("SELECT %s FROM %s", cols, "claimed_mail")
	q("WHERE token = %$", token)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.ClaimedMail{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = c.db.QueryRow(ctx, stmt, token).Scan(
		&claimed.ClaimId, &claimed.WorkspaceId, &claimed.CustomerId, &claimed.Email,
		&claimed.HasConflict, &claimed.ExpiresAt, &claimed.Token,
		&claimed.IsMailSent, &claimed.Platform, &claimed.SenderId,
		&claimed.SenderStatus, &claimed.SentAt,
		&claimed.CreatedAt, &claimed.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ClaimedMail{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.ClaimedMail{}, ErrQuery
	}
	return claimed, nil
}

// DeleteCustomerClaimedMail deletes the claimed email token for workspace customer by email.
func (c *CustomerDB) DeleteCustomerClaimedMail(
	ctx context.Context, workspaceId string, customerId string, email string) error {

	q := builq.New()
	q("DELETE FROM %s", "claimed_mail")
	q("WHERE workspace_id = %$", workspaceId)
	q("AND customer_id = %$", customerId)
	q("AND email = %$", email)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	_, err = c.db.Exec(ctx, stmt, workspaceId, customerId, email)
	if err != nil {
		slog.Error("failed to delete claimed mail verification", slog.Any("error", err))
		return ErrQuery
	}
	return nil
}

func (c *CustomerDB) LookupLatestClaimedMail(
	ctx context.Context, workspaceId string, customerId string) (models.ClaimedMail, error) {
	var claimed models.ClaimedMail
	q := builq.New()
	cols := claimedMailCols()

	q("SELECT %s FROM %s", cols, "claimed_mail")
	q("WHERE workspace_id = %$", workspaceId)
	q("AND customer_id = %$", customerId)
	q("ORDER BY created_at DESC LIMIT 1")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.ClaimedMail{}, ErrQuery
	}

	err = c.db.QueryRow(ctx, stmt, workspaceId, customerId).Scan(
		&claimed.ClaimId, &claimed.WorkspaceId, &claimed.CustomerId, &claimed.Email,
		&claimed.HasConflict, &claimed.ExpiresAt, &claimed.Token,
		&claimed.IsMailSent, &claimed.Platform, &claimed.SenderId,
		&claimed.SenderStatus, &claimed.SentAt,
		&claimed.CreatedAt, &claimed.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ClaimedMail{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.ClaimedMail{}, ErrQuery
	}
	return claimed, nil
}

func (c *CustomerDB) InsertEvent(
	ctx context.Context, event models.CustomerEvent) (models.CustomerEvent, error) {
	q := builq.New()
	cols := customerEventCols()

	q("INSERT INTO %s (%s)", "customer_event", cols)
	q("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$)",
		event.EventId, event.CustomerId, event.ThreadId, event.Event, event.EventBody,
		event.Severity, event.EventTimestamp, event.NotificationStatus, event.IdempotencyKey,
		event.CreatedAt, event.UpdatedAt,
	)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.CustomerEvent{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = c.db.QueryRow(ctx, stmt,
		event.EventId, event.CustomerId, event.ThreadId, event.Event, event.EventBody,
		event.Severity, event.EventTimestamp, event.NotificationStatus, event.IdempotencyKey,
		event.CreatedAt, event.UpdatedAt,
	).Scan(
		&event.EventId, &event.CustomerId, &event.ThreadId, &event.Event, &event.EventBody,
		&event.Severity, &event.EventTimestamp, &event.NotificationStatus, &event.IdempotencyKey,
		&event.CreatedAt, &event.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.CustomerEvent{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.CustomerEvent{}, ErrQuery
	}
	return event, nil
}

func (c *CustomerDB) FetchEventsByCustomerId(
	ctx context.Context, customerId string) ([]models.CustomerEvent, error) {
	var event models.CustomerEvent
	limit := 11
	events := make([]models.CustomerEvent, 0, limit)

	q := builq.New()
	cols := customerEventCols()

	q("SELECT %s FROM %s", cols, "customer_event")
	q("WHERE customer_id = %$ ORDER BY event_timestamp DESC", customerId)
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.CustomerEvent{}, ErrQuery
	}

	rows, _ := c.db.Query(ctx, stmt, customerId)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&event.EventId, &event.CustomerId, &event.ThreadId, &event.Event, &event.EventBody,
		&event.Severity, &event.EventTimestamp, &event.NotificationStatus, &event.IdempotencyKey,
		&event.CreatedAt, &event.UpdatedAt,
	}, func() error {
		events = append(events, event)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.CustomerEvent{}, ErrQuery
	}

	return events, nil
}
