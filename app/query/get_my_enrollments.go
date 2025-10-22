package query

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/enrollment"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GetMyEnrollments struct {
	UserID string
}

type GetMyEnrollmentsHandler decorator.QueryHandler[GetMyEnrollments, []*enrollment.Enrollment]

type getMyEnrollmentsHandler struct {
	enrollmentRepository enrollment.EnrollmentRepository
}

func NewGetMyEnrollmentsHandler(
	enrollmentRepository enrollment.EnrollmentRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetMyEnrollmentsHandler {
	if enrollmentRepository == nil {
		panic("enrollment repository is required")
	}

	return decorator.ApplyQueryDecorators(
		getMyEnrollmentsHandler{
			enrollmentRepository: enrollmentRepository,
		},
		logger,
		metricsClient,
	)
}

func (h getMyEnrollmentsHandler) Handle(ctx context.Context, query GetMyEnrollments) ([]*enrollment.Enrollment, error) {
	if query.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	enrollments, err := h.enrollmentRepository.GetAllByUserID(ctx, query.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get enrollments")
	}

	return enrollments, nil
}
