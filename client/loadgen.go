package client

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
)

const (
	DefaultCallTimeout       = 1 * time.Minute
	DefaultStatsPollInterval = 5 * time.Second
	UntilStopDuration        = 99999 * time.Hour
)

var (
	ErrNoCfg        = errors.New("config is nil")
	ErrNoImpl       = errors.New("either \"gun\" or \"instance\" implementation must provided")
	ErrCallTimeout  = errors.New("generator request call timeout")
	ErrStartFrom    = errors.New("StartFrom must be > 0")
	ErrScheduleType = errors.New("schedule type must be \"rps_schedule\" or \"instance_schedule\"")
)

// LoadTestable is basic interface to run limited load with a contract call and save all transactions
type LoadTestable interface {
	Call(l *LoadGenerator) CallResult
}

// LoadInstance is basic interface to run load instances
type LoadInstance interface {
	Run(l *LoadGenerator)
}

// CallResult represents basic call result info
type CallResult struct {
	Failed     bool          `json:"failed"`
	Timeout    bool          `json:"timeout"`
	Duration   time.Duration `json:"duration"`
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Data       interface{}   `json:"data"`
	Error      string        `json:"error,omitempty"`
}

const (
	RPSScheduleType       string = "rps_schedule"
	InstancesScheduleType string = "instance_schedule"
)

// LoadSchedule load test schedule
type LoadSchedule struct {
	Type          string
	StartFrom     int64
	Increase      int64
	StageInterval time.Duration
	Limit         int64
}

func (ls *LoadSchedule) Validate() error {
	if ls.Type != RPSScheduleType && ls.Type != InstancesScheduleType {
		return ErrScheduleType
	}
	if ls.StartFrom <= 0 {
		return ErrStartFrom
	}
	return nil
}

// LoadGeneratorConfig is for shared load test data and configuration
type LoadGeneratorConfig struct {
	T                 *testing.T
	Labels            map[string]string
	LokiConfig        *LokiConfig
	Schedule          *LoadSchedule
	Duration          time.Duration
	StatsPollInterval time.Duration
	CallTimeout       time.Duration
	Gun               LoadTestable
	Instance          LoadInstance
	Logger            zerolog.Logger
	SharedData        interface{}
}

func (lgc *LoadGeneratorConfig) Validate() error {
	if lgc.Duration == 0 {
		lgc.Duration = UntilStopDuration
	}
	if lgc.CallTimeout == 0 {
		lgc.CallTimeout = DefaultCallTimeout
	}
	if lgc.StatsPollInterval == 0 {
		lgc.StatsPollInterval = DefaultStatsPollInterval
	}
	if lgc.Gun == nil && lgc.Instance == nil {
		return ErrNoImpl
	}
	return nil
}

