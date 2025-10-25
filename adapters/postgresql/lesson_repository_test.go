package postgresql

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/lesson"
)

type LessonRepositoryTest struct {
	Name       string
	Repository *LessonRepository
}

func TestLessonRepository(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UTC().UnixNano())

	repositories := createLessonRepositories(t)

	for i := range repositories {
		r := repositories[i]

		t.Run(r.Name, func(t *testing.T) {
			t.Parallel()

			t.Run("CreateAndGet", func(t *testing.T) {
				t.Parallel()
				testLessonCreateAndGet(t, r.Repository)
			})
			t.Run("Update", func(t *testing.T) {
				t.Parallel()
				testLessonUpdate(t, r.Repository)
			})
			t.Run("Delete", func(t *testing.T) {
				t.Parallel()
				testLessonDelete(t, r.Repository)
			})
			t.Run("GetByModuleID", func(t *testing.T) {
				t.Parallel()
				testLessonGetByModuleID(t, r.Repository)
			})
			t.Run("Exists", func(t *testing.T) {
				t.Parallel()
				testLessonExists(t, r.Repository)
			})
			t.Run("UpdateOrder", func(t *testing.T) {
				t.Parallel()
				testLessonUpdateOrder(t, r.Repository)
			})
			t.Run("ReorderLessons", func(t *testing.T) {
				t.Parallel()
				testLessonReorder(t, r.Repository)
			})
			t.Run("LessonWithVideo", func(t *testing.T) {
				t.Parallel()
				testLessonWithVideo(t, r.Repository)
			})
			t.Run("UpdateLessonContent", func(t *testing.T) {
				t.Parallel()
				testLessonUpdateContent(t, r.Repository)
			})
		})
	}
}

func createLessonRepositories(t *testing.T) []LessonRepositoryTest {
	return []LessonRepositoryTest{
		{
			Name:       "PostgreSQL",
			Repository: newPostgreSQLLessonRepository(t),
		},
	}
}

func testLessonCreateAndGet(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	l, err := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Introduction to Testing",
		"Learn the basics",
		"Full content here",
		"video-123",
		900,
		1,
	)
	if err != nil {
		t.Fatalf("failed to create lesson domain model: %v", err)
	}

	if err := repository.Create(ctx, l); err != nil {
		t.Fatalf("failed to create lesson: %v", err)
	}

	retrieved, err := repository.Get(ctx, l.ID())
	if err != nil {
		t.Fatalf("failed to get lesson: %v", err)
	}

	assertLessonEqual(t, l, retrieved)
}

func testLessonUpdate(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	l, _ := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Original Title",
		"Original overview",
		"Original content",
		"video-123",
		900,
		1,
	)

	repository.Create(ctx, l)

	l.UpdateTitle("Updated Title")
	l.UpdateOverview("Updated overview")
	l.UpdateContent("Updated content")
	l.UpdateVideoID("video-456")
	l.UpdateDuration(1800)
	l.UpdateOrder(5)

	if err := repository.Update(ctx, l); err != nil {
		t.Fatalf("failed to update lesson: %v", err)
	}

	retrieved, err := repository.Get(ctx, l.ID())
	if err != nil {
		t.Fatalf("failed to get updated lesson: %v", err)
	}

	if retrieved.Title() != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", retrieved.Title())
	}
	if retrieved.Overview() != "Updated overview" {
		t.Errorf("expected overview 'Updated overview', got '%s'", retrieved.Overview())
	}
	if retrieved.Content() != "Updated content" {
		t.Errorf("expected content 'Updated content', got '%s'", retrieved.Content())
	}
	if retrieved.VideoID() != "video-456" {
		t.Errorf("expected videoID 'video-456', got '%s'", retrieved.VideoID())
	}
	if retrieved.Duration() != 1800 {
		t.Errorf("expected duration 1800, got %d", retrieved.Duration())
	}
	if retrieved.Order() != 5 {
		t.Errorf("expected order 5, got %d", retrieved.Order())
	}
}

