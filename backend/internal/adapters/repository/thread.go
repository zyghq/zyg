package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg/internal/domain"
)

// creates and returns a new thread chat for customer
// a customer must exist to create a thread chat
func (tc *ThreadChatDB) CreateThreadChat(ctx context.Context, th domain.ThreadChat, msg string,
) (domain.ThreadChat, domain.ThreadChatMessage, error) {
	var message domain.ThreadChatMessage

	tx, err := tc.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start db tx", "error", err)
		return th, message, ErrQuery
	}

	defer tx.Rollback(ctx)

	thId := th.GenId()
	th.Status = domain.ThreadStatus{}.Todo() // set status
	err = tx.QueryRow(ctx, `INSERT INTO
		thread_chat(workspace_id, customer_id, thread_chat_id,
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
		return domain.ThreadChat{}, domain.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.ThreadChat{}, domain.ThreadChatMessage{}, ErrQuery
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
		return domain.ThreadChat{}, domain.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.ThreadChat{}, domain.ThreadChatMessage{}, ErrQuery

	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("failed to commit query", "error", err)
		return domain.ThreadChat{}, domain.ThreadChatMessage{}, ErrTxQuery
	}

	return th, message, nil
}

// returns a thread chat for the workspace and thread chat id
// a thread chat must exist in the workspace
func (tc *ThreadChatDB) GetByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string,
) (domain.ThreadChat, error) {
	var th domain.ThreadChat

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
		WHERE th.workspace_id = $1 AND th.thread_chat_id = $2`

	// QueryRow automatically closes the connection
	err := tc.db.QueryRow(ctx, stmt, workspaceId, threadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
		&th.CreatedAt, &th.UpdatedAt,
	)

	// check if query returned a row
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ThreadChat{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.ThreadChat{}, ErrQuery
	}

	return th, nil
}

// returns a list of thread chats with latest message for customer in workspace
// a thread chat message cannot exist without thread chat
func (tc *ThreadChatDB) GetListByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string,
) ([]domain.ThreadChatWithMessage, error) {
	var th domain.ThreadChat
	var message domain.ThreadChatMessage

	ths := make([]domain.ThreadChatWithMessage, 0, 100)
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

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, customerId)

	defer rows.Close()

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
		thm := domain.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

