package pkg

import "net/http"

type PaymentMethod struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" validate:"required"`
	AccountID string `json:"account_id" validate:"required"`
}

type Service interface {
	CreatePaymentMethod(PaymentMethod) (*PaymentMethod, error)
}

type service struct {
	validator  Validator
	repository Repository
}

func NewService(validator Validator, repository Repository) *service {
	return &service{validator, repository}
}

func (s *service) CreatePaymentMethod(data PaymentMethod) (*PaymentMethod, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	method, err := s.repository.CreatePaymentMethod(data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error creating payment method",
			"something went wrong while creating payment method",
		)
	}

	return method, nil
}
