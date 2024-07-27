package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/models"
)

// creates and returns a new thread chat for customer
// a customer must exist to create a thread chat
func (tc *ThreadChatDB) InsertThreadChat(ctx context.Context, th models.ThreadChat, msg string,
) (models.ThreadChat, models.ThreadChatMessage, error) {
	var message models.ThreadChatMessage

	// start transaction
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return th, message, ErrQuery
	}

	defer tx.Rollback(ctx)

	thId := th.GenId()

	// TODO: @sanchitrk - move this logic to service layer.
	if th.Status == "" {
		th.Status = models.ThreadStatus{}.DefaultStatus() // set status
	}

	if th.Priority == "" {
		th.Priority = models.ThreadPriority{}.DefaultPriority() // set priority
	}

	stmt := `WITH ins AS (
		INSERT INTO thread_chat (
			workspace_id,
			customer_id,
			thread_chat_id,
			title,
			summary, 
			status,
			priority
		)
		VALUES 
			(
				$1, $2, $3, $4, $5, $6, $7
			)
		RETURNING
			workspace_id,
			customer_id,
			assignee_id,
			thread_chat_id,
			title,
			summary,
			sequence,
			status,
			read,
			replied,
			priority,
			created_at,
			updated_at
		) SELECT ins.workspace_id as workspace_id,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS assignee_id,
			m.name AS assignee_name,
			ins.thread_chat_id,
			ins.title,
			ins.summary,
			ins.sequence,
			ins.status,
			ins.read,
			ins.replied,
			ins.priority,
			ins.created_at,
			ins.updated_at
		FROM ins
		INNER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.assignee_id = m.member_id`

	err = tx.QueryRow(ctx, stmt,
		th.WorkspaceId, th.CustomerId, thId,
		th.Title, th.Summary, th.Status, th.Priority).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChat{}, models.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.ThreadChat{}, models.ThreadChatMessage{}, ErrQuery
	}

	thmId := message.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO
		thread_chat_message(thread_chat_id, thread_chat_message_id, body, customer_id)
		VALUES ($1, $2, $3, $4)
		RETURNING
		thread_chat_id, thread_chat_message_id, body, sequence, customer_id, member_id,
		created_at, updated_at`,
		th.ThreadChatId, thmId, msg, th.CustomerId).Scan(
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CustomerId, &message.MemberId,
		&message.CreatedAt, &message.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChat{}, models.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.ThreadChat{}, models.ThreadChatMessage{}, ErrQuery

	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", "error", err)
		return models.ThreadChat{}, models.ThreadChatMessage{}, ErrTxQuery
	}

	return th, message, nil
}

// update a thread chat
// @sanchitrk!: fix
func (tc *ThreadChatDB) ModifyThreadChatById(ctx context.Context, th models.ThreadChat, fields []string,
) (models.ThreadChat, error) {

	args := make([]interface{}, 0, len(fields))

	ups := `UPDATE thread_chat SET`

	for i, field := range fields {
		if field == "priority" {
			args = append(args, th.Priority)
			ups += fmt.Sprintf(" %s = $%d,", "priority", i+1)
		} else if field == "assignee" {
			args = append(args, th.AssigneeId)
			ups += fmt.Sprintf(" %s = $%d,", "assignee_id", i+1)
		} else if field == "status" {
			args = append(args, th.Status)
			ups += fmt.Sprintf(" %s = $%d,", "status", i+1)
		}
	}

	ups += " updated_at = NOW()"
	ups += fmt.Sprintf(" WHERE thread_chat_id = $%d", len(fields)+1)
	args = append(args, th.ThreadChatId)

	stmt := `WITH ups AS (
		%s
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, read, replied, priority,
		created_at, updated_at
	) SELECT
		ups.workspace_id AS workspace_id,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS assignee_id,
		m.name AS assignee_name,
		ups.thread_chat_id AS thread_chat_id,
		ups.title AS title,
		ups.summary AS summary,
		ups.sequence AS sequence,
		ups.status AS status,
		ups.read AS read,
		ups.replied AS replied,
		ups.priority AS priority,
		ups.created_at AS created_at,
		ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	stmt = fmt.Sprintf(stmt, ups)

	err := tc.db.QueryRow(ctx, stmt, args...).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChat{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.ThreadChat{}, ErrQuery
	}

	return th, nil
}

