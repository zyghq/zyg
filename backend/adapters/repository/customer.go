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
