package logwatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/wasp"
)

type LogTarget int

const (
	Loki LogTarget = 1 << iota
	File
)

type HandleLogTarget interface {
	Handle(*ContainerLogConsumer, LogContent) error
	PrintLogLocation(*LogWatch)
}

func getDefaultLogHandlers() map[LogTarget]HandleLogTarget {
	handlers := make(map[LogTarget]HandleLogTarget)
	handlers[Loki] = LokiLogHandler{
		shouldSkipLogging: make(map[string]bool),
	}
	handlers[File] = FileLogHandler{
		testLogFolders:    make(map[string]string),
		shouldSkipLogging: make(map[string]bool),
	}

	return handlers
}

// streams logs to local files
type FileLogHandler struct {
	testLogFolders    map[string]string
	shouldSkipLogging map[string]bool
}

func (h FileLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if val, ok := h.shouldSkipLogging[content.TestName]; val && ok {
		return nil
	}

	folder, err := h.getOrCreateLogFolder(content.TestName)
	if err != nil {
		h.shouldSkipLogging[content.TestName] = true

		return errors.Wrap(err, "failed to create logs folder. File logging stopped")
	}

	logFileName := filepath.Join(folder, fmt.Sprintf("%s.log", content.ContainerName))
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		h.shouldSkipLogging[content.TestName] = true

		return errors.Wrap(err, "failed to open log file. File logging stopped")
	}

	defer logFile.Close()

	if _, err := logFile.WriteString(string(content.Content)); err != nil {
		h.shouldSkipLogging[content.TestName] = true

		return errors.Wrap(err, "failed to write to log file. File logging stopped")
	}

	return nil
}

func (h FileLogHandler) PrintLogLocation(l *LogWatch) {
	for testname, folder := range h.testLogFolders {
		l.log.Info().Str("Test", testname).Str("Folder", folder).Msg("Logs saved to folder:")
	}
}

func (h FileLogHandler) getOrCreateLogFolder(testname string) (string, error) {
	var folder string
	if _, ok := h.testLogFolders[testname]; !ok {
		folder = fmt.Sprintf("./logs/%s-%s", testname, time.Now().Format("2006-01-02T15-04-05"))
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return "", err
		}
		h.testLogFolders[testname] = folder
	}
	folder = h.testLogFolders[testname]

	return folder, nil
}

// streams logs to Loki
type LokiLogHandler struct {
	shouldSkipLogging map[string]bool
}

func (h LokiLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if val, ok := h.shouldSkipLogging[content.TestName]; val && ok {
		c.lw.log.Warn().Str("Test", content.TestName).Msg("Skipping pushing logs to Loki for this test")
		return nil
	}

	if c.lw.loki == nil {
		loki, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
		if err != nil {
			c.lw.log.Error().Err(err).Msg("Failed to create Loki client")
			h.shouldSkipLogging[content.TestName] = true

			return err
		}
		c.lw.loki = loki
	}
	// we can notify more than one time if it matches, but we push only once
	_ = c.lw.loki.Handle(model.LabelSet{
		"type":      "log_watch",
		"test":      model.LabelValue(content.TestName),
		"container": model.LabelValue(content.ContainerName),
	}, time.Now(), string(content.Content))

	return nil
}

func (h LokiLogHandler) PrintLogLocation(l *LogWatch) {
	queries := make([]GrafanaExploreQuery, 0)

	rangeFrom := time.Now()
	rangeTo := time.Now()

	for _, c := range l.consumers {
		if c.hasLogTarget(Loki) {
			queries = append(queries, GrafanaExploreQuery{
				refId:     c.name,
				container: c.name,
			})
		}

		// lets find the oldest log message to know when to start the range from
		if len(c.Messages) > 0 {
			var firstMsg struct {
				Ts string `json:"ts"`
			}
			if err := json.Unmarshal([]byte(c.Messages[0]), &firstMsg); err != nil {
				l.log.Error().Err(err).Str("container", c.name).Msg("Failed to unmarshal first log message")
			} else {
				firstTs, err := time.Parse(time.RFC3339, firstMsg.Ts)
				if err != nil {
					l.log.Error().Err(err).Str("container", c.name).Msg("Failed to parse first log message timestamp")
				} else {
					if firstTs.Before(rangeFrom) {
						rangeFrom = firstTs
					}
				}
			}
		}
	}

	grafanaUrl := GrafanaExploreUrl{
		baseurl:    os.Getenv("GRAFANA_URL"),
		datasource: os.Getenv("GRAFANA_DATASOURCE"),
		queries:    queries,
		rangeFrom:  rangeFrom.UnixMilli(),
		rangeTo:    rangeTo.UnixMilli() + 60000, //just to make sure we get the last message
	}.getUrl()

	l.log.Info().Str("URL", string(grafanaUrl)).Msg("Loki logs can be found in Grafana at (will only work when you unescape quotes):")

	fmt.Printf("Loki logs can be found in Grafana at: %s\n", grafanaUrl)
}
