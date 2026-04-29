package user

import "time"

type AccessTokenData struct {
	UserID    string
	Role      string
	ExpiresAt time.Time
	TokenID   string
}

type RefreshTokenData struct {
	UserID    string
	ExpiresAt time.Time
	TokenID   string
}
