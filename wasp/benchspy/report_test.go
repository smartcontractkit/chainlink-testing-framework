package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var lokiConfig = &wasp.LokiConfig{
	URL:       "http://localhost:3100",
	TenantID:  "test",
	BasicAuth: "user:pass",
}

func TestBenchSpy_NewStandardReport(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
			LokiConfig: lokiConfig,
		},
	}

	t.Run("successful creation (loki)", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &LokiQueryExecutor{}, report.QueryExecutors[0])
	})

	t.Run("successful creation (generator)", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Direct), WithGenerators(basicGen))
		require.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &DirectQueryExecutor{}, report.QueryExecutors[0])
	})

	t.Run("missing branch label", func(t *testing.T) {
		invalidGen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen",
				Labels: map[string]string{
					"commit": "abc123",
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
				},
			},
		}
		_, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(invalidGen))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing branch or commit labels")
	})

	t.Run("missing commit label", func(t *testing.T) {
		invalidGen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen",
				Labels: map[string]string{
					"branch": "abc123",
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
				},
			},
		}
		_, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(invalidGen))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing branch or commit labels")
	})

	t.Run("missing loki config", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen",
				Labels: map[string]string{
					"branch": "main",
					"commit": "abc123",
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
				},
				// no loki config
			},
		}

		report, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(gen))
		require.Nil(t, report)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loki config is missing")
	})

	t.Run("nil generator", func(t *testing.T) {
		_, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki))
		require.Error(t, err)
	})
}

func TestBenchSpy_NewStandardReportWithPrometheus(t *testing.T) {
	baseTime := time.Now()
	promConfig := &PrometheusConfig{
		Url:               "http://localhost:9090",
		NameRegexPatterns: []string{"node"},
	}

	validGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
			LokiConfig: lokiConfig,
		},
	}

	t.Run("successful prometheus creation", func(t *testing.T) {
		report, err := NewStandardReport("test-commit",
			WithStandardQueries(StandardQueryExecutor_Prometheus, StandardQueryExecutor_Loki),
			WithGenerators(validGen),
			WithPrometheusConfig(promConfig))
		require.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 2, len(report.QueryExecutors))
		assert.IsType(t, &LokiQueryExecutor{}, report.QueryExecutors[0])
		assert.IsType(t, &PrometheusQueryExecutor{}, report.QueryExecutors[1])
	})

	t.Run("successful prometheus creation (multiple name regex)", func(t *testing.T) {
		multiPromConfig := &PrometheusConfig{
			Url:               "http://localhost:9090",
			NameRegexPatterns: []string{"node", "eth"},
		}

		report, err := NewStandardReport("test-commit",
			WithStandardQueries(StandardQueryExecutor_Prometheus),
			WithGenerators(validGen),
			WithPrometheusConfig(multiPromConfig))
		require.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 2, len(report.QueryExecutors))
		assert.IsType(t, &PrometheusQueryExecutor{}, report.QueryExecutors[0])
		assert.IsType(t, &PrometheusQueryExecutor{}, report.QueryExecutors[1])
		firstAsProm := report.QueryExecutors[0].(*PrometheusQueryExecutor)
		assert.Equal(t, 6, len(firstAsProm.Queries))
		secondAsProm := report.QueryExecutors[0].(*PrometheusQueryExecutor)
		assert.Equal(t, 6, len(secondAsProm.Queries))
	})

	t.Run("invalid prometheus config (mising url)", func(t *testing.T) {
		invalidPromConfig := &PrometheusConfig{
			NameRegexPatterns: []string{"node"},
		}

		_, err := NewStandardReport("test-commit",
			WithStandardQueries(StandardQueryExecutor_Prometheus),
			WithGenerators(validGen),
			WithPrometheusConfig(invalidPromConfig),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "prometheus url is not set")
	})

	t.Run("invalid prometheus config (mising name regex)", func(t *testing.T) {
		invalidPromConfig := &PrometheusConfig{
			Url: "http://localhost:9090",
		}

		_, err := NewStandardReport("test-commit",
			WithStandardQueries(StandardQueryExecutor_Prometheus),
			WithGenerators(validGen),
			WithPrometheusConfig(invalidPromConfig),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "prometheus name regex patterns are not set")
	})
}

