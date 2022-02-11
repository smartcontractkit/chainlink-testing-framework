package testreporters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
)

// KeeperBlockTimeTestReporter enables reporting
type KeeperBlockTimeTestReporter struct {
	Reports                        []KeeperBlockTimeTestReport `json:"reports"`
	ReportMutex                    sync.Mutex
	AttemptedChainlinkTransactions []*client.TransactionsData `json:"attemptedChainlinkTransactions"`
}

// KeeperBlockTimeTestReport holds a report information for a single Upkeep Consumer contract
type KeeperBlockTimeTestReport struct {
	ContractAddress        string  `json:"contractAddress"`
	TotalExpectedUpkeeps   int64   `json:"totalExpectedUpkeeps"`
	TotalSuccessfulUpkeeps int64   `json:"totalSuccessfulUpkeeps"`
	AllMissedUpkeeps       []int64 `json:"allMissedUpkeeps"` // List of each time an upkeep was missed, represented by how many blocks it was missed by
}

func (k *KeeperBlockTimeTestReporter) WriteReport(folderPath string) error {
	keeperReportFile, err := os.Create(filepath.Join(folderPath, "block_time_report.csv"))
	if err != nil {
		return err
	}
	defer keeperReportFile.Close()

	keeperReportWriter := csv.NewWriter(keeperReportFile)
	err = keeperReportWriter.Write([]string{
		"Contract Index",
		"Contract Address",
		"Total Expected Upkeeps",
		"Total Successful Upkeeps",
		"Total Missed Upkeeps",
		"Average Blocks Missed",
		"Largest Missed Upkeep",
		"Percent Successful",
	})
	if err != nil {
		return err
	}
	for contractIndex, report := range k.Reports {
		avg, max := int64AvgMax(report.AllMissedUpkeeps)
		err = keeperReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.ContractAddress,
			fmt.Sprint(report.TotalExpectedUpkeeps),
			fmt.Sprint(report.TotalSuccessfulUpkeeps),
			fmt.Sprint(len(report.AllMissedUpkeeps)),
			fmt.Sprint(avg),
			fmt.Sprint(max),
			fmt.Sprintf("%.2f%%", (float64(report.TotalSuccessfulUpkeeps)/float64(report.TotalExpectedUpkeeps))*100),
		})
		if err != nil {
			return err
		}
	}
	keeperReportWriter.Flush()

	txs, err := json.Marshal(k.AttemptedChainlinkTransactions)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(folderPath, "attempted_transactions_report.json"), txs, 0600)
	if err != nil {
		return err
	}

	log.Info().Str("Report Location", folderPath).Msg("Successfully wrote report on Keeper Block Timing")
	return nil
}

// int64AvgMax helper calculates the avg and the max values in a list
func int64AvgMax(in []int64) (float64, int64) {
	var sum int64
	var max int64
	if len(in) == 0 {
		return 0, 0
	}
	for _, num := range in {
		sum += num
		max = int64(math.Max(float64(max), float64(num)))
	}
	return float64(sum) / float64(len(in)), max
}
