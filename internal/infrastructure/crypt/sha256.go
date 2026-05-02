package crypt

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

var _ user.TokenHasher = (*SHA256)(nil)

var (
	ErrTokenMismatch = errors.New("crypt.sha256: token mismatch")
)

type SHA256 struct{}

func NewSHA256() *SHA256 {
	return &SHA256{}
}

func (s SHA256) Hash(plainToken string) string {
	sum := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(sum[:])
}

func (s SHA256) Compare(hashedToken, token string) error {
	calculated := s.Hash(token)
	if subtle.ConstantTimeCompare([]byte(hashedToken), []byte(calculated)) != 1 {
		return ErrTokenMismatch
	}
	return nil
}