func TestBenchSpy_StandardReport_FetchData_WithMockExecutors(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}
	ctx := context.Background()

	t.Run("successful parallel execution", func(t *testing.T) {
		// Create mock executors that simulate successful execution
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.NoError(t, err)
		assert.True(t, exec1.ExecuteCalled)
		assert.True(t, exec2.ExecuteCalled)
	})

	t.Run("one executor fails execution", func(t *testing.T) {
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return fmt.Errorf("execution failed") },
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "execution failed")
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)

		exec1 := &MockQueryExecutor{
			ExecuteFn: func(ctx context.Context) error {
				cancel()
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn: func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			},
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("executor returns error", func(t *testing.T) {
		exec1 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		exec2 := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return fmt.Errorf("execution failed") },
			ValidateFn: func() error { return nil },
		}

		report := &StandardReport{
			BasicData: BasicData{
				TestStart:        baseTime,
				TestEnd:          baseTime.Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"basic": basicGen.Cfg},
			},
			QueryExecutors: []QueryExecutor{exec1, exec2},
		}

		err := report.FetchData(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "execution failed")
	})
}

// MockQueryExecutor implements QueryExecutor interface for testing
type MockQueryExecutor struct {
	ExecuteFn      func(context.Context) error
	ValidateFn     func() error
	TimeRangeFn    func(time.Time, time.Time)
	ValidateCalled bool
	ExecuteCalled  bool
	ResultsFn      func() map[string]interface{}
	KindFn         func() string
}

func (m *MockQueryExecutor) Kind() string {
	return m.KindFn()
}

func (m *MockQueryExecutor) Execute(ctx context.Context) error {
	m.ExecuteCalled = true
	return m.ExecuteFn(ctx)
}

func (m *MockQueryExecutor) Validate() error {
	m.ValidateCalled = true
	return m.ValidateFn()
}

func (m *MockQueryExecutor) TimeRange(start, end time.Time) {
	if m.TimeRangeFn != nil {
		m.TimeRangeFn(start, end)
	}
}

func (m *MockQueryExecutor) IsComparable(other QueryExecutor) error {
	return nil
}

func (m *MockQueryExecutor) Results() map[string]interface{} {
	return m.ResultsFn()
}

