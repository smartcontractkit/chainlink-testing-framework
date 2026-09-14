package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const rulerBody = `{
  "Example-Zone-A": [
    {
      "name": "Gateway",
      "rules": [
        {
          "for": "5m",
          "grafana_alert": {
            "title": "Example No Gateways Available",
            "uid": "rule0000006a",
            "namespace_uid": "folder0000006",
            "intervalSeconds": 60,
            "no_data_state": "OK",
            "exec_err_state": "OK",
            "is_paused": false
          }
        }
      ]
    }
  ]
}`

func healthBody(version string) string {
	return fmt.Sprintf(`{"database":"ok","version":%q,"commit":"abc123"}`, version)
}

func grafanaTestServer(t *testing.T, version string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/health":
			_, _ = w.Write([]byte(healthBody(version)))
		case "/api/ruler/grafana/api/v1/rules":
			_, _ = w.Write([]byte(rulerBody))
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestRunList_HappyPath(t *testing.T) {
	srv := grafanaTestServer(t, "13.1.0")
	t.Setenv("GRAFANA_URL", srv.URL)
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	code := run([]string{"list"}, &stdout, &stderr)
	require.Equal(t, 0, code)
	out := stdout.String()
	require.Contains(t, out, "rule0000006a")
	require.Contains(t, out, "Example No Gateways Available")
	require.Contains(t, out, "grafana-managed")
}

func TestRunList_UnsupportedVersion(t *testing.T) {
	srv := grafanaTestServer(t, "12.5.0")
	t.Setenv("GRAFANA_URL", srv.URL)
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	code := run([]string{"list"}, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), "12.5.0")
}

func TestRunList_RejectsArgs(t *testing.T) {
	t.Setenv("GRAFANA_URL", "http://example.invalid")
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	code := run([]string{"list", "extra"}, &stdout, &stderr)
	require.Equal(t, 2, code)
}