// returns thread chat for the workspace
func (tc *ThreadChatDB) LookupByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string,
) (models.ThreadChat, error) {
	var th models.ThreadChat

	stmt := `SELECT th.workspace_id AS workspace_id,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		a.member_id AS assignee_id,
		a.name AS assignee_name,
		th.thread_chat_id AS thread_chat_id,
		th.title AS title,
		th.summary AS summary,
		th.sequence AS sequence,
		th.status AS status,
		th.read AS read,
		th.replied AS replied,
		th.priority AS priority,
		th.created_at AS created_at,
		th.updated_at AS updated_at
		FROM thread_chat th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member a ON th.assignee_id = a.member_id
		WHERE th.workspace_id = $1 AND th.thread_chat_id = $2`

	// QueryRow automatically closes the connection
	err := tc.db.QueryRow(ctx, stmt, workspaceId, threadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChat{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.ThreadChat{}, ErrQuery
	}

	return th, nil
}

// returns list of thread chats with latest message for the customer
func (tc *ThreadChatDB) FetchByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string,
) ([]models.ThreadChatWithMessage, error) {
	var th models.ThreadChat
	var message models.ThreadChatMessage

	ths := make([]models.ThreadChatWithMessage, 0, 100)
	// (@sanchitrk)
	//
	// shall we try this, for rendering list of labels.
	// https://dba.stackexchange.com/questions/173831/convert-right-side-of-join-of-many-to-many-into-array
	//
	// shall we do query profiling?
	stmt := `SELECT
			th.workspace_id AS workspace_id,	
			thc.customer_id AS thread_customer_id,
			thc.name AS thread_customer_name,
			tha.member_id AS thread_assignee_id,
			tha.name AS thread_assignee_name,
			th.thread_chat_id AS thread_chat_id,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.created_at AS created_at,
			th.updated_at AS updated_at,
			thm.thread_chat_id AS message_thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS message_body,
			thm.sequence AS message_sequence,
			thm.created_at AS message_created_at,
			thm.updated_at AS message_updated_at,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread_chat th
		INNER JOIN thread_chat_message thm ON th.thread_chat_id = thm.thread_chat_id
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON thm.customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON thm.member_id = thmm.member_id
		INNER JOIN (
			SELECT thread_chat_id, MAX(sequence) AS sequence
			FROM thread_chat_message
			GROUP BY
			thread_chat_id
		) latest ON thm.thread_chat_id = latest.thread_chat_id
		AND thm.sequence = latest.sequence
		WHERE th.workspace_id = $1 AND th.customer_id = $2
		ORDER BY message_sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, customerId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := models.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

// assign member to a existing thread chat
// a member exist in the workspace
func (tc *ThreadChatDB) UpdateAssignee(ctx context.Context, threadChatId string, assigneeId string,
) (models.ThreadChat, error) {
	var th models.ThreadChat
	stmt := `WITH ups AS (
		UPDATE thread_chat
		SET assignee_id = $1, updated_at = NOW()
		WHERE thread_chat_id = $2
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, read, replied, priority,
		created_at, updated_at
	) SELECT
		ups.workspace_id AS workspace_id,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS assignee_id,
		m.name AS assignee_name,
		ups.thread_chat_id AS thread_chat_id,
		ups.title AS title,
		ups.summary AS summary,
		ups.sequence AS sequence,
		ups.status AS status,
		ups.read AS read,
		ups.replied AS replied,
		ups.priority AS priority,
		ups.created_at AS created_at,
		ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	err := tc.db.QueryRow(ctx, stmt, assigneeId, threadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return th, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return th, ErrQuery
	}

	return th, nil
}

