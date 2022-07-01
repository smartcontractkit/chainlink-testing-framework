package testreporters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/onsi/ginkgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
)

// KeeperBenchmarkTestReporter enables reporting on the keeper benchmark test
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
	AllCheckDelays         []int64 `json:"allCheckDelays"` // List of the delays since checkUpkeep for all performs
}

func (k *KeeperBenchmarkTestReporter) SetNamespace(namespace string) {
	k.namespace = namespace
}

func (k *KeeperBenchmarkTestReporter) WriteReport(folderLocation string) error {
	k.keeperReportFile = filepath.Join(folderLocation, "./benchmark_report.csv")
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
		"Average Perform Delay",
		"Largest Perform Delay",
		"Percent Successful",
	})
	if err != nil {
		return err
	}
	var totalExpected, totalSuccessful int64
	var allDelays []int64
	for contractIndex, report := range k.Reports {
		avg, max := int64AvgMax(report.AllCheckDelays)

		err = keeperReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.ContractAddress,
			fmt.Sprint(report.TotalExpectedUpkeeps),
			fmt.Sprint(report.TotalSuccessfulUpkeeps),
			fmt.Sprint(avg),
			fmt.Sprint(max),
			fmt.Sprintf("%.2f%%", (float64(report.TotalSuccessfulUpkeeps)/float64(report.TotalExpectedUpkeeps))*100),
		})
		totalExpected += report.TotalExpectedUpkeeps
		totalSuccessful += report.TotalSuccessfulUpkeeps
		if err != nil {
			return err
		}
		allDelays = append(allDelays, report.AllCheckDelays...)
	}
	keeperReportWriter.Flush()

	err = keeperReportWriter.Write([]string{"Full Test Summary"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{"Total Expected", "Total Successful", "Average Perform Delay", "Largest Perform Delay", "Percent Successful"})
	if err != nil {
		return err
	}
	avg, max := int64AvgMax(allDelays)
	err = keeperReportWriter.Write([]string{
		fmt.Sprint(totalExpected),
		fmt.Sprint(totalSuccessful),
		fmt.Sprint(avg),
		fmt.Sprint(max),
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

	log.Info().Msg("Successfully wrote report on Keeper Benchmark")
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
		InitialComment:  fmt.Sprintf("Keeper Benchmark Test Report %s", k.namespace),
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
