package client

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPositiveOneRequest(t *testing.T) {
	t.Parallel()
	gen := NewLoadGenerator(&LoadGeneratorConfig{
		RPS: 1,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	gen.Run()
	time.Sleep(40 * time.Millisecond)
	_, failed := gen.Stop()
	require.Equal(t, false, failed)
	gs := &GeneratorStats{}
	gs.Success.Add(2)
	require.Equal(t, gs, gen.Stats())

	okData, failData := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.Empty(t, failData)
	require.Empty(t, gen.Errors())
}

func TestFailedOneRequest(t *testing.T) {
	t.Parallel()
	gen := NewLoadGenerator(&LoadGeneratorConfig{
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

	okData, failData := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, []string{"failedCallData", "error", "failedCallData", "error"}, failData)
	require.Equal(t, []error{errors.New("error"), errors.New("error")}, gen.Errors())
}

func TestLoadGenCallTimeout(t *testing.T) {
	t.Parallel()
	gen := NewLoadGenerator(&LoadGeneratorConfig{
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

	okData, failData := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, []string{"generator request call timeout"}, failData)
	require.Equal(t, []error{ErrCallTimeout}, gen.Errors())
}

func TestLoadGenCallTimeoutWait(t *testing.T) {
	t.Parallel()
	gen := NewLoadGenerator(&LoadGeneratorConfig{
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

	okData, failData := convertResponsesData(gen.GetData())
	require.Empty(t, okData)
	require.Equal(t, []string{"generator request call timeout"}, failData)
	require.Equal(t, []error{ErrCallTimeout}, gen.Errors())
}

func TestCancelledByDeadlineWait(t *testing.T) {
	t.Parallel()
	gen := NewLoadGenerator(&LoadGeneratorConfig{
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
	okData, failData := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.Empty(t, failData)
	require.Empty(t, gen.Errors())
}

func TestCancelledBeforeDeadline(t *testing.T) {
	t.Parallel()
	gen := NewLoadGenerator(&LoadGeneratorConfig{
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

	okData, failData := convertResponsesData(gen.GetData())
	require.Equal(t, []string{"successCallData", "successCallData"}, okData)
	require.Empty(t, failData)
	require.Empty(t, gen.Errors())
}

func TestStaticRPSSchedulePrecision(t *testing.T) {
	gen := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:      1000,
		Duration: 1 * time.Second,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, false, failed)
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(999))
	require.LessOrEqual(t, gen.Stats().Success.Load(), int64(1001))
	require.Equal(t, gen.Stats().Failed.Load(), int64(0))
	require.Equal(t, gen.Stats().CallTimeout.Load(), int64(0))

	okData, failData := convertResponsesData(gen.GetData())
	require.GreaterOrEqual(t, len(okData), 999)
	require.LessOrEqual(t, len(okData), 1001)
	require.Empty(t, failData)
}

func TestStaticRPSScheduleIsNotBlocking(t *testing.T) {
	gen := NewLoadGenerator(&LoadGeneratorConfig{
		RPS:      1000,
		Duration: 1 * time.Second,
		Gun: NewMockGun(&MockGunConfig{
			CallSleep: 1 * time.Second,
		}),
	})
	gen.Run()
	_, failed := gen.Wait()
	require.Equal(t, false, failed)
	require.GreaterOrEqual(t, gen.Stats().Success.Load(), int64(999))
	require.LessOrEqual(t, gen.Stats().Success.Load(), int64(1001))
	require.Equal(t, gen.Stats().Failed.Load(), int64(0))
	require.Equal(t, gen.Stats().CallTimeout.Load(), int64(0))

	okData, failData := convertResponsesData(gen.GetData())
	require.GreaterOrEqual(t, len(okData), 999)
	require.LessOrEqual(t, len(okData), 1001)
	require.Empty(t, failData)
}
