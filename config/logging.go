package config

import (
	_ "embed"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/net"
)

//go:embed tomls/logging_default.toml
var DefaultLoggingConfig []byte

type LoggingConfig struct {
	TestLogCollect *bool            `toml:"test_log_collect"`
	LokiTenantId   *string          `toml:"loki_tenant_id"`
	LokiUrl        *string          `toml:"loki_url"`
	LokiBasicAuth  *string          `toml:"loki_basic_auth"`
	RunId          *string          `toml:"run_id"`
	Grafana        *GrafanaConfig   `toml:"Grafana"`
	LogStream      *LogStreamConfig `toml:"LogStream"`
}

type LogStreamConfig struct {
	LogTargets            []string                `toml:"log_targets"`
	LogProducerTimeout    *blockchain.StrDuration `toml:"log_producer_timeout"`
	LogProducerRetryLimit *uint                   `toml:"log_producer_retry_limit"`
}

func (l *LoggingConfig) Validate() error {
	if l.LokiUrl != nil {
		if !net.IsValidURL(*l.LokiUrl) {
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

func (l *LoggingConfig) ApplyOverrides(from *LoggingConfig) error {
	if from == nil {
		return nil
	}
	if from.TestLogCollect != nil {
		l.TestLogCollect = from.TestLogCollect
	}
	if from.LogStream != nil {
		l.LogStream = from.LogStream
	}
	if from.LokiTenantId != nil {
		l.LokiTenantId = from.LokiTenantId
	}
	if from.LokiUrl != nil {
		l.LokiUrl = from.LokiUrl
	}
	if from.LokiBasicAuth != nil {
		l.LokiBasicAuth = from.LokiBasicAuth
	}
	if from.Grafana != nil {
		l.Grafana = from.Grafana
	}
	if from.RunId != nil {
		l.RunId = from.RunId
	}

	return nil
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
		if !net.IsValidURL(*c.GrafanaUrl) {
			return errors.Errorf("invalid grafana url %s", *c.GrafanaUrl)
		}
	}

	return nil
}

func (c *GrafanaConfig) Default() error {
	return nil
}
