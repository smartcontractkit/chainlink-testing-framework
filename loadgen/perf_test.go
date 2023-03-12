package loadgen

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/stretchr/testify/require"

	"context"
	//nolint
	_ "net/http/pprof"
	"runtime"

	"github.com/pyroscope-io/client/pyroscope"
)

/* This tests can also be used as a performance validation of a tool itself or as a dashboard data filler */

func stdPyro(t *testing.T) {
	t.Helper()
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "loadgen",
		ServerAddress:   "http://localhost:4040",
		Logger:          pyroscope.StandardLogger,
		Tags:            map[string]string{"test": "loadgen-trace-1"},

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

func TestLocalTrace(t *testing.T) {
	t.Skip("Local tracing test")
	// run like
	// go test -run TestLocalTrace -trace trace.out
	// to have all in one, then
	// go tool trace trace.out
	stdPyro(t)
	t.Parallel()
	t.Run("trace test", func(t *testing.T) {
		t.Parallel()
		pyroscope.TagWrapper(context.Background(), pyroscope.Labels("scope", "loadgen_impl"), func(c context.Context) {
			gen, err := NewLoadGenerator(&Config{
				T: t,
				LokiConfig: client.NewDefaultLokiConfig(
					os.Getenv("LOKI_URL"),
					os.Getenv("LOKI_TOKEN")),
				Labels: map[string]string{
					"cluster":    "sdlc",
					"namespace":  "load-dummy-test",
					"app":        "dummy",
					"test_group": "generator_healthcheck",
					"test_id":    "dummy-healthcheck-pyro-1",
				},
				CallTimeout: 100 * time.Millisecond,
				LoadType:    RPSScheduleType,
				Schedule:    Plain(2000, 10*time.Second),
				Gun: NewMockGun(&MockGunConfig{
					TimeoutRatio: 1,
					CallSleep:    50 * time.Millisecond,
				}),
			})
			require.NoError(t, err)
			//nolint
			gen.Run()
			_, _ = gen.Wait()
		})
	})
}

func TestLokiRPSRun(t *testing.T) {
	t.Skip("This test is for manual run and dashboard development, you need LOKI_URL and LOKI_TOKEN to run")
	t.Parallel()
	t.Run("can_report_to_loki", func(t *testing.T) {
		t.Parallel()
		gen, err := NewLoadGenerator(&Config{
			T: t,
			LokiConfig: client.NewDefaultLokiConfig(
				os.Getenv("LOKI_URL"),
				os.Getenv("LOKI_TOKEN")),
			Labels: map[string]string{
				"test_group": "generator_healthcheck",
				"cluster":    "sdlc",
				"app":        "dummy",
				"namespace":  "load-dummy-test",
				"test_id":    "dummy-healthcheck-rps-1",
			},
			CallTimeout: 100 * time.Millisecond,
			LoadType:    RPSScheduleType,
			Schedule: CombineAndRepeat(
				2,
				Line(1, 100, 30*time.Second),
				Plain(200, 30*time.Second),
				Line(100, 1, 30*time.Second),
			),
			Gun: NewMockGun(&MockGunConfig{
				TimeoutRatio: 1,
				CallSleep:    50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run()
		_, _ = gen.Wait()
	})
}

func TestLokiInstancesRun(t *testing.T) {
	t.Skip("This test is for manual run and dashboard development, you need LOKI_URL and LOKI_TOKEN to run")
	t.Parallel()
	t.Run("can_report_to_loki", func(t *testing.T) {
		t.Parallel()
		gen, err := NewLoadGenerator(&Config{
			T: t,
			LokiConfig: client.NewDefaultLokiConfig(
				os.Getenv("LOKI_URL"),
				os.Getenv("LOKI_TOKEN")),
			Labels: map[string]string{
				"test_group": "generator_healthcheck",
				"cluster":    "sdlc",
				"app":        "dummy",
				"namespace":  "load-dummy-test",
				"test_id":    "dummy-healthcheck-instances-1",
			},
			CallTimeout: 100 * time.Millisecond,
			LoadType:    InstancesScheduleType,
			Schedule: CombineAndRepeat(
				2,
				Line(1, 20, 30*time.Second),
				Plain(30, 30*time.Second),
				Line(20, 1, 30*time.Second),
			),
			Instance: NewMockInstance(MockInstanceConfig{
				CallSleep: 100 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run()
		_, _ = gen.Wait()
	})
}

func TestLokiSpikeMaxLoadRun(t *testing.T) {
	t.Skip("This test is for manual run with or without Loki to measure max RPS")
	t.Parallel()
	t.Run("max_spike", func(t *testing.T) {
		t.Parallel()
		gen, err := NewLoadGenerator(&Config{
			T: t,
			LokiConfig: client.NewDefaultLokiConfig(
				os.Getenv("LOKI_URL"),
				os.Getenv("LOKI_TOKEN")),
			Labels: map[string]string{
				"cluster":    "sdlc",
				"namespace":  "load-dummy-test",
				"app":        "dummy",
				"test_group": "generator_healthcheck",
				"test_id":    "dummy-healthcheck-max-1",
			},
			CallTimeout: 100 * time.Millisecond,
			LoadType:    RPSScheduleType,
			Schedule:    Plain(5000, 20*time.Second),
			Gun: NewMockGun(&MockGunConfig{
				CallSleep: 50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run()
		_, _ = gen.Wait()
	})
}

func TestWS(t *testing.T) {
	t.Skip("This test is for manual run to measure max WS messages/s")
	s := httptest.NewServer(MockWSServer{
		sleep: 50 * time.Millisecond,
		logf:  t.Logf,
	})
	defer s.Close()

	gen, err := NewLoadGenerator(&Config{
		T: t,
		LokiConfig: client.NewDefaultLokiConfig(
			os.Getenv("LOKI_URL"),
			os.Getenv("LOKI_TOKEN")),
		Labels: map[string]string{
			"cluster":    "sdlc",
			"namespace":  "ws-dummy-test",
			"app":        "dummy",
			"test_group": "generator_healthcheck",
			"test_id":    "dummy-healthcheck-ws-instances-stages-25ms-answer",
		},
		LoadType: InstancesScheduleType,
		Schedule: []*Segment{
			{
				From:         10,
				Increase:     20,
				Steps:        10,
				StepDuration: 10 * time.Second,
			},
		},
		Instance: NewWSMockInstance(WSMockConfig{TargetURl: s.URL}),
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}

func TestHTTP(t *testing.T) {
	t.Skip("This test is for manual run to measure max HTTP RPS")
	srv := NewHTTPMockServer(50 * time.Millisecond)
	srv.Run()

	gen, err := NewLoadGenerator(&Config{
		T: t,
		LokiConfig: client.NewDefaultLokiConfig(
			os.Getenv("LOKI_URL"),
			os.Getenv("LOKI_TOKEN")),
		Labels: map[string]string{
			"cluster":    "sdlc",
			"namespace":  "dummy",
			"app":        "dummy",
			"test_group": "dummy",
			"test_id":    "http",
		},
		LoadType: RPSScheduleType,
		Schedule: Line(10, 400, 500*time.Second),
		Gun:      NewHTTPMockGun(&MockHTTPGunConfig{TargetURL: "http://localhost:8080"}),
	})
	require.NoError(t, err)
	gen.Run()
	_, _ = gen.Wait()
}
