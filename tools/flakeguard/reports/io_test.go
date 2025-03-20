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
	splunkToken   = "test-token"
	splunkEvent   = "test"
	reportID      = "123"
	totalTestRuns = 270
	testRunCount  = 15
	uniqueTests   = 19
)

func TestAggregateResultFilesSplunk(t *testing.T) {
	t.Parallel()

	var (
		reportRequestsReceived int
		resultRequestsReceived int
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Content-Type"), "unexpected content type")
		require.Equal(t, fmt.Sprintf("Splunk %s", splunkToken), r.Header.Get("Authorization"), "unexpected authorization header")

		// Figure out what kind of splunk data the payload is
		bodyBytes, err := io.ReadAll(r.Body)
		require.NoError(t, err, "error reading request")
		defer r.Body.Close()

		report := SplunkTestReport{}
		err = json.Unmarshal(bodyBytes, &report)
		if err == nil {
			require.Equal(t, SplunkSourceType, report.SourceType, "source type mismatch")
			require.Equal(t, Report, report.Event.Type, "event type mismatch")
			require.Equal(t, splunkEvent, report.Event.Event, "event mismatch")
			require.Equal(t, reportID, report.Event.Data.ID, "report ID mismatch")
			require.False(t, report.Event.Incomplete, "report should not be incomplete")
			require.NotNil(t, report.Event.Data.SummaryData, "report summary data is nil")
			require.Len(t, report.Event.Data.Results, 0, "shouldn't send all result data to splunk")
			reportRequestsReceived++
		} else {
			results, err := unBatchSplunkResults(bodyBytes)
			require.NoError(t, err, "error parsing splunk results data")
			require.NotNil(t, results, "expected some results")
			require.NotZero(t, len(results), "expected some results")
			for _, result := range results {
				require.Equal(t, SplunkSourceType, result.SourceType, "source type mismatch")
				require.Equal(t, Result, result.Event.Type, "event type mismatch")
				require.Equal(t, splunkEvent, result.Event.Event, "event mismatch")
				resultRequestsReceived++
			}
		}

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	results, err := LoadAndAggregate("./testdata")
	require.NoError(t, err, "LoadAndAggregate failed")

	report, err := NewTestReport(results, WithReportID(reportID))
	require.NoError(t, err, "NewTestReport failed")

	err = SendTestReportToSplunk(srv.URL, splunkToken, splunkEvent, report)
	require.NoError(t, err, "SendReportToSplunk failed")
	verifyAggregatedReport(t, report)
	assert.Equal(t, 1, reportRequestsReceived, "unexpected number of report requests")
	assert.Equal(t, uniqueTests, resultRequestsReceived, "unexpected number of report requests")
}

func TestAggregateResultFiles(t *testing.T) {
	t.Parallel()

	results, err := LoadAndAggregate("./testdata")
	require.NoError(t, err, "LoadAndAggregate failed")

	report, err := NewTestReport(results, WithReportID(reportID))
	require.NoError(t, err, "NewTestReport failed")

	verifyAggregatedReport(t, report)
}

func verifyAggregatedReport(t *testing.T, report TestReport) {
	require.NotNil(t, report, "report is nil")
	require.Equal(t, reportID, report.ID, "report ID mismatch")
	require.Equal(t, uniqueTests, len(report.Results), "report results count mismatch")
	require.Equal(t, totalTestRuns, report.SummaryData.TotalRuns, "report test total runs mismatch")
	require.Equal(t, false, report.RaceDetection, "race detection should be false")

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
		require.Equal(t, testFailName, testFail.TestName, "TestFail not found")
		assert.False(t, testFail.Panic, "TestFail should not panic")
		assert.False(t, testFail.Skipped, "TestFail should not be skipped")
		assert.Equal(t, testRunCount, testFail.Runs, "TestFail should run every time")
		assert.Zero(t, testFail.Skips, "TestFail should not be skipped")
		assert.Equal(t, testRunCount, testFail.Failures, "TestFail should fail every time")
		assert.Len(t, testFail.Durations, testRunCount, "TestFail should have durations")
	})

	t.Run("verify TestSkipped", func(t *testing.T) {
		require.Equal(t, testSkippedName, testSkipped.TestName, "TestSkip not found")
		assert.False(t, testSkipped.Panic, "TestSkipped should not panic")
		assert.Zero(t, testSkipped.Runs, "TestSkipped should not pass")
		assert.True(t, testSkipped.Skipped, "TestSkipped should be skipped")
		assert.Equal(t, testRunCount, testSkipped.Skips, "TestSkipped should be skipped entirely")
		assert.Empty(t, testSkipped.Durations, "TestSkipped should not have durations")
	})

	t.Run("verify TestPass", func(t *testing.T) {
		require.Equal(t, testPassName, testPass.TestName, "TestPass not found")
		assert.False(t, testPass.Panic, "TestPass should not panic")
		assert.Equal(t, testRunCount, testPass.Runs, "TestPass should run every time")
		assert.False(t, testPass.Skipped, "TestPass should not be skipped")
		assert.Zero(t, testPass.Skips, "TestPass should not be skipped")
		assert.Equal(t, testRunCount, testPass.Successes, "TestPass should pass every time")
		assert.Len(t, testPass.Durations, testRunCount, "TestPass should have durations")
	})
}

func BenchmarkAggregateResultFiles(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for i := 0; i < b.N; i++ {
		_, err := LoadAndAggregate("./testdata")
		require.NoError(b, err, "LoadAndAggregate failed")
	}
}

func BenchmarkAggregateResultFilesSplunk(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	b.Cleanup(srv.Close)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadAndAggregate("./testdata")
		require.NoError(b, err, "LoadAndAggregate failed")
	}
}
