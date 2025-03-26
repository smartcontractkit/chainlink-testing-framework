package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func TestOTELSendMetricsToOTELCollector(t *testing.T) {
	t.Skip("run manually to debug otel-collector")
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-metrics-service"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "test"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	metrics := generateMetricsData(res)
	jsonData, _ := json.MarshalIndent(metrics, "", "  ")
	t.Logf("Sending metrics payload:\n%s\n", string(jsonData))

	err = sendMetricsToCollector(metrics)
	if err != nil {
		t.Fatalf("failed to send metrics: %v", err)
	}
	t.Log("Metrics sent successfully to OTEL Collector!")
}

func TestOTELSendLogsToOTELCollector(t *testing.T) {
	t.Skip("run manually to debug otel-collector")
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-logging-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	logs := generateLogRecords(res)
	jsonData, _ := json.MarshalIndent(logs, "", "  ")
	t.Logf("Sending logs payload:\n%s\n", string(jsonData))
	err = sendLogsToCollector(logs)
	if err != nil {
		t.Fatalf("failed to send logs: %v", err)
	}
	t.Log("Logs sent successfully!")
}

func TestOTELSendTracesToOTELCollector(t *testing.T) {
	t.Skip("run manually to debug otel-collector")
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-tracing-service"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "test"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	traceData := generateTraceData(res)
	jsonData, _ := json.MarshalIndent(traceData, "", "  ")
	t.Logf("Sending traces payload:\n%s\n", string(jsonData))
	err = sendTracesToCollector(traceData)
	if err != nil {
		t.Fatalf("failed to send traces: %v", err)
	}
	t.Log("Traces sent successfully to OTEL Collector!")
}

