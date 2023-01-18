package pkg

import (
	"context"
	"net/http"
	"time"
)

type Rent struct {
	ID                 string            `json:"id" bson:"_id,omitempty"`
	PeriodID           string            `json:"period_id" validate:"required"`
	PaymentMethodID    string            `json:"payment_method_id" validate:"required,payment_method"`
	PaymentMethod      *PaymentMethod    `json:"payment_method,omitempty"`
	PaymentConditionID string            `json:"payment_condition_id" validate:"required,payment_condition"`
	PaymentCondition   *PaymentCondition `json:"payment_condition,omitempty"`
	PaymentTypeID      string            `json:"payment_type_id" validate:"required,payment_type"`
	PaymentType        *PaymentType      `json:"payment_type,omitempty"`
	CarrierID          string            `json:"carrier_id" validate:"required"`
	CustomerID         string            `json:"customer_id" validate:"required,customer"`
	Customer           *Customer         `json:"customer,omitempty"`
	StartDate          time.Time         `json:"start_date" validate:"required"`
	EndDate            time.Time         `json:"end_date" validate:"required"`
	Items              []Item            `json:"items" validate:"required,dive"`
	QtyDays            int               `json:"qty_days" validate:"required"`
	Discount           float64           `json:"discount" validate:"omitempty,gt=0"`
	PaidValue          float64           `json:"paid_value" validate:"omitempty,gt=0"`
	Bill               float64           `json:"bill" validate:"required_with=PaidValue,ltefield=PaidValue"`
	Observations       string            `json:"observations"`
	CheckInfo          string            `json:"check_info"`
	DeliveryValue      float64           `json:"delivery_value" validate:"omitempty,required_with=CarrierID"`
	DeliveryAddress    string            `json:"delivery_address" validate:"omitempty,required_with=CarrierID"`
	UsageAddress       string            `json:"usage_address"`
}

type Item struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	EquipmentID string `json:"equipment_id" validate:"required"`
	Qty         int    `json:"qty" validate:"required,gt=0"`
}

type PaymentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PaymentMethod struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Account *Account `json:"account,omitempty"`
}

type PaymentCondition struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Increment    float32      `json:"increment"`
	Installments []int32      `json:"installments"`
	PaymentType  *PaymentType `json:"payment_type"`
}

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email" validate:"omitempty,email"`
	CpfCnpj   string `json:"cpf_cnpj" validate:"required,cpf_cnpj"`
	RgInscEst string `json:"rg_insc_est"`
	Phone     string `json:"phone"`
	Cellphone string `json:"cellphone"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	CreateRent(Rent) (*Rent, error)
}

type service struct {
	validator  Validator
	repository Repository
}

func NewService(validator Validator, repository Repository) *service {
	return &service{validator, repository}
}

func (s *service) CreateRent(data Rent) (*Rent, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	rent, err := s.repository.CreateRent(data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error creating rent",
			"something went wrong creating rent",
		)
	}

	return rent, nil
}
