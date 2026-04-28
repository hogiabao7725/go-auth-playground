package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user/vo"
	"github.com/hogiabao7725/go-auth-playground/internal/infrastructure/persistence/sqlc"
	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("infrastructure.persistence.user_repo.FindByEmail: %w", err)
	}
	return toDomainUser(&dbUser), nil
}

func toDomainUser(dbUser *sqlc.User) *user.User {
	name := vo.ReconstituteName(dbUser.Name)
	email := vo.ReconstituteEmail(dbUser.Email)
	hashedPwd := vo.ReconstituteHashedPassword(dbUser.Password)

	role := vo.ReconstituteRole(dbUser.Role)

	return user.ReconstructUser(
		dbUser.ID,
		name,
		email,
		hashedPwd,
		role,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)
}
