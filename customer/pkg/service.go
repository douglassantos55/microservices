package pkg

type Customer struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty" validate:"required"`
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
