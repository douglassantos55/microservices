package pkg

import "math"

type Carrier interface {
	String() string
	GetQuote(origin, destination string, items []Item) (*Quote, error)
}

type localCarrier struct {
	name      string
	gasPrice  float64
	distLiter float64

	router Router
	coord  Coordinator
}

func NewLocalCarrier(price, dist float64, router Router, coord Coordinator) *localCarrier {
	return &localCarrier{"local", price, dist, router, coord}
}

func (c *localCarrier) String() string {
	return c.name
}

func (c *localCarrier) GetQuote(from, to string, items []Item) (*Quote, error) {
	origin, err := c.coord.GetCoordinates(from)
	if err != nil {
		return nil, err
	}

	dest, err := c.coord.GetCoordinates(to)
	if err != nil {
		return nil, err
	}

	routes, err := c.router.GetRoutes(origin, dest)
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
