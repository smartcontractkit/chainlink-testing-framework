package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrometheusQueryClient_Query(t *testing.T) {
	tests := []struct {
		name           string
		response       string
		expectedStatus string
		expectedCount  int
		validateResult func(t *testing.T, result *QueryResponse)
	}{
		{
			name: "successful query with multiple metrics",
			response: `{
				"status": "success",
				"data": {
					"resultType": "vector",
					"result": [
						{
							"metric": {
								"__name__": "go_gc_duration_seconds",
								"instance": "cadvisor:8080",
								"job": "cadvisor",
								"quantile": "0.5"
							},
							"value": [1753701299.664, "0.000187874"]
						},
						{
							"metric": {
								"__name__": "go_gc_duration_seconds",
								"instance": "postgres_exporter:9187",
								"job": "postgres",
								"quantile": "0.5"
							},
							"value": [1753701299.664, "0.000257292"]
						}
					]
				}
			}`,
			expectedStatus: "success",
			expectedCount:  2,
			validateResult: func(t *testing.T, result *QueryResponse) {
				assert.Equal(t, "vector", result.Data.ResultType)
				assert.Equal(t, "go_gc_duration_seconds", result.Data.Result[0].Metric["__name__"])
				assert.Equal(t, "cadvisor:8080", result.Data.Result[0].Metric["instance"])
				assert.Equal(t, "cadvisor", result.Data.Result[0].Metric["job"])
				assert.Equal(t, "0.5", result.Data.Result[0].Metric["quantile"])
				assert.Equal(t, 1753701299.664, result.Data.Result[0].Value[0].(float64))
				assert.Equal(t, "0.000187874", result.Data.Result[0].Value[1].(string))
				assert.Equal(t, "go_gc_duration_seconds", result.Data.Result[1].Metric["__name__"])
				assert.Equal(t, "postgres_exporter:9187", result.Data.Result[1].Metric["instance"])
				assert.Equal(t, "postgres", result.Data.Result[1].Metric["job"])
				assert.Equal(t, "0.5", result.Data.Result[1].Metric["quantile"])
				assert.Equal(t, 1753701299.664, result.Data.Result[1].Value[0].(float64))
				assert.Equal(t, "0.000257292", result.Data.Result[1].Value[1].(string))
			},
		},
		{
			name: "query with single metric",
			response: `{
				"status": "success",
				"data": {
					"resultType": "vector",
					"result": [
						{
							"metric": {
								"__name__": "go_goroutines",
								"instance": "localhost:9090",
								"job": "prometheus"
							},
							"value": [1753701299.664, "42"]
						}
					]
				}
			}`,
			expectedStatus: "success",
			expectedCount:  1,
			validateResult: func(t *testing.T, result *QueryResponse) {
				assert.Equal(t, "go_goroutines", result.Data.Result[0].Metric["__name__"])
				assert.Equal(t, "localhost:9090", result.Data.Result[0].Metric["instance"])
				assert.Equal(t, "prometheus", result.Data.Result[0].Metric["job"])
				assert.Equal(t, 1753701299.664, result.Data.Result[0].Value[0].(float64))
				assert.Equal(t, "42", result.Data.Result[0].Value[1].(string))
			},
		},
		{
			name: "empty result",
			response: `{
				"status": "success",
				"data": {
					"resultType": "vector",
					"result": []
				}
			}`,
			expectedStatus: "success",
			expectedCount:  0,
			validateResult: func(t *testing.T, result *QueryResponse) {
				assert.Empty(t, result.Data.Result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/query", r.URL.Path)
				assert.Equal(t, "go_gc_duration_seconds", r.URL.Query().Get("query"))

				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(tt.response))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := NewQueryClient(mockServer.URL)
			timestamp := time.Unix(1753701299, 664000000)
			result, err := client.Query("go_gc_duration_seconds", timestamp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, result.Status)
			assert.Equal(t, tt.expectedCount, len(result.Data.Result))
			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestPrometheusQueryClient_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		response      string
		expectedError string
	}{
		{
			name:          "bad request",
			statusCode:    http.StatusBadRequest,
			response:      `{"status":"error","errorType":"bad_data","error":"invalid query expression"}`,
			expectedError: "unexpected status code: 400",
		},
		{
			name:          "internal server error",
			statusCode:    http.StatusInternalServerError,
			response:      `{"status":"error","errorType":"server_error","error":"internal server error"}`,
			expectedError: "unexpected status code: 500",
		},
		{
			name:          "invalid json",
			statusCode:    http.StatusOK,
			response:      "invalid json",
			expectedError: "failed to unmarshal response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.response))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := NewQueryClient(mockServer.URL)
			result, err := client.Query("go_gc_duration_seconds", time.Now())

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestPrometheusQueryClient_NetworkError(t *testing.T) {
	client := NewQueryClient("http://invalid-url:9090")
	result, err := client.Query("go_gc_duration_seconds", time.Now())

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")
}

func TestPrometheusQueryClientQueryRange(t *testing.T) {
	tests := []struct {
		name           string
		response       string
		expectedStatus string
		expectedCount  int
		validateResult func(t *testing.T, result *QueryRangeResponse)
	}{
		{
			name: "successful range query with multiple data points",
			response: `{
				"status": "success",
				"data": {
					"resultType": "matrix",
					"result": [
						{
							"metric": {
								"__name__": "http_requests_total",
								"job": "api-server",
								"code": "200"
							},
							"values": [
								[1435781451.781, "943"],
								[1435781465.781, "955"],
								[1435781470.781, "971"]
							]
						}
					]
				}
			}`,
			expectedStatus: "success",
			expectedCount:  1,
			validateResult: func(t *testing.T, result *QueryRangeResponse) {
				assert.Equal(t, "matrix", result.Data.ResultType)
				assert.Equal(t, "http_requests_total", result.Data.Result[0].Metric["__name__"])
				assert.Equal(t, "api-server", result.Data.Result[0].Metric["job"])
				assert.Equal(t, "200", result.Data.Result[0].Metric["code"])
				values := result.Data.Result[0].Values
				assert.Len(t, values, 3)
				assert.Equal(t, 1435781451.781, values[0][0].(float64))
				assert.Equal(t, "943", values[0][1].(string))
				assert.Equal(t, 1435781470.781, values[2][0].(float64))
				assert.Equal(t, "971", values[2][1].(string))
			},
		},
		{
			name: "range query with multiple time series",
			response: `{
				"status": "success",
				"data": {
					"resultType": "matrix",
					"result": [
						{
							"metric": {"__name__": "cpu_usage", "instance": "server1"},
							"values": [
								[1435781451.781, "0.45"],
								[1435781465.781, "0.48"]
							]
						},
						{
							"metric": {"__name__": "cpu_usage", "instance": "server2"},
							"values": [
								[1435781451.781, "0.62"],
								[1435781465.781, "0.65"]
							]
						}
					]
				}
			}`,
			expectedStatus: "success",
			expectedCount:  2,
			validateResult: func(t *testing.T, result *QueryRangeResponse) {
				assert.Equal(t, "matrix", result.Data.ResultType)
				assert.Equal(t, "cpu_usage", result.Data.Result[0].Metric["__name__"])
				assert.Equal(t, "server1", result.Data.Result[0].Metric["instance"])
				assert.Equal(t, "0.45", result.Data.Result[0].Values[0][1].(string))
				assert.Equal(t, "cpu_usage", result.Data.Result[1].Metric["__name__"])
				assert.Equal(t, "server2", result.Data.Result[1].Metric["instance"])
				assert.Equal(t, "0.62", result.Data.Result[1].Values[0][1].(string))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/v1/query_range", r.URL.Path)
				assert.Equal(t, "http_requests_total", r.URL.Query().Get("query"))
				assert.Equal(t, "1435781430.000", r.URL.Query().Get("start"))
				assert.Equal(t, "1435781490.000", r.URL.Query().Get("end"))
				assert.Equal(t, "15s", r.URL.Query().Get("step"))

				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(tt.response))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := NewQueryClient(mockServer.URL)
			params := QueryRangeParams{
				Query: "http_requests_total",
				Start: time.Unix(1435781430, 0),
				End:   time.Unix(1435781490, 0),
				Step:  15 * time.Second,
			}
			result, err := client.QueryRange(params)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, result.Status)
			assert.Equal(t, tt.expectedCount, len(result.Data.Result))
			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestPrometheusQueryClientQueryRange_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		response      string
		expectedError string
	}{
		{
			name:          "bad request",
			statusCode:    http.StatusBadRequest,
			response:      `{"status":"error","errorType":"bad_data","error":"invalid query expression"}`,
			expectedError: "unexpected status code: 400",
		},
		{
			name:          "invalid step",
			statusCode:    http.StatusBadRequest,
			response:      `{"status":"error","errorType":"bad_data","error":"zero or negative query resolution step widths are not accepted"}`,
			expectedError: "unexpected status code: 400",
		},
		{
			name:          "timeout",
			statusCode:    http.StatusServiceUnavailable,
			response:      `{"status":"error","errorType":"timeout","error":"query timed out"}`,
			expectedError: "unexpected status code: 503",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.response))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()

			client := NewQueryClient(mockServer.URL)
			params := QueryRangeParams{
				Query: "http_requests_total",
				Start: time.Now().Add(-30 * time.Minute),
				End:   time.Now(),
				Step:  15 * time.Second,
			}
			result, err := client.QueryRange(params)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestPrometheusQueryClient_ResultToLabelsMap(t *testing.T) {
	tests := []struct {
		name     string
		input    *QueryResponse
		expected map[string][]interface{}
	}{
		{
			name: "single metric with multiple labels",
			input: &QueryResponse{
				Data: struct {
					ResultType string `json:"resultType"`
					Result     []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					} `json:"result"`
				}{
					Result: []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					}{
						{
							Metric: map[string]string{
								"__name__": "http_requests_total",
								"job":      "api-server",
								"code":     "200",
							},
							Value: []interface{}{float64(1435781451.781), "943"},
						},
					},
				},
			},
			expected: map[string][]interface{}{
				"__name__:http_requests_total": {"943"},
				"job:api-server":               {"943"},
				"code:200":                     {"943"},
			},
		},
		{
			name: "multiple metrics with shared labels",
			input: &QueryResponse{
				Data: struct {
					ResultType string `json:"resultType"`
					Result     []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					} `json:"result"`
				}{
					Result: []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					}{
						{
							Metric: map[string]string{
								"__name__": "http_requests_total",
								"job":      "api-server",
								"code":     "200",
							},
							Value: []interface{}{float64(1435781451.781), "943"},
						},
						{
							Metric: map[string]string{
								"__name__": "http_requests_total",
								"job":      "api-server",
								"code":     "500",
							},
							Value: []interface{}{float64(1435781451.781), "42"},
						},
					},
				},
			},
			expected: map[string][]interface{}{
				"__name__:http_requests_total": {"943", "42"},
				"job:api-server":               {"943", "42"},
				"code:200":                     {"943"},
				"code:500":                     {"42"},
			},
		},
		{
			name: "empty result",
			input: &QueryResponse{
				Data: struct {
					ResultType string `json:"resultType"`
					Result     []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					} `json:"result"`
				}{
					Result: []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					}{},
				},
			},
			expected: map[string][]interface{}{},
		},
		{
			name: "metric with no labels",
			input: &QueryResponse{
				Data: struct {
					ResultType string `json:"resultType"`
					Result     []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					} `json:"result"`
				}{
					Result: []struct {
						Metric map[string]string `json:"metric"`
						Value  []interface{}     `json:"value"`
					}{
						{
							Metric: map[string]string{},
							Value:  []interface{}{float64(1435781451.781), "1"},
						},
					},
				},
			},
			expected: map[string][]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLabelsMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
