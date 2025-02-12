package parrot

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

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

func TestHealthy(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	healthCount := 0
	targetCount := 3

	ticker := time.NewTicker(time.Millisecond * 10)
	timeout := time.NewTimer(time.Second)
	t.Cleanup(func() {
		ticker.Stop()
		timeout.Stop()
	})

	for {
		select {
		case <-ticker.C:
			if err := p.Healthy(); err == nil {
				healthCount++
			}
			if healthCount >= targetCount {
				return
			}
		case <-timeout.C:
			require.GreaterOrEqual(t, targetCount, healthCount, "parrot never became healthy")
		}
	}
}

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	testCases := []struct {
		name  string
		route *Route
	}{
		{
			name: "get route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/hello",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "json route",
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/json",
				ResponseBody:       map[string]any{"message": "Squawk"},
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "post route",
			route: &Route{
				Method:             http.MethodPost,
				Path:               "/post",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: 201,
			},
		},
		{
			name: "put route",
			route: &Route{
				Method:             http.MethodPut,
				Path:               "/put",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "delete route",
			route: &Route{
				Method:             http.MethodDelete,
				Path:               "/delete",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "patch route",
			route: &Route{
				Method:             http.MethodPatch,
				Path:               "/patch",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
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
			err := p.Register(tc.route)
			require.NoError(t, err, "error registering route")

			resp, err := p.Call(tc.route.Method, tc.route.Path)
			require.NoError(t, err, "error calling parrot")

			assert.Equal(t, tc.route.ResponseStatusCode, resp.StatusCode())
			if tc.route.ResponseBody != nil {
				jsonBody, err := json.Marshal(tc.route.ResponseBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(jsonBody), string(resp.Body()))
			} else {
				assert.Equal(t, tc.route.RawResponseBody, string(resp.Body()))
			}
		})
	}
}

func TestGetRoutes(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	routes := []*Route{
		{
			Method:             http.MethodGet,
			Path:               "/hello",
			RawResponseBody:    "Squawk",
			ResponseStatusCode: http.StatusOK,
		},
		{
			Method:             http.MethodPost,
			Path:               "/goodbye",
			RawResponseBody:    "Squeak",
			ResponseStatusCode: 201,
		},
	}

	for _, route := range routes {
		err := p.Register(route)
		require.NoError(t, err, "error registering route")
	}

	registeredRoutes := p.Routes()
	require.Len(t, registeredRoutes, len(routes))
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
			name:  "no protected paths",
			paths: []string{HealthRoute, RoutesRoute, RecorderRoute, fmt.Sprintf("%s/%s", RoutesRoute, "route-id"), fmt.Sprintf("%s/%s", HealthRoute, "recorder-id"), fmt.Sprintf("%s/%s", RecorderRoute, "recorder-id")},
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

func TestPreRegisterRoutesRecorders(t *testing.T) {
	t.Parallel()

	routes := []*Route{
		{
			Method:             http.MethodGet,
			Path:               "/hello",
			RawResponseBody:    "Squawk",
			ResponseStatusCode: http.StatusOK,
		},
		{
			Method:             http.MethodPost,
			Path:               "/goodbye",
			RawResponseBody:    "Squeak",
			ResponseStatusCode: 201,
		},
	}
	recorders := []string{
		"http://localhost:8080",
		"http://localhost:8081",
	}

	saveFile := t.Name() + ".json"
	p, err := NewServer(WithSaveFile(saveFile), WithRecorders(recorders...), WithRoutes(routes), WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")

	t.Cleanup(func() {
		err := p.Shutdown(context.Background())
		assert.NoError(t, err, "error shutting down parrot")
		p.WaitShutdown()
		os.Remove(saveFile)
	})

	foundRoutes := p.Routes()
	require.Len(t, foundRoutes, len(routes))
	foundRecorders := p.Recorders()
	require.Len(t, foundRecorders, len(recorders))
}

func TestCustomLogFile(t *testing.T) {
	t.Parallel()

	logFile := t.Name() + ".log"
	saveFile := t.Name() + ".json"
	p, err := NewServer(WithLogFile(logFile), WithSaveFile(saveFile), WithLogLevel(zerolog.InfoLevel))
	require.NoError(t, err, "error waking parrot")

	t.Cleanup(func() {
		err := p.Shutdown(context.Background())
		assert.NoError(t, err, "error shutting down parrot")
		p.WaitShutdown()
		os.Remove(logFile)
		os.Remove(saveFile)
	})

	// Call a route to generate some logs
	route := &Route{
		Method:             http.MethodGet,
		Path:               "/hello",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}
	err = p.Register(route)
	require.NoError(t, err, "error registering route")

	_, err = p.Call(route.Method, route.Path)
	require.NoError(t, err, "error calling parrot")

	require.FileExists(t, logFile, "expected log file to exist")
	logData, err := os.ReadFile(logFile)
	require.NoError(t, err, "error reading log file")
	require.Contains(t, string(logData), "GET:/hello", "expected log file to contain route call")
}

func TestBadRegisterRoute(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

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
			err:  ErrInvalidMethod,
			route: &Route{
				Path:               "/hello",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "no path",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "base path",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "invalid path",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "invalid path",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "no response",
			err:  ErrNoResponse,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/hello",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "invalid url",
			err:  ErrInvalidPath,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "http://example.com",
				RawResponseBody:    "Squawk",
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "multiple responses",
			err:  ErrOnlyOneResponse,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/hello",
				RawResponseBody:    "Squawk",
				ResponseBody:       map[string]any{"message": "Squawk"},
				ResponseStatusCode: http.StatusOK,
			},
		},
		{
			name: "bad JSON",
			err:  ErrResponseMarshal,
			route: &Route{
				Method:             http.MethodGet,
				Path:               "/json",
				ResponseBody:       map[string]any{"message": make(chan int)},
				ResponseStatusCode: http.StatusOK,
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

func TestBadRecorder(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	err := p.Record("")
	require.ErrorIs(t, err, ErrNoRecorderURL, "expected error recording parrot")

	err = p.Record("invalid url")
	require.ErrorIs(t, err, ErrInvalidRecorderURL, "expected error recording parrot")
}

func TestUnregisteredRoute(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	resp, err := p.Call(http.MethodGet, "/unregistered")
	require.NoError(t, err, "error calling parrot")
	require.NotNil(t, resp, "response should not be nil")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
}

func TestDelete(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	route := &Route{
		Method:             http.MethodPost,
		Path:               "/hello",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}

	err := p.Register(route)
	require.NoError(t, err, "error registering route")

	resp, err := p.Call(route.Method, route.Path)
	require.NoError(t, err, "error calling parrot")

	assert.Equal(t, resp.StatusCode(), route.ResponseStatusCode)
	assert.Equal(t, route.RawResponseBody, string(resp.Body()))

	p.Delete(route)

	resp, err = p.Call(route.Method, route.Path)
	require.NoError(t, err, "error calling parrot")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
}

func TestSaveLoad(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	routes := []*Route{
		{
			Method:             http.MethodGet,
			Path:               "/hello",
			RawResponseBody:    "Squawk",
			ResponseStatusCode: http.StatusOK,
		},
		{
			Method:             http.MethodPost,
			Path:               "/goodbye",
			RawResponseBody:    "Squeak",
			ResponseStatusCode: 201,
		},
	}

	recorders := []string{ // Dummy recorder URLs
		"http://localhost:8080",
		"http://localhost:8081",
	}

	for _, route := range routes {
		err := p.Register(route)
		require.NoError(t, err, "error registering route")
	}

	for _, recorder := range recorders {
		err := p.Record(recorder)
		require.NoError(t, err, "error recording parrot")
	}

	err := p.save()
	require.NoError(t, err)

	require.FileExists(t, t.Name()+".json")
	err = p.load()
	require.NoError(t, err)

	for _, route := range routes {
		resp, err := p.Call(route.Method, route.Path)
		require.NoError(t, err, "error calling parrot")

		assert.Equal(t, route.ResponseStatusCode, resp.StatusCode(), "unexpected status code for route %s", route.ID())
		assert.Equal(t, route.RawResponseBody, string(resp.Body()))
	}

	registeredRecorders := p.Recorders()
	require.Len(t, registeredRecorders, len(recorders), "unexpected number of recorders")
}

func TestShutDown(t *testing.T) {
	t.Parallel()

	fileName := t.Name() + ".json"
	p, err := NewServer(WithSaveFile(fileName), WithLogLevel(testLogLevel))
	require.NoError(t, err, "error waking parrot")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = p.Shutdown(ctx)
	require.NoError(t, err, "error shutting down parrot")

	p.WaitShutdown() // Wait for shutdown to complete

	_, err = p.Call(http.MethodGet, "/hello")
	require.ErrorIs(t, err, ErrServerShutdown, "expected error calling parrot after shutdown")

	err = p.Record("http://localhost:8080")
	require.ErrorIs(t, err, ErrServerShutdown, "expected error recording parrot after shutdown")

	testRoute := &Route{
		Method:             http.MethodGet,
		Path:               "/hello",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}
	err = p.Register(testRoute)
	require.ErrorIs(t, err, ErrServerShutdown, "expected error registering route after shutdown")

	err = p.Shutdown(context.Background())
	require.ErrorIs(t, err, ErrServerShutdown, "expected error shutting down parrot after shutdown")
}

func TestJSONLogger(t *testing.T) {
	t.Parallel()

	logFileName := t.Name() + ".log"
	fileName := t.Name() + ".json"
	p, err := NewServer(WithSaveFile(fileName), WithLogLevel(zerolog.DebugLevel), WithLogFile(logFileName), WithJSONLogs(), DisableConsoleLogs())
	require.NoError(t, err, "error waking parrot")
	t.Cleanup(func() {
		os.Remove(fileName)
		os.Remove(logFileName)
	})

	route := &Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}

	err = p.Register(route)
	assert.NoError(t, err, "error registering route")

	_, err = p.Call(route.Method, route.Path)
	assert.NoError(t, err, "error calling parrot")

	err = p.Shutdown(context.Background())
	assert.NoError(t, err, "error shutting down parrot")
	p.WaitShutdown()
	require.FileExists(t, logFileName, "expected log file to exist")
	logFile, err := os.Open(logFileName)
	require.NoError(t, err, "error opening log file")
	logs, err := io.ReadAll(logFile)
	require.NoError(t, err, "error reading log file")
	require.NotNil(t, logs, "expected logs to be read from file")
	require.NotEmpty(t, logs, "expected logs to be written to file")
	require.Contains(t, string(logs), fmt.Sprintf(`"Route ID":"%s"`, route.ID()), "expected log file to contain route call in JSON format")
}

func newParrot(tb testing.TB) *Server {
	tb.Helper()

	logFileName := tb.Name() + ".log"
	saveFileName := tb.Name() + ".json"
	p, err := NewServer(WithSaveFile(saveFileName), WithLogFile(logFileName), WithLogLevel(testLogLevel))
	require.NoError(tb, err, "error waking parrot")
	tb.Cleanup(func() {
		err := p.Shutdown(context.Background())
		assert.NoError(tb, err, "error shutting down parrot")
		p.WaitShutdown() // Wait for shutdown to complete and file to be written
		os.Remove(saveFileName)
		os.Remove(logFileName)
	})
	return p
}
