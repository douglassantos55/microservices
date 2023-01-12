package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Carrier interface {
	GetQuote(origin, destination string, items []Item) (*Quote, error)
}

type coordinates struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func (c coordinates) String() string {
	return fmt.Sprintf("%s,%s", c.Lon, c.Lat)
}

type Route struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

type localCarrier struct {
	name      string
	gasPrice  float64
	distLiter float64
}

func NewLocalCarrier(price, dist float64) *localCarrier {
	return &localCarrier{"local", price, dist}
}

func (c *localCarrier) GetQuote(origin, dest string, items []Item) (*Quote, error) {
	routes, err := c.GetRoutes(origin, dest)
	if err != nil {
		return nil, err
	}

	quote := &Quote{Carrier: c.name, Value: math.Inf(0)}

	for _, route := range routes {
		baseValue := route.Distance / 1000 / c.distLiter * c.gasPrice * 1.15
		value := math.Round(baseValue*100) / 100

		if value < quote.Value {
			quote.Value = value
		}
	}

	return quote, nil
}

func (c *localCarrier) GetRoutes(origin, dest string) ([]Route, error) {
	originCoord, err := c.GetCoordinates(origin)
	if err != nil {
		return nil, err
	}

	destCoord, err := c.GetCoordinates(dest)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	format := "https://www.mapeia.com.br/route/v1/driving/%s;%s?overview=false&alternatives=true&steps=false&hints=;"
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(format, originCoord, destCoord), nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var path struct {
		Routes []Route `json:"routes"`
	}

	if err := json.NewDecoder(res.Body).Decode(&path); err != nil {
		return nil, err
	}

	return path.Routes, nil
}

func (c *localCarrier) GetCoordinates(address string) (*coordinates, error) {
	if strings.TrimSpace(address) == "" {
		return nil, errors.New("empty address")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	format := "https://www.mapeia.com.br/search?q=%s&addressdetails=1&namedetails=1&accept-language=pt-BR&countrycodes=br&format=jsonv2&limit=20"
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(format, url.QueryEscape(address)), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var coordinates []*coordinates
	if err := json.NewDecoder(res.Body).Decode(&coordinates); err != nil {
		return nil, err
	}

	if len(coordinates) == 0 {
		return nil, errors.New("address not found")
	}

	return coordinates[0], nil
}
