package parrot_test

import (
	"fmt"
	"io"
	"net/http"

	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
)

func ExampleServer() {
	p, err := parrot.Wake()
	if err != nil {
		panic(err)
	}

	route := &parrot.Route{
		Method:              http.MethodGet,
		Path:                "/test",
		RawResponseBody:     "Squawk",
		ResponseStatusCode:  200,
		ResponseContentType: "text/plain",
	}

	err = p.Register(route)
	if err != nil {
		panic(err)
	}

	resp, err := p.Call(http.MethodGet, "/test")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
	// text/plain
	// Squawk
}
