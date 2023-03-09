package loadgen

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/client"

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
	ErrNoImpl       = errors.New("either \"gun\" or \"instanceTemplate\" implementation must provided")
	ErrCallTimeout  = errors.New("generator request call timeout")
	ErrStartFrom    = errors.New("from must be > 0")
	ErrScheduleType = errors.New("schedule type must be \"rps_schedule\" or \"instance_schedule\"")
	ErrNoGun        = errors.New("rps load schedule selected but gun implementation is nil")
	ErrNoInstance   = errors.New("instanceTemplate load schedule selected but instanceTemplate implementation is nil")
)

// Gun is basic interface to run limited load with a contract call and save all transactions
type Gun interface {
	Call(l *Generator) CallResult
}

// Instance is basic interface to run load instances
type Instance interface {
	Run(l *Generator)
	Stop(l *Generator)
}

// CallResult represents basic call result info
type CallResult struct {
	Failed     bool          `json:"failed,omitempty"`
	Timeout    bool          `json:"timeout,omitempty"`
	Duration   time.Duration `json:"duration"`
	StartedAt  *time.Time    `json:"started_at,omitempty"`
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
	Data       interface{}   `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
}

const (
	RPSScheduleType       string = "rps_schedule"
	InstancesScheduleType string = "instance_schedule"
)

// LoadSchedule load test schedule
type LoadSchedule struct {
	Type         string
	From         int64
	Increase     int64
	StepDuration time.Duration
	Limit        int64
}

func (ls *LoadSchedule) Validate() error {
	if ls.Type != RPSScheduleType && ls.Type != InstancesScheduleType {
		return ErrScheduleType
	}
	if ls.From <= 0 {
		return ErrStartFrom
	}
	switch ls.Type {
	case RPSScheduleType:
		if ls.Increase < 0 && ls.Limit < 1 {
			ls.Limit = 1
		}
	}
	return nil
}

// LoadGeneratorConfig is for shared load test data and configuration
type LoadGeneratorConfig struct {
	T                 *testing.T
	Labels            map[string]string
	LokiConfig        *client.LokiConfig
	Schedule          *LoadSchedule
	Duration          time.Duration
	StatsPollInterval time.Duration
	CallTimeout       time.Duration
	Gun               Gun
	Instance          Instance
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
	CurrentRPS             atomic.Int64 `json:"currentRPS"`
	CurrentInstances       atomic.Int64 `json:"currentInstances"`
	CurrentScheduleSegment atomic.Int64 `json:"current_schedule_segment"`
	CurrentScheduleStep    atomic.Int64 `json:"current_schedule_step"`
	RunStopped             atomic.Bool  `json:"runStopped"`
	RunFailed              atomic.Bool  `json:"runFailed"`
	Success                atomic.Int64 `json:"success"`
	Failed                 atomic.Int64 `json:"failed"`
	CallTimeout            atomic.Int64 `json:"callTimeout"`
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

// Generator generates load with some RPS
type Generator struct {
	cfg                *LoadGeneratorConfig
	Log                zerolog.Logger
	labels             model.LabelSet
	rl                 ratelimit.Limiter
	schedule           *LoadSchedule
	ResponsesWaitGroup *sync.WaitGroup
	dataWaitGroup      *sync.WaitGroup
	ResponsesCtx       context.Context
	responsesCancel    context.CancelFunc
	dataCtx            context.Context
	dataCancel         context.CancelFunc
	gun                Gun
	instanceTemplate   Instance
	instances          []Instance
	ResponsesChan      chan CallResult
	responsesData      *ResponseData
	errsMu             *sync.Mutex
	errs               []string
	stats              *GeneratorStats
	loki               client.ExtendedLokiClient
	lokiResponsesChan  chan CallResult
}

// NewLoadGenerator creates a new instanceTemplate for a contract,
// shoots for scheduled RPS until timeout, test logic is defined through Gun
func NewLoadGenerator(cfg *LoadGeneratorConfig) (*Generator, error) {
	if cfg == nil {
		return nil, ErrNoCfg
	}
	if err := cfg.Schedule.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.Schedule.Type == RPSScheduleType && cfg.Gun == nil {
		return nil, ErrNoGun
	}
	if cfg.Schedule.Type == InstancesScheduleType && cfg.Instance == nil {
		return nil, ErrNoInstance
	}
	rl := ratelimit.New(int(cfg.Schedule.From))

	// creating logger from *testing.T context or using a global logger
	var l zerolog.Logger
	if cfg.T != nil {
		l = zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(cfg.T))).With().Timestamp().Logger()
	} else {
		l = log.Logger
	}

	var loki client.ExtendedLokiClient
	var ls model.LabelSet
	var err error
	if cfg.LokiConfig != nil {
		if cfg.LokiConfig.URL == "" {
			l.Warn().Msg("Loki config is set but URL is empty, data is collected but won't be pushed anywhere!")
			loki = client.NewMockPromtailClient()
			if err != nil {
				return nil, err
			}
		} else {
			loki, err = client.NewLokiClient(cfg.LokiConfig)
			if err != nil {
				return nil, err
			}
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
	return &Generator{
		cfg:                cfg,
		schedule:           cfg.Schedule,
		rl:                 rl,
		ResponsesWaitGroup: &sync.WaitGroup{},
		dataWaitGroup:      &sync.WaitGroup{},
		ResponsesCtx:       responsesCtx,
		responsesCancel:    responsesCancel,
		dataCtx:            dataCtx,
		dataCancel:         dataCancel,
		gun:                cfg.Gun,
		instanceTemplate:   cfg.Instance,
		ResponsesChan:      make(chan CallResult),
		labels:             ls,
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
		Log:               l,
		lokiResponsesChan: make(chan CallResult, 50000),
	}, nil
}

func (l *Generator) setupSchedule() {
	switch l.cfg.Schedule.Type {
	case RPSScheduleType:
		l.ResponsesWaitGroup.Add(1)
		l.rl = ratelimit.New(int(l.schedule.From))
		l.stats.CurrentRPS.Store(l.schedule.From)

		// we run pacedCall controlled by stats.CurrentRPS
		go func() {
			for {
				select {
				case <-l.ResponsesCtx.Done():
					l.ResponsesWaitGroup.Done()
					l.Log.Info().Msg("RPS generator stopped")
					return
				default:
					l.pacedCall()
				}
			}
		}()
	case InstancesScheduleType:
		l.stats.CurrentInstances.Store(l.schedule.From)
		// we start all instances once
		instances := l.stats.CurrentInstances.Load()
		for i := 0; i < int(instances); i++ {
			inst := l.instanceTemplate
			inst.Run(l)
			l.instances = append(l.instances, inst)
		}
	}
}

// runSchedule runs scheduling loop, changes Generator.currentRPS according to a load schedule
func (l *Generator) runSchedule() {
	l.ResponsesWaitGroup.Add(1)
	go func() {
		defer l.ResponsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				l.Log.Info().Msg("Scheduler exited")
				return
			default:
				time.Sleep(l.schedule.StepDuration)
				log.Info().
					Int64("Instances", l.stats.CurrentInstances.Load()).
					Int64("RPS", l.stats.CurrentRPS.Load()).
					Msg("Scheduler step")
				switch l.cfg.Schedule.Type {
				case RPSScheduleType:
					newRPS := l.stats.CurrentRPS.Load() + l.schedule.Increase
					if newRPS > l.schedule.Limit {
						return
					}
					l.rl = ratelimit.New(int(newRPS))
					l.stats.CurrentRPS.Store(newRPS)
				case InstancesScheduleType:
					if l.schedule.Increase == 0 {
						log.Info().Msg("No instances changes, skipping scheduling")
					}
					if l.schedule.Increase > 0 {
						for i := 0; i < int(l.schedule.Increase); i++ {
							if l.stats.CurrentInstances.Load()+l.schedule.Increase > l.schedule.Limit {
								log.Info().Msg("Upper instances limit reached")
								continue
							}
							instance := l.instanceTemplate
							instance.Run(l)
							l.instances = append(l.instances, instance)
							l.stats.CurrentInstances.Store(l.stats.CurrentInstances.Load() + 1)
						}
					} else {
						log.Info().Int("Instances len", len(l.instances)).Send()
						for i := 0; i < int(math.Abs(float64(l.schedule.Increase))); i++ {
							if l.stats.CurrentInstances.Load()+l.schedule.Increase < l.schedule.Limit {
								log.Info().Msg("Lower instances limit reached")
								continue
							}
							if len(l.instances) == 0 {
								log.Info().Msg("No more instances left")
								continue
							}
							l.instances[0].Stop(l)
							l.instances = l.instances[1:]
							l.stats.CurrentInstances.Store(l.stats.CurrentInstances.Load() - 1)
						}
					}
				}
			}
		}
	}()
}

// handleCallResult stores local metrics for CallResult, pushed them to Loki stream too if Loki is on
func (l *Generator) handleCallResult(res CallResult) {
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

		l.Log.Error().Str("Err", res.Error).Msg("load generator request failed")
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

// collectResults collects CallResult from all the Instances
func (l *Generator) collectResults() {
	if l.cfg.Schedule.Type == RPSScheduleType {
		return
	}
	l.dataWaitGroup.Add(1)
	go func() {
		defer l.dataWaitGroup.Done()
		for {
			select {
			case <-l.dataCtx.Done():
				l.Log.Info().Msg("Collect data exited")
				return
			case res := <-l.ResponsesChan:
				if res.StartedAt.IsZero() {
					log.Error().Msg("StartedAt is not set in instanceTemplate implementation")
					return
				}
				tn := time.Now()
				res.FinishedAt = &tn
				res.Duration = time.Since(*res.StartedAt)
				l.handleCallResult(res)
			}
		}
	}()
}

// pacedCall calls a gun according to a schedule or plain RPS
func (l *Generator) pacedCall() {
	l.rl.Take()
	result := make(chan CallResult)
	requestCtx, cancel := context.WithTimeout(context.Background(), l.cfg.CallTimeout)
	callStartTS := time.Now()
	l.ResponsesWaitGroup.Add(1)
	go func() {
		defer l.ResponsesWaitGroup.Done()
		select {
		case result <- l.gun.Call(l):
		case <-requestCtx.Done():
			ts := time.Now()
			cr := CallResult{Duration: time.Since(callStartTS), FinishedAt: &ts, Timeout: true, Error: ErrCallTimeout.Error()}
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
	l.ResponsesWaitGroup.Add(1)
	go func() {
		defer l.ResponsesWaitGroup.Done()
		select {
		case <-requestCtx.Done():
			return
		case res := <-result:
			defer close(result)
			res.Duration = time.Since(callStartTS)
			ts := time.Now()
			res.FinishedAt = &ts
			l.handleCallResult(res)
		}
		cancel()
	}()
}

// Run runs load loop until timeout or stop
func (l *Generator) Run() {
	l.Log.Info().Msg("Load generator started")
	l.printStatsLoop()
	if l.cfg.LokiConfig != nil {
		l.runLokiPromtailResponses()
		l.runLokiPromtailStats()
	}
	l.collectResults()
	l.setupSchedule()
	l.runSchedule()
}

// Stop stops load generator, waiting for all calls for either finish or timeout
func (l *Generator) Stop() (interface{}, bool) {
	l.responsesCancel()
	return l.Wait()
}

// Wait waits until test ends
func (l *Generator) Wait() (interface{}, bool) {
	l.Log.Info().Msg("Waiting for all responses to finish")
	l.ResponsesWaitGroup.Wait()
	if l.cfg.LokiConfig != nil {
		l.handleLokiStatsPayload()
		l.dataCancel()
		l.dataWaitGroup.Wait()
		l.stopLokiStream()
	}
	return l.GetData(), l.stats.RunFailed.Load()
}

// InputSharedData returns the SharedData passed in Generator config
func (l *Generator) InputSharedData() interface{} {
	return l.cfg.SharedData
}

// Errors get all calls errors
func (l *Generator) Errors() []string {
	return l.errs
}

// GetData get all calls data
func (l *Generator) GetData() *ResponseData {
	return l.responsesData
}

// Stats get all load stats
func (l *Generator) Stats() *GeneratorStats {
	return l.stats
}

/* Loki's methods to handle CallResult/Stats and stream it to Loki */

// stopLokiStream stops the Loki stream client
func (l *Generator) stopLokiStream() {
	if l.cfg.LokiConfig != nil && l.cfg.LokiConfig.URL != "" {
		l.Log.Info().Msg("Stopping Loki")
		l.loki.Stop()
		l.Log.Info().Msg("Loki exited")
	}
}

// handleLokiResponsePayload handles CallResult payload with adding default labels
func (l *Generator) handleLokiResponsePayload(cr CallResult) {
	ls := l.labels.Merge(model.LabelSet{
		"go_test_name":   model.LabelValue(l.cfg.T.Name()),
		"test_data_type": "responses",
	})
	// we are removing time.Time{} because when it marshalled to string it creates N responses for some Loki queries
	// and to minimize the payload, duration is already calculated at that point
	ts := cr.FinishedAt
	cr.StartedAt = nil
	cr.FinishedAt = nil
	err := l.loki.HandleStruct(ls, *ts, cr)
	if err != nil {
		l.Log.Err(err).Send()
	}
}

// handleLokiStatsPayload handles StatsJSON payload with adding default labels
func (l *Generator) handleLokiStatsPayload() {
	ls := l.labels.Merge(model.LabelSet{
		"go_test_name":   model.LabelValue(l.cfg.T.Name()),
		"test_data_type": "stats",
	})
	err := l.loki.HandleStruct(ls, time.Now(), l.StatsJSON())
	if err != nil {
		l.Log.Err(err).Send()
	}
}

// runLokiPromtailResponses pushes CallResult to Loki
func (l *Generator) runLokiPromtailResponses() {
	l.Log.Info().
		Str("URL", l.cfg.LokiConfig.URL).
		Interface("DefaultLabels", l.cfg.Labels).
		Msg("Streaming data to Loki")
	l.dataWaitGroup.Add(1)
	go func() {
		defer l.dataWaitGroup.Done()
		for {
			select {
			case <-l.dataCtx.Done():
				l.Log.Info().Msg("Loki responses exited")
				return
			case r := <-l.lokiResponsesChan:
				l.handleLokiResponsePayload(r)
			}
		}
	}()
}

// runLokiPromtailStats pushes Stats payloads to Loki
func (l *Generator) runLokiPromtailStats() {
	l.dataWaitGroup.Add(1)
	go func() {
		defer l.dataWaitGroup.Done()
		for {
			select {
			case <-l.dataCtx.Done():
				l.Log.Info().Msg("Loki stats exited")
				return
			default:
				time.Sleep(l.cfg.StatsPollInterval)
				l.handleLokiStatsPayload()
			}
		}
	}()
}

/* Local logging methods */

// StatsJSON get all load stats for export
func (l *Generator) StatsJSON() map[string]interface{} {
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
func (l *Generator) printStatsLoop() {
	l.ResponsesWaitGroup.Add(1)
	go func() {
		defer l.ResponsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				l.Log.Info().Msg("Stats loop exited")
				return
			default:
				time.Sleep(l.cfg.StatsPollInterval)
				l.Log.Info().
					Int64("Success", l.stats.Success.Load()).
					Int64("Failed", l.stats.Failed.Load()).
					Int64("CallTimeout", l.stats.CallTimeout.Load()).
					Msg("Load stats")
			}
		}
	}()
}
