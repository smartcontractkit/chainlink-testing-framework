package benchspy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBenchSpy_NewGeneratorQueryExecutor(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				GenName: "my_gen",
				Schedule: []*wasp.Segment{
					{
						Type:     "plain",
						From:     1,
						Duration: time.Minute,
					},
				},
			},
		}
		executor, err := NewGeneratorQueryExecutor(gen)
		assert.NoError(t, err)
		assert.Equal(t, "generator", executor.KindName)
		assert.Equal(t, gen, executor.Generator)
		assert.NotEmpty(t, executor.Queries)
		assert.NotNil(t, executor.QueryResults)

		err = executor.Validate()
		assert.NoError(t, err)
	})

	t.Run("nil generator", func(t *testing.T) {
		executor, err := NewGeneratorQueryExecutor(nil)
		assert.NoError(t, err)

		err = executor.Validate()
		assert.Error(t, err)
	})
}

func TestBenchSpy_GeneratorQueryExecutor_Results(t *testing.T) {
	expected := map[string]interface{}{
		"test": "result",
	}
	executor := &GeneratorQueryExecutor{
		QueryResults: expected,
	}
	assert.Equal(t, expected, executor.Results())
}

func TestBenchSpy_GeneratorQueryExecutor_IsComparable(t *testing.T) {
	baseGen := &wasp.Generator{
		Cfg: &wasp.Config{
			GenName: "my_gen",
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: time.Minute,
				},
			},
		},
	}

	t.Run("same configs", func(t *testing.T) {
		exec1, _ := NewGeneratorQueryExecutor(baseGen)
		exec2, _ := NewGeneratorQueryExecutor(baseGen)
		assert.NoError(t, exec1.IsComparable(exec2))
	})

	t.Run("different query count", func(t *testing.T) {
		exec1, _ := NewGeneratorQueryExecutor(baseGen)
		exec2, _ := NewGeneratorQueryExecutor(baseGen)
		exec2.Queries = map[string]GeneratorQueryFn{"test": nil}

		assert.Error(t, exec1.IsComparable(exec2))
	})

	t.Run("different queries", func(t *testing.T) {
		exec1, _ := NewGeneratorQueryExecutor(baseGen)
		exec2, _ := NewGeneratorQueryExecutor(baseGen)
		exec2.Queries = map[string]GeneratorQueryFn{"test": nil, "test2": nil, "test3": nil}

		assert.Error(t, exec1.IsComparable(exec2))
	})

	t.Run("different configs", func(t *testing.T) {
		exec1, _ := NewGeneratorQueryExecutor(baseGen)
		differentGen := &wasp.Generator{
			Cfg: &wasp.Config{
				GenName: "my_gen",
				Schedule: []*wasp.Segment{
					{
						Type:     "plain",
						From:     2,
						Duration: time.Minute,
					},
				},
			},
		}
		exec2, _ := NewGeneratorQueryExecutor(differentGen)
		assert.Error(t, exec1.IsComparable(exec2))
	})

	t.Run("different types", func(t *testing.T) {
		exec1, _ := NewGeneratorQueryExecutor(baseGen)
		exec2 := &LokiQueryExecutor{}
		assert.Error(t, exec1.IsComparable(exec2))
	})
}

func TestBenchSpy_GeneratorQueryExecutor_Validate(t *testing.T) {
	t.Run("valid case", func(t *testing.T) {
		executor, _ := NewGeneratorQueryExecutor(&wasp.Generator{
			Cfg: &wasp.Config{},
		})
		assert.NoError(t, executor.Validate())
	})

	t.Run("missing generator", func(t *testing.T) {
		executor := &GeneratorQueryExecutor{
			Queries: map[string]GeneratorQueryFn{"test": nil},
		}
		assert.Error(t, executor.Validate())
	})

	t.Run("empty queries", func(t *testing.T) {
		executor := &GeneratorQueryExecutor{
			Generator: &wasp.Generator{},
			Queries:   map[string]GeneratorQueryFn{},
		}
		assert.Error(t, executor.Validate())
	})
}

