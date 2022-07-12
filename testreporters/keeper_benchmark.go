package testreporters

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
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
	NumRevertedUpkeeps             int64

	namespace                 string
	keeperReportFile          string
	attemptedTransactionsFile string
}

// KeeperBenchmarkTestReport holds a report information for a single Upkeep Consumer contract
type KeeperBenchmarkTestReport struct {
	ContractAddress       string  `json:"contractAddress"`
	TotalEligibleCount    int64   `json:"totalEligibleCount"`
	TotalSLAMissedUpkeeps int64   `json:"totalSLAMissedUpkeeps"`
	TotalPerformedUpkeeps int64   `json:"totalPerformedUpkeeps"`
	AllCheckDelays        []int64 `json:"allCheckDelays"` // List of the delays since checkUpkeep for all performs
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
	var totalEligibleCount, totalPerformed, totalMissedSLA, totalReverted int64
	var allDelays []int64
	for _, report := range k.Reports {
		totalEligibleCount += report.TotalEligibleCount
		totalPerformed += report.TotalPerformedUpkeeps
		totalMissedSLA += report.TotalSLAMissedUpkeeps

		allDelays = append(allDelays, report.AllCheckDelays...)
	}
	totalReverted = k.NumRevertedUpkeeps
	pctWithinSLA := (1.0 - float64(totalMissedSLA)/float64(totalEligibleCount)) * 100
	var pctReverted float64
	if totalPerformed > 0 {
		pctReverted = (float64(totalReverted) / float64(totalPerformed)) * 100
	}

	err = keeperReportWriter.Write([]string{"Full Test Summary"})
	if err != nil {
		return err
	}
	err = keeperReportWriter.Write([]string{
		"Total Times Eligible",
		"Total Performed",
		"Total Reverted",
		"Average Perform Delay",
		"Median Perform Delay",
		"90th pct Perform Delay",
		"99th pct Perform Delay",
		"Max Perform Delay",
		"Percent Within SLA",
		"Percent Revert",
	})
	if err != nil {
		return err
	}
	avg, median, ninetyPct, ninetyNinePct, max := intListStats(allDelays)
	err = keeperReportWriter.Write([]string{
		fmt.Sprint(totalEligibleCount),
		fmt.Sprint(totalPerformed),
		fmt.Sprint(totalReverted),
		fmt.Sprintf("%.2f", avg),
		fmt.Sprint(median),
		fmt.Sprint(ninetyPct),
		fmt.Sprint(ninetyNinePct),
		fmt.Sprint(max),
		fmt.Sprintf("%.2f%%", pctWithinSLA),
		fmt.Sprintf("%.2f%%", pctReverted),
	})
	if err != nil {
		return err
	}
	keeperReportWriter.Flush()
	log.Info().
		Int64("Total Times Eligible", totalEligibleCount).
		Int64("Total Performed", totalPerformed).
		Int64("Total Reverted", totalReverted).
		Float64("Average Perform Delay", avg).
		Int64("Median Perform Delay", median).
		Int64("90th pct Perform Delay", ninetyPct).
		Int64("99th pct Perform Delay", ninetyNinePct).
		Int64("Max Perform Delay", max).
		Float64("Percent Within SLA", pctWithinSLA).
		Float64("Percent Reverted", pctReverted).
		Msg("Calculated Aggregate Results")

	err = keeperReportWriter.Write([]string{
		"Contract Index",
		"Contract Address",
		"Total Times Eligible",
		"Total Performed Upkeeps",
		"Average Perform Delay",
		"Median Perform Delay",
		"90th pct Perform Delay",
		"99th pct Perform Delay",
		"Largest Perform Delay",
		"Percent Within SLA",
	})
	if err != nil {
		return err
	}

	for contractIndex, report := range k.Reports {
		avg, median, ninetyPct, ninetyNinePct, max := intListStats(report.AllCheckDelays)
		err = keeperReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.ContractAddress,
			fmt.Sprint(report.TotalEligibleCount),
			fmt.Sprint(report.TotalPerformedUpkeeps),
			fmt.Sprintf("%.2f", avg),
			fmt.Sprint(median),
			fmt.Sprint(ninetyPct),
			fmt.Sprint(ninetyNinePct),
			fmt.Sprint(max),
			fmt.Sprintf("%.2f%%", (1.0-float64(report.TotalSLAMissedUpkeeps)/float64(report.TotalEligibleCount))*100),
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

// intListStats helper calculates some statistics on an int list: avg, median, 90pct, 99pct, max
func intListStats(in []int64) (float64, int64, int64, int64, int64) {
	length := len(in)
	if length == 0 {
		return 0, 0, 0, 0, 0
	}
	sort.Slice(in, func(i, j int) bool { return in[i] < in[j] })
	var sum int64
	for _, num := range in {
		sum += num
	}
	return float64(sum) / float64(length), in[int(math.Floor(float64(length)*0.5))], in[int(math.Floor(float64(length)*0.9))], in[int(math.Floor(float64(length)*0.99))], in[length-1]
}
