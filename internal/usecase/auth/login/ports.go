package login

import (
	"context"

	userDomain "github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

type LoginUseCase interface {
	Login(ctx context.Context, cmd Command) (*userDomain.User, error)
}
