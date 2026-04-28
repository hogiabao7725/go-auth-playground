package crypt

import (
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var _ user.PasswordHasher = (*Bcrypt)(nil)

var (
	ErrMismatched = errors.New("password does not match")
)

type Bcrypt struct{}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (Bcrypt) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("infrastructure.crypt.bcrypt.Hash: %w", err)
	}
	return string(bytes), nil
}

func (Bcrypt) Compare(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrMismatched
		}
		return fmt.Errorf("infrastructure.crypt.bcrypt.Compare: %w", err)
	}
	return nil
}
