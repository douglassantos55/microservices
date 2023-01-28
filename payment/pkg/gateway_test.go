package pkg_test

import (
	"testing"

	"reconcip.com.br/microservices/payment/pkg"
)

func TestStripeGateway(t *testing.T) {
	t.Run("get customer", func(t *testing.T) {
		gateway := pkg.NewStripeGateway("sk_test_4eC39HqLyjWDarjtT1zdp7dc")

		expected := "cus_NFVWErxwNNoa84"
		got, err := gateway.GetCustomer("somecustomerid")

		if err != nil {
			t.Fatalf("did not expect error, got: %v", err)
		}

		if expected != got {
			t.Errorf("expected id %v, got %v", expected, got)
		}
	})

	t.Run("create customer", func(t *testing.T) {
		gateway := pkg.NewStripeGateway("sk_test_4eC39HqLyjWDarjtT1zdp7dc")

		expected, err := gateway.CreateCustomer(&pkg.Customer{
			ID:    "johndoeid",
			Name:  "john doe",
			Email: "johndoe@email.com",
			Phone: "+33 52305222",
		})

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if expected == "" {
			t.Error("expected customer id")
		}
	})

	t.Run("process payment", func(t *testing.T) {
		gateway := pkg.NewStripeGateway("sk_test_4eC39HqLyjWDarjtT1zdp7dc")

		err := gateway.ProcessPayment(&pkg.Invoice{
			ID:         "someinvoiceid",
			CustomerID: "somecustomerid",
		})

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}
