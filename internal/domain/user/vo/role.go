package vo

import (
	"errors"
	"strings"
)

var (
	ErrInvalidRole = errors.New("invalid role")
)

// Possible values: "user", "organizer", "admin"
type Role struct {
	value string
}

// Predefined roles
var (
	RoleUser      = Role{"user"}
	RoleOrganizer = Role{"organizer"}
	RoleAdmin     = Role{"admin"}
)

func NewRole(value string) (Role, error) {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	switch trimmed {
	case "user", "organizer", "admin":
		return Role{value: trimmed}, nil
	default:
		return Role{}, ErrInvalidRole
	}
}

func ReconstituteRole(value string) Role {
	return Role{value: value}
}

func (r Role) String() string {
	return r.value
}

func (r Role) Equal(other Role) bool {
	return r == other
}

func (r Role) IsAdmin() bool {
	return r.value == "admin"
}

func (r Role) IsOrganizer() bool {
	return r.value == "organizer" || r.value == "admin"
}
