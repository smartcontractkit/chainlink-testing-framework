package config

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/net"
)

type LoggingConfig struct {
	ShowHTMLCoverageReport *bool          `toml:"show_html_coverage_report,omitempty"` // Show reports with go coverage data
	Loki                   *LokiConfig    `toml:"-"`
	Grafana                *GrafanaConfig `toml:"Grafana,omitempty"`
}

// Validate executes config validation for Grafana and Loki
func (l *LoggingConfig) Validate() error {
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
	BaseUrlCI    *string `toml:"base_url_github_ci"` // URL of GAP proxy used on CI for Grafana
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
