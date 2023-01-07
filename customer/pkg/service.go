package pkg

import (
	"math"
	"net/http"
	"time"
)

type Customer struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string    `bson:"name" json:"name" validate:"required"`
	Email        string    `bson:"email" json:"email" validate:"omitempty,email"`
	Birthdate    time.Time `bson:"birthdate" json:"birthdate" validate:"omitempty"`
	CpfCnpj      string    `bson:"cpf_cnpj" json:"cpf_cnpj" validate:"required"`
	RgInscEst    string    `bson:"rg_insc_est" json:"rg_insc_est"`
	Phone        string    `bson:"phone" json:"phone"`
	Cellphone    string    `bson:"cellphone" json:"cellphone"`
	Ocupation    string    `bson:"ocupation" json:"ocupation"`
	Address      Address   `bson:"inline" json:"address" validate:"required"`
	Observations string    `bson:"observations" json:"observations"`
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

type ListResult struct {
	Items      []any `json:"items"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

type Service interface {
	List(page, perPage int64) (*ListResult, error)
	Create(Customer) (*Customer, error)
	Update(id string, data Customer) (*Customer, error)
}

type service struct {
	validator  Validator
	repository Repository
}

func NewService(validator Validator, repository Repository) Service {
	return &service{validator, repository}
}

func (s *service) List(page, perPage int64) (*ListResult, error) {
	customers, total, err := s.repository.List(page, perPage)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error fetching customers",
			"something went wrong while fetching customers",
		)
	}

	totalPages := int64(math.Max(1, float64(total/perPage)))
	if page >= totalPages {
		return nil, NewError(
			http.StatusBadRequest,
			"invalid page",
			"page exceeds the maximum number of pages available",
		)
	}

	items := make([]any, len(customers))
	for i, customer := range customers {
		items[i] = customer
	}

	return &ListResult{
		Items:      items,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (s *service) Create(customer Customer) (*Customer, error) {
	if err := s.validator.Validate(customer); err != nil {
		return nil, err
	}
	return s.repository.Create(customer)
}

func (s *service) Update(id string, data Customer) (*Customer, error) {
	_, err := s.repository.Get(id)
	if err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"customer not found",
			"could not find the customer you're trying to edit",
		)
	}
	return s.repository.Update(id, data)
}
