package magento2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	contentType      = "application/json"
	defaultUserAgent = "go-magento2"
)

// Client is a Magento 2 REST API client
type Client struct {
	baseURL    string
	apiKey     string
	userAgent  string
	httpClient *http.Client
}

type opt func(*Client)

// APIError represents the error return from the Magento 2 REST API
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

// NewClient constructs and returns a Magento 2 REST API client, the default
// values can be overridden using the opt functions
func NewClient(baseURL string, opts ...opt) *Client {
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: defaultUserAgent,
	}

	for _, f := range opts {
		f(c)
	}

	return c
}

// WithHTTPClientOpt allows using a custom http.Client for sending the requests
func WithHTTPClientOpt(httpClient *http.Client) opt {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithAPIKeyOpt allows setting an API key which is used in the Authorization
// header for the requests
func WithAPIKeyOpt(apiKey string) opt {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithUserAgent allows setting the value for the User-Agent header in requests
func WithUserAgent(userAgent string) opt {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// newRequest will marshal request's payload and set appropriate headers before
// returning a *http.Request
func (c *Client) newRequest(
	method, url string,
	payload interface{},
) (*http.Request, error) {
	var req *http.Request
	var err error

	var reqPayload *bytes.Buffer
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, fmt.Errorf("could not marshal the payload to json: %w", err)
		}

		reqPayload = bytes.NewBuffer(b)
	}

	if reqPayload != nil {
		req, err = http.NewRequest(method, c.baseURL+url, reqPayload)
	} else {
		req, err = http.NewRequest(method, c.baseURL+url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create new http.Request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+c.apiKey)
	}

	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Accept", contentType)
	if payload != nil {
		req.Header.Add("Content-Type", contentType)
	}

	return req, nil
}

// do will make the actual requests agains the Magento 2 REST API. It will
// unmarshal any expected response body and also unmarshal into APIErr in case
// of errors returned from the API
func (c *Client) do(
	req *http.Request,
	expectedStatusCode int,
	expectedRespBody interface{},
) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request %v: %w", req, err)
	}
	defer resp.Body.Close()

	readBody := func(into interface{}) error {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %w", err)
		}

		err = json.Unmarshal(body, &into)
		if err != nil {
			return fmt.Errorf("could not unmarshal body into expected type: %w", err)
		}

		return nil
	}

	if expectedStatusCode != resp.StatusCode {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
		}
		err := readBody(apiErr)
		if err != nil {
			return err
		}

		return apiErr
	}

	if expectedRespBody != nil {
		return readBody(expectedRespBody)
	}

	return nil
}
