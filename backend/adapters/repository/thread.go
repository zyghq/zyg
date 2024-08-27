package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (tc *ThreadChatDB) InsertInboundThreadChat(
	ctx context.Context, inboundMessage models.InboundMessage,
	thread models.Thread, chat models.Chat) (models.Thread, models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("error", err))
		return models.Thread{}, chat, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	messageId := inboundMessage.GenId()
	stmt := `
		with ins as (
			insert into inbound_message (message_id, customer_id, first_seq_id, last_seq_id, preview_text)
				values ($1, $2, $3, $4, $5)
			returning
				message_id, customer_id,
				first_seq_id, last_seq_id,
				preview_text,
				created_at, updated_at
		) select
			i.message_id,
			c.customer_id,
			c.name,
			i.first_seq_id,
			i.last_seq_id,
			i.preview_text,
			i.created_at,
			i.updated_at
		from ins i
		inner join customer c on i.customer_id = c.customer_id
	`

	err = tx.QueryRow(ctx, stmt, messageId, inboundMessage.CustomerId,
		inboundMessage.FirstSeqId, inboundMessage.LastSeqId, inboundMessage.PreviewText).Scan(
		&inboundMessage.MessageId,
		&inboundMessage.CustomerId, &inboundMessage.CustomerName,
		&inboundMessage.FirstSeqId, &inboundMessage.LastSeqId,
		&inboundMessage.PreviewText,
		&inboundMessage.CreatedAt, &inboundMessage.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrQuery
	}

	var (
		inboundMessageId    sql.NullString
		inboundCustomerId   sql.NullString
		inboundCustomerName sql.NullString
		inboundPreviewText  sql.NullString
		inboundFirstSeqId   sql.NullString
		inboundLastSeqId    sql.NullString
		inboundCreatedAt    sql.NullTime
		inboundUpdatedAt    sql.NullTime
	)

	threadId := thread.GenId()
	stmt = `
		with ins as (
			insert into thread (
				thread_id, workspace_id, customer_id, assignee_id,
				title, description, status, read, replied,
				priority, spam, channel, preview_text,
				inbound_message_id 
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			returning
				thread_id, workspace_id, customer_id, assignee_id,
				title, description, sequence, status, read, replied,
				priority, spam, channel, preview_text,
				inbound_message_id, outbound_message_id,
				created_at, updated_at
		) select
			ins.thread_id as thread_id,
			ins.workspace_id as workspace_id,
			c.customer_id as customer_id,
			c.name as customer_name,
			m.member_id as assignee_id,
			m.name as assignee_name,
			ins.title as title,
			ins.description as description,
			ins.sequence as sequence,
			ins.status as status,
			ins.read as read,
			ins.replied as replied,
			ins.priority as priority,
			ins.spam as spam,
			ins.channel as channel,
			ins.preview_text as preview_text,
			inb.message_id,
			inbc.customer_id, 
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
            inb.updated_at,
			ins.created_at as created_at,
			ins.updated_at as updated_at
		from ins
		inner join customer c on ins.customer_id = c.customer_id
		left outer join member m on ins.assignee_id = m.member_id
		left outer join inbound_message inb on ins.inbound_message_id = inb.message_id
		left outer join customer inbc on inb.customer_id = inbc.customer_id
	`

	err = tx.QueryRow(ctx, stmt, threadId, thread.WorkspaceId, thread.CustomerId, thread.AssigneeId,
		thread.Title, thread.Description, thread.Status, thread.Read, thread.Replied,
		thread.Priority, thread.Spam, thread.Channel, thread.PreviewText,
		inboundMessage.MessageId).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, chat, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, chat, ErrQuery
	}

	// set the inbound message if an inbound message exists.
	if inboundMessageId.Valid {
		thread.AddInboundMessage(inboundMessageId.String,
			inboundCustomerId.String, inboundCustomerName.String,
			inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
			inboundCreatedAt.Time,
			inboundUpdatedAt.Time,
		)
	} else {
		thread.ClearInboundMessage()
	}

	// insert thread chat
	stmt = `
		with ins as (
			insert into chat (
				chat_id, thread_id, body, sequence,
				customer_id, member_id, is_head
			)
			values ($1, $2, $3, $4, $5, $6, $7)
			returning
				chat_id, thread_id, body, sequence,
				customer_id, member_id, is_head, created_at,
				updated_at
		) select ins.chat_id as chat_id,
			ins.thread_id as thread_id,
			ins.body as body,
			ins.sequence as sequence,
			c.customer_id as customer_id,
			c.name as customer_name,
			m.member_id as member_id,
			m.name as member_name,
			ins.is_head as is_head,
			ins.created_at as created_at,
			ins.updated_at as updated_at
		from ins
		left outer join customer c on ins.customer_id = c.customer_id
		left outer join member m on ins.member_id = m.member_id
	`

	chatId := chat.GenId()
	err = tx.QueryRow(ctx, stmt, chatId, thread.ThreadId, chat.Body, thread.Sequence,
		chat.CustomerId, chat.MemberId, chat.IsHead).Scan(
		&chat.ChatId, &chat.ThreadId, &chat.Body, &chat.Sequence,
		&chat.CustomerId,
		&chat.CustomerName,
		&chat.MemberId,
		&chat.MemberName,
		&chat.IsHead, &chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrQuery
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrTxQuery
	}
	return thread, chat, nil
}

