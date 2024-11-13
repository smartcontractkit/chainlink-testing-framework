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

// NewResponses creates a new Responses instance with the provided channel for Response pointers.
// It returns a pointer to the Responses struct initialized with the given channel.
func NewResponses(ch chan *Response) *Responses {
	return &Responses{ch}
}

// OK sends a successful HTTP response to the channel m.ch.
// It constructs a Response with the duration of the request,
// the specified group, and the response body data.
func (m *Responses) OK(r *resty.Response, group string) {
	m.ch <- &Response{
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}

// Err processes a failed HTTP response and sends it to the channel m.ch.
// It constructs a Response with failure status, error message, response duration,
// group identifier, and response body.
func (m *Responses) Err(r *resty.Response, group string, err error) {
	m.ch <- &Response{
		Failed:   true,
		Error:    err.Error(),
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}
