package refresh

import (
	"context"
	"time"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

type Result struct {
	AccessToken  string
	ExpiresIn    int64
	RefreshToken string
	RefreshTTL   time.Duration
}

type Interactor struct {
	idGen         user.IdentifierGenerator
	tokenProvider user.TokenProvider
	tokenHasher   user.TokenHasher
	refreshRepo   user.RefreshTokenRepository
}

func NewInteractor(idGen user.IdentifierGenerator, tokenProvider user.TokenProvider, tokenHasher user.TokenHasher, refreshRepo user.RefreshTokenRepository) RefreshUseCase {
	return &Interactor{
		idGen:         idGen,
		tokenProvider: tokenProvider,
		tokenHasher:   tokenHasher,
		refreshRepo:   refreshRepo,
	}
}

func (i *Interactor) Refresh(ctx context.Context, cmd Command) (*Result, error) {
	// 1. Hash the incoming raw refresh token
	hashStr := i.tokenHasher.Hash(cmd.RawRefreshToken)

	// 2. Find the refresh token record by hash
	record, err := i.refreshRepo.FindByHash(ctx, hashStr)
	if err != nil {
		return nil, err
	}

	// 3. Check expired
	if record.IsExpired() {
		_ = i.refreshRepo.DeleteByHash(ctx, hashStr) // Cleanup expired token
		return nil, user.ErrTokenExpired
	}

	// 4. Rotation
	if err := i.refreshRepo.DeleteByHash(ctx, hashStr); err != nil {
		return nil, err
	}

	accessToken, err := i.tokenProvider.GenerateAccessToken(record.UserID, record.Role)
	if err != nil {
		return nil, err
	}

	// 5. Generate new refresh token
	rawRefreshToken := i.tokenProvider.GenerateRefreshTokenRaw()
	newRecord := user.RefreshTokenRecord{
		ID:        i.idGen.Generate(),
		UserID:    record.UserID,
		TokenHash: i.tokenHasher.Hash(rawRefreshToken),
		ExpiresAt: time.Now().Add(i.tokenProvider.RefreshTTL()),
		CreatedAt: time.Now(),
	}
	if err := i.refreshRepo.Save(ctx, &newRecord); err != nil {
		return nil, err
	}

	return &Result{
		AccessToken:  accessToken,
		ExpiresIn:    int64(i.tokenProvider.AccessTTL().Seconds()),
		RefreshToken: rawRefreshToken,
		RefreshTTL:   i.tokenProvider.RefreshTTL(),
	}, nil
}
