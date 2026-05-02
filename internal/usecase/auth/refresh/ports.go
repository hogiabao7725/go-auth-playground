package refresh

import "context"

type RefreshUseCase interface {
	Refresh(ctx context.Context, cmd Command) (*Result, error)
}
