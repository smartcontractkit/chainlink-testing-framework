package client

import (
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"
)

// APIClient handles basic request sending logic and cookie handling
type APIClient struct {
	RC     *resty.Client
	Header http.Header
}

// NewAPIClient returns new basic resty client configured with an base URL
func NewAPIClient(baseURL string) *APIClient {
	rc := resty.New()
	rc.SetBaseURL(baseURL)
	return &APIClient{
		RC: rc,
	}
}

func (c *APIClient) WithHeader(header http.Header) *APIClient {
	c.Header = header
	return c
}

func (c *APIClient) Request(method,
	endpoint string,
	body interface{},
	obj interface{},
	expectedStatusCode int,
) (*resty.Response, error) {
	req := c.RC.R()
	req.Method = method
	req.URL = endpoint
	resp, err := req.
		SetHeaderMultiValues(c.Header).
		SetBody(body).
		SetResult(&obj).
		Send()
	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return resp, fmt.Errorf(
			"unexpected response code, got %d",
			resp.StatusCode(),
		)
	} else if resp.StatusCode() != expectedStatusCode {
		return resp, fmt.Errorf(
			"unexpected response code, got %d, expected 200\nURL: %s\nresponse received: %s",
			resp.StatusCode(),
			resp.Request.URL,
			resp.String(),
		)
	}
	return resp, err
}
