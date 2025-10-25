package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql/database"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/exercise"
	"github.com/pkg/errors"
)

type ExerciseRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewExerciseRepository(db *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{
		db:      db,
		queries: database.New(db),
	}
}

// Create implements exercise.ExerciseRepository
func (r *ExerciseRepository) Create(ctx context.Context, e *exercise.Exercise) error {
	params := database.CreateExerciseParams{
		ID:            e.ID(),
		LessonID:      e.LessonID(),
		Question:      e.Question(),
		Answers:       e.Answers(),
		CorrectAnswer: e.CorrectAnswerForStorage(), // Special method to get correct answer (only for storage)
		OrderIndex:    int32(e.Order()),
	}

	if err := r.queries.CreateExercise(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create exercise")
	}

	return nil
}

// Update implements exercise.ExerciseRepository
func (r *ExerciseRepository) Update(ctx context.Context, e *exercise.Exercise) error {
	params := database.UpdateExerciseParams{
		ID:            e.ID(),
		Question:      e.Question(),
		Answers:       e.Answers(),
		CorrectAnswer: e.CorrectAnswerForStorage(), // Special method to get correct answer (only for storage)
		OrderIndex:    int32(e.Order()),
	}

	if err := r.queries.UpdateExercise(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update exercise")
	}

	return nil
}

// Delete implements exercise.ExerciseRepository
func (r *ExerciseRepository) Delete(ctx context.Context, id string) error {
	if err := r.queries.DeleteExercise(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete exercise")
	}

	return nil
}

// Get implements exercise.ExerciseRepository
func (r *ExerciseRepository) Get(ctx context.Context, id string) (*exercise.Exercise, error) {
	dbExercise, err := r.queries.GetExerciseByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exercise")
	}

	return r.toDomainExercise(dbExercise)
}

// GetByLessonID implements exercise.ExerciseRepository
func (r *ExerciseRepository) GetByLessonID(ctx context.Context, lessonID string) ([]*exercise.Exercise, error) {
	dbExercises, err := r.queries.GetExercisesByLessonID(ctx, lessonID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exercises by lesson")
	}

	exercises := make([]*exercise.Exercise, 0, len(dbExercises))
	for _, dbExercise := range dbExercises {
		domainExercise, err := r.toDomainExercise(dbExercise)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, domainExercise)
	}

	return exercises, nil
}

// Exists implements exercise.ExerciseRepository
func (r *ExerciseRepository) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.queries.ExerciseExists(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to check exercise existence")
	}
	return exists, nil
}

// UpdateOrder implements exercise.ExerciseRepository
func (r *ExerciseRepository) UpdateOrder(ctx context.Context, id string, order int) error {
	params := database.UpdateExerciseOrderParams{
		ID:         id,
		OrderIndex: int32(order),
	}

	if err := r.queries.UpdateExerciseOrder(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update exercise order")
	}

	return nil
}

// ReorderExercises implements exercise.ExerciseRepository
func (r *ExerciseRepository) ReorderExercises(ctx context.Context, exerciseOrders map[string]int) error {
	// Start a transaction for atomic update
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Update each exercise's order
	for exerciseID, order := range exerciseOrders {
		params := database.UpdateExerciseOrderParams{
			ID:         exerciseID,
			OrderIndex: int32(order),
		}
		if err := qtx.UpdateExerciseOrder(ctx, params); err != nil {
			return errors.Wrap(err, "failed to update exercise order")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Helper methods

func (r *ExerciseRepository) toDomainExercise(dbExercise database.Exercise) (*exercise.Exercise, error) {
	return exercise.NewExercise(
		dbExercise.ID,
		dbExercise.LessonID,
		dbExercise.Question,
		dbExercise.Answers,
		dbExercise.CorrectAnswer,
		int(dbExercise.OrderIndex),
	)
}
