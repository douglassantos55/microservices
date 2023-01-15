package pkg

import "net/http"

type PaymentMethod struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" validate:"required"`
	AccountID string `json:"account_id" validate:"required"`
}

type PaymentType struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" validate:"required"`
}

type Condition struct {
	ID            string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string  `json:"name" validate:"required"`
	Increment     float64 `json:"increment" validate:"min=0"`
	PaymentTypeID string  `json:"payment_type_id" validate:"required,exists=payment_types"`
	Installments  []int   `json:"installments,dive,gt=0"`
}

type Service interface {
	CreatePaymentMethod(PaymentMethod) (*PaymentMethod, error)
	ListPaymentMethods() ([]*PaymentMethod, error)
	UpdatePaymentMethod(string, PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) error
	GetPaymentMethod(string) (*PaymentMethod, error)

	CreatePaymentType(PaymentType) (*PaymentType, error)
	ListPaymentTypes() ([]*PaymentType, error)
	UpdatePaymentType(string, PaymentType) (*PaymentType, error)
	DeletePaymentType(string) error
	GetPaymentType(string) (*PaymentType, error)

	CreatePaymentCondition(Condition) (*Condition, error)
	ListPaymentConditions() ([]*Condition, error)
	UpdatePaymentCondition(string, Condition) (*Condition, error)
	DeletePaymentCondition(string) error
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
			"error deleting payment method",
			"something went wrong while deleting payment method",
		)
	}

	return nil
}

func (s *service) GetPaymentMethod(id string) (*PaymentMethod, error) {
	method, err := s.repository.GetPaymentMethod(id)
	if err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"payment method not found",
			"could not find payment method",
		)
	}
	return method, nil
}

func (s *service) CreatePaymentType(data PaymentType) (*PaymentType, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	paymentType, err := s.repository.CreatePaymentType(data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error creating payment type",
			"could not create payment type due to an error",
		)
	}

	return paymentType, nil
}

func (s *service) ListPaymentTypes() ([]*PaymentType, error) {
	return s.repository.ListPaymentTypes()
}

func (s *service) UpdatePaymentType(id string, data PaymentType) (*PaymentType, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	method, err := s.repository.UpdatePaymentType(id, data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error updating payment type",
			"something went wrong updating payment type",
		)
	}

	return method, err
}

func (s *service) DeletePaymentType(id string) error {
	if _, err := s.repository.GetPaymentType(id); err != nil {
		return NewError(
			http.StatusNotFound,
			"payment type not found",
			"could not find the payment type you're trying to delete",
		)
	}

	if err := s.repository.DeletePaymentType(id); err != nil {
		return NewError(
			http.StatusInternalServerError,
			"error deleting payment type",
			"something went wrong while deleting payment type",
		)
	}

	return nil
}

func (s *service) GetPaymentType(id string) (*PaymentType, error) {
	paymentType, err := s.repository.GetPaymentType(id)
	if err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"payment method not found",
			"could not find payment method",
		)
	}
	return paymentType, nil
}

func (s *service) CreatePaymentCondition(data Condition) (*Condition, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	condition, err := s.repository.CreatePaymentCondition(data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error creating condition",
			"something went wrong creating condition",
		)
	}

	return condition, nil
}

func (s *service) ListPaymentConditions() ([]*Condition, error) {
	return s.repository.ListPaymentConditions()
}

func (s *service) UpdatePaymentCondition(id string, data Condition) (*Condition, error) {
	if _, err := s.repository.GetPaymentCondition(id); err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"condition not found",
			"could not find the condition you're trying to update",
		)
	}

	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	condition, err := s.repository.UpdatePaymentCondition(id, data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error updating condition",
			"something went wrong updating condition",
		)
	}

	return condition, nil
}

func (s *service) DeletePaymentCondition(id string) error {
	if _, err := s.repository.GetPaymentCondition(id); err != nil {
		return NewError(
			http.StatusNotFound,
			"condition not found",
			"could not find condition you're trying to delete",
		)
	}
	return s.repository.DeletePaymentCondition(id)
}
