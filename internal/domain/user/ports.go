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
	// Access Token (JWT)
	GenerateAccessToken(userID, role string) (string, error)
	ParseAccessToken(token string) (*AccessTokenData, error)
	ValidateAccessToken(token string) bool
	AccessTTL() time.Duration

	// Refresh Token (Opaque)
	GenerateRefreshTokenRaw() string
	RefreshTTL() time.Duration
}

type TokenHasher interface {
	Hash(plainToken string) string
	Compare(hashedToken, plainToken string) error
}
