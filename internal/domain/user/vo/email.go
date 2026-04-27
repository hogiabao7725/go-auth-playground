package vo

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrEmptyEmail   = errors.New("email cannot be empty")
	ErrInvalidEmail = errors.New("invalid email format")
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	emailStr, err := validateAndNormalizeEmail(raw)
	if err != nil {
		return Email{}, err
	}
	return Email{value: emailStr}, nil
}

func ReconstituteEmail(raw string) Email {
	return Email{value: raw}
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equal(other Email) bool {
	return e == other
}

func validateAndNormalizeEmail(raw string) (string, error) {
	emailStr := strings.TrimSpace(raw)
	emailStr = strings.ToLower(emailStr)

	if emailStr == "" {
		return "", ErrEmptyEmail
	}

	addrEmail, err := mail.ParseAddress(emailStr)
	if err != nil {
		return "", ErrInvalidEmail
	}
	return addrEmail.Address, nil
}
