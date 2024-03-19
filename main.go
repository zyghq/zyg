package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	//
	// from: https://github.com/golang/go/issues/4799
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("Got a %s request for : %v", r.Method, r.URL.Path)
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func run() error {
	var err error

	pgConnStr, pgConnStatus := os.LookupEnv("POSTGRES_URI")
	if !pgConnStatus {
		return fmt.Errorf("env POSTGRES_URI not set")
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

	log.Printf("datbase ready with datetime: %s", dbTime)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleGetIndex)

	srv := &http.Server{
		Addr:              *addr,
		Handler:           mux,
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
