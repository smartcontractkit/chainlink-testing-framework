package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	splunkToken               = "test-token"
	splunkEvent   SplunkEvent = "test"
	numberReports             = 3
	reportID                  = "123"
	testRunCount              = 15
	uniqueTests               = 18
)

func TestAggregateResultsSplunk(t *testing.T) {
	t.Parallel()

	srv := splunkServer(t)
	t.Cleanup(srv.Close)

	report, err := LoadAndAggregate("./testdata", WithReportID(reportID), WithSplunk(srv.URL, splunkToken, splunkEvent))
	require.NoError(t, err, "LoadAndAggregate failed")
	verifyAggregatedReport(t, report)
}

func TestAggregateResults(t *testing.T) {
	t.Parallel()

	report, err := LoadAndAggregate("./testdata", WithReportID(reportID))
	require.NoError(t, err, "LoadAndAggregate failed")
	verifyAggregatedReport(t, report)
}

func splunkServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Content-Type"), "unexpected content type")
		require.Equal(t, fmt.Sprintf("Splunk %s", splunkToken), r.Header.Get("Authorization"), "unexpected authorization header")

		// Figure out what kind of splunk data the payload is
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
			var report SplunkTestReport
			err := json.Unmarshal(bodyBytes, &report)
			require.NoError(t, err, "error parsing report data")
			require.NotNil(t, report, "error parsing report data")
			assert.Equal(t, SplunkSourceType, report.SourceType, "source type mismatch")
			assert.Equal(t, splunkEvent, report.Event.Event, "event mismatch")
			assert.Equal(t, reportID, report.Event.Data.ID, "ID mismatch")
		} else if eventType == string(Result) {
			var result SplunkTestResult
			err := json.Unmarshal(bodyBytes, &result)
			require.NoError(t, err, "error parsing results data")
			require.NotNil(t, result, "error parsing results data")
			assert.Equal(t, SplunkSourceType, result.SourceType, "source type mismatch")
			assert.Equal(t, splunkEvent, result.Event.Event, "event mismatch")
		} else {
			t.Errorf("unexpected splunk event type: %s", eventType)
		}

		w.WriteHeader(http.StatusOK)
	}))
}

func verifyAggregatedReport(t *testing.T, report *TestReport) {
	require.NotNil(t, report, "report is nil")
	assert.Equal(t, reportID, report.ID, "report ID mismatch")
	assert.Equal(t, uniqueTests, len(report.Results), "report results count mismatch")
	assert.Equal(t, testRunCount, report.TestRunCount, "report test run count mismatch")
	assert.Equal(t, false, report.RaceDetection, "race detection should be false")

	var (
		testFail, testSkipped, testPass TestResult
		testFailName                    = "TestFail"
		testSkippedName                 = "TestSkipped"
		testPassName                    = "TestPass"
	)
	for _, result := range report.Results {
		if result.TestName == testFailName {
			testFail = result
		}
		if result.TestName == testSkippedName {
			testSkipped = result
		}
		if result.TestName == testPassName {
			testPass = result
		}
	}

	t.Run("verify TestFail", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, testFailName, testFail.TestName, "TestFail not found")
		assert.False(t, testFail.Panic, "TestFail should not panic")
		assert.False(t, testFail.Skipped, "TestFail should not be skipped")
		assert.Equal(t, testRunCount, testFail.Runs, "TestFail should run every time")
		assert.Zero(t, testFail.Skips, "TestFail should not be skipped")
		assert.Equal(t, testRunCount, testFail.Failures, "TestFail should fail every time")
		assert.Len(t, testFail.Durations, testRunCount, "TestFail should have durations")
	})

	t.Run("verify TestSkipped", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, testSkippedName, testSkipped.TestName, "TestSkip not found")
		assert.False(t, testSkipped.Panic, "TestSkipped should not panic")
		assert.Zero(t, testSkipped.Runs, "TestSkipped should not pass")
		assert.True(t, testSkipped.Skipped, "TestSkipped should be skipped")
		assert.Equal(t, testRunCount, testSkipped.Skips, "TestSkipped should be skipped entirely")
		assert.Empty(t, testSkipped.Durations, "TestSkipped should not have durations")
	})

	t.Run("verify TestPass", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, testPassName, testPass.TestName, "TestPass not found")
		assert.False(t, testPass.Panic, "TestPass should not panic")
		assert.Equal(t, testRunCount, testPass.Runs, "TestPass should run every time")
		assert.False(t, testPass.Skipped, "TestPass should not be skipped")
		assert.Zero(t, testPass.Skips, "TestPass should not be skipped")
		assert.Equal(t, testRunCount, testPass.Successes, "TestPass should pass every time")
		assert.Len(t, testPass.Durations, testRunCount, "TestPass should have durations")
	})
}

func BenchmarkTestAggregateResults(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for i := 0; i < b.N; i++ {
		_, err := LoadAndAggregate("./testdata", WithReportID(reportID))
		if err != nil {
			b.Fatalf("LoadAndAggregate failed: %v", err)
		}
	}
}

func BenchmarkTestAggregateResultsSplunk(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadAndAggregate("./testdata", WithReportID(reportID), WithSplunk(srv.URL, splunkToken, "test"))
		if err != nil {
			b.Fatalf("LoadAndAggregate failed: %v", err)
		}
	}
}
