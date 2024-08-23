package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (tc *ThreadChatDB) InsertInboundThreadChat(
	ctx context.Context, inbound models.InboundMessage,
	thread models.Thread, chat models.Chat) (models.Thread, models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return models.Thread{}, chat, ErrQuery
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	messageId := inbound.GenId()
	stmt := `
		with ins as (
			insert into inbound_message (message_id, customer_id)
			values ($1, $2)
			returning
				message_id, customer_id,
				first_sequence, last_sequence,
				created_at, updated_at
		) select ins.message_id as message_id,
			ins.customer_id as customer_id,
			ins.first_sequence as first_sequence,
			ins.last_sequence as last_sequence,
			ins.created_at as created_at,
			ins.updated_at as updated_at
		from ins
	`

	err = tx.QueryRow(ctx, stmt, messageId, inbound.CustomerId).Scan(
		&inbound.MessageId, &inbound.CustomerId,
		&inbound.FirstMessageTs, &inbound.LastMessageTs,
		&inbound.CreatedAt, &inbound.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.Thread{}, models.Chat{}, ErrQuery
	}

	threadId := thread.GenId()
	stmt = `
		with ins as (
			insert into thread (
				thread_id, workspace_id, customer_id, assignee_id,
				title, description, status, read, replied,
				priority, spam, channel, preview_text,
				inbound_message_id, sequence
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			returning
				thread_id, workspace_id, customer_id, assignee_id,
				title, description, sequence, status, read, replied,
				priority, spam, channel, preview_text,
				inbound_message_id, outbound_message_id,
				created_at, updated_at
		) select ins.thread_id as thread_id,
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
			ing.message_id as inbound_message_id,
			ing.first_sequence as ingress_first_seq,
			ing.last_sequence as ingress_last_seq,
			ingc.customer_id as ingress_customer_id,
			ingc.name as ingress_customer_name,
			eg.message_id as outbound_message_id,
			eg.first_sequence as egress_first_seq,
			eg.last_sequence as egress_last_seq,
			egm.member_id as egress_member_id,
			egm.name as egress_member_name,
			ins.created_at as created_at,
			ins.updated_at as updated_at
		from ins
		inner join customer c on ins.customer_id = c.customer_id
		left outer join member m on ins.assignee_id = m.member_id
		left outer join inbound_message ing on ins.inbound_message_id = ing.message_id
		left outer join outbound_message eg on ins.outbound_message_id = eg.message_id
		left outer join customer ingc on ing.customer_id = ingc.customer_id
		left outer join member egm on eg.member_id = egm.member_id
	`

	err = tx.QueryRow(ctx, stmt, threadId, thread.WorkspaceId, thread.CustomerId, thread.AssigneeId,
		thread.Title, thread.Description, thread.Status, thread.Read, thread.Replied,
		thread.Priority, thread.Spam, thread.Channel, thread.PreviewText,
		inbound.MessageId, inbound.FirstMessageTs).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
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

	stmt = `
		WITH ins AS (
			INSERT INTO chat (
				chat_id, thread_id, body, sequence,
				customer_id, member_id, is_head
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING
				chat_id, thread_id, body, sequence,
				customer_id, member_id, is_head, created_at,
				updated_at
		) SELECT ins.chat_id AS chat_id,
			ins.thread_id AS thread_id,
			ins.body AS body,
			ins.sequence AS sequence,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS member_id,
			m.name AS member_name,
			ins.is_head AS is_head,
			ins.created_at AS created_at,
			ins.updated_at AS updated_at
		FROM ins
		LEFT OUTER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.member_id = m.member_id
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
		case "previewText":
			args = append(args, thread.PreviewText)
			ups += fmt.Sprintf(" %s = $%d,", "preview_text", i+1)
		}
	}

	ups += " updated_at = NOW()"
	ups += fmt.Sprintf(" WHERE thread_id = $%d", len(fields)+1)
	args = append(args, thread.ThreadId)

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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON ups.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON ups.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
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
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
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
		ing.message_id AS inbound_message_id,
		ing.first_sequence AS ingress_first_seq,
		ing.last_sequence AS ingress_last_seq,
		ingc.customer_id AS ingress_customer_id,
		ingc.name AS ingress_customer_name,
		eg.message_id AS outbound_message_id,
		eg.first_sequence AS egress_first_seq,
		eg.last_sequence AS egress_last_seq,
		egm.member_id AS egress_member_id,
		egm.name AS egress_member_name,
		th.created_at AS created_at,
		th.updated_at AS updated_at
	FROM thread th
	INNER JOIN customer c ON th.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
	LEFT OUTER JOIN inbound_message ing ON th.inbound_message_id = ing.message_id
	LEFT OUTER JOIN outbound_message eg ON th.outbound_message_id = eg.message_id
	LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
	LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
	WHERE th.workspace_id = $1 AND th.thread_id = $2`

	args = append(args, workspaceId)
	args = append(args, threadId)

	if channel != nil {
		stmt += " AND th.channel = $3"
		args = append(args, *channel)
	}

	err := tc.db.QueryRow(ctx, stmt, args...).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON th.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON th.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
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

	stmt += " ORDER BY ingress_last_seq DESC LIMIT 100"

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
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON ups.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON ups.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
	`

	err := tc.db.QueryRow(ctx, stmt, assigneeId, threadId).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
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

	return thread, nil
}

func (tc *ThreadChatDB) UpdateRepliedState(
	ctx context.Context, threadId string, replied bool) (models.Thread, error) {
	var thread models.Thread
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON ups.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON ups.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
	`

	err := tc.db.QueryRow(ctx, stmt, replied, threadId).Scan(
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Description, &thread.Sequence,
		&thread.Status, &thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam, &thread.Channel,
		&thread.PreviewText,
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
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

	return thread, nil
}

func (tc *ThreadChatDB) FetchThreadsByWorkspaceId(
	ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON th.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON th.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
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

	stmt += " ORDER BY ingress_last_seq DESC LIMIT 100"

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
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON th.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON th.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
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

	stmt += " ORDER BY ingress_last_seq DESC LIMIT 100"

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
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON th.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON th.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
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

	stmt += " ORDER BY ingress_last_seq DESC LIMIT 100"

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
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
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
			ing.message_id AS inbound_message_id,
			ing.first_sequence AS ingress_first_seq,
			ing.last_sequence AS ingress_last_seq,
			ingc.customer_id AS ingress_customer_id,
			ingc.name AS ingress_customer_name,
			eg.message_id AS outbound_message_id,
			eg.first_sequence AS egress_first_seq,
			eg.last_sequence AS egress_last_seq,
			egm.member_id AS egress_member_id,
			egm.name AS egress_member_name,
			th.created_at AS created_at,
			th.updated_at AS updated_at
		FROM thread th
		INNER JOIN customer c ON th.customer_id = c.customer_id	
		INNER JOIN thread_label tl ON th.thread_id = tl.thread_id
		LEFT OUTER JOIN member m ON th.assignee_id = m.member_id
		LEFT OUTER JOIN inbound_message ing ON th.inbound_message_id = ing.message_id
		LEFT OUTER JOIN outbound_message eg ON th.outbound_message_id = eg.message_id
		LEFT OUTER JOIN customer ingc ON ing.customer_id = ingc.customer_id
		LEFT OUTER JOIN member egm ON eg.member_id = egm.member_id
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

	stmt += " ORDER BY ingress_last_seq DESC LIMIT 100"

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
		&thread.IngressMessageId, &thread.IngressFirstSeq,
		&thread.IngressLastSeq, &thread.IngressCustomerId,
		&thread.IngressCustomerName,
		&thread.EgressMessageId, &thread.EgressFirstSeq,
		&thread.EgressLastSeq, &thread.EgressMemberId,
		&thread.EgressMemberName,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
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
	ctx context.Context, inboundMessageId *string, chat models.Chat) (models.Chat, error) {
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

	// check if the inbound message id is provided
	// if not, create a new one and update the thread
	// else update the existing one
	if inboundMessageId != nil {
		stmt = `update inbound_message set
			last_sequence = $2,
			updated_at = now()
			where message_id = $1`

		_, err = tx.Exec(ctx, stmt, *inboundMessageId, chat.Sequence)
		if err != nil {
			slog.Error("failed to update inbound message", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}

		// update thread with preview text
		stmt = `update thread set
			preview_text = $2,
			updated_at = now()
			where thread_id = $1`

		_, err = tx.Exec(ctx, stmt, chat.ThreadId, chat.PreviewText())
		if err != nil {
			slog.Error("failed to update thread", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}
	} else {
		var inbound models.InboundMessage
		messageId := inbound.GenId()
		stmt := `
			with ins as (
				insert into inbound_message (message_id, customer_id)
				values ($1, $2)
				returning
					message_id, customer_id,
					first_sequence, last_sequence,
					created_at, updated_at
			) select ins.message_id as message_id,
				ins.customer_id as customer_id,
				ins.first_sequence as first_sequence,
				ins.last_sequence as last_sequence,
				ins.created_at as created_at,
				ins.updated_at as updated_at
			from ins`

		err = tx.QueryRow(ctx, stmt, messageId, chat.CustomerId).Scan(
			&inbound.MessageId, &inbound.CustomerId,
			&inbound.FirstMessageTs, &inbound.LastMessageTs,
			&inbound.CreatedAt, &inbound.UpdatedAt,
		)

		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("no rows returned", slog.Any("err", err))
			return models.Chat{}, ErrEmpty
		}

		if err != nil {
			slog.Error("failed to insert query", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}

		// update thread with inbound message id and preview text
		stmt = `update thread set
			inbound_message_id = $2, preview_text = $3,
			updated_at = now()
			where thread_id = $1`

		_, err = tx.Exec(ctx, stmt, chat.ThreadId, inbound.MessageId, chat.PreviewText())
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
	ctx context.Context, outboundMessageId *string, chat models.Chat) (models.Chat, error) {
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

	// check if outbound message id is provided
	// if not, create a new one and update the thread
	// else update the existing one
	if outboundMessageId != nil {
		stmt = `update outbound_message set
		last_sequence = $2,
		updated_at = now()
		where message_id = $1`

		_, err = tx.Exec(ctx, stmt, *outboundMessageId, chat.Sequence)
		if err != nil {
			slog.Error("failed to update outbound message", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}

		// update thread with outbound message id and preview text
		stmt = `update thread set
			preview_text = $2,
			updated_at = now()
			where thread_id = $1`

		_, err = tx.Exec(ctx, stmt, chat.ThreadId, chat.PreviewText())
		if err != nil {
			slog.Error("failed to update thread", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}
	} else {
		var outbound models.OutboundMessage
		messageId := outbound.GenId()
		stmt := `
			WITH ins AS (
				INSERT INTO outbound_message (message_id, member_id)
				VALUES ($1, $2)
				RETURNING
					message_id, member_id,
					first_sequence, last_sequence,
					created_at, updated_at
			) SELECT ins.message_id AS message_id,
				ins.member_id AS member_id,
				ins.first_sequence AS first_sequence,
				ins.last_sequence AS last_sequence,
				ins.created_at AS created_at,
				ins.updated_at AS updated_at
			FROM ins`

		err = tx.QueryRow(ctx, stmt, messageId, chat.MemberId).Scan(
			&outbound.MessageId, &outbound.MemberId,
			&outbound.FirstMessageTs, &outbound.LastMessageTs,
			&outbound.CreatedAt, &outbound.UpdatedAt,
		)

		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("no rows returned", slog.Any("err", err))
			return models.Chat{}, ErrEmpty
		}

		if err != nil {
			slog.Error("failed to insert query", slog.Any("err", err))
			return models.Chat{}, ErrQuery
		}

		// update thread with outbound message id and preview text
		stmt = `update thread set
			outbound_message_id = $2, preview_text = $3,
			updated_at = now()
			where thread_id = $1`

		_, err = tx.Exec(ctx, stmt, chat.ThreadId, outbound.MessageId, chat.PreviewText())
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
