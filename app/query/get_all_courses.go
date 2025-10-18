package query

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/core/domain/course"
	"github.com/sirupsen/logrus"
)

type GetAllCourses struct {
	// Add filtering options if needed
	Domain string // optional filter
	Level  string // optional filter
	Tag    string // optional filter
}

type GetAllCoursesHandler decorator.QueryHandler[GetAllCourses, []*course.Course]

type getAllCoursesHandler struct {
	courseRepository course.CourseRepository
}

func NewGetAllCoursesHandler(
	courseRepository course.CourseRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetAllCoursesHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}

	return decorator.ApplyQueryDecorators(
		getAllCoursesHandler{
			courseRepository: courseRepository,
		},
		logger,
		metricsClient,
	)
}

// !!! Change in the future
func (h getAllCoursesHandler) Handle(ctx context.Context, query GetAllCourses) ([]*course.Course, error) {
	// For MVP, return all courses
	// In production, you'd implement filtering logic here or in repository
	courses, err := h.courseRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Optional: filter courses based on query parameters
	if query.Domain != "" || query.Level != "" || query.Tag != "" {
		filtered := make([]*course.Course, 0)
		for _, c := range courses {
			if h.matchesFilters(c, query) {
				filtered = append(filtered, c)
			}
		}
		return filtered, nil
	}

	return courses, nil
}

func (h getAllCoursesHandler) matchesFilters(c *course.Course, query GetAllCourses) bool {
	if query.Domain != "" && c.Domain().String() != query.Domain {
		return false
	}
	if query.Level != "" && c.Level().String() != query.Level {
		return false
	}
	if query.Tag != "" {
		tag, err := course.NewTagFromString(query.Tag)
		if err != nil || !c.HasTag(tag) {
			return false
		}
	}
	return true
}
