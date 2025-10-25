package postgresql

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/module"
)

type ModuleRepositoryTest struct {
	Name       string
	Repository *ModuleRepository
}

func TestModuleRepository(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UTC().UnixNano())

	repositories := createModuleRepositories(t)

	for i := range repositories {
		r := repositories[i]

		t.Run(r.Name, func(t *testing.T) {
			t.Parallel()

			t.Run("CreateAndGet", func(t *testing.T) {
				t.Parallel()
				testModuleCreateAndGet(t, r.Repository)
			})
			t.Run("Update", func(t *testing.T) {
				t.Parallel()
				testModuleUpdate(t, r.Repository)
			})
			t.Run("Delete", func(t *testing.T) {
				t.Parallel()
				testModuleDelete(t, r.Repository)
			})
			t.Run("GetByCourseID", func(t *testing.T) {
				t.Parallel()
				testModuleGetByCourseID(t, r.Repository)
			})
			t.Run("Exists", func(t *testing.T) {
				t.Parallel()
				testModuleExists(t, r.Repository)
			})
			t.Run("UpdateOrder", func(t *testing.T) {
				t.Parallel()
				testModuleUpdateOrder(t, r.Repository)
			})
			t.Run("ReorderModules", func(t *testing.T) {
				t.Parallel()
				testModuleReorder(t, r.Repository)
			})
		})
	}
}

func createModuleRepositories(t *testing.T) []ModuleRepositoryTest {
	return []ModuleRepositoryTest{
		{
			Name:       "PostgreSQL",
			Repository: newPostgreSQLModuleRepository(t),
		},
	}
}

func testModuleCreateAndGet(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	// Create parent course first
	courseID := createTestCourse(t, ctx, repository.db)

	m, err := module.NewModule(
		generateModuleID(),
		courseID,
		"Introduction to Testing",
		1,
	)
	if err != nil {
		t.Fatalf("failed to create module domain model: %v", err)
	}

	if err := repository.Create(ctx, m); err != nil {
		t.Fatalf("failed to create module: %v", err)
	}

	retrieved, err := repository.Get(ctx, m.ID())
	if err != nil {
		t.Fatalf("failed to get module: %v", err)
	}

	assertModuleEqual(t, m, retrieved)
}

func testModuleUpdate(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	courseID := createTestCourse(t, ctx, repository.db)

	m, _ := module.NewModule(
		generateModuleID(),
		courseID,
		"Original Title",
		1,
	)

	repository.Create(ctx, m)

	m.UpdateTitle("Updated Title")
	m.UpdateOrder(5)

	if err := repository.Update(ctx, m); err != nil {
		t.Fatalf("failed to update module: %v", err)
	}

	retrieved, err := repository.Get(ctx, m.ID())
	if err != nil {
		t.Fatalf("failed to get updated module: %v", err)
	}

	if retrieved.Title() != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", retrieved.Title())
	}
	if retrieved.Order() != 5 {
		t.Errorf("expected order 5, got %d", retrieved.Order())
	}
}

func testModuleDelete(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	courseID := createTestCourse(t, ctx, repository.db)

	m, _ := module.NewModule(
		generateModuleID(),
		courseID,
		"Module to Delete",
		1,
	)

	repository.Create(ctx, m)

	exists, _ := repository.Exists(ctx, m.ID())
	if !exists {
		t.Fatal("module should exist before deletion")
	}

	if err := repository.Delete(ctx, m.ID()); err != nil {
		t.Fatalf("failed to delete module: %v", err)
	}

	exists, _ = repository.Exists(ctx, m.ID())
	if exists {
		t.Error("module should not exist after deletion")
	}
}

func testModuleGetByCourseID(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	courseID := createTestCourse(t, ctx, repository.db)

	modules := make([]*module.Module, 3)
	for i := 0; i < 3; i++ {
		m, _ := module.NewModule(
			generateModuleID(),
			courseID,
			fmt.Sprintf("Module %d", i+1),
			i+1,
		)
		repository.Create(ctx, m)
		modules[i] = m
	}

	otherModule, _ := module.NewModule(
		generateModuleID(),
		createTestCourse(t, ctx, repository.db),
		"Other Module",
		1,
	)
	repository.Create(ctx, otherModule)

	courseModules, err := repository.GetByCourseID(ctx, courseID)
	if err != nil {
		t.Fatalf("failed to get modules by course: %v", err)
	}

	if len(courseModules) != 3 {
		t.Errorf("expected 3 modules for course, got %d", len(courseModules))
	}

	for _, m := range courseModules {
		if m.CourseID() != courseID {
			t.Errorf("expected courseID '%s', got '%s'", courseID, m.CourseID())
		}
	}

	for i, m := range courseModules {
		expectedOrder := i + 1
		if m.Order() != expectedOrder {
			t.Errorf("expected module %d to have order %d, got %d", i, expectedOrder, m.Order())
		}
	}
}

