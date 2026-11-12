package chipingressset

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v72/github"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, handler http.HandlerFunc) *github.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	base, err := url.Parse(srv.URL + "/")
	require.NoError(t, err)

	client := github.NewClient(nil)
	client.BaseURL = base
	return client
}

func noRetryWait(t *testing.T) {
	t.Helper()
	old := retryWait
	retryWait = func(context.Context, time.Duration) error { return nil }
	t.Cleanup(func() { retryWait = old })
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprint(w, body)
}

func TestResolveRefSHATag(t *testing.T) {
	commitCalls := 0
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/ref/tags/"):
			writeJSON(w, http.StatusOK, `{"ref":"refs/tags/v1.2.3","object":{"sha":"tag-sha"}}`)
		case strings.Contains(r.URL.Path, "/commits/"):
			commitCalls++
			http.NotFound(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	sha, err := resolveRefSHA(context.Background(), client, "o", "r", "v1.2.3")
	require.NoError(t, err)
	require.Equal(t, "tag-sha", sha)
	require.Zero(t, commitCalls)
}

func TestResolveRefSHABranch(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/ref/tags/"):
			http.NotFound(w, r)
		case strings.Contains(r.URL.Path, "/git/ref/heads/"):
			writeJSON(w, http.StatusOK, `{"ref":"refs/heads/main","object":{"sha":"branch-sha"}}`)
		default:
			http.NotFound(w, r)
		}
	})

	sha, err := resolveRefSHA(context.Background(), client, "o", "r", "main")
	require.NoError(t, err)
	require.Equal(t, "branch-sha", sha)
}

func TestResolveRefSHAShortCommitSHA(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/git/ref/"):
			http.NotFound(w, r)
		case strings.Contains(r.URL.Path, "/commits/fa1d775e2138"):
			writeJSON(w, http.StatusOK, `{"sha":"fa1d775e21381de1c33baa11e9b378935c8e3a0e"}`)
		default:
			http.NotFound(w, r)
		}
	})

	sha, err := resolveRefSHA(context.Background(), client, "o", "r", "fa1d775e2138")
	require.NoError(t, err)
	require.Equal(t, "fa1d775e21381de1c33baa11e9b378935c8e3a0e", sha)
}

func TestResolveRefSHANotFound(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	_, err := resolveRefSHA(context.Background(), client, "o", "r", "deadbeef")
	require.Error(t, err)
	require.Contains(t, err.Error(), `ref "deadbeef" not found`)
}

func TestResolveRefSHASurfacesServerError(t *testing.T) {
	noRetryWait(t)
	attempts := 0
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/commits/") {
			attempts++
			writeJSON(w, http.StatusInternalServerError, `{"message":"Server Error"}`)
			return
		}
		http.NotFound(w, r)
	})

	_, err := resolveRefSHA(context.Background(), client, "o", "r", "fa1d775e2138")
	require.Error(t, err)
	require.Contains(t, err.Error(), "500")
	require.NotContains(t, err.Error(), "not found")
	require.Equal(t, commitLookupAttempts, attempts)
}

func TestResolveRefSHARetriesTransientError(t *testing.T) {
	noRetryWait(t)
	attempts := 0
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/commits/") {
			attempts++
			if attempts == 1 {
				writeJSON(w, http.StatusBadGateway, `{"message":"Bad Gateway"}`)
				return
			}
			writeJSON(w, http.StatusOK, `{"sha":"full-sha"}`)
			return
		}
		http.NotFound(w, r)
	})

	sha, err := resolveRefSHA(context.Background(), client, "o", "r", "fa1d775e2138")
	require.NoError(t, err)
	require.Equal(t, "full-sha", sha)
	require.Equal(t, 2, attempts)
}
