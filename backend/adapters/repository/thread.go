package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cristalhq/builq"
	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
)

// Returns the required columns for the thread table.
// The order of the columns matters when returning the results.
func threadCols() builq.Columns {
	return builq.Columns{
		"thread_id",    // PK
		"workspace_id", // FK to workspace
		"customer_id",  // FK to customer
		"assignee_id",  // FK Nullable to member
		"assigned_at",  // Nullable
		"title",
		"description",
		"status",
		"status_changed_at",
		"status_changed_by_id",
		"stage",
		"replied",
		"priority",
		"channel",
		"last_inbound_at",  // Nullable last inbound at
		"last_outbound_at", // Nullable last outbound at
		"created_by_id",    // FK to member
		"updated_by_id",    // FK to member
		"created_at",
		"updated_at",
	}
}

// Returns the required columns and joined for the thread table.
// The order of the columns matters when returning the results.
func threadJoinedCols() builq.Columns {
	return builq.Columns{
		"th.thread_id",
		"th.workspace_id",
		"c.customer_id",
		"c.name",
		"am.member_id",
		"am.name",
		"th.assigned_at",
		"th.title",
		"th.description",
		"th.status",
		"th.status_changed_at",
		"scm.member_id",
		"scm.name",
		"th.stage",
		"th.replied",
		"th.priority",
		"th.channel",
		"th.last_inbound_at",
		"th.last_outbound_at",
		"mc.member_id",
		"mc.name",
		"mu.member_id",
		"mu.name",
		"th.created_at",
		"th.updated_at",
	}
}

func threadActivityCols() builq.Columns {
	return builq.Columns{
		"activity_id", // PK
		"thread_id",   // FK to thread
		"activity_type",
		"customer_id", // FK Nullable to customer
		"member_id",   // FK Nullable to member
		"body",
		"created_at",
		"updated_at",
	}
}

func threadActivityJoinedCols() builq.Columns {
	return builq.Columns{
		"act.activity_id",
		"act.thread_id",
		"act.activity_type",
		"c.customer_id",
		"c.name",
		"m.member_id",
		"m.name",
		"act.body",
		"act.created_at",
		"act.updated_at",
	}
}

func postmarkMessageLogCols() builq.Columns {
	return builq.Columns{
		"activity_id", // PK
		"payload",
		"postmark_message_id",
		"mail_message_id",
		"reply_mail_message_id",
		"has_error",
		"submitted_at",
		"error_code",
		"postmark_message",
		"message_event",
		"acknowledged",
		"message_type",
		"created_at",
		"updated_at",
	}
}

func messageAttachmentCols() builq.Columns {
	return builq.Columns{
		"attachment_id",
		"activity_id",
		"name",
		"content_type",
		"content_key",
		"content_url",
		"spam",
		"has_error",
		"error",
		"md5_hash",
		"created_at",
		"updated_at",
	}
}

// upsertThreadTx upserts a thread within transaction.
func upsertThreadTx(ctx context.Context, tx pgx.Tx, thread *models.Thread) (*models.Thread, error) {
	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	if thread.LastInboundAt != nil {
		lastInboundAt = sql.NullTime{
			Time:  *thread.LastInboundAt,
			Valid: true,
		}
	}
	if thread.LastOutboundAt != nil {
		lastOutboundAt = sql.NullTime{
			Time:  *thread.LastOutboundAt,
			Valid: true,
		}
	}

	// Check if the thread is assigned to a member.
	// If assigned, then set the assigned member ID and assigned at for db insert values.
	// Otherwise, by default assigned member ID and assigned at will be NULL.
	if thread.AssignedMember != nil {
		assignedMemberId = sql.NullString{String: thread.AssignedMember.MemberId, Valid: true}
		assignedAt = sql.NullTime{Time: thread.AssignedMember.AssignedAt, Valid: true}
	}

	cols := threadCols()
	insertB := builq.Builder{}
	insertParams := []any{
		thread.ThreadId, thread.WorkspaceId, thread.Customer.CustomerId,
		assignedMemberId, assignedAt,
		thread.Title, thread.Description,
		thread.ThreadStatus.Status, thread.ThreadStatus.StatusChangedAt,
		thread.ThreadStatus.StatusChangedBy.MemberId,
		thread.ThreadStatus.Stage,
		thread.Replied, thread.Priority, thread.Channel,
		lastInboundAt,
		lastOutboundAt,
		thread.CreatedBy.MemberId,
		thread.UpdatedBy.MemberId,
		thread.CreatedAt,
		thread.UpdatedAt,
	}

	insertB.Addf("INSERT INTO thread (%s)", cols)
	insertB.Addf(
		"VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$)",
		insertParams...,
	)
	insertB.Addf("ON CONFLICT (thread_id) DO UPDATE SET")
	insertB.Addf("workspace_id = EXCLUDED.workspace_id,")
	insertB.Addf("customer_id = EXCLUDED.customer_id,")
	insertB.Addf("assignee_id = EXCLUDED.assignee_id,")
	insertB.Addf("assigned_at = EXCLUDED.assigned_at,")
	insertB.Addf("title = EXCLUDED.title,")
	insertB.Addf("description = EXCLUDED.description,")
	insertB.Addf("status = EXCLUDED.status,")
	insertB.Addf("status_changed_at = EXCLUDED.status_changed_at,")
	insertB.Addf("status_changed_by_id = EXCLUDED.status_changed_by_id,")
	insertB.Addf("stage = EXCLUDED.stage,")
	insertB.Addf("replied = EXCLUDED.replied,")
	insertB.Addf("priority = EXCLUDED.priority,")
	insertB.Addf("channel = EXCLUDED.channel,")
	insertB.Addf("last_inbound_at = EXCLUDED.last_inbound_at,")
	insertB.Addf("last_outbound_at = EXCLUDED.last_outbound_at,")
	insertB.Addf("updated_by_id = EXCLUDED.updated_by_id,")
	insertB.Addf("updated_at = EXCLUDED.updated_at")
	insertB.Addf("RETURNING %s", cols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build upsert query", slog.Any("err", err))
		return &models.Thread{}, ErrQuery
	}

	// Build the select query required after upsert
	q := builq.New()
	cols = threadJoinedCols()

	q("WITH ups AS (%s)", insertQuery)
	q("SELECT %s FROM %s", cols, "ups th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return &models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return &models.Thread{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to upsert query", slog.Any("err", err))
		return &models.Thread{}, ErrQuery
	}

	if assignedMemberId.Valid {
		memberActor := models.MemberActor{
			MemberId: assignedMemberId.String,
			Name:     assignedMemberName.String,
		}
		thread.AssignMember(memberActor, assignedAt.Time)
	} else {
		thread.ClearAssignedMember()
	}

	// Set last inbound and outbound time if not nil
	if lastInboundAt.Valid {
		thread.LastInboundAt = &lastInboundAt.Time
	}
	if lastOutboundAt.Valid {
		thread.LastOutboundAt = &lastOutboundAt.Time
	}
	return thread, nil
}

