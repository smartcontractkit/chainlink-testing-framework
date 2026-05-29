package wasp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// OTELConfig is a simplified subset of an OTEL OTLP log client configuration.
type OTELConfig struct {
	// Endpoint is the OTLP HTTP endpoint host:port (e.g. "localhost:4318").
	Endpoint string `yaml:"endpoint"`
	// URLPath is the OTLP HTTP path for log ingestion (defaults to "/v1/logs").
	URLPath string `yaml:"url_path"`
	// Insecure disables TLS for the OTLP HTTP exporter.
	Insecure bool `yaml:"insecure"`
	// Headers are additional request headers (e.g. for auth tokens).
	Headers map[string]string `yaml:"headers"`
	// ServiceName populates the OTEL service.name resource attribute.
	ServiceName string `yaml:"service_name"`
	// BatchTimeout is the max time the batch processor waits before flushing.
	BatchTimeout time.Duration
	// ExportTimeout is the per-batch export timeout.
	ExportTimeout time.Duration
}

// DefaultOTELConfig returns an OTELConfig initialized with reasonable defaults
// pointing at the local compose-victoria-metrics OTEL Collector (OTLP HTTP on :4318),
// which forwards logs to VictoriaLogs. Override Endpoint for production use.
func DefaultOTELConfig() *OTELConfig {
	return &OTELConfig{
		Endpoint:      "localhost:4318",
		URLPath:       "/v1/logs",
		Insecure:      true,
		ServiceName:   "wasp",
		BatchTimeout:  3 * time.Second,
		ExportTimeout: 20 * time.Second,
	}
}

// NewEnvOTELConfig creates an OTELConfig populated from environment variables.
// Recognized: OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_INSECURE, OTEL_SERVICE_NAME.
// Returns nil unless LOG_SEND_METHOD is unset or "otel" — this lets callers set both
// LokiConfig and OTELConfig on wasp.Config and have only the selected backend become non-nil.
func NewEnvOTELConfig() *OTELConfig {
	if logSendMethod() != LogSendMethodOTEL {
		return nil
	}
	d := DefaultOTELConfig()
	if v := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); v != "" {
		d.Endpoint = v
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "false" {
		d.Insecure = false
	}
	if v := os.Getenv("OTEL_SERVICE_NAME"); v != "" {
		d.ServiceName = v
	}
	return d
}

// OTELClient is a thin wrapper over an OTLP HTTP log exporter + batch provider.
type OTELClient struct {
	provider *sdklog.LoggerProvider
	logger   otellog.Logger
}

// NewOTELClient initializes an OTELClient that pushes logs via OTLP HTTP to the
// configured endpoint. Use StopNow to flush and shut down cleanly.
func NewOTELClient(extCfg *OTELConfig) (*OTELClient, error) {
	if extCfg.Endpoint == "" {
		return nil, fmt.Errorf("OTEL endpoint is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(extCfg.Endpoint),
		otlploghttp.WithURLPath(extCfg.URLPath),
	}
	if extCfg.Insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}
	if len(extCfg.Headers) > 0 {
		opts = append(opts, otlploghttp.WithHeaders(extCfg.Headers))
	}
	if extCfg.ExportTimeout > 0 {
		opts = append(opts, otlploghttp.WithTimeout(extCfg.ExportTimeout))
	}

	exp, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP log exporter: %w", err)
	}

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(extCfg.ServiceName),
	))
	if err != nil {
		return nil, fmt.Errorf("create OTEL resource: %w", err)
	}

	processorOpts := []sdklog.BatchProcessorOption{}
	if extCfg.BatchTimeout > 0 {
		processorOpts = append(processorOpts, sdklog.WithExportInterval(extCfg.BatchTimeout))
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp, processorOpts...)),
	)

	return &OTELClient{
		provider: provider,
		logger:   provider.Logger("wasp"),
	}, nil
}

// Handle emits a log record with the given labels (as attributes), timestamp, and body string.
func (m *OTELClient) Handle(ls model.LabelSet, t time.Time, body string) error {
	var rec otellog.Record
	observedTime := time.Now()
	rec.SetTimestamp(t)
	rec.SetObservedTimestamp(observedTime)
	rec.SetSeverity(otellog.SeverityInfo)
	rec.SetBody(otellog.StringValue(body))
	log.Trace().
		Interface("Labels", ls).
		Time("Time", t).
		Time("ObservedTime", observedTime).
		Str("Data", body).
		Msg("Sending data to OTEL")
	for k, v := range ls {
		rec.AddAttributes(otellog.String(string(k), string(v)))
	}
	m.logger.Emit(context.Background(), rec)
	return nil
}

// HandleStruct marshals the given struct to JSON and emits it as an OTEL log record body.
func (m *OTELClient) HandleStruct(ls model.LabelSet, t time.Time, st any) error {
	d, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %v", st)
	}
	return m.Handle(ls, t, string(d))
}

// StopNow flushes and shuts down the OTEL provider.
func (m *OTELClient) StopNow() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = m.provider.Shutdown(ctx)
}
