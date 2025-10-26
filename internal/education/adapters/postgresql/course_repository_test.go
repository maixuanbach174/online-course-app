package postgresql

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
)

type CourseRepositoryTest struct {
	Name       string
	Repository *CourseRepository
}

func TestCourseRepository(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UTC().UnixNano())

	repositories := createCourseRepositories(t)

	for i := range repositories {
		// When you are looping over slice and later using iterated value in goroutine (here because of t.Parallel()),
		// you need to always create variable scoped in loop body!
		// More info here: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		r := repositories[i]

		t.Run(r.Name, func(t *testing.T) {
			// It's always a good idea to build all non-unit tests to be able to work in parallel.
			// Thanks to that, your tests will be always fast and you will not be afraid to add more tests because of slowdown.
			t.Parallel()

			t.Run("CreateAndGet", func(t *testing.T) {
				t.Parallel()
				testCourseCreateAndGet(t, r.Repository)
			})
			t.Run("Update", func(t *testing.T) {
				t.Parallel()
				testCourseUpdate(t, r.Repository)
			})
			t.Run("Delete", func(t *testing.T) {
				t.Parallel()
				testCourseDelete(t, r.Repository)
			})
			t.Run("GetAll", func(t *testing.T) {
				t.Parallel()
				testCourseGetAll(t, r.Repository)
			})
			t.Run("GetAllByTeacherID", func(t *testing.T) {
				t.Parallel()
				testCourseGetAllByTeacherID(t, r.Repository)
			})
			t.Run("Exists", func(t *testing.T) {
				t.Parallel()
				testCourseExists(t, r.Repository)
			})
			t.Run("CourseWithTags", func(t *testing.T) {
				t.Parallel()
				testCourseWithTags(t, r.Repository)
			})
			t.Run("UpdateCourseTags", func(t *testing.T) {
				t.Parallel()
				testUpdateCourseTags(t, r.Repository)
			})
		})
	}
}

func createCourseRepositories(t *testing.T) []CourseRepositoryTest {
	return []CourseRepositoryTest{
		{
			Name:       "PostgreSQL",
			Repository: newPostgreSQLCourseRepository(t),
		},
	}
}

func testCourseCreateAndGet(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	// Create a test course
	c, err := course.NewCourse(
		generateID(),
		"teacher-"+generateID(),
		"Test Course",
		"Test Description",
		"test-thumbnail.jpg",
		3600,
		course.DomainProgramming,
		[]course.Tag{course.TagBackend},
		4.5,
		course.Beginner,
	)
	if err != nil {
		t.Fatalf("failed to create course domain model: %v", err)
	}

	// Save to repository
	if err := repository.Create(ctx, c); err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Retrieve from repository
	retrieved, err := repository.Get(ctx, c.ID())
	if err != nil {
		t.Fatalf("failed to get course: %v", err)
	}

	// Verify
	assertCourseEqual(t, c, retrieved)
}

func testCourseUpdate(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	// Create a test course
	c, _ := course.NewCourse(
		generateID(),
		"teacher-"+generateID(),
		"Original Title",
		"Original Description",
		"original.jpg",
		3600,
		course.DomainProgramming,
		[]course.Tag{course.TagBackend},
		4.0,
		course.Beginner,
	)

	repository.Create(ctx, c)

	// Update course
	c.UpdateBasicInfo("Updated Title", "Updated Description", "updated.jpg")
	c.UpdateDuration(7200)
	c.UpdateRating(4.8)

	if err := repository.Update(ctx, c); err != nil {
		t.Fatalf("failed to update course: %v", err)
	}

	// Retrieve and verify
	retrieved, err := repository.Get(ctx, c.ID())
	if err != nil {
		t.Fatalf("failed to get updated course: %v", err)
	}

	if retrieved.Title() != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", retrieved.Title())
	}
	if retrieved.Description() != "Updated Description" {
		t.Errorf("expected description 'Updated Description', got '%s'", retrieved.Description())
	}
	if retrieved.Thumbnail() != "updated.jpg" {
		t.Errorf("expected thumbnail 'updated.jpg', got '%s'", retrieved.Thumbnail())
	}
	if retrieved.Duration() != 7200 {
		t.Errorf("expected duration 7200, got %d", retrieved.Duration())
	}
	if retrieved.Rating() != 4.8 {
		t.Errorf("expected rating 4.8, got %f", retrieved.Rating())
	}
}

