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
	// Start a transaction for creating course with all nested entities
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

	// Create modules with lessons and exercises
	for _, module := range c.Modules() {
		if err := r.createModule(ctx, qtx, c.ID(), module); err != nil {
			return err
		}

		// Create lessons for this module
		for _, lesson := range module.Lessons() {
			if err := r.createLesson(ctx, qtx, module.ID(), lesson); err != nil {
				return err
			}

			// Create exercises for this lesson
			for _, exercise := range lesson.Exercises() {
				if err := r.createExercise(ctx, qtx, lesson.ID(), exercise); err != nil {
					return err
				}
			}
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
	if err := qtx.DeleteCourseTags(ctx, c.ID()); err != nil {
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

	// Delete existing modules (cascade will delete lessons and exercises)
	if err := qtx.DeleteModulesByCourseID(ctx, c.ID()); err != nil {
		return errors.Wrap(err, "failed to delete modules")
	}

	// Recreate modules with lessons and exercises
	for _, module := range c.Modules() {
		if err := r.createModule(ctx, qtx, c.ID(), module); err != nil {
			return err
		}

		for _, lesson := range module.Lessons() {
			if err := r.createLesson(ctx, qtx, module.ID(), lesson); err != nil {
				return err
			}

			for _, exercise := range lesson.Exercises() {
				if err := r.createExercise(ctx, qtx, lesson.ID(), exercise); err != nil {
					return err
				}
			}
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

func (r *CourseRepository) createModule(ctx context.Context, q *database.Queries, courseID string, m course.Module) error {
	params := database.CreateModuleParams{
		ID:         m.ID(),
		CourseID:   courseID,
		Title:      m.Title(),
		OrderIndex: int32(m.Order()),
	}

	if err := q.CreateModule(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create module")
	}

	return nil
}

func (r *CourseRepository) createLesson(ctx context.Context, q *database.Queries, moduleID string, l course.Lesson) error {
	overview := pgtype.Text{String: l.Overview(), Valid: l.Overview() != ""}
	content := pgtype.Text{String: l.Content(), Valid: l.Content() != ""}
	videoID := pgtype.Text{String: l.VideoID(), Valid: l.VideoID() != ""}

	params := database.CreateLessonParams{
		ID:         l.ID(),
		ModuleID:   moduleID,
		Title:      l.Title(),
		Overview:   overview,
		Content:    content,
		VideoID:    videoID,
		OrderIndex: int32(l.Order()),
	}

	if err := q.CreateLesson(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create lesson")
	}

	return nil
}

func (r *CourseRepository) createExercise(ctx context.Context, q *database.Queries, lessonID string, e course.Exercise) error {
	params := database.CreateExerciseParams{
		ID:            e.ID(),
		LessonID:      lessonID,
		Question:      e.Question(),
		Answers:       e.Answers(),
		CorrectAnswer: "", // Note: We don't expose correct answer from domain
		OrderIndex:    int32(e.Order()),
	}

	if err := q.CreateExercise(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create exercise")
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

	// Get modules
	dbModules, err := r.queries.GetModulesByCourseID(ctx, dbCourse.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get modules")
	}

	modules := make([]course.Module, 0, len(dbModules))
	for _, dbModule := range dbModules {
		domainModule, err := r.toDomainModule(ctx, dbModule)
		if err != nil {
			return nil, err
		}
		modules = append(modules, *domainModule)
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
		modules,
	)
}

func (r *CourseRepository) toDomainModule(ctx context.Context, dbModule database.Module) (*course.Module, error) {
	// Get lessons for this module
	dbLessons, err := r.queries.GetLessonsByModuleID(ctx, dbModule.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lessons")
	}

	lessons := make([]course.Lesson, 0, len(dbLessons))
	for _, dbLesson := range dbLessons {
		domainLesson, err := r.toDomainLesson(ctx, dbLesson)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, *domainLesson)
	}

	return course.NewModule(
		dbModule.ID,
		dbModule.Title,
		lessons,
		int(dbModule.OrderIndex),
	)
}

func (r *CourseRepository) toDomainLesson(ctx context.Context, dbLesson database.Lesson) (*course.Lesson, error) {
	// Get exercises for this lesson
	dbExercises, err := r.queries.GetExercisesByLessonID(ctx, dbLesson.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exercises")
	}

	exercises := make([]course.Exercise, 0, len(dbExercises))
	for _, dbExercise := range dbExercises {
		domainExercise, err := r.toDomainExercise(dbExercise)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, *domainExercise)
	}

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

	return course.NewLesson(
		dbLesson.ID,
		dbLesson.Title,
		overview,
		content,
		videoID,
		exercises,
		int(dbLesson.OrderIndex),
	)
}

func (r *CourseRepository) toDomainExercise(dbExercise database.Exercise) (*course.Exercise, error) {
	return course.NewExercise(
		dbExercise.ID,
		dbExercise.Question,
		dbExercise.Answers,
		dbExercise.CorrectAnswer,
		int(dbExercise.OrderIndex),
	)
}
