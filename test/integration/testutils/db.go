package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestDB starts an ephemeral Postgres container, runs migrations, and returns a pgxpool.
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	// Spin up PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("ticketer_test"),
		postgres.WithUsername("ticketer_test"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err, "failed to start testcontainer")

	// Cleanup when test completes
	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Run Migrations
	// Assuming tests are run from internal package or test/integration
	migrationsPath := "file://../../../migrations"
	m, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		// Fallback for running from root
		migrationsPath = "file://../../migrations"
		m, err = migrate.New(migrationsPath, connStr)
		require.NoError(t, err, "failed to initialize migrations from either path")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err, "failed to apply migrations")
	}

	// Connect pgxpool
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
