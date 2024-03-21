package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

// some variations of the below
// https://www.alexedwards.net/blog/organising-database-access#:~:text=func%20booksIndex(env%20*Env)%20http.HandlerFunc%20%7B
func handleGetWorkspaces(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a %s request for : %v", r.Method, r.URL.Path)

		rows, err := db.Query(`
			SELECT 
			workspace_id, account_id, 
			slug, name, 
			created_at, updated_at
			FROM workspace LIMIT 100
		`)
		if err != nil {
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
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			workspaces = append(workspaces, workspace)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(workspaces); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func run() error {
	var err error

	pgConnStr, pgConnStatus := os.LookupEnv("POSTGRES_URI")
	if !pgConnStatus {
		return fmt.Errorf("env `POSTGRES_URI` is not set")
	}

	db, err := sql.Open("pgx", pgConnStr)

	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	rows, err := db.Query("SELECT NOW()")
	if err != nil {
		return fmt.Errorf("failed to query database: %v", err)
	}

	defer rows.Close()

	var dbTime string
	for rows.Next() {
		err = rows.Scan(&dbTime)
		if err != nil {
			return fmt.Errorf("failed to scan database: %v", err)
		}
	}

	log.Printf("database ready with date time: %s", dbTime)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handleGetIndex)
	mux.Handle("GET /workspaces/{$}", handleGetWorkspaces(db))

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
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