// marks a thread chat as replied or un-replied
func (tc *ThreadChatDB) UpdateRepliedStatus(ctx context.Context, threadChatId string, replied bool,
) (models.ThreadChat, error) {
	var th models.ThreadChat
	stmt := `WITH ups AS (
		UPDATE thread_chat
		SET replied = $1, updated_at = NOW()
		WHERE thread_chat_id = $2
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, read, replied, priority,
		created_at, updated_at
	) SELECT
		ups.workspace_id AS workspace_id,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS assignee_id,
		m.name AS assignee_name,
		ups.thread_chat_id AS thread_chat_id,
		ups.title AS title,
		ups.summary AS summary,
		ups.sequence AS sequence,
		ups.status AS status,
		ups.read AS read,
		ups.replied AS replied,
		ups.priority AS priority,
		ups.created_at AS created_at,
		ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	err := tc.db.QueryRow(ctx, stmt, replied, threadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return th, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return th, ErrQuery
	}

	return th, nil
}

func (tc *ThreadChatDB) RetrieveByWorkspaceId(ctx context.Context, workspaceId string,
) ([]models.ThreadChatWithMessage, error) {
	var th models.ThreadChat
	var message models.ThreadChatMessage
	cr := models.Customer{}.Engaged()
	ths := make([]models.ThreadChatWithMessage, 0, 100)
	stmt := `SELECT
			th.workspace_id AS workspace_id,	
			thc.customer_id AS thread_customer_id,
			thc.name AS thread_customer_name,
			tha.member_id AS thread_assignee_id,
			tha.name AS thread_assignee_name,
			th.thread_chat_id AS thread_chat_id,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.created_at AS created_at,
			th.updated_at AS updated_at,
			thm.thread_chat_id AS message_thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS message_body,
			thm.sequence AS message_sequence,
			thm.created_at AS message_created_at,
			thm.updated_at AS message_updated_at,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread_chat th
		INNER JOIN thread_chat_message thm ON th.thread_chat_id = thm.thread_chat_id
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON thm.customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON thm.member_id = thmm.member_id
		INNER JOIN (
			SELECT thread_chat_id, MAX(sequence) AS sequence
			FROM thread_chat_message
			GROUP BY
			thread_chat_id
		) latest ON thm.thread_chat_id = latest.thread_chat_id
		AND thm.sequence = latest.sequence
		WHERE th.workspace_id = $1 AND thc.role = $2
		ORDER BY message_sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, cr)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := models.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

func (tc *ThreadChatDB) FetchAssignedThreadsByMember(ctx context.Context, workspaceId string, memberId string,
) ([]models.ThreadChatWithMessage, error) {
	var th models.ThreadChat
	var message models.ThreadChatMessage
	cr := models.Customer{}.Engaged()
	ths := make([]models.ThreadChatWithMessage, 0, 100)
	stmt := `SELECT
			th.workspace_id AS workspace_id,	
			thc.customer_id AS thread_customer_id,
			thc.name AS thread_customer_name,
			tha.member_id AS thread_assignee_id,
			tha.name AS thread_assignee_name,
			th.thread_chat_id AS thread_chat_id,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.created_at AS created_at,
			th.updated_at AS updated_at,
			thm.thread_chat_id AS message_thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS message_body,
			thm.sequence AS message_sequence,
			thm.created_at AS message_created_at,
			thm.updated_at AS message_updated_at,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread_chat th
		INNER JOIN thread_chat_message thm ON th.thread_chat_id = thm.thread_chat_id
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON thm.customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON thm.member_id = thmm.member_id
		INNER JOIN (
			SELECT thread_chat_id, MAX(sequence) AS sequence
			FROM thread_chat_message
			GROUP BY
			thread_chat_id
		) latest ON thm.thread_chat_id = latest.thread_chat_id
		AND thm.sequence = latest.sequence
		WHERE th.workspace_id = $1 AND th.assignee_id = $2 AND thc.role = $3
		ORDER BY member_sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, memberId, cr)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := models.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

func (tc *ThreadChatDB) RetrieveUnassignedThreads(ctx context.Context, workspaceId string,
) ([]models.ThreadChatWithMessage, error) {
	var th models.ThreadChat
	var message models.ThreadChatMessage
	cr := models.Customer{}.Engaged()
	ths := make([]models.ThreadChatWithMessage, 0, 100)
	stmt := `SELECT
			th.workspace_id AS workspace_id,	
			thc.customer_id AS thread_customer_id,
			thc.name AS thread_customer_name,
			tha.member_id AS thread_assignee_id,
			tha.name AS thread_assignee_name,
			th.thread_chat_id AS thread_chat_id,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.created_at AS created_at,
			th.updated_at AS updated_at,
			thm.thread_chat_id AS message_thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS message_body,
			thm.sequence AS message_sequence,
			thm.created_at AS message_created_at,
			thm.updated_at AS message_updated_at,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread_chat th
		INNER JOIN thread_chat_message thm ON th.thread_chat_id = thm.thread_chat_id
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON thm.customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON thm.member_id = thmm.member_id
		INNER JOIN (
			SELECT thread_chat_id, MAX(sequence) AS sequence
			FROM thread_chat_message
			GROUP BY
			thread_chat_id
		) latest ON thm.thread_chat_id = latest.thread_chat_id
		AND thm.sequence = latest.sequence
		WHERE th.workspace_id = $1 AND th.assignee_id IS NULL AND thc.role = $2
		ORDER BY message_sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, cr)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := models.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

