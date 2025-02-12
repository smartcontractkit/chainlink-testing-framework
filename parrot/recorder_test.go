package parrot

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestResponseWriterRecorder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		responseFunc       http.HandlerFunc
		expectedRespCode   int
		expectedRespBody   string
		expectedRespHeader http.Header
	}{
		{
			name: "good response",
			responseFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("Squawk"))
				require.NoError(t, err, "error writing response")
			},
			expectedRespCode: http.StatusOK,
			expectedRespBody: "Squawk",
			expectedRespHeader: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
		{
			name: "error response",
			responseFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Squawk", http.StatusInternalServerError)
			},
			expectedRespCode: http.StatusInternalServerError,
			expectedRespBody: "Squawk\n", // http.Error adds a newline
			expectedRespHeader: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			writerRecorder := newResponseWriterRecorder(recorder)
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			handler := http.HandlerFunc(tc.responseFunc)
			handler.ServeHTTP(writerRecorder, req)

			actualResp := recorder.Result()
			recordedResp := writerRecorder.Result()
			t.Cleanup(func() {
				_ = actualResp.Body.Close()
				_ = recordedResp.Body.Close()
			})

			actualBody, err := io.ReadAll(actualResp.Body)
			require.NoError(t, err, "error reading actual response body")
			recordedBody, err := io.ReadAll(recordedResp.Body)
			require.NoError(t, err, "error reading recorded response body")

			assert.Equal(t, tc.expectedRespCode, actualResp.StatusCode, "actual response has unexpected status code")
			assert.Equal(t, tc.expectedRespCode, recordedResp.StatusCode, "recorded response has unexpected status code")
			assert.Equal(t, tc.expectedRespBody, string(actualBody), "actual response has unexpected body")
			assert.Equal(t, tc.expectedRespBody, string(recordedBody), "recorded response has unexpected body")
		})
	}
}

func TestRecorder(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	recorder, err := NewRecorder()
	require.NoError(t, err, "error creating recorder")
	t.Cleanup(func() {
		_ = recorder.Close()
	})

	err = p.Record(recorder.URL())
	require.NoError(t, err, "error recording parrot")
	t.Cleanup(func() {
		_ = recorder.Close()
	})

	route := &Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}
	err = p.Register(route)
	require.NoError(t, err, "error registering route")

	responseCount := 5

	var eg errgroup.Group
	for i := 0; i < responseCount; i++ {
		eg.Go(func() error {
			resp, err := p.Call(http.MethodGet, route.Path)
			if err != nil {
				return err
			}
			if resp.StatusCode() != route.ResponseStatusCode {
				return fmt.Errorf("unexpected status code calling '/test': %d", resp.StatusCode())
			}
			if string(resp.Body()) != route.RawResponseBody {
				return fmt.Errorf("unexpected body calling '/test': %s", string(resp.Body()))
			}
			return nil
		})
	}

	require.NoError(t, eg.Wait(), "error calling parrot")

	for i := 0; i < responseCount; i++ {
		select {
		case recordedRouteCall := <-recorder.Record():
			assert.Equal(t, route.ID(), recordedRouteCall.RouteID, "recorded response has unexpected route ID")
			assert.Equal(t, http.StatusOK, recordedRouteCall.Response.StatusCode, "recorded response has unexpected status code")
			assert.Equal(t, route.RawResponseBody, string(recordedRouteCall.Response.Body), "recorded response has unexpected body")
			assert.Equal(t, route.Path, recordedRouteCall.Request.URL.Path, "recorded request has unexpected path")
			assert.Equal(t, http.MethodGet, recordedRouteCall.Request.Method, "recorded request has unexpected method")
		case err := <-recorder.Err():
			require.NoError(t, err, "error recording route call")
		case <-time.After(time.Second):
			require.Fail(t, "timed out waiting for recorded route call")
		}
	}
}

func TestMultipleRecorders(t *testing.T) {
	t.Parallel()

	p := newParrot(t)

	var (
		numRecorders = 10
		numCalls     = 5
	)
	recorders := make([]*Recorder, numRecorders)
	for i := 0; i < numRecorders; i++ {
		recorder, err := NewRecorder()
		require.NoError(t, err, "error creating recorder")
		recorders[i] = recorder
	}
	t.Cleanup(func() {
		for _, recorder := range recorders {
			require.NoError(t, recorder.Close())
		}
	})

	for _, recorder := range recorders {
		err := p.Record(recorder.URL())
		require.NoError(t, err, "error recording parrot")
	}

	route := &Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}
	err := p.Register(route)
	require.NoError(t, err, "error registering route")

	var eg errgroup.Group
	for i := 0; i < numCalls; i++ {
		eg.Go(func() error {
			resp, err := p.Call(http.MethodGet, route.Path)
			if err != nil {
				return err
			}
			if resp.StatusCode() != route.ResponseStatusCode {
				return fmt.Errorf("unexpected status code calling '/test': %d", resp.StatusCode())
			}
			if string(resp.Body()) != route.RawResponseBody {
				return fmt.Errorf("unexpected body calling '/test': %s", string(resp.Body()))
			}
			return nil
		})
	}

	require.NoError(t, eg.Wait(), "error calling parrot")

	for _, recorder := range recorders {
		for i := 0; i < numCalls; i++ {
			select {
			case recordedRouteCall := <-recorder.Record():
				assert.Equal(t, route.ID(), recordedRouteCall.RouteID, "recorded response has unexpected route ID for recorder %d", i)
				assert.Equal(t, http.StatusOK, recordedRouteCall.Response.StatusCode, "recorded response has unexpected status code for recorder %d", i)
				assert.Equal(t, "Squawk", string(recordedRouteCall.Response.Body), "recorded response has unexpected body for recorder %d", i)
				assert.Equal(t, "/test", recordedRouteCall.Request.URL.Path, "recorded request has unexpected path for recorder %d", i)
				assert.Equal(t, http.MethodGet, recordedRouteCall.Request.Method, "recorded request has unexpected method for recorder %d", i)
			case err := <-recorder.Err():
				require.NoError(t, err, "error recording route call")
			case <-time.After(time.Second):
				require.Fail(t, "timed out waiting for recorder %d", i)
			}
		}
	}
}
