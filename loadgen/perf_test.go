package loadgen

import (
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
			gen, err := NewLoadGenerator(&LoadGeneratorConfig{
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
				Duration:    10 * time.Second,
				Schedule: &LoadSchedule{
					Type:      RPSScheduleType,
					StartFrom: 2000,
				},
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
		gen, err := NewLoadGenerator(&LoadGeneratorConfig{
			T: t,
			LokiConfig: client.NewDefaultLokiConfig(
				os.Getenv("LOKI_URL"),
				os.Getenv("LOKI_TOKEN")),
			Labels: map[string]string{
				"cluster":    "sdlc",
				"namespace":  "load-dummy-test",
				"app":        "dummy",
				"test_group": "generator_healthcheck",
				"test_id":    "dummy-healthcheck-rps-1",
			},
			CallTimeout: 100 * time.Millisecond,
			Duration:    10 * time.Second,
			Schedule: &LoadSchedule{
				Type:      RPSScheduleType,
				StartFrom: 1000,
			},
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
		gen, err := NewLoadGenerator(&LoadGeneratorConfig{
			T: t,
			LokiConfig: client.NewDefaultLokiConfig(
				os.Getenv("LOKI_URL"),
				os.Getenv("LOKI_TOKEN")),
			Labels: map[string]string{
				"cluster":    "sdlc",
				"namespace":  "load-dummy-test",
				"app":        "dummy",
				"test_group": "generator_healthcheck",
				"test_id":    "dummy-healthcheck-instances-1",
			},
			CallTimeout: 100 * time.Millisecond,
			Duration:    30 * time.Second,
			Schedule: &LoadSchedule{
				Type:          InstancesScheduleType,
				StartFrom:     1,
				Increase:      3,
				StageInterval: 10 * time.Second,
				Limit:         30,
			},
			Instance: NewMockInstance(&MockInstanceConfig{
				FailRatio: 5,
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
		gen, err := NewLoadGenerator(&LoadGeneratorConfig{
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
			Duration:    20 * time.Second,
			Schedule: &LoadSchedule{
				Type: RPSScheduleType,
				// TODO: tune Loki for 5k+ RPS responses
				StartFrom: 5000,
			},
			Gun: NewMockGun(&MockGunConfig{
				//TimeoutRatio: 1,
				CallSleep: 50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run()
		_, _ = gen.Wait()
	})
}
