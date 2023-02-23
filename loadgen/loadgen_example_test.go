package loadgen

import (
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/stretchr/testify/require"
)

/* This tests can also be used as a performance validation of a tool itself */

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
				StartFrom: 5000,
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
