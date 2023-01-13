package pkg

import "net/http"

type PaymentMethod struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" validate:"required"`
	AccountID string `json:"account_id" validate:"required"`
}

type Service interface {
	CreatePaymentMethod(PaymentMethod) (*PaymentMethod, error)
	ListPaymentMethods() ([]*PaymentMethod, error)
	UpdatePaymentMethod(string, PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) error
}

type service struct {
	validator  Validator
	repository Repository
}

func NewService(validator Validator, repository Repository) Service {
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

func (s *service) UpdatePaymentMethod(id string, data PaymentMethod) (*PaymentMethod, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	method, err := s.repository.UpdatePaymentMethod(id, data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error updating payment method",
			"something went wrong updating payment method",
		)
	}

	return method, err
}

func (s *service) ListPaymentMethods() ([]*PaymentMethod, error) {
	return s.repository.ListPaymentMethods()
}

func (s *service) DeletePaymentMethod(id string) error {
	if _, err := s.repository.GetPaymentMethod(id); err != nil {
		return NewError(
			http.StatusNotFound,
			"payment method not found",
			"could not find the payment method you're trying to delete",
		)
	}

	if err := s.repository.DeletePaymentMethod(id); err != nil {
		return NewError(
			http.StatusInternalServerError,
			"error deleting equipment",
			"something went wrong while deleting equipment",
		)
	}

	return nil
}
