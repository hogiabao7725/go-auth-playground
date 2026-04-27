package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/persistence/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ user.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Save(ctx context.Context, newUser *user.User) error {
	params := sqlc.CreateUserParams{
		ID:        newUser.ID(),
		Name:      newUser.Name().String(),
		Email:     newUser.Email().String(),
		Password:  newUser.PasswordHash(),
		Role:      newUser.Role().String(),
		CreatedAt: newUser.CreatedAt(),
		UpdatedAt: newUser.UpdatedAt(),
	}

	err := r.queries.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return user.ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("infrastructure.persistence.user_repo.Save: %w", err)
	}
	return nil
}
