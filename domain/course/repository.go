package course

import "context"

type CourseRepository interface {
	Create(ctx context.Context, course *Course) error
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Course, error)
	GetAll(ctx context.Context) ([]*Course, error)
}
