package pkg

import (
	"context"
	"encoding/json"
	"fmt"
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
		errors[error.Namespace()] = getErrorMessage(error)
	}
	return ValidationError{errors}
}

func getErrorMessage(error validator.FieldError) string {
	switch error.Tag() {
	case "required":
		return "this field is required"
	case "numeric", "number":
		return "this field contain a number"
	case "ltecsfield":
		return fmt.Sprintf("this field must be less than or equals to %s", error.Param())
	case "payment_type":
		return "invalid payment type"
	case "payment_condition":
		return "invalid payment condition"
	case "payment_method":
		return "invalid payment method"
	case "customer":
		return "invalid customer"
	case "equipment":
		return "invalid equipment"
	default:
		return "something is not right about this field"
	}
}

type paymentMethodRule struct {
	cc *grpc.ClientConn
}

func NewPaymentMethodRule(cc *grpc.ClientConn) paymentMethodRule {
	return paymentMethodRule{cc}
}

func (r paymentMethodRule) Tag() string {
	return "payment_method"
}

func (r paymentMethodRule) Valid(value string) bool {
	endpoint := getPaymentMethodEndpoint(r.cc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	_, err := endpoint(ctx, value)

	return err == nil
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

type paymentConditionRule struct {
	cc *grpc.ClientConn
}

func NewPaymentConditionRule(cc *grpc.ClientConn) paymentConditionRule {
	return paymentConditionRule{cc}
}

func (r paymentConditionRule) Tag() string {
	return "payment_condition"
}

func (r paymentConditionRule) Valid(value string) bool {
	endpoint := getPaymentConditionEndpoint(r.cc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	_, err := endpoint(ctx, value)

	return err == nil
}

type customerRule struct {
	cc *grpc.ClientConn
}

func NewCustomerRule(cc *grpc.ClientConn) customerRule {
	return customerRule{cc}
}

func (r customerRule) Tag() string {
	return "customer"
}

func (r customerRule) Valid(value string) bool {
	endpoint := getCustomerEndpoint(r.cc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	_, err := endpoint(ctx, value)

	return err == nil
}
