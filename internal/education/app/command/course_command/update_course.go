package course_command

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UpdateCourse struct {
	CourseID    string
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

type UpdateCourseHandler decorator.CommandHandler[UpdateCourse]

type updateCourseHandler struct {
	courseRepository course.CourseRepository
}

func NewUpdateCourseHandler(
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) UpdateCourseHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyCommandDecorators(
		updateCourseHandler{
			courseRepository: courseRepository,
		},
		logger,
		metricsClient,
	)
}

func (h updateCourseHandler) Handle(ctx context.Context, cmd UpdateCourse) error {
	// Validate input
	if cmd.CourseID == "" {
		return errors.New("course ID is required")
	}
	if cmd.TeacherID == "" {
		return errors.New("teacher ID is required")
	}
	if cmd.Title == "" {
		return errors.New("course title is required")
	}

	// Check if course exists
	exists, err := h.courseRepository.Exists(ctx, cmd.CourseID)
	if err != nil {
		return errors.Wrap(err, "failed to check course existence")
	}
	if !exists {
		return errors.New("course not found")
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

	// Create course entity with updated data
	updatedCourse, err := course.NewCourse(
		cmd.CourseID,
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
	if err := h.courseRepository.Update(ctx, updatedCourse); err != nil {
		return errors.Wrap(err, "failed to update course")
	}

	return nil
}
