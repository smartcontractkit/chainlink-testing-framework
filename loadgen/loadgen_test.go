package loadgen

import (
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestConcurrentGenerators(t *testing.T) {
	t.Parallel()
	gens := make([]*Generator, 0)
	for i := 0; i < 2; i++ {
		gen, err := NewLoadGenerator(&LoadGeneratorConfig{
			T:        t,
			Duration: 1 * time.Second,
			Schedule: &LoadSchedule{
				Type:      RPSScheduleType,
				StartFrom: 1,
			},
			Gun: NewMockGun(&MockGunConfig{
				CallSleep: 50 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run()
		gens = append(gens, gen)
	}
	for _, gen := range gens {
		_, failed := gen.Wait()
		require.Equal(t, false, failed)
		stats := gen.Stats()
		require.Equal(t, int64(1), stats.CurrentRPS.Load())
		// we do not check exact RPS, because ratelimit.Limiter implementation
		// compensate RPS only after several requests
		// see example_test.go for precision tests of long runs
		// https://github.com/uber-go/ratelimit/blob/a12885fa6127db0aa3c29d33fc8ddeeb1fa1530c/limiter_atomic.go#L54
		require.GreaterOrEqual(t, stats.Success.Load(), int64(2))

		okData, okResponses, failResponses := convertResponsesData(gen.GetData())
		require.Contains(t, okData, "successCallData")
		require.GreaterOrEqual(t, len(okData), 2)
		require.GreaterOrEqual(t, okResponses[0].Duration, 50*time.Millisecond)
		require.Equal(t, okResponses[0].Data.(string), "successCallData")
		require.Equal(t, okResponses[1].Data.(string), "successCallData")
		require.GreaterOrEqual(t, len(okResponses), 2)
		require.Empty(t, failResponses)
		require.Empty(t, gen.Errors())
	}
}

func TestPositiveOneRequest(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	time.Sleep(50 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
	gs.CurrentInstances.Store(1)
	gs.CurrentRPS.Store(1)
	gs.Success.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, okResponses, failResponses := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.GreaterOrEqual(t, okResponses[0].Duration, 50*time.Millisecond)
	require.Equal(t, okResponses[0].Data.(string), "successCallData")
	require.Equal(t, okResponses[1].Data.(string), "successCallData")
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}

func TestFailedOneRequest(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		Gun: NewMockGun(&MockGunConfig{
			FailRatio: 100,
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	time.Sleep(40 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, true, failed)
	gs := &GeneratorStats{}
	gs.RunFailed.Store(true)
	gs.CurrentRPS.Store(1)
	gs.CurrentInstances.Store(1)
	gs.Failed.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.GreaterOrEqual(t, failResponses[0].Duration, 50*time.Millisecond)
	require.GreaterOrEqual(t, failResponses[1].Duration, 50*time.Millisecond)
	require.Equal(t, failResponses[0].Data.(string), "failedCallData")
	require.Equal(t, failResponses[0].Error, "error")
	require.Equal(t, failResponses[1].Data.(string), "failedCallData")
	require.Equal(t, failResponses[1].Error, "error")
	require.Equal(t, []string{"error", "error"}, gen.Errors())
}

func TestLoadGenCallTimeout(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		CallTimeout: 50 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 55 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	time.Sleep(99 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, true, failed)
	gs := &GeneratorStats{}
	gs.CurrentRPS.Store(1)
	gs.CurrentInstances.Store(1)
	gs.RunFailed.Store(true)
	gs.CallTimeout.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, failResponses[0].Data, nil)
	require.Equal(t, failResponses[0].Error, "generator request call timeout")
	require.Equal(t, []string{ErrCallTimeout.Error(), ErrCallTimeout.Error()}, gen.Errors())
}

func TestLoadGenCallTimeoutWait(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		CallTimeout: 50 * time.Millisecond,
		Duration:    1 * time.Second,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 55 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, true, failed)
	stats := gen.Stats()
	require.Equal(t, int64(1), stats.CurrentRPS.Load())
	require.GreaterOrEqual(t, stats.CallTimeout.Load(), int64(2))
	require.Equal(t, true, stats.RunFailed.Load())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, failResponses[0].Data, nil)
	require.Equal(t, failResponses[0].Error, "generator request call timeout")
	require.Contains(t, gen.Errors(), ErrCallTimeout.Error())
	require.GreaterOrEqual(t, len(gen.Errors()), 2)
}

func TestCancelledByDeadlineWait(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		Duration: 40 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	before := time.Now()
	_, failed := gen.Wait()
	after := time.Now()
	elapsed := after.Sub(before)
	// execution time + last request
	require.Greater(t, elapsed, 1050*time.Millisecond)
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
	gs.CurrentInstances.Store(1)
	gs.CurrentRPS.Store(1)
	gs.Success.Add(2)
	require.Equal(t, gs, gen.Stats())

	// in case of gen.Stop() if we don't have test duration or if gen.Wait() and we have a deadline
	// we are waiting for all requests, so result in that case must be successful
	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}

func TestCancelledBeforeDeadline(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		Duration: 40 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	before := time.Now()
	time.Sleep(20 * time.Millisecond)
	_, failed := gen.Stop()
	after := time.Now()
	elapsed := after.Sub(before)

	require.Greater(t, elapsed, 1050*time.Millisecond)
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
	gs.CurrentInstances.Store(1)
	gs.CurrentRPS.Store(1)
	gs.Success.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}

func TestStaticRPSSchedulePrecision(t *testing.T) {
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T:        t,
		Duration: 1 * time.Second,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1000,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, false, failed)
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(995))
	require.LessOrEqual(t, gen.Stats().Success.Load(), int64(1003))
	require.Equal(t, gen.Stats().Failed.Load(), int64(0))
	require.Equal(t, gen.Stats().CallTimeout.Load(), int64(0))

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.GreaterOrEqual(t, len(okData), 995)
	require.LessOrEqual(t, len(okData), 1002)
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}

func TestStaticRPSScheduleIsNotBlocking(t *testing.T) {
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T:        t,
		Duration: 1 * time.Second,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1000,
		},
		Gun: NewMockGun(&MockGunConfig{
			// call time must not affect the load schedule
			CallSleep: 1 * time.Second,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, false, failed)
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(998))
	require.LessOrEqual(t, gen.Stats().Success.Load(), int64(1002))
	require.Equal(t, gen.Stats().Failed.Load(), int64(0))
	require.Equal(t, gen.Stats().CallTimeout.Load(), int64(0))

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.GreaterOrEqual(t, len(okData), 998)
	require.LessOrEqual(t, len(okData), 1002)
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}

func TestLoadSchedule(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T:                 t,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			Type:          RPSScheduleType,
			StartFrom:     1,
			Increase:      1,
			StageInterval: 1 * time.Second,
			Limit:         5,
		},
		Duration: 7000 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, false, failed)
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(28))
}

