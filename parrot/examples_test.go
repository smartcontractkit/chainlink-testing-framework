package parrot_test

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
)

func ExampleServer() {
	// Create a new parrot instance with no logging
	p, err := parrot.Wake(parrot.WithLogLevel(zerolog.NoLevel))
	if err != nil {
		panic(err)
	}

	// Create a new route /test that will return a 200 status code with a text/plain response body of "Squawk"
	route := &parrot.Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: 200,
	}

	// Register the route with the parrot instance
	err = p.Register(route)
	if err != nil {
		panic(err)
	}

	// Call the route
	resp, err := p.Call(http.MethodGet, "/test")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// 200
	// Squawk
}
