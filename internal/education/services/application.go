package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/common/metrics"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql"
	"github.com/maixuanbach174/online-course-app/internal/education/app"
	"github.com/maixuanbach174/online-course-app/internal/education/app/command"
	"github.com/maixuanbach174/online-course-app/internal/education/app/command/course_command"
	"github.com/maixuanbach174/online-course-app/internal/education/app/query/course_query"
	"github.com/sirupsen/logrus"
)

// ApplicationContainer holds the application and its resources for proper cleanup
type ApplicationContainer struct {
	App     app.Application
	pool    *pgxpool.Pool
	logger  *logrus.Entry
	cleanup func()
}

// Close gracefully closes all application resources
func (ac *ApplicationContainer) Close() {
	if ac.cleanup != nil {
		ac.cleanup()
	}
}

// NewApplication creates a new application with proper resource management
// Returns an ApplicationContainer that must be closed when done
func NewApplication(ctx context.Context) (*ApplicationContainer, error) {
	config := LoadConfig()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	// Create a single connection pool shared by all repositories
	pool, err := newConnectionPool(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify the connection is working
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create repositories using the shared pool
	userRepository := postgresql.NewUserRepository(pool)
	courseRepository := postgresql.NewCourseRepository(pool)

	application := app.Application{
		Commands: app.Commands{
			RegisterUser: command.NewRegisterUserHandler(userRepository, logger, metricsClient),
			CreateCourse: course_command.NewCreateCourseHandler(courseRepository, logger, metricsClient),
			DeleteCourse: course_command.NewDeleteCourseHandler(courseRepository, logger, metricsClient),
			UpdateCourse: course_command.NewUpdateCourseHandler(courseRepository, logger, metricsClient),
		},
		Queries: app.Queries{
			GetAllCourses:    course_query.NewGetAllCoursesHandler(courseRepository, logger, metricsClient),
			GetCourseDetails: course_query.NewGetCourseDetailsHandler(courseRepository, logger, metricsClient),
			CoursesByTeacher: course_query.NewCourseByTeacherHandler(courseRepository, logger, metricsClient),
		},
	}

	container := &ApplicationContainer{
		App:    application,
		pool:   pool,
		logger: logger,
		cleanup: func() {
			logger.Info("Closing database connection pool")
			pool.Close()
		},
	}

	return container, nil
}

// newConnectionPool creates a new database connection pool with proper configuration
func newConnectionPool(ctx context.Context, config *Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.PostgresURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings for better resource management
	poolConfig.MaxConns = 25         // Maximum number of connections
	poolConfig.MinConns = 5          // Minimum number of connections to keep open
	poolConfig.MaxConnLifetime = 0   // Connections live forever (unless explicitly closed)
	poolConfig.MaxConnIdleTime = 0   // Keep idle connections
	poolConfig.HealthCheckPeriod = 0 // Disable health checks (can be enabled if needed)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