func (tc *ThreadChatDB) FetchThreadsByLabel(ctx context.Context, workspaceId string, labelId string,
) ([]models.ThreadChatWithMessage, error) {
	var th models.ThreadChat
	var message models.ThreadChatMessage
	cr := models.Customer{}.Engaged()
	ths := make([]models.ThreadChatWithMessage, 0, 100)
	stmt := `SELECT
			th.workspace_id AS workspace_id,	
			thc.customer_id AS thread_customer_id,
			thc.name AS thread_customer_name,
			tha.member_id AS thread_assignee_id,
			tha.name AS thread_assignee_name,
			th.thread_chat_id AS thread_chat_id,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.read AS read,
			th.replied AS replied,
			th.priority AS priority,
			th.created_at AS created_at,
			th.updated_at AS updated_at,
			thm.thread_chat_id AS message_thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS message_body,
			thm.sequence AS message_sequence,
			thm.created_at AS message_created_at,
			thm.updated_at AS message_updated_at,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread_chat th
		INNER JOIN thread_chat_message thm ON th.thread_chat_id = thm.thread_chat_id
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON thm.customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON thm.member_id = thmm.member_id
		INNER JOIN (
			SELECT thread_chat_id, MAX(sequence) AS sequence
			FROM thread_chat_message
			GROUP BY
			thread_chat_id
		) latest ON thm.thread_chat_id = latest.thread_chat_id
		AND thm.sequence = latest.sequence
		INNER JOIN thread_chat_label tcl ON th.thread_chat_id = tcl.thread_chat_id
		WHERE th.workspace_id = $1 AND tcl.label_id = $2 AND thc.role = $3
		ORDER BY message_sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, labelId, cr)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId,
		&th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied, &th.Priority,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := models.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

func (tc *ThreadChatDB) CheckExistenceByWorkspaceThreadChatId(
	ctx context.Context, workspaceId string, threadChatId string) (bool, error) {
	var isExist bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM thread_chat
		WHERE workspace_id = $1 AND thread_chat_id = $2
	)`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, threadChatId).Scan(&isExist)
	if err != nil {
		slog.Error("failed to query", "error", err)
		return false, ErrQuery
	}

	return isExist, nil
}

// add a label to a thread chat
func (tc *ThreadChatDB) AttachLabelToThread(
	ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error) {
	var IsCreated bool
	id := thl.GenId()

	stmt := `WITH ins AS (
		INSERT INTO thread_chat_label (thread_chat_label_id, thread_chat_id, label_id, addedby)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (thread_chat_id, label_id) DO NOTHING
		RETURNING thread_chat_label_id, thread_chat_id, label_id, addedby,
		created_at, updated_at, TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT thread_chat_label_id, thread_chat_id, label_id, addedby,
	created_at, updated_at, FALSE AS is_created FROM thread_chat_label
	WHERE thread_chat_id = $2 AND label_id = $3 AND NOT EXISTS (SELECT 1 FROM ins)`

	err := tc.db.QueryRow(ctx, stmt, id, thl.ThreadChatId, thl.LabelId, thl.AddedBy).Scan(
		&thl.ThreadChatLabelId, &thl.ThreadChatId, &thl.LabelId, &thl.AddedBy,
		&thl.CreatedAt, &thl.UpdatedAt, &IsCreated,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChatLabel{}, IsCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.ThreadChatLabel{}, IsCreated, ErrQuery
	}
	return thl, IsCreated, nil
}

