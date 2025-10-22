package command

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/enrollment"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CompleteLesson struct {
	UserID   string
	CourseID string
	LessonID string
}

type CompleteLessonHandler decorator.CommandHandler[CompleteLesson]

type completeLessonHandler struct {
	enrollmentRepository enrollment.EnrollmentRepository
}

func NewCompleteLessonHandler(
	enrollmentRepository enrollment.EnrollmentRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CompleteLessonHandler {
	if enrollmentRepository == nil {
		panic("enrollment repository is required")
	}

	return decorator.ApplyCommandDecorators(
		completeLessonHandler{
			enrollmentRepository: enrollmentRepository,
		},
		logger,
		metricsClient,
	)
}

func (h completeLessonHandler) Handle(ctx context.Context, cmd CompleteLesson) error {
	// Validate input
	if cmd.UserID == "" {
		return errors.New("user ID is required")
	}
	if cmd.CourseID == "" {
		return errors.New("course ID is required")
	}
	if cmd.LessonID == "" {
		return errors.New("lesson ID is required")
	}

	// Get enrollment
	enroll, err := h.enrollmentRepository.GetByUserAndCourse(ctx, cmd.UserID, cmd.CourseID)
	if err != nil {
		return errors.Wrap(err, "enrollment not found - user not enrolled in course")
	}

	// Mark lesson as completed
	if err := enroll.CompleteLesson(cmd.LessonID); err != nil {
		return errors.Wrap(err, "failed to complete lesson")
	}

	// Update enrollment
	if err := h.enrollmentRepository.Update(ctx, enroll); err != nil {
		return errors.Wrap(err, "failed to update enrollment")
	}

	return nil
}
