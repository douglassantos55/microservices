package pkg

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name" validate:"required"`
	Email        string             `bson:"email" json:"email" validate:"omitempty,email"`
	Birthdate    time.Time          `bson:"birthdate" json:"birthdate" validate:"omitempty"`
	CpfCnpj      string             `bson:"cpf_cnpj" json:"cpf_cnpj" validate:"required"`
	RgInscEst    string             `bson:"rg_insc_est" json:"rg_insc_est"`
	Phone        string             `bson:"phone" json:"phone"`
	Cellphone    string             `bson:"cellphone" json:"cellphone"`
	Ocupation    string             `bson:"ocupation" json:"ocupation"`
	Address      Address            `bson:"inline" json:"address" validate:"required"`
	Observations string             `bson:"observations" json:"observations"`
}

type Address struct {
	Address      string `bson:"address" json:"address"`
	Number       string `bson:"number" json:"number"`
	Complement   string `bson:"complement" json:"complement"`
	Neighborhood string `bson:"neighborhood" json:"neighborhood"`
	City         string `bson:"city" json:"city"`
	State        string `bson:"state" json:"state"`
	Postcode     string `bson:"postcode" json:"postcode"`
}

type Service interface {
	List(page, perPage int64) ([]*Customer, error)
	Create(Customer) (*Customer, error)
}

type service struct {
	validator  Validator
	repository Repository
}

func NewService(validator Validator, repository Repository) Service {
	return &service{validator, repository}
}

func (s *service) List(page, perPage int64) ([]*Customer, error) {
	return s.repository.List(page, perPage)
}

func (s *service) Create(customer Customer) (*Customer, error) {
	if err := s.validator.Validate(customer); err != nil {
		return nil, err
	}
	return s.repository.Create(customer)
}
