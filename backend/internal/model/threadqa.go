package model

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

func (tq ThreadQA) MarshalJSON() ([]byte, error) {
	var pth *string
	if tq.ParentThreadId.Valid {
		pth = &tq.ParentThreadId.String
	}
	aux := &struct {
		WorkspaceId    string  `json:"workspaceId"`
		CustomerId     string  `json:"customerId"`
		ThreadId       string  `json:"threadId"`
		ParentThreadId *string `json:"parentThreadId"`
		Query          string  `json:"query"`
		Title          string  `json:"title"`
		Summary        string  `json:"summary"`
		Sequence       int     `json:"sequence"`
		CreatedAt      string  `json:"createdAt"`
		UpdatedAt      string  `json:"updatedAt"`
	}{
		WorkspaceId:    tq.WorkspaceId,
		CustomerId:     tq.CustomerId,
		ThreadId:       tq.ThreadId,
		ParentThreadId: pth,
		Query:          tq.Query,
		Title:          tq.Title,
		Summary:        tq.Summary,
		Sequence:       tq.Sequence,
		CreatedAt:      tq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      tq.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (tq ThreadQA) GenId() string {
	return "tq_" + xid.New().String()
}

func (tq ThreadQA) Create(ctx context.Context, db *pgxpool.Pool) (ThreadQA, error) {
	var thread ThreadQA

	tqId := tq.GenId()
	stmt := `INSERT INTO 
		thread_qa(workspace_id, customer_id, thread_id, parent_thread_id, query, title, summary)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING 
		workspace_id, customer_id, thread_id, parent_thread_id,
		query, title, summary, sequence,
		created_at, updated_at`

	err := db.QueryRow(ctx, stmt, tq.WorkspaceId, tq.CustomerId, tqId, tq.ParentThreadId, tq.Query, tq.Title, tq.Summary).Scan(
		&thread.WorkspaceId, &thread.CustomerId, &thread.ThreadId, &thread.ParentThreadId,
		&thread.Query, &thread.Title, &thread.Summary, &thread.Sequence,
		&thread.CreatedAt, &thread.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return thread, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return thread, ErrQuery
	}

	return thread, nil
}

func (tqa ThreadQAA) MarshalJSON() ([]byte, error) {
	var eval *int32
	if tqa.Eval.Valid {
		eval = &tqa.Eval.Int32
	}
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		ThreadQAId  string `json:"threadQAId"`
		AnswerId    string `json:"answerId"`
		Answer      string `json:"answer"`
		Sequence    int    `json:"sequence"`
		Eval        *int32 `json:"eval"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: tqa.WorkspaceId,
		ThreadQAId:  tqa.ThreadQAId,
		AnswerId:    tqa.AnswerId,
		Answer:      tqa.Answer,
		Sequence:    tqa.Sequence,
		Eval:        eval,
		CreatedAt:   tqa.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   tqa.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (tqa ThreadQAA) GenId() string {
	return "tqa_" + xid.New().String()
}

func (tqa ThreadQAA) Create(ctx context.Context, db *pgxpool.Pool) (ThreadQAA, error) {
	var thread ThreadQAA

	tqaId := tqa.GenId()
	stmt := `INSERT INTO 
		thread_qa_answer(workspace_id, thread_qa_id, answer_id, answer)
		VALUES ($1, $2, $3, $4)
		RETURNING 
		workspace_id, thread_qa_id, answer_id, answer, 
		eval, sequence, created_at, updated_at`

	err := db.QueryRow(ctx, stmt, tqa.WorkspaceId, tqa.ThreadQAId, tqaId, tqa.Answer).Scan(
		&thread.WorkspaceId, &thread.ThreadQAId, &thread.AnswerId, &thread.Answer,
		&thread.Eval, &thread.Sequence, &thread.CreatedAt, &thread.UpdatedAt,
	)

	// check if the query returned no rows
	if errors.Is(err, pgx.ErrNoRows) {
		return thread, ErrEmpty
	}

	// check if the query returned an error
	if err != nil {
		slog.Error("failed to insert query", "error", err)
		return thread, ErrQuery
	}

	return thread, nil
}
