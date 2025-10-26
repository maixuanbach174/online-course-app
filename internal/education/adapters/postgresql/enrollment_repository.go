package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql/database"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/enrollment"
	"github.com/pkg/errors"
)

type EnrollmentRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewEnrollmentRepository(db *pgxpool.Pool) *EnrollmentRepository {
	return &EnrollmentRepository{
		db:      db,
		queries: database.New(db),
	}
}

// Create implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) Create(ctx context.Context, e *enrollment.Enrollment) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Create enrollment
	if err := r.createEnrollment(ctx, qtx, e); err != nil {
		return err
	}

	// Create module progress
	for _, mp := range e.ModuleProgress() {
		if err := r.createModuleProgress(ctx, qtx, e.ID(), mp); err != nil {
			return err
		}
	}

	// Create lesson progress
	for _, lp := range e.LessonProgress() {
		if err := r.createLessonProgress(ctx, qtx, e.ID(), lp); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Update implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) Update(ctx context.Context, e *enrollment.Enrollment) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Update enrollment
	if err := r.updateEnrollment(ctx, qtx, e); err != nil {
		return err
	}

	// Delete existing progress records
	if err := qtx.DeleteModuleProgressByEnrollmentID(ctx, e.ID()); err != nil {
		return errors.Wrap(err, "failed to delete module progress")
	}

	if err := qtx.DeleteLessonProgressByEnrollmentID(ctx, e.ID()); err != nil {
		return errors.Wrap(err, "failed to delete lesson progress")
	}

	// Recreate module progress
	for _, mp := range e.ModuleProgress() {
		if err := r.createModuleProgress(ctx, qtx, e.ID(), mp); err != nil {
			return err
		}
	}

	// Recreate lesson progress
	for _, lp := range e.LessonProgress() {
		if err := r.createLessonProgress(ctx, qtx, e.ID(), lp); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Delete implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) Delete(ctx context.Context, id string) error {
	if err := r.queries.DeleteEnrollment(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete enrollment")
	}

	return nil
}

// Get implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) Get(ctx context.Context, id string) (*enrollment.Enrollment, error) {
	dbEnrollment, err := r.queries.GetEnrollmentByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get enrollment")
	}

	return r.toDomainEnrollment(ctx, dbEnrollment)
}

// GetAll implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) GetAll(ctx context.Context) ([]*enrollment.Enrollment, error) {
	dbEnrollments, err := r.queries.GetAllEnrollments(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all enrollments")
	}

	enrollments := make([]*enrollment.Enrollment, 0, len(dbEnrollments))
	for _, dbEnrollment := range dbEnrollments {
		domainEnrollment, err := r.toDomainEnrollment(ctx, dbEnrollment)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, domainEnrollment)
	}

	return enrollments, nil
}

// GetByUserAndCourse implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) GetByUserAndCourse(ctx context.Context, userID, courseID string) (*enrollment.Enrollment, error) {
	dbEnrollment, err := r.queries.GetEnrollmentByUserAndCourse(ctx, database.GetEnrollmentByUserAndCourseParams{
		UserID:   userID,
		CourseID: courseID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get enrollment by user and course")
	}

	return r.toDomainEnrollment(ctx, dbEnrollment)
}

// GetAllByUserID implements enrollment.EnrollmentRepository
func (r *EnrollmentRepository) GetAllByUserID(ctx context.Context, userID string) ([]*enrollment.Enrollment, error) {
	dbEnrollments, err := r.queries.GetEnrollmentsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get enrollments by user")
	}

	enrollments := make([]*enrollment.Enrollment, 0, len(dbEnrollments))
	for _, dbEnrollment := range dbEnrollments {
		domainEnrollment, err := r.toDomainEnrollment(ctx, dbEnrollment)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, domainEnrollment)
	}

	return enrollments, nil
}

// Helper methods

func (r *EnrollmentRepository) createEnrollment(ctx context.Context, q *database.Queries, e *enrollment.Enrollment) error {
	enrolledAt := pgtype.Timestamp{Time: e.EnrolledAt(), Valid: true}

	startedAt := pgtype.Timestamp{Valid: !e.StartedAt().IsZero()}
	if startedAt.Valid {
		startedAt.Time = e.StartedAt()
	}

	completedAt := pgtype.Timestamp{Valid: !e.CompletedAt().IsZero()}
	if completedAt.Valid {
		completedAt.Time = e.CompletedAt()
	}

	var courseProgressPercentage pgtype.Numeric
	if err := courseProgressPercentage.Scan(e.CourseProgress().Progress().ProgressPercentage()); err != nil {
		return errors.Wrap(err, "failed to convert course progress percentage")
	}

	params := database.CreateEnrollmentParams{
		ID:                       e.ID(),
		UserID:                   e.UserID(),
		CourseID:                 e.CourseID(),
		EnrolledAt:               enrolledAt,
		StartedAt:                startedAt,
		CompletedAt:              completedAt,
		CourseProgressPercentage: courseProgressPercentage,
		CourseProgressStatus:     e.CourseProgress().Progress().Status().String(),
	}

	if err := q.CreateEnrollment(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create enrollment")
	}

	return nil
}

