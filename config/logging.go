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
	RunId          *string          `toml:"run_id"`
	Loki           *LokiConfig      `toml:"Loki"`
	Grafana        *GrafanaConfig   `toml:"Grafana"`
	LogStream      *LogStreamConfig `toml:"LogStream"`
}

func (l *LoggingConfig) Validate() error {
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
	if from.LogStream != nil && l.LogStream == nil {
		l.LogStream = from.LogStream
	} else if from.LogStream != nil && l.LogStream != nil {
		if err := l.LogStream.ApplyOverrides(from.LogStream); err != nil {
			return errors.Wrapf(err, "error applying overrides to log stream config")
		}
	}
	if from.Loki != nil && l.Loki == nil {
		l.Loki = from.Loki
	} else if from.Loki != nil && l.Loki != nil {
		if err := l.Loki.ApplyOverrides(from.Loki); err != nil {
			return errors.Wrapf(err, "error applying overrides to loki config")
		}
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

type LogStreamConfig struct {
	LogTargets            []string                `toml:"log_targets"`
	LogProducerTimeout    *blockchain.StrDuration `toml:"log_producer_timeout"`
	LogProducerRetryLimit *uint                   `toml:"log_producer_retry_limit"`
}

func (l *LogStreamConfig) ApplyOverrides(from *LogStreamConfig) error {
	if from == nil {
		return nil
	}
	if from.LogTargets != nil {
		l.LogTargets = from.LogTargets
	}
	if from.LogProducerTimeout != nil {
		l.LogProducerTimeout = from.LogProducerTimeout
	}
	if from.LogProducerRetryLimit != nil {
		l.LogProducerRetryLimit = from.LogProducerRetryLimit
	}

	return nil
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

type LokiConfig struct {
	TenantId  *string `toml:"tenant_id"`
	Url       *string `toml:"url"`
	BasicAuth *string `toml:"basic_auth"`
}

func (l *LokiConfig) Validate() error {
	if l.Url != nil {
		if !net.IsValidURL(*l.Url) {
			return errors.Errorf("invalid loki url %s", *l.Url)
		}
	}

	return nil
}

func (l *LokiConfig) ApplyOverrides(from *LokiConfig) error {
	if from == nil {
		return nil
	}
	if from.TenantId != nil {
		l.TenantId = from.TenantId
	}
	if from.Url != nil {
		l.Url = from.Url
	}
	if from.BasicAuth != nil {
		l.BasicAuth = from.BasicAuth
	}

	return nil
}

type GrafanaConfig struct {
	Url *string `toml:"url"`
}

func (c *GrafanaConfig) ApplyOverrides(from *GrafanaConfig) error {
	if from == nil {
		return nil
	}
	if from.Url != nil {
		c.Url = from.Url
	}

	return nil
}

func (c *GrafanaConfig) Validate() error {
	if c.Url != nil {
		if !net.IsValidURL(*c.Url) {
			return errors.Errorf("invalid grafana url %s", *c.Url)
		}
	}

	return nil
}

func (c *GrafanaConfig) Default() error {
	return nil
}
