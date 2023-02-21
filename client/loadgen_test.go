package client

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

func TestLokiReporting(t *testing.T) {
	t.Skip("This test is for manual run and dashboard development")
	t.Parallel()
	t.Run("can_report_to_loki", func(t *testing.T) {
		t.Parallel()
		gen, err := NewLoadGenerator(&LoadGeneratorConfig{
			T:          t,
			LokiConfig: NewDefaultLokiConfig("http://localhost:3030/loki/api/v1/push"),
			Labels: map[string]string{
				"cluster":    "sdlc",
				"app":        "chainlink",
				"env":        "chainlink-test",
				"test_group": "stress",
				"test_id":    "zxc-11",
			},
			CallTimeout: 300 * time.Millisecond,
			Duration:    60 * time.Second,
			Schedule: &LoadSchedule{
				Type:          RPSScheduleType,
				StartFrom:     10,
				Increase:      1,
				StageInterval: 1 * time.Second,
				Limit:         60,
			},
			Gun: NewMockGun(&MockGunConfig{
				TimeoutRatio: 15,
				CallSleep:    130 * time.Millisecond,
			}),
		})
		require.NoError(t, err)
		gen.Run()
		_, _ = gen.Wait()
	})
}

func TestInstancesRequests(t *testing.T) {
	t.Parallel()
	gen, err := NewLoadGenerator(&LoadGeneratorConfig{
		T:        t,
		Duration: 10 * time.Second,
		Schedule: &LoadSchedule{
			Type:          InstancesScheduleType,
			StartFrom:     1,
			Increase:      2,
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
	require.Greater(t, len(okResponses), 3720)
	require.Greater(t, len(okData), 3720)
	require.Equal(t, okResponses[0].Data.(string), "successCallData")
	require.Equal(t, okResponses[3720].Data.(string), "successCallData")
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
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
	// because of go.uber.org/ratelimit implementation, if RPS = 1 it waits for one tick before start,
	// so 1 sec to start the schedule + 50ms because we are waiting for request to finish after the test is finished
	// it also fires 2 requests from the beginning to compensate that, so RPS schedule is accurate +-1 request
	// see TestStaticRPSSchedulePrecision
	require.Greater(t, elapsed, 1050*time.Millisecond)
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
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
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(998))
	require.LessOrEqual(t, gen.Stats().Success.Load(), int64(1003))
	require.Equal(t, gen.Stats().Failed.Load(), int64(0))
	require.Equal(t, gen.Stats().CallTimeout.Load(), int64(0))

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.GreaterOrEqual(t, len(okData), 998)
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
	//_, err = NewLoadGenerator(&LoadGeneratorConfig{
	//	T:                 t,
	//	RPS:               1,
	//	StatsPollInterval: 1 * time.Second,
	//	Schedule: &LoadSchedule{
	//		Type:      RPSScheduleType,
	//		StartFrom: 1,
	//		Increase:  1,
	//		Limit:     5,
	//	},
	//	Gun: NewMockGun(&MockGunConfig{
	//		CallSleep: 10 * time.Millisecond,
	//	}),
	//})
	//require.Equal(t, ErrStageInterval, err)
	//_, err = NewLoadGenerator(&LoadGeneratorConfig{
	//	T:                 t,
	//	RPS:               1,
	//	StatsPollInterval: 1 * time.Second,
	//	Schedule: &LoadSchedule{
	//		Type:          RPSScheduleType,
	//		StartFrom:     1,
	//		StageInterval: 1 * time.Second,
	//		Limit:         5,
	//	},
	//	Gun: NewMockGun(&MockGunConfig{
	//		CallSleep: 10 * time.Millisecond,
	//	}),
	//})
	//require.Equal(t, ErrIncrease, err)
	//_, err = NewLoadGenerator(&LoadGeneratorConfig{
	//	T:                 t,
	//	RPS:               1,
	//	StatsPollInterval: 1 * time.Second,
	//	Schedule: &LoadSchedule{
	//		Type:          RPSScheduleType,
	//		StartFrom:     1,
	//		Increase:      1,
	//		StageInterval: 1 * time.Second,
	//	},
	//	Gun: NewMockGun(&MockGunConfig{
	//		CallSleep: 10 * time.Millisecond,
	//	}),
	//})
	//require.Equal(t, ErrScheduleLimit, err)
	//_, err = NewLoadGenerator(&LoadGeneratorConfig{
	//	T: t,
	//	Gun: NewMockGun(&MockGunConfig{
	//		CallSleep: 10 * time.Millisecond,
	//	}),
	//})
	//require.Equal(t, ErrStaticRPS, err)
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
