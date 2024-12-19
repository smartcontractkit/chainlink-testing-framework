package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	splunkToken   = "test-token"
	numberReports = 3
	reportID      = "123"
	testsRun      = 15
)

func TestAggregateResultsSplunk(t *testing.T) {
	srv := splunkServer(t)
	t.Cleanup(srv.Close)

	report, err := LoadAndAggregate("./testdata", WithReportID(reportID), WithSplunk(srv.URL, splunkToken, "test"))
	require.NoError(t, err, "LoadAndAggregate failed")
	require.NotNil(t, report, "report is nil")
}

func TestAggregateResults(t *testing.T) {
	report, err := LoadAndAggregate("./testdata", WithReportID(reportID))
	require.NoError(t, err, "LoadAndAggregate failed")
	require.NotNil(t, report, "report is nil")
	assert.Equal(t, reportID, report.ID, "report ID mismatch")
	assert.Equal(t, testsRun, len(report.Results), "report results count mismatch")
	assert.Equal(t, testsRun, report.TestRunCount, "report test run count mismatch")
}

func splunkServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, fmt.Sprintf("Splunk %s", splunkToken), r.Header.Get("Authorization"))

		// Figure out what the payload is
		bodyBytes, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		var payload map[string]any
		err = json.Unmarshal(bodyBytes, &payload)
		require.NoError(t, err, "error parsing splunk event data")
		require.NotNil(t, payload, "error parsing splunk event data")
		require.NotNil(t, payload["event"], "unable to find event while parsing splunk data")
		event := payload["event"].(map[string]any)
		require.NotNil(t, event, "error parsing splunk event data")
		require.NotNil(t, event["type"], "unable to find inner event type while parsing splunk data")
		eventType := event["type"].(string)
		require.NotNil(t, eventType, "error parsing splunk event type")

		if eventType == string(Report) {
			var report TestReport
			err := json.Unmarshal(bodyBytes, &report)
			require.NoError(t, err, "error parsing report data")
			require.NotNil(t, report, "error parsing report data")
			assert.Equal(t, reportID, report.ID, "report ID mismatch")
			assert.Equal(t, testsRun, len(report.Results), "report results count mismatch")
			assert.Equal(t, testsRun, report.TestRunCount, "report test run count mismatch")
		} else if eventType == string(Result) {
			var result SplunkTestResult
			err := json.Unmarshal(bodyBytes, &result)
			require.NoError(t, err, "error parsing results data")
			require.NotNil(t, result, "error parsing results data")
		} else {
			t.Errorf("unexpected event type: %s", eventType)
		}

		w.WriteHeader(http.StatusOK)
	}))
}
