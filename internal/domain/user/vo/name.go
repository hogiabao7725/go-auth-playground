package vo

import (
	"errors"
	"strings"
)

var (
	ErrEmptyName = errors.New("name cannot be empty")
)

type Name struct {
	value string
}

func NewName(raw string) (Name, error) {
	normalizedName, err := validateAndNormalizeName(raw)
	if err != nil {
		return Name{}, err
	}
	return Name{value: normalizedName}, nil
}

func ReconstituteName(raw string) Name {
	return Name{value: raw}
}

func (n Name) String() string {
	return n.value
}

func (n Name) Equal(other Name) bool {
	return n == other
}

func validateAndNormalizeName(name string) (string, error) {
	normalize := strings.Join(strings.Fields(name), " ")
	if normalize == "" {
		return "", ErrEmptyName
	}
	return normalize, nil
}
