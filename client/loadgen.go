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
	ErrNoCfg                 = errors.New("config is nil")
	ErrNoGun                 = errors.New("no gun implementation provided")
	ErrStaticRPS             = errors.New("static RPS must be > 0")
	ErrCallTimeout           = errors.New("generator request call timeout")
	ErrStartRPS              = errors.New("StartRPS must be > 0")
	ErrIncreaseRPS           = errors.New("IncreaseRPS must be > 0")
	ErrIncreaseAfterDuration = errors.New("IncreaseAfter must be > 1sec")
	ErrHoldRPS               = errors.New("HoldRPS must be > 0")
)

// LoadTestable is basic interface to run limited load with a contract call and save all transactions
type LoadTestable interface {
	Call(data interface{}) CallResult
}

// CallResult represents basic call result info
type CallResult struct {
	Duration time.Duration
	Data     interface{}
	Error    error
}

// LoadSchedule load test schedule
type LoadSchedule struct {
	StartRPS      int
	IncreaseRPS   int
	IncreaseAfter time.Duration
	HoldRPS       int
}

func (ls *LoadSchedule) Validate() error {
	if ls.StartRPS <= 0 {
		return ErrStartRPS
	}
	if ls.IncreaseRPS <= 0 {
		return ErrIncreaseRPS
	}
	if ls.HoldRPS <= 0 {
		return ErrHoldRPS
	}
	if ls.IncreaseAfter < 1 {
		return ErrIncreaseAfterDuration
	}
	return nil
}

// LoadGeneratorConfig is for shared load test data and configuration
type LoadGeneratorConfig struct {
	RPS                  int
	Schedule             *LoadSchedule
	Duration             time.Duration
	StatsPollInterval    time.Duration
	CallFailThreshold    int64
	CallTimeoutThreshold int64
	CallTimeout          time.Duration
	Gun                  LoadTestable
	SharedData           interface{}
}

func (lgc *LoadGeneratorConfig) Validate() error {
	if lgc.RPS == 0 {
		return ErrStaticRPS
	}
	if lgc.Duration == 0 {
		lgc.Duration = UntilStopDuration
	}
	if lgc.CallTimeout == 0 {
		lgc.CallTimeout = DefaultCallTimeout
	}
	if lgc.StatsPollInterval == 0 {
		lgc.StatsPollInterval = DefaultStatsPollInterval
	}
	if lgc.Gun == nil {
		return ErrNoGun
	}
	return nil
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
	okDataMu        *sync.Mutex
	OKData          []interface{}
	okResponsesMu   *sync.Mutex
	OKResponses     []CallResult
	failResponsesMu *sync.Mutex
	FailResponses   []CallResult
}

// LoadGenerator generates load on chain with some RPS
type LoadGenerator struct {
	cfg           *LoadGeneratorConfig
	rl            ratelimit.Limiter
	currentRPS    int
	schedule      *LoadSchedule
	wg            *sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	gun           LoadTestable
	responsesData *ResponseData
	errsMu        *sync.Mutex
	errs          []error
	stopped       atomic.Bool
	failed        atomic.Bool
	stats         *GeneratorStats
}

// NewLoadGenerator creates a new instance for a contract,
// shoots for scheduled RPS until timeout, test logic is defined through LoadTestable
func NewLoadGenerator(cfg *LoadGeneratorConfig) (*LoadGenerator, error) {
	if cfg == nil {
		return nil, ErrNoCfg
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.Schedule != nil {
		if err := cfg.Schedule.Validate(); err != nil {
			return nil, err
		}
	}
	rl := ratelimit.New(cfg.RPS)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration)
	return &LoadGenerator{
		cfg:      cfg,
		schedule: cfg.Schedule,
		rl:       rl,
		wg:       &sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
		gun:      cfg.Gun,
		responsesData: &ResponseData{
			okDataMu:        &sync.Mutex{},
			OKData:          make([]interface{}, 0),
			okResponsesMu:   &sync.Mutex{},
			OKResponses:     make([]CallResult, 0),
			failResponsesMu: &sync.Mutex{},
			FailResponses:   make([]CallResult, 0),
		},
		errsMu: &sync.Mutex{},
		errs:   make([]error, 0),
		stats:  &GeneratorStats{},
	}, nil
}

