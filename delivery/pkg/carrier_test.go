package pkg_test

import (
	"strconv"
	"testing"

	"reconcip.com.br/microservices/delivery/pkg"
)

type fakeCoordinator struct{}

func (c *fakeCoordinator) GetCoordinates(source string) (*pkg.Coordinates, error) {
	if source == "rua monte alegre do sul, mogi guaçu, São Paulo, SP" {
		return &pkg.Coordinates{"0", "0"}, nil
	} else if source == "avenida john boyd dunlop, campinas, São Paulo, SP" {
		return &pkg.Coordinates{"1", "1"}, nil
	}
	return &pkg.Coordinates{"2", "2"}, nil
}

type fakeRouter struct{}

func (r *fakeRouter) GetRoutes(from, to *pkg.Coordinates) ([]pkg.Route, error) {
	dest, _ := strconv.ParseInt(to.Lat, 0, 0)
	origin, _ := strconv.ParseInt(from.Lat, 0, 0)

	distance := dest - origin

	routes := make([]pkg.Route, distance)
	for i := int(distance) - 1; i >= 0; i-- {
		routes[i] = pkg.Route{Distance: float64((i + 1) * 10000)}
	}

	return routes, nil
}

func TestCarrier(t *testing.T) {
	t.Run("GetQuote", func(t *testing.T) {
		t.Run("one path", func(t *testing.T) {
			router := &fakeRouter{}
			coordinator := &fakeCoordinator{}
			carrier := pkg.NewLocalCarrier(6, 10, router, coordinator)

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "avenida john boyd dunlop, campinas, São Paulo, SP"

			items := []pkg.Item{
				{Qty: 10, Weight: 1, Width: 1, Height: 1, Depth: 1},
				{Qty: 100, Weight: 1, Width: 1, Height: 1, Depth: 1},
				{Qty: 20, Weight: 1, Width: 1, Height: 1, Depth: 1},
				{Qty: 50, Weight: 1, Width: 1, Height: 1, Depth: 1},
			}

			quote, err := carrier.GetQuote(origin, dest, items)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if quote == nil {
				t.Fatal("expected quote, got nothing")
			}

			if quote.Carrier != "local" {
				t.Errorf("expected carrier local, got %v", quote.Carrier)
			}

			expectedValue := 6.9
			if quote.Value != expectedValue {
				t.Errorf("expected value %v, got %v", expectedValue, quote.Value)
			}
		})

		t.Run("multiple paths", func(t *testing.T) {
			router := &fakeRouter{}
			coordinator := &fakeCoordinator{}
			carrier := pkg.NewLocalCarrier(5, 10, router, coordinator)

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"

			items := []pkg.Item{
				{Qty: 10, Weight: 1, Width: 1, Height: 1, Depth: 1},
				{Qty: 100, Weight: 1, Width: 1, Height: 1, Depth: 1},
				{Qty: 20, Weight: 1, Width: 1, Height: 1, Depth: 1},
				{Qty: 50, Weight: 1, Width: 1, Height: 1, Depth: 1},
			}

			quote, err := carrier.GetQuote(origin, dest, items)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if quote == nil {
				t.Fatal("expected quote, got nothing")
			}

			if quote.Carrier != "local" {
				t.Errorf("expected carrier local, got %v", quote.Carrier)
			}

			expectedValue := 5.75
			if quote.Value != expectedValue {
				t.Errorf("expected value %v, got %v", expectedValue, quote.Value)
			}
		})
	})
}
