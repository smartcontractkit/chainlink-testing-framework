package config

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/net"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

type LoggingConfig struct {
	TestLogCollect         *bool            `toml:"test_log_collect,omitempty"`
	ShowHTMLCoverageReport *bool            `toml:"show_html_coverage_report,omitempty"` // Show reports with go coverage data
	RunId                  *string          `toml:"run_id,omitempty"`
	Loki                   *LokiConfig      `toml:"Loki,omitempty"`
	Grafana                *GrafanaConfig   `toml:"Grafana,omitempty"`
	LogStream              *LogStreamConfig `toml:"LogStream,omitempty"`
}

func (c *LoggingConfig) LoadFromEnv() error {
	if c.Loki == nil {
		c.Loki = &LokiConfig{}
	}
	err := c.Loki.LoadFromEnv()
	if err != nil {
		return errors.Wrap(err, "error loading Loki config from env")
	}
	if c.Grafana == nil {
		c.Grafana = &GrafanaConfig{}
	}
	err = c.Grafana.LoadFromEnv()
	if err != nil {
		return errors.Wrap(err, "error loading Grafana config from env")
	}
	return nil
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
	TenantId    *string `toml:"tenant_id"`
	Endpoint    *string `toml:"endpoint"`
	BasicAuth   *string `toml:"basic_auth_secret"`
	BearerToken *string `toml:"bearer_token_secret"`
}

func (l *LokiConfig) LoadFromEnv() error {
	logger := logging.GetTestLogger(nil)

	if l.TenantId == nil {
		tenantId, err := readEnvVarValue("E2E_TEST_LOKI_TENANT_ID", String)
		if err != nil {
			return err
		}
		if tenantId != nil && tenantId.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_LOKI_TENANT_ID env var to override Loki.TenantId")
			l.TenantId = ptr.Ptr(tenantId.(string))
		}
	}
	if l.Endpoint == nil {
		endpoint, err := readEnvVarValue("E2E_TEST_LOKI_ENDPOINT", String)
		if err != nil {
			return err
		}
		if endpoint != nil && endpoint.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_LOKI_ENDPOINT env var to override Loki.Endpoint")
			l.Endpoint = ptr.Ptr(endpoint.(string))
		}
	}
	if l.BasicAuth == nil {
		basicAuth, err := readEnvVarValue("E2E_TEST_LOKI_BASIC_AUTH", String)
		if err != nil {
			return err
		}
		if basicAuth != nil && basicAuth.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_LOKI_BASIC_AUTH env var to override Loki.BasicAuth")
			l.BasicAuth = ptr.Ptr(basicAuth.(string))
		}
	}
	if l.BearerToken == nil {
		bearerToken, err := readEnvVarValue("E2E_TEST_LOKI_BEARER_TOKEN", String)
		if err != nil {
			return err
		}
		if bearerToken != nil && bearerToken.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_LOKI_BEARER_TOKEN env var to override Loki.BearerToken")
			l.BearerToken = ptr.Ptr(bearerToken.(string))
		}
	}
	return nil
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
	DashboardUrl *string `toml:"dashboard_url"`
	DashboardUID *string `toml:"dashboard_uid"` // UID of the dashboard to put annotations on
	BearerToken  *string `toml:"bearer_token_secret"`
}

func (l *GrafanaConfig) LoadFromEnv() error {
	logger := logging.GetTestLogger(nil)

	if l.BaseUrl == nil {
		baseUrl, err := readEnvVarValue("E2E_TEST_GRAFANA_BASE_URL", String)
		if err != nil {
			return err
		}
		if baseUrl != nil && baseUrl.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_GRAFANA_BASE_URL env var to override Grafana.BaseUrl")
			l.BaseUrl = ptr.Ptr(baseUrl.(string))
		}
	}
	if l.DashboardUrl == nil {
		dashboardUrl, err := readEnvVarValue("E2E_TEST_GRAFANA_DASHBOARD_URL", String)
		if err != nil {
			return err
		}
		if dashboardUrl != nil && dashboardUrl.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_GRAFANA_DASHBOARD_URL env var to override Grafana.DashboardUrl")
			l.DashboardUrl = ptr.Ptr(dashboardUrl.(string))
		}
	}
	if l.BearerToken == nil {
		bearerToken, err := readEnvVarValue("E2E_TEST_GRAFANA_BEARER_TOKEN", String)
		if err != nil {
			return err
		}
		if bearerToken != nil && bearerToken.(string) != "" {
			logger.Debug().Msg("Using E2E_TEST_GRAFANA_BEARER_TOKEN env var to override Grafana.BearerToken")
			l.BearerToken = ptr.Ptr(bearerToken.(string))
		}
	}
	return nil
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