// runSchedule runs scheduling loop, changes LoadGenerator.currentRPS according to a load schedule
func (l *LoadGenerator) runSchedule() {
	if l.schedule == nil {
		return
	}
	l.rl = ratelimit.New(l.schedule.StartRPS)
	l.currentRPS = l.schedule.StartRPS
	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				time.Sleep(l.schedule.IncreaseAfter)
				newRPS := l.currentRPS + l.schedule.IncreaseRPS
				if newRPS > l.schedule.HoldRPS {
					log.Info().Int("RPS", l.currentRPS).Msg("Holding RPS")
					continue
				}
				l.rl = ratelimit.New(newRPS)
				l.currentRPS = newRPS
				log.Info().Int("RPS", l.currentRPS).Msg("Increasing RPS")
			}
		}
	}()
}

// pacedCall calls a gun according to a schedule or plain RPS
func (l *LoadGenerator) pacedCall() {
	l.rl.Take()
	if l.stopped.Load() {
		return
	}
	l.wg.Add(1)
	result := make(chan CallResult)
	ctx, cancel := context.WithTimeout(context.Background(), l.cfg.CallTimeout)
	callStartTS := time.Now()
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

			l.errsMu.Lock()
			defer l.errsMu.Unlock()
			l.errs = append(l.errs, ErrCallTimeout)
			l.responsesData.failResponsesMu.Lock()
			defer l.responsesData.failResponsesMu.Unlock()
			l.responsesData.FailResponses = append(l.responsesData.FailResponses, CallResult{Duration: time.Since(callStartTS), Error: ErrCallTimeout})
			log.Err(ctx.Err()).Msg("load generator transaction timeout")
		case res := <-result:
			defer close(result)
			res.Duration = time.Since(callStartTS)
			if res.Error != nil {
				l.stats.Failed.Add(1)

				l.errsMu.Lock()
				defer l.errsMu.Unlock()
				l.errs = append(l.errs, res.Error)
				l.responsesData.failResponsesMu.Lock()
				defer l.responsesData.failResponsesMu.Unlock()
				l.responsesData.FailResponses = append(l.responsesData.FailResponses, res)

				log.Err(res.Error).Msg("load generator request failed")
			} else {
				l.stats.Success.Add(1)
				l.responsesData.okDataMu.Lock()
				defer l.responsesData.okDataMu.Unlock()
				l.responsesData.OKData = append(l.responsesData.OKData, res.Data)
				l.responsesData.okResponsesMu.Lock()
				defer l.responsesData.okResponsesMu.Unlock()
				l.responsesData.OKResponses = append(l.responsesData.OKResponses, res)
			}
		}
		cancel()
		l.wg.Done()
	}()
}

// Run runs load loop until timeout or stop
func (l *LoadGenerator) Run() {
	log.Info().Msg("Load generator started")
	l.printStatsLoop()
	l.wg.Add(1)
	go l.runSchedule()
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
				l.pacedCall()
			}
		}
	}()
}

// Stop stops load generator, waiting for all calls for either finish or timeout
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

// Errors get all calls errors
func (l *LoadGenerator) Errors() []error {
	return l.errs
}

// GetData get all calls data
func (l *LoadGenerator) GetData() *ResponseData {
	return l.responsesData
}

// Stats get all load stats
func (l *LoadGenerator) Stats() *GeneratorStats {
	return l.stats
}

// PrintStats prints some runtime LoadGenerator.stats
func (l *LoadGenerator) PrintStats() {
	log.Info().
		Int64("Success", l.stats.Success.Load()).
		Int64("Failed", l.stats.Failed.Load()).
		Int64("CallTimeout", l.stats.CallTimeout.Load()).
		Msg("On-chain load stats")
}

// printStatsLoop prints stats periodically, with LoadGeneratorConfig.StatsPollInterval
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
