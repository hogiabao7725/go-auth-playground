package user

import "time"

// Access Token (JWT)

type AccessTokenData struct {
	UserID    string
	Role      string
	ExpiresAt time.Time
	TokenID   string
}

// Refresh Token (Opaque)

type RefreshTokenRecord struct {
	ID        string
	UserID    string
	Role      string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (rt *RefreshTokenRecord) Validate() error {
	if time.Now().After(rt.ExpiresAt) {
		return ErrTokenExpired
	}
	return nil
}

func (rt *RefreshTokenRecord) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}
