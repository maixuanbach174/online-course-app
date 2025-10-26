package lesson

import "context"

// LessonRepository manages Lesson aggregate persistence
type LessonRepository interface {
	// Create saves a new lesson to the database
	Create(ctx context.Context, lesson *Lesson) error

	// Update modifies an existing lesson
	Update(ctx context.Context, lesson *Lesson) error

	// Delete removes a lesson and all associated exercises via CASCADE
	Delete(ctx context.Context, id string) error

	// Get retrieves a lesson by ID
	Get(ctx context.Context, id string) (*Lesson, error)

	// GetByModuleID retrieves all lessons for a specific module, ordered by the order field
	GetByModuleID(ctx context.Context, moduleID string) ([]*Lesson, error)

	// Exists checks if a lesson with the given ID exists
	Exists(ctx context.Context, id string) (bool, error)

	// ReorderLessons updates the order of multiple lessons in a single transaction
	ReorderLessons(ctx context.Context, lessonOrders map[string]int) error
}
