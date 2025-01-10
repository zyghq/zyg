package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/zyghq/zyg/models"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/esync"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/services"
)

var workspaceID string

// AppServices holds all the application services
type AppServices struct {
	AuthService      *services.AuthService
	AccountService   *services.AccountService
	WorkspaceService *services.WorkspaceService
	CustomerService  *services.CustomerService
	ThreadService    *services.ThreadService
	SyncService      *services.SyncService
}

// AppConnections holds database and redis connections
type AppConnections struct {
	DB     *pgxpool.Pool
	SyncDB *pgxpool.Pool
	Redis  *redis.Client
}

func initConnections(ctx context.Context) (*AppConnections, error) {
	log.Info().Msg("Initializing database and redis connections...")

	// Initialize main database
	pgConnStr, err := zyg.GetEnv("DATABASE_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to get DATABASE_URL env: %w", err)
	}
	log.Debug().Msg("Connecting to main database...")
	db, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create app pg connection pool: %w", err)
	}
	log.Debug().Msg("Successfully connected to main database")

	// Initialize sync database
	syncPGConnStr, err := zyg.GetEnv("SYNC_DATABASE_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to get SYNC_DATABASE_URL env: %w", err)
	}
	log.Debug().Msg("Connecting to sync database...")
	syncDB, err := pgxpool.New(ctx, syncPGConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync pg connection pool: %w", err)
	}
	log.Debug().Msg("Successfully connected to sync database")

	// Verify database connections
	var tm time.Time
	if err := db.QueryRow(ctx, "SELECT NOW()").Scan(&tm); err != nil {
		return nil, fmt.Errorf("failed db query got error: %w", err)
	}
	log.Info().Msgf("app database time: %s", tm.Format(time.RFC1123))

	if err := syncDB.QueryRow(ctx, "SELECT NOW()").Scan(&tm); err != nil {
		return nil, fmt.Errorf("failed db query got error: %w", err)
	}
	log.Info().Msgf("sync database time: %s", tm.Format(time.RFC1123))

	// Redis options
	opts := &redis.Options{
		Addr:     zyg.RedisAddr(),
		Username: zyg.RedisUsername(),
		Password: zyg.RedisPassword(),
		DB:       0,
	}

	if zyg.RedisTLSEnabled() {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	log.Debug().Msg("Connecting to redis...")
	rdb := redis.NewClient(opts)

	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis got error: %v", err)
	}
	log.Info().Msgf("redis status PING: %s", status)

	// Verify Redis connection
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis check failed: %w", err)
	}
	log.Debug().Msg("Successfully connected to redis")

	log.Info().Msg("Database and redis connections initialized successfully.")
	return &AppConnections{
		DB:     db,
		SyncDB: syncDB,
		Redis:  rdb,
	}, nil
}

// initServices initializes all application services
func initServices(conn *AppConnections) *AppServices {
	log.Info().Msg("Initializing application services...")
	// Initialize application stores
	accountStore := repository.NewAccountDB(conn.DB)
	workspaceStore := repository.NewWorkspaceDB(conn.DB)
	memberStore := repository.NewMemberDB(conn.DB)
	customerStore := repository.NewCustomerDB(conn.DB)
	threadStore := repository.NewThreadDB(conn.DB, conn.Redis)

	// Initialize sync store
	syncStore := esync.NewSyncDB(conn.SyncDB)

	// Initialize services
	app := &AppServices{
		AuthService:      services.NewAuthService(accountStore, memberStore),
		AccountService:   services.NewAccountService(accountStore, workspaceStore),
		WorkspaceService: services.NewWorkspaceService(workspaceStore, memberStore, customerStore),
		CustomerService:  services.NewCustomerService(customerStore),
		ThreadService:    services.NewThreadService(threadStore),
		SyncService:      services.NewSyncService(syncStore),
	}
	log.Info().Msg("Application services initialized successfully.")
	return app
}

// initSentry initializes the Sentry client
func initSentry() error {
	log.Info().Msg("Initializing Sentry...")
	err := sentry.Init(sentry.ClientOptions{
		Debug:         zyg.SentryDebugEnabled(),
		EnableTracing: true,
		Environment:   zyg.SentryEnv(),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize Sentry")
		return err
	}
	log.Info().Msg("Sentry initialized successfully.")
	return nil
}

// cleanup handles graceful shutdown of connections
func cleanup(conn *AppConnections) {
	log.Info().Msg("Cleaning up connections...")
	if conn.DB != nil {
		log.Debug().Msg("Closing main database connection...")
		conn.DB.Close()
		log.Debug().Msg("Main database connection closed.")
	}
	if conn.SyncDB != nil {
		log.Debug().Msg("Closing sync database connection...")
		conn.SyncDB.Close()
		log.Debug().Msg("Sync database connection closed.")
	}
	if conn.Redis != nil {
		log.Debug().Msg("Closing redis connection...")
		if err := conn.Redis.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close redis")
		}
		log.Debug().Msg("Redis connection closed.")
	}
	sentry.Flush(2 * time.Second)
	log.Info().Msg("Cleanup complete.")
}

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI for Zyg platform operations",
	Long:  `A command line interface for managing Zyg workspaces and performing sync operations.`,
}

var syncCmd = &cobra.Command{
	Use:   "sync workspace",
	Short: "Sync a workspace",
	Long:  `Sync a workspace using its workspace ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		log.Info().Msg("Starting workspace sync command...")

		// Initialize Sentry
		if err := initSentry(); err != nil {
			return fmt.Errorf("failed to initialize sentry: %w", err)
		}

		// Initialize connections
		conn, err := initConnections(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize connections: %w", err)
		}
		defer cleanup(conn)

		// Initialize services
		services := initServices(conn)

		log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
		workspace, err := services.WorkspaceService.GetWorkspace(ctx, workspaceID)
		if err != nil {
			return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
		}

		u, _ := uuid.NewUUID()
		shape := models.WorkspaceShape{
			WorkspaceID: workspace.WorkspaceId,
			Name:        workspace.Name,
			PublicName:  workspace.Name,
			CreatedAt:   workspace.CreatedAt,
			UpdatedAt:   workspace.UpdatedAt,
			SyncedAt:    time.Now().UTC(),
			VersionID:   u.String(),
		}

		log.Info().Msgf("Syncing workspace with ID: %s", workspaceID)
		synced, err := services.SyncService.SyncWorkspace(ctx, shape)
		if err != nil {
			return fmt.Errorf("failed to sync workspace %s: %w", workspaceID, err)
		}
		log.Info().Msgf("Successfully synced workspace with ID: %s, versionID: %s", synced.WorkspaceID, synced.VersionID)
		log.Info().Msg("Workspace sync command completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVar(&workspaceID, "id", "", "Workspace ID (required)")
	err := syncCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to mark id flag as required")
		return
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to execute command")
		os.Exit(1)
	}
}