// GeneratorStats basic generator load stats
type GeneratorStats struct {
	CurrentRPS       atomic.Int64 `json:"currentRPS"`
	CurrentInstances atomic.Int64 `json:"currentInstances"`
	RunStopped       atomic.Bool  `json:"runStopped"`
	RunFailed        atomic.Bool  `json:"runFailed"`
	Success          atomic.Int64 `json:"success"`
	Failed           atomic.Int64 `json:"failed"`
	CallTimeout      atomic.Int64 `json:"callTimeout"`
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

// LoadGenerator generates load with some RPS
type LoadGenerator struct {
	cfg                  *LoadGeneratorConfig
	log                  zerolog.Logger
	labels               model.LabelSet
	rl                   ratelimit.Limiter
	schedule             *LoadSchedule
	responsesWaitGroup   *sync.WaitGroup
	ResponsesCtx         context.Context
	responsesCancel      context.CancelFunc
	dataCtx              context.Context
	dataCancel           context.CancelFunc
	gun                  LoadTestable
	instance             LoadInstance
	instanceResponseChan chan CallResult
	responsesData        *ResponseData
	errsMu               *sync.Mutex
	errs                 []string
	stats                *GeneratorStats
	loki                 *LokiClient
	lokiResponsesChan    chan CallResult
}

// NewLoadGenerator creates a new instance for a contract,
// shoots for scheduled RPS until timeout, test logic is defined through LoadTestable
func NewLoadGenerator(cfg *LoadGeneratorConfig) (*LoadGenerator, error) {
	if cfg == nil {
		return nil, ErrNoCfg
	}
	if err := cfg.Schedule.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	rl := ratelimit.New(int(cfg.Schedule.StartFrom))

	var loki *LokiClient
	var ls model.LabelSet
	var err error
	if cfg.LokiConfig != nil {
		loki, err = NewLokiClient(cfg.LokiConfig)
		if err != nil {
			return nil, err
		}
		ls = model.LabelSet{}
		for k, v := range cfg.Labels {
			ls[model.LabelName(k)] = model.LabelValue(v)
		}
	}
	// context for all requests/responses and instances
	responsesCtx, responsesCancel := context.WithTimeout(context.Background(), cfg.Duration)
	// context for all the collected data
	dataCtx, dataCancel := context.WithCancel(context.Background())
	// creating logger from *testing.T context or using a global logger
	var l zerolog.Logger
	if cfg.T != nil {
		l = zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(cfg.T))).With().Timestamp().Logger()
	} else {
		l = log.Logger
	}
	return &LoadGenerator{
		cfg:                  cfg,
		schedule:             cfg.Schedule,
		rl:                   rl,
		responsesWaitGroup:   &sync.WaitGroup{},
		ResponsesCtx:         responsesCtx,
		responsesCancel:      responsesCancel,
		dataCtx:              dataCtx,
		dataCancel:           dataCancel,
		gun:                  cfg.Gun,
		instance:             cfg.Instance,
		instanceResponseChan: make(chan CallResult),
		labels:               ls,
		responsesData: &ResponseData{
			okDataMu:        &sync.Mutex{},
			OKData:          make([]interface{}, 0),
			okResponsesMu:   &sync.Mutex{},
			OKResponses:     make([]CallResult, 0),
			failResponsesMu: &sync.Mutex{},
			FailResponses:   make([]CallResult, 0),
		},
		errsMu:            &sync.Mutex{},
		errs:              make([]string, 0),
		stats:             &GeneratorStats{},
		loki:              loki,
		log:               l,
		lokiResponsesChan: make(chan CallResult, 50000),
	}, nil
}

func (l *LoadGenerator) handleLokiResponsePayload(cr CallResult) {
	ls := l.labels.Merge(model.LabelSet{
		"go_test_name":   model.LabelValue(l.cfg.T.Name()),
		"test_data_type": "responses",
	})
	// generates new labels, Duration is already calculated at that point
	cr.StartedAt = time.Time{}
	cr.FinishedAt = time.Time{}
	err := l.loki.HandleStruct(ls, time.Now(), cr)
	if err != nil {
		l.log.Err(err).Send()
	}
}

func (l *LoadGenerator) handleLokiStatsPayload() {
	ls := l.labels.Merge(model.LabelSet{
		"go_test_name":   model.LabelValue(l.cfg.T.Name()),
		"test_data_type": "stats",
	})
	err := l.loki.HandleStruct(ls, time.Now(), l.StatsJSON())
	if err != nil {
		l.log.Err(err).Send()
	}
}

func (l *LoadGenerator) runLokiPromtailResponses() {
	l.log.Info().
		Str("URL", l.cfg.LokiConfig.URL).
		Interface("DefaultLabels", l.cfg.Labels).
		Msg("Streaming data to Loki")
	go func() {
		for {
			select {
			case <-l.dataCtx.Done():
				l.log.Info().Msg("Loki responses exited")
				return
			case r := <-l.lokiResponsesChan:
				l.handleLokiResponsePayload(r)
			}
		}
	}()
}

