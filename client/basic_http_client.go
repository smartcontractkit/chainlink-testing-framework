package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// ErrNotFound Error for not found
var ErrNotFound = errors.New("unexpected response code, got 404")

// ErrUnprocessableEntity Error for and unprocessable entity
var ErrUnprocessableEntity = errors.New("unexpected response code, got 422")

// BasicHTTPClient handles basic http sending logic and cookie handling
type BasicHTTPClient struct {
	BaseURL    string
	HttpClient *http.Client
	Cookies    []*http.Cookie
	Header     http.Header
}

// NewBasicHTTPClient returns new basic http client configured with an base URL
func NewBasicHTTPClient(c *http.Client, baseURL string) *BasicHTTPClient {
	return &BasicHTTPClient{
		BaseURL:    baseURL,
		HttpClient: c,
		Cookies:    make([]*http.Cookie, 0),
	}
}

func (em *BasicHTTPClient) do(
	method,
	endpoint string,
	body interface{},
	obj interface{},
	expectedStatusCode int,
) (*http.Response, error) {
	b, err := json.Marshal(body)
	if body != nil && err != nil {
		return nil, err
	}
	return em.doRaw(method, endpoint, b, obj, expectedStatusCode)
}

func (em *BasicHTTPClient) doRaw(
	method,
	endpoint string,
	body []byte, obj interface{},
	expectedStatusCode int,
) (*http.Response, error) {
	client := em.HttpClient
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", em.BaseURL, endpoint),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	for _, cookie := range em.Cookies {
		req.AddCookie(cookie)
	}

	req.Header = em.Header

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"error while reading response: %v\nURL: %s\nresponse received: %s",
			err,
			em.BaseURL,
			string(b),
		)
	}
	if resp.StatusCode == http.StatusNotFound {
		return resp, ErrNotFound
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		return resp, ErrUnprocessableEntity
	} else if resp.StatusCode != expectedStatusCode {
		return resp, fmt.Errorf(
			"unexpected response code, got %d, expected 200\nURL: %s\nresponse received: %s",
			resp.StatusCode,
			em.BaseURL,
			string(b),
		)
	}

	if obj == nil {
		return resp, err
	}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, fmt.Errorf(
			"error while unmarshaling response: %v\nURL: %s\nresponse received: %s",
			err,
			em.BaseURL,
			string(b),
		)
	}
	return resp, err
}