func TestBenchSpy_StandardReport_UnmarshalJSON(t *testing.T) {
	t.Run("valid loki executor", func(t *testing.T) {
		jsonData := `{
			"test_name": "test1",
			"commit_or_tag": "abc123",
			"query_executors": [{
				"kind": "loki",
				"queries": {
					"test query": "some query"
				},
                "query_results": {
                    "test query": ["1", "2", "3"]
                }
			}]
		}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.NoError(t, err)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &LokiQueryExecutor{}, report.QueryExecutors[0])
		asLoki := report.QueryExecutors[0].(*LokiQueryExecutor)
		assert.NotNil(t, asLoki.Queries["test query"])
		assert.Equal(t, "some query", asLoki.Queries["test query"])
		assert.Equal(t, 1, len(report.QueryExecutors[0].Results()))
		assert.IsType(t, []string{}, report.QueryExecutors[0].Results()["test query"])
		asStringSlice, err := ResultsAs([]string{}, report.QueryExecutors, StandardQueryExecutor_Loki, "test query")
		require.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3"}, asStringSlice["test query"])
	})

	t.Run("valid generator executor", func(t *testing.T) {
		jsonData := `{
            "test_name": "test1",
            "commit_or_tag": "abc123",
            "query_executors": [{
                "kind": "direct",
                "queries": [
					"test generator query"
				],
                "query_results": {
                    "test generator query": "1"
                }
            }]
        }`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.NoError(t, err)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &DirectQueryExecutor{}, report.QueryExecutors[0])
		asGenerator := report.QueryExecutors[0].(*DirectQueryExecutor)

		assert.Equal(t, 1, len(asGenerator.Queries))
		_, keyExists := asGenerator.Queries["test generator query"]
		assert.True(t, keyExists, "map should contain the key")
		assert.Nil(t, asGenerator.Queries["test generator query"])
		assert.Equal(t, 1, len(asGenerator.Results()))
		assert.IsType(t, "string", asGenerator.Results()["test generator query"])
		asStringSlice, err := ResultsAs("string", report.QueryExecutors, StandardQueryExecutor_Direct, "test generator query")
		require.NoError(t, err)
		assert.Equal(t, "1", asStringSlice["test generator query"])
	})

	t.Run("valid prometheus executor (vector)", func(t *testing.T) {
		jsonData := `{
    "test_name": "test1",
    "commit_or_tag": "abc123",
    "query_executors": [{
        "kind": "prometheus",
        "queries": {
            "rate": "rate(test_metric[5m])"
        },
        "query_results": {
            "rate": {
                "value": [{
                    "metric": {
                        "container_label_framework": "ctf",
                        "container_label_logging": "promtail",
                        "container_label_org_opencontainers_image_created": "2024-10-10T15:34:11Z",
                        "container_label_org_opencontainers_image_description": "\"node of the decentralized oracle network, bridging on and off-chain computation\"",
                        "container_label_org_opencontainers_image_licenses": "\"MIT\"",
                        "container_label_org_opencontainers_image_ref_name": "ubuntu",
                        "container_label_org_opencontainers_image_revision": "5ebb63266ca697f0649633641bbccb436c2c18bb",
                        "container_label_org_opencontainers_image_source": "\"https://github.com/smartcontractkit/chainlink\"",
                        "container_label_org_opencontainers_image_title": "chainlink",
                        "container_label_org_opencontainers_image_url": "\"https://github.com/smartcontractkit/chainlink\"",
                        "container_label_org_opencontainers_image_version": "2.17.0",
                        "container_label_org_testcontainers": "true",
                        "container_label_org_testcontainers_lang": "go",
                        "container_label_org_testcontainers_reap": "true",
                        "container_label_org_testcontainers_sessionId": "0e438f13ded27fcd3f85123134091358f9ce5c575c0b14a6f3c8998b4d2e7d14",
                        "container_label_org_testcontainers_version": "0.34.0",
                        "cpu": "total",
                        "id": "/docker/26c2319b0e3f5c4f6103b15c9e656fad1635356c39222d5cbd1076d4d49a1375",
                        "image": "public.ecr.aws/chainlink/chainlink:v2.17.0-arm64",
                        "instance": "cadvisor:8080",
                        "job": "cadvisor",
                        "name": "node3"
                    },
                    "value": [
                        1733919891.792,
                        "0.39449004082983885"
                        ]
					}],
                "metric_type": "vector"
        	}
    	}
	}]
}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.NoError(t, err)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &PrometheusQueryExecutor{}, report.QueryExecutors[0])
		asProm := report.QueryExecutors[0].(*PrometheusQueryExecutor)
		assert.NotNil(t, asProm.Queries["rate"])
		assert.Equal(t, "rate(test_metric[5m])", asProm.Queries["rate"])
		assert.Equal(t, 1, len(report.QueryExecutors[0].Results()))

		asValue := asProm.MustResultsAsValue()

		assert.IsType(t, model.Vector{}, asValue["rate"])
		asVector := asValue["rate"].(model.Vector)

		assert.Equal(t, 1, len(asVector))
		assert.Equal(t, 0.39449004082983885, float64(asVector[0].Value))
	})

	t.Run("valid prometheus executor (matrix)", func(t *testing.T) {
		jsonData := `{
	"test_name":"test1",
	"commit_or_tag":"abc123",
	"query_executors":[
		{
			"kind":"prometheus",
			"queries":{
				"rate":"rate(test_metric[5m])"
			},
			"query_results":{
				"rate":{
					"value":[{
								"metric": {
									"container_label_framework": "ctf",
									"container_label_logging": "promtail",
									"container_label_org_opencontainers_image_created": "2024-10-10T15:34:11Z",
									"container_label_org_opencontainers_image_description": "\"node of the decentralized oracle network, bridging on and off-chain computation\"",
									"container_label_org_opencontainers_image_licenses": "\"MIT\"",
									"container_label_org_opencontainers_image_ref_name": "ubuntu",
									"container_label_org_opencontainers_image_revision": "5ebb63266ca697f0649633641bbccb436c2c18bb",
									"container_label_org_opencontainers_image_source": "\"https://github.com/smartcontractkit/chainlink\"",
									"container_label_org_opencontainers_image_title": "chainlink",
									"container_label_org_opencontainers_image_url": "\"https://github.com/smartcontractkit/chainlink\"",
									"container_label_org_opencontainers_image_version": "2.17.0",
									"container_label_org_testcontainers": "true",
									"container_label_org_testcontainers_lang": "go",
									"container_label_org_testcontainers_reap": "true",
									"container_label_org_testcontainers_sessionId": "0e438f13ded27fcd3f85123134091358f9ce5c575c0b14a6f3c8998b4d2e7d14",
									"container_label_org_testcontainers_version": "0.34.0",
									"cpu": "total",
									"id": "/docker/a2f7ab689d98d06f941732a3c04eae867ce56de687d87c7fd1bb1ac8c36a415a",
									"image": "public.ecr.aws/chainlink/chainlink:v2.17.0-arm64",
									"instance": "cadvisor:8080",
									"job": "cadvisor",
									"name": "node1"
								},
								"values":[[
										1734010920,
										"0.004647525277233765"
								]]
							}],
					"metric_type":"matrix"
					}
				}
			}
		]
	}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.NoError(t, err)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &PrometheusQueryExecutor{}, report.QueryExecutors[0])
		asProm := report.QueryExecutors[0].(*PrometheusQueryExecutor)
		assert.NotNil(t, asProm.Queries["rate"])
		assert.Equal(t, "rate(test_metric[5m])", asProm.Queries["rate"])
		assert.Equal(t, 1, len(report.QueryExecutors[0].Results()))

		asValue := asProm.MustResultsAsValue()

		assert.IsType(t, model.Matrix{}, asValue["rate"])
		asMatrix := asValue["rate"].(model.Matrix)

		assert.Equal(t, 1, len(asMatrix))
		assert.Equal(t, 1, len(asMatrix[0].Values))
		assert.Equal(t, 0.004647525277233765, float64(asMatrix[0].Values[0].Value))
	})

	t.Run("valid prometheus executor (scalar)", func(t *testing.T) {
		jsonData := `{
    "test_name": "test1",
    "commit_or_tag": "abc123",
    "query_executors": [{
        "kind": "prometheus",
        "queries": {
            "rate": "scalar(quantile(0.95, rate(container_cpu_usage_seconds_total{}[5m])) * 100)"
        },
        "query_results": {
            "rate": {
				"value": [
					1734012682.065,
					"0.4631138973188853"
				],
				"metric_type": "scalar"
			}
		}
	}]
}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.NoError(t, err)
		assert.Equal(t, 1, len(report.QueryExecutors))
		assert.IsType(t, &PrometheusQueryExecutor{}, report.QueryExecutors[0])
		asProm := report.QueryExecutors[0].(*PrometheusQueryExecutor)
		assert.NotNil(t, asProm.Queries["rate"])
		assert.Equal(t, "scalar(quantile(0.95, rate(container_cpu_usage_seconds_total{}[5m])) * 100)", asProm.Queries["rate"])
		assert.Equal(t, 1, len(report.QueryExecutors[0].Results()))

		asValue := asProm.MustResultsAsValue()

		assert.IsType(t, &model.Scalar{}, asValue["rate"])
		asScalar := asValue["rate"].(*model.Scalar)

		assert.Equal(t, 0.4631138973188853, float64(asScalar.Value))
	})

	t.Run("unknown executor type", func(t *testing.T) {
		jsonData := `{
			"test_name": "test1",
			"commit_or_tag": "abc123",
			"query_executors": [{
				"kind": "unknown",
				"query": "test query"
			}]
		}`

		var report StandardReport
		err := json.Unmarshal([]byte(jsonData), &report)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown query executor type")
	})
}

