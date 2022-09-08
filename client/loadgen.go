package client

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
)

const (
	DefaultCallTimeout       = 1 * time.Minute
	DefaultStatsPollInterval = 10 * time.Second
	UntilStopDuration        = 99999 * time.Hour
)

var (
	ErrCallTimeout = errors.New("generator request call timeout")
)

// LoadTestable is basic interface to run limited load with a contract call and save all transactions
type LoadTestable interface {
	Call(data interface{}) CallResult
}

type CallResult struct {
	Data  interface{}
	Error error
}

type VerifyResult struct {
	Data  interface{}
	Error error
}

// LoadGeneratorConfig is for shared load test data and configuration
type LoadGeneratorConfig struct {
	RPS                  int
	Duration             time.Duration
	StatsPollInterval    time.Duration
	CallFailThreshold    int64
	CallTimeoutThreshold int64
	CallTimeout          time.Duration
	Gun                  LoadTestable
	SharedData           interface{}
}

// GeneratorStats basic generator load stats
type GeneratorStats struct {
	Success     atomic.Int64
	Failed      atomic.Int64
	CallTimeout atomic.Int64
}

// ResponseData includes any request/response data that a gun might store
// ok* slices usually contains successful responses and their verifications if their done async
// fail* slices contains CallResult with response data and an error
type ResponseData struct {
	okDataMu   *sync.Mutex
	OKData     []interface{}
	failDataMu *sync.Mutex
	FailData   []CallResult
}

// LoadGenerator generates load on chain with some RPS
type LoadGenerator struct {
	cfg           *LoadGeneratorConfig
	rl            ratelimit.Limiter
	wg            *sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	gun           LoadTestable
	responsesData *ResponseData
	errs          []error
	stopped       atomic.Bool
	failed        atomic.Bool
	stats         *GeneratorStats
}

// NewLoadGenerator creates a new instance for a contract,
// shoots for scheduled RPS until timeout, test logic is defined through LoadTestable
func NewLoadGenerator(cfg *LoadGeneratorConfig) *LoadGenerator {
	rl := ratelimit.New(cfg.RPS)
	if cfg.Duration == 0 {
		cfg.Duration = UntilStopDuration
	}
	if cfg.CallTimeout == 0 {
		cfg.CallTimeout = DefaultCallTimeout
	}
	if cfg.StatsPollInterval == 0 {
		cfg.StatsPollInterval = DefaultStatsPollInterval
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration)
	return &LoadGenerator{
		cfg:    cfg,
		rl:     rl,
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
		gun:    cfg.Gun,
		responsesData: &ResponseData{
			okDataMu:   &sync.Mutex{},
			OKData:     make([]interface{}, 0),
			failDataMu: &sync.Mutex{},
			FailData:   make([]CallResult, 0),
		},
		errs:  make([]error, 0),
		stats: &GeneratorStats{},
	}
}

func (l *LoadGenerator) Errors() []error {
	return l.errs
}

func (l *LoadGenerator) GetData() *ResponseData {
	return l.responsesData
}

func (l *LoadGenerator) Stats() *GeneratorStats {
	return l.stats
}

// PrintStats prints some runtime stats
func (l *LoadGenerator) PrintStats() {
	log.Info().
		Int64("Success", l.stats.Success.Load()).
		Int64("Failed", l.stats.Failed.Load()).
		Int64("CallTimeout", l.stats.CallTimeout.Load()).
		Msg("On-chain load stats")
}

func (l *LoadGenerator) printStatsLoop() {
	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				time.Sleep(l.cfg.StatsPollInterval)
				log.Info().
					Int64("Success", l.stats.Success.Load()).
					Int64("Failed", l.stats.Failed.Load()).
					Int64("CallTimeout", l.stats.CallTimeout.Load()).
					Msg("Load stats")
			}
		}
	}()
}

// scheduledCall schedules a rate limited contract call
func (l *LoadGenerator) scheduledCall() {
	l.rl.Take()
	if l.stopped.Load() {
		return
	}
	l.wg.Add(1)
	result := make(chan CallResult)
	ctx, cancel := context.WithTimeout(context.Background(), l.cfg.CallTimeout)
	go func() {
		result <- l.gun.Call(l.cfg.SharedData)
	}()
	go func() {
		select {
		case <-ctx.Done():
			l.stopped.Store(true)
			l.failed.Store(true)
			l.stats.CallTimeout.Add(1)
			l.stats.Failed.Add(1)

			l.errs = append(l.errs, ErrCallTimeout)
			l.responsesData.failDataMu.Lock()
			defer l.responsesData.failDataMu.Unlock()
			l.responsesData.FailData = append(l.responsesData.FailData, CallResult{Error: ErrCallTimeout})
			log.Err(ctx.Err()).Msg("load generator transaction timeout")
		case res := <-result:
			defer close(result)
			if res.Error != nil {
				l.stats.Failed.Add(1)

				l.errs = append(l.errs, res.Error)
				l.responsesData.failDataMu.Lock()
				defer l.responsesData.failDataMu.Unlock()
				l.responsesData.FailData = append(l.responsesData.FailData, res)

				log.Err(res.Error).Msg("load generator request failed")
			} else {
				l.stats.Success.Add(1)
				l.responsesData.okDataMu.Lock()
				defer l.responsesData.okDataMu.Unlock()
				l.responsesData.OKData = append(l.responsesData.OKData, res.Data)
			}
		}
		cancel()
		l.wg.Done()
	}()
}

// Run runs load until timeout
func (l *LoadGenerator) Run() {
	log.Info().Msg("Load generator started")
	l.printStatsLoop()
	l.wg.Add(1)
	go func() {
		for {
			select {
			case <-l.ctx.Done():
				log.Info().Msg("Load generator stopped, waiting for requests to finish")
				l.wg.Done()
				l.wg.Wait()
				log.Info().Msg("Load generator exited")
				l.PrintStats()
				return
			default:
				if l.stats.Failed.Load() > l.cfg.CallFailThreshold || l.stats.CallTimeout.Load() > l.cfg.CallTimeoutThreshold {
					l.cancel()
					l.failed.Store(true)
					log.Info().Msg("Test reached failed requests threshold")
				}
				l.scheduledCall()
			}
		}
	}()
}

// Stop stops load generator, waiting for all txns for either finish or timeout
func (l *LoadGenerator) Stop() (interface{}, bool) {
	l.cancel()
	l.wg.Wait()
	return l.GetData(), l.failed.Load()
}

// Wait waits until test ends
func (l *LoadGenerator) Wait() (interface{}, bool) {
	l.wg.Wait()
	return l.GetData(), l.failed.Load()
}
