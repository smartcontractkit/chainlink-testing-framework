package parrot

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testLogLevel = zerolog.NoLevel

func TestMain(m *testing.M) {
	testLogLevelFlag := ""
	flag.StringVar(&testLogLevelFlag, "testLogLevel", "", "a zerolog log level to use for tests")
	flag.Parse()
	var err error
	testLogLevel, err = zerolog.ParseLevel(testLogLevelFlag)
	if err != nil {
		fmt.Println("error parsing test log level:", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestWake(t *testing.T) {
	t.Parallel()

	p, err := Wake(WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")
	require.NotNil(t, p)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	p, err := Wake(WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")

	testCases := []struct {
		name  string
		route *Route
	}{
		// {
		// 	name: "get route",
		// 	route: &Route{
		// 		Method:             http.MethodGet,
		// 		Path:               "/hello",
		// 		RawResponseBody:    "Squawk",
		// 		ResponseStatusCode: 200,
		// 	},
		// },
		// {
		// 	name: "json route",
		// 	route: &Route{
		// 		Method:             http.MethodGet,
		// 		Path:               "/json",
		// 		ResponseBody:       map[string]any{"message": "Squawk"},
		// 		ResponseStatusCode: 200,
		// 	},
		// },
		// {
		// 	name: "post route",
		// 	route: &Route{
		// 		Method:             http.MethodPost,
		// 		Path:               "/post",
		// 		RawResponseBody:    "Squawk",
		// 		ResponseStatusCode: 201,
		// 	},
		// },
		// {
		// 	name: "put route",
		// 	route: &Route{
		// 		Method:             http.MethodPut,
		// 		Path:               "/put",
		// 		RawResponseBody:    "Squawk",
		// 		ResponseStatusCode: 200,
		// 	},
		// },
		// {
		// 	name: "delete route",
		// 	route: &Route{
		// 		Method:             http.MethodDelete,
		// 		Path:               "/delete",
		// 		RawResponseBody:    "Squawk",
		// 		ResponseStatusCode: 200,
		// 	},
		// },
		// {
		// 	name: "patch route",
		// 	route: &Route{
		// 		Method:             http.MethodPatch,
		// 		Path:               "/patch",
		// 		RawResponseBody:    "Squawk",
		// 		ResponseStatusCode: 200,
		// 	},
		// },
		{
			name: "error route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/error",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: 500,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := p.Register(tc.route)
			require.NoError(t, err, "error registering route")

			resp, err := p.Call(tc.route.Method, tc.route.Path)
			require.NoError(t, err, "error calling parrot")
			defer resp.Body.Close()

			assert.Equal(t, tc.route.ResponseStatusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			if tc.route.ResponseBody != nil {
				jsonBody, err := json.Marshal(tc.route.ResponseBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(jsonBody), string(body))
			} else {
				assert.Equal(t, tc.route.RawResponseBody, string(body))
			}
			resp.Body.Close()
		})
	}
}

func TestIsValidPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		paths []string
		valid bool
	}{
		{
			name:  "valid paths",
			paths: []string{"/hello"},
			valid: true,
		},
		{
			name:  "no register",
			paths: []string{"/register", "/register/", "/register//", "/register/other_stuff"},
			valid: false,
		},
		{
			name:  "invalid paths",
			paths: []string{"", "/", " ", " /", "/ ", " / ", "invalid path"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, path := range tc.paths {
				valid := isValidPath(path)
				assert.Equal(t, tc.valid, valid)
			}
		})
	}
}

func TestBadRegisterRoute(t *testing.T) {
	t.Parallel()

	p, err := Wake(WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")

	testCases := []struct {
		name  string
		err   error
		route *Route
	}{
		{
			name:  "nil route",
			err:   ErrNilRoute,
			route: nil,
		},
		{
			name: "no method",
			err:  ErrNoMethod,
			route: &Route{
				Path:               "/hello",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: 200,
			},
		},
		{
			name: "no path",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				RawResponseBody:    "Squawk",
				ResponseStatusCode: 200,
			},
		},
		{
			name: "base path",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: 200,
			},
		},
		{
			name: "invalid path",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "invalid path",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: 200,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := p.Register(tc.route)
			require.Error(t, err, "expected error registering route")
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestUnregisteredRoute(t *testing.T) {
	t.Parallel()

	p, err := Wake(WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")

	resp, err := p.Call(http.MethodGet, "/unregistered")
	require.NoError(t, err, "error calling parrot")
	require.NotNil(t, resp, "response should not be nil")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUnregister(t *testing.T) {
	t.Parallel()

	p, err := Wake(WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")

	route := &Route{
		Method:             http.MethodPost,
		Path:               "/hello",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: 200,
	}

	err = p.Register(route)
	require.NoError(t, err, "error registering route")

	resp, err := p.Call(route.Method, route.Path)
	require.NoError(t, err, "error calling parrot")

	assert.Equal(t, resp.StatusCode, route.ResponseStatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, route.RawResponseBody, string(body))
	resp.Body.Close()

	err = p.Unregister(route.ID())
	require.NoError(t, err, "error unregistering route")

	resp, err = p.Call(route.Method, route.Path)
	require.NoError(t, err, "error calling parrot")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSaveLoad(t *testing.T) {
	t.Parallel()

	saveFile := "save_test.json"
	p, err := Wake(WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(t, err, "error waking parrot")

	routes := []*Route{
		{
			Method:             "GET",
			Path:               "/hello",
			RawResponseBody:    "Squawk",
			ResponseStatusCode: 200,
		},
		{
			Method:             "Post",
			Path:               "/goodbye",
			RawResponseBody:    "Squeak",
			ResponseStatusCode: 201,
		},
	}

	for _, route := range routes {
		err = p.Register(route)
		require.NoError(t, err, "error registering route")
	}

	t.Cleanup(func() {
		os.Remove(saveFile)
	})

	err = p.save()
	require.NoError(t, err)

	require.FileExists(t, saveFile)
	err = p.load()
	require.NoError(t, err)

	for _, route := range routes {
		resp, err := p.Call(route.Method, route.Path)
		require.NoError(t, err, "error calling parrot")

		assert.Equal(t, route.ResponseStatusCode, resp.StatusCode, "unexpected status code for route %s", route.ID())
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, route.RawResponseBody, string(body))
		resp.Body.Close()
	}
}

func BenchmarkRegisterRoute(b *testing.B) {
	p, err := Wake(WithLogLevel(zerolog.Disabled))
	require.NoError(b, err)

	route := &Route{
		Method:             "GET",
		Path:               "/bench",
		RawResponseBody:    "Benchmark Response",
		ResponseStatusCode: 200,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := p.Register(route)
		require.NoError(b, err)
	}
}

func BenchmarkRouteResponse(b *testing.B) {
	p, err := Wake(WithLogLevel(zerolog.Disabled))
	require.NoError(b, err)

	route := &Route{
		Method:             "GET",
		Path:               "/bench",
		RawResponseBody:    "Benchmark Response",
		ResponseStatusCode: 200,
	}
	err = p.Register(route)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Call(route.Method, route.Path)
		require.NoError(b, err)
	}
}

func BenchmarkSaveRoutes(b *testing.B) {
	var (
		routes   = []*Route{}
		saveFile = "bench_save_routes.json"
	)

	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    fmt.Sprintf("Squawk %d", i),
			ResponseStatusCode: 200,
		})
	}
	p, err := Wake(WithRoutes(routes), WithLogLevel(zerolog.Disabled), WithSaveFile(saveFile))
	require.NoError(b, err)

	b.Cleanup(func() {
		os.Remove(saveFile)
	})

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.save()
		require.NoError(b, err)
	}
}

func BenchmarkLoadRoutes(b *testing.B) {
	var (
		routes   = []*Route{}
		saveFile = "bench_load_routes.json"
	)
	b.Cleanup(func() {
		os.Remove(saveFile)
	})

	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    fmt.Sprintf("Squawk %d", i),
			ResponseStatusCode: 200,
		})
	}
	p, err := Wake(WithRoutes(routes), WithLogLevel(zerolog.Disabled), WithSaveFile(saveFile))
	require.NoError(b, err, "error waking parrot")
	err = p.save()
	require.NoError(b, err, "error saving routes")

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.load()
		require.NoError(b, err)
	}
}
