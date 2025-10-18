package command

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/core/domain/user"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RegisterUser struct {
	UserID   string
	Username string
	Email    string
	Role     string
	Profile  string
}

type RegisterUserHandler decorator.CommandHandler[RegisterUser]

type registerUserHandler struct {
	userRepository user.UserRepository
}

func NewRegisterUserHandler(
	userRepository user.UserRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RegisterUserHandler {
	if userRepository == nil {
		panic("user repository is required")
	}

	return decorator.ApplyCommandDecorators(
		registerUserHandler{
			userRepository: userRepository,
		},
		logger,
		metricsClient,
	)
}

func (h registerUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	// Validate input
	if cmd.UserID == "" {
		return errors.New("user ID is required")
	}
	if cmd.Username == "" {
		return errors.New("username is required")
	}
	if cmd.Email == "" {
		return errors.New("email is required")
	}

	// Parse role
	role, err := user.NewRoleFromString(cmd.Role)
	if err != nil {
		return errors.Wrap(err, "invalid role")
	}

	// Create user entity
	newUser, err := user.NewUser(cmd.UserID, cmd.Username, cmd.Email, role, cmd.Profile)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	// Persist to repository
	if err := h.userRepository.Create(ctx, newUser); err != nil {
		return errors.Wrap(err, "failed to save user")
	}

	return nil
}
