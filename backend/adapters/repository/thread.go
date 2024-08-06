package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

func (tc *ThreadChatDB) InsertInAppThreadChat(
	ctx context.Context, th models.Thread, chat models.Chat) (models.Thread, models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return models.Thread{}, chat, ErrQuery
	}

	defer tx.Rollback(ctx)

	threadId := th.GenId()
	stmt := `
		WITH ins AS (
			INSERT INTO thread (thread_id, workspace_id, customer_id, assignee_id,
				title, summary, status, read, replied, priority, spam,
				channel, message_body, message_customer_id, message_member_id
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING
				thread_id, workspace_id, customer_id, assignee_id,
				title, summary, sequence, status, read, replied, priority, spam,
				channel, messsage_body, message_sequence,
				message_customer_id, message_member_id,
				created_at, updated_at
		) SELECT ins.thread_id AS thread_id,
			ins.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ins.title AS title,
			ins.summary AS summary,
			ins.sequence AS sequence,
			ins.status AS status,
			ins.read AS read,
			ins.replied AS replied,
			ins.priority AS priority,
			ins.spam AS spam,
			ins.channel AS channel,
			ins.message_body AS message_body,
			ins.message_sequence AS message_sequence,
			mc.customer_id AS message_customer_id,
			mc.name AS message_customer_name,
			mm.member_id AS message_member_id,
			mm.name AS message_member_name,
			ins.created_at AS created_at,
			ins.updated_at AS updated_at
		FROM ins
		INNER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.assignee_id = m.member_id
		LEFT OUTER JOIN customer mc ON ins.message_customer_id = mc.customer_id
		LEFT OUTER JOIN member mm ON ins.message_member_id = mm.member_id
	`

	err = tx.QueryRow(ctx, stmt, threadId, th.WorkspaceId, th.CustomerId, th.AssigneeId,
		th.Title, th.Summary, th.Status, th.Read, th.Replied, th.Priority,
		th.Spam, th.Channel, th.MessageBody, th.MessageCustomerId, th.MessageMemberId).Scan(
		&th.ThreadId, &th.WorkspaceId,
		&th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.Title, &th.Summary,
		&th.Sequence, &th.Status,
		&th.Read, &th.Replied,
		&th.Priority, &th.Spam,
		&th.Channel, &th.MessageBody,
		&th.MessageSequence,
		&th.MessageCustomerId, &th.MessageCustomerName,
		&th.MessageMemberId, &th.MessageMemberName,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Thread{}, chat, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.Thread{}, chat, ErrQuery
	}

	chatId := chat.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO
		chat (chat_id, thread_id, body, sequence, customer_id, member_id, is_head)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			chat_id, thread_id, body, sequence,
			customer_id, member_id, is_head,
			created_at, updated_at`,
		chatId, th.ThreadId, chat.Body, th.MessageSequence,
		chat.CustomerId, chat.MemberId, chat.IsHead).Scan(
		&chat.ChatId, &chat.ThreadId, &chat.Body,
		&chat.Sequence, &chat.CustomerId, &chat.MemberId, &chat.IsHead,
		&chat.CreatedAt, &chat.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Thread{}, models.Chat{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.Thread{}, models.Chat{}, ErrQuery
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", "error", err)
		return models.Thread{}, models.Chat{}, ErrTxQuery
	}

	return th, chat, nil
}

func (tc *ThreadChatDB) ModifyThreadById(
	ctx context.Context, thread models.Thread, fields []string) (models.Thread, error) {

	args := make([]interface{}, 0, len(fields))
	ups := `UPDATE thread SET`

	for i, field := range fields {
		if field == "priority" {
			args = append(args, thread.Priority)
			ups += fmt.Sprintf(" %s = $%d,", "priority", i+1)
		} else if field == "assignee" {
			args = append(args, thread.AssigneeId)
			ups += fmt.Sprintf(" %s = $%d,", "assignee_id", i+1)
		} else if field == "status" {
			args = append(args, thread.Status)
			ups += fmt.Sprintf(" %s = $%d,", "status", i+1)
		}
	}

	ups += " updated_at = NOW()"
	ups += fmt.Sprintf(" WHERE thread_id = $%d", len(fields)+1)
	args = append(args, thread.ThreadId)

	stmt := `WITH ups AS (
		%s
		RETURNING
			thread_id, workspace_id, customer_id, assignee_id,
			title, summary, sequence, status, read, replied, priority, spam,
			channel, messsage_body, message_sequence,
			message_customer_id, message_member_id,
			created_at, updated_at
		) SELECT
	 		ups.thread_id AS thread_id,
			ups.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ups.title AS title,
			ups.summary AS summary,
			ups.sequence AS sequence,
			ups.status AS status,
			ups.read AS read,
			ups.replied AS replied,
			ups.priority AS priority,
			ups.spam AS spam,
			ups.channel AS channel,
			ups.message_body AS message_body,
			ups.message_sequence AS message_sequence,
			mc.customer_id AS message_customer_id,
			mc.name AS message_customer_name,
			mm.member_id AS message_member_id,
			mm.name AS message_member_name,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN customer mc ON ups.message_customer_id = mc.customer_id
		LEFT OUTER JOIN member mm ON ups.message_member_id = mm.member_id
	`

	stmt = fmt.Sprintf(stmt, ups)

	err := tc.db.QueryRow(ctx, stmt, args...).Scan(
		&thread.WorkspaceId, &thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.ThreadId, &thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status, &thread.Read,
		&thread.Replied, &thread.Priority, &thread.CreatedAt,
		&thread.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Thread{}, ErrEmpty
	}

	if err != nil {

		slog.Error("failed to update query", "error", err)
		return models.Thread{}, ErrQuery
	}

	return thread, nil
}

func (tc *ThreadChatDB) LookupByWorkspaceThreadId(
	ctx context.Context, workspaceId string, threadId string) (models.Thread, error) {
	var thread models.Thread
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			thc.customer_id AS customer_id,
			thc.name AS customer_name,
			tha.member_id AS assignee_id,
			tha.name AS assignee_name,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.message_body AS message_body,
			th.message_sequence AS message_sequence,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread th
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON th.message_customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON th.message_member_id = thmm.member_id
		WHERE th.workspace_id = $1 AND th.thread_id = $2
	`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, threadId).Scan(
		&thread.WorkspaceId, &thread.ThreadId, &thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName, &thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status, &thread.Read,
		&thread.Replied, &thread.Priority, &thread.Spam,
		&thread.Channel, &thread.MessageBody, &thread.MessageSequence,
		&thread.MessageCustomerId, &thread.MessageCustomerName,
		&thread.MessageMemberId, &thread.MessageMemberName,
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

func (tc *ThreadChatDB) RetrieveWorkspaceThChatsByCustomerId(
	ctx context.Context, workspaceId string, customerId string,
) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	channel := models.ThreadChannel{}.Chat()
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			thc.customer_id AS customer_id,
			thc.name AS customer_name,
			tha.member_id AS assignee_id,
			tha.name AS assignee_name,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.message_body AS message_body,
			th.message_sequence AS message_sequence,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread th
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON th.message_customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON th.message_member_id = thmm.member_id
		WHERE th.workspace_id = $1 AND th.customer_id = $2 AND channel = $3
		ORDER BY message_sequence DESC LIMIT 100
	`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, customerId, channel)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status,
		&thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam,
		&thread.Channel, &thread.MessageBody,
		&thread.MessageSequence,
		&thread.MessageCustomerId, &thread.MessageCustomerName,
		&thread.MessageMemberId, &thread.MessageMemberName,
		&thread.CreatedAt, &thread.UpdatedAt,
	}, func() error {
		threads = append(threads, thread)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.Thread{}, ErrQuery
	}

	return threads, nil
}

func (tc *ThreadChatDB) UpdateAssignee(
	ctx context.Context, threadId string, assigneeId string) (models.Thread, error) {
	var thread models.Thread
	stmt := `WITH ups AS (
			UPDATE thread
			SET assignee_id = $1, updated_at = NOW()
			WHERE thread_id = $2
			RETURNING
			thread_id, workspace_id, customer_id, assignee_id,
			title, summary, sequence, status, read, replied, priority, spam,
			channel, messsage_body, message_sequence,
			message_customer_id, message_member_id,
			created_at, updated_at
		) SELECT
			ups.thread_id AS thread_id,
			ups.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ups.title AS title,
			ups.summary AS summary,
			ups.sequence AS sequence,
			ups.status AS status,
			ups.read AS read,
			ups.replied AS replied,
			ups.priority AS priority,
			ups.spam AS spam,
			ups.channel AS channel,
			ups.message_body AS message_body,
			ups.message_sequence AS message_sequence,
			mc.customer_id AS message_customer_id,
			mc.name AS message_customer_name,
			mm.member_id AS message_member_id,
			mm.name AS message_member_name,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN customer mc ON ups.message_customer_id = mc.customer_id
		LEFT OUTER JOIN member mm ON ups.message_member_id = mm.member_id
	`

	err := tc.db.QueryRow(ctx, stmt, assigneeId, threadId).Scan(
		&thread.WorkspaceId, &thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.ThreadId, &thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status, &thread.Read,
		&thread.Replied, &thread.Priority, &thread.CreatedAt,
		&thread.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return thread, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return thread, ErrQuery
	}

	return thread, nil
}

func (tc *ThreadChatDB) UpdateRepliedStatus(ctx context.Context, threadId string, replied bool,
) (models.Thread, error) {
	var thread models.Thread
	stmt := `WITH ups AS (
			UPDATE thread
			SET replied = $1, updated_at = NOW()
			WHERE thread_id = $2
			RETURNING
			thread_id, workspace_id, customer_id, assignee_id,
			title, summary, sequence, status, read, replied, priority, spam,
			channel, messsage_body, message_sequence,
			message_customer_id, message_member_id,
			created_at, updated_at
		) SELECT
			ups.thread_id AS thread_id,
			ups.workspace_id AS workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ups.title AS title,
			ups.summary AS summary,
			ups.sequence AS sequence,
			ups.status AS status,
			ups.read AS read,
			ups.replied AS replied,
			ups.priority AS priority,
			ups.spam AS spam,
			ups.channel AS channel,
			ups.message_body AS message_body,
			ups.message_sequence AS message_sequence,
			mc.customer_id AS message_customer_id,
			mc.name AS message_customer_name,
			mm.member_id AS message_member_id,
			mm.name AS message_member_name,
			ups.created_at AS created_at,
			ups.updated_at AS updated_at
		FROM ups
		INNER JOIN customer c ON ups.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id
		LEFT OUTER JOIN customer mc ON ups.message_customer_id = mc.customer_id
		LEFT OUTER JOIN member mm ON ups.message_member_id = mm.member_id
	`

	err := tc.db.QueryRow(ctx, stmt, replied, threadId).Scan(
		&thread.WorkspaceId, &thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.ThreadId, &thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status, &thread.Read,
		&thread.Replied, &thread.Priority, &thread.CreatedAt,
		&thread.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Thread{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.Thread{}, ErrQuery
	}

	return thread, nil
}

func (tc *ThreadChatDB) FetchThChatsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			thc.customer_id AS customer_id,
			thc.name AS customer_name,
			tha.member_id AS assignee_id,
			tha.name AS assignee_name,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.message_body AS message_body,
			th.message_sequence AS message_sequence,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread th
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON th.message_customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON th.message_member_id = thmm.member_id
		WHERE th.workspace_id = $1 AND channel = $2 AND thc.role = $3
		ORDER BY message_sequence DESC LIMIT 100
	`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, channel, role)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status,
		&thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam,
		&thread.Channel, &thread.MessageBody,
		&thread.MessageSequence,
		&thread.MessageCustomerId, &thread.MessageCustomerName,
		&thread.MessageMemberId, &thread.MessageMemberName,
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

