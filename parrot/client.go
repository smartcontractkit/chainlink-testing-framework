package parrot

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Client interacts with a parrot server
type Client struct {
	restyClient *resty.Client
}

// NewClient creates a new client for a parrot server running at the given url.
func NewClient(url string) *Client {
	restyC := resty.New()
	restyC.SetBaseURL(url)
	return &Client{
		restyClient: restyC,
	}
}

// Health returns the health of the server
func (c *Client) Healthy() (bool, error) {
	resp, err := c.restyClient.R().Get(HealthRoute)
	if err != nil {
		return false, err
	}
	return resp.StatusCode() == http.StatusOK, nil
}

// Routes returns all the routes registered on the server
func (c *Client) Routes() ([]*Route, error) {
	routes := []*Route{}
	resp, err := c.restyClient.R().SetResult(&routes).Get(RoutesRoute)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get routes, got %d status code: %s", resp.StatusCode(), string(resp.Body()))
	}
	return routes, nil
}

// CallRoute calls a route on the server
func (c *Client) CallRoute(method, path string) (*resty.Response, error) {
	return c.restyClient.R().Execute(method, path)
}

// RegisterRoute registers a route on the server
func (c *Client) RegisterRoute(route *Route) error {
	resp, err := c.restyClient.R().SetBody(route).Post(RoutesRoute)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to register route, got %d status code: %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

// DeleteRoute deletes a route on the server
func (c *Client) DeleteRoute(route *Route) error {
	resp, err := c.restyClient.R().SetBody(route).Delete(RoutesRoute)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed to delete route, got %d status code: %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

// RegisterRecorder registers a recorder on the server
func (c *Client) RegisterRecorder(recorder *Recorder) error {
	resp, err := c.restyClient.R().SetBody(recorder).Post(RecorderRoute)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to register recorder, got %d status code: %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

// Recorders returns all the recorders registered on the server
func (c *Client) Recorders() ([]string, error) {
	recorders := []string{}
	resp, err := c.restyClient.R().SetResult(&recorders).Get(RecorderRoute)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get recorders, got %d status code: %s", resp.StatusCode(), string(resp.Body()))
	}
	return recorders, nil
}
