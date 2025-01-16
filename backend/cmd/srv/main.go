package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/zyghq/zyg/adapters/esync"
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
		return fmt.Errorf("unable to create app pg connection pool: %v", err)
	}
	defer db.Close()

	syncPGConnStr, err := zyg.GetEnv("SYNC_DATABASE_URL")
	if err != nil {
		return fmt.Errorf("failed to get SYNC_DATABASE_URL env got error: %v", err)
	}

	syncDB, err := pgxpool.New(ctx, syncPGConnStr)
	if err != nil {
		return fmt.Errorf("unable to create sync pg connection pool: %v", err)
	}
	defer syncDB.Close()

	// make sure db is up and running
	var tm time.Time
	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)
	if err != nil {
		return fmt.Errorf("db query failed got error: %v", err)
	}
	slog.Info("app database", slog.String("time", tm.Format(time.RFC1123)))

	err = syncDB.QueryRow(ctx, "SELECT NOW()").Scan(&tm)
	if err != nil {
		return fmt.Errorf("db query failed got error: %v", err)
	}
	slog.Info("sync database", slog.String("time", tm.Format(time.RFC1123)))

	// Redis options
	//opts := &redis.Options{
	//	Addr:     zyg.RedisAddr(),
	//	Username: zyg.RedisUsername(),
	//	Password: zyg.RedisPassword(),
	//	DB:       0,
	//}
	//
	//if zyg.RedisTLSEnabled() {
	//	opts.TLSConfig = &tls.Config{
	//		InsecureSkipVerify: true,
	//	}
	//}
	//
	//rdb := redis.NewClient(opts)
	//
	//defer func(rdb *redis.Client) {
	//	err := rdb.Close()
	//	if err != nil {
	//		slog.Error("failed to close redis client", slog.Any("err", err))
	//	}
	//}(rdb)

	// Perform basic diagnostic to check if the connection is working
	// Expected result > ping: PONG
	// If Redis is not running, error case is taken instead
	//status, err := rdb.Ping(ctx).Result()
	//if err != nil {
	//	return fmt.Errorf("failed to ping redis got error: %v", err)
	//}
	//slog.Info("redis", slog.Any("status", status))

	// setup sentry
	if err := sentry.Init(sentry.ClientOptions{
		Debug:         zyg.SentryDebugEnabled(),
		EnableTracing: true,
		Environment:   zyg.SentryEnv(),
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			// Don't sample Index
			if ctx.Span.Name == "GET /" {
				return 0.0
			}
			return 1.0
		},
	}); err != nil {
		slog.Error(
			"sentry init failed logging the error and continue...",
			slog.Any("err", err))
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	// Initialize application DB stores
	// Not to be confused with sync DB store.
	accountStore := repository.NewAccountDB(db)
	workspaceStore := repository.NewWorkspaceDB(db)
	memberStore := repository.NewMemberDB(db)
	customerStore := repository.NewCustomerDB(db)
	threadStore := repository.NewThreadDB(db)

	// Initialize sync DB store.
	// Not to be confused with application DB store.
	syncStore := esync.NewSyncDB(syncDB)

	// Initialize application services with application DB stores.
	authService := services.NewAuthService(accountStore, memberStore)
	accountService := services.NewAccountService(accountStore, workspaceStore)
	workspaceService := services.NewWorkspaceService(workspaceStore, memberStore, customerStore)
	customerService := services.NewCustomerService(customerStore)
	threadService := services.NewThreadService(threadStore)

	// Initialize sync service with sync DB store.
	syncService := services.NewSyncService(syncStore)

	// init server
	srv := handler.NewServer(
		authService,
		accountService,
		workspaceService,
		customerService,
		threadService,
		syncService,
	)

	// wrap sentry
	sentryHandler := sentryhttp.New(sentryhttp.Options{}).Handle(srv)

	addr = fmt.Sprintf("%s:%s", *host, *port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           sentryHandler,
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
