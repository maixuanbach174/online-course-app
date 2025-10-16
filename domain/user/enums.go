package user

import "github.com/pkg/errors"

var (
	RoleAdmin   = Role{r: "admin"}
	RoleStudent = Role{r: "student"}
	RoleTeacher = Role{r: "teacher"}
)

var roleValues = []Role{
	RoleAdmin,
	RoleStudent,
	RoleTeacher,
}

type Role struct {
	r string
}

func (r Role) NewRoleFromString(roleStr string) (Role, error) {
	for _, role := range roleValues {
		if role.String() == roleStr {
			return role, nil
		}
	}
	return Role{}, errors.Errorf("unknown '%s' role", roleStr)
}

func (r Role) String() string {
	return r.r
}
