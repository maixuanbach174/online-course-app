package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql/database"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/user"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:      db,
		queries: database.New(db),
	}
}

// Create implements user.UserRepository
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	profile := pgtype.Text{
		String: u.Profile(),
		Valid:  u.Profile() != "",
	}

	params := database.CreateUserParams{
		ID:       u.ID(),
		Username: u.Username(),
		Email:    u.Email(),
		Role:     u.Role().String(),
		Profile:  profile,
	}

	if err := r.queries.CreateUser(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

// Update implements user.UserRepository
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	profile := pgtype.Text{
		String: u.Profile(),
		Valid:  u.Profile() != "",
	}

	params := database.UpdateUserParams{
		ID:       u.ID(),
		Username: u.Username(),
		Email:    u.Email(),
		Role:     u.Role().String(),
		Profile:  profile,
	}

	if err := r.queries.UpdateUser(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

// Delete implements user.UserRepository
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}

// Get implements user.UserRepository
func (r *UserRepository) Get(ctx context.Context, id string) (*user.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return r.toDomainUser(dbUser)
}

// GetAll implements user.UserRepository
func (r *UserRepository) GetAll(ctx context.Context) ([]*user.User, error) {
	dbUsers, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}

	users := make([]*user.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		domainUser, err := r.toDomainUser(dbUser)
		if err != nil {
			return nil, err
		}
		users = append(users, domainUser)
	}

	return users, nil
}

// toDomainUser converts database.User to domain user.User
func (r *UserRepository) toDomainUser(dbUser database.User) (*user.User, error) {
	role, err := user.NewRoleFromString(dbUser.Role)
	if err != nil {
		return nil, errors.Wrap(err, "invalid role")
	}

	profile := ""
	if dbUser.Profile.Valid {
		profile = dbUser.Profile.String
	}

	domainUser, err := user.NewUser(
		dbUser.ID,
		dbUser.Username,
		dbUser.Email,
		role,
		profile,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create domain user")
	}

	return domainUser, nil
}
