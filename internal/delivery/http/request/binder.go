package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hogiabao7725/go-auth-playground/internal/delivery/http/response"
)

var v *validator.Validate

func init() {
	v = validator.New()

	// Register custom tag name function to use JSON tags in error messages
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
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

func BindJSON(w http.ResponseWriter, r *http.Request, req any) bool {
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.BadRequest(w, "invalid JSON", nil)
		return false
	}

	// Validate struct
	if err := v.Struct(req); err != nil {
		handleValidationError(w, err)
		return false
	}

	return true
}

func handleValidationError(w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve))
		for _, fe := range ve {
			fields[fe.Field()] = validationMessage(fe)
		}

		response.ValidationError(w, fields)
		return
	}
	response.BadRequest(w, "invalid request", nil)
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
