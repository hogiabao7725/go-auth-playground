package user

type PasswordHasher interface {
	Hash(plainPassword string) (string, error)
	Compare(hashedPassword, plainPassword string) error
}

type IdentifierGenerator interface {
	Generate() string
}

type TokenGenerator interface {
}
