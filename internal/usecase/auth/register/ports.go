package register

import (
	"context"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, cmd Command) (*user.User, error)
}
