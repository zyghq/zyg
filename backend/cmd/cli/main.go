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
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/services"
)

var workspaceID string
var customerID string
var memberID string
var threadID string

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
	threadStore := repository.NewThreadDB(conn.DB)

	// Initialize services
	app := &AppServices{
		AuthService:      services.NewAuthService(accountStore, memberStore),
		AccountService:   services.NewAccountService(accountStore, workspaceStore),
		WorkspaceService: services.NewWorkspaceService(workspaceStore, memberStore, customerStore),
		CustomerService:  services.NewCustomerService(customerStore),
		ThreadService:    services.NewThreadService(threadStore),
		SyncService:      services.NewSyncService(),
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

func runSyncWorkspace(cmd *cobra.Command, _ []string) error {
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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
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
	synced, err := app.SyncService.SyncWorkspaceRPC(ctx, shape)
	if err != nil {
		return fmt.Errorf("failed to sync workspace %s: %w", workspaceID, err)
	}
	log.Info().Msgf("Successfully synced workspace with ID: %s, versionID: %s", synced.WorkspaceID, synced.VersionID)
	log.Info().Msg("Workspace sync command completed successfully.")
	return nil
}

func runSyncWorkspaceCustomers(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log.Info().Msgf("Starting sync customers WorkspaceID: %s...", workspaceID)

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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
	}

	customers, err := app.WorkspaceService.ListCustomers(ctx, workspace.WorkspaceId)
	if err != nil {
		return fmt.Errorf("failed to get customers for workspaceID %s: %w", workspaceID, err)
	}
	shapes := make([]models.CustomerShape, 0, len(customers))
	var recentSyncedShape *models.CustomerShape
	var syncCount int

	for _, customer := range customers {
		var externalId *string
		var email *string
		var phone *string
		if customer.ExternalId.Valid {
			externalId = &customer.ExternalId.String
		}
		if customer.Email.Valid {
			email = &customer.Email.String
		}
		if customer.Phone.Valid {
			phone = &customer.Phone.String
		}
		u, _ := uuid.NewUUID()
		shape := models.CustomerShape{
			CustomerID:      customer.CustomerId,
			WorkspaceID:     customer.WorkspaceId,
			ExternalID:      externalId,
			Email:           email,
			Phone:           phone,
			Name:            customer.Name,
			Role:            customer.Role,
			AvatarURL:       customer.AvatarUrl(),
			IsEmailVerified: customer.IsEmailVerified,
			CreatedAt:       customer.CreatedAt,
			UpdatedAt:       customer.UpdatedAt,
			SyncedAt:        time.Now().UTC(),
			VersionID:       u.String(),
		}
		shapes = append(shapes, shape)
	}

	for _, shape := range shapes {
		log.Info().Msgf("Syncing customer with ID: %s", shape.CustomerID)
		synced, err := app.SyncService.SyncCustomer(ctx, shape)
		if err != nil {
			return fmt.Errorf("failed to sync customer %s: %w", shape.CustomerID, err)
		}
		log.Info().Msgf(
			"Successfully synced customerID: %s, versionID: %s", synced.CustomerID, synced.VersionID)
		recentSyncedShape = &shape
		syncCount++
	}

	log.Info().Msgf("Synced total %d of %d customers", syncCount, len(shapes))
	if recentSyncedShape != nil {
		log.Info().Msgf("Last synced customer: %s", recentSyncedShape.CustomerID)
	}
	return nil
}

