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
		"inbound_message_id",  // FK Nullable to inbound_message
		"outbound_message_id", // FK Nullable to outbound_message
		"created_by_id",       // FK to member
		"updated_by_id",       // FK to member
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
		"inb.message_id",
		"inbc.customer_id",
		"inbc.name",
		"inb.preview_text",
		"inb.first_seq_id",
		"inb.last_seq_id",
		"inb.created_at",
		"inb.updated_at",
		"oub.message_id",
		"oubm.member_id",
		"oubm.name",
		"oub.preview_text",
		"oub.first_seq_id",
		"oub.last_seq_id",
		"oub.created_at",
		"oub.updated_at",
		"mc.member_id",
		"mc.name",
		"mu.member_id",
		"mu.name",
		"th.created_at",
		"th.updated_at",
	}
}

// Returns the required columns for the inbound message table.
// The order of the columns matters when returning the results.
func inboundMessageCols() builq.Columns {
	return builq.Columns{
		"message_id", "customer_id",
		"preview_text",
		"first_seq_id", "last_seq_id",
		"created_at", "updated_at",
	}
}

// Returns the required columns for the outbound message table.
// The order of the columns matters when returning the results.
func outboundMessageCols() builq.Columns {
	return builq.Columns{
		"message_id", "member_id",
		"preview_text",
		"first_seq_id", "last_seq_id",
		"created_at", "updated_at",
	}
}

// Returns the required columns and joined for the inbound message table.
// The order of the columns matters when returning the results.
func inboundMessageJoinedCols() builq.Columns {
	cols := builq.Columns{
		"im.message_id",
		"c.customer_id",
		"c.name",
		"im.preview_text",
		"im.first_seq_id",
		"im.last_seq_id",
		"im.created_at",
		"im.updated_at",
	}
	return cols
}

// Returns the required columns and joined for the inbound message table.
// The order of the columns matters when returning the results.
func outboundMessageJoinedCols() builq.Columns {
	cols := builq.Columns{
		"om.message_id",
		"m.member_id",
		"m.name",
		"om.preview_text",
		"om.first_seq_id",
		"om.last_seq_id",
		"om.created_at",
		"om.updated_at",
	}
	return cols
}

func threadMessageCols() builq.Columns {
	return builq.Columns{
		"message_id", // PK
		"thread_id",  // FK to thread
		"text_body",
		"body",
		"customer_id", // FK Nullable to customer
		"member_id",   // FK Nullable to member
		"channel",
		"created_at",
		"updated_at",
	}
}

func threadMessageJoinedCols() builq.Columns {
	return builq.Columns{
		"msg.message_id",
		"msg.thread_id",
		"msg.text_body",
		"msg.body",
		"c.customer_id",
		"c.name",
		"m.member_id",
		"m.name",
		"msg.channel",
		"msg.created_at",
		"msg.updated_at",
	}
}

func postmarkInboundMessageCols() builq.Columns {
	return builq.Columns{
		"message_id", // PK
		"payload",
		"postmark_message_id",
		"mail_message_id",
		"reply_mail_message_id",
		"created_at",
		"updated_at",
	}
}

