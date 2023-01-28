package pkg

import (
	"net/http"
	"time"
)

type Method struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" validate:"required"`
	AccountID string `json:"account_id" validate:"required"`
}

type Type struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" validate:"required"`
}

type Condition struct {
	ID            string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string  `json:"name" validate:"required"`
	Increment     float32 `json:"increment" validate:"min=0"`
	PaymentTypeID string  `json:"payment_type_id" validate:"required,paymenttype"`
	PaymentType   *Type   `json:"payment_type,omitempty"`
	Installments  []int32 `json:"installments,dive,gt=0"`
}

type Invoice struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	CustomerID string    `json:"customer_id" validate:"required"`
	Customer   *Customer `json:"customer,omitempty"`
	DueDate    time.Time `json:"due_date" validate:"required,gt"`
	Total      float64   `json:"total" validate:"required,gt=0"`
	Items      []Item    `json:"items" validate:"required,dive"`
}

type Item struct {
	Description string  `json:"description" validate:"required"`
	Total       float64 `json:"total" validate:"required,gt=0"`
}

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CpfCnpj   string `json:"cpf_cnpj"`
	RgInscEst string `json:"rg_insc_est"`
	Phone     string `json:"phone"`
	Cellphone string `json:"cellphone"`
}

type Service interface {
	CreatePaymentMethod(Method) (*Method, error)
	ListPaymentMethods() ([]*Method, error)
	UpdatePaymentMethod(string, Method) (*Method, error)
	DeletePaymentMethod(string) error
	GetPaymentMethod(string) (*Method, error)

	CreatePaymentType(Type) (*Type, error)
	ListPaymentTypes() ([]*Type, error)
	UpdatePaymentType(string, Type) (*Type, error)
	DeletePaymentType(string) error
	GetPaymentType(string) (*Type, error)

	CreatePaymentCondition(Condition) (*Condition, error)
	ListPaymentConditions() ([]*Condition, error)
	UpdatePaymentCondition(string, Condition) (*Condition, error)
	DeletePaymentCondition(string) error
	GetPaymentCondition(string) (*Condition, error)

	CreateInvoice(Invoice) (*Invoice, error)
	ListInvoices(page, perPage int64) ([]*Invoice, int64, error)
	UpdateInvoice(string, Invoice) (*Invoice, error)
	DeleteInvoice(string) error
	GetInvoice(string) (*Invoice, error)
}

type service struct {
	validator  Validator
	repository Repository
	gateway    Gateway
}

func NewService(validator Validator, repository Repository, gateway Gateway) Service {
	return &service{validator, repository, gateway}
}

func (s *service) CreatePaymentMethod(data Method) (*Method, error) {
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

func (s *service) UpdatePaymentMethod(id string, data Method) (*Method, error) {
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

func (s *service) ListPaymentMethods() ([]*Method, error) {
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

func (s *service) GetPaymentMethod(id string) (*Method, error) {
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

func (s *service) CreatePaymentType(data Type) (*Type, error) {
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

func (s *service) ListPaymentTypes() ([]*Type, error) {
	return s.repository.ListPaymentTypes()
}

func (s *service) UpdatePaymentType(id string, data Type) (*Type, error) {
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

func (s *service) GetPaymentType(id string) (*Type, error) {
	paymentType, err := s.repository.GetPaymentType(id)
	if err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"payment type not found",
			"could not find payment type",
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

func (s *service) GetPaymentCondition(id string) (*Condition, error) {
	condition, err := s.repository.GetPaymentCondition(id)
	if err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"condition not found",
			"could not find condition",
		)
	}
	return condition, nil
}

func (s *service) CreateInvoice(data Invoice) (*Invoice, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	invoice, err := s.repository.CreateInvoice(data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"could not create invoice",
			"there was an error while creating invoice",
		)
	}

	if err := s.gateway.ProcessPayment(invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *service) ListInvoices(page, perPage int64) ([]*Invoice, int64, error) {
	return s.repository.ListInvoices(page, perPage)
}

func (s *service) UpdateInvoice(id string, data Invoice) (*Invoice, error) {
	if _, err := s.repository.GetInvoice(id); err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"invoice not found",
			"could not find invoice",
		)
	}

	invoice, err := s.repository.UpdateInvoice(id, data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"could not update invoice",
			"there was an error updating invoice",
		)
	}

	return invoice, nil
}

func (s *service) DeleteInvoice(id string) error {
	if _, err := s.repository.GetInvoice(id); err != nil {
		return NewError(
			http.StatusNotFound,
			"invoice not found",
			"could not find invoice",
		)
	}
	return s.repository.DeleteInvoice(id)
}

func (s *service) GetInvoice(id string) (*Invoice, error) {
	invoice, err := s.repository.GetInvoice(id)
	if err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"invoice not found",
			"could not find invoice",
		)
	}
	return invoice, nil
}
