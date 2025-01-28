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
	p, err := Wake(WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)

	defer func() { // Cleanup
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := p.Shutdown(ctx)
		cancel()
		require.NoError(b, err, "error shutting down parrot")
		p.WaitShutdown()
		os.Remove(saveFile)
	}()

	route := &Route{
		Method:             "GET",
		Path:               "/bench",
		RawResponseBody:    "Benchmark Response",
		ResponseStatusCode: http.StatusOK,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := p.Register(route)
		require.NoError(b, err)
	}
	b.StopTimer()
}

func BenchmarkRouteResponse(b *testing.B) {
	saveFile := b.Name() + ".json"
	p, err := Wake(WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)

	defer func() { // Cleanup
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := p.Shutdown(ctx)
		cancel()
		require.NoError(b, err, "error shutting down parrot")
		p.WaitShutdown()
		os.Remove(saveFile)
	}()

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

func BenchmarkSave(b *testing.B) {
	var (
		routes   = []*Route{}
		saveFile = "bench_save_routes.json"
	)

	for i := 0; i < 1000; i++ {
		routes = append(routes, &Route{
			Method:             "GET",
			Path:               fmt.Sprintf("/bench%d", i),
			RawResponseBody:    fmt.Sprintf("Squawk %d", i),
			ResponseStatusCode: http.StatusOK,
		})
	}
	p, err := Wake(WithRoutes(routes), WithLogLevel(testLogLevel), WithSaveFile(saveFile))
	require.NoError(b, err)
	defer func() { // Cleanup
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err = p.Shutdown(ctx)
		cancel()
		require.NoError(b, err, "error shutting down parrot")
		p.WaitShutdown()
		os.Remove(saveFile)
	}()

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.save()
		require.NoError(b, err)
	}
	b.StopTimer()
}

func BenchmarkLoad(b *testing.B) {
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
			ResponseStatusCode: http.StatusOK,
		})
	}
	p, err := Wake(WithRoutes(routes), WithLogLevel(zerolog.Disabled), WithSaveFile(saveFile))
	require.NoError(b, err, "error waking parrot")
	defer func() { // Cleanup
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err = p.Shutdown(ctx)
		cancel()
		require.NoError(b, err, "error shutting down parrot")
		p.WaitShutdown()
		os.Remove(saveFile)
	}()

	err = p.save()
	require.NoError(b, err, "error saving routes")

	b.ResetTimer() // Start measuring time
	for i := 0; i < b.N; i++ {
		err := p.load()
		require.NoError(b, err)
	}
	b.StopTimer()
}
