package profile

import (
	"context"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

var _ ProfileUseCase = (*Interactor)(nil)

type ProfileInfo struct {
	UserID   string
	Username string
	Email    string
	Role     string
}

type Interactor struct {
	userRepo user.UserRepository
}

func NewInteractor(userRepo user.UserRepository) ProfileUseCase {
	return &Interactor{
		userRepo: userRepo,
	}
}

func (i *Interactor) GetProfile(ctx context.Context, cmd Command) (*ProfileInfo, error) {
	user, err := i.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("usecase.profile.Interactor.GetProfile.FindByID: %w", err)
	}

	return &ProfileInfo{
		UserID:   user.ID(),
		Username: user.Name().String(),
		Email:    user.Email().String(),
		Role:     user.Role().String(),
	}, nil
}