func TestBenchSpy_StandardReport_FetchData(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
		},
	}

	t.Run("valid fetch", func(t *testing.T) {
		basicGen.Cfg.LokiConfig = lokiConfig
		report, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)

		mockExec := &MockQueryExecutor{
			ExecuteFn:  func(ctx context.Context) error { return nil },
			ValidateFn: func() error { return nil },
		}
		report.QueryExecutors = []QueryExecutor{mockExec}

		err = report.FetchData(context.Background())
		require.NoError(t, err)
		assert.True(t, mockExec.ExecuteCalled)
	})
}

func TestBenchSpy_StandardReport_IsComparable(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
			LokiConfig: lokiConfig,
		},
	}

	t.Run("matching reports", func(t *testing.T) {
		report1, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)
		report2, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)

		err = report1.IsComparable(report2)
		require.NoError(t, err)
	})

	t.Run("different report types", func(t *testing.T) {
		report1, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)

		// Create a mock reporter that implements Reporter interface
		mockReport := &MockReport{}

		err = report1.IsComparable(mockReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected type *StandardReport")
	})

	t.Run("different executors", func(t *testing.T) {
		report1, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)

		// Create second report with different executor
		diffGen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "test-gen-2", // different name
				Labels: map[string]string{
					"branch": "main",
					"commit": "def456", // different commit
				},
				Schedule: []*wasp.Segment{
					{StartTime: baseTime, EndTime: baseTime.Add(2 * time.Hour)}, // different duration
				},
				LokiConfig: lokiConfig,
			},
		}
		report2, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(diffGen))
		require.NoError(t, err)

		err = report1.IsComparable(report2)
		require.Error(t, err)
	})
}

