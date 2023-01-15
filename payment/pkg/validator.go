package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type Validator interface {
	Validate(any) error
}

type ValidationRule interface {
	Tag() string
	Valid(string) bool
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
	rules     []ValidationRule
	validator *validator.Validate
}

func NewValidator(rules []ValidationRule) *validate {
	validate := &validate{rules, validator.New()}
	validate.registerRules()

	return validate
}

func (v *validate) registerRules() {
	for _, rule := range v.rules {
		v.validator.RegisterValidation(rule.Tag(), func(fl validator.FieldLevel) bool {
			return rule.Valid(fl.Field().String())
		})
	}
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
	case "paymenttype":
		return "invalid payment type"
	default:
		return "something is not right about this field"
	}
}

type PaymentTypeSource interface {
	GetPaymentType(string) (*Type, error)
}

type PaymentTypeRule struct {
	service PaymentTypeSource
}

func NewPaymentTypeRule(svc PaymentTypeSource) *PaymentTypeRule {
	return &PaymentTypeRule{svc}
}

func (r PaymentTypeRule) Tag() string {
	return "paymenttype"
}

func (r PaymentTypeRule) Valid(value string) bool {
	_, err := r.service.GetPaymentType(value)
	return err == nil
}
