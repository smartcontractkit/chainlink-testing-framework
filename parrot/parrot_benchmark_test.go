package parrot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func BenchmarkRegisterRoute(b *testing.B) {
	saveFile := b.Name() + ".json"
	p, err := NewServer(WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)

	defer benchmarkCleanup(b, p, saveFile)

	routes := make([]*Route, b.N)
	for i := 0; i < b.N; i++ {
		routes[i] = &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    "Benchmark Response",
			ResponseStatusCode: http.StatusOK,
		}
		err := p.Register(routes[i])
		require.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := p.Register(routes[i])
		require.NoError(b, err)
	}
	b.StopTimer()
}

func BenchmarkRouteResponse(b *testing.B) {
	saveFile := b.Name() + ".json"
	p, err := NewServer(WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)

	defer benchmarkCleanup(b, p, saveFile)

	routes := make([]*Route, b.N)
	for i := 0; i < b.N; i++ {
		routes[i] = &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    "Benchmark Response",
			ResponseStatusCode: http.StatusOK,
		}
		err := p.Register(routes[i])
		require.NoError(b, err)
	}

	route := &Route{
		Method:             "GET",
		Path:               "/bench",
		RawResponseBody:    "Benchmark Response",
		ResponseStatusCode: http.StatusOK,
	}
	err = p.Register(route)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Call(route.Method, route.Path)
		require.NoError(b, err)
	}
	b.StopTimer()
}

func BenchmarkGetRoutes(b *testing.B) {
	saveFile := b.Name() + ".json"
	p, err := NewServer(WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)

	defer benchmarkCleanup(b, p, saveFile)

	routes := make([]*Route, b.N)
	for i := 0; i < b.N; i++ {
		routes[i] = &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    "Benchmark Response",
			ResponseStatusCode: http.StatusOK,
		}
		err := p.Register(routes[i])
		require.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		routes := p.Routes()
		require.Len(b, routes, b.N)
	}
	b.StopTimer()
}

func BenchmarkSave(b *testing.B) {
	var (
		routes    = []*Route{}
		recorders = []string{}
		saveFile  = "bench_save_routes.json"
	)

	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    fmt.Sprintf("Squawk %d", i),
			ResponseStatusCode: http.StatusOK,
		})
		recorders = append(recorders, fmt.Sprintf("http://recorder%d", i))
	}
	p, err := NewServer(WithRoutes(routes), WithRecorders(recorders...), WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)

	defer benchmarkCleanup(b, p, saveFile)

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.save()
		require.NoError(b, err)
	}
	b.StopTimer()
}

func BenchmarkLoad(b *testing.B) {
	var (
		routes    = []*Route{}
		recorders = []string{}
		saveFile  = "bench_load_routes.json"
	)
	b.Cleanup(func() {
		os.Remove(saveFile)
	})

	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    fmt.Sprintf("Squawk %d", i),
			ResponseStatusCode: http.StatusOK,
		})
		recorders = append(recorders, fmt.Sprintf("http://recorder%d", i))
	}
	p, err := NewServer(WithRoutes(routes), WithRecorders(recorders...), WithLogLevel(zerolog.Disabled), WithSaveFile(saveFile))
	require.NoError(b, err, "error waking parrot")

	defer benchmarkCleanup(b, p, saveFile)

	err = p.save()
	require.NoError(b, err, "error saving routes")

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.load()
		require.NoError(b, err)
	}
	b.StopTimer()
}

func benchmarkCleanup(b *testing.B, p *Server, saveFile string) {
	b.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err := p.Shutdown(ctx)
	cancel()
	require.NoError(b, err, "error shutting down parrot")
	p.WaitShutdown()
	os.Remove(saveFile)
}
