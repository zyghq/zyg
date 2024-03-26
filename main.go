package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

type Workspace struct {
	WorkspaceId string    `json:"workspaceId"`
	AccountId   string    `json:"accountId"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got a %s request for : %v", r.Method, r.URL.Path)
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleGetWorkspaces(ctx context.Context, db *pgx.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a %s request for : %v", r.Method, r.URL.Path)

		rows, err := db.Query(ctx, `SELECT
			workspace_id, account_id,
			slug, name,
			created_at, updated_at
			FROM workspace LIMIT 100
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
				&workspace.Slug, &workspace.Name,
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

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var err error

	pgConnStr, pgConnStatus := os.LookupEnv("POSTGRES_URI")
	if !pgConnStatus {
		return fmt.Errorf("env `POSTGRES_URI` is not set")
	}

	db, err := pgx.Connect(ctx, pgConnStr)

	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	defer db.Close(ctx)

	var tm time.Time

	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)

	if err != nil {
		return fmt.Errorf("failed to query database: %v", err)
	}

	log.Printf("database ready with db time: %s\n", tm.Format(time.RFC1123))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handleGetIndex)
	mux.Handle("GET /workspaces/{$}", handleGetWorkspaces(ctx, db))

	srv := &http.Server{
		Addr:              *addr,
		Handler:           mux,
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
