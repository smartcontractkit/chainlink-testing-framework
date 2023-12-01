package logwatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/wasp"
)

type LogTarget string

const (
	Loki LogTarget = "loki"
	File LogTarget = "file"
)

type HandleLogTarget interface {
	Handle(*ContainerLogConsumer, LogContent) error
	GetLogLocation(map[string]*ContainerLogConsumer) (string, error)
	GetTarget() LogTarget
	SetRunId(string)
	GetRunId() string
}

func getDefaultLogHandlers() map[LogTarget]HandleLogTarget {
	handlers := make(map[LogTarget]HandleLogTarget)
	handlers[Loki] = &LokiLogHandler{}
	handlers[File] = &FileLogHandler{}

	return handlers
}

// streams logs to local files
type FileLogHandler struct {
	logFolder         string
	shouldSkipLogging bool
	runId             string
}

func (h *FileLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if h.shouldSkipLogging {
		return nil
	}

	folder, err := h.getOrCreateLogFolder(content.TestName)
	if err != nil {
		h.shouldSkipLogging = true

		return errors.Wrap(err, "failed to create logs folder. File logging stopped")
	}

	logFileName := filepath.Join(folder, fmt.Sprintf("%s.log", content.ContainerName))
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		h.shouldSkipLogging = true

		return errors.Wrap(err, "failed to open log file. File logging stopped")
	}

	defer logFile.Close()

	if _, err := logFile.WriteString(string(content.Content)); err != nil {
		h.shouldSkipLogging = true

		return errors.Wrap(err, "failed to write to log file. File logging stopped")
	}

	return nil
}

func (h FileLogHandler) GetLogLocation(_ map[string]*ContainerLogConsumer) (string, error) {
	return h.logFolder, nil
}

func (h *FileLogHandler) getOrCreateLogFolder(testname string) (string, error) {
	if h.logFolder == "" {
		folder := fmt.Sprintf("./logs/%s-%s-%s", testname, time.Now().Format("2006-01-02T15-04-05"), h.runId)
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return "", err
		}
		h.logFolder = folder
	}

	return h.logFolder, nil
}

func (h FileLogHandler) GetTarget() LogTarget {
	return File
}

func (h *FileLogHandler) SetRunId(executionId string) {
	h.runId = executionId
}

func (h *FileLogHandler) GetRunId() string {
	return h.runId
}

// streams logs to Loki
type LokiLogHandler struct {
	grafanaUrl        string
	shouldSkipLogging bool
	runId             string
}

func (h *LokiLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if h.shouldSkipLogging {
		c.lw.log.Warn().Str("Test", content.TestName).Msg("Skipping pushing logs to Loki for this test")
		return nil
	}

	if c.lw.loki == nil {
		loki, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
		if err != nil {
			c.lw.log.Error().Err(err).Msg("Failed to create Loki client")
			h.shouldSkipLogging = true

			return err
		}
		c.lw.loki = loki
	}
	// we can notify more than one time if it matches, but we push only once
	_ = c.lw.loki.Handle(model.LabelSet{
		"type":         "log_watch",
		"test":         model.LabelValue(content.TestName),
		"container_id": model.LabelValue(content.ContainerName),
		"run_id":       model.LabelValue(h.runId),
	}, time.Now(), string(content.Content))

	return nil
}

func (h *LokiLogHandler) GetLogLocation(consumers map[string]*ContainerLogConsumer) (string, error) {
	if h.grafanaUrl != "" {
		return h.grafanaUrl, nil
	}

	grafanaBaseUrl := os.Getenv("GRAFANA_URL")
	if grafanaBaseUrl == "" {
		return "", errors.New("GRAFANA_URL env var is not set")
	}

	grafanaBaseUrl = strings.TrimSuffix(grafanaBaseUrl, "/")

	rangeFrom := time.Now()
	rangeTo := time.Now().Add(time.Minute) //just to make sure we get the last message

	var sb strings.Builder
	sb.WriteString(grafanaBaseUrl)
	sb.WriteString("/d/ddf75041-1e39-42af-aa46-361fe4c36e9e/ci-e2e-tests-logs?orgId=1&")
	sb.WriteString(fmt.Sprintf("var-run_id=%s", h.runId))

	if len(consumers) == 0 {
		return "", errors.New("no Loki consumers found")
	}

	for _, c := range consumers {
		if c.hasLogTarget(Loki) {
			sb.WriteString(fmt.Sprintf("&var-container_id=%s", c.name))
		}

		// lets find the oldest log message to know when to start the range from
		if len(c.Messages) > 0 {
			var firstMsg struct {
				Ts string `json:"ts"`
			}

			if err := json.Unmarshal([]byte(c.Messages[0]), &firstMsg); err != nil {
				return "", errors.Errorf("failed to unmarshal first log message for container '%s'", c.name)
			}

			firstTs, err := time.Parse(time.RFC3339, firstMsg.Ts)
			if err != nil {
				return "", errors.Errorf("failed to parse first log message's timestamp '%+v' for container '%s'", firstTs, c.name)
			}

			if firstTs.Before(rangeFrom) {
				rangeFrom = firstTs
			}
		}
	}

	sb.WriteString(fmt.Sprintf("&from=%d&to=%d", rangeFrom.UnixMilli(), rangeTo.UnixMilli()))
	h.grafanaUrl = sb.String()

	return h.grafanaUrl, nil
}

func (h LokiLogHandler) GetTarget() LogTarget {
	return Loki
}

func (h *LokiLogHandler) SetRunId(executionId string) {
	h.runId = executionId
}

func (h *LokiLogHandler) GetRunId() string {
	return h.runId
}
