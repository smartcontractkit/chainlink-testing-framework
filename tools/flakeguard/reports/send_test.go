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

	report, err := LoadAndAggregate("./testdata", WithReportID(reportID), WithSplunk(srv.URL, splunkToken, splunkEvent))
	require.NoError(t, err, "LoadAndAggregate failed")
	verifyAggregatedReport(t, report)
	assert.Equal(t, 1, reportRequestsReceived, "unexpected number of report requests")
	assert.Equal(t, uniqueTests, resultRequestsReceived, "unexpected number of report requests")
}
