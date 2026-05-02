package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/persistence/sqlc"
	"github.com/jackc/pgx/v5"
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
		Role:      record.Role,
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
	dbRecord, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("infrastructure.persistence.refresh_token_repo.FindByHash: %w", err)
	}
	return toDomainRefreshTokenRecord(&dbRecord), nil
}

func (r *RefreshTokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	if err := r.queries.DeleteRefreshTokenByTokenHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("infrastructure.persistence.refresh_token_repo.DeleteByHash: %w", err)
	}
	return nil
}

func toDomainRefreshTokenRecord(dbRecord *sqlc.GetRefreshTokenByHashRow) *user.RefreshTokenRecord {
	return &user.RefreshTokenRecord{
		ID:        dbRecord.ID,
		UserID:    dbRecord.UserID,
		TokenHash: dbRecord.TokenHash,
		ExpiresAt: dbRecord.ExpiresAt,
		CreatedAt: dbRecord.CreatedAt,
		Role:      dbRecord.Role,
	}
}
