package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/handler"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/services"
)

var host = flag.String("host", "0.0.0.0", "host")
var port = flag.String("port", "8080", "port")

var addr string

func run(ctx context.Context) error {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// get postgres connection string from env
	pgConnStr, err := zyg.GetEnv("DATABASE_URL")
	if err != nil {
		return fmt.Errorf("failed to get DATABASE_URL env got error: %v", err)
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

	slog.Info("database", slog.Any("db time", tm.Format(time.RFC1123)))

	rdb := redis.NewClient(&redis.Options{
		Addr:     zyg.RedisAddr(),
		Password: zyg.RedisPassword(),
		DB:       0,
	})

	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			slog.Error("failed to close redis client", slog.Any("err", err))
		}
	}(rdb)

	// Perform basic diagnostic to check if the connection is working
	// Expected result > ping: PONG
	// If Redis is not running, error case is taken instead
	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		// return fmt.Errorf("failed to ping redis got error: %v", err)
		slog.Error("failed to connect redis - JUST LOGGING FOR NOW!", slog.Any("err", err))
	}
	slog.Info("redis", slog.Any("redis status", status))

	// init stores
	accountStore := repository.NewAccountDB(db)
	workspaceStore := repository.NewWorkspaceDB(db)
	memberStore := repository.NewMemberDB(db)
	customerStore := repository.NewCustomerDB(db)
	threadStore := repository.NewThreadDB(db, rdb)

	// init services
	authService := services.NewAuthService(accountStore, memberStore)
	accountService := services.NewAccountService(accountStore, workspaceStore)
	workspaceService := services.NewWorkspaceService(workspaceStore, memberStore, customerStore)
	customerService := services.NewCustomerService(customerStore)
	threadService := services.NewThreadService(threadStore)

	// init server
	srv := handler.NewServer(
		authService,
		accountService,
		workspaceService,
		customerService,
		threadService,
	)

	addr = fmt.Sprintf("%s:%s", *host, *port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           srv,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	slog.Info("server up and running", slog.String("addr", addr))

	err = httpServer.ListenAndServe()
	return err
}

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%s\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
