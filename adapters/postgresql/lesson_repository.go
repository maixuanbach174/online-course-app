package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql/database"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/lesson"
	"github.com/pkg/errors"
)

type LessonRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewLessonRepository(db *pgxpool.Pool) *LessonRepository {
	return &LessonRepository{
		db:      db,
		queries: database.New(db),
	}
}

// Create implements lesson.LessonRepository
func (r *LessonRepository) Create(ctx context.Context, l *lesson.Lesson) error {
	overview := pgtype.Text{String: l.Overview(), Valid: l.Overview() != ""}
	content := pgtype.Text{String: l.Content(), Valid: l.Content() != ""}
	videoID := pgtype.Text{String: l.VideoID(), Valid: l.VideoID() != ""}

	params := database.CreateLessonParams{
		ID:         l.ID(),
		ModuleID:   l.ModuleID(),
		Title:      l.Title(),
		Overview:   overview,
		Content:    content,
		VideoID:    videoID,
		Duration:   int32(l.Duration()),
		OrderIndex: int32(l.Order()),
	}

	if err := r.queries.CreateLesson(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create lesson")
	}

	return nil
}

// Update implements lesson.LessonRepository
func (r *LessonRepository) Update(ctx context.Context, l *lesson.Lesson) error {
	overview := pgtype.Text{String: l.Overview(), Valid: l.Overview() != ""}
	content := pgtype.Text{String: l.Content(), Valid: l.Content() != ""}
	videoID := pgtype.Text{String: l.VideoID(), Valid: l.VideoID() != ""}

	params := database.UpdateLessonParams{
		ID:         l.ID(),
		Title:      l.Title(),
		Overview:   overview,
		Content:    content,
		VideoID:    videoID,
		Duration:   int32(l.Duration()),
		OrderIndex: int32(l.Order()),
	}

	if err := r.queries.UpdateLesson(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update lesson")
	}

	return nil
}

// Delete implements lesson.LessonRepository
func (r *LessonRepository) Delete(ctx context.Context, id string) error {
	if err := r.queries.DeleteLesson(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete lesson")
	}

	return nil
}

// Get implements lesson.LessonRepository
func (r *LessonRepository) Get(ctx context.Context, id string) (*lesson.Lesson, error) {
	dbLesson, err := r.queries.GetLessonByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lesson")
	}

	return r.toDomainLesson(dbLesson)
}

// GetByModuleID implements lesson.LessonRepository
func (r *LessonRepository) GetByModuleID(ctx context.Context, moduleID string) ([]*lesson.Lesson, error) {
	dbLessons, err := r.queries.GetLessonsByModuleID(ctx, moduleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lessons by module")
	}

	lessons := make([]*lesson.Lesson, 0, len(dbLessons))
	for _, dbLesson := range dbLessons {
		domainLesson, err := r.toDomainLesson(dbLesson)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, domainLesson)
	}

	return lessons, nil
}

// Exists implements lesson.LessonRepository
func (r *LessonRepository) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.queries.LessonExists(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to check lesson existence")
	}
	return exists, nil
}

// UpdateOrder implements lesson.LessonRepository
func (r *LessonRepository) UpdateOrder(ctx context.Context, id string, order int) error {
	params := database.UpdateLessonOrderParams{
		ID:         id,
		OrderIndex: int32(order),
	}

	if err := r.queries.UpdateLessonOrder(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update lesson order")
	}

	return nil
}

// ReorderLessons implements lesson.LessonRepository
func (r *LessonRepository) ReorderLessons(ctx context.Context, lessonOrders map[string]int) error {
	// Start a transaction for atomic update
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Update each lesson's order
	for lessonID, order := range lessonOrders {
		params := database.UpdateLessonOrderParams{
			ID:         lessonID,
			OrderIndex: int32(order),
		}
		if err := qtx.UpdateLessonOrder(ctx, params); err != nil {
			return errors.Wrap(err, "failed to update lesson order")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Helper methods

func (r *LessonRepository) toDomainLesson(dbLesson database.Lesson) (*lesson.Lesson, error) {
	overview := ""
	if dbLesson.Overview.Valid {
		overview = dbLesson.Overview.String
	}

	content := ""
	if dbLesson.Content.Valid {
		content = dbLesson.Content.String
	}

	videoID := ""
	if dbLesson.VideoID.Valid {
		videoID = dbLesson.VideoID.String
	}

	return lesson.NewLesson(
		dbLesson.ID,
		dbLesson.ModuleID,
		dbLesson.Title,
		overview,
		content,
		videoID,
		int(dbLesson.Duration),
		int(dbLesson.OrderIndex),
	)
}
