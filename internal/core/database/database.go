package database

import (
	"context"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDB(lc fx.Lifecycle, logger *zap.Logger) (*pgxpool.Pool, error) {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://ticketer:password@localhost:5432/ticketer?sslmode=disable"
	}

	// 1. Parse the standard pgx config
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// 2. Check for GCP Cloud SQL environment variable and attach the secure dialer
	instanceConnName := os.Getenv("CLOUD_SQL_INSTANCE_NAME")
	if instanceConnName != "" {
		logger.Info("Configuring Cloud SQL Go Connector", zap.String("instance", instanceConnName))
		
		dialer, err := cloudsqlconn.NewDialer(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Cloud SQL dialer: %w", err)
		}

		// Intercept the network dialer to route through GCP
		config.ConnConfig.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
			return dialer.Dial(ctx, instanceConnName)
		}

		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				logger.Info("Closing Cloud SQL dialer")
				return dialer.Close()
			},
		})
	} else {
		logger.Info("Using standard TCP database connection", zap.String("url", dbURL))
	}

	// ---------------------------------------------------------
	// MIGRATIONS USING THE SHARED CONNECTION TUNNEL
	// ---------------------------------------------------------
	logger.Info("Running database migrations")

	// Convert the custom pgx config into a standard library *sql.DB
	// This ensures golang-migrate uses the exact same Cloud SQL tunnel (or local TCP)
	sqlDB := stdlib.OpenDB(*config.ConnConfig)
	
	// Create the migrate database driver using our configured sqlDB
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		// Ensure we don't leak the temporary connection on failure
		sqlDB.Close()
		return nil, fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", 
		driver,
	)
	if err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to initialize migrate: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	
	// Clean up migration resources. m.Close() closes the underlying sqlDB for us.
	if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
		logger.Warn("Error closing migration resources", zap.Error(dbErr))
	} else {
		logger.Info("Database migrations applied successfully")
	}
	// ---------------------------------------------------------

	// 3. Create the main application connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 4. Register pool closure with fx lifecycle
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection pool")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}