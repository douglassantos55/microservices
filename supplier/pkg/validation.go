package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(data any) error
}

type validate struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	validator := &validate{validator.New()}
	validator.registerRules()
	return validator
}

func (v *validate) registerRules() {
	v.validator.RegisterValidation("cnpj", func(fl validator.FieldLevel) bool {
		return IsCNPJ(fl.Field().String())
	})
}

func (v *validate) Validate(data any) error {
	err := v.validator.Struct(data)
	if err == nil {
		return err
	}
	return parseErrors(err.(validator.ValidationErrors))
}

type ValidationError struct {
	errors map[string]string
}

func (e ValidationError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func (e ValidationError) Error() string {
	bts, err := json.Marshal(e.errors)
	if err != nil {
		return err.Error()
	}
	return string(bts)
}

func parseErrors(err validator.ValidationErrors) ValidationError {
	errors := make(map[string]string)
	for _, error := range err {
		errors[error.Namespace()] = getMessage(error.Tag())
	}
	return ValidationError{errors}
}

func getMessage(tag string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "this field must be a valid email"
	case "cnpj":
		return "this field must be a valid cnpj"
	default:
		return "something wrong is not right with this field"
	}
}