// insertThreadActivityTx inserts a thread activity within transaction.
func insertThreadActivityTx(ctx context.Context, tx pgx.Tx, activity *models.Activity) (*models.Activity, error) {
	var (
		customerId   sql.NullString
		customerName sql.NullString
		memberId     sql.NullString
		memberName   sql.NullString
	)

	if activity.Customer != nil {
		customerId = sql.NullString{String: activity.Customer.CustomerId, Valid: true}
	}
	if activity.Member != nil {
		memberId = sql.NullString{String: activity.Member.MemberId, Valid: true}
	}

	// Persist the message with referenced thread ID
	insertB := builq.Builder{}
	insertCols := threadActivityCols()
	insertParams := []any{
		activity.ActivityID, activity.ThreadID, activity.ActivityType,
		customerId, memberId, activity.Body, activity.CreatedAt, activity.UpdatedAt,
	}

	insertB.Addf("INSERT INTO activity (%s)", insertCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return activity, ErrQuery
	}

	// Build the select query required after insert
	q := builq.New()
	joinedCols := threadActivityJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins act")
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return activity, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&activity.ActivityID, &activity.ThreadID, &activity.ActivityType,
		&customerId, &customerName,
		&memberId, &memberName,
		&activity.Body, &activity.CreatedAt, &activity.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return activity, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return activity, ErrQuery
	}

	if customerId.Valid {
		activity.Customer = &models.CustomerActor{
			CustomerId: customerId.String,
			Name:       customerName.String,
		}
	}
	if memberId.Valid {
		activity.Member = &models.MemberActor{
			MemberId: memberId.String,
			Name:     memberName.String,
		}
	}
	return activity, nil
}

