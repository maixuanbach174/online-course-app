package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql/database"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/pkg/errors"
)

type CourseRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{
		db:      db,
		queries: database.New(db),
	}
}

// Create implements course.CourseRepository
func (r *CourseRepository) Create(ctx context.Context, c *course.Course) error {
	// Start a transaction for creating course with tags
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Create course
	if err := r.createCourse(ctx, qtx, c); err != nil {
		return err
	}

	// Create course tags
	for _, tag := range c.Tags() {
		if err := qtx.CreateCourseTag(ctx, database.CreateCourseTagParams{
			CourseID: c.ID(),
			Tag:      tag.String(),
		}); err != nil {
			return errors.Wrap(err, "failed to create course tag")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Update implements course.CourseRepository
func (r *CourseRepository) Update(ctx context.Context, c *course.Course) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Update course
	if err := r.updateCourse(ctx, qtx, c); err != nil {
		return err
	}

	// Delete existing tags and recreate
	if err := qtx.DeleteAllCourseTags(ctx, c.ID()); err != nil {
		return errors.Wrap(err, "failed to delete course tags")
	}

	for _, tag := range c.Tags() {
		if err := qtx.CreateCourseTag(ctx, database.CreateCourseTagParams{
			CourseID: c.ID(),
			Tag:      tag.String(),
		}); err != nil {
			return errors.Wrap(err, "failed to create course tag")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Delete implements course.CourseRepository
func (r *CourseRepository) Delete(ctx context.Context, id string) error {
	if err := r.queries.DeleteCourse(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete course")
	}

	return nil
}

// Get implements course.CourseRepository
func (r *CourseRepository) Get(ctx context.Context, id string) (*course.Course, error) {
	dbCourse, err := r.queries.GetCourseByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get course")
	}

	return r.toDomainCourse(ctx, dbCourse)
}

// GetAll implements course.CourseRepository
func (r *CourseRepository) GetAll(ctx context.Context) ([]*course.Course, error) {
	dbCourses, err := r.queries.GetAllCourses(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all courses")
	}

	courses := make([]*course.Course, 0, len(dbCourses))
	for _, dbCourse := range dbCourses {
		domainCourse, err := r.toDomainCourse(ctx, dbCourse)
		if err != nil {
			return nil, err
		}
		courses = append(courses, domainCourse)
	}

	return courses, nil
}

// GetAllByTeacherID implements course.CourseRepository
func (r *CourseRepository) GetAllByTeacherID(ctx context.Context, teacherID string) ([]*course.Course, error) {
	dbCourses, err := r.queries.GetCoursesByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get courses by teacher")
	}

	courses := make([]*course.Course, 0, len(dbCourses))
	for _, dbCourse := range dbCourses {
		domainCourse, err := r.toDomainCourse(ctx, dbCourse)
		if err != nil {
			return nil, err
		}
		courses = append(courses, domainCourse)
	}

	return courses, nil
}

// Exists implements course.CourseRepository
func (r *CourseRepository) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.queries.CourseExists(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to check course existence")
	}
	return exists, nil
}

// Helper methods

func (r *CourseRepository) createCourse(ctx context.Context, q *database.Queries, c *course.Course) error {
	description := pgtype.Text{String: c.Description(), Valid: c.Description() != ""}
	thumbnail := pgtype.Text{String: c.Thumbnail(), Valid: c.Thumbnail() != ""}

	var rating pgtype.Numeric
	if err := rating.Scan(c.Rating()); err != nil {
		return errors.Wrap(err, "failed to convert rating")
	}

	params := database.CreateCourseParams{
		ID:          c.ID(),
		TeacherID:   c.TeacherID(),
		Title:       c.Title(),
		Description: description,
		Thumbnail:   thumbnail,
		Duration:    int32(c.Duration()),
		Domain:      c.Domain().String(),
		Rating:      rating,
		Level:       c.Level().String(),
	}

	if err := q.CreateCourse(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create course")
	}

	return nil
}

func (r *CourseRepository) updateCourse(ctx context.Context, q *database.Queries, c *course.Course) error {
	description := pgtype.Text{String: c.Description(), Valid: c.Description() != ""}
	thumbnail := pgtype.Text{String: c.Thumbnail(), Valid: c.Thumbnail() != ""}

	var rating pgtype.Numeric
	if err := rating.Scan(c.Rating()); err != nil {
		return errors.Wrap(err, "failed to convert rating")
	}

	params := database.UpdateCourseParams{
		ID:          c.ID(),
		TeacherID:   c.TeacherID(),
		Title:       c.Title(),
		Description: description,
		Thumbnail:   thumbnail,
		Duration:    int32(c.Duration()),
		Domain:      c.Domain().String(),
		Rating:      rating,
		Level:       c.Level().String(),
	}

	if err := q.UpdateCourse(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update course")
	}

	return nil
}

func (r *CourseRepository) toDomainCourse(ctx context.Context, dbCourse database.Course) (*course.Course, error) {
	// Get tags
	dbTags, err := r.queries.GetCourseTagsByCourseID(ctx, dbCourse.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get course tags")
	}

	tags := make([]course.Tag, 0, len(dbTags))
	for _, tagStr := range dbTags {
		tag, err := course.NewTagFromString(tagStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid tag")
		}
		tags = append(tags, tag)
	}

	// Convert domain and level
	domain, err := course.NewDomainFromString(dbCourse.Domain)
	if err != nil {
		return nil, errors.Wrap(err, "invalid domain")
	}

	level, err := course.NewCourseLevelFromString(dbCourse.Level)
	if err != nil {
		return nil, errors.Wrap(err, "invalid level")
	}

	// Convert optional fields
	description := ""
	if dbCourse.Description.Valid {
		description = dbCourse.Description.String
	}

	thumbnail := ""
	if dbCourse.Thumbnail.Valid {
		thumbnail = dbCourse.Thumbnail.String
	}

	rating := 0.0
	if dbCourse.Rating.Valid {
		r, err := dbCourse.Rating.Float64Value()
		if err == nil {
			rating = r.Float64
		}
	}

	return course.NewCourse(
		dbCourse.ID,
		dbCourse.TeacherID,
		dbCourse.Title,
		description,
		thumbnail,
		int(dbCourse.Duration),
		domain,
		tags,
		rating,
		level,
	)
}
