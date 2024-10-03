package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

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
	BasicAuth   string
	QueryParams LokiQueryParams
}

// LokiQueryParams holds the parameters required for querying Loki
type LokiQueryParams struct {
	Query     string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
}

// NewLokiClient creates a new Loki client with the given parameters
func NewLokiClient(baseURL, tenantID, basicAuth string, queryParams LokiQueryParams) *LokiClient {
	return &LokiClient{
		BaseURL:     baseURL,
		TenantID:    tenantID,
		BasicAuth:   basicAuth,
		QueryParams: queryParams,
	}
}

// QueryLogs queries Loki logs based on the query parameters and returns the raw log entries
func (lc *LokiClient) QueryLogs(ctx context.Context) ([]LokiLogEntry, error) {
	client := resty.New()

	// Build query parameters
	params := map[string]string{
		"query": lc.QueryParams.Query,
		"start": fmt.Sprintf("%d", lc.QueryParams.StartTime.UnixNano()),
		"end":   fmt.Sprintf("%d", lc.QueryParams.EndTime.UnixNano()),
		"limit": fmt.Sprintf("%d", lc.QueryParams.Limit),
	}

	// Send request using Resty
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("X-Scope-OrgID", lc.TenantID).
		SetBasicAuth(lc.BasicAuth, "").
		SetQueryParams(params).
		Get(lc.BaseURL + "/loki/api/v1/query_range")

	if err != nil {
		return nil, fmt.Errorf("error querying Loki: %w", err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		log.Printf("Loki API returned status code: %d", resp.StatusCode())
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode())
	}

	// Parse the response into the LokiResponse struct
	var lokiResp LokiResponse
	if err := json.Unmarshal(resp.Body(), &lokiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Extract log entries from the response
	logEntries := lc.extractRawLogEntries(lokiResp)
	return logEntries, nil
}

// extractRawLogEntries processes the LokiResponse and returns raw log entries
func (lc *LokiClient) extractRawLogEntries(lokiResp LokiResponse) []LokiLogEntry {
	var logEntries []LokiLogEntry

	for _, result := range lokiResp.Data.Result {
		for _, entry := range result.Values {
			timestamp := entry[0].(string)
			logLine := entry[1].(string)
			logEntries = append(logEntries, LokiLogEntry{
				Timestamp: timestamp,
				Log:       logLine,
			})
		}
	}

	return logEntries
}
