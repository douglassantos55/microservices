package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type Validator interface {
	Validate(any) error
}

type ValidationError struct {
	errors map[string]string
}

func (v ValidationError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func (e ValidationError) Error() string {
	b, err := json.Marshal(e.errors)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type validate struct {
	validator *validator.Validate
}

func NewValidator() *validate {
	return &validate{validator.New()}
}

func (v *validate) Validate(data any) error {
	err := v.validator.Struct(data)
	if err == nil {
		return nil
	}
	return makeValidationError(err.(validator.ValidationErrors))
}

func makeValidationError(validationErrors validator.ValidationErrors) ValidationError {
	errors := make(map[string]string)
	for _, error := range validationErrors {
		errors[error.Namespace()] = getErrorMessage(error.Tag())
	}
	return ValidationError{errors}
}

func getErrorMessage(error string) string {
	switch error {
	case "required":
		return "this field is required"
	case "numeric", "number":
		return "this field contain a number"
	default:
		return "something is not right about this field"
	}
}
