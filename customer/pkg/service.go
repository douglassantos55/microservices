package pkg

import "time"

type Customer struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name" validate:"required"`
	Email        string     `json:"email" validate:"omitempty,email"`
	Birthdate    *time.Time `json:"birthdate"`
	CpfCnpj      string     `json:"cpf_cnpj" validate:"required"`
	RgInscEst    string     `json:"rg_insc_est"`
	Phone        string     `json:"phone"`
	Cellphone    string     `json:"cellphone"`
	Ocupation    string     `json:"ocupation"`
	Address      Address    `json:"address" validate:"required"`
	Observations string     `json:"observations"`
}

type Address struct {
	Address      string `json:"address"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
}

type Service interface {
	Create(Customer) (*Customer, error)
}

type service struct {
	validator Validator
}

func NewService(validator Validator) Service {
	return &service{validator}
}

func (s *service) Create(customer Customer) (*Customer, error) {
	if err := s.validator.Validate(customer); err != nil {
		return nil, err
	}
	customer.ID = "aoeu"
	return &customer, nil
}