// MockReport implements Reporter interface for testing
type MockReport struct{}

func (m *MockReport) FetchData(ctx context.Context) error      { return nil }
func (m *MockReport) Store() (string, error)                   { return "", nil }
func (m *MockReport) Load(testName, commitOrTag string) error  { return nil }
func (m *MockReport) LoadLatest(testName string) error         { return nil }
func (m *MockReport) IsComparable(other Reporter) error        { return nil }
func (m *MockReport) FetchResources(ctx context.Context) error { return nil }

func TestBenchSpy_StandardReport_Store_Load(t *testing.T) {
	// Setup test directory with git repo
	tmpDir := t.TempDir()

	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
			},
			LokiConfig: lokiConfig,
		},
	}

	t.Run("store and load", func(t *testing.T) {
		report, err := NewStandardReport("test-commit", WithStandardQueries(StandardQueryExecutor_Loki), WithGenerators(basicGen))
		require.NoError(t, err)

		storage := LocalStorage{Directory: tmpDir}
		report.LocalStorage = storage

		_, err = report.Store()
		require.NoError(t, err)

		loadedReport := &StandardReport{LocalStorage: storage}
		err = loadedReport.Load(report.TestName, report.CommitOrTag)
		require.NoError(t, err)

		assert.Equal(t, report.TestName, loadedReport.TestName)
		assert.Equal(t, report.CommitOrTag, loadedReport.CommitOrTag)
	})
}

func TestBenchSpy_StandardReport_ResultsAs(t *testing.T) {
	t.Run("successful type conversion for float64", func(t *testing.T) {
		mockExecutor := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": float64(123.45),
					"query2": float64(678.90),
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		results, err := ResultsAs(float64(0), []QueryExecutor{mockExecutor}, StandardQueryExecutor_Loki)
		require.NoError(t, err)
		assert.Equal(t, float64(123.45), results["query1"])
		assert.Equal(t, float64(678.90), results["query2"])
	})

	t.Run("successful type conversion with specific query names", func(t *testing.T) {
		mockExecutor := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": "result1",
					"query2": "result2",
					"query3": "result3",
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		results, err := ResultsAs("", []QueryExecutor{mockExecutor}, StandardQueryExecutor_Loki, "query1", "query3")
		require.NoError(t, err)
		assert.Equal(t, "result1", results["query1"])
		assert.Equal(t, "result3", results["query3"])
		assert.Empty(t, results["query2"])
	})

	t.Run("failed type conversion", func(t *testing.T) {
		mockExecutor := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": "result1",
					"query2": 2,
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		results, err := ResultsAs("", []QueryExecutor{mockExecutor}, StandardQueryExecutor_Loki, "query1", "query2")
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "failed to cast result to type string")
	})

	t.Run("duplicated query name", func(t *testing.T) {
		firstMockExecutor := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": "result1",
					"query2": "result2",
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		secondMockExecutor := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query2": "result2",
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		results, err := ResultsAs("", []QueryExecutor{firstMockExecutor, secondMockExecutor}, StandardQueryExecutor_Loki)
		require.Error(t, err)
		require.Nil(t, results)
		require.Contains(t, err.Error(), "results for query 'query2' already exist")
	})
}