func runSyncWorkspaceMembers(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log.Info().Msgf("Starting sync members WorkspaceID: %s...", workspaceID)

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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
	}

	members, err := app.WorkspaceService.ListMembers(ctx, workspace.WorkspaceId)
	if err != nil {
		return fmt.Errorf("failed to get members for workspaceID %s: %w", workspaceID, err)
	}
	shapes := make([]models.MemberShape, 0, len(members))
	var recentSyncedShape *models.MemberShape
	var syncCount int

	for _, member := range members {
		u, _ := uuid.NewUUID()
		permissions := make(map[string]interface{})
		shape := models.MemberShape{
			MemberID:    member.MemberId,
			WorkspaceID: member.WorkspaceId,
			Name:        member.Name,
			PublicName:  member.Name,
			Role:        member.Role,
			Permissions: permissions,
			AvatarURL:   member.AvatarUrl(),
			CreatedAt:   member.CreatedAt,
			UpdatedAt:   member.UpdatedAt,
			SyncedAt:    time.Now().UTC(),
			VersionID:   u.String(),
		}
		shapes = append(shapes, shape)
	}

	for _, shape := range shapes {
		log.Info().Msgf("Syncing member with ID: %s", shape.MemberID)
		synced, err := app.SyncService.SyncMember(ctx, shape)
		if err != nil {
			return fmt.Errorf("failed to sync customer %s: %w", shape.MemberID, err)
		}
		log.Info().Msgf(
			"Successfully synced memberID: %s, versionID: %s", synced.MemberID, synced.VersionID)
		recentSyncedShape = &shape
		syncCount++
	}

	log.Info().Msgf("Synced total %d of %d customers", syncCount, len(shapes))
	if recentSyncedShape != nil {
		log.Info().Msgf("Last synced memberID: %s", recentSyncedShape.MemberID)
	}
	return nil
}

func runSyncWorkspaceCustomer(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log.Info().Msgf("Starting sync customer WorkspaceID: %s, CustomerID: %s...", workspaceID, customerID)

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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
	}
	customer, err := app.WorkspaceService.GetCustomer(ctx, workspace.WorkspaceId, customerID, nil)
	if err != nil {
		return fmt.Errorf("failed to get customer %s: %w", customerID, err)
	}

	var externalId *string
	var email *string
	var phone *string

	if customer.ExternalId.Valid {
		externalId = &customer.ExternalId.String
	}
	if customer.Email.Valid {
		email = &customer.Email.String
	}
	if customer.Phone.Valid {
		phone = &customer.Phone.String
	}

	u, _ := uuid.NewUUID()
	shape := models.CustomerShape{
		CustomerID:      customer.CustomerId,
		WorkspaceID:     customer.WorkspaceId,
		ExternalID:      externalId,
		Email:           email,
		Phone:           phone,
		Name:            customer.Name,
		Role:            customer.Role,
		AvatarURL:       customer.AvatarUrl(),
		IsEmailVerified: customer.IsEmailVerified,
		CreatedAt:       customer.CreatedAt,
		UpdatedAt:       customer.UpdatedAt,
		SyncedAt:        time.Now().UTC(),
		VersionID:       u.String(),
	}
	log.Info().Msgf("Syncing customer with ID: %s", customerID)
	synced, err := app.SyncService.SyncCustomer(ctx, shape)
	if err != nil {
		return fmt.Errorf("failed to sync customer %s: %w", customerID, err)
	}
	log.Info().Msgf(
		"Successfully synced customerID: %s, versionID: %s", synced.CustomerID, synced.VersionID)
	log.Info().Msg("Customer sync command completed successfully.")
	return nil
}

func runSyncWorkspaceMember(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log.Info().Msgf("Starting sync member WorkspaceID: %s, MemberID: %s...", workspaceID, memberID)

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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
	}
	member, err := app.WorkspaceService.GetMember(ctx, workspace.WorkspaceId, memberID)
	if err != nil {
		return fmt.Errorf("failed to get member %s: %w", memberID, err)
	}

	u, _ := uuid.NewUUID()
	permissions := make(map[string]interface{})
	shape := models.MemberShape{
		MemberID:    member.MemberId,
		WorkspaceID: member.WorkspaceId,
		Name:        member.Name,
		PublicName:  member.Name,
		Role:        member.Role,
		Permissions: permissions,
		AvatarURL:   member.AvatarUrl(),
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
		SyncedAt:    time.Now().UTC(),
		VersionID:   u.String(),
	}

	log.Info().Msgf("Syncing member with ID: %s", memberID)
	synced, err := app.SyncService.SyncMember(ctx, shape)
	if err != nil {
		return fmt.Errorf("failed to sync member %s: %w", memberID, err)
	}
	log.Info().Msgf(
		"Successfully synced memberID: %s, versionID: %s", synced.MemberID, synced.VersionID)
	log.Info().Msg("Member sync command completed successfully.")
	return nil
}

