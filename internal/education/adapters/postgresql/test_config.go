package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver for database/sql
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgreSQLContainer wraps the testcontainer for PostgreSQL
type PostgreSQLContainer struct {
	*pgcontainer.PostgresContainer
	ConnectionString string
}

// SetupTestDatabase creates a PostgreSQL testcontainer with migrations applied
func SetupTestDatabase(t *testing.T) (*PostgreSQLContainer, func()) {
	ctx := context.Background()

	// Create PostgreSQL container
	pgContainer, err := pgcontainer.Run(ctx,
		"postgres:16-alpine",
		pgcontainer.WithDatabase("online_course_test"),
		pgcontainer.WithUsername("test"),
		pgcontainer.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	container := &PostgreSQLContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}

	// Run migrations
	if err := runMigrations(connStr, t); err != nil {
		pgContainer.Terminate(ctx)
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	}

	return container, cleanup
}

// runMigrations applies all migrations from the migrations directory
func runMigrations(connStr string, t *testing.T) error {
	// Open database connection for migrations
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get the migrations directory path
	migrationsPath := getMigrationsPath()
	t.Logf("Using migrations from: %s", migrationsPath)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run all up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in dirty state at version %d", version)
	}

	t.Logf("Migrations applied successfully (version: %d)", version)
	return nil
}

// getMigrationsPath returns the absolute path to the migrations directory
func getMigrationsPath() string {
	// Get the current file's directory
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	// Navigate up to the education module root, then to migrations
	// test_config.go is in: internal/education/adapters/postgresql/
	// migrations is in: internal/education/migrations/
	migrationsDir := filepath.Join(currentDir, "..", "..", "migrations")

	// Get absolute path
	absPath, _ := filepath.Abs(migrationsDir)
	return absPath
}

// RunMigrationsDown rolls back all migrations (useful for cleanup in some test scenarios)
func RunMigrationsDown(connStr string) error {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	migrationsPath := getMigrationsPath()
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

// GetConnectionPool creates a pgxpool connection pool for tests
func GetConnectionPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
