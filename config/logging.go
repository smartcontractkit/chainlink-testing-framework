package config

import (
	_ "embed"
	"net/url"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

//go:embed tomls/default.toml
var DefaultLoggingConfig []byte

type LoggingConfig struct {
	TestLogCollect *bool            `toml:"test_log_collect"`
	TestLogLevel   *string          `toml:"test_log_level"`
	LokiTenantId   *string          `toml:"loki_tenant_id"`
	LokiUrl        *string          `toml:"loki_url"`
	LokiBasicAuth  *string          `toml:"loki_basic_auth"`
	Grafana        *GrafanaConfig   `toml:"Grafana"`
	RunId          *string          `toml:"run_id"`
	LogStream      *LogStreamConfig `toml:"log_stream"`
}

type LogStreamConfig struct {
	LogTargets            []string                    `toml:"log_targets"`
	LogProducerTimeout    *blockchain.JSONStrDuration `toml:"log_producer_timeout"`
	LogProducerRetryLimit *uint                       `toml:"log_producer_retry_limit"`
}

func (l *LoggingConfig) Validate() error {
	// LokiUrl is a valid URL, but only if log target includes loki
	// GrafanaUrl is a valid URL, but only if log target includes loki
	// GrafanaDataSource is not "", but only if log target includes loki

	if l.TestLogLevel != nil {
		validLevels := []string{"trace", "debug", "info", "warn", "error", "panic", "fatal"}
		valid := false
		for _, level := range validLevels {
			if *l.TestLogLevel == strings.ToUpper(level) {
				valid = true
				break
			}
		}
		if !valid {
			return errors.Errorf("invalid test log level %s", *l.TestLogLevel)
		}
	}

	if l.LokiUrl != nil {
		if !isValidURL(*l.LokiUrl) {
			return errors.Errorf("invalid loki url %s", *l.LokiUrl)
		}
	}

	if l.LogStream != nil {
		if err := l.LogStream.Validate(); err != nil {
			return errors.Wrapf(err, "invalid log stream config")
		}
	}

	if l.Grafana != nil {
		if err := l.Grafana.Validate(); err != nil {
			return errors.Wrapf(err, "invalid grafana config")
		}
	}

	return nil
}

func (l *LoggingConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case LoggingConfig:
		if asCfg.TestLogLevel != nil {
			l.TestLogLevel = asCfg.TestLogLevel
		}
		if asCfg.TestLogCollect != nil {
			l.TestLogCollect = asCfg.TestLogCollect
		}
		if asCfg.LogStream != nil {
			l.LogStream = asCfg.LogStream
		}
		if asCfg.LokiTenantId != nil {
			l.LokiTenantId = asCfg.LokiTenantId
		}
		if asCfg.LokiUrl != nil {
			l.LokiUrl = asCfg.LokiUrl
		}
		if asCfg.LokiBasicAuth != nil {
			l.LokiBasicAuth = asCfg.LokiBasicAuth
		}
		if asCfg.Grafana != nil {
			l.Grafana = asCfg.Grafana
		}
		if asCfg.RunId != nil {
			l.RunId = asCfg.RunId
		}

		return nil
	case *LoggingConfig:
		if asCfg == nil {
			return nil
		}
		if asCfg.TestLogLevel != nil {
			l.TestLogLevel = asCfg.TestLogLevel
		}
		if asCfg.TestLogCollect != nil {
			l.TestLogCollect = asCfg.TestLogCollect
		}
		if asCfg.LogStream != nil {
			l.LogStream = asCfg.LogStream
		}
		if asCfg.LokiTenantId != nil {
			l.LokiTenantId = asCfg.LokiTenantId
		}
		if asCfg.LokiUrl != nil {
			l.LokiUrl = asCfg.LokiUrl
		}
		if asCfg.LokiBasicAuth != nil {
			l.LokiBasicAuth = asCfg.LokiBasicAuth
		}
		if asCfg.Grafana.GrafanaUrl != nil {
			l.Grafana.GrafanaUrl = asCfg.Grafana.GrafanaUrl
		}
		if asCfg.RunId != nil {
			l.RunId = asCfg.RunId
		}

		return nil
	default:
		return errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (l *LoggingConfig) Default() error {
	if err := toml.Unmarshal(DefaultLoggingConfig, l); err != nil {
		return errors.Wrapf(err, "error unmarshaling config")
	}

	return nil
}

func (l *LogStreamConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case LogStreamConfig:
		if asCfg.LogTargets != nil {
			l.LogTargets = asCfg.LogTargets
		}
		if asCfg.LogProducerTimeout != nil {
			l.LogProducerTimeout = asCfg.LogProducerTimeout
		}
		if asCfg.LogProducerRetryLimit != nil {
			l.LogProducerRetryLimit = asCfg.LogProducerRetryLimit
		}

		return nil
	case *LogStreamConfig:
		if asCfg == nil {
			return nil
		}
		if asCfg.LogTargets != nil {
			l.LogTargets = asCfg.LogTargets
		}
		if asCfg.LogProducerTimeout != nil {
			l.LogProducerTimeout = asCfg.LogProducerTimeout
		}
		if asCfg.LogProducerRetryLimit != nil {
			l.LogProducerRetryLimit = asCfg.LogProducerRetryLimit
		}

		return nil
	default:
		return errors.Errorf("cannot apply overrides to log stream config from unknown type %T", from)
	}
}

func (l *LogStreamConfig) Validate() error {
	if len(l.LogTargets) > 0 {
		for _, target := range l.LogTargets {
			if target != "loki" && target != "file" && target != "in-memory" {
				return errors.Errorf("invalid log target %s", target)
			}
		}
	}

	if l.LogProducerTimeout != nil {
		if l.LogProducerTimeout.Duration == 0 {
			return errors.Errorf("log producer timeout must be greater than 0")
		}
	}

	return nil
}

func (l *LogStreamConfig) Default() error {
	return nil
}

type GrafanaConfig struct {
	GrafanaUrl *string `toml:"grafana_url"`
}

func (c *GrafanaConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case GrafanaConfig:
		if asCfg.GrafanaUrl != nil {
			c.GrafanaUrl = asCfg.GrafanaUrl
		}

		return nil
	case *GrafanaConfig:
		if asCfg == nil {
			return nil
		}
		if asCfg.GrafanaUrl != nil {
			c.GrafanaUrl = asCfg.GrafanaUrl
		}

		return nil
	default:
		return errors.Errorf("cannot apply overrides to grafana config from unknown type %T", from)
	}
}

func (c *GrafanaConfig) Validate() error {
	if c.GrafanaUrl != nil {
		if !isValidURL(*c.GrafanaUrl) {
			return errors.Errorf("invalid grafana url %s", *c.GrafanaUrl)
		}
	}

	return nil
}

func (c *GrafanaConfig) Default() error {
	return nil
}

func isValidURL(testURL string) bool {
	parsedURL, err := url.Parse(testURL)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}
