package pkg_test

import (
	"testing"

	"reconcip.com.br/microservices/delivery/pkg"
)

func TestMapeiaCoordinator(t *testing.T) {
	t.Run("returns proper", func(t *testing.T) {
		coordinator := pkg.NewMapeiaCoordinator()

		coord, err := coordinator.GetCoordinates("rua monte alegre do sul, mogi guaçu, São Paulo, SP")
		if err != nil {
			t.Fatalf("did not expect error: %v", err)
		}

		expected := "-46.9489845,-22.3560953"
		if coord.String() != expected {
			t.Errorf("expected coordinate %v, got %v", expected, coord)
		}
	})

	t.Run("invalid address", func(t *testing.T) {
		coordinator := pkg.NewMapeiaCoordinator()
		_, err := coordinator.GetCoordinates("rua que nao e para existir...")

		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestMapeiaRouter(t *testing.T) {
	t.Run("invalid origin", func(t *testing.T) {
		router := pkg.NewMapeiaRouter()
		dest := &pkg.Coordinates{"-46.9489845", "-22.3560953"}

		_, err := router.GetRoutes(nil, dest)

		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("invalid destination", func(t *testing.T) {
		router := pkg.NewMapeiaRouter()
		origin := &pkg.Coordinates{"-46.9489845", "-22.3560953"}

		_, err := router.GetRoutes(origin, nil)

		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("same origin and dest", func(t *testing.T) {
		router := pkg.NewMapeiaRouter()

		origin := &pkg.Coordinates{"-46.9489845", "-22.3560953"}
		dest := &pkg.Coordinates{"-46.9489845", "-22.3560953"}

		routes, err := router.GetRoutes(origin, dest)

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
		router := pkg.NewMapeiaRouter()

		origin := &pkg.Coordinates{"-46.9489845", "-22.3560953"}
		dest := &pkg.Coordinates{"-47.167642", "-22.9390489"}

		routes, err := router.GetRoutes(origin, dest)
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
		router := pkg.NewMapeiaRouter()

		origin := &pkg.Coordinates{"-46.9489845", "-22.3560953"}
		dest := &pkg.Coordinates{"-46.333889", "-23.960833"}

		routes, err := router.GetRoutes(origin, dest)
		if err != nil {
			t.Fatalf("did not expect error: %v", err)
		}

		if len(routes) != 2 {
			t.Fatalf("expected %v routes, got %v", 2, len(routes))
		}
	})
}
