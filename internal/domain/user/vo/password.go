package vo

import "errors"

var (
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrWeakPassword  = errors.New("password is too weak")
)

type PlainPassword struct {
	value string
}

func NewPlainPassword(raw string) (PlainPassword, error) {
	if raw == "" {
		return PlainPassword{}, ErrEmptyPassword
	}
	if len(raw) < 6 {
		return PlainPassword{}, ErrWeakPassword
	}
	return PlainPassword{value: raw}, nil
}

func (pp PlainPassword) Value() string {
	return pp.value
}

// --------------------------------------------------------------------------

type HashedPassword struct {
	value string
}

func NewHashedPassword(hash string) HashedPassword {
	return HashedPassword{value: hash}
}

func ReconstituteHashedPassword(hash string) HashedPassword {
	return HashedPassword{value: hash}
}

func (p HashedPassword) Value() string {
	return p.value
}

func (p HashedPassword) Equal(other HashedPassword) bool {
	return p == other
}
