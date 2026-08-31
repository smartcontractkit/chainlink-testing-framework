package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	if code != 0 {
		t.Fatalf("code = %d, want 0; stderr = %q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "rule0000006a") {
		t.Errorf("stdout = %q, want it to list rule0000006a", out)
	}
	if !strings.Contains(out, "Example No Gateways Available") {
		t.Errorf("stdout = %q, want it to list the rule title", out)
	}
	if !strings.Contains(out, "grafana-managed") {
		t.Errorf("stdout = %q, want it to name the rule kind", out)
	}
}

func TestRunList_UnsupportedVersion(t *testing.T) {
	srv := grafanaTestServer(t, "12.5.0")
	t.Setenv("GRAFANA_URL", srv.URL)
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	code := run([]string{"list"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "12.5.0") {
		t.Fatalf("stderr = %q, want it to name the unsupported version", stderr.String())
	}
}

func TestRunList_RejectsArgs(t *testing.T) {
	t.Setenv("GRAFANA_URL", "http://example.invalid")
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	code := run([]string{"list", "extra"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}
