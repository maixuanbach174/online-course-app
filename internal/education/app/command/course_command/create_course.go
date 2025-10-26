package course_command

import (
	"context"

	"github.com/google/uuid"
	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CreateCourse struct {
	TeacherID   string
	Title       string
	Description string
	Thumbnail   string
	Duration    int
	Domain      string
	Tags        []string
	Rating      float64
	Level       string
}

type CreateCourseHandler decorator.CommandHandler[CreateCourse]

type createCourseHandler struct {
	courseRepository course.CourseRepository
}

func NewCreateCourseHandler(
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateCourseHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyCommandDecorators(
		createCourseHandler{
			courseRepository: courseRepository,
		},
		logger,
		metricsClient,
	)
}

func (h createCourseHandler) Handle(ctx context.Context, cmd CreateCourse) error {
	// Validate input
	if cmd.TeacherID == "" {
		return errors.New("teacher ID is required")
	}
	if cmd.Title == "" {
		return errors.New("course title is required")
	}

	// Parse domain
	domain, err := course.NewDomainFromString(cmd.Domain)
	if err != nil {
		return errors.Wrap(err, "invalid domain")
	}

	// Parse level
	level, err := course.NewCourseLevelFromString(cmd.Level)
	if err != nil {
		return errors.Wrap(err, "invalid level")
	}

	// Parse tags
	tags := make([]course.Tag, 0, len(cmd.Tags))
	for _, tagStr := range cmd.Tags {
		tag, err := course.NewTagFromString(tagStr)
		if err != nil {
			return errors.Wrapf(err, "invalid tag: %s", tagStr)
		}
		tags = append(tags, tag)
	}

	// Create course entity
	newCourse, err := course.NewCourse(
		uuid.New().String(),
		cmd.TeacherID,
		cmd.Title,
		cmd.Description,
		cmd.Thumbnail,
		cmd.Duration,
		domain,
		tags,
		cmd.Rating,
		level,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create course entity")
	}

	// Persist to repository
	if err := h.courseRepository.Create(ctx, newCourse); err != nil {
		return errors.Wrap(err, "failed to save course")
	}

	return nil
}
