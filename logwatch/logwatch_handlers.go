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

type LogTarget string

const (
	Loki LogTarget = "loki"
	File LogTarget = "file"
)

type HandleLogTarget interface {
	Handle(*ContainerLogConsumer, LogContent) error
	GetLogLocation(map[string]*ContainerLogConsumer) (string, error)
	GetTarget() LogTarget
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
		folder := fmt.Sprintf("./logs/%s-%s", testname, time.Now().Format("2006-01-02T15-04-05"))
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

// streams logs to Loki
type LokiLogHandler struct {
	grafanaUrl        string
	shouldSkipLogging bool
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
		"type":      "log_watch",
		"test":      model.LabelValue(content.TestName),
		"container": model.LabelValue(content.ContainerName),
	}, time.Now(), string(content.Content))

	return nil
}

func (h *LokiLogHandler) GetLogLocation(consumers map[string]*ContainerLogConsumer) (string, error) {
	if h.grafanaUrl != "" {
		return h.grafanaUrl, nil
	}

	queries := make([]GrafanaExploreQuery, 0)

	rangeFrom := time.Now()
	rangeTo := time.Now().Add(time.Minute) //just to make sure we get the last message

	for _, c := range consumers {
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

	if len(queries) == 0 {
		return "", errors.New("no Loki consumers found")
	}

	h.grafanaUrl = GrafanaExploreUrl{
		baseurl:    os.Getenv("GRAFANA_URL"),
		datasource: os.Getenv("GRAFANA_DATASOURCE"),
		queries:    queries,
		rangeFrom:  rangeFrom.UnixMilli(),
		rangeTo:    rangeTo.UnixMilli(),
	}.getUrl()

	return h.grafanaUrl, nil
}

func (h LokiLogHandler) GetTarget() LogTarget {
	return Loki
}
