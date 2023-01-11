package pkg

import "net/http"

type Equipment struct {
	ID             string         `json:"id" bson:"_id,omitempty" validate:"omitempty,required"`
	Description    string         `json:"description" validate:"required"`
	Stock          int            `json:"in_stock" validate:"omitempty,number"`
	EffectiveStock int            `json:"effective_qty" validate:"omitempty,number"`
	Weight         float64        `json:"weight" validate:"omitempty,numeric"`
	UnitValue      float64        `json:"unit_value" validate:"omitempty,numeric"`
	PurchaseValue  float64        `json:"purchase_value" validate:"omitempty,numeric"`
	ReplaceValue   float64        `json:"replace_value" validate:"omitempty,numeric"`
	MinQty         int            `json:"min_qty" validate:"omitempty,number"`
	SupplierID     string         `json:"supplier_id,omitempty"`
	Supplier       *Supplier      `json:"supplier"`
	RentingValues  []RentingValue `json:"renting_values" validate:"required,dive"`
}

type Supplier struct {
	ID         string  `json:"id"`
	SocialName string  `json:"social_name"`
	LegalName  string  `json:"legal_name"`
	Email      string  `json:"email"`
	Website    string  `json:"website"`
	Cnpj       string  `json:"cnpj"`
	InscEst    string  `json:"insc_est"`
	Phone      string  `json:"phone"`
	Address    Address `json:"address"`
}

type Address struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
}

type RentingValue struct {
	PeriodID string  `json:"period_id" validate:"required"`
	Value    float64 `json:"value"`
}

type Service interface {
	CreateEquipment(Equipment) (*Equipment, error)
	ListEquipment(page, perPage int) ([]*Equipment, int, error)
	UpdateEquipment(string, Equipment) (*Equipment, error)
}

type service struct {
	validator  Validator
	repository Repository
}

func NewService(validator Validator, repository Repository) Service {
	return &service{validator, repository}
}

func (s *service) CreateEquipment(data Equipment) (*Equipment, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}
	return s.repository.Create(data)
}

func (s *service) ListEquipment(page, perPage int) ([]*Equipment, int, error) {
	return s.repository.List(page, perPage)
}

func (s *service) UpdateEquipment(id string, data Equipment) (*Equipment, error) {
	if _, err := s.repository.Get(id); err != nil {
		return nil, NewError(
			http.StatusNotFound,
			"equipment not found",
			"could not find the equipment you're trying to edit",
		)
	}

	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	equipment, err := s.repository.Update(id, data)
	if err != nil {
		return nil, NewError(
			http.StatusInternalServerError,
			"error updating equipment",
			"something went wrong while updating equipment",
		)
	}

	return equipment, nil
}
