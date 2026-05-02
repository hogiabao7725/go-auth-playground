package login

import (
	"context"
	"errors"
	"time"

	userDomain "github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user/vo"
)

type Result struct {
	AccessToken  string
	ExpiresIn    time.Time
	RefreshToken string
	User         *userDomain.User
}

type Interactor struct {
	passHasher        userDomain.PasswordHasher
	repoUser          userDomain.UserRepository
	tokenProvider     userDomain.TokenProvider
	tokenHasher       userDomain.TokenHasher
	idGen             userDomain.IdentifierGenerator
	refreshRepository userDomain.RefreshTokenRepository
}

func NewInteractor(
	passHasher userDomain.PasswordHasher,
	repoUser userDomain.UserRepository,
	tokenProvider userDomain.TokenProvider,
	tokenHasher userDomain.TokenHasher,
	idGen userDomain.IdentifierGenerator,
	refreshRepository userDomain.RefreshTokenRepository) LoginUseCase {
	return &Interactor{
		passHasher:        passHasher,
		repoUser:          repoUser,
		tokenProvider:     tokenProvider,
		tokenHasher:       tokenHasher,
		idGen:             idGen,
		refreshRepository: refreshRepository,
	}
}

func (i *Interactor) Login(ctx context.Context, cmd Command) (*Result, error) {
	email, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	pass, err := vo.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
	}

	user, err := i.repoUser.FindByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, userDomain.ErrUserNotFound) {
			return nil, userDomain.ErrInvalidCredentials
		}
		return nil, err
	}

	// compare password
	if err := i.passHasher.Compare(user.PasswordHash(), pass.Value()); err != nil {
		return nil, err
	}

	// Generate access token (JWT)
	accessToken, err := i.tokenProvider.GenerateAccessToken(user.ID(), user.Role().String())
	if err != nil {
		return nil, err
	}

	// Generate and save refresh token
	rawRefreshToken := i.tokenProvider.GenerateRefreshTokenRaw() // Random string for hashing
	refreshTokenData := userDomain.RefreshTokenRecord{
		ID:        i.idGen.Generate(),
		UserID:    user.ID(),
		TokenHash: i.tokenHasher.Hash(rawRefreshToken),
		ExpiresAt: time.Now().Add(i.tokenProvider.RefreshTTL()),
		CreatedAt: time.Now(),
	}

	if err := i.refreshRepository.Save(ctx, &refreshTokenData); err != nil {
		return nil, err
	}

	return &Result{
		AccessToken:  accessToken,
		ExpiresIn:    time.Now().Add(i.tokenProvider.AccessTTL()),
		RefreshToken: rawRefreshToken,
		User:         user,
	}, nil
}