func TestBenchSpy_ConvertQueryResults(t *testing.T) {
	t.Run("successful type conversions", func(t *testing.T) {
		type customType struct{}

		t.Run("int conversion", func(t *testing.T) {
			input := map[string]interface{}{
				"int_val": 123,
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, 0, result["int_val"])
			assert.Equal(t, 123, result["int_val"])
		})

		t.Run("string conversion", func(t *testing.T) {
			input := map[string]interface{}{
				"str_val": "test string",
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, "", result["str_val"])
			assert.Equal(t, "test string", result["str_val"])
		})

		t.Run("float64 conversion", func(t *testing.T) {
			input := map[string]interface{}{
				"float_val": 123.45,
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, float64(0), result["float_val"])
			assert.Equal(t, 123.45, result["float_val"])
		})

		t.Run("[]int conversion", func(t *testing.T) {
			input := map[string]interface{}{
				"int_slice": []interface{}{1, 2, 3},
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, []int{}, result["int_slice"])
			assert.Equal(t, []int{1, 2, 3}, result["int_slice"])
		})

		t.Run("[]string conversion", func(t *testing.T) {
			input := map[string]interface{}{
				"str_slice": []interface{}{"a", "b", "c"},
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, []string{}, result["str_slice"])
			assert.Equal(t, []string{"a", "b", "c"}, result["str_slice"])
		})

		t.Run("float64 slice conversion", func(t *testing.T) {
			input := map[string]interface{}{
				"float_slice": []interface{}{1.1, 2.2, 3.3},
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, []float64{}, result["float_slice"])
			assert.Equal(t, []float64{1.1, 2.2, 3.3}, result["float_slice"])
		})

		t.Run("no conversion (not all are int)", func(t *testing.T) {
			input := map[string]interface{}{
				"interface_slice": []interface{}{1, "a", 2.2},
				"other":           customType{},
			}

			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, customType{}, result["other"])
			assert.Equal(t, customType{}, result["other"])
			assert.Equal(t, []interface{}{1, "a", 2.2}, result["interface_slice"])
			assert.IsType(t, []interface{}{}, result["interface_slice"])
		})

		t.Run("no conversion (not all are string)", func(t *testing.T) {
			input := map[string]interface{}{
				"interface_slice": []interface{}{"a", 2.2},
			}
			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.Equal(t, []interface{}{"a", 2.2}, result["interface_slice"])
			assert.IsType(t, []interface{}{}, result["interface_slice"])
		})

		t.Run("no conversion (not all are float64)", func(t *testing.T) {
			input := map[string]interface{}{
				"interface_slice": []interface{}{1.2, "a", 2},
			}

			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.Equal(t, []interface{}{1.2, "a", 2}, result["interface_slice"])
			assert.IsType(t, []interface{}{}, result["interface_slice"])
		})

		t.Run("no conversion (custom type)", func(t *testing.T) {
			input := map[string]interface{}{
				"other": customType{},
			}

			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, customType{}, result["other"])
			assert.Equal(t, customType{}, result["other"])
		})

		t.Run("partial", func(t *testing.T) {
			input := map[string]interface{}{
				"other":     customType{},
				"str_slice": []interface{}{"a", "b", "c"},
			}

			result, err := convertQueryResults(input)
			require.NoError(t, err)
			assert.IsType(t, customType{}, result["other"])
			assert.Equal(t, customType{}, result["other"])
			assert.IsType(t, []string{}, result["str_slice"])
			assert.Equal(t, []string{"a", "b", "c"}, result["str_slice"])
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("nil input", func(t *testing.T) {
			result, err := convertQueryResults(nil)
			require.NoError(t, err)
			assert.Equal(t, map[string]interface{}{}, result)
		})
	})
}

