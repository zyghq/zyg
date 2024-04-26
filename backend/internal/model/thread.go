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
	var message ThreadChatMessage
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

	rows, _ := db.Query(ctx, stmt, thm.ThreadChatId)

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
		return []ThreadChatMessage{}, ErrQuery
	}

	return messages, nil
}

func (th ThreadChat) GenId() string {
	return "th_" + xid.New().String()
}

// creates a customer thread chat
func (th ThreadChat) CreateCustomerThChat(
	ctx context.Context, db *pgxpool.Pool, msg string,
) (ThreadChat, ThreadChatMessage, error) {
	var message ThreadChatMessage

	tx, err := db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return th, message, ErrQuery
	}

	defer tx.Rollback(ctx)

	thId := th.GenId()
	th.Status = ThreadStatus{}.Todo() // set status
	err = tx.QueryRow(ctx, `INSERT INTO
		thread_chat(
			workspace_id, customer_id, thread_chat_id,
			title, summary, status, read, replied
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
		workspace_id, customer_id, assignee_id,
		thread_chat_id, title, summary, sequence, status, read,
		created_at, updated_at`,
		th.WorkspaceId, th.CustomerId, thId,
		th.Title, th.Summary, th.Status, th.Read, th.Replied).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.AssigneeId,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return th, message, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return th, message, ErrQuery
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
		return th, message, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return th, message, ErrQuery
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", "error", err)
		return th, message, ErrQuery
	}

	// set customer attributes we already have
	message.CustomerName = th.CustomerName
	return th, message, nil
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
		th.read AS read,
		th.replied AS replied,
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
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
		&th.CreatedAt, &th.UpdatedAt,
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
	var message ThreadChatMessage
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
			th.read AS read,
			th.replied AS replied,
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

	rows, _ := db.Query(ctx, stmt, th.WorkspaceId, th.CustomerId)

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []ThreadChatWithMessage{}, ErrQuery
	}

	defer rows.Close()

	return ths, nil
}

func (th ThreadChat) AssignMember(ctx context.Context, db *pgxpool.Pool) (ThreadChat, error) {
	stmt := `WITH ups AS (
		UPDATE thread_chat
		SET assignee_id = $1, updated_at = NOW()
		WHERE thread_chat_id = $2
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, read, replied,
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
		ups.created_at AS created_at,
		ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	err := db.QueryRow(ctx, stmt, th.AssigneeId, th.ThreadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
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

func (th ThreadChat) MarkReplied(ctx context.Context, db *pgxpool.Pool) (ThreadChat, error) {
	stmt := `WITH ups AS (
		UPDATE thread_chat
		SET replied = $1, updated_at = NOW()
		WHERE thread_chat_id = $2
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, read, replied,
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
		ups.created_at AS created_at,
		ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	err := db.QueryRow(ctx, stmt, th.Replied, th.ThreadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
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

func (th ThreadChat) GetListByWorkspace(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatWithMessage, error) {
	var message ThreadChatMessage
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
			th.read AS read,
			th.replied AS replied,
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

	rows, _ := db.Query(ctx, stmt, th.WorkspaceId)

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []ThreadChatWithMessage{}, ErrQuery
	}

	defer rows.Close()

	return ths, nil
}

func (th ThreadChat) IsExistInWorkspaceById(ctx context.Context, db *pgxpool.Pool) (bool, error) {
	var isExist bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM thread_chat
		WHERE thread_chat_id = $1 AND workspace_id = $2
	)`

	err := db.QueryRow(ctx, stmt, th.ThreadChatId, th.WorkspaceId).Scan(&isExist)

	if err != nil {
		slog.Error("failed to query", "error", err)
		return false, ErrQuery
	}

	return isExist, nil
}

func (l Label) GenId() string {
	return "l_" + xid.New().String()
}

func (l Label) GetOrCreate(ctx context.Context, db *pgxpool.Pool) (Label, bool, error) {
	var isCreated bool
	lId := l.GenId()
	stmt := `WITH ins AS (
		INSERT INTO label (label_id, workspace_id, name, icon)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (workspace_id, name) DO NOTHING
		RETURNING label_id, workspace_id, name, icon, created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT label_id, workspace_id, name, icon,
	created_at, updated_at, FALSE AS is_created FROM label
	WHERE workspace_id = $2 AND name = $3 AND NOT EXISTS (SELECT 1 FROM ins)`

	err := db.QueryRow(ctx, stmt, lId, l.WorkspaceId, l.Name, l.Icon).Scan(
		&l.LabelId, &l.WorkspaceId, &l.Name, &l.Icon,
		&l.CreatedAt, &l.UpdatedAt, &isCreated,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return l, isCreated, ErrEmpty
	}

	// query error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return l, isCreated, ErrQuery
	}

	return l, isCreated, nil
}

func (thl ThreadChatLabel) GenId() string {
	// no prefix required
	return xid.New().String()
}

func (thl ThreadChatLabel) Add(ctx context.Context, db *pgxpool.Pool) (ThreadChatLabel, bool, error) {
	var IsCreated bool
	thId := thl.GenId()

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

	err := db.QueryRow(ctx, stmt, thId, thl.ThreadChatId, thl.LabelId, thl.AddedBy).Scan(
		&thl.ThreadChatLabelId, &thl.ThreadChatId, &thl.LabelId, &thl.AddedBy,
		&thl.CreatedAt, &thl.UpdatedAt, &IsCreated,
	)

	// no rows returned
	if errors.Is(err, pgx.ErrNoRows) {
		return thl, IsCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return thl, IsCreated, ErrQuery
	}

	return thl, IsCreated, nil
}

func (thl ThreadChatLabel) GetListByThreadChatId(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatLabelled, error) {
	var l ThreadChatLabelled
	labels := make([]ThreadChatLabelled, 0, 100)
	stmt := `SELECT
			tcl.thread_chat_label_id AS thread_chat_label_id,
			tcl.thread_chat_id AS thread_chat_id,
			tcl.label_id AS label_id,
			l.name AS name,
			l.icon AS icon,
			tcl.addedby AS addedby,
			tcl.created_at AS created_at,
			tcl.updated_at AS updated_at
		FROM thread_chat_label AS tcl
		INNER JOIN label AS l ON tcl.label_id = l.label_id
		WHERE tcl.thread_chat_id = $1
		ORDER BY tcl.created_at DESC LIMIT 100`

	rows, _ := db.Query(ctx, stmt, thl.ThreadChatId)

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
		return []ThreadChatLabelled{}, ErrQuery
	}

	defer rows.Close()

	return labels, nil
}
