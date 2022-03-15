package testreporters

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

type OCRSoakTestReporter struct {
	Reports     map[string]*OCRSoakTestReport // contractAddress: Report
	namespace   string
	csvLocation string
}

type OCRSoakTestReport struct {
	ContractAddress string
	TotalRounds     uint

	averageRoundTime  time.Duration
	LongestRoundTime  time.Duration
	ShortestRoundTime time.Duration
	totalRoundTimes   time.Duration

	averageRoundBlocks  uint
	LongestRoundBlocks  uint
	ShortestRoundBlocks uint
	totalBlockLengths   uint
}

// SetNamespace sets the namespace of the report for clean reports
func (o *OCRSoakTestReporter) SetNamespace(namespace string) {
	o.namespace = namespace
}

// WriteReport writes OCR Soak test report to logs
func (o *OCRSoakTestReporter) WriteReport(folderLocation string) error {
	for _, report := range o.Reports {
		report.averageRoundBlocks = report.totalBlockLengths / report.TotalRounds
		report.averageRoundTime = time.Duration(report.totalRoundTimes.Nanoseconds() / int64(report.TotalRounds))
	}
	if err := o.writeCSV(folderLocation); err != nil {
		return err
	}

	log.Info().Msg("OCR Soak Test Report")
	log.Info().Msg("--------------------")
	for contractAddress, report := range o.Reports {
		log.Info().
			Str("Contract Address", report.ContractAddress).
			Uint("Total Rounds Processed", report.TotalRounds).
			Str("Average Round Time", fmt.Sprint(report.averageRoundTime)).
			Str("Longest Round Time", fmt.Sprint(report.LongestRoundTime)).
			Str("Shortest Round Time", fmt.Sprint(report.ShortestRoundTime)).
			Uint("Average Round Blocks", report.averageRoundBlocks).
			Uint("Longest Round Blocks", report.LongestRoundBlocks).
			Uint("Shortest Round Blocks", report.ShortestRoundBlocks).
			Msg(contractAddress)
	}
	log.Info().Msg("--------------------")
	return nil
}

// SendNotification sends a slack message to a slack webhook and uploads test artifacts
func (o *OCRSoakTestReporter) SendSlackNotification(slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(slackAPIKey)
	}

	testFailed := ginkgo.CurrentSpecReport().Failed()
	headerText := ":white_check_mark: OCR Soak Test PASSED :white_check_mark:"
	messageBlocks := []slack.Block{}
	if testFailed {
		headerText = ":x: OCR Soak Test FAILED :x:"
	}
	messageBlocks = append(messageBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	messageBlocks = append(messageBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", o.namespace, false, false)))
	messageBlocks = append(messageBlocks, slack.NewDividerBlock())
	messageBlocks = append(messageBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("Test ran for %s\nSummary CSV created on _remote-test-runner_ at _%s_\nNotifying <@%s>",
			ginkgo.CurrentSpecReport().RunTime.Truncate(time.Second), o.csvLocation, slackUserID), false, true), nil, nil))
	if testFailed {
		messageBlocks = append(messageBlocks,
			slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
		messageBlocks = append(messageBlocks, slack.NewDividerBlock())
		messageBlocks = append(messageBlocks, slack.NewSectionBlock(
			slack.NewTextBlockObject("plain_text", ginkgo.CurrentSpecReport().FailureMessage(), false, false), nil, nil))
	}
	ts, err := sendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}

	return uploadSlackFile(slackClient, slack.FileUploadParameters{
		Title:           fmt.Sprintf("OCR Soak Test Report %s", o.namespace),
		Filetype:        "csv",
		Filename:        fmt.Sprintf("ocr_soak_%s.csv", o.namespace),
		File:            o.csvLocation,
		InitialComment:  fmt.Sprintf("OCR Soak Test Report %s.", o.namespace),
		Channels:        []string{slackChannel},
		ThreadTimestamp: ts,
	})
}

// UpdateReport updates the report based on the latest info
func (o *OCRSoakTestReport) UpdateReport(roundTime time.Duration, blockLength uint) {
	// Updates min values from default 0
	if o.ShortestRoundBlocks == 0 {
		o.ShortestRoundBlocks = blockLength
	}
	if o.ShortestRoundTime == 0 {
		o.ShortestRoundTime = roundTime
	}
	o.TotalRounds++
	o.totalRoundTimes += roundTime
	o.totalBlockLengths += blockLength
	if roundTime >= o.LongestRoundTime {
		o.LongestRoundTime = roundTime
	}
	if roundTime <= o.ShortestRoundTime {
		o.ShortestRoundTime = roundTime
	}
	if blockLength >= o.LongestRoundBlocks {
		o.LongestRoundBlocks = blockLength
	}
	if blockLength <= o.ShortestRoundBlocks {
		o.ShortestRoundBlocks = blockLength
	}
}

// writes a CSV report on the test runner
func (o *OCRSoakTestReporter) writeCSV(folderLocation string) error {
	reportLocation := filepath.Join(folderLocation, "./ocr_soak_report.csv")
	log.Debug().Str("Location", reportLocation).Msg("Writing OCR report")
	o.csvLocation = reportLocation
	ocrReportFile, err := os.Create(reportLocation)
	if err != nil {
		return err
	}
	defer ocrReportFile.Close()

	ocrReportWriter := csv.NewWriter(ocrReportFile)
	err = ocrReportWriter.Write([]string{
		"Contract Index",
		"Contract Address",
		"Total Rounds Processed",
		"Average Round Time",
		"Longest Round Time",
		"Shortest Round Time",
		"Average Round Blocks",
		"Longest Round Blocks",
		"Shortest Round Blocks",
	})
	if err != nil {
		return err
	}
	for contractIndex, report := range o.Reports {
		err = ocrReportWriter.Write([]string{
			fmt.Sprint(contractIndex),
			report.ContractAddress,
			fmt.Sprint(report.TotalRounds),
			fmt.Sprint(report.averageRoundTime),
			fmt.Sprint(report.LongestRoundTime),
			fmt.Sprint(report.ShortestRoundTime),
			fmt.Sprint(report.averageRoundBlocks),
			fmt.Sprint(report.LongestRoundBlocks),
			fmt.Sprint(report.ShortestRoundBlocks),
		})
		if err != nil {
			return err
		}
	}
	ocrReportWriter.Flush()

	log.Info().Str("Location", reportLocation).Msg("Wrote CSV file")
	return nil
}
