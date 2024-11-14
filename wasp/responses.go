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

// NewResponses creates a new Responses that processes responses from the provided channel ch.
// It returns a pointer to a Responses instance.
func NewResponses(ch chan *Response) *Responses {
	return &Responses{ch}
}

// OK processes a successful resty.Response by sending a Response struct to the Responses channel.
// It records the duration of the response, associates it with the specified group, and includes the response body data.
func (m *Responses) OK(r *resty.Response, group string) {
	m.ch <- &Response{
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}

// Err sends a failed Response to the Responses channel with the provided error, group, and resty.Response.
// It marks the Response as failed, includes the error message, duration, group name, and response body.
func (m *Responses) Err(r *resty.Response, group string, err error) {
	m.ch <- &Response{
		Failed:   true,
		Error:    err.Error(),
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}
