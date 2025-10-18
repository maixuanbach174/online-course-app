package enrollment

import "context"

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *Enrollment) error
	Update(ctx context.Context, enrollment *Enrollment) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Enrollment, error)
	GetAll(ctx context.Context) ([]*Enrollment, error)
	GetByUserAndCourse(ctx context.Context, userID, courseID string) (*Enrollment, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*Enrollment, error)
}
