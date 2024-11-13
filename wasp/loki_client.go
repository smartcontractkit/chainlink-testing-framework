package wasp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"errors"
	"strings"

	"github.com/grafana/dskit/backoff"
	dskit "github.com/grafana/dskit/flagext"
	lokiAPI "github.com/grafana/loki/clients/pkg/promtail/api"
	lokiClient "github.com/grafana/loki/clients/pkg/promtail/client"
	lokiProto "github.com/grafana/loki/pkg/logproto"
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

// NewLokiLogWrapper creates a new instance of LokiLogWrapper with a specified maximum number of errors.
// The maxErrors parameter determines the limit of errors to be tracked, where -1 indicates no limit.
func NewLokiLogWrapper(maxErrors int) *LokiLogWrapper {
	return &LokiLogWrapper{
		MaxErrors: maxErrors,
		errors:    make([]error, 0),
	}
}

// SetClient assigns a LokiClient to the LokiLogWrapper. 
// This allows the LokiLogWrapper to interact with the specified LokiClient for logging operations.
func (m *LokiLogWrapper) SetClient(c *LokiClient) {
	m.client = c
}

// Log processes a variable number of key-value pairs, kvars, and logs them using the Loki client.
// If the number of errors exceeds MaxErrors, the function returns without logging.
// If kvars has fewer than 13 elements, it logs an error message indicating a malformed log message.
// If the 14th element of kvars is an error, it appends it to the errors slice and logs the error details.
// It also logs a trace message with the provided kvars.
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

// Handle sends a log entry to Loki with the specified label set, timestamp, and string data.
// It returns an error if the maximum number of errors allowed is exceeded.
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

// HandleStruct marshals the provided struct into JSON and sends it to Loki using the specified label set and timestamp.
// It returns an error if the struct cannot be marshaled or if there is an issue sending the data to Loki.
func (m *LokiClient) HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error {
	d, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal struct in response: %v", st)
	}
	return m.Handle(ls, t, string(d))
}

// StopNow immediately stops the LokiClient by invoking the StopNow method on its underlying client. 
// It is typically used to halt the Loki logging stream when the Loki configuration is present and active.
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

// DefaultLokiConfig returns a pointer to a LokiConfig struct with default settings.
// These settings include parameters such as MaxErrors, BatchWait, BatchSize, Timeout,
// and others, which are initialized to predefined values suitable for typical use cases.
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

// NewEnvLokiConfig creates a new LokiConfig instance with default settings.
// It populates the TenantID, URL, Token, and BasicAuth fields using environment variables
// LOKI_TENANT_ID, LOKI_URL, LOKI_TOKEN, and LOKI_BASIC_AUTH, respectively.
func NewEnvLokiConfig() *LokiConfig {
	d := DefaultLokiConfig()
	d.TenantID = os.Getenv("LOKI_TENANT_ID")
	d.URL = os.Getenv("LOKI_URL")
	d.Token = os.Getenv("LOKI_TOKEN")
	d.BasicAuth = os.Getenv("LOKI_BASIC_AUTH")
	return d
}

// NewLokiConfig creates a new LokiConfig with the specified endpoint, tenant, basicAuth, and token.
// It initializes the configuration with default values and overrides them with the provided parameters if they are not nil.
// The function returns a pointer to the configured LokiConfig instance.
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

// NewLokiClient initializes and returns a new LokiClient based on the provided LokiConfig.
// It validates the configuration, sets up authentication if needed, and configures the client
// with specified parameters such as batch size, timeout, and headers. It returns an error if
// the configuration is invalid or if the client cannot be created.
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