// insertPostmarkMessageLogTx inserts postmark message log within transaction.
func insertPostmarkMessageLogTx(
	ctx context.Context, tx pgx.Tx, activityID string, messageLog *models.PostmarkMessageLog,
) (*models.PostmarkMessageLog, error) {
	q := builq.New()
	logCols := postmarkMessageLogCols()
	insertParams := []any{
		activityID, messageLog.Payload, messageLog.PostmarkMessageId,
		messageLog.MailMessageId, messageLog.ReplyMailMessageId,
		messageLog.HasError, messageLog.SubmittedAt, messageLog.ErrorCode,
		messageLog.PostmarkMessage, messageLog.MessageEvent, messageLog.Acknowledged,
		messageLog.MessageType, messageLog.CreatedAt, messageLog.UpdatedAt,
	}

	q("INSERT INTO postmark_message_log (%s)", logCols)
	q("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	q("RETURNING %s", logCols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return &models.PostmarkMessageLog{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&messageLog.ActivityID, &messageLog.Payload, &messageLog.PostmarkMessageId,
		&messageLog.MailMessageId, &messageLog.ReplyMailMessageId,
		&messageLog.HasError, &messageLog.SubmittedAt, &messageLog.ErrorCode,
		&messageLog.PostmarkMessage, &messageLog.MessageEvent, &messageLog.Acknowledged,
		&messageLog.MessageType, &messageLog.CreatedAt, &messageLog.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return &models.PostmarkMessageLog{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return &models.PostmarkMessageLog{}, ErrQuery
	}
	return messageLog, nil
}

func (th *ThreadDB) SavePostmarkThreadActivity(
	ctx context.Context, thread *models.Thread, activity *models.Activity,
	postmarkMessageLog *models.PostmarkMessageLog) (*models.Thread, *models.Activity, error) {
	tx, err := th.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", slog.Any("err", err))
		return &models.Thread{}, &models.Activity{}, ErrTxQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	// 1. upsert thread
	// 2. insert new thread activity
	// 3. insert new postmark message log

	thread, err = upsertThreadTx(ctx, tx, thread)
	if err != nil {
		slog.Error("failed to upsert thread", slog.Any("err", err))
		return &models.Thread{}, &models.Activity{}, ErrQuery
	}

	activity, err = insertThreadActivityTx(ctx, tx, activity)
	if err != nil {
		slog.Error("failed to insert thread activity", slog.Any("err", err))
		return &models.Thread{}, &models.Activity{}, ErrQuery
	}

	_, err = insertPostmarkMessageLogTx(ctx, tx, activity.ActivityID, postmarkMessageLog)
	if err != nil {
		slog.Error("failed to insert postmark message log", slog.Any("err", err))
		return &models.Thread{}, &models.Activity{}, ErrQuery
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return &models.Thread{}, &models.Activity{}, ErrTxQuery
	}
	return thread, activity, nil
}

// ModifyThreadById modifies thread for selective fields
func (th *ThreadDB) ModifyThreadById(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {
	upsertQ := builq.New()
	upsertParams := make([]any, 0, len(fields)+1) // updates + thread ID
	cols := threadCols()

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	upsertQ("UPDATE thread SET")
	for _, field := range fields {
		switch field {
		case "priority":
			upsertQ("priority = %$,", thread.Priority)
			upsertParams = append(upsertParams, thread.Priority)
		case "assignee":
			if thread.AssignedMember != nil {
				assignedMemberId = sql.NullString{String: thread.AssignedMember.MemberId, Valid: true}
			}
			upsertQ("assignee_id = %$,", assignedMemberId)
			upsertParams = append(upsertParams, assignedMemberId)
		case "stage":
			upsertQ("stage = %$,", thread.ThreadStatus.Stage)
			upsertQ("status = %$,", thread.ThreadStatus.Status)
			upsertQ("status_changed_at = %$,", thread.ThreadStatus.StatusChangedAt)
			upsertQ("status_changed_by_id = %$,", thread.ThreadStatus.StatusChangedBy.MemberId)
			upsertParams = append(upsertParams, thread.ThreadStatus.Stage)
			upsertParams = append(upsertParams, thread.ThreadStatus.Status)
			upsertParams = append(upsertParams, thread.ThreadStatus.StatusChangedAt)
			upsertParams = append(upsertParams, thread.ThreadStatus.StatusChangedBy.MemberId)
		}
	}

	upsertQ("updated_at = NOW()")
	upsertQ("WHERE thread_id = %$", thread.ThreadId)
	upsertParams = append(upsertParams, thread.ThreadId)

	upsertQ("RETURNING %s", cols)

	stmt, _, err := upsertQ.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	q := builq.New()
	joinedCols := threadJoinedCols()

	q("WITH ups AS (%s)", stmt)
	q("SELECT %s FROM %s", joinedCols, "ups th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = th.db.QueryRow(ctx, stmt, upsertParams...).Scan(
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	// Sets the assigned member if a valid assigned member exists,
	// otherwise clears the assigned member.
	if assignedMemberId.Valid {
		memberActor := models.MemberActor{
			MemberId: assignedMemberId.String,
			Name:     assignedMemberName.String,
		}
		thread.AssignMember(memberActor, assignedAt.Time)
	} else {
		thread.ClearAssignedMember()
	}

	// Set last inbound and outbound time if not nil
	if lastInboundAt.Valid {
		thread.LastInboundAt = &lastInboundAt.Time
	}
	if lastOutboundAt.Valid {
		thread.LastOutboundAt = &lastOutboundAt.Time
	}
	return thread, nil
}

// LookupByWorkspaceThreadId looks up thread in workspace for provided thread ID.
func (th *ThreadDB) LookupByWorkspaceThreadId(
	ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error) {
	var thread models.Thread

	params := []any{workspaceId, threadId}
	cols := threadJoinedCols()
	q := builq.New()
	q("SELECT %s FROM %s", cols, "thread th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	q("WHERE th.workspace_id = %$ AND th.thread_id = %$", workspaceId, threadId)
	if channel != nil {
		q("AND th.channel = %$", *channel)
		params = append(params, *channel)
	}

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	err = th.db.QueryRow(ctx, stmt, params...).Scan(
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	// Sets the assigned member if a valid assigned member exists,
	// otherwise clears the assigned member.
	if assignedMemberId.Valid {
		memberActor := models.MemberActor{
			MemberId: assignedMemberId.String,
			Name:     assignedMemberName.String,
		}
		thread.AssignMember(memberActor, assignedAt.Time)
	} else {
		thread.ClearAssignedMember()
	}

	// Sets the assigned member if a valid assigned member exists,
	// otherwise clears the assigned member.
	if assignedMemberId.Valid {
		memberActor := models.MemberActor{
			MemberId: assignedMemberId.String,
			Name:     assignedMemberName.String,
		}
		thread.AssignMember(memberActor, assignedAt.Time)
	} else {
		thread.ClearAssignedMember()
	}

	// Set last inbound and outbound time if not nil
	if lastInboundAt.Valid {
		thread.LastInboundAt = &lastInboundAt.Time
	}
	if lastOutboundAt.Valid {
		thread.LastOutboundAt = &lastOutboundAt.Time
	}
	return thread, nil
}

func (th *ThreadDB) FetchThreadsByWorkspaceId(
	ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	limit := 100
	threads := make([]models.Thread, 0, limit)

	params := []any{workspaceId}
	cols := threadJoinedCols()
	q := builq.New()
	q("SELECT %s FROM %s", cols, "thread th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	q("WHERE th.workspace_id = %$", workspaceId)
	if channel != nil {
		q("AND th.channel = %$", *channel)
		params = append(params, *channel)
	}
	if role != nil {
		q("AND c.role = %$", *role)
		params = append(params, *role)
	}
	// ignore threads from visitors
	q("AND c.role <> %$", models.Customer{}.Visitor())
	params = append(params, models.Customer{}.Visitor())

	// Sort by earliest created threads.
	q("ORDER BY th.created_at ASC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	rows, _ := th.db.Query(ctx, stmt, params...)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// Sets the assigned member if a valid assigned member exists,
		// otherwise clears the assigned member.
		if assignedMemberId.Valid {
			memberActor := models.MemberActor{
				MemberId: assignedMemberId.String,
				Name:     assignedMemberName.String,
			}
			thread.AssignMember(memberActor, assignedAt.Time)
		} else {
			thread.ClearAssignedMember()
		}

		// Set last inbound and outbound time if not nil
		if lastInboundAt.Valid {
			thread.LastInboundAt = &lastInboundAt.Time
		}
		if lastOutboundAt.Valid {
			thread.LastOutboundAt = &lastOutboundAt.Time
		}
		threads = append(threads, thread)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Thread{}, ErrQuery
	}
	return threads, nil
}

func (th *ThreadDB) FetchThreadsByAssignedMemberId(
	ctx context.Context, memberId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	limit := 100
	threads := make([]models.Thread, 0, limit)

	params := []any{memberId}
	cols := threadJoinedCols()
	q := builq.New()
	q("SELECT %s FROM %s", cols, "thread th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	q("WHERE th.assignee_id = %$", memberId)
	if channel != nil {
		q("AND th.channel = %$", *channel)
		params = append(params, *channel)
	}
	if role != nil {
		q("AND c.role = %$", *role)
		params = append(params, *role)
	}
	// ignore threads from visitors
	q("AND c.role <> %$", models.Customer{}.Visitor())
	params = append(params, models.Customer{}.Visitor())

	// Sort by earliest created threads.
	q("ORDER BY th.created_at ASC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	rows, _ := th.db.Query(ctx, stmt, params...)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// Sets the assigned member if a valid assigned member exists,
		// otherwise clears the assigned member.
		if assignedMemberId.Valid {
			memberActor := models.MemberActor{
				MemberId: assignedMemberId.String,
				Name:     assignedMemberName.String,
			}
			thread.AssignMember(memberActor, assignedAt.Time)
		} else {
			thread.ClearAssignedMember()
		}

		// Set last inbound and outbound time if not nil
		if lastInboundAt.Valid {
			thread.LastInboundAt = &lastInboundAt.Time
		}
		if lastOutboundAt.Valid {
			thread.LastOutboundAt = &lastOutboundAt.Time
		}
		threads = append(threads, thread)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.Thread{}, ErrQuery
	}
	return threads, nil
}

func (th *ThreadDB) FetchThreadsByMemberUnassigned(
	ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	limit := 100
	threads := make([]models.Thread, 0, limit)

	params := []any{workspaceId}
	cols := threadJoinedCols()
	q := builq.New()
	q("SELECT %s FROM %s", cols, "thread th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	q("WHERE th.workspace_id = %$", workspaceId)
	q("AND th.assignee_id IS NULL")
	if channel != nil {
		q("AND th.channel = %$", *channel)
		params = append(params, *channel)
	}
	if role != nil {
		q("AND c.role = %$", *role)
		params = append(params, *role)
	}
	// ignore threads from visitors
	q("AND c.role <> %$", models.Customer{}.Visitor())
	params = append(params, models.Customer{}.Visitor())

	// Sort by earliest created threads.
	q("ORDER BY th.created_at ASC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	rows, _ := th.db.Query(ctx, stmt, params...)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// Sets the assigned member if a valid assigned member exists,
		// otherwise clears the assigned member.
		if assignedMemberId.Valid {
			memberActor := models.MemberActor{
				MemberId: assignedMemberId.String,
				Name:     assignedMemberName.String,
			}
			thread.AssignMember(memberActor, assignedAt.Time)
		} else {
			thread.ClearAssignedMember()
		}
		// Sets the assigned member if a valid assigned member exists,
		// otherwise clears the assigned member.
		if assignedMemberId.Valid {
			memberActor := models.MemberActor{
				MemberId: assignedMemberId.String,
				Name:     assignedMemberName.String,
			}
			thread.AssignMember(memberActor, assignedAt.Time)
		} else {
			thread.ClearAssignedMember()
		}

		// Set last inbound and outbound time if not nil
		if lastInboundAt.Valid {
			thread.LastInboundAt = &lastInboundAt.Time
		}
		if lastOutboundAt.Valid {
			thread.LastOutboundAt = &lastOutboundAt.Time
		}
		threads = append(threads, thread)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.Thread{}, ErrQuery
	}
	return threads, nil
}

func (th *ThreadDB) FetchThreadsByLabelId(
	ctx context.Context, labelId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	limit := 100
	threads := make([]models.Thread, 0, limit)

	params := []any{labelId}
	cols := threadJoinedCols()
	q := builq.New()
	q("SELECT %s FROM %s", cols, "thread th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")
	q("INNER JOIN thread_label tl ON th.thread_id = tl.thread_id")

	q("WHERE tl.label_id = %$", labelId)
	if channel != nil {
		q("AND th.channel = %$", *channel)
		params = append(params, *channel)
	}
	if role != nil {
		q("AND c.role = %$", *role)
		params = append(params, *role)
	}
	// ignore threads from visitors
	q("AND c.role <> %$", models.Customer{}.Visitor())
	params = append(params, models.Customer{}.Visitor())

	// Sort by earliest created threads.
	q("ORDER BY th.created_at ASC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return []models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	rows, _ := th.db.Query(ctx, stmt, params...)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// Sets the assigned member if a valid assigned member exists,
		// otherwise clears the assigned member.
		if assignedMemberId.Valid {
			memberActor := models.MemberActor{
				MemberId: assignedMemberId.String,
				Name:     assignedMemberName.String,
			}
			thread.AssignMember(memberActor, assignedAt.Time)
		} else {
			thread.ClearAssignedMember()
		}
		// Sets the assigned member if a valid assigned member exists,
		// otherwise clears the assigned member.
		if assignedMemberId.Valid {
			memberActor := models.MemberActor{
				MemberId: assignedMemberId.String,
				Name:     assignedMemberName.String,
			}
			thread.AssignMember(memberActor, assignedAt.Time)
		} else {
			thread.ClearAssignedMember()
		}

		// Set last inbound and outbound time if not nil
		if lastInboundAt.Valid {
			thread.LastInboundAt = &lastInboundAt.Time
		}
		if lastOutboundAt.Valid {
			thread.LastOutboundAt = &lastOutboundAt.Time
		}
		threads = append(threads, thread)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return []models.Thread{}, ErrQuery
	}
	return threads, nil
}

func (th *ThreadDB) CheckThreadInWorkspaceExists(
	ctx context.Context, workspaceId string, threadId string) (bool, error) {
	var isExist bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM thread
		WHERE workspace_id = $1 AND thread_id= $2
	)`

	err := th.db.QueryRow(ctx, stmt, workspaceId, threadId).Scan(&isExist)
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return false, ErrQuery
	}
	return isExist, nil
}

func (th *ThreadDB) SetThreadLabel(
	ctx context.Context, threadLabel models.ThreadLabel) (models.ThreadLabel, bool, error) {
	var IsCreated bool
	thLabelId := threadLabel.GenId()
	stmt := `
		WITH ins AS (
			INSERT INTO thread_label (thread_label_id, thread_id, label_id, addedby)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (thread_id, label_id) DO NOTHING
			RETURNING thread_label_id, thread_id, label_id, addedby,
			created_at, updated_at, TRUE AS is_created
		)
		SELECT
			ins.thread_label_id,
			ins.thread_id,
			ins.label_id,
			l.name,
			l.icon,
			ins.addedby,
			ins.created_at,
			ins.updated_at,
			ins.is_created
		FROM ins
		JOIN label l ON ins.label_id = l.label_id
		UNION ALL
		SELECT
			tl.thread_label_id,
			tl.thread_id,
			tl.label_id,
			l.name,
			l.icon,
			tl.addedby,
			tl.created_at,
			tl.updated_at,
			FALSE AS is_created
		FROM thread_label tl
		INNER JOIN label l ON tl.label_id = l.label_id
		WHERE tl.thread_id = $2 AND l.label_id = $3 AND NOT EXISTS (SELECT 1 FROM ins)
	`

	err := th.db.QueryRow(ctx, stmt, thLabelId, threadLabel.ThreadId, threadLabel.LabelId, threadLabel.AddedBy).Scan(
		&threadLabel.ThreadLabelId, &threadLabel.ThreadId, &threadLabel.LabelId,
		&threadLabel.Name, &threadLabel.Icon, &threadLabel.AddedBy,
		&threadLabel.CreatedAt, &threadLabel.UpdatedAt, &IsCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.ThreadLabel{}, IsCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.ThreadLabel{}, IsCreated, ErrQuery
	}
	return threadLabel, IsCreated, nil
}

func (th *ThreadDB) DeleteThreadLabelById(
	ctx context.Context, threadId string, labelId string) error {

	q := builq.New()
	q("DELETE FROM thread_label")
	q("WHERE thread_id = %$ AND label_id = %$", threadId, labelId)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	_, err = th.db.Exec(ctx, stmt, threadId, labelId)
	if err != nil {
		slog.Error("failed to delete query", slog.Any("err", err))
		return ErrQuery
	}
	return nil
}

func (th *ThreadDB) FetchAttachedLabelsByThreadId(
	ctx context.Context, threadId string) ([]models.ThreadLabel, error) {
	var label models.ThreadLabel
	labels := make([]models.ThreadLabel, 0, 100)
	stmt := `SELECT tl.thread_label_id AS thread_label_id,
		tl.thread_id AS thread_id,
		tl.label_id AS label_id,
		l.name AS name, l.icon AS icon,
		tl.addedby AS addedby,
		tl.created_at AS created_at,
		tl.updated_at AS updated_at
		FROM thread_label AS tl
		INNER JOIN label AS l ON tl.label_id = l.label_id
		WHERE tl.thread_id = $1
		ORDER BY tl.created_at DESC LIMIT 100`

	rows, _ := th.db.Query(ctx, stmt, threadId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&label.ThreadLabelId, &label.ThreadId, &label.LabelId,
		&label.Name, &label.Icon, &label.AddedBy,
		&label.CreatedAt, &label.UpdatedAt,
	}, func() error {
		labels = append(labels, label)
		return nil
	})

	if err != nil {
		slog.Error("failed to scan", slog.Any("err", err))
		return []models.ThreadLabel{}, ErrQuery
	}

	return labels, nil
}

func (th *ThreadDB) FetchMessagesByThreadId(
	ctx context.Context, threadId string) ([]models.Activity, error) {
	var activity models.Activity
	activities := make([]models.Activity, 0, 100)

	var (
		customerId   sql.NullString
		customerName sql.NullString
		memberId     sql.NullString
		memberName   sql.NullString
	)

	q := builq.New()
	activitiesJoinedCols := threadActivityJoinedCols()
	q("SELECT %s FROM activity act", activitiesJoinedCols)
	q("LEFT OUTER JOIN customer c ON act.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON act.member_id = m.member_id")
	q("WHERE act.thread_id = %$ AND activity_type = %$", threadId, models.ActivityThreadMessage)

	q("ORDER BY act.created_at ASC")
	q("LIMIT 100")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return nil, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	rows, _ := th.db.Query(ctx, stmt, threadId)
	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&activity.ActivityID, &activity.ThreadID, &activity.ActivityType,
		&customerId, &customerName,
		&memberId, &memberName,
		&activity.Body,
		&activity.CreatedAt, &activity.UpdatedAt,
	}, func() error {
		if customerId.Valid {
			activity.Customer = &models.CustomerActor{
				CustomerId: customerId.String,
				Name:       customerName.String,
			}
		}
		if memberId.Valid {
			activity.Member = &models.MemberActor{
				MemberId: memberId.String,
				Name:     memberName.String,
			}
		}
		activities = append(activities, activity)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Activity{}, ErrQuery
	}
	return activities, nil
}

func (th *ThreadDB) FetchMessagesWithAttachmentsByThreadId(
	ctx context.Context, threadId string) ([]models.ActivityWithAttachments, error) {
	var activity models.ActivityWithAttachments
	limit := 100
	activities := make([]models.ActivityWithAttachments, 0, limit)

	var (
		customerId   sql.NullString
		customerName sql.NullString
		memberId     sql.NullString
		memberName   sql.NullString
	)

	cols := threadActivityJoinedCols()
	stmt := `SELECT
		%s,
		COALESCE(
			(
				SELECT JSON_AGG(
					JSON_BUILD_OBJECT(
						'attachmentId', aa.attachment_id,
						'activityId', aa.activity_id,
						'name', aa.name,
						'contentType', aa.content_type,
						'contentKey', aa.content_key,
						'contentUrl', aa.content_url,
						'spam', aa.spam,
						'hasError', aa.has_error,
						'error', aa.error,
						'md5Hash', aa.md5_hash,
						'createdAt', aa.created_at AT TIME ZONE 'UTC',
						'updatedAt', aa.updated_at AT TIME ZONE 'UTC'
					)
				)
				FROM (
					SELECT *
					FROM activity_attachment
					WHERE activity_id = act.activity_id
					ORDER BY created_at ASC
					LIMIT 10
				) aa
			), 
			'[]'::json
		) as attachments
	FROM activity act
	LEFT OUTER JOIN customer c ON act.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON act.member_id = m.member_id`

	stmt = fmt.Sprintf(stmt, cols)

	q := builq.New()
	q("%s", stmt)
	q("WHERE act.thread_id = %$ AND act.activity_type = %$", threadId, models.ActivityThreadMessage)
	q("ORDER BY act.created_at ASC")
	q("LIMIT %d", limit)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return nil, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var attachmentsJson []byte

	rows, _ := th.db.Query(ctx, stmt, threadId)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&activity.ActivityID, &activity.ThreadID, &activity.ActivityType,
		&customerId, &customerName,
		&memberId, &memberName,
		&activity.Body,
		&activity.CreatedAt, &activity.UpdatedAt,
		&attachmentsJson,
	}, func() error {
		if customerId.Valid {
			activity.Customer = &models.CustomerActor{
				CustomerId: customerId.String,
				Name:       customerName.String,
			}
		}
		if memberId.Valid {
			activity.Member = &models.MemberActor{
				MemberId: memberId.String,
			}
		}
		var attachments []models.ActivityAttachment
		if len(attachmentsJson) > 0 {
			if err := json.Unmarshal(attachmentsJson, &attachments); err != nil {
				slog.Error("failed to unmarshal attachments",
					slog.String("activityID", activity.ActivityID))
				slog.Any("err", err)
				return err
			}
			activity.Attachments = attachments
		} else {
			activity.Attachments = []models.ActivityAttachment{}
		}
		activities = append(activities, activity)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.ActivityWithAttachments{}, ErrQuery
	}
	return activities, nil
}

func (th *ThreadDB) ComputeStatusMetricsByWorkspaceId(
	ctx context.Context, workspaceId string) (models.ThreadMetrics, error) {
	var metrics models.ThreadMetrics
	stmt := `SELECT
		COALESCE(SUM(CASE WHEN status = 'todo' THEN 1 ELSE 0 END), 0) AS active_count,
		COALESCE(SUM(CASE WHEN status = 'todo' AND stage = 'needs_first_response' THEN 1 ELSE 0 END), 0)
		AS needs_first_response_count,
		COALESCE(SUM(CASE WHEN status = 'todo' AND stage = 'waiting_on_customer' THEN 1 ELSE 0 END), 0)
		AS waiting_on_customer_count,
		COALESCE(SUM(CASE WHEN status = 'todo' AND stage = 'hold' THEN 1 ELSE 0 END), 0)
		AS hold_count,
		COALESCE(SUM(CASE WHEN status = 'todo' AND stage = 'needs_next_response' THEN 1 ELSE 0 END), 0)
		AS needs_next_response_count
	FROM 
		thread th
	INNER JOIN customer c ON th.customer_id = c.customer_id
	WHERE 
		th.workspace_id = $1 AND c.role <> 'visitor'`

	err := th.db.QueryRow(ctx, stmt, workspaceId).Scan(
		&metrics.ActiveCount, &metrics.NeedsFirstResponseCount,
		&metrics.WaitingOnCustomerCount, &metrics.HoldCount,
		&metrics.NeedsNextResponseCount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.ThreadMetrics{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.ThreadMetrics{}, ErrQuery
	}
	return metrics, nil
}

func (th *ThreadDB) ComputeAssigneeMetricsByMember(
	ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error) {
	var metrics models.ThreadAssigneeMetrics
	stmt := `SELECT
			COALESCE(SUM(
				CASE WHEN assignee_id = $2 AND status = 'todo' THEN 1 ELSE 0 END), 0) AS member_assigned_count,
			COALESCE(SUM(
				CASE WHEN assignee_id IS NULL AND status = 'todo' THEN 1 ELSE 0 END), 0) AS unassigned_count,
			COALESCE(SUM(
				CASE WHEN assignee_id IS NOT NULL AND assignee_id <> $2 AND status = 'todo' THEN 1 ELSE 0 END), 0) AS other_assigned_count
		FROM
			thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		WHERE
			th.workspace_id = $1 AND c.role <> 'visitor'`

	err := th.db.QueryRow(ctx, stmt, workspaceId, memberId).Scan(
		&metrics.MeCount, &metrics.UnAssignedCount, &metrics.OtherAssignedCount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.ThreadAssigneeMetrics{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.ThreadAssigneeMetrics{}, ErrQuery
	}
	return metrics, nil
}

func (th *ThreadDB) ComputeLabelMetricsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error) {
	var metric models.ThreadLabelMetric
	metrics := make([]models.ThreadLabelMetric, 0, 100)

	stmt := `SELECT l.label_id,
			l.name AS label_name, l.icon AS label_icon,
			COUNT(tl.thread_id) AS count
		FROM
			label l
		LEFT JOIN
			thread_label tl ON l.label_id = tl.label_id
		WHERE
			l.workspace_id = $1
		GROUP BY
			l.label_id, l.name
		ORDER BY MAX(tl.updated_at) DESC
		LIMIT 100
	`

	rows, _ := th.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&metric.LabelId, &metric.Name, &metric.Icon, &metric.Count,
	}, func() error {
		metrics = append(metrics, metric)
		return nil
	})

	if err != nil {
		slog.Error("failed to scan", slog.Any("err", err))
		return []models.ThreadLabelMetric{}, ErrQuery
	}

	return metrics, nil
}

func (th *ThreadDB) FindThreadByPostmarkReplyMessageId(
	ctx context.Context, workspaceId string, mailMessageId string) (models.Thread, error) {
	var thread models.Thread

	var selectB builq.Builder
	selectB.Addf("SELECT act.thread_id AS thread_id")
	selectB.Addf("FROM postmark_message_log p")
	selectB.Addf("INNER JOIN activity act ON p.activity_id = act.activity_id")
	selectB.Addf("WHERE p.mail_message_id = $2")

	selectQuery, _, err := selectB.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	q := builq.New()
	cols := threadJoinedCols()
	q("WITH cte AS (%s)", selectQuery)
	q("SELECT %s FROM %s", cols, "thread th")

	q("INNER JOIN cte cte ON cte.thread_id = th.thread_id")

	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	q("WHERE th.workspace_id = $1")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var (
		assignedMemberId   sql.NullString
		assignedMemberName sql.NullString
		assignedAt         sql.NullTime
		lastInboundAt      sql.NullTime
		lastOutboundAt     sql.NullTime
	)

	err = th.db.QueryRow(ctx, stmt, workspaceId, mailMessageId).Scan(
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&lastInboundAt, &lastOutboundAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	// Sets the assigned member if a valid assigned member exists,
	// otherwise clears the assigned member.
	if assignedMemberId.Valid {
		memberActor := models.MemberActor{
			MemberId: assignedMemberId.String,
			Name:     assignedMemberName.String,
		}
		thread.AssignMember(memberActor, assignedAt.Time)
	} else {
		thread.ClearAssignedMember()
	}

	// Set last inbound and outbound time if not nil
	if lastInboundAt.Valid {
		thread.LastInboundAt = &lastInboundAt.Time
	}
	if lastOutboundAt.Valid {
		thread.LastOutboundAt = &lastOutboundAt.Time
	}
	return thread, nil
}

func (th *ThreadDB) CheckPostmarkInboundExists(ctx context.Context, pmMessageId string) (bool, error) {
	var isExist bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM postmark_message_log WHERE postmark_message_id = $1 AND message_type = 'inbound'
	)`

	err := th.db.QueryRow(ctx, stmt, pmMessageId).Scan(&isExist)
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return false, ErrQuery
	}
	return isExist, nil
}

func (th *ThreadDB) InsertMessageAttachment(
	ctx context.Context, attachment models.ActivityAttachment) (models.ActivityAttachment, error) {
	cols := messageAttachmentCols()
	q := builq.New()
	insertParams := []any{
		attachment.AttachmentId, attachment.ActivityID, attachment.Name,
		attachment.ContentType, attachment.ContentKey, attachment.ContentUrl,
		attachment.Spam, attachment.HasError, attachment.Error, attachment.MD5Hash,
		attachment.CreatedAt, attachment.UpdatedAt,
	}
	q("INSERT INTO activity_attachment (%s)", cols)
	q("VALUES (%+$)", insertParams)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.ActivityAttachment{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = th.db.QueryRow(ctx, stmt, insertParams...).Scan(
		&attachment.AttachmentId, &attachment.ActivityID, &attachment.Name,
		&attachment.ContentType, &attachment.ContentKey, &attachment.ContentUrl,
		&attachment.Spam, &attachment.HasError, &attachment.Error, &attachment.MD5Hash,
		&attachment.CreatedAt, &attachment.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.ActivityAttachment{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.ActivityAttachment{}, ErrQuery
	}
	return attachment, nil
}

func (th *ThreadDB) FetchMessageAttachmentById(
	ctx context.Context, messageId, attachmentId string) (models.ActivityAttachment, error) {
	var attachment models.ActivityAttachment
	cols := messageAttachmentCols()

	q := builq.New()
	q("SELECT %s FROM message_attachment", cols)
	q("WHERE activity_id = %$ AND attachment_id = %$", messageId, attachmentId)

	stmt, _, err := q.Build()

	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.ActivityAttachment{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = th.db.QueryRow(ctx, stmt, messageId, attachmentId).Scan(
		&attachment.AttachmentId, &attachment.ActivityID, &attachment.Name,
		&attachment.ContentType, &attachment.ContentKey, &attachment.ContentUrl,
		&attachment.Spam, &attachment.HasError, &attachment.Error, &attachment.MD5Hash,
		&attachment.CreatedAt, &attachment.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.ActivityAttachment{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.ActivityAttachment{}, ErrQuery
	}
	return attachment, nil
}

func (th *ThreadDB) GetRecentMailMessageIdByThreadId(ctx context.Context, threadId string) (string, error) {
	var mailMessageId string
	q := builq.New()
	q("SELECT pml.mail_message_id FROM postmark_message_log pml")
	q("INNER JOIN activity act ON act.activity_id = pml.activity_id")
	q("WHERE act.thread_id = %$", threadId)
	q("ORDER BY act.created_at DESC")
	q("LIMIT 1")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return "", ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = th.db.QueryRow(ctx, stmt, threadId).Scan(&mailMessageId)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return "", ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return "", ErrQuery
	}
	return mailMessageId, nil
}
