package pkg_test

import (
	"testing"

	"reconcip.com.br/microservices/delivery/pkg"
)

func TestCarrier(t *testing.T) {
	t.Run("GetCoordinates", func(t *testing.T) {
		t.Run("returns proper", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)
			coord, err := carrier.GetCoordinates("rua monte alegre do sul, mogi guaçu, São Paulo, SP")
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			expected := "-46.9489845,-22.3560953"
			if coord.String() != expected {
				t.Errorf("expected coordinate %v, got %v", expected, coord)
			}
		})

		t.Run("invalid address", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)
			_, err := carrier.GetCoordinates("rua que nao e para existir...")

			if err == nil {
				t.Fatal("expected error")
			}
		})

		t.Run("empty address", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)
			_, err := carrier.GetCoordinates("  ")

			if err == nil {
				t.Fatal("expected error")
			}
		})
	})

	t.Run("GetRoutes", func(t *testing.T) {
		t.Run("invalid origin", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

			origin := "address that should not be found...."
			dest := "avenida john boyd dunlop, campinas, São Paulo, SP"

			_, err := carrier.GetRoutes(origin, dest)

			if err == nil {
				t.Fatal("expected error")
			}
		})

		t.Run("invalid destination", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

			dest := "address that should not be found...."
			origin := "avenida john boyd dunlop, campinas, São Paulo, SP"

			_, err := carrier.GetRoutes(origin, dest)

			if err == nil {
				t.Fatal("expected error")
			}
		})

		t.Run("same origin and dest", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

			dest := "avenida john boyd dunlop, campinas, São Paulo, SP"
			origin := "avenida john boyd dunlop, campinas, São Paulo, SP"

			routes, err := carrier.GetRoutes(origin, dest)

			if err != nil {
				t.Fatal("did not expect error")
			}

			if len(routes) == 0 {
				t.Fatal("expected routes, got nothing")
			}

			if routes[0].Distance != 0 {
				t.Errorf("expected 0 distance, got %v", routes[0].Distance)
			}
		})

		t.Run("success", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "avenida john boyd dunlop, campinas, São Paulo, SP"

			routes, err := carrier.GetRoutes(origin, dest)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if len(routes) == 0 {
				t.Fatal("expected routes, got nothing")
			}

			route := routes[0]
			expectedDistance := 80390.8
			if route.Distance != expectedDistance {
				t.Errorf("expected distance %v, got %v", expectedDistance, route.Distance)
			}

			expectedDuration := 3947.9
			if route.Duration != expectedDuration {
				t.Errorf("expected distance %v, got %v", expectedDuration, route.Duration)
			}
		})

		t.Run("multiple routes", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"

			routes, err := carrier.GetRoutes(origin, dest)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if len(routes) != 2 {
				t.Fatalf("expected %v routes, got %v", 2, len(routes))
			}
		})
	})

	t.Run("GetQuote", func(t *testing.T) {
		t.Run("one path", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

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

			expectedValue := 55.47
			if quote.Value != expectedValue {
				t.Errorf("expected value %v, got %v", expectedValue, quote.Value)
			}
		})

		t.Run("multiple paths", func(t *testing.T) {
			carrier := pkg.NewLocalCarrier(6, 10)

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

			expectedValue := 169.02
			if quote.Value != expectedValue {
				t.Errorf("expected value %v, got %v", expectedValue, quote.Value)
			}
		})
	})
}