func testCourseDelete(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	// Create a test course
	c, _ := course.NewCourse(
		generateID(),
		"teacher-"+generateID(),
		"Course to Delete",
		"",
		"",
		0,
		course.DomainProgramming,
		nil,
		0,
		course.Beginner,
	)

	repository.Create(ctx, c)

	// Verify it exists
	exists, _ := repository.Exists(ctx, c.ID())
	if !exists {
		t.Fatal("course should exist before deletion")
	}

	// Delete
	if err := repository.Delete(ctx, c.ID()); err != nil {
		t.Fatalf("failed to delete course: %v", err)
	}

	// Verify it doesn't exist
	exists, _ = repository.Exists(ctx, c.ID())
	if exists {
		t.Error("course should not exist after deletion")
	}
}

func testCourseGetAll(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	// Create multiple courses with unique teacher IDs
	teacherID := "teacher-" + generateID()
	courses := make([]*course.Course, 3)
	for i := 0; i < 3; i++ {
		c, _ := course.NewCourse(
			generateID(),
			teacherID,
			fmt.Sprintf("Course %d", i+1),
			"",
			"",
			0,
			course.DomainProgramming,
			nil,
			0,
			course.Beginner,
		)
		repository.Create(ctx, c)
		courses[i] = c
	}

	// Get all courses
	allCourses, err := repository.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all courses: %v", err)
	}

	// Verify our courses are in the result (there might be other courses from parallel tests)
	foundCount := 0
	for _, created := range courses {
		for _, retrieved := range allCourses {
			if created.ID() == retrieved.ID() {
				foundCount++
				break
			}
		}
	}

	if foundCount != 3 {
		t.Errorf("expected to find 3 courses, found %d", foundCount)
	}
}

func testCourseGetAllByTeacherID(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	teacherID := "teacher-" + generateID()

	// Create courses for this teacher
	for i := 0; i < 3; i++ {
		c, _ := course.NewCourse(
			generateID(),
			teacherID,
			fmt.Sprintf("Teacher Course %d", i+1),
			"",
			"",
			0,
			course.DomainProgramming,
			nil,
			0,
			course.Beginner,
		)
		repository.Create(ctx, c)
	}

	// Create a course for different teacher
	otherCourse, _ := course.NewCourse(
		generateID(),
		"other-teacher-"+generateID(),
		"Other Course",
		"",
		"",
		0,
		course.DomainProgramming,
		nil,
		0,
		course.Beginner,
	)
	repository.Create(ctx, otherCourse)

	// Get courses by teacher ID
	teacherCourses, err := repository.GetAllByTeacherID(ctx, teacherID)
	if err != nil {
		t.Fatalf("failed to get courses by teacher: %v", err)
	}

	if len(teacherCourses) != 3 {
		t.Errorf("expected 3 courses for teacher, got %d", len(teacherCourses))
	}

	// Verify all courses belong to the teacher
	for _, c := range teacherCourses {
		if c.TeacherID() != teacherID {
			t.Errorf("expected teacherID '%s', got '%s'", teacherID, c.TeacherID())
		}
	}
}

func testCourseExists(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	courseID := generateID()

	// Check non-existent course
	exists, err := repository.Exists(ctx, courseID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Error("course should not exist")
	}

	// Create course
	c, _ := course.NewCourse(
		courseID,
		"teacher-"+generateID(),
		"Existence Test",
		"",
		"",
		0,
		course.DomainProgramming,
		nil,
		0,
		course.Beginner,
	)
	repository.Create(ctx, c)

	// Check existing course
	exists, err = repository.Exists(ctx, courseID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("course should exist")
	}
}

