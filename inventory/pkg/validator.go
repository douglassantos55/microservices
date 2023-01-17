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
	Validate(data any) error
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

func (v ValidationError) Error() string {
	data, err := json.Marshal(v.errors)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

type validate struct {
	rules     []ValidationRule
	validator *validator.Validate
}

func NewValidator(rules []ValidationRule) Validator {
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
	if err != nil {
		return makeValidationError(err.(validator.ValidationErrors))
	}
	return nil
}

func makeValidationError(validationErrors validator.ValidationErrors) error {
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
	case "supplier":
		return "invalid supplier"
	default:
		return "something is not right about this field"
	}
}

type supplierRule struct {
	cc *grpc.ClientConn
}

func NewSupplierRule(cc *grpc.ClientConn) *supplierRule {
	return &supplierRule{cc}
}

func (r supplierRule) Tag() string {
	return "supplier"
}

func (r supplierRule) Valid(value string) bool {
	endpoint := makeFetchSupplierEndpoint(r.cc)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	_, err := endpoint(ctx, value)

	return err == nil
}
