package client

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPositiveOneRequest(t *testing.T) {
	t.Parallel()
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS: 1,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	gen.Run()
	time.Sleep(50 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
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
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS: 1,
		Gun: NewMockGun(&MockGunConfig{
			Fail:      true,
			CallSleep: 50 * time.Millisecond,
		}),
	})
	gen.Run()
	time.Sleep(40 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
	gs.Failed.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.GreaterOrEqual(t, failResponses[0].Duration, 50*time.Millisecond)
	require.GreaterOrEqual(t, failResponses[1].Duration, 50*time.Millisecond)
	require.Equal(t, failResponses[0].Data.(string), "failedCallData")
	require.Equal(t, failResponses[0].Error, errors.New("error"))
	require.Equal(t, failResponses[1].Data.(string), "failedCallData")
	require.Equal(t, failResponses[1].Error, errors.New("error"))
	require.Equal(t, []error{errors.New("error"), errors.New("error")}, gen.Errors())
}

func TestLoadGenCallTimeout(t *testing.T) {
	t.Parallel()
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:         1,
		CallTimeout: 50 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 55 * time.Millisecond,
		}),
	})
	gen.Run()
	time.Sleep(99 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, true, failed)
	gs := &GeneratorStats{}
	gs.Failed.Add(1)
	gs.CallTimeout.Add(1)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, failResponses[0].Data, nil)
	require.Equal(t, failResponses[0].Error, errors.New("generator request call timeout"))
	require.Equal(t, []error{ErrCallTimeout}, gen.Errors())
}

func TestLoadGenCallTimeoutWait(t *testing.T) {
	t.Parallel()
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:         1,
		CallTimeout: 50 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 55 * time.Millisecond,
		}),
	})
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, true, failed)
	gs := &GeneratorStats{}
	gs.Failed.Add(1)
	gs.CallTimeout.Add(1)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, failResponses[0].Data, nil)
	require.Equal(t, failResponses[0].Error, errors.New("generator request call timeout"))
	require.Equal(t, []error{ErrCallTimeout}, gen.Errors())
}

func TestCancelledByDeadlineWait(t *testing.T) {
	t.Parallel()
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:      1,
		Duration: 40 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
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
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:      1,
		Duration: 40 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	gen.Run()
	before := time.Now()
	time.Sleep(20 * time.Millisecond)
	_, failed := gen.Stop()
	after := time.Now()
	elapsed := after.Sub(before)

	require.Greater(t, elapsed, 1050*time.Millisecond)
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
	gs.Success.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, _, failResponses := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.Empty(t, failResponses)
	require.Empty(t, gen.Errors())
}

func TestStaticRPSSchedulePrecision(t *testing.T) {
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:      1000,
		Duration: 1 * time.Second,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
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

func TestStaticRPSScheduleIsNotBlocking(t *testing.T) {
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:      1000,
		Duration: 1 * time.Second,
		Gun: NewMockGun(&MockGunConfig{
			// call time must not affect the load schedule
			CallSleep: 1 * time.Second,
		}),
	})
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
	gen, _ := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:               1,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			StartRPS:      1,
			IncreaseRPS:   1,
			IncreaseAfter: 1 * time.Second,
			HoldRPS:       5,
		},
		Duration: 7000 * time.Millisecond,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, false, failed)
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(29))
}

func TestValidation(t *testing.T) {
	t.Parallel()
	_, err := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:               1,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			StartRPS:      0,
			IncreaseRPS:   1,
			IncreaseAfter: 1 * time.Second,
			HoldRPS:       5,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.Equal(t, ErrStartRPS, err)
	_, err = NewLoadGenerator(&LoadGeneratorConfig{
		RPS:               1,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			StartRPS:    1,
			IncreaseRPS: 1,
			HoldRPS:     5,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.Equal(t, ErrIncreaseAfterDuration, err)
	_, err = NewLoadGenerator(&LoadGeneratorConfig{
		RPS:               1,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			StartRPS:      1,
			IncreaseAfter: 1 * time.Second,
			HoldRPS:       5,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.Equal(t, ErrIncreaseRPS, err)
	_, err = NewLoadGenerator(&LoadGeneratorConfig{
		RPS:               1,
		StatsPollInterval: 1 * time.Second,
		Schedule: &LoadSchedule{
			StartRPS:      1,
			IncreaseRPS:   1,
			IncreaseAfter: 1 * time.Second,
		},
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.Equal(t, ErrHoldRPS, err)
	_, err = NewLoadGenerator(&LoadGeneratorConfig{
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 10 * time.Millisecond,
		}),
	})
	require.Equal(t, ErrStaticRPS, err)
	_, err = NewLoadGenerator(nil)
	require.Equal(t, ErrNoCfg, err)
	_, err = NewLoadGenerator(&LoadGeneratorConfig{
		RPS: 1,
		Gun: nil,
	})
	require.Equal(t, ErrNoGun, err)
}