func (l *LoadGenerator) runLokiPromtailStats() {
	go func() {
		for {
			select {
			case <-l.dataCtx.Done():
				l.log.Info().Msg("Loki stats exited")
				return
			default:
				time.Sleep(l.cfg.StatsPollInterval)
				l.handleLokiStatsPayload()
			}
		}
	}()
}

// runSchedule runs scheduling loop, changes LoadGenerator.currentRPS according to a load schedule
func (l *LoadGenerator) runSchedule() {
	l.rl = ratelimit.New(int(l.schedule.StartFrom))
	l.stats.CurrentRPS.Store(l.schedule.StartFrom)
	l.responsesWaitGroup.Add(1)
	go func() {
		defer l.responsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				l.log.Info().Msg("Scheduler exited")
				return
			default:
				time.Sleep(l.schedule.StageInterval)
				switch l.cfg.Schedule.Type {
				case RPSScheduleType:
					newRPS := l.stats.CurrentRPS.Load() + l.schedule.Increase
					if newRPS > l.schedule.Limit {
						return
					}
					l.rl = ratelimit.New(int(newRPS))
					l.stats.CurrentRPS.Store(newRPS)
				case InstancesScheduleType:
					newInstances := l.stats.CurrentInstances.Load() + l.schedule.Increase
					if newInstances > l.schedule.Limit {
						return
					}
					l.stats.CurrentInstances.Store(newInstances)
					for i := 0; i < int(newInstances); i++ {
						l.responsesWaitGroup.Add(1)
						l.instance.Run(l)
					}
				}
			}
		}
	}()
}

func (l *LoadGenerator) handleCallResult(res CallResult) {
	if l.cfg.LokiConfig != nil {
		l.lokiResponsesChan <- res
	}
	if res.Error != "" {
		l.stats.RunFailed.Store(true)
		l.stats.Failed.Add(1)

		l.errsMu.Lock()
		l.responsesData.failResponsesMu.Lock()
		l.errs = append(l.errs, res.Error)
		l.responsesData.FailResponses = append(l.responsesData.FailResponses, res)
		l.errsMu.Unlock()
		l.responsesData.failResponsesMu.Unlock()

		l.log.Error().Str("Err", res.Error).Msg("load generator request failed")
	} else {
		l.stats.Success.Add(1)
		l.responsesData.okDataMu.Lock()
		l.responsesData.OKData = append(l.responsesData.OKData, res.Data)
		l.responsesData.okResponsesMu.Lock()
		l.responsesData.OKResponses = append(l.responsesData.OKResponses, res)
		l.responsesData.okDataMu.Unlock()
		l.responsesData.okResponsesMu.Unlock()
	}
}

func (l *LoadGenerator) collectData() {
	if l.cfg.Schedule.Type == RPSScheduleType {
		return
	}
	go func() {
		for {
			select {
			case <-l.dataCtx.Done():
				l.log.Info().Msg("Collect data exited")
				return
			case res := <-l.instanceResponseChan:
				if res.StartedAt.IsZero() {
					log.Error().Msg("StartedAt is not set in instance implementation")
					return
				}
				res.FinishedAt = time.Now()
				res.Duration = time.Since(res.StartedAt)
				l.handleCallResult(res)
			}
		}
	}()
}