func messageAttachmentCols() builq.Columns {
	return builq.Columns{
		"attachment_id",
		"message_id",
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

// InsertInboundThreadMessage inserts a new inbound thread message for the customer in a transaction.
// First, insert the inbound message.
// Then, insert the Thread with in persisted inbound message ID.
// Finally, insert the chat with in persisted thread ID.
//
// The IDs are already generated within the time space.
func (th *ThreadDB) InsertInboundThreadMessage(
	ctx context.Context, inbound models.ThreadMessage) (models.Thread, models.Message, error) {
	// start transaction
	// If fails then stop the execution and return the error.
	tx, err := th.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	// insert thread linked inbound message
	// inserts thread
	// inserts message
	// in that order.

	thread := inbound.Thread
	message := inbound.Message

	// Checks if the thread has an inbound message.
	// If not, adding inbound thread is not allowed.
	if thread.InboundMessage == nil {
		slog.Error("thread inbound message cannot be empty", slog.Any("thread", thread))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	inboundMessage := thread.InboundMessage
	// Persist the inbound message.
	// Do insert an inbound message first before inserting thread.
	// Thread will reference the inbound message ID.
	var insertB builq.Builder
	messageCols := inboundMessageCols()
	insertParams := []any{
		inboundMessage.MessageId, inboundMessage.Customer.CustomerId,
		inboundMessage.PreviewText,
		inboundMessage.FirstSeqId, inboundMessage.LastSeqId,
		inboundMessage.CreatedAt, inboundMessage.UpdatedAt,
	}

	// Build the insert query.
	insertB.Addf("INSERT INTO inbound_message (%s)", messageCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", messageCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	// Build the select query required after insert.
	q := builq.New()
	messageJoinedCols := inboundMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", messageJoinedCols, "ins im") // inserted table, with alias im
	q("INNER JOIN customer c ON im.customer_id = c.customer_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	// Make the insert query.
	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&inboundMessage.MessageId, &inboundMessage.Customer.CustomerId, &inboundMessage.Customer.Name,
		&inboundMessage.PreviewText, &inboundMessage.FirstSeqId, &inboundMessage.LastSeqId,
		&inboundMessage.CreatedAt, &inboundMessage.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	// hold db values for nullables.
	var (
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
	)

	// Persisted inbound message ID.
	inboundMessageId = sql.NullString{String: inboundMessage.MessageId, Valid: true}

	// Check if the thread is assigned to a member.
	// If assigned, then set the assigned member ID and assigned at for db insert values.
	// Otherwise, by default assigned member ID and assigned at will be NULL.
	if thread.AssignedMember != nil {
		assignedMemberId = sql.NullString{String: thread.AssignedMember.MemberId, Valid: true}
		assignedAt = sql.NullTime{Time: thread.AssignedMember.AssignedAt, Valid: true}
	}

	// Persist the thread with referenced inbound message ID.
	insertB = builq.Builder{}
	messageCols = threadCols()
	insertParams = []any{
		thread.ThreadId, thread.WorkspaceId, thread.Customer.CustomerId,
		assignedMemberId, assignedAt,
		thread.Title, thread.Description,
		thread.ThreadStatus.Status, thread.ThreadStatus.StatusChangedAt,
		thread.ThreadStatus.StatusChangedBy.MemberId,
		thread.ThreadStatus.Stage,
		thread.Replied, thread.Priority, thread.Channel,
		inboundMessageId,
		outboundMessageId,
		thread.CreatedBy.MemberId,
		thread.UpdatedBy.MemberId,
		thread.CreatedAt,
		thread.UpdatedAt,
	}

	insertB.Addf("INSERT INTO thread (%s)", messageCols)
	insertB.Addf(
		"VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$)",
		insertParams...,
	)
	insertB.Addf("RETURNING %s", messageCols)

	insertQuery, _, err = insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	// Build the select query required after insert.
	q = builq.New()
	messageJoinedCols = threadJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", messageJoinedCols, "ins th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
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

	// Sets the inbound message if a valid inbound message exists,
	// otherwise clears the inbound message.
	if inboundMessageId.Valid {
		customer := models.CustomerActor{
			CustomerId: inboundCustomerId.String,
			Name:       inboundCustomerName.String,
		}
		thread.InboundMessage = &models.InboundMessage{
			MessageId:   inboundMessageId.String,
			Customer:    customer,
			PreviewText: inboundPreviewText.String,
			FirstSeqId:  inboundFirstSeqId.String,
			LastSeqId:   inboundLastSeqId.String,
			CreatedAt:   inboundCreatedAt.Time,
			UpdatedAt:   inboundUpdatedAt.Time,
		}
	} else {
		thread.ClearInboundMessage()
	}

	// Sets the outbound message if a valid outbound message exists,
	// otherwise clears the outbound message.
	if outboundMessageId.Valid {
		member := models.MemberActor{
			MemberId: outboundMemberId.String,
			Name:     outboundMemberName.String,
		}
		thread.OutboundMessage = &models.OutboundMessage{
			MessageId:   outboundMessageId.String,
			Member:      member,
			PreviewText: outboundPreviewText.String,
			FirstSeqId:  outboundFirstSeqId.String,
			LastSeqId:   outboundLastSeqId.String,
			CreatedAt:   outboundCreatedAt.Time,
			UpdatedAt:   outboundUpdatedAt.Time,
		}
	} else {
		thread.ClearOutboundMessage()
	}

	// hold db nullables.
	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString
	if message.Customer != nil {
		customerId = sql.NullString{String: message.Customer.CustomerId, Valid: true}
	}
	if message.Member != nil {
		memberId = sql.NullString{String: message.Member.MemberId, Valid: true}
	}

	// Persist the message with referenced thread ID
	insertB = builq.Builder{}
	messageCols = threadMessageCols()
	insertParams = []any{
		message.MessageId, message.ThreadId, message.TextBody, message.Body,
		customerId, memberId, message.Channel, message.CreatedAt, message.UpdatedAt,
	}

	insertB.Addf("INSERT INTO message (%s)", messageCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", messageCols)

	insertQuery, _, err = insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	// Build the select query required after insert
	q = builq.New()
	messageJoinedCols = threadMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", messageJoinedCols, "ins msg")
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&message.MessageId, &message.ThreadId, &message.TextBody, &message.Body,
		&customerId, &customerName,
		&memberId, &memberName,
		&message.Channel, &message.CreatedAt, &message.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	if customerId.Valid {
		message.Customer = &models.CustomerActor{
			CustomerId: customerId.String,
			Name:       customerName.String,
		}
	}
	if memberId.Valid {
		message.Member = &models.MemberActor{
			MemberId: memberId.String,
			Name:     memberName.String,
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrTxQuery
	}
	return *thread, *message, nil
}

func InsertInboundThreadTx(
	ctx context.Context, tx pgx.Tx, thread *models.Thread) (*models.Thread, error) {
	// Make sure InboundMessage exists for inbound Thread.
	if thread.InboundMessage == nil {
		slog.Error("thread inbound message cannot be empty", slog.Any("thread", thread))
		return thread, ErrQuery
	}

	inboundMessage := thread.InboundMessage
	// Persist the inbound message.
	// Do insert an inbound message first before inserting thread.
	// Thread will reference the inbound message ID.
	var insertB builq.Builder
	insertCols := inboundMessageCols()
	insertParams := []any{
		inboundMessage.MessageId, inboundMessage.Customer.CustomerId,
		inboundMessage.PreviewText,
		inboundMessage.FirstSeqId, inboundMessage.LastSeqId,
		inboundMessage.CreatedAt, inboundMessage.UpdatedAt,
	}

	// Build the insert query.
	insertB.Addf("INSERT INTO inbound_message (%s)", insertCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return thread, ErrQuery
	}

	// Build the select query required after insert.
	q := builq.New()
	joinedCols := inboundMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins im") // inserted table, with alias im
	q("INNER JOIN customer c ON im.customer_id = c.customer_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return thread, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	// Make the insert query.
	err = tx.QueryRow(ctx, stmt, inboundMessage.MessageId, inboundMessage.Customer.CustomerId,
		inboundMessage.PreviewText,
		inboundMessage.FirstSeqId, inboundMessage.LastSeqId,
		inboundMessage.CreatedAt, inboundMessage.UpdatedAt).Scan(
		&inboundMessage.MessageId, &inboundMessage.Customer.CustomerId, &inboundMessage.Customer.Name,
		&inboundMessage.PreviewText, &inboundMessage.FirstSeqId, &inboundMessage.LastSeqId,
		&inboundMessage.CreatedAt, &inboundMessage.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return thread, ErrQuery
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return thread, ErrQuery
	}

	// hold db values for nullables.
	var (
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
	)

	// Persisted inbound message ID.
	inboundMessageId = sql.NullString{String: inboundMessage.MessageId, Valid: true}

	// Check if the thread is assigned to a member.
	// If assigned, then set the assigned member ID and assigned at for db insert values.
	// Otherwise, by default assigned member ID and assigned at will be NULL.
	if thread.AssignedMember != nil {
		assignedMemberId = sql.NullString{String: thread.AssignedMember.MemberId, Valid: true}
		assignedAt = sql.NullTime{Time: thread.AssignedMember.AssignedAt, Valid: true}
	}

	// Persist the thread with referenced inbound message ID.
	insertB = builq.Builder{}
	insertCols = threadCols()
	insertParams = []any{
		thread.ThreadId, thread.WorkspaceId, thread.Customer.CustomerId,
		assignedMemberId, assignedAt,
		thread.Title, thread.Description,
		thread.ThreadStatus.Status, thread.ThreadStatus.StatusChangedAt,
		thread.ThreadStatus.StatusChangedBy.MemberId,
		thread.ThreadStatus.Stage,
		thread.Replied, thread.Priority, thread.Channel,
		inboundMessageId,
		outboundMessageId,
		thread.CreatedBy.MemberId,
		thread.UpdatedBy.MemberId,
		thread.CreatedAt,
		thread.UpdatedAt,
	}

	insertB.Addf("INSERT INTO thread (%s)", insertCols)
	insertB.Addf(
		"VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$, %$)",
		insertParams...,
	)
	insertB.Addf("RETURNING %s", insertCols)

	insertQuery, _, err = insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return thread, ErrQuery
	}

	// Build the select query required after insert.
	q = builq.New()
	joinedCols = threadJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return thread, ErrQuery
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedBy.MemberId, &thread.CreatedBy.Name,
		&thread.UpdatedBy.MemberId, &thread.UpdatedBy.Name,
		&thread.CreatedAt, &thread.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return thread, ErrQuery
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return thread, ErrQuery
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

	// Sets the inbound message if a valid inbound message exists,
	// otherwise clears the inbound message.
	if inboundMessageId.Valid {
		customer := models.CustomerActor{
			CustomerId: inboundCustomerId.String,
			Name:       inboundCustomerName.String,
		}
		thread.InboundMessage = &models.InboundMessage{
			MessageId:   inboundMessageId.String,
			Customer:    customer,
			PreviewText: inboundPreviewText.String,
			FirstSeqId:  inboundFirstSeqId.String,
			LastSeqId:   inboundLastSeqId.String,
			CreatedAt:   inboundCreatedAt.Time,
			UpdatedAt:   inboundUpdatedAt.Time,
		}
	} else {
		thread.ClearInboundMessage()
	}

	// Sets the outbound message if a valid outbound message exists,
	// otherwise clears the outbound message.
	if outboundMessageId.Valid {
		member := models.MemberActor{
			MemberId: outboundMemberId.String,
			Name:     outboundMemberName.String,
		}
		thread.OutboundMessage = &models.OutboundMessage{
			MessageId:   outboundMessageId.String,
			Member:      member,
			PreviewText: outboundPreviewText.String,
			FirstSeqId:  outboundFirstSeqId.String,
			LastSeqId:   outboundLastSeqId.String,
			CreatedAt:   outboundCreatedAt.Time,
			UpdatedAt:   outboundUpdatedAt.Time,
		}
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func InsertThreadMessageTx(
	ctx context.Context, tx pgx.Tx, message *models.Message) (*models.Message, error) {
	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString

	if message.Customer != nil {
		customerId = sql.NullString{String: message.Customer.CustomerId, Valid: true}
	}
	if message.Member != nil {
		memberId = sql.NullString{String: message.Member.MemberId, Valid: true}
	}

	// Persist the message with referenced thread ID
	insertB := builq.Builder{}
	insertCols := threadMessageCols()
	insertParams := []any{
		message.MessageId, message.ThreadId, message.TextBody, message.Body,
		customerId, memberId, message.Channel, message.CreatedAt, message.UpdatedAt,
	}

	insertB.Addf("INSERT INTO message (%s)", insertCols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return message, ErrQuery
	}

	// Build the select query required after insert
	q := builq.New()
	joinedCols := threadMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins msg")
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return message, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&message.MessageId, &message.ThreadId, &message.TextBody, &message.Body,
		&customerId, &customerName,
		&memberId, &memberName,
		&message.Channel, &message.CreatedAt, &message.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return message, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return message, ErrQuery
	}

	if customerId.Valid {
		message.Customer = &models.CustomerActor{
			CustomerId: customerId.String,
			Name:       customerName.String,
		}
	}
	if memberId.Valid {
		message.Member = &models.MemberActor{
			MemberId: memberId.String,
			Name:     memberName.String,
		}
	}
	return message, nil
}

func (th *ThreadDB) InsertPostmarkInboundThreadMessage(
	ctx context.Context, inbound models.PostmarkInboundThreadMessage) (models.Thread, models.Message, error) {
	tx, err := th.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrTxQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	postmarkInbound := inbound.PostmarkInboundMessage

	// Insert inbound thread
	thread, err := InsertInboundThreadTx(ctx, tx, inbound.Thread)
	if err != nil {
		slog.Error("failed to insert inbound thread", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	// Insert thread message
	message, err := InsertThreadMessageTx(ctx, tx, inbound.Message)
	if err != nil {
		slog.Error("failed to insert thread message", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	insertCols := postmarkInboundMessageCols()
	insertParams := []any{
		message.MessageId, postmarkInbound.Payload, postmarkInbound.PostmarkMessageId,
		postmarkInbound.MailMessageId, postmarkInbound.ReplyMailMessageId,
		message.CreatedAt, message.UpdatedAt,
	}

	var insertB builq.Builder
	insertB.Addf("INSERT INTO postmark_inbound_message (%s)", insertCols)
	insertB.Addf("VALUES (%+$)", insertParams)
	insertB.Addf("RETURNING %s", insertCols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := insertB.DebugBuild()
		debugQuery(debug)
	}

	var throwablePk string
	err = tx.QueryRow(ctx, insertQuery, insertParams...).Scan(
		&throwablePk, &postmarkInbound.Payload, &postmarkInbound.PostmarkMessageId,
		&postmarkInbound.MailMessageId, &postmarkInbound.ReplyMailMessageId,
		&postmarkInbound.CreatedAt, &postmarkInbound.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrQuery
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Thread{}, models.Message{}, ErrTxQuery
	}
	return *thread, *message, nil
}

func (th *ThreadDB) ModifyThreadById(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {
	upsertQ := builq.New()
	upsertParams := make([]any, 0, len(fields)+1) // updates + thread ID
	threadCols := threadCols()

	upsertQ("UPDATE thread SET")
	var assignedMemberId sql.NullString
	for _, field := range fields {
		switch field {
		case "priority":
			upsertQ("priority = %$,", thread.Priority)
			upsertParams = append(upsertParams, thread.Priority)
		case "assignee":
			if thread.AssignedMember == nil {
				upsertQ("assignee_id = %$,", assignedMemberId)
				upsertParams = append(upsertParams, assignedMemberId)
			} else {
				upsertQ("assignee_id = %$,", thread.AssignedMember.MemberId)
				upsertParams = append(upsertParams, thread.AssignedMember.MemberId)
			}
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

	upsertQ("RETURNING %s", threadCols)

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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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

	var (
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
	)

	err = th.db.QueryRow(ctx, stmt, upsertParams...).Scan(
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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

	// Sets the inbound message if a valid inbound message exists,
	// otherwise clears the inbound message.
	if inboundMessageId.Valid {
		customer := models.CustomerActor{
			CustomerId: inboundCustomerId.String,
			Name:       inboundCustomerName.String,
		}
		thread.InboundMessage = &models.InboundMessage{
			MessageId:   inboundMessageId.String,
			Customer:    customer,
			PreviewText: inboundPreviewText.String,
			FirstSeqId:  inboundFirstSeqId.String,
			LastSeqId:   inboundLastSeqId.String,
			CreatedAt:   inboundCreatedAt.Time,
			UpdatedAt:   inboundUpdatedAt.Time,
		}
	} else {
		thread.ClearInboundMessage()
	}

	// Sets the outbound message if a valid outbound message exists,
	// otherwise clears the outbound message.
	if outboundMessageId.Valid {
		member := models.MemberActor{
			MemberId: outboundMemberId.String,
			Name:     outboundMemberName.String,
		}
		thread.OutboundMessage = &models.OutboundMessage{
			MessageId:   outboundMessageId.String,
			Member:      member,
			PreviewText: outboundPreviewText.String,
			FirstSeqId:  outboundFirstSeqId.String,
			LastSeqId:   outboundLastSeqId.String,
			CreatedAt:   outboundCreatedAt.Time,
			UpdatedAt:   outboundUpdatedAt.Time,
		}
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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

	// Sets the inbound message if a valid inbound message exists,
	// otherwise clears the inbound message.
	if inboundMessageId.Valid {
		customer := models.CustomerActor{
			CustomerId: inboundCustomerId.String,
			Name:       inboundCustomerName.String,
		}
		thread.InboundMessage = &models.InboundMessage{
			MessageId:   inboundMessageId.String,
			Customer:    customer,
			PreviewText: inboundPreviewText.String,
			FirstSeqId:  inboundFirstSeqId.String,
			LastSeqId:   inboundLastSeqId.String,
			CreatedAt:   inboundCreatedAt.Time,
			UpdatedAt:   inboundUpdatedAt.Time,
		}
	} else {
		thread.ClearInboundMessage()
	}

	// Sets the outbound message if a valid outbound message exists,
	// otherwise clears the outbound message.
	if outboundMessageId.Valid {
		member := models.MemberActor{
			MemberId: outboundMemberId.String,
			Name:     outboundMemberName.String,
		}
		thread.OutboundMessage = &models.OutboundMessage{
			MessageId:   outboundMessageId.String,
			Member:      member,
			PreviewText: outboundPreviewText.String,
			FirstSeqId:  outboundFirstSeqId.String,
			LastSeqId:   outboundLastSeqId.String,
			CreatedAt:   outboundCreatedAt.Time,
			UpdatedAt:   outboundUpdatedAt.Time,
		}
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func (th *ThreadDB) FetchThreadsByCustomerId(
	ctx context.Context, customerId string, channel *string) ([]models.Thread, error) {
	var thread models.Thread
	limit := 100
	threads := make([]models.Thread, 0, limit)

	params := []any{customerId}
	cols := threadJoinedCols()
	q := builq.New()
	q("SELECT %s FROM %s", cols, "thread th")
	q("INNER JOIN customer c ON th.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member am ON th.assignee_id = am.member_id")
	q("INNER JOIN member scm ON th.status_changed_by_id = scm.member_id")
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
	q("INNER JOIN member mc ON th.created_by_id = mc.member_id")
	q("INNER JOIN member mu ON th.updated_by_id = mu.member_id")

	q("WHERE th.customer_id = %$", customerId)
	if channel != nil {
		q("AND th.channel = %$", *channel)
		params = append(params, *channel)
	}

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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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
		// Sets the inbound message an if valid inbound message exists,
		// otherwise clears the inbound message.
		if inboundMessageId.Valid {
			customer := models.CustomerActor{
				CustomerId: inboundCustomerId.String,
				Name:       inboundCustomerName.String,
			}
			thread.InboundMessage = &models.InboundMessage{
				MessageId:   inboundMessageId.String,
				Customer:    customer,
				PreviewText: inboundPreviewText.String,
				FirstSeqId:  inboundFirstSeqId.String,
				LastSeqId:   inboundLastSeqId.String,
				CreatedAt:   inboundCreatedAt.Time,
				UpdatedAt:   inboundUpdatedAt.Time,
			}
		} else {
			thread.ClearInboundMessage()
		}
		// Sets the outbound message if a valid outbound message exists,
		// otherwise clears the outbound message.
		if outboundMessageId.Valid {
			member := models.MemberActor{
				MemberId: outboundMemberId.String,
				Name:     outboundMemberName.String,
			}
			thread.OutboundMessage = &models.OutboundMessage{
				MessageId:   outboundMessageId.String,
				Member:      member,
				PreviewText: outboundPreviewText.String,
				FirstSeqId:  outboundFirstSeqId.String,
				LastSeqId:   outboundLastSeqId.String,
				CreatedAt:   outboundCreatedAt.Time,
				UpdatedAt:   outboundUpdatedAt.Time,
			}
		} else {
			thread.ClearOutboundMessage()
		}
		threads = append(threads, thread)
		return nil
	})

	if err != nil {
		slog.Error("failed to scan", slog.Any("err", err))
		return []models.Thread{}, ErrQuery
	}
	return threads, nil
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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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
		// Sets the inbound message an if valid inbound message exists,
		// otherwise clears the inbound message.
		if inboundMessageId.Valid {
			customer := models.CustomerActor{
				CustomerId: inboundCustomerId.String,
				Name:       inboundCustomerName.String,
			}
			thread.InboundMessage = &models.InboundMessage{
				MessageId:   inboundMessageId.String,
				Customer:    customer,
				PreviewText: inboundPreviewText.String,
				FirstSeqId:  inboundFirstSeqId.String,
				LastSeqId:   inboundLastSeqId.String,
				CreatedAt:   inboundCreatedAt.Time,
				UpdatedAt:   inboundUpdatedAt.Time,
			}
		} else {
			thread.ClearInboundMessage()
		}
		// Sets the outbound message if a valid outbound message exists,
		// otherwise clears the outbound message.
		if outboundMessageId.Valid {
			member := models.MemberActor{
				MemberId: outboundMemberId.String,
				Name:     outboundMemberName.String,
			}
			thread.OutboundMessage = &models.OutboundMessage{
				MessageId:   outboundMessageId.String,
				Member:      member,
				PreviewText: outboundPreviewText.String,
				FirstSeqId:  outboundFirstSeqId.String,
				LastSeqId:   outboundLastSeqId.String,
				CreatedAt:   outboundCreatedAt.Time,
				UpdatedAt:   outboundUpdatedAt.Time,
			}
		} else {
			thread.ClearOutboundMessage()
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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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
		// Sets the inbound message an if valid inbound message exists,
		// otherwise clears the inbound message.
		if inboundMessageId.Valid {
			customer := models.CustomerActor{
				CustomerId: inboundCustomerId.String,
				Name:       inboundCustomerName.String,
			}
			thread.InboundMessage = &models.InboundMessage{
				MessageId:   inboundMessageId.String,
				Customer:    customer,
				PreviewText: inboundPreviewText.String,
				FirstSeqId:  inboundFirstSeqId.String,
				LastSeqId:   inboundLastSeqId.String,
				CreatedAt:   inboundCreatedAt.Time,
				UpdatedAt:   inboundUpdatedAt.Time,
			}
		} else {
			thread.ClearInboundMessage()
		}
		// Sets the outbound message if a valid outbound message exists,
		// otherwise clears the outbound message.
		if outboundMessageId.Valid {
			member := models.MemberActor{
				MemberId: outboundMemberId.String,
				Name:     outboundMemberName.String,
			}

			thread.OutboundMessage = &models.OutboundMessage{
				MessageId:   outboundMessageId.String,
				Member:      member,
				PreviewText: outboundPreviewText.String,
				FirstSeqId:  outboundFirstSeqId.String,
				LastSeqId:   outboundLastSeqId.String,
				CreatedAt:   outboundCreatedAt.Time,
				UpdatedAt:   outboundUpdatedAt.Time,
			}
		} else {
			thread.ClearOutboundMessage()
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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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
		// Sets the inbound message an if valid inbound message exists,
		// otherwise clears the inbound message.
		if inboundMessageId.Valid {
			customer := models.CustomerActor{
				CustomerId: inboundCustomerId.String,
				Name:       inboundCustomerName.String,
			}
			thread.InboundMessage = &models.InboundMessage{
				MessageId:   inboundMessageId.String,
				Customer:    customer,
				PreviewText: inboundPreviewText.String,
				FirstSeqId:  inboundFirstSeqId.String,
				LastSeqId:   inboundLastSeqId.String,
				CreatedAt:   inboundCreatedAt.Time,
				UpdatedAt:   inboundUpdatedAt.Time,
			}
		} else {
			thread.ClearInboundMessage()
		}
		// Sets the outbound message if a valid outbound message exists,
		// otherwise clears the outbound message.
		if outboundMessageId.Valid {
			member := models.MemberActor{
				MemberId: outboundMemberId.String,
				Name:     outboundMemberName.String,
			}
			thread.OutboundMessage = &models.OutboundMessage{
				MessageId:   outboundMessageId.String,
				Member:      member,
				PreviewText: outboundPreviewText.String,
				FirstSeqId:  outboundFirstSeqId.String,
				LastSeqId:   outboundLastSeqId.String,
				CreatedAt:   outboundCreatedAt.Time,
				UpdatedAt:   outboundUpdatedAt.Time,
			}
		} else {
			thread.ClearOutboundMessage()
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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
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
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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
		// Sets the inbound message an if valid inbound message exists,
		// otherwise clears the inbound message.
		if inboundMessageId.Valid {
			customer := models.CustomerActor{
				CustomerId: inboundCustomerId.String,
				Name:       inboundCustomerName.String,
			}
			thread.InboundMessage = &models.InboundMessage{
				MessageId:   inboundMessageId.String,
				Customer:    customer,
				PreviewText: inboundPreviewText.String,
				FirstSeqId:  inboundFirstSeqId.String,
				LastSeqId:   inboundLastSeqId.String,
				CreatedAt:   inboundCreatedAt.Time,
				UpdatedAt:   inboundUpdatedAt.Time,
			}
		} else {
			thread.ClearInboundMessage()
		}
		// Sets the outbound message if a valid outbound message exists,
		// otherwise clears the outbound message.
		if outboundMessageId.Valid {
			member := models.MemberActor{
				MemberId: outboundMemberId.String,
				Name:     outboundMemberName.String,
			}
			thread.OutboundMessage = &models.OutboundMessage{
				MessageId:   outboundMessageId.String,
				Member:      member,
				PreviewText: outboundPreviewText.String,
				FirstSeqId:  outboundFirstSeqId.String,
				LastSeqId:   outboundLastSeqId.String,
				CreatedAt:   outboundCreatedAt.Time,
				UpdatedAt:   outboundUpdatedAt.Time,
			}
		} else {
			thread.ClearOutboundMessage()
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

func (th *ThreadDB) AppendInboundThreadMessage(
	ctx context.Context, inbound models.ThreadMessage,
) (models.Message, error) {
	tx, err := th.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", slog.Any("err", err))
		return models.Message{}, ErrTxQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	// upsert thread linked inbound message
	// based on upsert(created) flag insert to thread
	// insert message

	thread := inbound.Thread
	message := inbound.Message

	if thread.InboundMessage == nil {
		slog.Error("thread inbound message cannot be empty", slog.Any("thread", thread))
		return models.Message{}, ErrQuery
	}

	inboundMessage := thread.InboundMessage

	var insertB builq.Builder
	cols := inboundMessageCols()
	insertParams := []any{
		inboundMessage.MessageId, inboundMessage.Customer.CustomerId,
		inboundMessage.PreviewText, inboundMessage.FirstSeqId, inboundMessage.LastSeqId,
		inboundMessage.CreatedAt, inboundMessage.UpdatedAt,
	}

	// Build the upsert query to insert thread inbound message
	insertB.Addf("INSERT INTO inbound_message (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (message_id)")
	insertB.Addf("DO UPDATE")
	insertB.Addf("SET")
	insertB.Addf("preview_text = EXCLUDED.preview_text,")
	insertB.Addf("last_seq_id = EXCLUDED.last_seq_id,")
	insertB.Addf("updated_at = EXCLUDED.updated_at")
	insertB.Addf("RETURNING %s, (xmax = 0) AS is_created", cols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build upsert query", slog.Any("error", err))
		return models.Message{}, ErrQuery
	}

	// Build the select query
	q := builq.New()
	joinedCols := inboundMessageJoinedCols()
	q("WITH ups AS (%s)", insertQuery)
	q("SELECT %s, im.is_created", joinedCols)
	q("FROM ups im")
	q("INNER JOIN customer c ON im.customer_id = c.customer_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var isCreated bool
	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&inboundMessage.MessageId,
		&inboundMessage.Customer.CustomerId, &inboundMessage.Customer.Name,
		&inboundMessage.PreviewText,
		&inboundMessage.FirstSeqId, &inboundMessage.LastSeqId,
		&inboundMessage.CreatedAt, &inboundMessage.UpdatedAt,
		&isCreated,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}
	// Check if the inbound message is create or updated
	if isCreated {
		// Link inbound message to the thread
		q = builq.New()
		updates := []any{inboundMessage.MessageId, thread.ThreadId}
		q("UPDATE thread SET")
		q("inbound_message_id = %$, updated_at = NOW() WHERE thread_id = %$", updates...)

		stmt, _, err = q.Build()
		if err != nil {
			slog.Error("failed to build query", slog.Any("error", err))
			return models.Message{}, ErrQuery
		}
		_, err = tx.Exec(ctx, stmt, updates...)
		if err != nil {
			slog.Error("failed to update query", slog.Any("err", err))
			return models.Message{}, ErrQuery
		}
	}

	// hold db nullables.
	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString
	if message.Customer != nil {
		customerId = sql.NullString{String: message.Customer.CustomerId, Valid: true}
	}
	if message.Member != nil {
		memberId = sql.NullString{String: message.Member.MemberId, Valid: true}
	}

	// Persist the message with referenced thread ID
	insertB = builq.Builder{}
	cols = threadMessageCols()
	insertParams = []any{
		message.MessageId, message.ThreadId, message.TextBody, message.Body,
		customerId, memberId, message.Channel, message.CreatedAt, message.UpdatedAt,
	}

	insertB.Addf("INSERT INTO message (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", cols)

	insertQuery, _, err = insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	// Build the select query required after insert
	q = builq.New()
	joinedCols = threadMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins msg")
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&message.MessageId, &message.ThreadId, &message.TextBody, &message.Body,
		&customerId, &customerName,
		&memberId, &memberName,
		&message.Channel, &message.CreatedAt, &message.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if customerId.Valid {
		message.Customer = &models.CustomerActor{
			CustomerId: customerId.String,
			Name:       customerName.String,
		}
	}
	if memberId.Valid {
		message.Member = &models.MemberActor{
			MemberId: memberId.String,
			Name:     memberName.String,
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Message{}, ErrTxQuery
	}
	return *message, nil
}

// AppendOutboundThreadMessage inserts a member chat into the database.
func (th *ThreadDB) AppendOutboundThreadMessage(
	ctx context.Context, outbound models.ThreadMessage,
) (models.Message, error) {
	tx, err := th.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", slog.Any("err", err))
		return models.Message{}, ErrTxQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	// upsert thread linked outbound message
	// based on upsert(created) flag insert to thread
	// insert message

	thread := outbound.Thread
	message := outbound.Message

	if thread.OutboundMessage == nil {
		slog.Error("thread outbound message cannot be empty", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}

	outboundMessage := thread.OutboundMessage

	var insertB builq.Builder
	cols := outboundMessageCols()
	insertParams := []any{
		outboundMessage.MessageId, outboundMessage.Member.MemberId,
		outboundMessage.PreviewText, outboundMessage.FirstSeqId, outboundMessage.LastSeqId,
		outboundMessage.CreatedAt, outboundMessage.UpdatedAt,
	}

	// Build the upsert query to insert thread outbound message
	insertB.Addf("INSERT INTO outbound_message (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (message_id)")
	insertB.Addf("DO UPDATE")
	insertB.Addf("SET")
	insertB.Addf("preview_text = EXCLUDED.preview_text,")
	insertB.Addf("last_seq_id = EXCLUDED.last_seq_id,")
	insertB.Addf("updated_at = EXCLUDED.updated_at")
	insertB.Addf("RETURNING %s, (xmax = 0) AS is_created", cols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	// Build the select query
	q := builq.New()
	joinedCols := outboundMessageJoinedCols()
	q("WITH ups AS (%s)", insertQuery)
	q("SELECT %s, om.is_created", joinedCols)
	q("FROM ups om")
	q("INNER JOIN member m ON om.member_id = m.member_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var isCreated bool
	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&outboundMessage.MessageId,
		&outboundMessage.Member.MemberId, &outboundMessage.Member.Name,
		&outboundMessage.PreviewText,
		&outboundMessage.FirstSeqId, &outboundMessage.LastSeqId,
		&outboundMessage.CreatedAt, &outboundMessage.UpdatedAt,
		&isCreated,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}
	// Check if the inbound message is create or updated
	if isCreated {
		// Link inbound message to the thread
		q = builq.New()
		updates := []any{outboundMessage.MessageId, thread.ThreadId}
		q("UPDATE thread SET")
		q("outbound_message_id = %$, updated_at = NOW() WHERE thread_id = %$", updates...)

		stmt, _, err = q.Build()
		if err != nil {
			slog.Error("failed to build query", slog.Any("error", err))
			return models.Message{}, ErrQuery
		}
		_, err = tx.Exec(ctx, stmt, updates...)
		if err != nil {
			slog.Error("failed to update query", slog.Any("err", err))
			return models.Message{}, ErrQuery
		}
	}

	// hold db nullables.
	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString
	if message.Customer != nil {
		customerId = sql.NullString{String: message.Customer.CustomerId, Valid: true}
	}
	if message.Member != nil {
		memberId = sql.NullString{String: message.Member.MemberId, Valid: true}
	}

	// Persist the message with referenced thread ID
	insertB = builq.Builder{}
	cols = threadMessageCols()
	insertParams = []any{
		message.MessageId, message.ThreadId, message.TextBody, message.Body,
		customerId, memberId, message.Channel, message.CreatedAt, message.UpdatedAt,
	}

	insertB.Addf("INSERT INTO message (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", cols)

	insertQuery, _, err = insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	// Build the select query required after insert
	q = builq.New()
	joinedCols = threadMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins msg")
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&message.MessageId, &message.ThreadId, &message.TextBody, &message.Body,
		&customerId, &customerName,
		&memberId, &memberName,
		&message.Channel, &message.CreatedAt, &message.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if customerId.Valid {
		message.Customer = &models.CustomerActor{
			CustomerId: customerId.String,
			Name:       customerName.String,
		}
	}
	if memberId.Valid {
		message.Member = &models.MemberActor{
			MemberId: memberId.String,
			Name:     memberName.String,
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Message{}, ErrTxQuery
	}
	return *message, nil
}

func (th *ThreadDB) FetchMessagesByThreadId(
	ctx context.Context, threadId string) ([]models.Message, error) {
	var message models.Message
	messages := make([]models.Message, 0, 100)

	q := builq.New()
	messagesJoinedCols := threadMessageJoinedCols()
	q("SELECT %s FROM message msg", messagesJoinedCols)
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")
	q("WHERE msg.thread_id = %$", threadId)

	q("ORDER BY msg.created_at ASC")
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

	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString

	rows, _ := th.db.Query(ctx, stmt, threadId)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&message.MessageId, &message.ThreadId, &message.TextBody,
		&message.Body, &customerId, &customerName,
		&memberId, &memberName,
		&message.Channel,
		&message.CreatedAt, &message.UpdatedAt,
	}, func() error {
		if customerId.Valid {
			message.Customer = &models.CustomerActor{
				CustomerId: customerId.String,
				Name:       customerName.String,
			}
		}
		if memberId.Valid {
			message.Member = &models.MemberActor{
				MemberId: memberId.String,
				Name:     memberName.String,
			}
		}
		messages = append(messages, message)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Message{}, ErrQuery
	}
	return messages, nil
}

func (th *ThreadDB) FetchMessagesWithAttachmentsByThreadId(
	ctx context.Context, threadId string) ([]models.MessageWithAttachments, error) {
	var message models.MessageWithAttachments

	limit := 100
	messages := make([]models.MessageWithAttachments, 0, limit)
	cols := threadMessageJoinedCols()

	stmt := `SELECT
		%s,
		COALESCE(
			(
				SELECT JSON_AGG(
					JSON_BUILD_OBJECT(
						'attachmentId', ma.attachment_id,
						'messageId', ma.message_id,
						'name', ma.name,
						'contentType', ma.content_type,
						'contentKey', ma.content_key,
						'contentUrl', ma.content_url,
						'spam', ma.spam,
						'hasError', ma.has_error,
						'error', ma.error,
						'md5Hash', ma.md5_hash,
						'createdAt', ma.created_at AT TIME ZONE 'UTC',
						'updatedAt', ma.updated_at AT TIME ZONE 'UTC'
					)
				)
				FROM (
					SELECT *
					FROM message_attachment
					WHERE message_id = msg.message_id
					ORDER BY created_at ASC
					LIMIT 10
				) ma
			), 
			'[]'::json
		) as attachments
	FROM message msg
	LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON msg.member_id = m.member_id`

	stmt = fmt.Sprintf(stmt, cols)

	q := builq.New()
	q("%s", stmt)
	q("WHERE msg.thread_id = %$", threadId)
	q("ORDER BY msg.created_at ASC")
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

	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString
	var attachmentsJson []byte

	rows, _ := th.db.Query(ctx, stmt, threadId)

	defer rows.Close()

	_, err = pgx.ForEachRow(rows, []any{
		&message.MessageId, &message.ThreadId, &message.TextBody,
		&message.Body, &customerId, &customerName,
		&memberId, &memberName,
		&message.Channel,
		&message.CreatedAt, &message.UpdatedAt,
		&attachmentsJson,
	}, func() error {
		if customerId.Valid {
			message.Customer = &models.CustomerActor{
				CustomerId: customerId.String,
				Name:       customerName.String,
			}
		}
		if memberId.Valid {
			message.Member = &models.MemberActor{
				MemberId: memberId.String,
			}
		}
		var attachments []models.MessageAttachment
		if len(attachmentsJson) > 0 {
			if err := json.Unmarshal(attachmentsJson, &attachments); err != nil {
				slog.Error("failed to unmarshal attachments",
					slog.String("messageId", message.MessageId))
				slog.Any("err", err)
				return err
			}
			message.Attachments = attachments
		} else {
			message.Attachments = []models.MessageAttachment{}
		}
		messageCopy := message
		messages = append(messages, messageCopy)
		return nil
	})
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.MessageWithAttachments{}, ErrQuery
	}
	return messages, nil
}

// ComputeStatusMetricsByWorkspaceId computes the thread count metrics for the workspace.
// Returns the count of active threads, needs first response threads, waiting on customer threads,
// hold threads, and needs next response threads.
// Ignores visitor customer threads.
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

// ComputeAssigneeMetricsByMember computes the thread count metrics for the member.
// Returns the count of member assigned threads, unassigned threads, and other assigned threads.
// Ignores visitor customer threads.
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

func (th *ThreadDB) FetchThreadByPostmarkInboundInReplyMessageId(
	ctx context.Context, workspaceId string, inReplyMessageId string) (models.Thread, error) {
	var thread models.Thread

	var selectB builq.Builder
	selectB.Addf("SELECT m.thread_id AS thread_id")
	selectB.Addf("FROM postmark_inbound_message p")
	selectB.Addf("INNER JOIN message m ON p.message_id = m.message_id")
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
	q("LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id")
	q("LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id")
	q("LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id")
	q("LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id")
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
		assignedMemberId    sql.NullString
		assignedMemberName  sql.NullString
		assignedAt          sql.NullTime
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
		outboundMessageId   sql.NullString
		outboundMemberId    sql.NullString
		outboundMemberName  sql.NullString
		outboundPreviewText sql.NullString
		outboundFirstSeqId  sql.NullString
		outboundLastSeqId   sql.NullString
		outboundCreatedAt   sql.NullTime
		outboundUpdatedAt   sql.NullTime
	)

	err = th.db.QueryRow(ctx, stmt, workspaceId, inReplyMessageId).Scan(
		&thread.ThreadId, &thread.WorkspaceId, &thread.Customer.CustomerId, &thread.Customer.Name,
		&assignedMemberId, &assignedMemberName, &assignedAt,
		&thread.Title, &thread.Description,
		&thread.ThreadStatus.Status,
		&thread.ThreadStatus.StatusChangedAt,
		&thread.ThreadStatus.StatusChangedBy.MemberId, &thread.ThreadStatus.StatusChangedBy.Name,
		&thread.ThreadStatus.Stage,
		&thread.Replied, &thread.Priority, &thread.Channel,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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

	// Sets the inbound message if a valid inbound message exists,
	// otherwise clears the inbound message.
	if inboundMessageId.Valid {
		customer := models.CustomerActor{
			CustomerId: inboundCustomerId.String,
			Name:       inboundCustomerName.String,
		}
		thread.InboundMessage = &models.InboundMessage{
			MessageId:   inboundMessageId.String,
			Customer:    customer,
			PreviewText: inboundPreviewText.String,
			FirstSeqId:  inboundFirstSeqId.String,
			LastSeqId:   inboundLastSeqId.String,
			CreatedAt:   inboundCreatedAt.Time,
			UpdatedAt:   inboundUpdatedAt.Time,
		}
	} else {
		thread.ClearInboundMessage()
	}

	// Sets the outbound message if a valid outbound message exists,
	// otherwise clears the outbound message.
	if outboundMessageId.Valid {
		member := models.MemberActor{
			MemberId: outboundMemberId.String,
			Name:     outboundMemberName.String,
		}
		thread.OutboundMessage = &models.OutboundMessage{
			MessageId:   outboundMessageId.String,
			Member:      member,
			PreviewText: outboundPreviewText.String,
			FirstSeqId:  outboundFirstSeqId.String,
			LastSeqId:   outboundLastSeqId.String,
			CreatedAt:   outboundCreatedAt.Time,
			UpdatedAt:   outboundUpdatedAt.Time,
		}
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func (th *ThreadDB) AppendPostmarkInboundThreadMessage(
	ctx context.Context, inbound models.PostmarkInboundThreadMessage) (models.Message, error) {

	// start transaction
	// If fails then stop the execution and return the error.
	tx, err := th.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", slog.Any("err", err))
		}
	}(tx, ctx)

	// upsert thread linked inbound log
	// based on upsert(created) flag insert to thread
	// insert message
	// insert postmark inbound message

	thread := inbound.Thread
	message := inbound.Message
	pmInboundMessage := inbound.PostmarkInboundMessage

	if thread.InboundMessage == nil {
		slog.Error("thread inbound message cannot be empty", slog.Any("thread", thread))
		return models.Message{}, ErrQuery
	}

	inboundMessage := thread.InboundMessage

	var insertB builq.Builder
	cols := inboundMessageCols()
	insertParams := []any{
		inboundMessage.MessageId, inboundMessage.Customer.CustomerId,
		inboundMessage.PreviewText, inboundMessage.FirstSeqId, inboundMessage.LastSeqId,
		inboundMessage.CreatedAt, inboundMessage.UpdatedAt,
	}

	// Build the upsert query to insert thread inbound message
	insertB.Addf("INSERT INTO inbound_message (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (message_id)")
	insertB.Addf("DO UPDATE")
	insertB.Addf("SET")
	insertB.Addf("preview_text = EXCLUDED.preview_text,")
	insertB.Addf("last_seq_id = EXCLUDED.last_seq_id,")
	insertB.Addf("updated_at = EXCLUDED.updated_at")
	insertB.Addf("RETURNING %s, (xmax = 0) AS is_created", cols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build upsert query", slog.Any("error", err))
		return models.Message{}, ErrQuery
	}

	// Build the select query
	q := builq.New()
	joinedCols := inboundMessageJoinedCols()
	q("WITH ups AS (%s)", insertQuery)
	q("SELECT %s, im.is_created", joinedCols)
	q("FROM ups im")
	q("INNER JOIN customer c ON im.customer_id = c.customer_id")

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("error", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var isCreated bool
	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&inboundMessage.MessageId,
		&inboundMessage.Customer.CustomerId, &inboundMessage.Customer.Name,
		&inboundMessage.PreviewText,
		&inboundMessage.FirstSeqId, &inboundMessage.LastSeqId,
		&inboundMessage.CreatedAt, &inboundMessage.UpdatedAt,
		&isCreated,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}
	// Check if the inbound message is create or updated
	if isCreated {
		// Link inbound message to the thread
		q = builq.New()
		updates := []any{inboundMessage.MessageId, thread.ThreadId}
		q("UPDATE thread SET")
		q("inbound_message_id = %$, updated_at = NOW() WHERE thread_id = %$", updates...)

		stmt, _, err = q.Build()
		if err != nil {
			slog.Error("failed to build query", slog.Any("error", err))
			return models.Message{}, ErrQuery
		}
		_, err = tx.Exec(ctx, stmt, updates...)
		if err != nil {
			slog.Error("failed to update query", slog.Any("err", err))
			return models.Message{}, ErrQuery
		}
	}

	// hold db nullables.
	var customerId, customerName sql.NullString
	var memberId, memberName sql.NullString
	if message.Customer != nil {
		customerId = sql.NullString{String: message.Customer.CustomerId, Valid: true}
	}
	if message.Member != nil {
		memberId = sql.NullString{String: message.Member.MemberId, Valid: true}
	}

	// Persist the message with referenced thread ID
	insertB = builq.Builder{}
	cols = threadMessageCols()
	insertParams = []any{
		message.MessageId, message.ThreadId, message.TextBody, message.Body,
		customerId, memberId, message.Channel, message.CreatedAt, message.UpdatedAt,
	}

	insertB.Addf("INSERT INTO message (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("RETURNING %s", cols)

	insertQuery, _, err = insertB.Build()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	// Build the select query required after insert
	q = builq.New()
	joinedCols = threadMessageJoinedCols()

	q("WITH ins AS (%s)", insertQuery)
	q("SELECT %s FROM %s", joinedCols, "ins msg")
	q("LEFT OUTER JOIN customer c ON msg.customer_id = c.customer_id")
	q("LEFT OUTER JOIN member m ON msg.member_id = m.member_id")

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&message.MessageId, &message.ThreadId, &message.TextBody, &message.Body,
		&customerId, &customerName,
		&memberId, &memberName,
		&message.Channel, &message.CreatedAt, &message.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if customerId.Valid {
		message.Customer = &models.CustomerActor{
			CustomerId: customerId.String,
			Name:       customerName.String,
		}
	}
	if memberId.Valid {
		message.Member = &models.MemberActor{
			MemberId: memberId.String,
			Name:     memberName.String,
		}
	}

	// Insert the Postmark inbound message
	q = builq.New()
	pmInboundCols := postmarkInboundMessageCols()

	insertParams = []any{
		message.MessageId, pmInboundMessage.Payload, pmInboundMessage.PostmarkMessageId,
		pmInboundMessage.MailMessageId, pmInboundMessage.ReplyMailMessageId,
		// consider the time space of the message rather than the postmark message.
		message.CreatedAt, message.UpdatedAt,
	}

	q("INSERT INTO postmark_inbound_message (%s)", pmInboundCols)
	q("VALUES (%$, %$, %$, %$, %$, %$, %$)", insertParams...)
	q("RETURNING %s", pmInboundCols)

	stmt, _, err = q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	var throwablePk string
	err = tx.QueryRow(ctx, stmt, insertParams...).Scan(
		&throwablePk, &pmInboundMessage.Payload, &pmInboundMessage.PostmarkMessageId,
		&pmInboundMessage.MailMessageId, &pmInboundMessage.ReplyMailMessageId,
		&pmInboundMessage.CreatedAt, &pmInboundMessage.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Message{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Message{}, ErrQuery
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Message{}, ErrTxQuery
	}
	return *message, nil
}

func (th *ThreadDB) CheckPostmarkInboundMessageExists(ctx context.Context, messageId string) (bool, error) {
	var isExist bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM postmark_inbound_message WHERE postmark_message_id = $1	
	)`

	err := th.db.QueryRow(ctx, stmt, messageId).Scan(&isExist)
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return false, ErrQuery
	}
	return isExist, nil
}

func (th *ThreadDB) InsertMessageAttachment(
	ctx context.Context, attachment models.MessageAttachment) (models.MessageAttachment, error) {
	cols := messageAttachmentCols()
	q := builq.New()
	insertParams := []any{
		attachment.AttachmentId, attachment.MessageId, attachment.Name,
		attachment.ContentType, attachment.ContentKey, attachment.ContentUrl,
		attachment.Spam, attachment.HasError, attachment.Error, attachment.MD5Hash,
		attachment.CreatedAt, attachment.UpdatedAt,
	}
	q("INSERT INTO message_attachment (%s)", cols)
	q("VALUES (%+$)", insertParams)
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.MessageAttachment{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = th.db.QueryRow(ctx, stmt, insertParams...).Scan(
		&attachment.AttachmentId, &attachment.MessageId, &attachment.Name,
		&attachment.ContentType, &attachment.ContentKey, &attachment.ContentUrl,
		&attachment.Spam, &attachment.HasError, &attachment.Error, &attachment.MD5Hash,
		&attachment.CreatedAt, &attachment.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.MessageAttachment{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.MessageAttachment{}, ErrQuery
	}
	return attachment, nil
}

func (th *ThreadDB) FetchMessageAttachmentById(
	ctx context.Context, messageId, attachmentId string) (models.MessageAttachment, error) {
	var attachment models.MessageAttachment
	cols := messageAttachmentCols()

	q := builq.New()
	q("SELECT %s FROM message_attachment", cols)
	q("WHERE message_id = %$ AND attachment_id = %$", messageId, attachmentId)

	stmt, _, err := q.Build()

	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.MessageAttachment{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = th.db.QueryRow(ctx, stmt, messageId, attachmentId).Scan(
		&attachment.AttachmentId, &attachment.MessageId, &attachment.Name,
		&attachment.ContentType, &attachment.ContentKey, &attachment.ContentUrl,
		&attachment.Spam, &attachment.HasError, &attachment.Error, &attachment.MD5Hash,
		&attachment.CreatedAt, &attachment.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.MessageAttachment{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.MessageAttachment{}, ErrQuery
	}
	return attachment, nil
}
