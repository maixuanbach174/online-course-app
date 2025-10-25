package course_command

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DeleteCourse struct {
	CourseID string
}

type DeleteCourseHandler decorator.CommandHandler[DeleteCourse]

type deleteCourseHandler struct {
	courseRepository course.CourseRepository
}

func NewDeleteCourseHandler(
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) DeleteCourseHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyCommandDecorators(
		deleteCourseHandler{
			courseRepository: courseRepository,
		},
		logger,
		metricsClient,
	)
}

func (h deleteCourseHandler) Handle(ctx context.Context, cmd DeleteCourse) error {
	// Validate input
	if cmd.CourseID == "" {
		return errors.New("course ID is required")
	}

	// Check if course exists
	exists, err := h.courseRepository.Exists(ctx, cmd.CourseID)
	if err != nil {
		return errors.Wrap(err, "failed to check course existence")
	}
	if !exists {
		return errors.New("course not found")
	}

	// Delete course (CASCADE will delete all modules, lessons, and exercises)
	if err := h.courseRepository.Delete(ctx, cmd.CourseID); err != nil {
		return errors.Wrap(err, "failed to delete course")
	}

	return nil
}
