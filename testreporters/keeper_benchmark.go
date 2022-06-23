package testreporters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/onsi/ginkgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
)

// KeeperBenchmarkTestReporter enables reporting on the keeper block time test
type KeeperBenchmarkTestReporter struct {
	Reports                        []KeeperBenchmarkTestReport `json:"reports"`
	ReportMutex                    sync.Mutex
	AttemptedChainlinkTransactions []*client.TransactionsData `json:"attemptedChainlinkTransactions"`

	namespace                 string
	keeperReportFile          string
	attemptedTransactionsFile string
}

// KeeperBenchmarkTestReport holds a report information for a single Upkeep Consumer contract
type KeeperBenchmarkTestReport struct {
	ContractAddress        string  `json:"contractAddress"`
	TotalExpectedUpkeeps   int64   `json:"totalExpectedUpkeeps"`
	TotalSuccessfulUpkeeps int64   `json:"totalSuccessfulUpkeeps"`
	AllMissedUpkeeps       []int64 `json:"allMissedUpkeeps"` // List of each time an upkeep was missed, represented by how many blocks it was missed by
}

func (k *KeeperBenchmarkTestReporter) SetNamespace(namespace string) {
	k.namespace = namespace
}

func (k *KeeperBenchmarkTestReporter) WriteReport(folderLocation string) error {
	k.keeperReportFile = filepath.Join(folderLocation, "./block_time_report.csv")
	k.attemptedTransactionsFile = filepath.Join(folderLocation, "./attempted_transactions_report.json")
	keeperReportFile, err := os.Create(k.keeperReportFile)
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
	var totalExpected, totalSuccessful, totalMissed, worstMiss int64
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
		totalExpected += report.TotalExpectedUpkeeps
		totalSuccessful += report.TotalSuccessfulUpkeeps
		totalMissed += int64(len(report.AllMissedUpkeeps))
		worstMiss = int64(math.Max(float64(max), float64(worstMiss)))
		if err != nil {
			return err
		}
	}
	keeperReportWriter.Flush()

	err = keeperReportWriter.Write([]string{"Full Test Summary"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{"Total Expected", "Total Successful", "Total Missed", "Worst Miss", "Total Percent"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{
		fmt.Sprint(totalExpected),
		fmt.Sprint(totalSuccessful),
		fmt.Sprint(totalMissed),
		fmt.Sprint(worstMiss),
		fmt.Sprintf("%.2f%%", (float64(totalSuccessful)/float64(totalExpected))*100)})
	if err != nil {
		return err
	}
	keeperReportWriter.Flush()

	txs, err := json.Marshal(k.AttemptedChainlinkTransactions)
	if err != nil {
		return err
	}
	err = os.WriteFile(k.attemptedTransactionsFile, txs, 0600)
	if err != nil {
		return err
	}

	log.Info().Msg("Successfully wrote report on Keeper Block Timing")
	return nil
}

// SendSlackNotification sends a slack notification on the results of the test
func (k *KeeperBenchmarkTestReporter) SendSlackNotification(slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(slackAPIKey)
	}

	testFailed := ginkgo.CurrentSpecReport().Failed()
	headerText := ":white_check_mark: Keeper Benchmark Test PASSED :white_check_mark:"
	if testFailed {
		headerText = ":x: Keeper Benchmark Test FAILED :x:"
	}
	messageBlocks := commonSlackNotificationBlocks(slackClient, headerText, k.namespace, k.keeperReportFile, slackUserID, testFailed)
	ts, err := sendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}

	if err := uploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("Keeper Benchmark Test Report %s", k.namespace),
		Filetype:        "csv",
		Filename:        fmt.Sprintf("keeper_benchmark_%s.csv", k.namespace),
		File:            k.keeperReportFile,
		InitialComment:  fmt.Sprintf("Keeper Block Time Test Report %s", k.namespace),
		Channels:        []string{slackChannel},
		ThreadTimestamp: ts,
	}); err != nil {
		return err
	}
	return uploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("Keeper Benchmark Attempted Chainlink Txs %s", k.namespace),
		Filetype:        "json",
		Filename:        fmt.Sprintf("attempted_cl_txs_%s.json", k.namespace),
		File:            k.attemptedTransactionsFile,
		InitialComment:  fmt.Sprintf("Keeper Benchmark Attempted Txs %s", k.namespace),
		Channels:        []string{slackChannel},
		ThreadTimestamp: ts,
	})
}
