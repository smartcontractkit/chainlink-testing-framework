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

// NewResponses creates a new instance of Responses, initialized with the provided channel for Response objects. 
// This function returns a pointer to the Responses struct, which can be used to manage and process responses 
// sent through the specified channel.
func NewResponses(ch chan *Response) *Responses {
	return &Responses{ch}
}

// OK sends a response object containing the duration, group, and data 
// from the provided resty.Response to a channel. 
// It captures the time taken for the request, associates it with the 
// specified group, and includes the response body data. 
// This function is intended for use in handling successful responses 
// in a concurrent environment.
func (m *Responses) OK(r *resty.Response, group string) {
	m.ch <- &Response{
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}

// Err sends an error response to a channel, encapsulating details about the failure. 
// It constructs a Response object that includes the error message, the duration of the request, 
// the specified group, and the response body. 
// This function does not return a value but communicates the failure through the channel.
func (m *Responses) Err(r *resty.Response, group string, err error) {
	m.ch <- &Response{
		Failed:   true,
		Error:    err.Error(),
		Duration: r.Time(),
		Group:    group,
		Data:     r.Body(),
	}
}
