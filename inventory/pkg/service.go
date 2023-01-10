package pkg

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
	SupplierID     string         `json:"supplier_id"`
	RentingValues  []RentingValue `json:"renting_values" validate:"required,dive"`
}

type RentingValue struct {
	PeriodID string  `json:"period_id" validate:"required"`
	Value    float64 `json:"value"`
}

type Service interface {
	CreateEquipment(Equipment) (*Equipment, error)
	ListEquipment(page, perPage int) ([]*Equipment, int, error)
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
