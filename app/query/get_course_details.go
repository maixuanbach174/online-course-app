package query

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GetCourseDetails struct {
	CourseID string
}

type GetCourseDetailsHandler decorator.QueryHandler[GetCourseDetails, *course.Course]

type getCourseDetailsHandler struct {
	courseRepository course.CourseRepository
}

func NewGetCourseDetailsHandler(
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetCourseDetailsHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyQueryDecorators(
		getCourseDetailsHandler{
			courseRepository: courseRepository,
		},
		logger,
		metricsClient,
	)
}

func (h getCourseDetailsHandler) Handle(ctx context.Context, query GetCourseDetails) (*course.Course, error) {
	if query.CourseID == "" {
		return nil, errors.New("course ID is required")
	}

	course, err := h.courseRepository.Get(ctx, query.CourseID)
	if err != nil {
		return nil, errors.Wrap(err, "course not found")
	}

	return course, nil
}