func TestBenchSpy_MustAllResults(t *testing.T) {
	t.Run("MustAllLokiResults", func(t *testing.T) {
		mockLokiExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": []string{"log1", "log2"},
					"query2": []string{"log3", "log4"},
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{mockLokiExec},
		}

		results := MustAllLokiResults(sr)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, []string{"log1", "log2"}, results["query1"])
		assert.Equal(t, []string{"log3", "log4"}, results["query2"])
	})

	t.Run("MustAllLokiResults - two executors", func(t *testing.T) {
		firstMockLokiExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": []string{"log1", "log2"},
					"query2": []string{"log3", "log4"},
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		secondMockLokiExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query3": []string{"log5", "log6"},
					"query4": []string{"log7", "log8"},
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{firstMockLokiExec, secondMockLokiExec},
		}

		results := MustAllLokiResults(sr)
		assert.Equal(t, 4, len(results))
		assert.Equal(t, []string{"log1", "log2"}, results["query1"])
		assert.Equal(t, []string{"log3", "log4"}, results["query2"])
		assert.Equal(t, []string{"log5", "log6"}, results["query3"])
		assert.Equal(t, []string{"log7", "log8"}, results["query4"])
	})

	t.Run("MustAllGeneratorResults", func(t *testing.T) {
		mockGenExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": 1.0,
					"query2": 2.0,
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Direct)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{mockGenExec},
		}

		results := MustAllDirectResults(sr)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, 1.0, results["query1"])
		assert.Equal(t, 2.0, results["query2"])
	})

	t.Run("MustAllPrometheusResults", func(t *testing.T) {
		vector1 := model.Vector{
			&model.Sample{
				Value: 1.23,
			},
		}
		vector2 := model.Vector{
			&model.Sample{
				Value: 4.56,
			},
		}

		mockPromExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": vector1,
					"query2": vector2,
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Prometheus)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{mockPromExec},
		}

		results := MustAllPrometheusResults(sr)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, model.Vector(vector1), results["query1"].(model.Vector))
		assert.Equal(t, model.Vector(vector2), results["query2"].(model.Vector))
	})

	t.Run("MustAllLokiResults panics on wrong type", func(t *testing.T) {
		mockExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": 123, // wrong type
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{mockExec},
		}

		assert.Panics(t, func() {
			MustAllLokiResults(sr)
		})
	})

	t.Run("MustAllGeneratorResults panics on wrong type", func(t *testing.T) {
		mockExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"query1": []string{"wrong", "type"},
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Direct)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{mockExec},
		}

		assert.Panics(t, func() {
			MustAllDirectResults(sr)
		})
	})

	t.Run("Results from mixed executors", func(t *testing.T) {
		mockLokiExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"loki_query": []string{"log1", "log2"},
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Loki)
			},
		}

		mockGenExec := &MockQueryExecutor{
			ResultsFn: func() map[string]interface{} {
				return map[string]interface{}{
					"gen_query": 1.0,
				}
			},
			KindFn: func() string {
				return string(StandardQueryExecutor_Direct)
			},
		}

		sr := &StandardReport{
			QueryExecutors: []QueryExecutor{mockLokiExec, mockGenExec},
		}

		lokiResults := MustAllLokiResults(sr)
		assert.Equal(t, 1, len(lokiResults))
		assert.Equal(t, []string{"log1", "log2"}, lokiResults["loki_query"])

		genResults := MustAllDirectResults(sr)
		assert.Equal(t, 1, len(genResults))
		assert.Equal(t, 1.0, genResults["gen_query"])
	})
}

