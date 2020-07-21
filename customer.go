package magento2

import (
	"fmt"
	"net/http"
)

// Customer represents a customer entity in Magento
type Customer struct {
	ID              int64             `json:"id"`
	GroupID         int64             `json:"group_id"`
	Firstname       string            `json:"firstname"`
	Lastname        string            `json:"lastname"`
	Email           string            `json:"email"`
	CreatedIn       string            `json:"created_in"`
	StoreID         int64             `json:"store_id"`
	WebsiteID       int64             `json:"website_id"`
	DefaultBilling  string            `json:"default_billing"`
	DefaultShipping string            `json:"default_shipping"`
	Addresses       []CustomerAddress `json:"addresses"`
}

// CustomerAddress represents an address for a customer in Magento
type CustomerAddress struct {
	CountryID       string   `json:"country_id"`
	Street          []string `json:"street"`
	Company         string   `json:"company"`
	Telephone       string   `json:"telephone"`
	Postcode        string   `json:"postcode"`
	Firstname       string   `json:"firstname"`
	Lastname        string   `json:"lastname"`
	City            string   `json:"city"`
	DefaultBilling  bool     `json:"default_billing"`  // Whether to use this address as default for billing address
	DefaultShipping bool     `json:"default_shipping"` // Whether to use this address as default for shipping address
}

const (
	customersEndpointPrefix = "/customers"
)

// CustomersClient performs requests specifically against the /customers
// endpoint
type CustomersClient struct {
	c *Client
}

// Customers returns a CustomersClient
func (c *Client) Customers() *CustomersClient {
	return &CustomersClient{c}
}

// Search performs a request against the /search endpoint for customers
func (c *CustomersClient) Search(criteria *SearchCriteria) ([]Customer, error) {
	endpoint := customersEndpointPrefix + "/search"
	req, err := c.c.newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("could not call client.newRequest(): %w", err)
	}

	v := req.URL.Query()
	criteria.SetQueryParams(v)
	req.URL.RawQuery = v.Encode()

	var respEnvelope struct {
		Items []Customer `json:"items"`
	}

	err = c.c.do(req, http.StatusOK, &respEnvelope)
	if err != nil {
		return nil, err
	}

	return respEnvelope.Items, nil
}