func testModuleExists(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	courseID := createTestCourse(t, ctx, repository.db)
	moduleID := generateModuleID()

	exists, err := repository.Exists(ctx, moduleID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Error("module should not exist")
	}

	m, _ := module.NewModule(
		moduleID,
		courseID,
		"Existence Test",
		1,
	)
	repository.Create(ctx, m)

	exists, err = repository.Exists(ctx, moduleID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("module should exist")
	}
}

func testModuleUpdateOrder(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	courseID := createTestCourse(t, ctx, repository.db)

	m, _ := module.NewModule(
		generateModuleID(),
		courseID,
		"Module for Order Update",
		1,
	)

	repository.Create(ctx, m)

	if err := repository.UpdateOrder(ctx, m.ID(), 10); err != nil {
		t.Fatalf("failed to update module order: %v", err)
	}

	retrieved, err := repository.Get(ctx, m.ID())
	if err != nil {
		t.Fatalf("failed to get module: %v", err)
	}

	if retrieved.Order() != 10 {
		t.Errorf("expected order 10, got %d", retrieved.Order())
	}
}

func testModuleReorder(t *testing.T, repository *ModuleRepository) {
	ctx := context.Background()

	courseID := createTestCourse(t, ctx, repository.db)

	modules := make([]*module.Module, 4)
	for i := 0; i < 4; i++ {
		m, _ := module.NewModule(
			generateModuleID(),
			courseID,
			fmt.Sprintf("Module %d", i+1),
			i+1,
		)
		repository.Create(ctx, m)
		modules[i] = m
	}

	newOrders := map[string]int{
		modules[0].ID(): 4,
		modules[1].ID(): 3,
		modules[2].ID(): 2,
		modules[3].ID(): 1,
	}

	if err := repository.ReorderModules(ctx, newOrders); err != nil {
		t.Fatalf("failed to reorder modules: %v", err)
	}

	for _, m := range modules {
		retrieved, err := repository.Get(ctx, m.ID())
		if err != nil {
			t.Fatalf("failed to get module: %v", err)
		}

		expectedOrder := newOrders[m.ID()]
		if retrieved.Order() != expectedOrder {
			t.Errorf("module %s: expected order %d, got %d", m.ID(), expectedOrder, retrieved.Order())
		}
	}

	courseModules, err := repository.GetByCourseID(ctx, courseID)
	if err != nil {
		t.Fatalf("failed to get modules by course: %v", err)
	}

	for i, m := range courseModules {
		expectedOrder := i + 1
		if m.Order() != expectedOrder {
			t.Errorf("module at position %d: expected order %d, got %d", i, expectedOrder, m.Order())
		}
	}
}

func assertModuleEqual(t *testing.T, expected, actual *module.Module) {
	if actual.ID() != expected.ID() {
		t.Errorf("expected ID '%s', got '%s'", expected.ID(), actual.ID())
	}
	if actual.CourseID() != expected.CourseID() {
		t.Errorf("expected CourseID '%s', got '%s'", expected.CourseID(), actual.CourseID())
	}
	if actual.Title() != expected.Title() {
		t.Errorf("expected Title '%s', got '%s'", expected.Title(), actual.Title())
	}
	if actual.Order() != expected.Order() {
		t.Errorf("expected Order %d, got %d", expected.Order(), actual.Order())
	}
}

func generateModuleID() string {
	return fmt.Sprintf("test-mod-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}

func newPostgreSQLModuleRepository(t *testing.T) *ModuleRepository {
	// Setup testcontainer for PostgreSQL
	container, cleanup := SetupTestDatabase(t)
	t.Cleanup(cleanup)

	pool, err := pgxpool.New(context.Background(), container.ConnectionString)
	if err != nil {
		t.Fatalf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("unable to ping database: %v", err)
	}

	return NewModuleRepository(pool)
}

// createTestCourse creates a test course in the database and returns its ID
func createTestCourse(t *testing.T, ctx context.Context, pool *pgxpool.Pool) string {
	courseID := fmt.Sprintf("test-course-%d-%d", time.Now().UnixNano(), rand.Intn(10000))

	query := `
		INSERT INTO courses (id, teacher_id, title, description, duration, domain, rating, level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := pool.Exec(ctx, query,
		courseID,
		"test-teacher-"+generateModuleID(),
		"Test Course",
		"Test Description",
		3600,
		"programming",
		4.5,
		"beginner",
	)
	if err != nil {
		t.Fatalf("failed to create test course: %v", err)
	}

	return courseID
}
