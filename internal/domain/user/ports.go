package user

type PasswordHasher interface {
	Hash(plainPassword string) (string, error)
}

type IdentifierGenerator interface {
	Generate() string
}
