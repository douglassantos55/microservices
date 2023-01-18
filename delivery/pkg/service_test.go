package pkg_test

import (
	"testing"

	"reconcip.com.br/microservices/delivery/pkg"
)

func TestService(t *testing.T) {
	t.Run("GetQuotes", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			router := &fakeRouter{}
			coordinator := &fakeCoordinator{}

			svc := pkg.NewService([]pkg.Carrier{
				pkg.NewLocalCarrier(5, 12, router, coordinator),
			})

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"
			items := []pkg.Item{
				{Qty: 1, Weight: 1, Width: 1, Height: 1, Depth: 1},
			}

			quotes, err := svc.GetQuotes(origin, dest, items)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if len(quotes) == 0 {
				t.Fatal("expected quotes, got nothing")
			}

			if quotes[0].Carrier != "local" {
				t.Errorf("expected carrier local, got %v", quotes[0].Carrier)
			}

			expectedValue := 4.79
			if quotes[0].Value != expectedValue {
				t.Errorf("expected value %v, got %v", expectedValue, quotes[0].Value)
			}
		})

		t.Run("no carriers", func(t *testing.T) {
			svc := pkg.NewService([]pkg.Carrier{})

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"
			items := []pkg.Item{
				{Qty: 1, Weight: 1, Width: 1, Height: 1, Depth: 1},
			}

			if _, err := svc.GetQuotes(origin, dest, items); err == nil {
				t.Fatalf("expected error %v, got nothing", pkg.ErrNoCarriers)
			}
		})

		t.Run("no items", func(t *testing.T) {
			svc := pkg.NewService([]pkg.Carrier{})

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"
			items := []pkg.Item{}

			if _, err := svc.GetQuotes(origin, dest, items); err == nil {
				t.Fatalf("expected error %v, got nothing", pkg.ErrNoItems)
			}
		})
	})

	t.Run("GetQuote", func(t *testing.T) {
		t.Run("invalid carrier", func(t *testing.T) {
			router := &fakeRouter{}
			coordinator := &fakeCoordinator{}

			svc := pkg.NewService([]pkg.Carrier{
				pkg.NewLocalCarrier(5, 12, router, coordinator),
			})

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"
			items := []pkg.Item{
				{Qty: 1, Weight: 1, Width: 1, Height: 1, Depth: 1},
			}

			if _, err := svc.GetQuote(origin, dest, "correios", items); err != pkg.ErrCarrierNotFound {
				t.Errorf("expected error %v, got %v", pkg.ErrCarrierNotFound, err)
			}
		})

		t.Run("no items", func(t *testing.T) {
			router := &fakeRouter{}
			coordinator := &fakeCoordinator{}

			svc := pkg.NewService([]pkg.Carrier{
				pkg.NewLocalCarrier(5, 12, router, coordinator),
			})

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"
			items := []pkg.Item{}

			if _, err := svc.GetQuote(origin, dest, "local", items); err != pkg.ErrNoItems {
				t.Errorf("expected error %v, got %v", pkg.ErrNoItems, err)
			}
		})
		t.Run("success", func(t *testing.T) {
			router := &fakeRouter{}
			coordinator := &fakeCoordinator{}

			svc := pkg.NewService([]pkg.Carrier{
				pkg.NewLocalCarrier(5, 12, router, coordinator),
			})

			origin := "rua monte alegre do sul, mogi guaçu, São Paulo, SP"
			dest := "santos, São Paulo, SP"
			items := []pkg.Item{
				{Qty: 1, Weight: 1, Width: 1, Height: 1, Depth: 1},
			}

			quote, err := svc.GetQuote(origin, dest, "local", items)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if quote == nil {
				t.Fatal("expected quote, got nothing")
			}

			if quote.Carrier != "local" {
				t.Errorf("expected carrier local, got %v", quote.Carrier)
			}

			expectedValue := 4.79
			if quote.Value != expectedValue {
				t.Errorf("expected value %v, got %v", expectedValue, quote.Value)
			}
		})
	})
}