func (tc *ThreadChatDB) ModifyThreadById(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {

	args := make([]interface{}, 0, len(fields))
	ups := `UPDATE thread SET`

	for i, field := range fields {
		switch field {
		case "priority":
			args = append(args, thread.Priority)
			ups += fmt.Sprintf(" %s = $%d,", "priority", i+1)
		case "assignee":
			args = append(args, thread.AssigneeId)
			ups += fmt.Sprintf(" %s = $%d,", "assignee_id", i+1)
		case "status":
			args = append(args, thread.Status)
			ups += fmt.Sprintf(" %s = $%d,", "status", i+1)
		case "read":
			args = append(args, thread.Read)
			ups += fmt.Sprintf(" %s = $%d,", "read", i+1)
		case "replied":
			args = append(args, thread.Replied)
			ups += fmt.Sprintf(" %s = $%d,", "replied", i+1)
		case "spam":
			args = append(args, thread.Spam)
			ups += fmt.Sprintf(" %s = $%d,", "spam", i+1)
		}
	}

	ups += " updated_at = NOW()"
	ups += fmt.Sprintf(" WHERE thread_id = $%d", len(fields)+1)
	args = append(args, thread.ThreadId)

	var (
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

	stmt := `WITH ups AS (
		%s
		RETURNING
			thread_id, workspace_id, customer_id, assignee_id,
			title, description, sequence, status, read, replied,
			priority, spam, channel, preview_text,
			inbound_message_id, outbound_message_id,
			created_at, updated_at
		) SELECT
		 	ups.thread_id AS thread_id,
			ups.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ups.title AS title,
			ups.description AS description,
			ups.sequence AS sequence,
			ups.status AS status,
			ups.read AS read,
			ups.replied AS replied,
			ups.priority AS priority,
			ups.spam AS spam,
			ups.channel AS channel,
			ups.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id,
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
            inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON ups.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON ups.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
	`

	stmt = fmt.Sprintf(stmt, ups)

	err := tc.db.QueryRow(ctx, stmt, args...).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	// set the inbound message if an inbound message exists.
	if inboundMessageId.Valid {
		thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
			inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
			inboundCreatedAt.Time,
			inboundUpdatedAt.Time,
		)
	} else {
		thread.ClearInboundMessage()
	}

	// set the outbound message if an outbound message exists.
	if outboundMessageId.Valid {
		thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
			outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
			outboundCreatedAt.Time,
			outboundUpdatedAt.Time,
		)
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func (tc *ThreadChatDB) LookupByWorkspaceThreadId(
	ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error) {
	var thread models.Thread

	args := make([]interface{}, 0, 3)
	stmt := `SELECT th.thread_id AS thread_id,
		th.workspace_id AS workspace_id,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS assignee_id,
		m.name AS assignee_name,
		th.title AS title,
		th.description AS description,
		th.sequence AS sequence,
		th.status AS status,
		th.read AS read,
		th.replied AS replied,
		th.priority AS priority,
		th.spam AS spam,
		th.channel AS channel,
		th.preview_text AS preview_text,
		inb.message_id,
		inbc.customer_id, 
		inbc.name,
		inb.preview_text,
		inb.first_seq_id,
		inb.last_seq_id,
		inb.created_at,
		inb.updated_at,
		oub.message_id,
		oubm.member_id,
		oubm.name,
		oub.preview_text,
		oub.first_seq_id,
		oub.last_seq_id,
		oub.created_at,
		oub.updated_at,
		th.created_at AS created_at,
		th.updated_at AS updated_at
	FROM thread th
	INNER JOIN customer c ON th.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
	LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id
	LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id
	LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
	LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
	WHERE th.workspace_id = $1 AND th.thread_id = $2`

	args = append(args, workspaceId)
	args = append(args, threadId)

	if channel != nil {
		stmt += " AND th.channel = $3"
		args = append(args, *channel)
	}

	var (
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

	err := tc.db.QueryRow(ctx, stmt, args...).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
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

	// set the inbound message if an inbound message exists.
	if inboundMessageId.Valid {
		thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
			inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
			inboundCreatedAt.Time,
			inboundUpdatedAt.Time,
		)
	} else {
		thread.ClearInboundMessage()
	}

	// set the outbound message if an outbound message exists.
	if outboundMessageId.Valid {
		thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
			outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
			outboundCreatedAt.Time,
			outboundUpdatedAt.Time,
		)
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func (tc *ThreadChatDB) FetchThreadsByCustomerId(
	ctx context.Context, customerId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	args := make([]interface{}, 0, 3)
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			th.title AS title,
			th.description AS description,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id,
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
            inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
		WHERE th.customer_id = $1
	`

	args = append(args, customerId)

	if channel != nil {
		stmt += " AND th.channel = $2"
		args = append(args, *channel)
	}

	if role != nil {
		stmt += " AND c.role = $3"
		args = append(args, *role)
	}

	stmt += " ORDER BY inb.last_seq_id DESC LIMIT 100"

	rows, _ := tc.db.Query(ctx, stmt, args...)

	defer rows.Close()

	var (
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

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// set the inbound message if an inbound message exists.
		if inboundMessageId.Valid {
			thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
				inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
				inboundCreatedAt.Time,
				inboundUpdatedAt.Time,
			)
		} else {
			thread.ClearInboundMessage()
		}
		// set the outbound message if an outbound message exists.
		if outboundMessageId.Valid {
			thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
				outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
				outboundCreatedAt.Time,
				outboundUpdatedAt.Time,
			)
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

func (tc *ThreadChatDB) UpdateAssignee(
	ctx context.Context, threadId string, assigneeId string) (models.Thread, error) {
	var thread models.Thread
	var (
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

	stmt := `
		WITH ups AS (
			UPDATE thread
			SET assignee_id = $1, updated_at = NOW()
			WHERE thread_id = $2
			RETURNING
				thread_id, workspace_id, customer_id, assignee_id,
				title, description, sequence, status, read, replied,
				priority, spam, channel, preview_text,
				inbound_message_id, outbound_message_id,
				created_at, updated_at
		) SELECT
		 	ups.thread_id AS thread_id,
			ups.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ups.title AS title,
			ups.description AS description,
			ups.sequence AS sequence,
			ups.status AS status,
			ups.read AS read,
			ups.replied AS replied,
			ups.priority AS priority,
			ups.spam AS spam,
			ups.channel AS channel,
			ups.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id,
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
            inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON ups.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON ups.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
	`

	err := tc.db.QueryRow(ctx, stmt, assigneeId, threadId).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	// set the inbound message if an inbound message exists.
	if inboundMessageId.Valid {
		thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
			inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
			inboundCreatedAt.Time,
			inboundUpdatedAt.Time,
		)
	} else {
		thread.ClearInboundMessage()
	}

	// set the outbound message if an outbound message exists.
	if outboundMessageId.Valid {
		thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
			outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
			outboundCreatedAt.Time,
			outboundUpdatedAt.Time,
		)
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func (tc *ThreadChatDB) UpdateRepliedState(
	ctx context.Context, threadId string, replied bool) (models.Thread, error) {
	var thread models.Thread

	var (
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

	stmt := `
		WITH ups AS (
			UPDATE thread
			SET replied = $1, updated_at = NOW()
			WHERE thread_id = $2
			RETURNING
				thread_id, workspace_id, customer_id, assignee_id,
				title, description, sequence, status, read, replied,
				priority, spam, channel, preview_text,
				inbound_message_id, outbound_message_id,
				created_at, updated_at
		) SELECT
		 	ups.thread_id AS thread_id,
			ups.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ups.title AS title,
			ups.description AS description,
			ups.sequence AS sequence,
			ups.status AS status,
			ups.read AS read,
			ups.replied AS replied,
			ups.priority AS priority,
			ups.spam AS spam,
			ups.channel AS channel,
			ups.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id,
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
            inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON ups.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON ups.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
	`

	err := tc.db.QueryRow(ctx, stmt, replied, threadId).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, ErrQuery
	}

	// set the inbound message if an inbound message exists.
	if inboundMessageId.Valid {
		thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
			inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
			inboundCreatedAt.Time,
			inboundUpdatedAt.Time,
		)
	} else {
		thread.ClearInboundMessage()
	}

	// set the outbound message if an outbound message exists.
	if outboundMessageId.Valid {
		thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
			outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
			outboundCreatedAt.Time,
			outboundUpdatedAt.Time,
		)
	} else {
		thread.ClearOutboundMessage()
	}
	return thread, nil
}

func (tc *ThreadChatDB) FetchThreadsByWorkspaceId(
	ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	args := make([]interface{}, 0, 3)
	var (
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

	stmt := `SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			th.title AS title,
			th.description AS description,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id, 
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
			inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
		WHERE th.workspace_id = $1
	`

	args = append(args, workspaceId)

	if channel != nil {
		stmt += " AND th.channel = $2"
		args = append(args, *channel)
	}

	if role != nil {
		stmt += " AND c.role = $3"
		args = append(args, *role)
	}

	stmt += " ORDER BY inb.last_seq_id DESC LIMIT 100"

	rows, _ := tc.db.Query(ctx, stmt, args...)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// set the inbound message if an inbound message exists.
		if inboundMessageId.Valid {
			thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
				inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
				inboundCreatedAt.Time,
				inboundUpdatedAt.Time,
			)
		} else {
			thread.ClearInboundMessage()
		}

		// set the outbound message if an outbound message exists.
		if outboundMessageId.Valid {
			thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
				outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
				outboundCreatedAt.Time,
				outboundUpdatedAt.Time,
			)
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

func (tc *ThreadChatDB) FetchThreadsByAssignedMemberId(
	ctx context.Context, memberId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	args := make([]interface{}, 0, 3)

	var (
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

	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			th.title AS title,
			th.description AS description,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id, 
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
			inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
		WHERE th.assignee_id = $1
	`

	args = append(args, memberId)

	if channel != nil {
		stmt += " AND th.channel = $2"
		args = append(args, *channel)
	}

	if role != nil {
		stmt += " AND c.role = $3"
		args = append(args, *role)
	}

	stmt += " ORDER BY inb.last_seq_id DESC LIMIT 100"

	rows, _ := tc.db.Query(ctx, stmt, args...)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// set the inbound message if an inbound message exists.
		if inboundMessageId.Valid {
			thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
				inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
				inboundCreatedAt.Time,
				inboundUpdatedAt.Time,
			)
		} else {
			thread.ClearInboundMessage()
		}

		// set the outbound message if an outbound message exists.
		if outboundMessageId.Valid {
			thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
				outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
				outboundCreatedAt.Time,
				outboundUpdatedAt.Time,
			)
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

func (tc *ThreadChatDB) FetchThreadsByMemberUnassigned(
	ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	args := make([]interface{}, 0, 3)

	var (
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

	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			th.title AS title,
			th.description AS description,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id, 
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
			inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
		WHERE th.workspace_id = $1 AND th.assignee_id IS NULL
	`

	args = append(args, workspaceId)

	if channel != nil {
		stmt += " AND th.channel = $2"
		args = append(args, *channel)
	}

	if role != nil {
		stmt += " AND c.role = $3"
		args = append(args, *role)
	}

	stmt += " ORDER BY inb.last_seq_id DESC LIMIT 100"

	rows, _ := tc.db.Query(ctx, stmt, args...)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// set the inbound message if an inbound message exists.
		if inboundMessageId.Valid {
			thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
				inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
				inboundCreatedAt.Time,
				inboundUpdatedAt.Time,
			)
		} else {
			thread.ClearInboundMessage()
		}

		// set the outbound message if an outbound message exists.
		if outboundMessageId.Valid {
			thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
				outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
				outboundCreatedAt.Time,
				outboundUpdatedAt.Time,
			)
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

func (tc *ThreadChatDB) FetchThreadsByLabelId(
	ctx context.Context, labelId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	args := make([]interface{}, 0, 3)
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			th.title AS title,
			th.description AS description,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.preview_text AS preview_text,
			inb.message_id,
			inbc.customer_id,
			inbc.name,
			inb.preview_text,
			inb.first_seq_id,
			inb.last_seq_id,
			inb.created_at,
            inb.updated_at,
			oub.message_id,
			oubm.member_id,
			oubm.name,
			oub.preview_text,
			oub.first_seq_id,
			oub.last_seq_id,
			oub.created_at,
			oub.updated_at,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id	
		INNER JOIN thread_label tl ON th.thread_id = tl.thread_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message inb ON th.inbound_message_id = inb.message_id
		LEFT OUTER JOIN outbound_message oub ON th.outbound_message_id = oub.message_id
		LEFT OUTER JOIN customer inbc ON inb.customer_id = inbc.customer_id
		LEFT OUTER JOIN member oubm ON oub.member_id = oubm.member_id
		WHERE tl.label_id = $1
	`

	args = append(args, labelId)

	if channel != nil {
		stmt += " AND th.channel = $2"
		args = append(args, *channel)
	}

	if role != nil {
		stmt += " AND c.role = $3"
		args = append(args, *role)
	}

	stmt += " ORDER BY inb.last_seq_id DESC LIMIT 100"

	rows, _ := tc.db.Query(ctx, stmt, args...)

	defer rows.Close()

	var (
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

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&inboundMessageId, &inboundCustomerId, &inboundCustomerName,
		&inboundPreviewText, &inboundFirstSeqId, &inboundLastSeqId,
		&inboundCreatedAt, &inboundUpdatedAt,
		&outboundMessageId, &outboundMemberId, &outboundMemberName,
		&outboundPreviewText, &outboundFirstSeqId, &outboundLastSeqId,
		&outboundCreatedAt, &outboundUpdatedAt,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		// set the inbound message if an inbound message exists.
		if inboundMessageId.Valid {
			thread.AddInboundMessage(inboundMessageId.String, inboundCustomerId.String, inboundCustomerName.String,
				inboundPreviewText.String, inboundFirstSeqId.String, inboundLastSeqId.String,
				inboundCreatedAt.Time,
				inboundUpdatedAt.Time,
			)
		} else {
			thread.ClearInboundMessage()
		}

		// set the outbound message if an outbound message exists.
		if outboundMessageId.Valid {
			thread.AddOutboundMessage(outboundMessageId.String, outboundMemberId.String, outboundMemberName.String,
				outboundPreviewText.String, outboundFirstSeqId.String, outboundLastSeqId.String,
				outboundCreatedAt.Time,
				outboundUpdatedAt.Time,
			)
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

func (tc *ThreadChatDB) CheckWorkspaceExistenceByThreadId(
	ctx context.Context, workspaceId string, threadId string) (bool, error) {
	var isExist bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM thread
		WHERE workspace_id = $1 AND thread_id= $2
	)`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, threadId).Scan(&isExist)
	if err != nil {
		slog.Error("failed to query", slog.Any("error", err))
		return false, ErrQuery
	}

	return isExist, nil
}

func (tc *ThreadChatDB) SetLabelToThread(
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

	err := tc.db.QueryRow(ctx, stmt, thLabelId, threadLabel.ThreadId, threadLabel.LabelId, threadLabel.AddedBy).Scan(
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

func (tc *ThreadChatDB) DeleteThreadLabelByCompId(
	ctx context.Context, threadId string, labelId string) error {
	stmt := `
		delete from thread_label
		where thread_id = $1 and label_id = $2`
	_, err := tc.db.Exec(ctx, stmt, threadId, labelId)
	if err != nil {
		slog.Error("failed to delete query", slog.Any("err", err))
		return ErrQuery
	}
	return nil
}

func (tc *ThreadChatDB) RetrieveLabelsByThreadId(
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

	rows, _ := tc.db.Query(ctx, stmt, threadId)

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

func (tc *ThreadChatDB) InsertCustomerChat(
	ctx context.Context, thread models.Thread, inboundMessage models.InboundMessage, chat models.Chat) (models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return chat, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	// insert new chat into the database
	chatId := chat.GenId()
	stmt := `
		with ins as (
			insert into chat (chat_id, thread_id, body, customer_id, is_head)
			values ($1, $2, $3, $4, $5)
			returning
				chat_id, thread_id, body, sequence, customer_id, member_id, is_head,
				created_at, updated_at
		) select ins.chat_id as chat_id,
			ins.thread_id as thread_id,
			ins.body as body,
			ins.sequence as sequence,
			c.customer_id as customer_id,
			c.name as customer_name,
			m.member_id as member_id,
			m.name as member_name,
			ins.is_head as is_head,
			ins.created_at as created_at,
			ins.updated_at as updated_at
		from ins
		left outer join customer c on ins.customer_id = c.customer_id
		left outer join member m on ins.member_id = m.member_id
	`

	err = tx.QueryRow(ctx, stmt, chatId, chat.ThreadId, chat.Body, chat.CustomerId, chat.IsHead).Scan(
		&chat.ChatId, &chat.ThreadId, &chat.Body,
		&chat.Sequence, &chat.CustomerId, &chat.CustomerName,
		&chat.MemberId, &chat.MemberName,
		&chat.IsHead, &chat.CreatedAt, &chat.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Chat{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Chat{}, ErrQuery
	}

	// insert or update the inbound message based on messageId
	stmt = `
		with ups as (
			insert into inbound_message (message_id, customer_id, first_seq_id, last_seq_id, preview_text)
				values ($1, $2, $3, $4, $5)
			on conflict (message_id) do update
				set last_seq_id = $4, preview_text = $5, updated_at = now()
			returning
				message_id, customer_id,
				first_seq_id, last_seq_id,
				preview_text,
				created_at, updated_at
		)
		select 
			u.message_id, 
			c.customer_id, 
			c.name,
			u.first_seq_id, 
			u.last_seq_id, 
			u.preview_text, 
			u.created_at, 
			u.updated_at
		from ups u
		inner join customer c on u.customer_id = c.customer_id
	`
	err = tx.QueryRow(ctx, stmt, inboundMessage.MessageId, inboundMessage.CustomerId,
		inboundMessage.FirstSeqId, inboundMessage.LastSeqId, inboundMessage.PreviewText).Scan(
		&inboundMessage.MessageId,
		&inboundMessage.CustomerId, &inboundMessage.CustomerName,
		&inboundMessage.FirstSeqId, &inboundMessage.LastSeqId,
		&inboundMessage.PreviewText,
		&inboundMessage.CreatedAt, &inboundMessage.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Chat{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Chat{}, ErrQuery
	}

	// check if the thread has the reference to inbound message,
	// if not then update thread with lastest inbound message ID.
	if thread.InboundMessage == nil {
		stmt = `update thread set
			inbound_message_id = $2, updated_at = now()
			where thread_id = $1`

		_, err = tx.Exec(ctx, stmt, chat.ThreadId, inboundMessage.MessageId)
		if err != nil {
			slog.Error("failed to update thread", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Chat{}, ErrTxQuery
	}

	return chat, nil
}

func (tc *ThreadChatDB) InsertMemberChat(
	ctx context.Context, thread models.Thread, outboundMessage models.OutboundMessage, chat models.Chat) (models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return chat, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	// insert new chat into the database
	chatId := chat.GenId()
	stmt := `
		with ins as (
			insert into chat (chat_id, thread_id, body, member_id, is_head)
			values ($1, $2, $3, $4, $5)
			returning
				chat_id, thread_id, body, sequence, customer_id, member_id, is_head,
				created_at, updated_at
		) select ins.chat_id as chat_id,
			ins.thread_id as thread_id,
			ins.body as body,
			ins.sequence as sequence,
			c.customer_id as customer_id,
			c.name as customer_name,
			m.member_id as member_id,
			m.name as member_name,
			ins.is_head as is_head,
			ins.created_at as created_at,
			ins.updated_at as updated_at
		from ins
		left outer join customer c on ins.customer_id = c.customer_id
		left outer join member m on ins.member_id = m.member_id
	`

	err = tx.QueryRow(ctx, stmt, chatId, chat.ThreadId, chat.Body, chat.MemberId, chat.IsHead).Scan(
		&chat.ChatId, &chat.ThreadId, &chat.Body,
		&chat.Sequence, &chat.CustomerId, &chat.CustomerName,
		&chat.MemberId, &chat.MemberName,
		&chat.IsHead, &chat.CreatedAt, &chat.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Chat{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Chat{}, ErrQuery
	}

	// insert or update the outbound message by message ID.
	stmt = `
		with ups as (
			insert into outbound_message (message_id, member_id, first_sequence, last_sequence, preview_text)
				values ($1, $2, $3, $4, $5)
			on conflict (message_id) do update
				set last_sequence = $4, preview_text = $5, updated_at = now()
			returning
				message_id, member_id, first_sequence, last_sequence, preview_text, created_at, updated_at
		)
		select
			u.message_id,
			m.member_id,
			m.name,
			u.first_seq_id,
			u.last_seq_id,
			u.preview_text,
			u.created_at,
			u.updated_at
		from ups u
		inner join member m on u.member_id = m.member_id
	`

	err = tx.QueryRow(ctx, stmt, outboundMessage.MessageId, outboundMessage.MemberId,
		outboundMessage.FirstSeqId, outboundMessage.LastSeqId, outboundMessage.PreviewText).Scan(
		&outboundMessage.MessageId,
		&outboundMessage.MemberId, &outboundMessage.MemberName,
		&outboundMessage.FirstSeqId, &outboundMessage.LastSeqId,
		&outboundMessage.PreviewText,
		&outboundMessage.CreatedAt, &outboundMessage.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Chat{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Chat{}, ErrQuery
	}

	// check if the thread has the reference to outbound message,
	// if not then update thread with lastest outbound message ID.
	if thread.OutboundMessage == nil {
		stmt = `update thread set
			outbound_message_id = $2, updated_at = now()
			where thread_id = $1`

		_, err = tx.Exec(ctx, stmt, chat.ThreadId, outboundMessage.MessageId)
		if err != nil {
			slog.Error("failed to update thread", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", slog.Any("err", err))
		return models.Chat{}, ErrTxQuery
	}
	return chat, nil
}

func (tc *ThreadChatDB) FetchThChatMessagesByThreadId(
	ctx context.Context, threadId string) ([]models.Chat, error) {
	var message models.Chat
	messages := make([]models.Chat, 0, 100)
	stmt := `SELECT
			ch.chat_id AS chat_id,
			ch.thread_id AS thread_id,
			ch.body AS body,
			ch.sequence AS sequence,
			chc.customer_id AS customer_id,
			chc.name AS customer_name,
			chm.member_id AS member_id,
			chm.name AS member_name,
			ch.is_head AS is_head,
			ch.created_at AS created_at,
			ch.updated_at AS updated_at
		FROM chat ch
		LEFT OUTER JOIN customer chc ON ch.customer_id = chc.customer_id
		LEFT OUTER JOIN member chm ON ch.member_id = chm.member_id
		WHERE ch.thread_id = $1
		ORDER BY sequence DESC LIMIT 100
	`

	rows, _ := tc.db.Query(ctx, stmt, threadId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&message.ChatId, &message.ThreadId, &message.Body,
		&message.Sequence, &message.CustomerId, &message.CustomerName,
		&message.MemberId, &message.MemberName, &message.IsHead,
		&message.CreatedAt, &message.UpdatedAt,
	}, func() error {
		messages = append(messages, message)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return []models.Chat{}, ErrQuery
	}

	return messages, nil
}

func (tc *ThreadChatDB) ComputeStatusMetricsByWorkspaceId(
	ctx context.Context, workspaceId string) (models.ThreadMetrics, error) {
	var metrics models.ThreadMetrics
	role := models.Customer{}.Engaged()
	stmt := `SELECT
		COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0) AS done_count,
		COALESCE(SUM(CASE WHEN status = 'todo' THEN 1 ELSE 0 END), 0) AS todo_count,
		COALESCE(SUM(CASE WHEN status = 'snoozed' THEN 1 ELSE 0 END), 0) AS snoozed_count,
		COALESCE(SUM(CASE WHEN status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS active_count
	FROM 
		thread th
	INNER JOIN customer c ON th.customer_id = c.customer_id
	WHERE 
		th.workspace_id = $1 AND c.role = $2`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, role).Scan(
		&metrics.DoneCount, &metrics.TodoCount,
		&metrics.SnoozedCount, &metrics.ActiveCount,
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

func (tc *ThreadChatDB) ComputeAssigneeMetricsByMember(
	ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error) {
	var metrics models.ThreadAssigneeMetrics
	role := models.Customer{}.Engaged()
	stmt := `SELECT
			COALESCE(SUM(
				CASE WHEN assignee_id = $2 AND status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS member_assigned_count,
			COALESCE(SUM(
				CASE WHEN assignee_id IS NULL AND status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS unassigned_count,
			COALESCE(SUM(
				CASE WHEN assignee_id IS NOT NULL AND assignee_id <> $2 AND status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS other_assigned_count
		FROM
			thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		WHERE
			th.workspace_id = $1 AND c.role = $3`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, memberId, role).Scan(
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

func (tc *ThreadChatDB) ComputeLabelMetricsByWorkspaceId(
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

	rows, _ := tc.db.Query(ctx, stmt, workspaceId)

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
