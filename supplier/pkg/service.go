package pkg

type Supplier struct {
	ID           string  `json:"id,omitempty" bson:"_id,omitempty"`
	SocialName   string  `json:"social_name" bson:"social_name" validate:"required"`
	LegalName    string  `json:"legal_name" bson:"legal_name" `
	Email        string  `json:"email" bson:"email" validate:"omitempty,email"`
	Website      string  `json:"website" bson:"website"`
	Cnpj         string  `json:"cnpj" bson:"cnpj" validate:"required,cnpj"`
	InscEst      string  `json:"insc_est" bson:"insc_est"`
	Phone        string  `json:"phone" bson:"phone"`
	Address      Address `json:"address" bson:"inline"`
	Observations string  `json:"observations" bson:"observations"`
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

type Service interface {
	Create(Supplier) (*Supplier, error)
}

func NewService(validator Validator, repository Repository) Service {
	return &service{validator, repository}
}

type service struct {
	validator  Validator
	repository Repository
}

func (s *service) Create(data Supplier) (*Supplier, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}
	return s.repository.Create(data)
}
