package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to register a new route
func registerRoute(tb testing.TB, route Route) {
	tb.Helper()

	body, _ := json.Marshal(route)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	registerRouteHandler(rec, req)

	resp := rec.Result()
	tb.Cleanup(func() {
		resp.Body.Close()
	})
	require.Equal(tb, http.StatusCreated, resp.StatusCode)
}

func TestRegisterRoute(t *testing.T) {
	t.Parallel()

	route := Route{
		Method:      "GET",
		Path:        "/test",
		Response:    "{\"message\":\"Test successful\"}",
		StatusCode:  200,
		ContentType: "application/json",
	}

	registerRoute(t, route)
}

func TestRegisteredRoute(t *testing.T) {
	t.Parallel()

	routes := []Route{
		{
			Method:      "GET",
			Path:        "/hello",
			Response:    "{\"message\":\"Hello, world!\"}",
			StatusCode:  200,
			ContentType: "application/json",
		},
		{
			Method:      "Post",
			Path:        "/goodbye",
			Response:    "{\"message\":\"Goodbye, world!\"}",
			StatusCode:  201,
			ContentType: "application/json",
		},
	}
	for _, r := range routes {
		route := r
		t.Run(route.Method+":"+route.Path, func(t *testing.T) {
			t.Parallel()

			registerRoute(t, route)

			req := httptest.NewRequest(route.Method, route.Path, nil)
			rec := httptest.NewRecorder()
			dynamicHandler(rec, req)
			resp := rec.Result()

			assert.Equal(t, resp.StatusCode, route.StatusCode)
			assert.Equal(t, resp.Header.Get("Content-Type"), route.ContentType)
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, string(body), route.Response)
			resp.Body.Close()
		})
	}
}

func TestUnregisteredRoute(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/unregistered", nil)
	rec := httptest.NewRecorder()

	dynamicHandler(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d but got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestSaveLoad(t *testing.T) {
	routes := []Route{
		{
			Method:      "GET",
			Path:        "/hello",
			Response:    "{\"message\":\"Hello, world!\"}",
			StatusCode:  200,
			ContentType: "application/json",
		},
		{
			Method:      "Post",
			Path:        "/goodbye",
			Response:    "{\"message\":\"Goodbye, world!\"}",
			StatusCode:  201,
			ContentType: "application/json",
		},
	}

	for _, route := range routes {
		registerRoute(t, route)
	}

	filename := "test_routes.json"
	t.Cleanup(func() {
		os.Remove(filename)
	})

	err := save()
	require.NoError(t, err)
	require.FileExists(t, filename)

	err = load()
	require.NoError(t, err)

	for _, route := range routes {
		req := httptest.NewRequest(route.Method, route.Path, nil)
		rec := httptest.NewRecorder()

		dynamicHandler(rec, req)
		resp := rec.Result()

		assert.Equal(t, resp.StatusCode, route.StatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), route.ContentType)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, string(body), route.Response)
		resp.Body.Close()
	}
}

func BenchmarkRegisterRoute(b *testing.B) {
	route := Route{
		Method:      "GET",
		Path:        "/bench",
		Response:    "Benchmark Response",
		StatusCode:  200,
		ContentType: "text/plain",
	}

	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(route)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		registerRouteHandler(rec, req)
	}
}

func BenchmarkRouteResponse(b *testing.B) {
	route := Route{
		Method:      "GET",
		Path:        "/bench",
		Response:    "Benchmark Response",
		StatusCode:  200,
		ContentType: "text/plain",
	}
	registerRoute(b, route)
	req := httptest.NewRequest(route.Method, route.Path, nil)
	rec := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dynamicHandler(rec, req)
	}
}

func BenchmarkSaveRoutes(b *testing.B) {
	for i := 0; i < 1000; i++ {
		route := Route{
			Method:      "GET",
			Path:        fmt.Sprintf("/bench%d", i),
			Response:    fmt.Sprintf("{\"message\":\"Response %d\"}", i),
			StatusCode:  200,
			ContentType: "text/plain",
		}
		routes[route.Method+":"+route.Path] = route
	}

	filename := "bench_routes.json"
	b.Cleanup(func() {
		os.Remove(filename)
	})

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		if err := save(); err != nil {
			b.Fatalf("SaveRoutes failed: %v", err)
		}
	}
}

func BenchmarkLoadRoutes(b *testing.B) {
	filename := "bench_routes.json"
	b.Cleanup(func() {
		os.Remove(filename)
	})

	benchmarkRoutes := make(map[string]Route)
	for i := 0; i < 1000; i++ {
		route := Route{
			Method:      "GET",
			Path:        fmt.Sprintf("/bench%d", i),
			Response:    fmt.Sprintf("{\"message\":\"Response %d\"}", i),
			StatusCode:  200,
			ContentType: "text/plain",
		}
		benchmarkRoutes[route.Method+":"+route.Path] = route
	}
	data, _ := json.MarshalIndent(benchmarkRoutes, "", "  ")
	if err := os.WriteFile(filename, data, 0644); err != nil { //nolint:gosec
		b.Fatalf("Failed to write benchmark file: %v", err)
	}

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		if err := load(); err != nil {
			b.Fatalf("Load failed: %v", err)
		}
	}
}
