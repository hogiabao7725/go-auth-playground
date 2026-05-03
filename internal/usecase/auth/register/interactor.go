package register

import (
	"context"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user/vo"
)

type Interactor struct {
	passHasher user.PasswordHasher
	idGen      user.IdentifierGenerator
	repoUser   user.UserRepository
}

func NewInteractor(passHasher user.PasswordHasher, idGen user.IdentifierGenerator, repoUser user.UserRepository) RegisterUseCase {
	return &Interactor{
		passHasher: passHasher,
		idGen:      idGen,
		repoUser:   repoUser,
	}
}

func (i *Interactor) Register(ctx context.Context, cmd Command) (*user.User, error) {
	name, err := vo.NewName(cmd.Name)
	if err != nil {
		return nil, err
	}

	email, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	plainPass, err := vo.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	hashStr, err := i.passHasher.Hash(plainPass.Value())
	if err != nil {
		return nil, fmt.Errorf("usecase.register.Interactor.Register.HashPassword: %w", err)
	}
	hashedPass := vo.NewHashedPassword(hashStr)

	id := i.idGen.Generate()

	userSave, err := user.NewUser(id, name, email, hashedPass, vo.RoleUser)
	if err != nil {
		return nil, err
	}

	if err := i.repoUser.Save(ctx, userSave); err != nil {
		return nil, fmt.Errorf("usecase.register.Interactor.Register.SaveUser: %w", err)
	}

	return userSave, nil
}
