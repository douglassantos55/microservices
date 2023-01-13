package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Router interface {
	GetRoutes(from, to *Coordinates) ([]Route, error)
}

type Coordinator interface {
	GetCoordinates(source string) (*Coordinates, error)
}

type Route struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

type Coordinates struct {
	Lon string `json:"lon"`
	Lat string `json:"lat"`
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%s,%s", c.Lon, c.Lat)
}

type mapeiaCoordinator struct {
	url string
}

func NewMapeiaCoordinator() *mapeiaCoordinator {
	return &mapeiaCoordinator{"https://www.mapeia.com.br/search?"}
}

func (c *mapeiaCoordinator) GetCoordinates(source string) (*Coordinates, error) {
	if strings.TrimSpace(source) == "" {
		return nil, errors.New("empty address")
	}

	values := url.Values{}
	values.Set("q", source)
	values.Set("addressdetails", "1")
	values.Set("namedetails", "1")
	values.Set("accept-language", "pt-BR")
	values.Set("contrycodes", "br")
	values.Set("format", "jsonv2")
	values.Set("limit", "20")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	url := c.url + values.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var coordinates []*Coordinates
	if err := json.NewDecoder(res.Body).Decode(&coordinates); err != nil {
		return nil, err
	}

	if len(coordinates) == 0 {
		return nil, errors.New("address not found")
	}

	return coordinates[0], nil
}

type mapeiaRouter struct {
	url string
}

func NewMapeiaRouter() *mapeiaRouter {
	return &mapeiaRouter{"https://www.mapeia.com.br/route/v1/driving/%s;%s?"}
}

func (r *mapeiaRouter) GetRoutes(from, to *Coordinates) ([]Route, error) {
	if from == nil || to == nil {
		return nil, errors.New("invalid coordinates")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	values := url.Values{}
	values.Set("overview", "false")
	values.Set("alternatives", "true")
	values.Set("steps", "false")
	values.Set("hints", ";")

	url := fmt.Sprintf(r.url, from, to) + values.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

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
