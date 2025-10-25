package course_query

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/course"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CourseByTeacherQuery struct {
	TeacherID string
}

type CourseByTeacherHandler decorator.QueryHandler[CourseByTeacherQuery, []*course.Course]

type courseByTeacherHandler struct {
	courseRepository course.CourseRepository
}

func NewCourseByTeacherHandler(
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CourseByTeacherHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyQueryDecorators(
		courseByTeacherHandler{
			courseRepository: courseRepository,
		},
		logger,
		metricsClient,
	)
}

func (h courseByTeacherHandler) Handle(ctx context.Context, query CourseByTeacherQuery) ([]*course.Course, error) {
	if query.TeacherID == "" {
		return nil, errors.New("teacher ID is required")
	}
	return h.courseRepository.GetAllByTeacherID(ctx, query.TeacherID)
}