func testLessonDelete(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	l, _ := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Lesson to Delete",
		"",
		"",
		"",
		0,
		1,
	)

	repository.Create(ctx, l)

	exists, _ := repository.Exists(ctx, l.ID())
	if !exists {
		t.Fatal("lesson should exist before deletion")
	}

	if err := repository.Delete(ctx, l.ID()); err != nil {
		t.Fatalf("failed to delete lesson: %v", err)
	}

	exists, _ = repository.Exists(ctx, l.ID())
	if exists {
		t.Error("lesson should not exist after deletion")
	}
}

func testLessonGetByModuleID(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	moduleID := createTestModule(t, ctx, repository.db)

	lessons := make([]*lesson.Lesson, 3)
	for i := 0; i < 3; i++ {
		l, _ := lesson.NewLesson(
			generateLessonID(),
			moduleID,
			fmt.Sprintf("Lesson %d", i+1),
			"",
			"",
			"",
			300*(i+1),
			i+1,
		)
		repository.Create(ctx, l)
		lessons[i] = l
	}

	otherLesson, _ := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Other Lesson",
		"",
		"",
		"",
		0,
		1,
	)
	repository.Create(ctx, otherLesson)

	moduleLessons, err := repository.GetByModuleID(ctx, moduleID)
	if err != nil {
		t.Fatalf("failed to get lessons by module: %v", err)
	}

	if len(moduleLessons) != 3 {
		t.Errorf("expected 3 lessons for module, got %d", len(moduleLessons))
	}

	for _, l := range moduleLessons {
		if l.ModuleID() != moduleID {
			t.Errorf("expected moduleID '%s', got '%s'", moduleID, l.ModuleID())
		}
	}

	for i, l := range moduleLessons {
		expectedOrder := i + 1
		if l.Order() != expectedOrder {
			t.Errorf("expected lesson %d to have order %d, got %d", i, expectedOrder, l.Order())
		}
	}
}

func testLessonExists(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	lessonID := generateLessonID()

	exists, err := repository.Exists(ctx, lessonID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Error("lesson should not exist")
	}

	l, _ := lesson.NewLesson(
		lessonID,
		createTestModule(t, ctx, repository.db),
		"Existence Test",
		"",
		"",
		"",
		0,
		1,
	)
	repository.Create(ctx, l)

	exists, err = repository.Exists(ctx, lessonID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("lesson should exist")
	}
}

func testLessonUpdateOrder(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	l, _ := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Lesson for Order Update",
		"",
		"",
		"",
		0,
		1,
	)

	repository.Create(ctx, l)

	if err := repository.UpdateOrder(ctx, l.ID(), 10); err != nil {
		t.Fatalf("failed to update lesson order: %v", err)
	}

	retrieved, err := repository.Get(ctx, l.ID())
	if err != nil {
		t.Fatalf("failed to get lesson: %v", err)
	}

	if retrieved.Order() != 10 {
		t.Errorf("expected order 10, got %d", retrieved.Order())
	}
}

func testLessonReorder(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	moduleID := createTestModule(t, ctx, repository.db)

	lessons := make([]*lesson.Lesson, 4)
	for i := 0; i < 4; i++ {
		l, _ := lesson.NewLesson(
			generateLessonID(),
			moduleID,
			fmt.Sprintf("Lesson %d", i+1),
			"",
			"",
			"",
			0,
			i+1,
		)
		repository.Create(ctx, l)
		lessons[i] = l
	}

	newOrders := map[string]int{
		lessons[0].ID(): 4,
		lessons[1].ID(): 3,
		lessons[2].ID(): 2,
		lessons[3].ID(): 1,
	}

	if err := repository.ReorderLessons(ctx, newOrders); err != nil {
		t.Fatalf("failed to reorder lessons: %v", err)
	}

	for _, l := range lessons {
		retrieved, err := repository.Get(ctx, l.ID())
		if err != nil {
			t.Fatalf("failed to get lesson: %v", err)
		}

		expectedOrder := newOrders[l.ID()]
		if retrieved.Order() != expectedOrder {
			t.Errorf("lesson %s: expected order %d, got %d", l.ID(), expectedOrder, retrieved.Order())
		}
	}

	moduleLessons, err := repository.GetByModuleID(ctx, moduleID)
	if err != nil {
		t.Fatalf("failed to get lessons by module: %v", err)
	}

	for i, l := range moduleLessons {
		expectedOrder := i + 1
		if l.Order() != expectedOrder {
			t.Errorf("lesson at position %d: expected order %d, got %d", i, expectedOrder, l.Order())
		}
	}
}