type fakeGun struct {
	maxSuccesses   int
	successCounter int
	maxFailures    int
	failureCounter int
	successMutex   sync.Mutex
	failureMutex   sync.Mutex

	schedule *wasp.Segment
}

func (f *fakeGun) Call(l *wasp.Generator) *wasp.Response {
	now := time.Now()
	f.successMutex.Lock()
	defer f.successMutex.Unlock()
	if f.successCounter < f.maxSuccesses {
		f.successCounter++
		d := time.Duration(150) * time.Millisecond
		time.Sleep(d)
		finishedAt := time.Now()
		return &wasp.Response{
			StartedAt:  &now,
			Duration:   d,
			FinishedAt: &finishedAt,
			Failed:     false,
			Data:       fmt.Sprintf("success[%d]", f.successCounter),
		}
	}

	f.failureMutex.Lock()
	defer f.failureMutex.Unlock()
	if f.failureCounter < f.maxFailures {
		f.failureCounter++
		d := time.Duration(200) * time.Millisecond
		time.Sleep(d)
		finishedAt := time.Now()
		return &wasp.Response{
			StartedAt:  &now,
			Duration:   d,
			FinishedAt: &finishedAt,
			Failed:     true,
			Data:       fmt.Sprintf("failure[%d]", f.failureCounter),
		}
	}

	panic(fmt.Sprintf("fakeGun.Call called too many times (%d vs %d). Expected maxFailures: %d. Expected maxSuccesses: %d. do adjust your settings (sum of maxSuccesses and maxFilure should be greater than the duration of the segment in seconds)", f.maxFailures+f.maxSuccesses, f.maxFailures+f.maxSuccesses+1, f.maxFailures, f.maxSuccesses))
}

