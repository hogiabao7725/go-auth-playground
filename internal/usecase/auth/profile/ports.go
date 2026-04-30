package profile

import "context"

type ProfileUseCase interface {
	GetProfile(ctx context.Context, cmd Command) (*ProfileInfo, error)
}
