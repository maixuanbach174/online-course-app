package exercise

import "context"

// ExerciseRepository manages Exercise aggregate persistence
type ExerciseRepository interface {
	// Create saves a new exercise to the database
	Create(ctx context.Context, exercise *Exercise) error

	// Update modifies an existing exercise
	Update(ctx context.Context, exercise *Exercise) error

	// Delete removes an exercise
	Delete(ctx context.Context, id string) error

	// Get retrieves an exercise by ID
	Get(ctx context.Context, id string) (*Exercise, error)

	// GetByLessonID retrieves all exercises for a specific lesson, ordered by the order field
	GetByLessonID(ctx context.Context, lessonID string) ([]*Exercise, error)

	// Exists checks if an exercise with the given ID exists
	Exists(ctx context.Context, id string) (bool, error)

	// ReorderExercises updates the order of multiple exercises in a single transaction
	ReorderExercises(ctx context.Context, exerciseOrders map[string]int) error
}