func (tc *ThreadChatDB) FetchAssignedThChatsByMemberId(
	ctx context.Context, memberId string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			thc.customer_id AS customer_id,
			thc.name AS customer_name,
			tha.member_id AS assignee_id,
			tha.name AS assignee_name,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.message_body AS message_body,
			th.message_sequence AS message_sequence,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread th
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON th.message_customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON th.message_member_id = thmm.member_id
		WHERE th.assignee_id = $1 AND channel = $2 AND thc.role = $3
		ORDER BY message_sequence DESC LIMIT 100
	`

	rows, _ := tc.db.Query(ctx, stmt, memberId, channel, role)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status,
		&thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam,
		&thread.Channel, &thread.MessageBody,
		&thread.MessageSequence,
		&thread.MessageCustomerId, &thread.MessageCustomerName,
		&thread.MessageMemberId, &thread.MessageMemberName,
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

func (tc *ThreadChatDB) FetchUnassignedThChatsByWorkspaceId(
	ctx context.Context, workspaceId string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			thc.customer_id AS customer_id,
			thc.name AS customer_name,
			tha.member_id AS assignee_id,
			tha.name AS assignee_name,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.message_body AS message_body,
			th.message_sequence AS message_sequence,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread th
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON th.message_customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON th.message_member_id = thmm.member_id
		WHERE th.workspace_id = $1 AND channel = $2 AND th.assignee_id IS NULL AND thc.role = $3
		ORDER BY message_sequence DESC LIMIT 100
	`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, channel, role)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status,
		&thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam,
		&thread.Channel, &thread.MessageBody,
		&thread.MessageSequence,
		&thread.MessageCustomerId, &thread.MessageCustomerName,
		&thread.MessageMemberId, &thread.MessageMemberName,
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

