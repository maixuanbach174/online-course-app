package command

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/enrollment"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/user"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EnrollInCourse struct {
	EnrollmentID string
	UserID       string
	CourseID     string
}

type EnrollInCourseHandler decorator.CommandHandler[EnrollInCourse]

type enrollInCourseHandler struct {
	enrollmentRepository enrollment.EnrollmentRepository
	userRepository       user.UserRepository
	courseRepository     course.CourseRepository
}

func NewEnrollInCourseHandler(
	enrollmentRepository enrollment.EnrollmentRepository,
	userRepository user.UserRepository,
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) EnrollInCourseHandler {
	if enrollmentRepository == nil {
		panic("enrollment repository is required")
	}
	if userRepository == nil {
		panic("user repository is required")
	}
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyCommandDecorators(
		enrollInCourseHandler{
			enrollmentRepository: enrollmentRepository,
			userRepository:       userRepository,
			courseRepository:     courseRepository,
		},
		logger,
		metricsClient,
	)
}

func (h enrollInCourseHandler) Handle(ctx context.Context, cmd EnrollInCourse) error {
	// Validate input
	if cmd.EnrollmentID == "" {
		return errors.New("enrollment ID is required")
	}
	if cmd.UserID == "" {
		return errors.New("user ID is required")
	}
	if cmd.CourseID == "" {
		return errors.New("course ID is required")
	}

	// Verify user exists and can enroll
	student, err := h.userRepository.Get(ctx, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "user not found")
	}
	if !student.CanEnroll() {
		return errors.New("only students can enroll in courses")
	}

	// Verify course exists
	_, err = h.courseRepository.Get(ctx, cmd.CourseID)
	if err != nil {
		return errors.Wrap(err, "course not found")
	}

	// Create enrollment entity
	newEnrollment, err := enrollment.NewEnrollment(cmd.EnrollmentID, cmd.UserID, cmd.CourseID)
	if err != nil {
		return errors.Wrap(err, "failed to create enrollment")
	}

	// Persist to repository
	if err := h.enrollmentRepository.Create(ctx, newEnrollment); err != nil {
		return errors.Wrap(err, "failed to save enrollment")
	}

	return nil
}
