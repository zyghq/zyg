package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/xid"
)

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

func GetEnv(key string) (string, error) {
	value, status := os.LookupEnv(key)
	if !status {
		return "", fmt.Errorf("env `%s` is not set", key)
	}
	return value, nil
}

type Workspace struct {
	WorkspaceId string    `json:"workspaceId"`
	AccountId   string    `json:"accountId"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LLMRRLog struct {
	WorkspaceId string        `json:"workspaceId"`
	RequestId   string        `json:"requestId"`
	Prompt      string        `json:"prompt"`
	Response    string        `json:"response"`
	Model       string        `json:"model"`
	Eval        sql.NullInt64 `json:"eval"`
}

type LLMRequestQuery struct {
	Q string `json:"q"`
}

type LLMRREval struct {
	Eval int `json:"eval"`
}

type LLM struct {
	WorkspaceId string
	Prompt      string
	RequestId   string
}

type LLMResponse struct {
	Text      string `json:"text"`
	RequestId string `json:"requestId"`
	Model     string `json:"model"`
}

// for now we are directly using the Ollama server
// later will update to use our `converse` server probably with grpc.
// similarly the `LLMResponse` will be updated to include other specific fields.
func (llm LLM) Generate() (LLMResponse, error) {

	var err error
	var response LLMResponse

	buf := new(bytes.Buffer)
	// for now this is specific to the Ollama server
	// will update to use our `converse` server probably with grpc.
	body := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  "llama2",
		Prompt: llm.Prompt,
		Stream: false,
	}

	err = json.NewEncoder(buf).Encode(&body)
	if err != nil {
		log.Printf("failed to encode LLM request body for requestId: %s with error: %v", llm.RequestId, err)
		return response, err
	}

	log.Printf("LLM request for workspaceId: %s with requestId: %s", llm.WorkspaceId, llm.RequestId)
	resp, err := http.Post("http://0.0.0.0:11434/api/generate", "application/json", buf)
	if err != nil {
		log.Printf("LLM request failed for requestId: %s with error: %v", llm.RequestId, err)
		return response, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected status %d; but got %d", http.StatusOK, resp.StatusCode)
	}

	// response structure from Ollama server
	// will be updated based on the `converse` server response.
	rb := struct {
		Model              string `json:"model"`
		CreatedAt          string `json:"created_at"`
		Response           string `json:"response"`
		Done               bool   `json:"done"`
		Context            []int  `json:"context"`
		TotalDuration      int    `json:"total_duration"`
		LoadDuration       int    `json:"load_duration"`
		PromptEvalCount    int    `json:"prompt_eval_count"`
		PromptEvalDuration int    `json:"prompt_eval_duration"`
		EvalCount          int    `json:"eval_count"`
		EvalDuration       int    `json:"eval_duration"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&rb)
	if err != nil {
		log.Printf("failed to decode LLM response for requestId: %s with error: %v", llm.RequestId, err)
		return response, err
	}

	return LLMResponse{
		Text:      rb.Response,
		RequestId: llm.RequestId,
		Model:     rb.Model,
	}, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("TODO: add authentication logic before invoking the next handler...")
		next.ServeHTTP(w, r)
	})
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleGetWorkspaces(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(ctx, `SELECT
			workspace_id, account_id,
			name, created_at, updated_at
			FROM workspace ORDER BY created_at
			DESC LIMIT 100
		`)
		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		workspaces := make([]Workspace, 0)
		for rows.Next() {
			var workspace Workspace
			err = rows.Scan(
				&workspace.WorkspaceId, &workspace.AccountId,
				&workspace.Name,
				&workspace.CreatedAt, &workspace.UpdatedAt,
			)
			if err != nil {
				log.Printf("error: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			workspaces = append(workspaces, workspace)
		}

		if err := rows.Err(); err != nil {
			log.Printf("error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(workspaces); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleLLMQuery(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var workspace Workspace

		var query LLMRequestQuery
		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		workspaceId := r.PathValue("workspaceId")

		row, err := db.Query(ctx, `SELECT workspace_id, account_id,
			name, created_at, updated_at
			FROM workspace WHERE workspace_id = $1`,
			workspaceId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer row.Close()

		if !row.Next() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = row.Scan(
			&workspace.WorkspaceId, &workspace.AccountId,
			&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		requestId := xid.New()
		wrkLLM := LLM{WorkspaceId: workspaceId, Prompt: query.Q, RequestId: requestId.String()}

		resp, err := wrkLLM.Generate()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(ctx, `INSERT INTO
			llm_rr_log(workspace_id, request_id, prompt, response, model)
			VALUES ($1, $2, $3, $4, $5)`,
			workspace.WorkspaceId, wrkLLM.RequestId, wrkLLM.Prompt, resp.Text, resp.Model)
		if err != nil {
			log.Printf("failed to insert into llm request response log with error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleLLMQueryEval(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var workspace Workspace

		var eval LLMRREval
		err := json.NewDecoder(r.Body).Decode(&eval)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		workspaceId := r.PathValue("workspaceId")

		row, err := db.Query(ctx, `SELECT workspace_id, account_id,
			name, created_at, updated_at
			FROM workspace WHERE workspace_id = $1`,
			workspaceId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer row.Close()

		if !row.Next() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = row.Scan(
			&workspace.WorkspaceId, &workspace.AccountId,
			&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		requestId := r.PathValue("requestId")

		_, err = db.Exec(ctx, `UPDATE llm_rr_log SET eval = $1
			WHERE workspace_id = $2 AND request_id = $3`,
			eval.Eval, workspace.WorkspaceId, requestId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var err error

	pgConnStr, err := GetEnv("POSTGRES_URI")
	if err != nil {
		return err
	}

	db, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		return fmt.Errorf("unable to create pg connection pool: %v", err)
	}

	defer db.Close()

	var tm time.Time

	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)

	if err != nil {
		return fmt.Errorf("failed to query database: %v", err)
	}

	log.Printf("database ready with db time: %s\n", tm.Format(time.RFC1123))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handleGetIndex)

	mux.Handle("GET /workspaces/{$}", authenticatedOnly(handleGetWorkspaces(ctx, db)))
	mux.Handle("POST /workspaces/{workspaceId}/-/queries/{$}", authenticatedOnly(handleLLMQuery(ctx, db)))
	mux.Handle("POST /workspaces/{workspaceId}/-/queries/{requestId}/{$}", authenticatedOnly(handleLLMQueryEval(ctx, db)))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	srv := &http.Server{
		Addr:              *addr,
		Handler:           LoggingMiddleware(c.Handler(mux)),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

	log.Printf("server up and running on %s", *addr)

	err = srv.ListenAndServe()

	return err

}

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
