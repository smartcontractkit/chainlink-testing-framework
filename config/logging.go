package config

import (
	_ "embed"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

//go:embed tomls/default.toml
var DefaultLoggingConfig []byte

type LoggingConfig struct {
	Logging *struct {
		TestLogCollect *bool            `toml:"test_log_collect"`
		TestLogLevel   *string          `toml:"test_log_level"`
		LokiTenantId   *string          `toml:"loki_tenant_id"`
		LokiUrl        *string          `toml:"loki_url"`
		LokiBasicAuth  *string          `toml:"loki_basic_auth"`
		Grafana        *GrafanaConfig   `toml:"Grafana"`
		RunId          *string          `toml:"run_id"`
		LogStream      *LogStreamConfig `toml:"log_stream"`
	} `toml:"Logging"`
}

type LogStreamConfig struct {
	LogTargets            []string                    `toml:"log_targets"`
	LogProducerTimeout    *blockchain.JSONStrDuration `toml:"log_producer_timeout"`
	LogProducerRetryLimit *uint                       `toml:"log_producer_retry_limit"`
}

type GrafanaConfig struct {
	GrafanaUrl *string `toml:"grafana_url"`
}

func (l *LoggingConfig) ReadSecrets() error {
	lokiBasicAuth, err := osutil.GetEnv("LOKI_BASIC_AUTH")
	if err != nil {
		return err
	}

	if lokiBasicAuth != "" {
		l.Logging.LokiBasicAuth = &lokiBasicAuth
	}

	return nil
}

func (l *LoggingConfig) Validate() error {
	// TestLogLevel in ["trace", "debug", "info", "warn", "error", "panic", "fatal"]
	// LogStreamLogTargets in ["loki", "file", "in-memory"] -- add method to LS to get valid targets
	// LokiUrl is a valid URL, but only if log target includes loki
	// GrafanaUrl is a valid URL, but only if log target includes loki
	// GrafanaDataSource is not "", but only if log target includes loki

	return nil
}

func (l *LoggingConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case LoggingConfig:
		if asCfg.Logging.TestLogLevel != nil {
			l.Logging.TestLogLevel = asCfg.Logging.TestLogLevel
		}
		if asCfg.Logging.TestLogCollect != nil {
			l.Logging.TestLogCollect = asCfg.Logging.TestLogCollect
		}
		if asCfg.Logging.LogStream != nil {
			l.Logging.LogStream = asCfg.Logging.LogStream
		}
		if asCfg.Logging.LokiTenantId != nil {
			l.Logging.LokiTenantId = asCfg.Logging.LokiTenantId
		}
		if asCfg.Logging.LokiUrl != nil {
			l.Logging.LokiUrl = asCfg.Logging.LokiUrl
		}
		if asCfg.Logging.LokiBasicAuth != nil {
			l.Logging.LokiBasicAuth = asCfg.Logging.LokiBasicAuth
		}
		if asCfg.Logging.Grafana.GrafanaUrl != nil {
			l.Logging.Grafana.GrafanaUrl = asCfg.Logging.Grafana.GrafanaUrl
		}
		if asCfg.Logging.RunId != nil {
			l.Logging.RunId = asCfg.Logging.RunId
		}

		return nil
	case *LoggingConfig:
		if asCfg.Logging.TestLogLevel != nil {
			l.Logging.TestLogLevel = asCfg.Logging.TestLogLevel
		}
		if asCfg.Logging.TestLogCollect != nil {
			l.Logging.TestLogCollect = asCfg.Logging.TestLogCollect
		}
		if asCfg.Logging.LogStream != nil {
			l.Logging.LogStream = asCfg.Logging.LogStream
		}
		if asCfg.Logging.LokiTenantId != nil {
			l.Logging.LokiTenantId = asCfg.Logging.LokiTenantId
		}
		if asCfg.Logging.LokiUrl != nil {
			l.Logging.LokiUrl = asCfg.Logging.LokiUrl
		}
		if asCfg.Logging.LokiBasicAuth != nil {
			l.Logging.LokiBasicAuth = asCfg.Logging.LokiBasicAuth
		}
		if asCfg.Logging.Grafana.GrafanaUrl != nil {
			l.Logging.Grafana.GrafanaUrl = asCfg.Logging.Grafana.GrafanaUrl
		}
		if asCfg.Logging.RunId != nil {
			l.Logging.RunId = asCfg.Logging.RunId
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

// func isValidURL(testURL string) bool {
// 	parsedURL, err := url.Parse(testURL)
// 	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
// }
