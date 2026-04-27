package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidJSON    = errors.New("invalid JSON")
	ErrInvalidRequest = errors.New("invalid request")
)

var v *validator.Validate

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation error"
}

func init() {
	v = validator.New()

	// Register custom tag name function to use JSON tags in error messages
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		// eg: `json:"email,omitempty"` -> "email,omitempty" → ["email", "omitempty"]
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validation "not_blank"
	v.RegisterValidation("not_blank", func(field validator.FieldLevel) bool {
		return strings.TrimSpace(field.Field().String()) != ""
	})
}

func BindJSON(w http.ResponseWriter, r *http.Request, req any) error {
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return ErrInvalidJSON
	}

	// Validate struct
	if err := v.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string, len(ve))
			for _, fe := range ve {
				fields[fe.Field()] = validationMessage(fe)
			}
			return &ValidationError{Fields: fields}
		}
		return ErrInvalidRequest
	}

	return nil
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters long"
	case "not_blank":
		return "cannot be blank"
	default:
		return "is not valid"
	}
}
