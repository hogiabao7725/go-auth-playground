package persistence

import (
	"context"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/persistence/sqlc"
)

var _ user.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

type RefreshTokenRepository struct {
	queries *sqlc.Queries
}

func NewRefreshTokenRepository(queries *sqlc.Queries) *RefreshTokenRepository {
	return &RefreshTokenRepository{queries: queries}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, record *user.RefreshTokenRecord) error {
	params := sqlc.CreateRefreshTokenParams{
		ID:        record.ID,
		UserID:    record.UserID,
		TokenHash: record.TokenHash,
		ExpiresAt: record.ExpiresAt,
		CreatedAt: record.CreatedAt,
	}
	if err := r.queries.CreateRefreshToken(ctx, params); err != nil {
		return fmt.Errorf("infrastructure.persistence.refresh_token_repo.Save: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*user.RefreshTokenRecord, error) {
	return nil, nil
}

func (r *RefreshTokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	return nil
}
