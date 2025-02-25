package wasp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grafana/dskit/backoff"
	dskit "github.com/grafana/dskit/flagext"
	lokiAPI "github.com/grafana/loki/v3/clients/pkg/promtail/api"
	lokiClient "github.com/grafana/loki/v3/clients/pkg/promtail/client"
	lokiProto "github.com/grafana/loki/v3/pkg/logproto"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
)

// LokiLogWrapper wraps Loki errors received through logs, handles them
type LokiLogWrapper struct {
	MaxErrors int
	errors    []error
	client    *LokiClient
}

// NewLokiLogWrapper initializes a LokiLogWrapper with the specified maximum number of errors.
// It is used to track and limit error occurrences within the LokiClient.
func NewLokiLogWrapper(maxErrors int) *LokiLogWrapper {
	return &LokiLogWrapper{
		MaxErrors: maxErrors,
		errors:    make([]error, 0),
	}
}

// SetClient assigns a LokiClient to the LokiLogWrapper.
// Use it to link the log wrapper with a specific Loki client instance for managing log operations.
func (m *LokiLogWrapper) SetClient(c *LokiClient) {
	m.client = c
}

// Log processes and forwards log entries to Loki, handling malformed messages and recording errors.
// It ensures error limits are respected and logs at appropriate levels for monitoring purposes.
func (m *LokiLogWrapper) Log(kvars ...interface{}) error {
	if len(m.errors) > m.MaxErrors {
		return nil
	}
	if len(kvars) < 13 {
		log.Error().
			Interface("Line", kvars).
			Msg("Malformed promtail log message, skipping")
		return nil
	}
	if kvars[13] != nil {
		if _, ok := kvars[13].(error); ok {
			m.errors = append(m.errors, kvars[13].(error))
			log.Error().
				Interface("Status", kvars[9]).
				Str("Error", kvars[13].(error).Error()).
				Msg("Loki error")
		}
	}
	log.Trace().Interface("Line", kvars).Msg("Loki client internal log")
	return nil
}

// LokiClient is a Loki/Promtail client wrapper
type LokiClient struct {
	logWrapper *LokiLogWrapper
	lokiClient.Client
}

// Handle sends a log entry with the given labels, timestamp, and message to the Loki service.
// It checks error thresholds and returns an error if sending fails or limits are exceeded.
func (m *LokiClient) Handle(ls model.LabelSet, t time.Time, s string) error {
	if m.logWrapper.MaxErrors != -1 && len(m.logWrapper.errors) > m.logWrapper.MaxErrors {
		return fmt.Errorf("can't send data to Loki, errors: %v", m.logWrapper.errors)
	}
	log.Trace().
		Interface("Labels", ls).
		Time("Time", t).
		Str("Data", s).
		Msg("Sending data to Loki")
	m.Client.Chan() <- lokiAPI.Entry{Labels: ls, Entry: lokiProto.Entry{Timestamp: t, Line: s}}
	return nil
}

// HandleStruct marshals the provided struct to JSON and sends it to Loki with the specified labels and timestamp.
// Use this function to log structured data in a decentralized logging system.
func (m *LokiClient) HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error {
	d, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal struct in response: %v", st)
	}
	return m.Handle(ls, t, string(d))
}

// StopNow gracefully terminates the Loki client, ensuring all active streams
// are closed and resources are released. Use this function to stop Loki when
// shutting down services or when Loki is no longer needed.
func (m *LokiClient) StopNow() {
	m.Client.StopNow()
}

// LokiConfig is simplified subset of a Promtail client configuration
type LokiConfig struct {
	// URL url to Loki endpoint
	URL string `yaml:"url"`
	// Token is Loki authorization token
	Token string `yaml:"token"`
	// BasicAuth is a basic login:password auth string
	BasicAuth string `yaml:"basic_auth"`
	// MaxErrors max amount of errors to ignore before exiting
	MaxErrors int
	// BatchWait max time to wait until sending a new batch
	BatchWait time.Duration
	// BatchSize size of a messages batch
	BatchSize int
	// Timeout is batch send timeout
	Timeout time.Duration
	// BackoffConfig backoff configuration
	BackoffConfig backoff.Config
	// Headers are additional request headers
	Headers map[string]string
	// The tenant ID to use when pushing logs to Loki (empty string means
	// single tenant mode)
	TenantID string
	// When enabled, Promtail will not retry batches that get a
	// 429 'Too Many Requests' response from the distributor. Helps
	// prevent HOL blocking in multitenant deployments.
	DropRateLimitedBatches bool
	// ExposePrometheusMetrics if enabled exposes Promtail Prometheus metrics
	ExposePrometheusMetrics bool
	MaxStreams              int
	MaxLineSize             int
	MaxLineSizeTruncate     bool
}

