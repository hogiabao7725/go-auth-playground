package crypt

import (
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var _ user.PasswordHasher = (*Bcrypt)(nil)

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
