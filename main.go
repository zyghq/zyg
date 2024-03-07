package main

import (
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

func run() error {

	var err error

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleGetRoot)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

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

	log.Println("listening...")
	err = srv.ListenAndServe()

	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
