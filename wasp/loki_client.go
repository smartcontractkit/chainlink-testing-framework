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

// NewLokiLogWrapper initializes a new LokiLogWrapper with a specified maximum number of errors allowed. 
// It returns a pointer to the LokiLogWrapper instance, which maintains an internal list of errors encountered during logging operations. 
// The MaxErrors parameter determines how many errors can be recorded before further errors are ignored. 
// If MaxErrors is set to -1, it indicates that all errors should be ignored.
func NewLokiLogWrapper(maxErrors int) *LokiLogWrapper {
	return &LokiLogWrapper{
		MaxErrors: maxErrors,
		errors:    make([]error, 0),
	}
}

// SetClient assigns the provided LokiClient to the LokiLogWrapper instance. 
// This allows the LokiLogWrapper to utilize the specified LokiClient for logging operations. 
// It does not return any value or error.
func (m *LokiLogWrapper) SetClient(c *LokiClient) {
	m.client = c
}

// Log processes and logs a variable number of key-value pairs representing a log message. 
// It checks for the number of provided arguments and logs an error if the message is malformed. 
// If an error is present in the arguments, it appends it to the internal error list and logs the error details. 
// Regardless of the input, it always logs the received log message at a trace level. 
// The function returns nil, indicating successful processing of the log message, 
// or does not modify the internal state if the maximum number of errors has been reached.
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

// Handle processes and sends log data to the Loki server. 
// It takes a set of labels, a timestamp, and a string message as input. 
// If the number of recorded errors exceeds a predefined limit, 
// it returns an error indicating that data cannot be sent. 
// Upon successful processing, it logs the details and sends the 
// log entry to the Loki client channel. 
// The function returns nil if the operation is successful, 
// or an error if any issues arise during the process.
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

// HandleStruct marshals the provided struct into a JSON string and sends it to the Loki logging service along with the associated labels and timestamp. 
// If the marshaling fails, it returns an error indicating the failure. 
// This function is typically used to log structured data responses or statistics to Loki for monitoring and analysis.
func (m *LokiClient) HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error {
	d, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal struct in response: %v", st)
	}
	return m.Handle(ls, t, string(d))
}

// StopNow stops the Loki client immediately. 
// It is typically called when there is a need to halt the logging process, 
// ensuring that no further log entries are sent to the Loki server. 
// This function does not return any value and is intended for use 
// in scenarios where an immediate shutdown of the logging service is required.
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

// DefaultLokiConfig returns a pointer to a new LokiConfig instance initialized with default values. 
// The configuration includes settings for maximum errors, batch wait time, batch size, timeout, 
// and various other parameters relevant to Loki's logging functionality. 
// This function is useful for creating a baseline configuration that can be modified as needed.
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

// NewEnvLokiConfig creates a new LokiConfig instance populated with values 
// from environment variables. It retrieves the TenantID, URL, Token, and 
// BasicAuth fields from the respective environment variables: 
// "LOKI_TENANT_ID", "LOKI_URL", "LOKI_TOKEN", and "LOKI_BASIC_AUTH". 
// The function returns a pointer to the newly created LokiConfig instance.
func NewEnvLokiConfig() *LokiConfig {
	d := DefaultLokiConfig()
	d.TenantID = os.Getenv("LOKI_TENANT_ID")
	d.URL = os.Getenv("LOKI_URL")
	d.Token = os.Getenv("LOKI_TOKEN")
	d.BasicAuth = os.Getenv("LOKI_BASIC_AUTH")
	return d
}

// NewLokiConfig creates a new LokiConfig instance with the specified parameters. 
// It initializes the configuration with default values and overrides them 
// with the provided endpoint, tenant, basicAuth, and token if they are not nil. 
// The function returns a pointer to the newly created LokiConfig.
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

// NewLokiClient initializes a new LokiClient using the provided LokiConfig. 
// It validates the configuration, checks the server URL, and sets up the necessary 
// parameters for the Loki client, including authentication if specified. 
// If successful, it returns a pointer to the newly created LokiClient. 
// In case of any errors during initialization, it returns nil and the corresponding error.
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
