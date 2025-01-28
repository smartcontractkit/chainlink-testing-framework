package parrot

import (
	"errors"
	"fmt"
)

var (
	ErrNilRoute        = errors.New("route is nil")
	ErrNoMethod        = errors.New("no method specified")
	ErrInvalidPath     = errors.New("invalid path")
	ErrNoResponse      = errors.New("route must have a handler or some response")
	ErrOnlyOneResponse = errors.New("route can only have one response type")
	ErrResponseMarshal = errors.New("unable to marshal response body to JSON")
	ErrRouteNotFound   = errors.New("route not found")

	ErrNoRecorderURL      = errors.New("no recorder URL specified")
	ErrInvalidRecorderURL = errors.New("invalid recorder URL")
	ErrRecorderNotFound   = errors.New("recorder not found")

	ErrServerShutdown = errors.New("parrot is already asleep")
)

// Custom error type to help add more detail to base errors
type dynamicError struct {
	Base  error  // Base error for comparison
	Extra string // Dynamic context (e.g., method name)
}

func (e *dynamicError) Error() string {
	return fmt.Sprintf("%s: %s", e.Base.Error(), e.Extra)
}

func (e *dynamicError) Unwrap() error {
	return e.Base
}

func newDynamicError(base error, detail string) error {
	return &dynamicError{
		Base:  base,
		Extra: detail,
	}
}
