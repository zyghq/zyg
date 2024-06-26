// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgxpool"

// 	"github.com/zyghq/zyg"
// 	"github.com/zyghq/zyg/internal/api"
// )

// var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

// func run(ctx context.Context) error {
// 	var err error
// 	ctx, cancel := context.WithCancel(ctx)

// 	defer cancel()

// 	pgConnStr, err := zyg.GetEnv("POSTGRES_URI")
// 	if err != nil {
// 		return fmt.Errorf("failed to get POSTGRES_URI env got error: %v", err)
// 	}

// 	db, err := pgxpool.New(ctx, pgConnStr)
// 	if err != nil {
// 		return fmt.Errorf("unable to create pg connection pool: %v", err)
// 	}

// 	defer db.Close()

// 	var tm time.Time
// 	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)

// 	if err != nil {
// 		return fmt.Errorf("db query failed got error: %v", err)
// 	}

// 	slog.Info("database ready", slog.Any("dbtime", tm.Format(time.RFC1123)))

// 	hanldler := api.NewHandler(ctx, db)
// 	srv := &http.Server{
// 		Addr:              *addr,
// 		Handler:           hanldler,
// 		ReadTimeout:       30 * time.Second,
// 		WriteTimeout:      90 * time.Second,
// 		IdleTimeout:       time.Minute,
// 		ReadHeaderTimeout: 30 * time.Second,
// 	}

// 	slog.Info("server up and running", slog.String("addr", *addr))

// 	err = srv.ListenAndServe()
// 	return err
// }

// func main() {
// 	flag.Parse()
// 	ctx := context.Background()
// 	if err := run(ctx); err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		os.Exit(1)
// 	}
// }

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
	"github.com/zyghq/zyg/internal/adapters/repository"
	"github.com/zyghq/zyg/internal/adapters/xhandler"
	"github.com/zyghq/zyg/internal/services"
)

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

func run(ctx context.Context) error {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// get POSTGRES_URI from env
	pgConnStr, err := zyg.GetEnv("POSTGRES_URI")
	if err != nil {
		return fmt.Errorf("failed to get POSTGRES_URI env got error: %v", err)
	}

	// create pg connection pool
	db, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		return fmt.Errorf("unable to create pg connection pool: %v", err)
	}

	defer db.Close()

	// make sure db is up and running
	var tm time.Time
	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)
	if err != nil {
		return fmt.Errorf("db query failed got error: %v", err)
	}

	slog.Info("database", slog.Any("dbtime", tm.Format(time.RFC1123)))

	// init respective stores
	workspaceStore := repository.NewWorkspaceDB(db)
	memberStore := repository.NewMemberDB(db)
	customerStore := repository.NewCustomerDB(db)
	threadChatStore := repository.NewThreadChatDB(db)

	// init respective services
	workspaceService := services.NewWorkspaceService(workspaceStore, memberStore, customerStore)
	customerService := services.NewCustomerService(customerStore)
	threadChatService := services.NewThreadChatService(threadChatStore)

	// init server
	srv := xhandler.NewServer(
		workspaceService,
		customerService,
		threadChatService,
	)

	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	slog.Info("server up and running", slog.String("addr", *addr))

	err = httpServer.ListenAndServe()
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