// returns list of labels added to a thread chat
func (tc *ThreadChatDB) RetrieveLabelsByThreadChatId(ctx context.Context, threadChatId string,
) ([]models.ThreadChatLabelled, error) {
	var l models.ThreadChatLabelled
	labels := make([]models.ThreadChatLabelled, 0, 100)
	stmt := `SELECT tcl.thread_chat_label_id AS thread_chat_label_id,
		tcl.thread_chat_id AS thread_chat_id,
		tcl.label_id AS label_id,
		l.name AS name, l.icon AS icon,
		tcl.addedby AS addedby,
		tcl.created_at AS created_at,
		tcl.updated_at AS updated_at
		FROM thread_chat_label AS tcl
		INNER JOIN label AS l ON tcl.label_id = l.label_id
		WHERE tcl.thread_chat_id = $1
		ORDER BY tcl.created_at DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, threadChatId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&l.ThreadChatLabelId,
		&l.ThreadChatId,
		&l.LabelId,
		&l.Name, &l.Icon,
		&l.AddedBy,
		&l.CreatedAt,
		&l.UpdatedAt}, func() error {
		labels = append(labels, l)
		return nil
	})

	if err != nil {
		slog.Error("failed to scan", "error", err)
		return []models.ThreadChatLabelled{}, ErrQuery
	}

	return labels, nil
}

// creates a thread chat message for a customer
// a thread chat message must belong to either a customer or a member not both
func (tc *ThreadChatDB) InsertCustomerMessage(
	ctx context.Context, threadChatId string, customerId string,
	msg string,
) (models.ThreadChatMessage, error) {
	var thm models.ThreadChatMessage
	id := thm.GenId()
	stmt := `WITH ins AS (
		INSERT INTO thread_chat_message (thread_chat_id, thread_chat_message_id, body, customer_id)
			VALUES ($1, $2, $3, $4)
		RETURNING
			thread_chat_id, thread_chat_message_id, body, sequence,
			customer_id, member_id, created_at, updated_at
		) SELECT ins.thread_chat_id AS thread_chat_id,
			ins.thread_chat_message_id AS thread_chat_message_id,
			ins.body AS body,
			ins.sequence AS sequence,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS member_id,
			m.name AS member_name,
			ins.created_at AS created_at,
			ins.updated_at AS updated_at
		FROM ins
		LEFT OUTER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.member_id = m.member_id`

	err := tc.db.QueryRow(ctx, stmt, threadChatId, id, msg, customerId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.ThreadChatMessage{}, ErrQuery
	}

	return thm, nil
}

// creates a thread chat message for a member
// a thread chat message must belong to either a customer or a member not both
func (tc *ThreadChatDB) InsertMemberMessage(
	ctx context.Context, threadChatId string, memberId string,
	msg string,
) (models.ThreadChatMessage, error) {
	var thm models.ThreadChatMessage
	id := thm.GenId()
	stmt := `WITH ins AS (
		INSERT INTO thread_chat_message (thread_chat_id, thread_chat_message_id, body, member_id)
			VALUES ($1, $2, $3, $4)
		RETURNING
			thread_chat_id, thread_chat_message_id, body, sequence,
			customer_id, member_id, created_at, updated_at
		) SELECT ins.thread_chat_id AS thread_chat_id,
			ins.thread_chat_message_id AS thread_chat_message_id,
			ins.body AS body,
			ins.sequence AS sequence,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS member_id,
			m.name AS member_name,
			ins.created_at AS created_at,
			ins.updated_at AS updated_at
		FROM ins
		LEFT OUTER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.member_id = m.member_id`

	err := tc.db.QueryRow(ctx, stmt, threadChatId, id, msg, memberId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return models.ThreadChatMessage{}, ErrQuery
	}

	return thm, nil
}

// returns list of messages in desc order for the thread chat
func (tc *ThreadChatDB) FetchMessagesByThreadChatId(ctx context.Context, threadChatId string,
) ([]models.ThreadChatMessage, error) {
	var message models.ThreadChatMessage
	messages := make([]models.ThreadChatMessage, 0, 100)
	stmt := `SELECT
			thm.thread_chat_id AS thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS body,
			thm.sequence AS sequence,
			thm.created_at AS created_at,
			thm.updated_at AS updated_at,
			c.customer_id AS customer_id,
			c.name AS customer_name,
			m.member_id AS member_id,
			m.name AS member_name
		FROM thread_chat_message AS thm
		LEFT OUTER JOIN customer AS c ON thm.customer_id = c.customer_id
		LEFT OUTER JOIN member AS m ON thm.member_id = m.member_id
		WHERE thm.thread_chat_id = $1
		ORDER BY sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, threadChatId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		messages = append(messages, message)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []models.ThreadChatMessage{}, ErrQuery
	}

	return messages, nil
}

