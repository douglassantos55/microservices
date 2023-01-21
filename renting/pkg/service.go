package pkg

import (
	"encoding/json"
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
	Items              []*Item           `json:"items" validate:"required,dive"`
	Discount           float64           `json:"discount" validate:"omitempty,gt=0"`
	PaidValue          float64           `json:"paid_value" validate:"omitempty,gt=0"`
	Bill               float64           `json:"bill" validate:"required_with=PaidValue,ltefield=PaidValue"`
	Observations       string            `json:"observations"`
	CheckInfo          string            `json:"check_info"`
	DeliveryValue      float64           `json:"delivery_value"`
	DeliveryAddress    string            `json:"delivery_address" validate:"required_with=CarrierID"`
	UsageAddress       string            `json:"usage_address"`
}

func (r *Rent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":                r.ID,
		"period":            r.PeriodID,
		"start_date":        r.StartDate.Local(),
		"end_date":          r.EndDate.Local(),
		"qty_days":          r.GetQtyDays(),
		"customer":          r.Customer,
		"carrier":           r.CarrierID,
		"observations":      r.Observations,
		"usage_address":     r.UsageAddress,
		"deliver_address":   r.DeliveryAddress,
		"delivery_value":    r.DeliveryValue,
		"subtotal":          r.GetSubtotal(),
		"discount":          r.Discount,
		"total":             r.GetTotal(),
		"paid_value":        r.PaidValue,
		"bill":              r.Bill,
		"change":            r.GetChange(),
		"remaining":         r.GetRemaining(),
		"check_info":        r.CheckInfo,
		"total_weight":      r.GetTotalWeight(),
		"total_unit_value":  r.GetTotalUnitValue(),
		"total_pieces":      r.GetTotalPieces(),
		"payment_method":    r.PaymentMethod,
		"payment_type":      r.PaymentType,
		"payment_condition": r.PaymentCondition,
		"items":             r.Items,
	})
}

func (r *Rent) GetQtyDays() int {
	return int(r.EndDate.Sub(r.StartDate).Hours() / 24)
}

func (r *Rent) GetTotal() float64 {
	return r.GetSubtotal() + r.DeliveryValue - r.Discount
}

func (r *Rent) GetChange() float64 {
	return r.Bill - r.PaidValue
}

func (r *Rent) GetRemaining() float64 {
	return r.GetTotal() - r.PaidValue
}

func (r *Rent) GetSubtotal() float64 {
	total := 0.0
	for _, item := range r.Items {
		total += item.GetSubtotal(r.PeriodID)
	}
	return total
}

func (r *Rent) GetTotalWeight() float64 {
	total := 0.0
	for _, item := range r.Items {
		total += item.GetSubtotalWeight()
	}
	return total
}

func (r *Rent) GetTotalUnitValue() float64 {
	total := 0.0
	for _, item := range r.Items {
		total += float64(item.Qty) * item.Equipment.UnitValue
	}
	return total
}

func (r *Rent) GetTotalPieces() int {
	total := 0
	for _, item := range r.Items {
		total += item.Qty
	}
	return total
}

type Item struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	EquipmentID string     `json:"equipment_id" validate:"required"`
	Equipment   *Equipment `json:"equipment"`
	Qty         int        `json:"qty" validate:"required,gt=0,ltecsfield=Equipment.EffectiveStock"`
}

func (i *Item) GetSubtotal(period string) float64 {
	return float64(i.Qty) * i.Equipment.GetRentingValue(period)
}

func (i *Item) GetSubtotalWeight() float64 {
	return float64(i.Qty) * i.Equipment.Weight
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

type Equipment struct {
	ID             string          `json:"id"`
	Description    string          `json:"description"`
	Weight         float64         `json:"weight"`
	UnitValue      float64         `json:"unit_value"`
	EffectiveStock int             `json:"effective_qty"`
	RentingValues  []*RentingValue `json:"renting_values" validate:"required,dive"`
}

func (e *Equipment) GetRentingValue(period string) float64 {
	for _, value := range e.RentingValues {
		if value.PeriodID == period {
			return value.Value
		}
	}
	return 0
}

type RentingValue struct {
	PeriodID string  `json:"period_id"`
	Period   *Period `json:"period,omitempty"`
	Value    float64 `json:"value"`
}

type Period struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	QtyDays int32  `json:"qty_days"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Quote struct {
	Carrier string
	Value   float64
}

type Service interface {
	CreateRent(Rent) (*Rent, error)
}

type DeliveryService interface {
	GetQuote(origin, dest, carrier string, items []*Item) (*Quote, error)
}

type InventoryService interface {
	ReduceStock(items []*Item) error
}

type service struct {
	validator  Validator
	repository Repository
	delivery   DeliveryService
}

func NewService(validator Validator, repository Repository, delivery DeliveryService) Service {
	return &service{validator, repository, delivery}
}

func (s *service) CreateRent(data Rent) (*Rent, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}

	if data.CarrierID != "" {
		origin := "rua monte alegre do sul, mogi guacu, sp"
		quote, err := s.delivery.GetQuote(origin, data.DeliveryAddress, data.CarrierID, data.Items)

		if err != nil {
			return nil, err
		}

		data.DeliveryValue = quote.Value
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
