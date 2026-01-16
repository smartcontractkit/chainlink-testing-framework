package framework

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// PrometheusQueryClient is a client for querying Prometheus metrics
type PrometheusQueryClient struct {
	client  *resty.Client
	baseURL string
}

// NewPrometheusQueryClient creates a new PrometheusQueryClient
func NewPrometheusQueryClient(baseURL string) *PrometheusQueryClient {
	isDebug := os.Getenv("RESTY_DEBUG") == "true"
	return &PrometheusQueryClient{
		client:  resty.New().SetDebug(isDebug),
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// PrometheusQueryResponse represents the response from Prometheus API
type PrometheusQueryResponse struct {
	Status string                 `json:"status"`
	Data   *PromQueryResponseData `json:"data"`
}

type PromQueryResponseData struct {
	ResultType string                    `json:"resultType"`
	Result     []PromQueryResponseResult `json:"result"`
}

type PromQueryResponseResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

// QueryRangeParams contains parameters for range queries
type QueryRangeParams struct {
	Query string
	Start time.Time
	End   time.Time
	Step  time.Duration
}

// Query executes an instant query against the Prometheus API
func (p *PrometheusQueryClient) Query(query string, timestamp time.Time) (*PrometheusQueryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/query", p.baseURL)
	resp, err := p.client.R().
		SetQueryParams(map[string]string{
			"query": query,
			"time":  formatPrometheusTime(timestamp),
		}).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	var result PrometheusQueryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if result.Status != "success" {
		return nil, fmt.Errorf("query failed with status: %s", result.Status)
	}
	return &result, nil
}

// QueryRange executes a range query against the Prometheus API
func (p *PrometheusQueryClient) QueryRange(params QueryRangeParams) (*PrometheusQueryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/query_range", p.baseURL)
	resp, err := p.client.R().
		SetQueryParams(map[string]string{
			"query": params.Query,
			"start": formatPrometheusTime(params.Start),
			"end":   formatPrometheusTime(params.End),
			"step":  formatDuration(params.Step),
		}).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to execute range query: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	var result PrometheusQueryResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if result.Status != "success" {
		return nil, fmt.Errorf("range query failed with status: %s", result.Status)
	}
	return &result, nil
}

// ToLabelsMap converts PrometheusQueryResponse.Data.Result into a map where keys are
// metric labels in "k:v" format and values are slices of all values with that label
func ToLabelsMap(response *PrometheusQueryResponse) map[string][]interface{} {
	resultMap := make(map[string][]interface{})
	for _, res := range response.Data.Result {
		for k, v := range res.Metric {
			label := fmt.Sprintf("%s:%s", k, v)
			resultMap[label] = append(resultMap[label], res.Value[1]) // Value[1] is the metric value
		}
	}
	return resultMap
}

// formatPrometheusTime formats time for Prometheus API
func formatPrometheusTime(t time.Time) string {
	return fmt.Sprintf("%.3f", float64(t.UnixNano())/1e9)
}

// formatDuration formats duration for Prometheus API
func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.0fs", d.Seconds())
}
