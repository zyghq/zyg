package model

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

func (thm ThreadChatMessage) GenId() string {
	return "thm_" + xid.New().String()
}

func (thm ThreadChatMessage) CreateCustomerThChatMessage(
	ctx context.Context, db *pgxpool.Pool, c Customer, msg string,
) (ThreadChatMessage, error) {
	thmId := thm.GenId()
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

	err := db.QueryRow(ctx, stmt, thm.ThreadChatId, thmId, msg, c.CustomerId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, sql.ErrNoRows) {
		return thm, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return thm, ErrQuery
	}

	return thm, nil
}

func (thm ThreadChatMessage) CreateMemberThChatMessage(
	ctx context.Context, db *pgxpool.Pool, m Member, msg string,
) (ThreadChatMessage, error) {
	thmId := thm.GenId()
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

	err := db.QueryRow(ctx, stmt, thm.ThreadChatId, thmId, msg, m.MemberId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return thm, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert got error", "error", err)
		return thm, err
	}

	return thm, nil
}

func (thm ThreadChatMessage) GetListByThreadChatId(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatMessage, error) {
	messages := make([]ThreadChatMessage, 0, 100)
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

	rows, err := db.Query(ctx, stmt, thm.ThreadChatId)

	// checks if the query was infact sent to db
	if err != nil {
		slog.Error("failed to query", "error", err)
		return messages, ErrQuery
	}

	defer rows.Close()

	if !rows.Next() {
		return messages, ErrEmpty
	}

	// got some rows iterate over them
	for rows.Next() {
		var m ThreadChatMessage
		err = rows.Scan(
			&m.ThreadChatId, &m.ThreadChatMessageId, &m.Body, &m.Sequence,
			&m.CreatedAt, &m.UpdatedAt, &m.CustomerId, &m.CustomerName,
			&m.MemberId, &m.MemberName,
		)
		if err != nil {
			slog.Error("failed to scan", "error", err)
			return messages, ErrQuery
		}
		messages = append(messages, m)
	}

	// checks if there was an error during scanning
	// we might have returned rows but might have failed to scan
	// check if there are any errors during collecting rows
	if err = rows.Err(); err != nil {
		slog.Error("failed to collect rows", "error", err)
		return messages, ErrQuery
	}

	return messages, nil
}

func (th ThreadChat) GenId() string {
	return "th_" + xid.New().String()
}

// TODO: fix to use th instead of passing w, c
func (th ThreadChat) CreateCustomerThChat(
	ctx context.Context, db *pgxpool.Pool, w Workspace, c Customer, m string,
) (ThreadChat, ThreadChatMessage, error) {
	var thm ThreadChatMessage

	tx, err := db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return th, thm, ErrQuery
	}

	defer tx.Rollback(ctx)

	thId := th.GenId()
	th.Status = ThreadStatus{}.Todo() // set status
	err = tx.QueryRow(ctx, `INSERT INTO thread_chat(workspace_id, customer_id, thread_chat_id, title, summary, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
		workspace_id, customer_id, assignee_id,
		thread_chat_id, title, summary, sequence, status, created_at, updated_at`,
		w.WorkspaceId, c.CustomerId, thId, th.Title, th.Summary, th.Status).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.AssigneeId,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return th, thm, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return th, thm, ErrQuery
	}

	thmId := thm.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO thread_chat_message(thread_chat_id, thread_chat_message_id, body, customer_id)
		VALUES ($1, $2, $3, $4)
		RETURNING
		thread_chat_id, thread_chat_message_id, body, sequence, customer_id, member_id, created_at, updated_at`,
		th.ThreadChatId, thmId, m, c.CustomerId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.MemberId,
		&thm.CreatedAt, &thm.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return th, thm, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return th, thm, ErrQuery
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", "error", err)
		return th, thm, ErrQuery
	}

	// set customer attributes we already have
	th.CustomerName = c.Name
	thm.CustomerName = c.Name
	return th, thm, nil
}

func (th ThreadChat) GetById(ctx context.Context, db *pgxpool.Pool) (ThreadChat, error) {
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
		th.created_at AS created_at,
		th.updated_at AS updated_at
		FROM thread_chat th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member a ON th.assignee_id = a.member_id
		WHERE th.thread_chat_id = $1`

	err := db.QueryRow(ctx, stmt, th.ThreadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return th, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return th, ErrQuery
	}

	return th, nil
}

