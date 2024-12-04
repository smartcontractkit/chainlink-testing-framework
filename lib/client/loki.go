package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

// LokiAPIError is a custom error type for handling non-200 responses from the Loki API
type LokiAPIError struct {
	StatusCode int
	Message    string
}

// Implement the `Error` interface for LokiAPIError
func (e *LokiAPIError) Error() string {
	return fmt.Sprintf("Loki API error: %s (status code: %d)", e.Message, e.StatusCode)
}

// LokiBasicAuth holds the authentication details for Loki
type LokiBasicAuth struct {
	Login    string
	Password string
}

// LokiResponse represents the structure of the response from Loki
type LokiResponse struct {
	Data struct {
		Result []struct {
			Stream map[string]string `json:"stream"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// LokiLogEntry represents a single log entry with a timestamp and raw log message
type LokiLogEntry struct {
	Timestamp string
	Log       string
}

// LokiClient represents a client to interact with Loki for querying logs
type LokiClient struct {
	BaseURL     string
	TenantID    string
	BasicAuth   LokiBasicAuth
	QueryParams LokiQueryParams
	Logger      logging.Logger
	RestyClient *resty.Client
}

// LokiQueryParams holds the parameters required for querying Loki
type LokiQueryParams struct {
	Query     string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
}

// NewLokiClient creates a new Loki client with the given parameters, initializes a logger, and configures Resty with debug mode
func NewLokiClient(baseURL, tenantID string, auth LokiBasicAuth, queryParams LokiQueryParams) *LokiClient {
	logging.Init()

	logger := logging.GetLogger(nil, "LOKI_CLIENT_LOG_LEVEL")
	logger.Info().
		Str("BaseURL", baseURL).
		Str("TenantID", tenantID).
		Msg("Initializing Loki Client")

	// Set debug mode for Resty if RESTY_DEBUG is enabled
	isDebug := os.Getenv("RESTY_DEBUG") == "true"

	restyClient := resty.New().
		SetDebug(isDebug)

	return &LokiClient{
		BaseURL:     baseURL,
		TenantID:    tenantID,
		BasicAuth:   auth,
		QueryParams: queryParams,
		Logger:      logger,
		RestyClient: restyClient,
	}
}

// QueryLogs queries Loki logs based on the query parameters and returns the raw log entries
func (lc *LokiClient) QueryLogs(ctx context.Context) ([]LokiLogEntry, error) {
	// Log request details
	lc.Logger.Info().
		Str("Query", lc.QueryParams.Query).
		Str("StartTime", lc.QueryParams.StartTime.Format(time.RFC3339Nano)).
		Str("EndTime", lc.QueryParams.EndTime.Format(time.RFC3339Nano)).
		Int("Limit", lc.QueryParams.Limit).
		Msg("Making request to Loki API")

	// Start tracking request duration
	start := time.Now()

	// Build query parameters
	params := map[string]string{
		"query": lc.QueryParams.Query,
		"start": lc.QueryParams.StartTime.Format(time.RFC3339Nano),
		"end":   lc.QueryParams.EndTime.Format(time.RFC3339Nano),
		"limit": fmt.Sprintf("%d", lc.QueryParams.Limit),
	}

	// Send request using the pre-configured Resty client
	resp, err := lc.RestyClient.R().
		SetContext(ctx).
		SetHeader("X-Scope-OrgID", lc.TenantID).
		SetBasicAuth(lc.BasicAuth.Login, lc.BasicAuth.Password).
		SetQueryParams(params).
		Get(lc.BaseURL + "/loki/api/v1/query_range")

	// Track request duration
	duration := time.Since(start)

	if err != nil {
		lc.Logger.Error().Err(err).Dur("duration", duration).Msg("Error querying Loki")
		return nil, err
	}

	// Log non-200 responses
	if resp.StatusCode() != 200 {
		bodySnippet := string(resp.Body())
		if len(bodySnippet) > 200 {
			bodySnippet = bodySnippet[:200] + "..."
		}
		lc.Logger.Error().
			Int("StatusCode", resp.StatusCode()).
			Dur("duration", duration).
			Str("ResponseBody", bodySnippet).
			Msg("Loki API returned non-200 status")
		return nil, &LokiAPIError{
			StatusCode: resp.StatusCode(),
			Message:    "unexpected status code from Loki API",
		}
	}

	// Log successful response
	lc.Logger.Info().
		Int("StatusCode", resp.StatusCode()).
		Dur("duration", duration).
		Msg("Successfully queried Loki API")

	// Parse the response into the LokiResponse struct
	var lokiResp LokiResponse
	if err := json.Unmarshal(resp.Body(), &lokiResp); err != nil {
		lc.Logger.Error().Err(err).Msg("Error decoding response from Loki")
		return nil, err
	}

	// Extract log entries from the response
	logEntries := lc.extractRawLogEntries(lokiResp)

	// Log the number of entries retrieved
	lc.Logger.Info().Int("LogEntries", len(logEntries)).Msg("Successfully retrieved logs from Loki")

	return logEntries, nil
}

// extractRawLogEntries processes the LokiResponse and returns raw log entries
func (lc *LokiClient) extractRawLogEntries(lokiResp LokiResponse) []LokiLogEntry {
	var logEntries []LokiLogEntry

	for _, result := range lokiResp.Data.Result {
		for _, entry := range result.Values {
			if len(entry) != 2 {
				lc.Logger.Error().Interface("Log entry", entry).Msgf("Error parsing log entry. Expected 2 elements, got %d", len(entry))
				continue
			}
			var timestamp string
			if entry[0] == nil {
				lc.Logger.Error().Msg("Error parsing timestamp. Entry at index 0, that should be a timestamp, is nil")
				continue
			}
			if timestampString, ok := entry[0].(string); ok {
				timestamp = timestampString
			} else if timestampInt, ok := entry[0].(int); ok {
				timestamp = fmt.Sprintf("%d", timestampInt)
			} else if timestampFloat, ok := entry[0].(float64); ok {
				timestamp = fmt.Sprintf("%f", timestampFloat)
			} else {
				lc.Logger.Error().Msgf("Error parsing timestamp. Expected string, int, or float64, got %T", entry[0])
				continue
			}
			logLine := entry[1].(string)
			logEntries = append(logEntries, LokiLogEntry{
				Timestamp: timestamp,
				Log:       logLine,
			})
		}
	}

	return logEntries
}