func (tc *ThreadChatDB) FetchThChatsByLabelId(
	ctx context.Context, labelId string) ([]models.Thread, error) {
	var thread models.Thread
	threads := make([]models.Thread, 0, 100)
	channel := models.ThreadChannel{}.Chat()
	role := models.Customer{}.Engaged()
	stmt := `
		SELECT th.thread_id AS thread_id,
			th.workspace_id AS workspace_id,
			thc.customer_id AS customer_id,
			thc.name AS customer_name,
			tha.member_id AS assignee_id,
			tha.name AS assignee_name,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.spam AS spam,
			th.channel AS channel,
			th.message_body AS message_body,
			th.message_sequence AS message_sequence,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread th
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON th.message_customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON th.message_member_id = thmm.member_id
		INNER JOIN thread_label tl ON th.thread_id = tl.thread_id
		WHERE tl.label_id = $1 AND th.channel = $2 AND thc.role = $3
		ORDER BY message_sequence DESC LIMIT 100
	`

	rows, _ := tc.db.Query(ctx, stmt, labelId, channel, role)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&thread.ThreadId, &thread.WorkspaceId,
		&thread.CustomerId, &thread.CustomerName,
		&thread.AssigneeId, &thread.AssigneeName,
		&thread.Title, &thread.Summary,
		&thread.Sequence, &thread.Status,
		&thread.Read, &thread.Replied,
		&thread.Priority, &thread.Spam,
		&thread.Channel, &thread.MessageBody,
		&thread.MessageSequence,
		&thread.MessageCustomerId, &thread.MessageCustomerName,
		&thread.MessageMemberId, &thread.MessageMemberName,
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
		) SELECT * FROM ins
		UNION ALL
		SELECT
			tl.thread_label_id AS thread_label_id,
			tl.thread_id AS thread_id,
			l.label_id AS label_id,
			l.name AS label_name,
			tl.addedby AS addedby,
			tl.created_at AS created_at,
			tl.updated_at AS updated_at,
			FALSE AS is_created
		FROM thread_label tl
		INNER JOIN label l ON tl.label_id = l.label_id
		WHERE tl.thread_id = $2 AND l.label_id = $3 AND NOT EXISTS (SELECT 1 FROM ins)
	`

	err := tc.db.QueryRow(ctx, stmt,
		thLabelId, threadLabel.ThreadId, threadLabel.LabelId, threadLabel.AddedBy).Scan(
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

func (tc *ThreadChatDB) InsertCustomerMessage(ctx context.Context, chat models.Chat) (models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return chat, ErrQuery
	}

	defer tx.Rollback(ctx)

	chatId := chat.GenId()
	stmt := `
		WITH ins AS (
			INSERT INTO chat (chat_id, thread_id, body, customer_id, is_head)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING
				chat_id, thread_id, body, sequence, customer_id, member_id, is_head,
				created_at, updated_at
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

	err = tx.QueryRow(ctx, stmt, chatId, chat.ThreadId, chat.Body, chat.CustomerId, chat.IsHead).Scan(
		&chat.ChatId, &chat.ThreadId, &chat.Body,
		&chat.Sequence, &chat.CustomerId, &chat.CustomerName,
		&chat.MemberId, &chat.MemberName,
		&chat.IsHead, &chat.CreatedAt, &chat.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Chat{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.Chat{}, ErrQuery
	}

	// update thread
	stmt = `UPDATE thread SET
		message_body = $1,
		message_sequence = $2,
		message_customer_id = $3,
		message_member_id = $4 WHERE thread_id = $5
	`
	previewBody := chat.PreviewBody()
	_, err = tx.Exec(ctx, stmt, previewBody, chat.Sequence, chat.CustomerId, chat.MemberId, chat.ThreadId)
	if err != nil {
		slog.Error("failed to update thread", "error", err)
		return models.Chat{}, ErrQuery
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", "error", err)
		return models.Chat{}, ErrTxQuery
	}

	return chat, nil
}

func (tc *ThreadChatDB) InsertThChatMemberMessage(
	ctx context.Context, chat models.Chat) (models.Chat, error) {
	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", slog.Any("err", err))
		return chat, ErrQuery
	}

	defer tx.Rollback(ctx)

	chatId := chat.GenId()
	stmt := `
		WITH ins AS (
			INSERT INTO chat (chat_id, thread_id, body, member_id, is_head)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING
				chat_id, thread_id, body, sequence, customer_id, member_id, is_head,
				created_at, updated_at
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

	// update thread
	stmt = `UPDATE thread SET
		message_body = $1,
		message_sequence = $2,
		message_customer_id = $3,
		message_member_id = $4 WHERE thread_id = $5
	`
	previewBody := chat.PreviewBody()
	_, err = tx.Exec(ctx, stmt, previewBody, chat.Sequence, chat.CustomerId, chat.MemberId, chat.ThreadId)
	if err != nil {
		slog.Error("failed to update thread", slog.Any("err", err))
		return models.Chat{}, ErrQuery
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

func (tc *ThreadChatDB) ComputeStatusMetricsByWorkspaceId(ctx context.Context, workspaceId string,
) (models.ThreadMetrics, error) {
	var metrics models.ThreadMetrics
	role := models.Customer{}.Engaged()
	stmt := `SELECT
		COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0) AS done,
		COALESCE(SUM(CASE WHEN status = 'todo' THEN 1 ELSE 0 END), 0) AS todo,
		COALESCE(SUM(CASE WHEN status = 'snoozed' THEN 1 ELSE 0 END), 0) AS snoozed,
		COALESCE(SUM(CASE WHEN status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS active
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
			COALESCE(SUM(CASE WHEN assignee_id = $2 THEN 1 ELSE 0 END), 0) AS member_assigned_count,
			COALESCE(SUM(CASE WHEN assignee_id IS NULL THEN 1 ELSE 0 END), 0) AS unassigned_count,
			COALESCE(SUM(CASE WHEN assignee_id IS NOT NULL AND assignee_id <> $2 THEN 1 ELSE 0 END), 0) AS other_assigned_count
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

	stmt := `SELECT
		l.label_id,
		l.name AS label_name,
		l.icon AS label_icon,
		COUNT(tl.thread_id) AS count
	FROM
		label l
	LEFT JOIN
		thread_label tl ON l.label_id = tl.label_id
	WHERE
		l.workspace_id = $1
	GROUP BY
		l.label_id, l.name
	ORDER BY MAX(tcl.updated_at) DESC
	LIMIT 100`

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
