package parrot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWake(t *testing.T) {
	t.Parallel()
	p, err := Wake()
	require.NoError(t, err, "error waking parrot")
	require.NotNil(t, p)
}

func TestRegisterRoute(t *testing.T) {
	t.Parallel()
	p, err := Wake()
	require.NoError(t, err)

	route := &Route{
		Method:              http.MethodGet,
		Path:                "/test",
		ResponseBody:        "{\"message\":\"Test successful\"}",
		ResponseStatusCode:  200,
		ResponseContentType: "application/json",
	}

	err = p.Register(route)
	require.NoError(t, err, "error registering route")
}

func TestRegisteredRoute(t *testing.T) {
	t.Parallel()

	p, err := Wake()
	require.NoError(t, err, "error waking parrot")

	routes := []*Route{
		{
			Method:              http.MethodPost,
			Path:                "/hello",
			ResponseBody:        "{\"message\":\"Hello, world!\"}",
			ResponseStatusCode:  200,
			ResponseContentType: "application/json",
		},
		{
			Method:              http.MethodPost,
			Path:                "/goodbye",
			ResponseBody:        "{\"message\":\"Goodbye, world!\"}",
			ResponseStatusCode:  201,
			ResponseContentType: "application/json",
		},
	}

	for _, r := range routes {
		route := r
		t.Run(route.Method+":"+route.Path, func(t *testing.T) {
			t.Parallel()

			err = p.Register(route)
			require.NoError(t, err, "error registering route")

			resp, err := p.Call(route.Method, route.Path)
			require.NoError(t, err, "error calling parrot")
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, route.ResponseStatusCode)
			assert.Equal(t, resp.Header.Get("Content-Type"), route.ResponseContentType)
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, string(body), route.ResponseBody)
			resp.Body.Close()
		})
	}
}

func TestUnregisteredRoute(t *testing.T) {
	t.Parallel()

	p, err := Wake()
	require.NoError(t, err, "error waking parrot")

	resp, err := p.Call(http.MethodGet, "/unregistered")
	require.Error(t, err, "error calling parrot")
	require.NotNil(t, resp, "response should not be nil")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSaveLoad(t *testing.T) {
	t.Parallel()

	saveFile := "save_test.json"
	p, err := Wake(WithSaveFile(saveFile))
	require.NoError(t, err, "error waking parrot")

	routes := []*Route{
		{
			Method:              "GET",
			Path:                "/hello",
			ResponseBody:        "{\"message\":\"Hello, world!\"}",
			ResponseStatusCode:  200,
			ResponseContentType: "application/json",
		},
		{
			Method:              "Post",
			Path:                "/goodbye",
			ResponseBody:        "{\"message\":\"Goodbye, world!\"}",
			ResponseStatusCode:  201,
			ResponseContentType: "application/json",
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

		assert.Equal(t, resp.StatusCode, route.ResponseStatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), route.ResponseContentType)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, string(body), route.ResponseBody)
		resp.Body.Close()
	}
}

func BenchmarkRegisterRoute(b *testing.B) {
	p, err := Wake()
	require.NoError(b, err)

	route := &Route{
		Method:              "GET",
		Path:                "/bench",
		ResponseBody:        "Benchmark Response",
		ResponseStatusCode:  200,
		ResponseContentType: "text/plain",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := p.Register(route)
		require.NoError(b, err)
	}
}

func BenchmarkRouteResponse(b *testing.B) {
	p, err := Wake()
	require.NoError(b, err)

	route := &Route{
		Method:              "GET",
		Path:                "/bench",
		ResponseBody:        "Benchmark Response",
		ResponseStatusCode:  200,
		ResponseContentType: "text/plain",
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
		saveFile = "bench_routes.json"
	)

	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:              "GET",
			Path:                fmt.Sprintf("/bench%d", i),
			ResponseBody:        fmt.Sprintf("{\"message\":\"Response %d\"}", i),
			ResponseStatusCode:  200,
			ResponseContentType: "text/plain",
		})
	}
	p, err := Wake(WithRoutes(routes))
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
		saveFile = "bench_routes.json"
	)
	b.Cleanup(func() {
		os.Remove(saveFile)
	})
	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:              "GET",
			Path:                fmt.Sprintf("/bench%d", i),
			ResponseBody:        fmt.Sprintf("{\"message\":\"Response %d\"}", i),
			ResponseStatusCode:  200,
			ResponseContentType: "text/plain",
		})
	}
	p, err := Wake(WithRoutes(routes))
	require.NoError(b, err, "error waking parrot")
	err = p.save()
	require.NoError(b, err, "error saving routes")

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.load()
		require.NoError(b, err)
	}
}
