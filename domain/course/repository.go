package course

import "context"

// CourseRepository manages Course aggregate persistence
type CourseRepository interface {
	// Create saves a new course to the database
	Create(ctx context.Context, course *Course) error

	// Update modifies an existing course metadata
	Update(ctx context.Context, course *Course) error

	// Delete removes a course and all associated data (modules, lessons, exercises) via CASCADE
	Delete(ctx context.Context, id string) error

	// Get retrieves a course by ID
	Get(ctx context.Context, id string) (*Course, error)

	// GetAll retrieves all courses
	GetAll(ctx context.Context) ([]*Course, error)

	// GetAllByTeacherID retrieves all courses for a specific teacher
	GetAllByTeacherID(ctx context.Context, teacherID string) ([]*Course, error)

	// Exists checks if a course with the given ID exists
	Exists(ctx context.Context, id string) (bool, error)
}