func TestBenchSpy_GeneratorQueryExecutor_Execute(t *testing.T) {
	t.Run("success case with mixed responses", func(t *testing.T) {
		cfg := &wasp.Config{
			GenName:  "my_gen",
			LoadType: wasp.RPS,
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: 5 * time.Second,
				},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 4,
			maxFailures:  3,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun

		gen, err := wasp.NewGenerator(cfg)
		assert.NoError(t, err)

		gen.Run(true)

		require.Equal(t, 4, len(gen.GetData().OKResponses.Data), "expected 4 successful responses")
		require.GreaterOrEqual(t, len(gen.GetData().FailResponses.Data), 2, "expected >=2 failed responses")

		actualFailures := len(gen.GetData().FailResponses.Data)

		executor, err := NewGeneratorQueryExecutor(gen)
		assert.NoError(t, err)

		err = executor.Execute(context.Background())
		assert.NoError(t, err)

		results := executor.Results()
		assert.NotEmpty(t, results)

		// 4 responses with ~150ms latency (150ms sleep + some execution overhead)
		// and 2-3 responses with ~200ms latency (200ms sleep + some execution overhead)
		// expected median latency: (150ms, 151ms>
		resultsAsStrings, err := ResultsAs("string", []QueryExecutor{executor}, StandardQueryExecutor_Generator, string(MedianLatency), string(Percentile95Latency), string(ErrorRate))
		assert.NoError(t, err)
		require.Equal(t, 3, len(resultsAsStrings))

		medianLatencyFloat, err := strconv.ParseFloat(resultsAsStrings[string(MedianLatency)], 64)
		assert.NoError(t, err)
		require.InDelta(t, 151.0, medianLatencyFloat, 1.0)

		// since we have 2-3 responses with 200-201ms latency, the 95th percentile should be (200ms, 201ms>
		p95LatencyFloat, err := strconv.ParseFloat(resultsAsStrings[string(Percentile95Latency)], 64)
		assert.NoError(t, err)
		require.InDelta(t, 201.0, p95LatencyFloat, 1.0)

		errorRate, exists := resultsAsStrings[string(ErrorRate)]
		assert.True(t, exists)

		// error rate is the number of failures divided by the total number of responses
		expectedErrorRate := float64(actualFailures) / (float64(fakeGun.maxSuccesses) + float64(actualFailures))
		assert.Equal(t, fmt.Sprintf("%.4f", expectedErrorRate), errorRate)
	})

	t.Run("all responses failed", func(t *testing.T) {
		cfg := &wasp.Config{
			GenName:  "my_gen",
			LoadType: wasp.RPS,
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: 5 * time.Second,
				},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 0,
			maxFailures:  10,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun

		gen, err := wasp.NewGenerator(cfg)
		assert.NoError(t, err)

		gen.Run(true)

		require.Equal(t, 0, len(gen.GetData().OKResponses.Data), "expected 0 successful responses")
		require.GreaterOrEqual(t, len(gen.GetData().FailResponses.Data), 5, "expected >=5 failed responses")

		executor, err := NewGeneratorQueryExecutor(gen)
		assert.NoError(t, err)

		err = executor.Execute(context.Background())
		assert.NoError(t, err)

		results := executor.Results()
		assert.NotEmpty(t, results)

		errorRate, exists := results[string(ErrorRate)]
		assert.True(t, exists)
		assert.Equal(t, "1.0000", errorRate)
	})

	t.Run("no responses", func(t *testing.T) {
		cfg := &wasp.Config{
			GenName:  "my_gen",
			LoadType: wasp.RPS,
			Schedule: []*wasp.Segment{
				{
					Type:     "plain",
					From:     1,
					Duration: 5 * time.Second,
				},
			},
		}

		fakeGun := &fakeGun{
			maxSuccesses: 0,
			maxFailures:  0,
			schedule:     cfg.Schedule[0],
		}

		cfg.Gun = fakeGun
		gen, err := wasp.NewGenerator(cfg)

		executor, _ := NewGeneratorQueryExecutor(gen)
		err = executor.Execute(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no responses found for generator")
	})
}

func TestBenchSpy_GeneratorQueryExecutor_MarshalJSON(t *testing.T) {
	t.Run("marshal/unmarshal round trip", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				GenName: "my_gen",
				Schedule: []*wasp.Segment{
					{
						Type:     "plain",
						From:     1,
						Duration: time.Minute,
					},
				},
			},
		}
		original, _ := NewGeneratorQueryExecutor(gen)
		original.QueryResults["test"] = "result"
		original.QueryResults["test2"] = "1"

		original.Queries = map[string]GeneratorQueryFn{
			"test": func(responses *wasp.SliceBuffer[wasp.Response]) (string, error) {
				return "result", nil
			},
			"test2": func(responses *wasp.SliceBuffer[wasp.Response]) (string, error) {
				return "1", nil
			},
		}

		data, err := json.Marshal(original)
		assert.NoError(t, err)

		var recovered GeneratorQueryExecutor
		err = json.Unmarshal(data, &recovered)
		assert.NoError(t, err)

		assert.Equal(t, original.KindName, recovered.KindName)
		assert.Equal(t, original.QueryResults, recovered.QueryResults)
		assert.Equal(t, len(original.Queries), len(recovered.Queries))
	})

	t.Run("marshal with nil generator", func(t *testing.T) {
		executor := &GeneratorQueryExecutor{}
		data, err := json.Marshal(executor)
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"generator_config":null`)
	})

	t.Run("unmarshal invalid JSON", func(t *testing.T) {
		var executor GeneratorQueryExecutor
		err := json.Unmarshal([]byte(`{invalid json}`), &executor)
		assert.Error(t, err)
	})
}

func TestBenchSpy_GeneratorQueryExecutor_TimeRange(t *testing.T) {
	executor := &GeneratorQueryExecutor{}
	start := time.Now()
	end := start.Add(time.Hour)

	// TimeRange should not modify the executor
	executor.TimeRange(start, end)
	assert.NotPanics(t, func() {
		executor.TimeRange(start, end)
	})
}
