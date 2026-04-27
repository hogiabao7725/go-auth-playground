package identifier

import (
	"github.com/google/uuid"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

var _ user.IdentifierGenerator = (*UUID)(nil)

type UUID struct{}

func (UUID) Generate() string {
	return uuid.New().String()
}
