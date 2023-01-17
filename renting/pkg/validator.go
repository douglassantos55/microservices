package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"google.golang.org/grpc"
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
		validationFunc := func(rule ValidationRule) validator.Func {
			return func(fl validator.FieldLevel) bool {
				return rule.Valid(fl.Field().String())
			}
		}(rule)
		v.validator.RegisterValidation(rule.Tag(), validationFunc)
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
	case "payment_type":
		return "invalid payment type"
	default:
		return "something is not right about this field"
	}
}

type PaymentMethodRule struct{}

func (r PaymentMethodRule) Tag() string {
	return "payment_method"
}

func (r PaymentMethodRule) Valid(value string) bool {
	return false
}

type paymentTypeRule struct {
	cc *grpc.ClientConn
}

func NewPaymentTypeRule(cc *grpc.ClientConn) paymentTypeRule {
	return paymentTypeRule{cc}
}

func (r paymentTypeRule) Tag() string {
	return "payment_type"
}

func (r paymentTypeRule) Valid(value string) bool {
	endpoint := getPaymentTypeEndpoint(r.cc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	_, err := endpoint(ctx, value)

	return err == nil
}

type PaymentConditionRule struct{}

func (r PaymentConditionRule) Tag() string {
	return "payment_condition"
}

func (r PaymentConditionRule) Valid(value string) bool {
	return false
}

type CustomerRule struct{}

func (r CustomerRule) Tag() string {
	return "customer"
}

func (r CustomerRule) Valid(value string) bool {
	return false
}
