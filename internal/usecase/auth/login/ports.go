package login

import (
	"context"
)

type LoginUseCase interface {
	Login(ctx context.Context, cmd Command) (*Result, error)
}
