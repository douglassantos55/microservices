package pkg

import (
	"errors"
	"log"
)

var (
	ErrNoItems    = errors.New("no items for delivery")
	ErrNoCarriers = errors.New("no carriers")
)

type Quote struct {
	Carrier string  `json:"company"`
	Value   float64 `json:"value"`
}

type Item struct {
	Qty    int     `json:"qty"`
	Weight float64 `json:"weight"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Depth  float64 `json:"depth"`
}

type Service interface {
	GetQuotes(origin, destination string, items []Item) ([]*Quote, error)
}

type service struct {
	carriers []Carrier
}

func NewService(carriers []Carrier) Service {
	return &service{carriers}
}

func (s *service) GetQuotes(origin, destination string, items []Item) ([]*Quote, error) {
	if len(s.carriers) == 0 {
		return nil, ErrNoCarriers
	}

	if len(items) == 0 {
		return nil, ErrNoItems
	}

	quotes := make([]*Quote, 0)
	for _, carrier := range s.carriers {
		quote, err := carrier.GetQuote(origin, destination, items)
		if err != nil {
			log.Printf("could not get quotes from carrier \"%s\": %s", carrier, err)
			continue
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}