func TestBenchSpy_FetchNewReportAndLoadLatestPrevious(t *testing.T) {
	baseTime := time.Now()
	basicGen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour), From: 1, Duration: 2 * time.Second},
			},
			LokiConfig: lokiConfig,
		},
	}

	t.Run("successful execution", func(t *testing.T) {
		tmpDir := t.TempDir()

		cfg := &wasp.Config{
			T:        t,
			LoadType: wasp.RPS,
			GenName:  "test-gen",
			Labels: map[string]string{
				"branch": "main",
				"commit": "abc123",
			},
			Schedule: []*wasp.Segment{
				{StartTime: baseTime, EndTime: baseTime.Add(time.Hour), From: 1, Duration: 2 * time.Second, Type: wasp.SegmentType_Plain},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 4,
			maxFailures:  3,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun

		gen, err := wasp.NewGenerator(cfg)
		require.NoError(t, err)

		gen.Run(true)

		prevReport, err := NewStandardReport("a7fc5826a572c09f8b93df3b9f674113372ce924",
			WithStandardQueries(StandardQueryExecutor_Direct),
			WithGenerators(gen),
			WithReportDirectory(tmpDir))
		require.NoError(t, err)
		_, err = prevReport.Store()
		require.NoError(t, err)

		newReport, prevLoadedReport, err := FetchNewStandardReportAndLoadLatestPrevious(
			context.Background(),
			"new-commit",
			WithStandardQueries(StandardQueryExecutor_Direct),
			WithGenerators(gen),
			WithReportDirectory(tmpDir),
		)
		require.NoError(t, err)

		assert.NotNil(t, newReport)
		assert.NotNil(t, prevLoadedReport)
		assert.Equal(t, "a7fc5826a572c09f8b93df3b9f674113372ce924", prevLoadedReport.CommitOrTag)
	})

	t.Run("no previous report", func(t *testing.T) {
		tmpDir := t.TempDir()

		basicGen.Cfg.T = t
		newReport, prevReport, err := FetchNewStandardReportAndLoadLatestPrevious(
			context.Background(),
			"new-commit-7",
			WithStandardQueries(StandardQueryExecutor_Direct),
			WithGenerators(basicGen),
			WithReportDirectory(tmpDir),
		)

		assert.Error(t, err)
		assert.Nil(t, newReport)
		assert.Nil(t, prevReport)
	})
}
func TestBenchSpy_CompareDirectWithThresholds(t *testing.T) {
	t.Run("metrics within thresholds", func(t *testing.T) {
		previousReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       100.0,
							string(Percentile95Latency): 200.0,
							string(MaxLatency):          300.0,
							string(ErrorRate):           1.0,
						}
					},
				},
			},
		}

		currentReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       105.0, // 5% increase
							string(Percentile95Latency): 210.0, // 5% increase
							string(MaxLatency):          315.0, // 5% increase
							string(ErrorRate):           1.05,  // 5% increase
						}
					},
				},
			},
		}

		failed, errs := CompareDirectWithThresholds(10.0, 10.0, 10.0, 10.0, previousReport, currentReport)
		assert.False(t, failed)
		assert.Empty(t, errs)
	})

	t.Run("one metric exceed thresholds", func(t *testing.T) {
		previousReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       100.0,
							string(Percentile95Latency): 200.0,
							string(MaxLatency):          300.0,
							string(ErrorRate):           1.0,
						}
					},
				},
			},
		}

		currentReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       150.0, // 50% increase
							string(Percentile95Latency): 200.0, // no increase
							string(MaxLatency):          300.0, // no increase
							string(ErrorRate):           1.0,   // no increase
						}
					},
				},
			},
		}

		failed, errs := CompareDirectWithThresholds(10.0, 1.0, 1.0, 1.0, previousReport, currentReport)
		assert.True(t, failed)
		assert.Len(t, errs, 1)
		for _, err := range errs {
			assert.Contains(t, err.Error(), "different, which is higher than the threshold")
		}
	})

	t.Run("all metrics exceed thresholds", func(t *testing.T) {
		previousReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       100.0,
							string(Percentile95Latency): 200.0,
							string(MaxLatency):          300.0,
							string(ErrorRate):           1.0,
						}
					},
				},
			},
		}

		currentReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       150.0, // 50% increase
							string(Percentile95Latency): 300.0, // 50% increase
							string(MaxLatency):          450.0, // 50% increase
							string(ErrorRate):           2.0,   // 100% increase
						}
					},
				},
			},
		}

		failed, errs := CompareDirectWithThresholds(10.0, 10.0, 10.0, 10.0, previousReport, currentReport)
		assert.True(t, failed)
		assert.Len(t, errs, 4)
		for _, err := range errs {
			assert.Contains(t, err.Error(), "different, which is higher than the threshold")
		}
	})

	t.Run("handle zero values", func(t *testing.T) {
		previousReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       0.0,
							string(Percentile95Latency): 0.0,
							string(MaxLatency):          0.0,
							string(ErrorRate):           0.0,
						}
					},
				},
			},
		}

		currentReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       0.0,
							string(Percentile95Latency): 0.0,
							string(MaxLatency):          0.0,
							string(ErrorRate):           0.0,
						}
					},
				},
			},
		}

		failed, errs := CompareDirectWithThresholds(10.0, 10.0, 10.0, 10.0, previousReport, currentReport)
		assert.False(t, failed)
		assert.Empty(t, errs)
	})

	t.Run("handle missing metrics", func(t *testing.T) {
		previousReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency): 100.0,
							// missing other metrics
						}
					},
				},
			},
		}

		currentReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency): 105.0,
							// missing other metrics
						}
					},
				},
			},
		}

		failed, errs := CompareDirectWithThresholds(10.0, 10.0, 10.0, 10.0, previousReport, currentReport)
		assert.True(t, failed)
		assert.Len(t, errs, 3) // Should have errors for missing P95, Max, and Error Rate
		for _, err := range errs {
			assert.Contains(t, err.Error(), "results were missing")
		}
	})

	t.Run("handle zero to non-zero transition", func(t *testing.T) {
		previousReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       0.0,
							string(Percentile95Latency): 0.0,
							string(MaxLatency):          0.0,
							string(ErrorRate):           0.0,
						}
					},
				},
			},
		}

		currentReport := &StandardReport{
			QueryExecutors: []QueryExecutor{
				&MockQueryExecutor{
					KindFn: func() string { return string(StandardQueryExecutor_Direct) },
					ResultsFn: func() map[string]interface{} {
						return map[string]interface{}{
							string(MedianLatency):       100.0,
							string(Percentile95Latency): 200.0,
							string(MaxLatency):          300.0,
							string(ErrorRate):           1.0,
						}
					},
				},
			},
		}

		failed, errs := CompareDirectWithThresholds(10.0, 10.0, 10.0, 10.0, previousReport, currentReport)
		assert.True(t, failed)
		assert.Len(t, errs, 4)
		for _, err := range errs {
			assert.Contains(t, err.Error(), "100.0000% different")
		}
	})
}
