package blockchain

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsRetryableFaucetErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"context canceled", context.Canceled, false},
		{"context deadline exceeded", context.DeadlineExceeded, false},
		{"4xx client error not retryable", &faucetStatusError{status: http.StatusBadRequest}, false},
		{"404 not retryable", &faucetStatusError{status: http.StatusNotFound}, false},
		{"401 unauthorized not retryable", &faucetStatusError{status: http.StatusUnauthorized}, false},
		{"429 rate limited retryable", &faucetStatusError{status: http.StatusTooManyRequests}, true},
		{"500 server error retryable", &faucetStatusError{status: http.StatusInternalServerError}, true},
		{"503 service unavailable retryable", &faucetStatusError{status: http.StatusServiceUnavailable}, true},
		{"bare error retryable (transport-classified)", errors.New("connection reset by peer"), true},
		{"io EOF retryable (transport-classified)", io.EOF, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isRetryableFaucetErr(tt.err))
		})
	}
}

// faucetStub is a configurable /gas handler that can fail the first N requests with
// either an HTTP 5xx or a connection reset (hijack+close), then succeed.
type faucetStub struct {
	failN      int32 // number of requests to fail
	resetMode  bool  // true: hijack+close (transport reset); false: 503
	hits       atomic.Int32
	successOK  atomic.Bool
}

func (f *faucetStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := f.hits.Add(1)
	if n <= f.failN {
		if f.resetMode {
			// Hijack and close the TCP conn mid-response → client sees a
			// connection reset / EOF (the readiness race we are fixing).
			if hj, ok := w.(http.Hijacker); ok {
				if conn, _, err := hj.Hijack(); err == nil {
					_ = conn.Close()
					return
				}
			}
			// Fallback if hijack is unavailable: force a transport-level abort.
			panic("faucetStub: hijack unavailable, cannot simulate connection reset")
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	f.successOK.Store(true)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"task":"ok"}`))
}

// shortCtx keeps the retry budget tight so a misbehaving stub can't hang the test.
// fundAccount wraps this in a 2min child, but the earlier parent deadline wins.
func shortCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestFundAccount_RetriesOn503ThenSucceeds(t *testing.T) {
	t.Parallel()

	stub := &faucetStub{failN: 2, resetMode: false}
	srv := httptest.NewServer(stub)
	t.Cleanup(srv.Close)

	err := fundAccount(shortCtx(t), srv.URL, "0xdeadbeef")
	require.NoError(t, err)
	assert.Equal(t, int32(3), stub.hits.Load(), "should fail twice then succeed on third")
	assert.True(t, stub.successOK.Load(), "success path should have run")
}

func TestFundAccount_RetriesOnConnectionResetThenSucceeds(t *testing.T) {
	t.Parallel()

	stub := &faucetStub{failN: 2, resetMode: true}
	srv := httptest.NewServer(stub)
	t.Cleanup(srv.Close)

	err := fundAccount(shortCtx(t), srv.URL, "0xdeadbeef")
	require.NoError(t, err)
	assert.Equal(t, int32(3), stub.hits.Load(), "should absorb two resets then succeed")
	assert.True(t, stub.successOK.Load(), "success path should have run")
}

func TestFundAccount_FastFailsOnNonRetryable4xx(t *testing.T) {
	t.Parallel()

	// A 4xx must stop immediately instead of burning the retry budget. A 5xx stub would
	// loop many times; this stub always returns 400 and must be hit exactly once.
	stub := &badRequestStub{}
	srv := httptest.NewServer(stub)
	t.Cleanup(srv.Close)

	err := fundAccount(shortCtx(t), srv.URL, "0xdeadbeef")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fund account via Sui faucet")
	assert.Contains(t, err.Error(), "status 400")
	assert.Equal(t, int32(1), stub.hits.Load(), "non-retryable 4xx must not be retried")
}

// badRequestStub always returns 400 to verify non-retryable fast-fail.
type badRequestStub struct {
	hits atomic.Int32
}

func (b *badRequestStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.hits.Add(1)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"error":"bad recipient"}`))
}