func TestValidation(t *testing.T) {
	t.Parallel()
	_, err := NewLoadGenerator(&LoadGeneratorConfig{
		T:                 t,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			Type:          RPSScheduleType,
			StartFrom:     0,
			Increase:      1,
			StageInterval: 1 * time.Second,
			Limit:         5,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.Equal(t, ErrStartFrom, err)
	_, err = NewLoadGenerator(nil)
	require.Equal(t, ErrNoCfg, err)
	_, err = NewLoadGenerator(&LoadGeneratorConfig{
		T: t,
		Schedule: &LoadSchedule{
			Type:      RPSScheduleType,
			StartFrom: 1,
		},
		Gun: nil,
	})
	require.Equal(t, ErrNoImpl, err)
}

func TestInstancesRequests(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T:        t,
		Duration: 10 * time.Second,
		Schedule: &LoadSchedule{
			Type:          InstancesScheduleType,
			StartFrom:     1,
			Increase:      1,
			StageInterval: 1 * time.Second,
			Limit:         10,
		},
		Instance: NewMockInstance(&MockInstanceConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)
	gen.Run()
	_, failed := gen.Wait()
	stats := gen.Stats()
	require.Equal(t, false, failed)
	require.Equal(t, int64(10), stats.CurrentInstances.Load())

	okData, okResponses, failResponses := convertResponsesData(gen.GetData())
	require.GreaterOrEqual(t, okResponses[0].Duration, 50*time.Millisecond)
	require.Greater(t, len(okResponses), 1000)
	require.Greater(t, len(okData), 1000)
	require.Equal(t, okResponses[0].Data.(string), "successCallData")
	require.Equal(t, okResponses[1000].Data.(string), "successCallData")
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}
