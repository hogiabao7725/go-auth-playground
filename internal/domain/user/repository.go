package user

import "context"

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, record *RefreshTokenRecord) error
	FindByHash(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error)
	DeleteByHash(ctx context.Context, tokenHash string) error
}
