package user

import (
	"time"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user/vo"
)

type User struct {
	id        string
	name      vo.Name
	email     vo.Email
	password  vo.HashedPassword
	role      vo.Role
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id string, name vo.Name, email vo.Email, hashedPassword vo.HashedPassword, role vo.Role) (*User, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	now := time.Now()

	return &User{
		id:        id,
		name:      name,
		email:     email,
		password:  hashedPassword,
		role:      role,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructUser(id string, name vo.Name, email vo.Email, hashedPassword vo.HashedPassword, role vo.Role, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		password:  hashedPassword,
		role:      role,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() string { return u.id }

func (u *User) Name() vo.Name { return u.name }

func (u *User) Email() vo.Email { return u.email }

func (u *User) Role() vo.Role { return u.role }

func (u *User) PasswordHash() string { return u.password.Value() }

func (u *User) CreatedAt() time.Time { return u.createdAt }

func (u *User) UpdatedAt() time.Time { return u.updatedAt }
