package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	//
	// from: https://github.com/golang/go/issues/4799
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	log.Printf("Got a %s request for : %v", r.Method, r.URL.Path)
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleSomething(ctx context.Context, db *sql.DB) http.Handler {
	fmt.Println("do something here...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var workspaceId string
		row := db.QueryRowContext(context.TODO(), "SELECT workspace_id FROM something WHERE id = $1", "57a76972641b4dec92b252e6fc28af97")
		log.Printf("Got a %s request for : %v", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("*** from route ***"))
	})
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var err error

	pgConnStr, pgConnStatus := os.LookupEnv("POSTGRES_URI")
	if !pgConnStatus {
		return fmt.Errorf("env `POSTGRES_URI` not found")
	}

	db, err := sql.Open("pgx", pgConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("database cannot be reached %v", err)

	}

	rows, err := db.Query("SELECT NOW()")
	if err != nil {
		return fmt.Errorf("failed to query database %v", err)
	}
	defer rows.Close()

	var dbTime string
	for rows.Next() {
		err := rows.Scan(&dbTime)
		if err != nil {
			return fmt.Errorf("failed to parse database results %v", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleGetRoot)
	mux.Handle("/something/", handleSomething(ctx, db))

	srv := &http.Server{
		Addr:              ":8080", // probaly read from env
		Handler:           mux,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

	log.Println("listening...")
	err = srv.ListenAndServe()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for {
			if <-c == os.Interrupt {
				_ = srv.Close()
				return
			}
		}
	}()

	return err
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
