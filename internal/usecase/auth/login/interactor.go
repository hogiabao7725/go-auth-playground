package login

import (
	"context"
	"errors"

	userDomain "github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user/vo"
)

type Interactor struct {
	passHasher userDomain.PasswordHasher
	repoUser   userDomain.UserRepository
}

func NewInteractor(passHasher userDomain.PasswordHasher, repoUser userDomain.UserRepository) *Interactor {
	return &Interactor{
		passHasher: passHasher,
		repoUser:   repoUser,
	}
}

func (i *Interactor) Execute(ctx context.Context, cmd Command) (*userDomain.User, error) {
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
		return nil, userDomain.ErrInvalidCredentials
	}

	return user, nil
}
