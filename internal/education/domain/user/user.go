package user

import "github.com/pkg/errors"

type User struct {
	id       string
	username string
	email    string
	role     Role
	profile  string
}

func NewUser(id string, username string, email string, role Role, profile string) (*User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if profile == "" && role == RoleTeacher {
		return nil, errors.New("profile is required")
	}
	return &User{
		id:       id,
		username: username,
		email:    email,
		role:     role,
		profile:  profile,
	}, nil
}

func (u *User) UpdateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}
	u.username = username
	return nil
}

func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	u.email = email
	return nil
}

func (u *User) UpdateProfile(profile string) error {
	if profile == "" && u.role == RoleTeacher {
		return errors.New("teacher profile is required")
	}
	u.profile = profile
	return nil
}

// Getters (read-only access for serialization/display)
func (u *User) ID() string       { return u.id }
func (u *User) Username() string { return u.username }
func (u *User) Email() string    { return u.email }
func (u *User) Role() Role       { return u.role }
func (u *User) Profile() string  { return u.profile }

// Behavior methods
func (u *User) HasRole(role Role) bool {
	return u.role == role
}

func (u *User) CanTeach() bool {
	return u.role == RoleTeacher || u.role == RoleAdmin
}

func (u *User) CanEnroll() bool {
	return u.role == RoleStudent
}
