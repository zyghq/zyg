package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/api"
)

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

func run(ctx context.Context) error {
	var err error
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	pgConnStr, err := zyg.GetEnv("POSTGRES_URI")
	if err != nil {
		return fmt.Errorf("failed to get POSTGRES_URI env got error: %v", err)
	}

	db, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		return fmt.Errorf("unable to create pg connection pool: %v", err)
	}

	defer db.Close()

	var tm time.Time
	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)

	if err != nil {
		return fmt.Errorf("db query failed got error: %v", err)
	}

	slog.Info("database ready", slog.Any("dbtime", tm.Format(time.RFC1123)))

	hanldler := api.NewHandler(ctx, db)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           hanldler,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

	slog.Info("server up and running", slog.String("addr", *addr))

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
