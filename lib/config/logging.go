package config

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/net"
)

type LoggingConfig struct {
	TestLogCollect         *bool            `toml:"test_log_collect,omitempty"`
	ShowHTMLCoverageReport *bool            `toml:"show_html_coverage_report,omitempty"` // Show reports with go coverage data
	RunId                  *string          `toml:"run_id,omitempty"`
	Loki                   *LokiConfig      `toml:"-"`
	Grafana                *GrafanaConfig   `toml:"Grafana,omitempty"`
	LogStream              *LogStreamConfig `toml:"LogStream,omitempty"`
}

// Validate executes config validation for LogStream, Grafana and Loki
func (l *LoggingConfig) Validate() error {
	if l.LogStream != nil {
		if err := l.LogStream.Validate(); err != nil {
			return fmt.Errorf("invalid log stream config: %w", err)
		}
	}

	if l.Grafana != nil {
		if err := l.Grafana.Validate(); err != nil {
			return fmt.Errorf("invalid grafana config: %w", err)
		}
	}

	if l.Loki != nil {
		if err := l.Loki.Validate(); err != nil {
			return fmt.Errorf("invalid loki config: %w", err)
		}
	}

	return nil
}

type LogStreamConfig struct {
	LogTargets            []string                `toml:"log_targets"`
	LogProducerTimeout    *blockchain.StrDuration `toml:"log_producer_timeout"`
	LogProducerRetryLimit *uint                   `toml:"log_producer_retry_limit"`
}

// Validate checks that the log stream config is valid, which means that
// log targets are valid and log producer timeout is greater than 0
func (l *LogStreamConfig) Validate() error {
	if len(l.LogTargets) > 0 {
		for _, target := range l.LogTargets {
			if target != "loki" && target != "file" && target != "in-memory" {
				return fmt.Errorf("invalid log target %s", target)
			}
		}
	}

	if l.LogProducerTimeout != nil {
		if l.LogProducerTimeout.Duration == 0 {
			return errors.New("log producer timeout must be greater than 0")
		}
	}

	return nil
}

type LokiConfig struct {
	TenantId    *string `toml:"-"`
	Endpoint    *string `toml:"-"`
	BasicAuth   *string `toml:"-"`
	BearerToken *string `toml:"-"`
}

// Validate checks that the loki config is valid, which means that
// endpoint is a valid URL and tenant id is not empty
func (l *LokiConfig) Validate() error {
	if l.Endpoint != nil {
		if !net.IsValidURL(*l.Endpoint) {
			return fmt.Errorf("invalid loki endpoint %s", *l.Endpoint)
		}
	}
	if l.TenantId == nil || *l.TenantId == "" {
		return errors.New("loki tenant id must be set")
	}

	return nil
}

type GrafanaConfig struct {
	BaseUrl      *string `toml:"base_url"`
	BaseUrlGap   *string `toml:"base_url_gap"` // Base URL for the dashboard via GAP proxy used on CI
	DashboardUrl *string `toml:"dashboard_url"`
	DashboardUID *string `toml:"dashboard_uid"` // UID of the dashboard to put annotations on
	BearerToken  *string `toml:"-"`
}

// Validate checks that the grafana config is valid, which means that
// base url is a valid URL and dashboard url and bearer token are not empty
// but that only applies if they are set
func (c *GrafanaConfig) Validate() error {
	if c.BaseUrl != nil {
		if !net.IsValidURL(*c.BaseUrl) {
			return fmt.Errorf("invalid grafana url %s", *c.BaseUrl)
		}
	}
	if c.DashboardUrl != nil {
		if *c.DashboardUrl == "" {
			return errors.New("if set, grafana dashboard url cannot be an empty string")
		}
	}
	if c.BearerToken != nil {
		if *c.BearerToken == "" {
			return errors.New("if set, grafana Bearer token cannot be an empty string")
		}
	}

	return nil
}
