package testreporters_test

import (
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/integrations-framework/testreporters"
	"github.com/stretchr/testify/require"
)

func TestSoakReport(t *testing.T) {
	t.Parallel()

	reporter := &testreporters.OCRSoakTestReporter{
		Reports: map[string]*testreporters.OCRSoakTestReport{
			"report": {
				ContractAddress: "0x0",
				TotalRounds:     1,
			},
		},
	}
	reporter.Reports["report"].UpdateReport(time.Minute, 1)

	// Create local logs folder if there isn't one
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		err = os.MkdirAll("./logs", os.ModePerm)
		require.NoError(t, err)
	}

	err := reporter.WriteReport("./logs")
	require.NoError(t, err)

	// Cleanup
	err = os.RemoveAll("./logs")
	require.NoError(t, err)
}
