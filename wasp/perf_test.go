package wasp

import (
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"

	"github.com/stretchr/testify/require"

	"context"
	//nolint
	_ "net/http/pprof"
	"runtime"

	"fmt"
	"net/http/httptest"
)

/* This tests can also be used as a performance validation of a tool itself or as a dashboard data filler */

var (
	labels = map[string]string{
		"branch": "generator_healthcheck",
		"commit": "generator_healthcheck",
	}
)

func stdPyro(t *testing.T) {
	t.Helper()
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "wasp",
		ServerAddress:   "http://localhost:4040",
		Logger:          pyroscope.StandardLogger,
		Tags:            map[string]string{"test": "wasp-trace-1"},

		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	require.NoError(t, err)
}

func TestPyroscopeLocalTraceRPSCalls(t *testing.T) {
	// run like
	// make test_pyro_rps or make test_pyro_vu
	// to have all in one, then
	// go tool trace trace.out
	stdPyro(t)
	t.Parallel()
	t.Run("trace test", func(t *testing.T) {
		t.Parallel()
		pyroscope.TagWrapper(context.Background(), pyroscope.Labels("scope", "loadgen_impl"), func(c context.Context) {
			gen, err := NewGenerator(&Config{
				T:           t,
				LokiConfig:  NewEnvLokiConfig(),
				Labels:      labels,
				CallTimeout: 100 * time.Millisecond,
				LoadType:    RPS,
				Schedule:    Plain(100, 10*time.Second),
				Gun: NewMockGun(&MockGunConfig{
					CallSleep: 50 * time.Millisecond,
				}),
			})
			require.NoError(t, err)
			//nolint
			gen.Run(true)
		})
	})
}

func TestPyroscopeLocalTraceVUCalls(t *testing.T) {
	stdPyro(t)
	t.Parallel()
	t.Run("trace test", func(t *testing.T) {
		t.Parallel()
		pyroscope.TagWrapper(context.Background(), pyroscope.Labels("scope", "loadgen_impl"), func(c context.Context) {
			gen, err := NewGenerator(&Config{
				T:           t,
				LokiConfig:  NewEnvLokiConfig(),
				Labels:      labels,
				CallTimeout: 100 * time.Millisecond,
				LoadType:    VU,
				Schedule:    Plain(10, 10*time.Second),
				VU: NewMockVU(&MockVirtualUserConfig{
					CallSleep: 50 * time.Millisecond,
				}),
			})
			require.NoError(t, err)
			//nolint
			gen.Run(true)
		})
	})
}

func TestPerfRenderLokiRPSRun(t *testing.T) {
	t.Parallel()
	t.Run("can_report_to_loki", func(t *testing.T) {
		t.Parallel()
		gen, err := NewGenerator(&Config{
			T:           t,
			LokiConfig:  NewEnvLokiConfig(),
			GenName:     "rps",
			Labels:      labels,
			CallTimeout: 100 * time.Millisecond,
			LoadType:    RPS,
			Schedule: CombineAndRepeat(
				2,
				Steps(10, 10, 10, 30*time.Second),
				Plain(200, 30*time.Second),
				Steps(100, -10, 10, 30*time.Second),
			),
			Gun: NewMockGun(&MockGunConfig{
				CallSleep: 50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run(true)
	})
}

func TestPerfRenderLokiVUsRun(t *testing.T) {
	t.Parallel()
	t.Run("can_report_to_loki", func(t *testing.T) {
		t.Parallel()
		gen, err := NewGenerator(&Config{
			T:           t,
			LokiConfig:  NewEnvLokiConfig(),
			GenName:     "vu",
			Labels:      labels,
			CallTimeout: 100 * time.Millisecond,
			LoadType:    VU,
			Schedule: CombineAndRepeat(
				2,
				Steps(10, 1, 10, 30*time.Second),
				Plain(30, 30*time.Second),
				Steps(20, -1, 10, 30*time.Second),
			),
			VU: NewMockVU(&MockVirtualUserConfig{
				CallSleep: 50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run(true)
	})
}

func TestRenderLokiParallelGenerators(t *testing.T) {
	t.Parallel()
	t.Run("can_report_to_loki", func(t *testing.T) {
		t.Parallel()
		p := NewProfile()
		for i := 0; i < 50; i++ {
			p.Add(NewGenerator(&Config{
				T:           t,
				LokiConfig:  NewEnvLokiConfig(),
				GenName:     fmt.Sprintf("rps-%d", i),
				Labels:      labels,
				CallTimeout: 100 * time.Millisecond,
				LoadType:    RPS,
				Schedule: Combine(
					Plain(10, 1*time.Minute),
				),
				Gun: NewMockGun(&MockGunConfig{
					CallSleep: 50 * time.Millisecond,
				}),
			}))
		}
		_, err := p.Run(true)
		require.NoError(t, err)
	})
}

func TestRenderLokiSpikeMaxLoadRun(t *testing.T) {
	t.Skip("This test is for manual run with or without Loki to measure max RPS")
	t.Parallel()
	t.Run("max_spike", func(t *testing.T) {
		t.Parallel()
		gen, err := NewGenerator(&Config{
			T:           t,
			LokiConfig:  NewEnvLokiConfig(),
			GenName:     "spike",
			Labels:      labels,
			CallTimeout: 100 * time.Millisecond,
			LoadType:    RPS,
			Schedule:    Plain(5000, 20*time.Second),
			Gun: NewMockGun(&MockGunConfig{
				CallSleep: 50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run(true)
	})
}

func TestRenderWS(t *testing.T) {
	t.Skip("This test is for manual run to measure max WS messages/s")
	s := httptest.NewServer(MockWSServer{
		Sleep: 50 * time.Millisecond,
		Logf:  t.Logf,
	})
	defer s.Close()

	gen, err := NewGenerator(&Config{
		T:          t,
		LokiConfig: NewEnvLokiConfig(),
		GenName:    "ws",
		Labels:     labels,
		LoadType:   VU,
		Schedule: []*Segment{
			{
				From:     10,
				Duration: 10 * time.Second,
			},
		},
		VU: NewWSMockVU(&WSMockVUConfig{TargetURl: s.URL}),
	})
	require.NoError(t, err)
	gen.Run(true)
}

func TestRenderHTTP(t *testing.T) {
	t.Skip("This test is for manual run to measure max HTTP RPS")
	srv := NewHTTPMockServer(nil)
	srv.Run()

	gen, err := NewGenerator(&Config{
		T:          t,
		LokiConfig: NewEnvLokiConfig(),
		GenName:    "http",
		Labels:     labels,
		LoadType:   RPS,
		Schedule:   Steps(10, 10, 10, 500*time.Second),
		Gun:        NewHTTPMockGun(&MockHTTPGunConfig{TargetURL: "http://localhost:8080"}),
	})
	require.NoError(t, err)
	gen.Run(true)
}