func (th ThreadChat) GetListByWorkspaceCustomerId(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatWithMessage, error) {
	ths := make([]ThreadChatWithMessage, 0, 100)

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
		ORDER BY sequence DESC LIMIT 100`

	rows, err := db.Query(ctx, stmt, th.WorkspaceId, th.CustomerId)

	// checks if the query was infact sent to db
	if err != nil {
		slog.Error("failed to query", "error", err)
		return ths, ErrQuery
	}

	defer rows.Close()

	// checks if there are no rows returned
	if !rows.Next() {
		return ths, ErrEmpty
	}

	// got some rows iterate over them
	for rows.Next() {
		var th ThreadChat
		var tc ThreadChatMessage
		err = rows.Scan(
			&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
			&th.AssigneeId, &th.AssigneeName,
			&th.ThreadChatId, &th.Title, &th.Summary,
			&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
			&tc.ThreadChatId, &tc.ThreadChatMessageId, &tc.Body,
			&tc.Sequence, &tc.CreatedAt, &tc.UpdatedAt,
			&tc.CustomerId, &tc.CustomerName, &tc.MemberId, &tc.MemberName,
		)
		if err != nil {
			slog.Error("failed to scan", "error", err)
			return ths, ErrQuery
		}

		thm := ThreadChatWithMessage{
			ThreadChat: th,
			Message:    tc,
		}
		ths = append(ths, thm)
	}

	// checks if there was an error during scanning
	// we might have returned rows but might have failed to scan
	// check if there are any errors during collecting rows
	if err = rows.Err(); err != nil {
		slog.Error("failed to collect rows", "error", err)
		return ths, ErrQuery
	}

	return ths, nil
}

func (th ThreadChat) AssignMember(ctx context.Context, db *pgxpool.Pool) (ThreadChat, error) {
	stmt := `WITH ups AS (
		UPDATE thread_chat SET assignee_id = $1
		WHERE thread_chat_id = $2
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, created_at, updated_at
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
		ups.created_at AS created_at,
		ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	err := db.QueryRow(ctx, stmt, th.AssigneeId, th.ThreadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
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

func (th ThreadChat) GetListByWorkspace(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatWithMessage, error) {
	ths := make([]ThreadChatWithMessage, 0, 100)
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
		WHERE th.workspace_id = $1 
		ORDER BY sequence DESC LIMIT 100`

	rows, err := db.Query(ctx, stmt, th.WorkspaceId)

	// checks if the query was infact sent to db
	if err != nil {
		slog.Error("failed to query", "error", err)
		return ths, ErrQuery
	}

	defer rows.Close()

	// checks if there are no rows returned
	if !rows.Next() {
		return ths, ErrEmpty
	}

	// got some rows iterate over them
	for rows.Next() {
		var th ThreadChat
		var tc ThreadChatMessage
		err = rows.Scan(
			&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
			&th.AssigneeId, &th.AssigneeName,
			&th.ThreadChatId, &th.Title, &th.Summary,
			&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
			&tc.ThreadChatId, &tc.ThreadChatMessageId, &tc.Body,
			&tc.Sequence, &tc.CreatedAt, &tc.UpdatedAt,
			&tc.CustomerId, &tc.CustomerName, &tc.MemberId, &tc.MemberName,
		)
		if err != nil {
			slog.Error("failed to scan", "error", err)
			return ths, ErrQuery
		}

		thm := ThreadChatWithMessage{
			ThreadChat: th,
			Message:    tc,
		}
		ths = append(ths, thm)
	}

	// checks if there was an error during scanning
	// we might have returned rows but might have failed to scan
	// check if there are any errors during collecting rows
	if err = rows.Err(); err != nil {
		slog.Error("failed to collect rows", "error", err)
		return ths, ErrQuery
	}

	return ths, nil
}