func testCourseWithTags(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	// Create course with multiple tags
	c, _ := course.NewCourse(
		generateID(),
		"teacher-"+generateID(),
		"Course with Tags",
		"",
		"",
		0,
		course.DomainProgramming,
		[]course.Tag{course.TagBackend, course.TagAPI, course.TagCloud, course.TagDatabase},
		4.2,
		course.Intermediate,
	)

	repository.Create(ctx, c)

	// Retrieve and verify tags
	retrieved, err := repository.Get(ctx, c.ID())
	if err != nil {
		t.Fatalf("failed to get course: %v", err)
	}

	if len(retrieved.Tags()) != 4 {
		t.Fatalf("expected 4 tags, got %d", len(retrieved.Tags()))
	}

	// Verify each tag
	expectedTags := map[string]bool{
		"backend":  false,
		"api":      false,
		"cloud":    false,
		"database": false,
	}

	for _, tag := range retrieved.Tags() {
		if _, ok := expectedTags[tag.String()]; ok {
			expectedTags[tag.String()] = true
		}
	}

	for tagName, found := range expectedTags {
		if !found {
			t.Errorf("expected tag '%s' not found", tagName)
		}
	}
}

func testUpdateCourseTags(t *testing.T, repository *CourseRepository) {
	ctx := context.Background()

	// Create course with initial tags
	c, _ := course.NewCourse(
		generateID(),
		"teacher-"+generateID(),
		"Course for Tag Update",
		"",
		"",
		0,
		course.DomainProgramming,
		[]course.Tag{course.TagBackend, course.TagAPI},
		0,
		course.Beginner,
	)

	repository.Create(ctx, c)

	// Update tags by removing one and adding new ones
	c.RemoveTag(course.TagBackend)
	c.AddTag(course.TagCloud)
	c.AddTag(course.TagDatabase)

	repository.Update(ctx, c)

	// Retrieve and verify updated tags
	retrieved, err := repository.Get(ctx, c.ID())
	if err != nil {
		t.Fatalf("failed to get course: %v", err)
	}

	if len(retrieved.Tags()) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(retrieved.Tags()))
	}

	if retrieved.HasTag(course.TagBackend) {
		t.Error("should not have TagBackend after removal")
	}
	if !retrieved.HasTag(course.TagAPI) {
		t.Error("should still have TagAPI")
	}
	if !retrieved.HasTag(course.TagCloud) {
		t.Error("should have TagCloud after adding")
	}
	if !retrieved.HasTag(course.TagDatabase) {
		t.Error("should have TagDatabase after adding")
	}
}

// Helper functions

func assertCourseEqual(t *testing.T, expected, actual *course.Course) {
	if actual.ID() != expected.ID() {
		t.Errorf("expected ID '%s', got '%s'", expected.ID(), actual.ID())
	}
	if actual.TeacherID() != expected.TeacherID() {
		t.Errorf("expected TeacherID '%s', got '%s'", expected.TeacherID(), actual.TeacherID())
	}
	if actual.Title() != expected.Title() {
		t.Errorf("expected Title '%s', got '%s'", expected.Title(), actual.Title())
	}
	if actual.Description() != expected.Description() {
		t.Errorf("expected Description '%s', got '%s'", expected.Description(), actual.Description())
	}
	if actual.Thumbnail() != expected.Thumbnail() {
		t.Errorf("expected Thumbnail '%s', got '%s'", expected.Thumbnail(), actual.Thumbnail())
	}
	if actual.Duration() != expected.Duration() {
		t.Errorf("expected Duration %d, got %d", expected.Duration(), actual.Duration())
	}
	if actual.Domain().String() != expected.Domain().String() {
		t.Errorf("expected Domain '%s', got '%s'", expected.Domain().String(), actual.Domain().String())
	}
	if actual.Rating() != expected.Rating() {
		t.Errorf("expected Rating %f, got %f", expected.Rating(), actual.Rating())
	}
	if actual.Level().String() != expected.Level().String() {
		t.Errorf("expected Level '%s', got '%s'", expected.Level().String(), actual.Level().String())
	}
	if len(actual.Tags()) != len(expected.Tags()) {
		t.Errorf("expected %d tags, got %d", len(expected.Tags()), len(actual.Tags()))
	}
}

func generateID() string {
	return fmt.Sprintf("test-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}

func newPostgreSQLCourseRepository(t *testing.T) *CourseRepository {
	// Setup testcontainer for PostgreSQL
	container, cleanup := SetupTestDatabase(t)
	t.Cleanup(cleanup)

	pool, err := pgxpool.New(context.Background(), container.ConnectionString)
	if err != nil {
		t.Fatalf("unable to create connection pool: %v", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("unable to ping database: %v", err)
	}

	return NewCourseRepository(pool)
}