func generateMetricsData(res *resource.Resource) map[string]interface{} {
	now := time.Now()
	nanos := now.UnixNano()
	startTime := now.Add(-5 * time.Minute).UnixNano()

	return map[string]interface{}{
		"resourceMetrics": []interface{}{
			map[string]interface{}{
				"resource": map[string]interface{}{
					"attributes": resourceToAttributes(res),
				},
				"scopeMetrics": []interface{}{
					map[string]interface{}{
						"scope": map[string]interface{}{
							"name":    "test-metrics",
							"version": "1.0",
						},
						"metrics": []interface{}{
							// Counter example
							map[string]interface{}{
								"name":        "http.requests",
								"description": "Count of HTTP requests",
								"unit":        "1",
								"sum": map[string]interface{}{
									"dataPoints": []interface{}{
										map[string]interface{}{
											"startTimeUnixNano": startTime,
											"timeUnixNano":      nanos,
											"asInt":             "42",
											"attributes": []interface{}{
												map[string]interface{}{
													"key":   "http.method",
													"value": map[string]interface{}{"stringValue": "GET"},
												},
												map[string]interface{}{
													"key":   "http.status",
													"value": map[string]interface{}{"stringValue": "200"},
												},
											},
										},
									},
									"aggregationTemporality": 2, // Cumulative
									"isMonotonic":            true,
								},
							},
							// Gauge example
							map[string]interface{}{
								"name":        "memory.usage",
								"description": "Memory usage in bytes",
								"unit":        "By",
								"gauge": map[string]interface{}{
									"dataPoints": []interface{}{
										map[string]interface{}{
											"timeUnixNano": nanos,
											"asDouble":     1024.42,
											"attributes": []interface{}{
												map[string]interface{}{
													"key":   "memory.type",
													"value": map[string]interface{}{"stringValue": "heap"},
												},
											},
										},
									},
								},
							},
							// Histogram example
							map[string]interface{}{
								"name":        "http.request.duration",
								"description": "HTTP request duration in seconds",
								"unit":        "s",
								"histogram": map[string]interface{}{
									"dataPoints": []interface{}{
										map[string]interface{}{
											"startTimeUnixNano": startTime,
											"timeUnixNano":      nanos,
											"count":             "5",
											"sum":               4.2,
											"bucketCounts":      []interface{}{"0", "1", "3", "1", "0"},
											"explicitBounds":    []interface{}{"0.1", "0.5", "1", "2"},
											"attributes": []interface{}{
												map[string]interface{}{
													"key":   "http.route",
													"value": map[string]interface{}{"stringValue": "/api/users"},
												},
											},
										},
									},
									"aggregationTemporality": 2, // Cumulative
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateLogRecords(res *resource.Resource) map[string]interface{} {
	now := time.Now()
	nanos := now.UnixNano()

	return map[string]interface{}{
		"resourceLogs": []interface{}{
			map[string]interface{}{
				"resource": map[string]interface{}{
					"attributes": resourceToAttributes(res),
				},
				"scopeLogs": []interface{}{
					map[string]interface{}{
						"scope": map[string]interface{}{
							"name":    "test-logger",
							"version": "1.0",
						},
						"logRecords": []interface{}{
							map[string]interface{}{
								"timeUnixNano":         nanos,
								"severityText":         "INFO",
								"severityNumber":       9,
								"body":                 map[string]interface{}{"stringValue": "Test log message"},
								"attributes":           []interface{}{},
								"observedTimeUnixNano": nanos,
							},
							map[string]interface{}{
								"timeUnixNano":   nanos + 1,
								"severityText":   "ERROR",
								"severityNumber": 17,
								"body":           map[string]interface{}{"stringValue": "Test error occurred"},
								"attributes": []interface{}{
									map[string]interface{}{
										"key":   "error.code",
										"value": map[string]interface{}{"stringValue": "500"},
									},
								},
								"observedTimeUnixNano": nanos + 1,
							},
						},
					},
				},
			},
		},
	}
}

func generateTraceData(res *resource.Resource) map[string]interface{} {
	now := time.Now()
	traceID := generateRandomHex(32)      // 16 bytes hex encoded
	spanID := generateRandomHex(16)       // 8 bytes hex encoded
	parentSpanID := generateRandomHex(16) // 8 bytes hex encoded

	return map[string]interface{}{
		"resourceSpans": []interface{}{
			map[string]interface{}{
				"resource": map[string]interface{}{
					"attributes": resourceToAttributes(res),
				},
				"scopeSpans": []interface{}{
					map[string]interface{}{
						"scope": map[string]interface{}{
							"name":    "test-tracer",
							"version": "1.0",
						},
						"spans": []interface{}{
							// Parent span
							map[string]interface{}{
								"traceId":           traceID,
								"spanId":            parentSpanID,
								"parentSpanId":      "",
								"name":              "test-parent-operation",
								"kind":              trace.SpanKindInternal,
								"startTimeUnixNano": now.UnixNano(),
								"endTimeUnixNano":   now.Add(100 * time.Millisecond).UnixNano(),
								"attributes": []interface{}{
									map[string]interface{}{
										"key":   "http.method",
										"value": map[string]interface{}{"stringValue": "GET"},
									},
									map[string]interface{}{
										"key":   "http.route",
										"value": map[string]interface{}{"stringValue": "/api/test"},
									},
								},
								"status": map[string]interface{}{
									"code": 0, // Unset
								},
								"traceState": "",
							},
							// Child span
							map[string]interface{}{
								"traceId":           traceID,
								"spanId":            spanID,
								"parentSpanId":      parentSpanID,
								"name":              "test-child-operation",
								"kind":              trace.SpanKindServer,
								"startTimeUnixNano": now.Add(10 * time.Millisecond).UnixNano(),
								"endTimeUnixNano":   now.Add(80 * time.Millisecond).UnixNano(),
								"attributes": []interface{}{
									map[string]interface{}{
										"key":   "db.system",
										"value": map[string]interface{}{"stringValue": "postgres"},
									},
									map[string]interface{}{
										"key":   "db.operation",
										"value": map[string]interface{}{"stringValue": "SELECT"},
									},
								},
								"status": map[string]interface{}{
									"code": 0,
								},
								"traceState": "",
								"events": []interface{}{
									map[string]interface{}{
										"timeUnixNano": now.Add(20 * time.Millisecond).UnixNano(),
										"name":         "test-cache-hit",
										"attributes": []interface{}{
											map[string]interface{}{
												"key":   "cache.key",
												"value": map[string]interface{}{"stringValue": "user:test123"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceToAttributes(res *resource.Resource) []interface{} {
	attrs := []interface{}{}
	for _, kv := range res.Attributes() {
		attrs = append(attrs, map[string]interface{}{
			"key":   string(kv.Key),
			"value": map[string]interface{}{"stringValue": kv.Value.AsString()},
		})
	}
	return attrs
}

func sendMetricsToCollector(metrics map[string]interface{}) error {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:4318/v1/metrics", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func sendLogsToCollector(logs map[string]interface{}) error {
	jsonData, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:4318/v1/logs", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func sendTracesToCollector(traces map[string]interface{}) error {
	jsonData, err := json.Marshal(traces)
	if err != nil {
		return fmt.Errorf("failed to marshal traces: %w", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:4318/v1/traces", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func generateRandomHex(length int) string {
	bytes := make([]byte, length/2) // 2 chars per byte
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
