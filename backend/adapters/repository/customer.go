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
		"is_email_verified",
		"role",
		"created_at",
		"updated_at",
	}
}

// claimedMailCols returns the required columns for the `claimed_mail` table.
func claimedMailCols() builq.Columns {
	return builq.Columns{
		"claim_id", "workspace_id", "customer_id", "email",
		"has_conflict", "expires_at", "token",
		"is_mail_sent", "platform", "sender_id",
		"sender_status", "sent_at",
		"created_at", "updated_at",
	}
}

// customerEventCols returns the required columns for the `customer_event` table.
func customerEventCols() builq.Columns {
	return builq.Columns{
		"event_id", "customer_id",
		"title", "severity", "timestamp", "components",
		"created_at", "updated_at",
	}
}

func customerEventJoinedCols() builq.Columns {
	return builq.Columns{
		"e.event_id", "c.customer_id", "c.name",
		"e.title", "e.severity", "e.timestamp", "e.components",
		"e.created_at", "e.updated_at",
	}
}

// LookupWorkspaceCustomerById returns the workspace customer by ID with optional role.
func (c *CustomerDB) LookupWorkspaceCustomerById(
	ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error) {
	var customer models.Customer

	cols := customerCols()
	q := builq.New()
	params := []any{workspaceId, customerId}

	q("SELECT %s FROM %s", cols, "customer")
	q("WHERE workspace_id = %$ AND customer_id = %$", workspaceId, customerId)
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
		&customer.Email, &customer.Phone, &customer.Name, &customer.IsEmailVerified, &customer.Role,
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
		&customer.Email, &customer.Phone, &customer.Name, &customer.IsEmailVerified, &customer.Role,
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
	var insertB builq.Builder
	insertCols := customerCols()
	insertParams := []any{
		cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsEmailVerified, customer.Role,
		customer.CreatedAt, customer.UpdatedAt,
	}

	// Build the insert query.
	insertB.Addf("INSERT INTO customer (%s)", insertCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (workspace_id, external_id) DO NOTHING")
	insertB.Addf("RETURNING %s, TRUE AS is_created", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("error", err))
		return models.Customer{}, false, ErrQuery
	}

	// Build the select query required after insert
	q := builq.New()
	q("WITH ins AS (%s)", insertQuery)
	q("SELECT * FROM ins")
	q("UNION ALL")
	q("SELECT %s, FALSE AS is_created FROM customer", insertCols)
	q("WHERE (workspace_id, external_id) = ($2, $3)")
	q("AND NOT EXISTS (SELECT 1 FROM ins)")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return models.Customer{}, false, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var isCreated bool
	err = c.db.QueryRow(ctx, stmt, insertParams...).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsEmailVerified, &customer.Role,
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
	var insertB builq.Builder
	insertCols := customerCols()
	insertParams := []any{
		cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsEmailVerified, customer.Role,
		customer.CreatedAt, customer.UpdatedAt,
	}

	// Build the insert query.
	insertB.Addf("INSERT INTO customer (%s)", insertCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (workspace_id, email) DO NOTHING")
	insertB.Addf("RETURNING %s, TRUE AS is_created", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("error", err))
		return models.Customer{}, false, ErrQuery
	}

	// Build the select query required after insert
	q := builq.New()
	q("WITH ins AS (%s)", insertQuery)
	q("SELECT * FROM ins")
	q("UNION ALL")
	q("SELECT %s, FALSE AS is_created FROM customer", insertCols)
	q("WHERE (workspace_id, email) = ($2, $4)")
	q("AND NOT EXISTS (SELECT 1 FROM ins)")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return models.Customer{}, false, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var isCreated bool
	err = c.db.QueryRow(ctx, stmt, insertParams...).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsEmailVerified, &customer.Role,
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
	var insertB builq.Builder
	insertCols := customerCols()
	insertParams := []any{
		cId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.IsEmailVerified, customer.Role,
		customer.CreatedAt, customer.UpdatedAt,
	}

	// Build the insert query.
	insertB.Addf("INSERT INTO customer (%s)", insertCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (workspace_id, phone) DO NOTHING")
	insertB.Addf("RETURNING %s, TRUE AS is_created", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("error", err))
		return models.Customer{}, false, ErrQuery
	}

	// Build the select query required after insert
	q := builq.New()
	q("WITH ins AS (%s)", insertQuery)
	q("SELECT * FROM ins")
	q("UNION ALL")
	q("SELECT %s, FALSE AS is_created FROM customer", insertCols)
	q("WHERE (workspace_id, phone) = ($2, $5)")
	q("AND NOT EXISTS (SELECT 1 FROM ins)")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return models.Customer{}, false, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var isCreated bool
	err = c.db.QueryRow(ctx, stmt, insertParams...).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsEmailVerified, &customer.Role,
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
	limit := 100
	customers := make([]models.Customer, 0, limit)

	cols := customerCols()
	q := builq.New()
	params := []any{workspaceId}

	q("SELECT %s FROM customer", cols)
	q("WHERE workspace_id = %$", workspaceId)
	q("AND role <> 'visitor'")

	if role != nil {
		q("AND role = %$", *role)
		params = append(params, *role)
	}

	q("ORDER BY created_at DESC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return []models.Customer{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	rows, _ := c.db.Query(ctx, stmt, params...)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email, &customer.Phone,
		&customer.Name, &customer.IsEmailVerified, &customer.Role,
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
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone, name, role, is_email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (customer_id) DO UPDATE SET
			external_id = $3,
			email = $4,
			phone = $5,
			name = $6,
			role = $7,
			is_email_verified = $8,
			updated_at = now()
		RETURNING customer_id, workspace_id,
		external_id, email, phone, name, role,
		is_email_verified,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone, name, role,
	is_email_verified, created_at, updated_at, FALSE AS is_created FROM customer
	WHERE customer_id = $1 AND NOT EXISTS (SELECT 1 FROM ins)`

	var isCreated bool
	err := c.db.QueryRow(
		ctx, stmt, customer.CustomerId, customer.WorkspaceId, customer.ExternalId, customer.Email, customer.Phone,
		customer.Name, customer.Role, customer.IsEmailVerified,
	).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.Role, &customer.IsEmailVerified,
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
	q := builq.New()
	cols := customerCols()
	updateParams := []any{
		customer.ExternalId,
		customer.Email,
		customer.Phone,
		customer.Name,
		customer.IsEmailVerified,
		customer.Role,
		customer.CustomerId,
	}

	q("UPDATE customer SET")
	q("external_id = %$,", customer.ExternalId)
	q("email = %$,", customer.Email)
	q("phone = %$,", customer.Phone)
	q("name = %$,", customer.Name)
	q("is_email_verified = %$,", customer.IsEmailVerified)
	q("role = %$,", customer.Role)
	q("updated_at = NOW()")
	q("WHERE customer_id = %$", customer.CustomerId)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build update query", slog.Any("error", err))
		return models.Customer{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = c.db.QueryRow(ctx, stmt, updateParams...).Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.IsEmailVerified, &customer.Role,
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
	stmt := `SELECT exists (
		SELECT 1
        	FROM customer
        WHERE workspace_id = $1 AND email = $2
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
	ctx context.Context, event models.Event) (models.Event, error) {

	q := builq.New()
	cols := customerEventCols()

	q("INSERT INTO customer_event (%s)", cols)
	q("VALUES (%$, %$, %$, %$, %$, %$, %$, %$)",
		event.EventID, event.Customer.CustomerId,
		event.Title, event.Severity, event.Timestamp, event.Components,
		event.CreatedAt, event.UpdatedAt,
	)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Event{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = c.db.QueryRow(ctx, stmt,
		event.EventID, event.Customer.CustomerId,
		event.Title, event.Severity, event.Timestamp, event.Components,
		event.CreatedAt, event.UpdatedAt,
	).Scan(
		&event.EventID, &event.Customer.CustomerId,
		&event.Title, &event.Severity, &event.Timestamp, &event.Components,
		&event.CreatedAt, &event.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("error", err))
		return models.Event{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return models.Event{}, ErrQuery
	}
	return event, nil
}

func (c *CustomerDB) FetchEventsByCustomerId(
	ctx context.Context, customerId string) ([]models.Event, error) {
	var event models.Event
	var events []models.Event
	limit := 11
	events = make([]models.Event, 0, limit)

	q := builq.New()
	cols := customerEventJoinedCols()

	q("SELECT %s FROM customer_event e", cols)
	q("INNER JOIN customer c ON e.customer_id = c.customer_id")
	q("WHERE e.customer_id = %$", customerId)
	q("ORDER BY timestamp DESC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.Event{}, ErrQuery
	}

	rows, _ := c.db.Query(ctx, stmt, customerId)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&event.EventID, &event.Customer.CustomerId, &event.Customer.Name,
		&event.Title, &event.Severity, &event.Timestamp, &event.Components,
		&event.CreatedAt, &event.UpdatedAt,
	}, func() error {
		events = append(events, event)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.Event{}, ErrQuery
	}
	return events, nil
}
