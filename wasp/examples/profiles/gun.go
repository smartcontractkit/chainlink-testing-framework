package main

import (
	"github.com/go-resty/resty/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type ExampleGun struct {
	target string
	client *resty.Client
	Data   []string
}

func NewExampleHTTPGun(target string) *ExampleGun {
	return &ExampleGun{
		client: resty.New(),
		target: target,
		Data:   make([]string, 0),
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *ExampleGun) Call(l *wasp.Generator) *wasp.Response {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		return &wasp.Response{Data: result, Error: err.Error()}
	}
	if r.Status() != "200 OK" {
		return &wasp.Response{Data: result, Error: "not 200", Failed: true}
	}
	return &wasp.Response{Data: result}
}
