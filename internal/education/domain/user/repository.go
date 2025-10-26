package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
}
