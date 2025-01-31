package parrot

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzRegisterPath(f *testing.F) {
	p := newParrot(f)

	baseRoute := Route{
		Method:             http.MethodGet,
		ResponseStatusCode: http.StatusOK,
		RawResponseBody:    "Squawk",
	}
	f.Add("/foo")
	f.Add("/foo/bar")
	f.Add("/*")
	f.Add("/foo/*")

	f.Fuzz(func(t *testing.T, path string) {
		route := baseRoute
		route.Path = path

		_ = p.Register(&route) // We just don't want panics
	})
}

func FuzzMethodAny(f *testing.F) {
	p := newParrot(f)

	route := &Route{
		Method:             MethodAny,
		Path:               "/any",
		ResponseStatusCode: http.StatusOK,
		RawResponseBody:    "Squawk",
	}

	err := p.Register(route)
	require.NoError(f, err)

	f.Add(http.MethodGet)
	f.Add(http.MethodPost)
	f.Add(http.MethodPut)
	f.Add(http.MethodPatch)
	f.Add(http.MethodDelete)
	f.Add(http.MethodOptions)
	f.Add(http.MethodConnect)
	f.Add(http.MethodTrace)

	f.Fuzz(func(t *testing.T, method string) {
		if !isValidMethod(method) {
			t.Skipf("invalid method '%s'", method)
		}
		resp, err := p.Call(method, route.Path)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode(), fmt.Sprintf("bad response code with method: '%s'", method))
		require.Equal(t, "Squawk", string(resp.Body()), fmt.Sprintf("bad response body with method: '%s'", method))
	})
}
