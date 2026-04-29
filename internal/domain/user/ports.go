package user

import "time"

type PasswordHasher interface {
	Hash(plainPassword string) (string, error)
	Compare(hashedPassword, plainPassword string) error
}

type IdentifierGenerator interface {
	Generate() string
}

type TokenProvider interface {
	GenerateAccessToken(userID, role string) (string, error)
	GenerateRefreshToken(userID string) (string, error)

	ParseAccessToken(token string) (*AccessTokenData, error)
	ParseRefreshToken(token string) (*RefreshTokenData, error)

	ValidateAccessToken(token string) bool
	ValidateRefreshToken(token string) bool

	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}
