package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Gateway interface {
	ProcessPayment(*Invoice) error
}

type StripeGateway struct {
	url    string
	apiKey string
}

func NewStripeGateway(apiKey string) *StripeGateway {
	return &StripeGateway{"https://api.stripe.com/v1", apiKey}
}

func (g *StripeGateway) ProcessPayment(invoice *Invoice) error {
	customerId, err := g.GetCustomer(invoice.CustomerID)
	if err != nil {
		_, err = g.CreateCustomer(invoice.Customer)
		if err != nil {
			return err
		}
	}

	res, err := g.request("POST", "/invoices", url.Values{
		"auto_advance":          []string{"true"},
		"customer":              []string{customerId},
		"metadata[internal_id]": []string{invoice.ID},
	})

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("could not create invoice")
	}

	return nil
}

func (g *StripeGateway) CreateCustomer(customer *Customer) (string, error) {
	res, err := g.request("POST", "/customers", url.Values{
		"name":                  []string{customer.Name},
		"email":                 []string{customer.Email},
		"phone":                 []string{customer.Phone},
		"metadata[internal_id]": []string{customer.ID},
	})

	if err != nil {
		return "", err
	}

	var entity struct {
		ID       string         `json:"id"`
		Metadata map[string]any `json:"metadata"`
	}

	if err := json.NewDecoder(res.Body).Decode(&entity); err != nil {
		return "", err
	}

	return entity.ID, nil
}

func (g *StripeGateway) GetCustomer(id string) (string, error) {
	params := url.Values{
		"query": []string{
			"metadata['internal_id']:'" + id + "'",
		},
	}

	res, err := g.request("GET", "/customers/search", params)
	if err != nil {
		return "", err
	}

	var customers struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&customers); err != nil {
		return "", err
	}

	if len(customers.Data) == 0 {
		return "", errors.New("no customer found")
	}

	if len(customers.Data) > 1 {
		return "", errors.New("found more than one customer")
	}

	return customers.Data[0].ID, nil
}

func (g *StripeGateway) request(method, endpoint string, params url.Values) (*http.Response, error) {
	endpoint, err := url.JoinPath(g.url, endpoint)
	if err != nil {
		return nil, err
	}

	var body io.Reader

	switch method {
	case "GET":
		endpoint = endpoint + "?" + params.Encode()
	case "POST", "PUT", "PATCH":
		body = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	return http.DefaultClient.Do(req)
}