// assigns a member to a thread chat
// a member exist in the workspace
func (tc *ThreadChatDB) SetAssignee(ctx context.Context, threadChatId string, assigneeId string,
) (domain.ThreadChat, error) {
	var th domain.ThreadChat
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
		ups.title AS title,9
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

	err := tc.db.QueryRow(ctx, stmt, assigneeId, threadChatId).Scan(
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

// marks a thread chat as replied or not replied
func (tc *ThreadChatDB) SetReplied(ctx context.Context, threadChatId string, replied bool,
) (domain.ThreadChat, error) {
	var th domain.ThreadChat
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

	err := tc.db.QueryRow(ctx, stmt, replied, threadChatId).Scan(
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

// returns a list of thread chats with latest message for workspace
// irrespective of customer
// this is different from `GetListByWorkspaceCustomerId`
func (tc *ThreadChatDB) GetListByWorkspaceId(ctx context.Context, workspaceId string,
) ([]domain.ThreadChatWithMessage, error) {
	var th domain.ThreadChat
	var message domain.ThreadChatMessage

	ths := make([]domain.ThreadChatWithMessage, 0, 100)
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

	rows, _ := tc.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

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
		thm := domain.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

// returns a list of member assigned thread chats with latest message in workspace
func (tc *ThreadChatDB) GetMemberAssignedListByWorkspaceId(ctx context.Context, workspaceId string, memberId string,
) ([]domain.ThreadChatWithMessage, error) {
	var th domain.ThreadChat
	var message domain.ThreadChatMessage

	ths := make([]domain.ThreadChatWithMessage, 0, 100)
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
		WHERE th.workspace_id = $1 AND th.assignee_id = $2
		ORDER BY sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, memberId)

	defer rows.Close()

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
		thm := domain.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

// returns a list of unassigned thread chats with latest message in workspace
func (tc *ThreadChatDB) GetUnassignedListByWorkspaceId(ctx context.Context, workspaceId string,
) ([]domain.ThreadChatWithMessage, error) {
	var th domain.ThreadChat
	var message domain.ThreadChatMessage

	ths := make([]domain.ThreadChatWithMessage, 0, 100)
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
		WHERE th.workspace_id = $1 AND th.assignee_id IS NULL
		ORDER BY sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId)

	defer rows.Close()

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
		thm := domain.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil
	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

// returns a list of labelled thread chats with latest message in workspace
func (tc *ThreadChatDB) GetLabelledListByWorkspaceId(ctx context.Context, workspaceId string, labelId string) ([]domain.ThreadChatWithMessage, error) {
	var th domain.ThreadChat
	var message domain.ThreadChatMessage

	ths := make([]domain.ThreadChatWithMessage, 0, 100)
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
		INNER JOIN thread_chat_label tcl ON th.thread_chat_id = tcl.thread_chat_id
		WHERE th.workspace_id = $1 AND tcl.label_id = $2
		ORDER BY sequence DESC LIMIT 100`

	rows, _ := tc.db.Query(ctx, stmt, workspaceId, labelId)

	defer rows.Close()

	_, err := pgx.ForEachRow(rows, []any{
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId,
		&th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.Read, &th.Replied,
		&th.CreatedAt, &th.UpdatedAt,
		&message.ThreadChatId, &message.ThreadChatMessageId, &message.Body,
		&message.Sequence, &message.CreatedAt, &message.UpdatedAt,
		&message.CustomerId, &message.CustomerName, &message.MemberId, &message.MemberName,
	}, func() error {
		thm := domain.ThreadChatWithMessage{
			ThreadChat: th,
			Message:    message,
		}
		ths = append(ths, thm)
		return nil

	})

	if err != nil {
		slog.Error("failed to query", "error", err)
		return []domain.ThreadChatWithMessage{}, ErrQuery
	}

	return ths, nil
}

// checks if a thread chat exists in the workspace
func (tc *ThreadChatDB) IsExistByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string,
) (bool, error) {
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

// adds a label to a thread chat
func (tc *ThreadChatDB) AddLabel(ctx context.Context, thl domain.ThreadChatLabel) (domain.ThreadChatLabel, bool, error) {
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
		return domain.ThreadChatLabel{}, IsCreated, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.ThreadChatLabel{}, IsCreated, ErrQuery
	}
	return thl, IsCreated, nil
}

// returns a list of labels added for a thread chat item
func (tc *ThreadChatDB) GetLabelListByThreadChatId(ctx context.Context, threadChatId string,
) ([]domain.ThreadChatLabelled, error) {
	var l domain.ThreadChatLabelled
	labels := make([]domain.ThreadChatLabelled, 0, 100)
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
		return []domain.ThreadChatLabelled{}, ErrQuery
	}

	return labels, nil
}

// creates a thread chat message for a customer
// a thread chat message must belong to either a customer or a member not both
func (tc *ThreadChatDB) CreateCustomerThChatMessage(
	ctx context.Context, threadChatId string, customerId string,
	msg string,
) (domain.ThreadChatMessage, error) {
	var thm domain.ThreadChatMessage
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
		return domain.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.ThreadChatMessage{}, ErrQuery
	}

	return thm, nil
}

// creates a thread chat message for a member
// a thread chat message must belong to either a customer or a member not both
func (tc *ThreadChatDB) CreateMemberThChatMessage(
	ctx context.Context, threadChatId string, memberId string,
	msg string,
) (domain.ThreadChatMessage, error) {
	var thm domain.ThreadChatMessage
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
		return domain.ThreadChatMessage{}, ErrEmpty
	}

	// check if query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return domain.ThreadChatMessage{}, ErrQuery
	}

	return thm, nil
}

// returns a list of thread chat messages for a thread chat
func (tc *ThreadChatDB) GetMessageListByThreadChatId(ctx context.Context, threadChatId string,
) ([]domain.ThreadChatMessage, error) {
	var message domain.ThreadChatMessage
	messages := make([]domain.ThreadChatMessage, 0, 100)
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
		return []domain.ThreadChatMessage{}, ErrQuery
	}

	return messages, nil
}

// returns count of thread chats by status in workspace
func (tc *ThreadChatDB) StatusMetricsByWorkspaceId(ctx context.Context, workspaceId string,
) (domain.ThreadMetrics, error) {
	var metrics domain.ThreadMetrics

	stmt := `SELECT
		COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0) AS done,
		COALESCE(SUM(CASE WHEN status = 'todo' THEN 1 ELSE 0 END), 0) AS todo,
		COALESCE(SUM(CASE WHEN status = 'snoozed' THEN 1 ELSE 0 END), 0) AS snoozed,
		COALESCE(SUM(CASE WHEN status = 'todo' OR status = 'snoozed' THEN 1 ELSE 0 END), 0) AS active
	FROM 
		thread_chat
	WHERE 
		workspace_id = $1`

	err := tc.db.QueryRow(ctx, stmt, workspaceId).Scan(
		&metrics.DoneCount, &metrics.TodoCount,
		&metrics.SnoozedCount, &metrics.ActiveCount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ThreadMetrics{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.ThreadMetrics{}, ErrQuery
	}

	return metrics, nil
}

// returns count of thread chats assigned to a member in workspace
func (tc *ThreadChatDB) MemberAssigneeMetricsByWorkspaceId(ctx context.Context, workspaceId string, memberId string,
) (domain.ThreadAssigneeMetrics, error) {
	var metrics domain.ThreadAssigneeMetrics

	stmt := `SELECT
			COALESCE(SUM(CASE WHEN assignee_id = $2 THEN 1 ELSE 0 END), 0) AS member_assigned_count,
			COALESCE(SUM(CASE WHEN assignee_id IS NULL THEN 1 ELSE 0 END), 0) AS unassigned_count,
			COALESCE(SUM(CASE WHEN assignee_id IS NOT NULL AND assignee_id <> $2 THEN 1 ELSE 0 END), 0) AS other_assigned_count
		FROM
			thread_chat
		WHERE
			workspace_id = $1`

	err := tc.db.QueryRow(ctx, stmt, workspaceId, memberId).Scan(
		&metrics.MeCount, &metrics.UnAssignedCount, &metrics.OtherAssignedCount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ThreadAssigneeMetrics{}, ErrEmpty
	}

	if err != nil {
		slog.Error("failed to query", "error", err)
		return domain.ThreadAssigneeMetrics{}, ErrQuery
	}

	return metrics, nil
}

// returns count of thread chat for each label in workspace
func (tc *ThreadChatDB) LabelMetricsByWorkspaceId(ctx context.Context, workspaceId string,
) ([]domain.ThreadLabelMetric, error) {
	var metric domain.ThreadLabelMetric
	metrics := make([]domain.ThreadLabelMetric, 0, 100)

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
		return []domain.ThreadLabelMetric{}, ErrQuery
	}

	return metrics, nil
}