// returns stats of thread chats by status in workspace
func (tc *ThreadChatDB) ComputeStatusMetricsByWorkspaceId(ctx context.Context, workspaceId string,
) (models.ThreadMetrics, error) {
	var metrics models.ThreadMetrics
	cr := models.Customer{}.Engaged()
	stmt := `SELECT
		COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0) AS done,
		COALESCE(SUM(CASE WHEN status = 'todo' THEN 1 ELSE 0 END), 0) AS todo,
		COALESCE(SUM(CASE WHEN status = 'snoozed' THEN 1 ELSE 0 END), 0) AS snoozed,
		COALESCE(SUM(CASE WHEN status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS active
	FROM 
		thread_chat th
	INNER JOIN customer c ON th.customer_id = c.customer_id
	WHERE 
		th.workspace_id = $1 AND c.role = $2`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, cr).Scan(
		&metrics.DoneCount, &metrics.TodoCount,
		&metrics.SnoozedCount, &metrics.ActiveCount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadMetrics{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.ThreadMetrics{}, ErrQuery
	}

	return metrics, nil
}

// returns stats of thread chats assigned to a member in workspace
func (tc *ThreadChatDB) CalculateAssigneeMetricsByMember(ctx context.Context, workspaceId string, memberId string,
) (models.ThreadAssigneeMetrics, error) {
	var metrics models.ThreadAssigneeMetrics
	cr := models.Customer{}.Engaged()
	stmt := `SELECT
			COALESCE(SUM(CASE WHEN assignee_id = $2 THEN 1 ELSE 0 END), 0) AS member_assigned_count,
			COALESCE(SUM(CASE WHEN assignee_id IS NULL THEN 1 ELSE 0 END), 0) AS unassigned_count,
			COALESCE(SUM(CASE WHEN assignee_id IS NOT NULL AND assignee_id <> $2 THEN 1 ELSE 0 END), 0) AS other_assigned_count
		FROM
			thread_chat th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		WHERE
			th.workspace_id = $1 AND c.role = $3`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, memberId, cr).Scan(
		&metrics.MeCount, &metrics.UnAssignedCount, &metrics.OtherAssignedCount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ThreadAssigneeMetrics{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return models.ThreadAssigneeMetrics{}, ErrQuery
	}

	return metrics, nil
}

// returns stats for labelled thread chats with atmost 100 labels
func (tc *ThreadChatDB) ComputeLabelMetricsByWorkspaceId(ctx context.Context, workspaceId string,
) ([]models.ThreadLabelMetric, error) {
	var metric models.ThreadLabelMetric
	metrics := make([]models.ThreadLabelMetric, 0, 100)

	stmt := `SELECT
		l.label_id,
		l.name AS label_name,
		l.icon AS label_icon,
		COUNT(tcl.thread_chat_id) AS count
	FROM
		label l
	LEFT JOIN
		thread_chat_label tcl ON l.label_id = tcl.label_id
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
		slog.Error("failed to query", "error", err)
		return []models.ThreadLabelMetric{}, ErrQuery
	}

	return metrics, nil
}