func testLessonWithVideo(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	l, _ := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Video Lesson",
		"Watch the video",
		"Additional content",
		"video-789",
		1200,
		1,
	)

	repository.Create(ctx, l)

	retrieved, err := repository.Get(ctx, l.ID())
	if err != nil {
		t.Fatalf("failed to get lesson: %v", err)
	}

	if !retrieved.HasVideo() {
		t.Error("lesson should have video")
	}
	if retrieved.VideoID() != "video-789" {
		t.Errorf("expected videoID 'video-789', got '%s'", retrieved.VideoID())
	}
}

func testLessonUpdateContent(t *testing.T, repository *LessonRepository) {
	ctx := context.Background()

	l, _ := lesson.NewLesson(
		generateLessonID(),
		createTestModule(t, ctx, repository.db),
		"Content Lesson",
		"Original overview",
		"Original content",
		"",
		0,
		1,
	)

	repository.Create(ctx, l)

	l.UpdateOverview("New overview with details")
	l.UpdateContent("This is a much longer and detailed content that explains everything in depth")

	repository.Update(ctx, l)

	retrieved, err := repository.Get(ctx, l.ID())
	if err != nil {
		t.Fatalf("failed to get lesson: %v", err)
	}

	if retrieved.Overview() != "New overview with details" {
		t.Errorf("unexpected overview: %s", retrieved.Overview())
	}
	if retrieved.Content() != "This is a much longer and detailed content that explains everything in depth" {
		t.Errorf("unexpected content: %s", retrieved.Content())
	}
}

func assertLessonEqual(t *testing.T, expected, actual *lesson.Lesson) {
	if actual.ID() != expected.ID() {
		t.Errorf("expected ID '%s', got '%s'", expected.ID(), actual.ID())
	}
	if actual.ModuleID() != expected.ModuleID() {
		t.Errorf("expected ModuleID '%s', got '%s'", expected.ModuleID(), actual.ModuleID())
	}
	if actual.Title() != expected.Title() {
		t.Errorf("expected Title '%s', got '%s'", expected.Title(), actual.Title())
	}
	if actual.Overview() != expected.Overview() {
		t.Errorf("expected Overview '%s', got '%s'", expected.Overview(), actual.Overview())
	}
	if actual.Content() != expected.Content() {
		t.Errorf("expected Content '%s', got '%s'", expected.Content(), actual.Content())
	}
	if actual.VideoID() != expected.VideoID() {
		t.Errorf("expected VideoID '%s', got '%s'", expected.VideoID(), actual.VideoID())
	}
	if actual.Duration() != expected.Duration() {
		t.Errorf("expected Duration %d, got %d", expected.Duration(), actual.Duration())
	}
	if actual.Order() != expected.Order() {
		t.Errorf("expected Order %d, got %d", expected.Order(), actual.Order())
	}
}

func generateLessonID() string {
	return fmt.Sprintf("test-lesson-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}

func newPostgreSQLLessonRepository(t *testing.T) *LessonRepository {
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

	return NewLessonRepository(pool)
}

// createTestModule creates a test module in the database and returns its ID
func createTestModule(t *testing.T, ctx context.Context, pool *pgxpool.Pool) string {
	// First create a course
	courseID := fmt.Sprintf("test-course-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
	query := `
		INSERT INTO courses (id, teacher_id, title, description, duration, domain, rating, level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := pool.Exec(ctx, query,
		courseID,
		"test-teacher-"+generateLessonID(),
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

	// Then create a module
	moduleID := fmt.Sprintf("test-module-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
	query2 := `
		INSERT INTO modules (id, course_id, title, order_index)
		VALUES ($1, $2, $3, $4)
	`
	_, err = pool.Exec(ctx, query2,
		moduleID,
		courseID,
		"Test Module",
		1,
	)
	if err != nil {
		t.Fatalf("failed to create test module: %v", err)
	}

	return moduleID
}