// DefaultLokiConfig returns a LokiConfig initialized with default parameters.
// It serves as a base configuration for Loki clients, allowing users to customize settings as needed.
func DefaultLokiConfig() *LokiConfig {
	return &LokiConfig{
		MaxErrors:               5,
		BatchWait:               3 * time.Second,
		BatchSize:               500 * 1024,
		Timeout:                 20 * time.Second,
		DropRateLimitedBatches:  false,
		ExposePrometheusMetrics: false,
		MaxStreams:              600,
		MaxLineSize:             999999,
		MaxLineSizeTruncate:     false,
	}
}

// NewEnvLokiConfig creates a LokiConfig populated with settings from environment variables.
func NewEnvLokiConfig() *LokiConfig {
	d := DefaultLokiConfig()
	d.TenantID = os.Getenv("LOKI_TENANT_ID")
	d.URL = os.Getenv("LOKI_URL")
	d.Token = os.Getenv("LOKI_TOKEN")
	d.BasicAuth = os.Getenv("LOKI_BASIC_AUTH")
	return d
}

// NewLokiConfig initializes a LokiConfig with the provided endpoint, tenant, basicAuth, and token.
// Use it to customize connection settings for the Loki logging service.
func NewLokiConfig(endpoint *string, tenant *string, basicAuth *string, token *string) *LokiConfig {
	d := DefaultLokiConfig()
	if endpoint != nil {
		d.URL = *endpoint
	}
	if tenant != nil {
		d.TenantID = *tenant
	}
	if basicAuth != nil {
		d.BasicAuth = *basicAuth
	}
	if token != nil {
		d.Token = *token
	}
	return d
}

// NewLokiClient initializes a new LokiClient with the given LokiConfig.
// It validates the configuration, sets up authentication, and prepares the client for interacting with Loki for logging purposes.
func NewLokiClient(extCfg *LokiConfig) (*LokiClient, error) {
	_, err := http.Get(extCfg.URL)
	if err != nil {
		return nil, err
	}
	serverURL := dskit.URLValue{}
	err = serverURL.Set(extCfg.URL)
	if err != nil {
		return nil, err
	}
	if extCfg.MaxErrors < -1 {
		return nil, errors.New("max errors should be 0..N, -1 to ignore errors")
	}
	cfg := lokiClient.Config{
		URL:                    serverURL,
		BatchWait:              extCfg.BatchWait,
		BatchSize:              extCfg.BatchSize,
		Timeout:                extCfg.Timeout,
		DropRateLimitedBatches: extCfg.DropRateLimitedBatches,
		BackoffConfig:          extCfg.BackoffConfig,
		Headers:                extCfg.Headers,
		TenantID:               extCfg.TenantID,
		Client: config.HTTPClientConfig{
			TLSConfig: config.TLSConfig{InsecureSkipVerify: true},
		},
	}
	if extCfg.BasicAuth != "" {
		logpass := strings.Split(extCfg.BasicAuth, ":")
		if len(logpass) != 2 {
			return nil, errors.New("basic auth should be in login:password format")
		}
		cfg.Client.BasicAuth = &config.BasicAuth{
			Username: logpass[0],
			Password: config.Secret(logpass[1]),
		}
	}
	if extCfg.Token != "" {
		cfg.Client.BearerToken = config.Secret(extCfg.Token)
	}
	ll := NewLokiLogWrapper(extCfg.MaxErrors)
	c, err := lokiClient.New(lokiClient.NewMetrics(nil), cfg, extCfg.MaxStreams, extCfg.MaxLineSize, extCfg.MaxLineSizeTruncate, ll)
	if err != nil {
		return nil, err
	}
	lc := &LokiClient{
		logWrapper: ll,
		Client:     c,
	}
	ll.SetClient(lc)
	return lc, nil
}