func runSyncWorkspaceThread(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log.Info().Msgf("Starting sync thread WorkspaceID: %s, threadID %s...", workspaceID, threadID)

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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
	}

	thread, err := app.ThreadService.GetWorkspaceThread(ctx, workspace.WorkspaceId, threadID, nil)
	if err != nil {
		return fmt.Errorf("failed to get thread %s: %w", threadID, err)
	}

	labels, err := app.ThreadService.ListThreadLabels(ctx, thread.ThreadId)
	if err != nil {
		return fmt.Errorf("failed to get labels for threadID %s: %w", threadID, err)
	}

	type ThreadLabel struct {
		LabelId   string    `json:"labelId"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	threadLabels := make([]ThreadLabel, 0, len(labels))
	for _, label := range labels {
		threadLabels = append(threadLabels, ThreadLabel{
			LabelId:   label.LabelId,
			Name:      label.Name,
			CreatedAt: label.CreatedAt,
			UpdatedAt: label.UpdatedAt,
		})
	}

	labelsMap := make(map[string]interface{})
	for _, label := range threadLabels {
		labelsMap[label.LabelId] = map[string]interface{}{
			"labelId":   label.LabelId,
			"name":      label.Name,
			"createdAt": label.CreatedAt,
			"updatedAt": label.UpdatedAt,
		}
	}

	var (
		assigneeId, inboundSeqId, outboundSeqId *string
	)

	var assignedAt *time.Time
	var previewText string

	if thread.AssignedMember != nil {
		assigneeId = &thread.AssignedMember.MemberId
		assignedAt = &thread.AssignedMember.AssignedAt
	}
	if thread.InboundMessage != nil {
		inboundSeqId = &thread.InboundMessage.LastSeqId
	}
	if thread.OutboundMessage != nil {
		outboundSeqId = &thread.OutboundMessage.LastSeqId
	}
	if thread.InboundMessage != nil {
		previewText = thread.InboundMessage.PreviewText
	} else if thread.OutboundMessage != nil {
		previewText = thread.OutboundMessage.PreviewText
	}

	u, _ := uuid.NewUUID()
	shape := models.ThreadShape{
		ThreadID:          thread.ThreadId,
		WorkspaceID:       thread.WorkspaceId,
		CustomerID:        thread.Customer.CustomerId,
		AssigneeID:        assigneeId,
		AssignedAt:        assignedAt,
		Title:             thread.Title,
		Description:       thread.Description,
		PreviewText:       previewText,
		Status:            thread.ThreadStatus.Status,
		StatusChangedAt:   thread.ThreadStatus.StatusChangedAt,
		StatusChangedByID: thread.ThreadStatus.StatusChangedBy.MemberId,
		Stage:             thread.ThreadStatus.Stage,
		Replied:           thread.Replied,
		Priority:          thread.Priority,
		Channel:           thread.Channel,
		CreatedByID:       thread.CreatedBy.MemberId,
		UpdatedByID:       thread.UpdatedBy.MemberId,
		Labels:            labelsMap,
		InboundSeqID:      inboundSeqId,
		OutboundSeqID:     outboundSeqId,
		CreatedAt:         thread.CreatedAt,
		UpdatedAt:         thread.UpdatedAt,
		SyncedAt:          time.Now().UTC(),
		VersionID:         u.String(),
	}

	fmt.Println(shape.Labels)

	log.Info().Msgf("Syncing thread with ID: %s", threadID)
	synced, err := app.SyncService.SyncThread(ctx, shape)
	if err != nil {
		return fmt.Errorf("failed to sync thread %s: %w", threadID, err)
	}
	log.Info().Msgf(
		"Successfully synced threadID: %s, versionID: %s", synced.ThreadID, synced.VersionID)
	log.Info().Msg("Thread sync command completed successfully.")
	return nil
}

func runSyncWorkspaceThreads(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	log.Info().Msgf("Starting sync thread WorkspaceID: %s, threadID %s...", workspaceID, threadID)

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
	app := initServices(conn)

	log.Info().Msgf("Fetching workspace with ID: %s", workspaceID)
	workspace, err := app.WorkspaceService.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
	}

	threads, err := app.ThreadService.ListWorkspaceThreads(ctx, workspace.WorkspaceId)
	if err != nil {
		return fmt.Errorf("failed to get threads for workspaceID %s: %w", workspaceID, err)
	}

	type ThreadLabel struct {
		LabelId   string    `json:"labelId"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	var recentSyncedThread *models.ThreadShape
	var syncCount int

	for _, thread := range threads {
		labels, err := app.ThreadService.ListThreadLabels(ctx, thread.ThreadId)
		if err != nil {
			return fmt.Errorf("failed to get labels for threadID %s: %w", threadID, err)
		}

		threadLabels := make([]ThreadLabel, 0, len(labels))
		for _, label := range labels {
			threadLabels = append(threadLabels, ThreadLabel{
				LabelId:   label.LabelId,
				Name:      label.Name,
				CreatedAt: label.CreatedAt,
				UpdatedAt: label.UpdatedAt,
			})
		}

		labelsMap := make(map[string]interface{})
		for _, label := range threadLabels {
			labelsMap[label.LabelId] = map[string]interface{}{
				"labelId":   label.LabelId,
				"name":      label.Name,
				"createdAt": label.CreatedAt,
				"updatedAt": label.UpdatedAt,
			}
		}

		var (
			assigneeId, inboundSeqId, outboundSeqId *string
		)

		var assignedAt *time.Time
		var previewText string

		if thread.AssignedMember != nil {
			assigneeId = &thread.AssignedMember.MemberId
			assignedAt = &thread.AssignedMember.AssignedAt
		}
		if thread.InboundMessage != nil {
			inboundSeqId = &thread.InboundMessage.LastSeqId
		}
		if thread.OutboundMessage != nil {
			outboundSeqId = &thread.OutboundMessage.LastSeqId
		}
		if thread.InboundMessage != nil {
			previewText = thread.InboundMessage.PreviewText
		} else if thread.OutboundMessage != nil {
			previewText = thread.OutboundMessage.PreviewText
		}

		u, _ := uuid.NewUUID()
		shape := models.ThreadShape{
			ThreadID:          thread.ThreadId,
			WorkspaceID:       thread.WorkspaceId,
			CustomerID:        thread.Customer.CustomerId,
			AssigneeID:        assigneeId,
			AssignedAt:        assignedAt,
			Title:             thread.Title,
			Description:       thread.Description,
			PreviewText:       previewText,
			Status:            thread.ThreadStatus.Status,
			StatusChangedAt:   thread.ThreadStatus.StatusChangedAt,
			StatusChangedByID: thread.ThreadStatus.StatusChangedBy.MemberId,
			Stage:             thread.ThreadStatus.Stage,
			Replied:           thread.Replied,
			Priority:          thread.Priority,
			Channel:           thread.Channel,
			CreatedByID:       thread.CreatedBy.MemberId,
			UpdatedByID:       thread.UpdatedBy.MemberId,
			Labels:            labelsMap,
			InboundSeqID:      inboundSeqId,
			OutboundSeqID:     outboundSeqId,
			CreatedAt:         thread.CreatedAt,
			UpdatedAt:         thread.UpdatedAt,
			SyncedAt:          time.Now().UTC(),
			VersionID:         u.String(),
		}

		log.Info().Msgf("Syncing thread with ID: %s", threadID)
		synced, err := app.SyncService.SyncThread(ctx, shape)
		if err != nil {
			return fmt.Errorf("failed to sync thread %s: %w", threadID, err)
		}
		log.Info().Msgf(
			"Successfully synced threadID: %s, versionID: %s", synced.ThreadID, synced.VersionID)
		recentSyncedThread = &shape
		syncCount++
	}
	log.Info().Msgf("Synced total %d of %d threads", syncCount, len(threads))
	if recentSyncedThread != nil {
		log.Info().Msgf("Last synced threadID: %s", recentSyncedThread.ThreadID)
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI for Zyg platform operations",
	Long:  `A command line interface for managing Zyg workspaces and performing sync operations.`,
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync operations",
}

var syncWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Sync workspace related data",
}

var workspaceSubCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Sync workspace data",
	Long:  `Sync a workspace using its workspace ID.`,
	RunE:  runSyncWorkspace,
}

var customersSubCmd = &cobra.Command{
	Use:   "customers",
	Short: "Sync workspace customers",
	RunE:  runSyncWorkspaceCustomers,
}

var membersSubCmd = &cobra.Command{
	Use:   "members",
	Short: "Sync workspace members",
	RunE:  runSyncWorkspaceMembers,
}

var customerSubCmd = &cobra.Command{
	Use:   "customer",
	Short: "Sync specific customer",
	RunE:  runSyncWorkspaceCustomer,
}

var memberSubCmd = &cobra.Command{
	Use:   "member",
	Short: "Sync specific member",
	RunE:  runSyncWorkspaceMember,
}

var threadSubCmd = &cobra.Command{
	Use:   "thread",
	Short: "Sync specific thread",
	RunE:  runSyncWorkspaceThread,
}

var threadsSubCmd = &cobra.Command{
	Use:   "threads",
	Short: "Sync workspace threads",
	RunE:  runSyncWorkspaceThreads,
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(syncWorkspaceCmd)

	syncWorkspaceCmd.AddCommand(workspaceSubCmd)

	syncWorkspaceCmd.AddCommand(customerSubCmd)
	syncWorkspaceCmd.AddCommand(customersSubCmd)

	syncWorkspaceCmd.AddCommand(memberSubCmd)
	syncWorkspaceCmd.AddCommand(membersSubCmd)

	syncWorkspaceCmd.AddCommand(threadSubCmd)
	syncWorkspaceCmd.AddCommand(threadsSubCmd)

	// Add the workspaceID flag
	syncWorkspaceCmd.PersistentFlags().StringVar(
		&workspaceID, "workspaceID", "", "Workspace ID (required)")
	if err := syncWorkspaceCmd.MarkPersistentFlagRequired("workspaceID"); err != nil {
		log.Fatal().Err(err).Msg("failed to mark workspace id flag as required")
	}

	// Add the customerID flag
	customerSubCmd.Flags().StringVar(&customerID, "customerID", "", "Customer ID (required)")
	if err := customerSubCmd.MarkFlagRequired("customerID"); err != nil {
		log.Fatal().Err(err).Msg("failed to mark customer id flag as required")
	}

	memberSubCmd.Flags().StringVar(&memberID, "memberID", "", "Member ID (required)")
	if err := memberSubCmd.MarkFlagRequired("memberID"); err != nil {
		log.Fatal().Err(err).Msg("failed to mark member id flag as required")
	}

	threadSubCmd.Flags().StringVar(&threadID, "threadID", "", "Thread ID (required)")
	if err := threadSubCmd.MarkFlagRequired("threadID"); err != nil {
		log.Fatal().Err(err).Msg("failed to mark thread id flag as required")
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