func (r *EnrollmentRepository) updateEnrollment(ctx context.Context, q *database.Queries, e *enrollment.Enrollment) error {
	enrolledAt := pgtype.Timestamp{Time: e.EnrolledAt(), Valid: true}

	startedAt := pgtype.Timestamp{Valid: !e.StartedAt().IsZero()}
	if startedAt.Valid {
		startedAt.Time = e.StartedAt()
	}

	completedAt := pgtype.Timestamp{Valid: !e.CompletedAt().IsZero()}
	if completedAt.Valid {
		completedAt.Time = e.CompletedAt()
	}

	var courseProgressPercentage pgtype.Numeric
	if err := courseProgressPercentage.Scan(e.CourseProgress().Progress().ProgressPercentage()); err != nil {
		return errors.Wrap(err, "failed to convert course progress percentage")
	}

	params := database.UpdateEnrollmentParams{
		ID:                       e.ID(),
		UserID:                   e.UserID(),
		CourseID:                 e.CourseID(),
		EnrolledAt:               enrolledAt,
		StartedAt:                startedAt,
		CompletedAt:              completedAt,
		CourseProgressPercentage: courseProgressPercentage,
		CourseProgressStatus:     e.CourseProgress().Progress().Status().String(),
	}

	if err := q.UpdateEnrollment(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update enrollment")
	}

	return nil
}

func (r *EnrollmentRepository) createModuleProgress(ctx context.Context, q *database.Queries, enrollmentID string, mp enrollment.ModuleProgress) error {
	var progressPercentage pgtype.Numeric
	if err := progressPercentage.Scan(mp.Progress().ProgressPercentage()); err != nil {
		return errors.Wrap(err, "failed to convert module progress percentage")
	}

	params := database.CreateModuleProgressParams{
		EnrollmentID:       enrollmentID,
		ModuleID:           mp.ModuleID(),
		ProgressPercentage: progressPercentage,
		ProgressStatus:     mp.Progress().Status().String(),
	}

	if err := q.CreateModuleProgress(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create module progress")
	}

	return nil
}

func (r *EnrollmentRepository) createLessonProgress(ctx context.Context, q *database.Queries, enrollmentID string, lp enrollment.LessonProgress) error {
	var progressPercentage pgtype.Numeric
	if err := progressPercentage.Scan(lp.Progress().ProgressPercentage()); err != nil {
		return errors.Wrap(err, "failed to convert lesson progress percentage")
	}

	var exerciseScore pgtype.Numeric
	if err := exerciseScore.Scan(lp.ExerciseScore()); err != nil {
		return errors.Wrap(err, "failed to convert exercise score")
	}

	params := database.CreateLessonProgressParams{
		EnrollmentID:       enrollmentID,
		LessonID:           lp.LessonID(),
		ProgressPercentage: progressPercentage,
		ProgressStatus:     lp.Progress().Status().String(),
		ExerciseScore:      exerciseScore,
	}

	if err := q.CreateLessonProgress(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create lesson progress")
	}

	return nil
}

func (r *EnrollmentRepository) toDomainEnrollment(ctx context.Context, dbEnrollment database.Enrollment) (*enrollment.Enrollment, error) {
	// Create basic enrollment
	domainEnrollment, err := enrollment.NewEnrollment(
		dbEnrollment.ID,
		dbEnrollment.UserID,
		dbEnrollment.CourseID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create domain enrollment")
	}

	// Note: We need to reconstruct the enrollment state from database
	// This is a simplified version - you may need to add setters to your domain model
	// or create a factory method that accepts all fields

	// Get module progress (for future use if needed)
	_, err = r.queries.GetModuleProgressByEnrollmentID(ctx, dbEnrollment.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get module progress")
	}

	// Get lesson progress
	dbLessonProgress, err := r.queries.GetLessonProgressByEnrollmentID(ctx, dbEnrollment.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lesson progress")
	}

	// Convert lesson progress and complete lessons
	for _, dbLP := range dbLessonProgress {
		// If lesson is completed, mark it in the enrollment
		if dbLP.ProgressStatus == "completed" {
			if err := domainEnrollment.CompleteLesson(dbLP.LessonID); err != nil {
				return nil, errors.Wrap(err, "failed to complete lesson in domain")
			}
		}
	}

	return domainEnrollment, nil
}
