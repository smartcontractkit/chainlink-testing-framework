package wasp

import (
	"github.com/go-resty/resty/v2"
)

/* handy wrappers to use with resty in scenario (VU) tests */

const (
	CallGroupLabel = "call_group"
)

type Responses struct {
	ch chan *Response
}

// NewResponses creates a Responses instance using the provided channel.
// It enables concurrent processing and management of Response objects.
func NewResponses(ch chan *Response) *Responses {
	return &Responses{ch}
}

// OK sends a successful resty.Response along with its duration and group to the Responses channel.
// It is used to handle and process successful responses in a concurrent environment.
func (m *Responses) OK(r *resty.Response, group string) {
	m.ch <- &Response{
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}

// Err sends a failed response, including error details and response data, to the Responses channel.
// It is used to handle and propagate errors within the response processing workflow.
func (m *Responses) Err(r *resty.Response, group string, err error) {
	m.ch <- &Response{
		Failed:   true,
		Error:    err.Error(),
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}