// pacedCall calls a gun according to a schedule or plain RPS
func (l *LoadGenerator) pacedCall() {
	l.rl.Take()
	result := make(chan CallResult)
	requestCtx, cancel := context.WithTimeout(context.Background(), l.cfg.CallTimeout)
	callStartTS := time.Now()
	l.responsesWaitGroup.Add(1)
	go func() {
		defer l.responsesWaitGroup.Done()
		select {
		case result <- l.gun.Call(l):
		case <-requestCtx.Done():
			cr := CallResult{Duration: time.Since(callStartTS), Timeout: true, Failed: true, Error: ErrCallTimeout.Error()}
			if l.cfg.LokiConfig != nil {
				l.lokiResponsesChan <- cr
			}
			l.stats.RunFailed.Store(true)
			l.stats.CallTimeout.Add(1)

			l.errsMu.Lock()
			defer l.errsMu.Unlock()
			l.errs = append(l.errs, ErrCallTimeout.Error())

			l.responsesData.failResponsesMu.Lock()
			defer l.responsesData.failResponsesMu.Unlock()
			l.responsesData.FailResponses = append(l.responsesData.FailResponses, cr)
			return
		}
	}()
	l.responsesWaitGroup.Add(1)
	go func() {
		defer l.responsesWaitGroup.Done()
		select {
		case <-requestCtx.Done():
			return
		case res := <-result:
			defer close(result)
			res.Duration = time.Since(callStartTS)
			l.handleCallResult(res)
		}
		cancel()
	}()
}

func (l *LoadGenerator) flushLokiStream() {
	if l.cfg.LokiConfig != nil {
		l.log.Info().Msg("Stopping Loki")
		l.loki.Stop()
		l.log.Info().Msg("Loki exited")
	}
}

// Run runs load loop until timeout or stop
func (l *LoadGenerator) Run() {
	l.log.Info().Msg("Load generator started")
	l.runSchedule()
	l.printStatsLoop()
	if l.cfg.LokiConfig != nil {
		l.runLokiPromtailResponses()
		l.runLokiPromtailStats()
	}
	l.collectData()

	if l.cfg.Schedule.Type == RPSScheduleType {
		l.responsesWaitGroup.Add(1)
		go func() {
			for {
				select {
				case <-l.ResponsesCtx.Done():
					l.responsesWaitGroup.Done()
					l.log.Info().Msg("RPS generator stopped, waiting for requests to finish")
					return
				default:
					l.pacedCall()
				}
			}
		}()
	}
}

// Stop stops load generator, waiting for all calls for either finish or timeout
func (l *LoadGenerator) Stop() (interface{}, bool) {
	l.responsesCancel()
	l.responsesWaitGroup.Wait()
	if l.cfg.LokiConfig != nil {
		l.handleLokiStatsPayload()
		l.dataCancel()
		l.flushLokiStream()
	}
	return l.GetData(), l.stats.RunFailed.Load()
}

// Wait waits until test ends
func (l *LoadGenerator) Wait() (interface{}, bool) {
	l.responsesWaitGroup.Wait()
	if l.cfg.LokiConfig != nil {
		l.handleLokiStatsPayload()
		l.dataCancel()
		l.flushLokiStream()
	}
	return l.GetData(), l.stats.RunFailed.Load()
}

// Errors get all calls errors
func (l *LoadGenerator) Errors() []string {
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

// StatsJSON get all load stats for export
func (l *LoadGenerator) StatsJSON() map[string]interface{} {
	return map[string]interface{}{
		"current_rps":       l.stats.CurrentRPS.Load(),
		"current_instances": l.stats.CurrentInstances.Load(),
		"run_stopped":       l.stats.RunStopped.Load(),
		"run_failed":        l.stats.RunFailed.Load(),
		"failed":            l.stats.Failed.Load(),
		"success":           l.stats.Success.Load(),
		"callTimeout":       l.stats.CallTimeout.Load(),
	}
}

// printStatsLoop prints stats periodically, with LoadGeneratorConfig.StatsPollInterval
func (l *LoadGenerator) printStatsLoop() {
	l.responsesWaitGroup.Add(1)
	go func() {
		defer l.responsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				l.log.Info().Msg("Stats loop exited")
				return
			default:
				time.Sleep(l.cfg.StatsPollInterval)
				l.log.Info().
					Int64("Success", l.stats.Success.Load()).
					Int64("Failed", l.stats.Failed.Load()).
					Int64("CallTimeout", l.stats.CallTimeout.Load()).
					Msg("Load stats")
			}
		}
	}()
}
