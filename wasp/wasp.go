package wasp

import (
	"context"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
)

const (
	DefaultCallTimeout           = 1 * time.Minute
	DefaultSetupTimeout          = 1 * time.Minute
	DefaultTeardownTimeout       = 1 * time.Minute
	DefaultStatsPollInterval     = 5 * time.Second
	DefaultRateLimitUnitDuration = 1 * time.Second
	DefaultCallResultBufLen      = 50000
	DefaultGenName               = "Generator"
)

var (
	ErrNoCfg                  = errors.New("config is nil")
	ErrNoImpl                 = errors.New("either \"gun\" or \"vu\" implementation must provided")
	ErrNoSchedule             = errors.New("no schedule segments were provided")
	ErrInvalidScheduleType    = errors.New("schedule type must be either of wasp.RPS, wasp.VU, use package constants")
	ErrCallTimeout            = errors.New("generator request call timeout")
	ErrSetupTimeout           = errors.New("generator request setup timeout")
	ErrSetup                  = errors.New("generator request setup error")
	ErrTeardownTimeout        = errors.New("generator request teardown timeout")
	ErrTeardown               = errors.New("generator request teardown error")
	ErrStartFrom              = errors.New("from must be > 0")
	ErrInvalidSegmentDuration = errors.New("SegmentDuration must be defined")
	ErrNoGun                  = errors.New("rps load scheduleSegments selected but gun implementation is nil")
	ErrNoVU                   = errors.New("vu load scheduleSegments selected but vu implementation is nil")
	ErrInvalidLabels          = errors.New("invalid Loki labels, labels should be [a-z][A-Z][0-9] and _")
)

// Gun is basic interface for some synthetic load test implementation
// Call performs one request according to some RPS schedule
type Gun interface {
	Call(l *Generator) *Response
}

// VirtualUser is basic interface to run virtual users load
// you should use it if:
// - your protocol is stateful, ex.: ws, grpc
// - you'd like to have some VirtualUser modelling, perform sequential requests
type VirtualUser interface {
	Call(l *Generator)
	Clone(l *Generator) VirtualUser
	Setup(l *Generator) error
	Teardown(l *Generator) error
	Stop(l *Generator)
	StopChan() chan struct{}
}

// NewVUControl creates and returns a new instance of VUControl.
// It initializes the stop channel with a buffer size of 1, allowing
// for controlled stopping of virtual user operations.
func NewVUControl() *VUControl {
	return &VUControl{stop: make(chan struct{}, 1)}
}

// VUControl is a base VU that allows us to control the schedule and bring VUs up and down
type VUControl struct {
	stop chan struct{}
}

// Stop sends a signal to the stop channel to halt the operation of the VUControl.
// It takes a pointer to a Generator as an argument, but does not utilize it.
func (m *VUControl) Stop(_ *Generator) {
	m.stop <- struct{}{}
}

// StopChan returns a channel of type struct{} that signals when the VUControl should stop.
// It provides a mechanism to gracefully shut down or interrupt the VUControl's operations.
func (m *VUControl) StopChan() chan struct{} {
	return m.stop
}

// Response represents basic result info
type Response struct {
	Failed     bool          `json:"failed,omitempty"`
	Timeout    bool          `json:"timeout,omitempty"`
	StatusCode string        `json:"status_code,omitempty"`
	Path       string        `json:"path,omitempty"`
	Duration   time.Duration `json:"duration"`
	StartedAt  *time.Time    `json:"started_at,omitempty"`
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
	Group      string        `json:"group"`
	Data       interface{}   `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
}

type ScheduleType string

const (
	RPS ScheduleType = "rps_schedule"
	VU  ScheduleType = "vu_schedule"
)

// Segment load test schedule segment
type Segment struct {
	From     int64
	Duration time.Duration
}

// Validate checks the Segment for validity by ensuring that the 'From' field is greater than zero and the 'Duration' field is non-zero. 
// It returns an error if any of these conditions are not met, indicating an invalid segment configuration.
func (ls *Segment) Validate() error {
	if ls.From <= 0 {
		return ErrStartFrom
	}
	if ls.Duration == 0 {
		return ErrInvalidSegmentDuration
	}
	return nil
}

// Config is for shared load test data and configuration
type Config struct {
	T                     *testing.T
	GenName               string
	LoadType              ScheduleType
	Labels                map[string]string
	LokiConfig            *LokiConfig
	Schedule              []*Segment
	RateLimitUnitDuration time.Duration
	CallResultBufLen      int
	StatsPollInterval     time.Duration
	CallTimeout           time.Duration
	SetupTimeout          time.Duration
	TeardownTimeout       time.Duration
	FailOnErr             bool
	Gun                   Gun
	VU                    VirtualUser
	Logger                zerolog.Logger
	SharedData            interface{}
	SamplerConfig         *SamplerConfig
	// calculated fields
	duration time.Duration
	// only available in cluster mode
	nodeID string
}

// Validate checks the Config object for required fields and assigns default values where necessary.
// It ensures that the configuration is complete and consistent, returning an error if any required
// fields are missing or invalid. This includes setting default timeouts, buffer lengths, and names,
// as well as verifying the presence of necessary components based on the load type.
func (lgc *Config) Validate() error {
	if lgc.CallTimeout == 0 {
		lgc.CallTimeout = DefaultCallTimeout
	}
	if lgc.SetupTimeout == 0 {
		lgc.SetupTimeout = DefaultSetupTimeout
	}
	if lgc.TeardownTimeout == 0 {
		lgc.TeardownTimeout = DefaultTeardownTimeout
	}
	if lgc.StatsPollInterval == 0 {
		lgc.StatsPollInterval = DefaultStatsPollInterval
	}
	if lgc.CallResultBufLen == 0 {
		lgc.CallResultBufLen = DefaultCallResultBufLen
	}
	if lgc.GenName == "" {
		lgc.GenName = DefaultGenName
	}
	if lgc.Gun == nil && lgc.VU == nil {
		return ErrNoImpl
	}
	if lgc.Schedule == nil {
		return ErrNoSchedule
	}
	if lgc.LoadType != RPS && lgc.LoadType != VU {
		return ErrInvalidScheduleType
	}
	if lgc.LoadType == RPS && lgc.Gun == nil {
		return ErrNoGun
	}
	if lgc.LoadType == VU && lgc.VU == nil {
		return ErrNoVU
	}
	if lgc.RateLimitUnitDuration == 0 {
		lgc.RateLimitUnitDuration = DefaultRateLimitUnitDuration
	}
	return nil
}

// Stats basic generator load stats
type Stats struct {
	// TODO: update json labels with dashboards on major release
	CurrentRPS      atomic.Int64 `json:"currentRPS"`
	CurrentTimeUnit int64        `json:"current_time_unit"`
	CurrentVUs      atomic.Int64 `json:"currentVUs"`
	LastSegment     atomic.Int64 `json:"last_segment"`
	CurrentSegment  atomic.Int64 `json:"current_schedule_segment"`
	SamplesRecorded atomic.Int64 `json:"samples_recorded"`
	SamplesSkipped  atomic.Int64 `json:"samples_skipped"`
	RunStarted      atomic.Bool  `json:"runStarted"`
	RunPaused       atomic.Bool  `json:"runPaused"`
	RunStopped      atomic.Bool  `json:"runStopped"`
	RunFailed       atomic.Bool  `json:"runFailed"`
	Success         atomic.Int64 `json:"success"`
	Failed          atomic.Int64 `json:"failed"`
	CallTimeout     atomic.Int64 `json:"callTimeout"`
	Duration        int64        `json:"load_duration"`
}

// ResponseData includes any request/response data that a gun might store
// ok* slices usually contains successful responses and their verifications if their done async
// fail* slices contains CallResult with response data and an error
type ResponseData struct {
	okDataMu        *sync.Mutex
	OKData          *SliceBuffer[any]
	okResponsesMu   *sync.Mutex
	OKResponses     *SliceBuffer[*Response]
	failResponsesMu *sync.Mutex
	FailResponses   *SliceBuffer[*Response]
}

// Generator generates load with some RPS
type Generator struct {
	Cfg                *Config
	sampler            *Sampler
	Log                zerolog.Logger
	labels             model.LabelSet
	rl                 atomic.Pointer[ratelimit.Limiter]
	scheduleSegments   []*Segment
	currentSegmentMu   *sync.Mutex
	currentSegment     *Segment
	ResponsesWaitGroup *sync.WaitGroup
	dataWaitGroup      *sync.WaitGroup
	ResponsesCtx       context.Context
	responsesCancel    context.CancelFunc
	dataCtx            context.Context
	dataCancel         context.CancelFunc
	gun                Gun
	vu                 VirtualUser
	vus                []VirtualUser
	ResponsesChan      chan *Response
	Responses          *Responses
	responsesData      *ResponseData
	errsMu             *sync.Mutex
	errs               *SliceBuffer[string]
	stats              *Stats
	loki               *LokiClient
	lokiResponsesChan  chan *Response
}

// NewGenerator initializes a new Generator instance using the provided configuration.
// It validates the configuration and its schedule, calculates the total duration,
// and sets up contexts for responses and data collection. It also initializes
// logging, labels, and response channels. If Loki configuration is provided,
// it creates a Loki client. It returns the Generator instance or an error if
// any validation or initialization fails.
func NewGenerator(cfg *Config) (*Generator, error) {
	if cfg == nil {
		return nil, ErrNoCfg
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	for _, s := range cfg.Schedule {
		if err := s.Validate(); err != nil {
			return nil, err
		}
	}
	for _, s := range cfg.Schedule {
		cfg.duration += s.Duration
	}
	l := GetLogger(cfg.T, cfg.GenName)

	ls := LabelsMapToModel(cfg.Labels)
	if cfg.T != nil {
		ls = ls.Merge(model.LabelSet{
			"go_test_name": model.LabelValue(cfg.T.Name()),
		})
	}
	if cfg.GenName != "" {
		ls = ls.Merge(model.LabelSet{
			"gen_name": model.LabelValue(cfg.GenName),
		})
	}
	if err := ls.Validate(); err != nil {
		return nil, ErrInvalidLabels
	}
	cfg.nodeID = os.Getenv("WASP_NODE_ID")
	// context for all requests/responses and vus
	responsesCtx, responsesCancel := context.WithTimeout(context.Background(), cfg.duration)
	// context for all the collected data
	dataCtx, dataCancel := context.WithCancel(context.Background())
	rch := make(chan *Response)
	g := &Generator{
		Cfg:                cfg,
		sampler:            NewSampler(cfg.SamplerConfig),
		scheduleSegments:   cfg.Schedule,
		ResponsesWaitGroup: &sync.WaitGroup{},
		dataWaitGroup:      &sync.WaitGroup{},
		ResponsesCtx:       responsesCtx,
		responsesCancel:    responsesCancel,
		dataCtx:            dataCtx,
		dataCancel:         dataCancel,
		gun:                cfg.Gun,
		vu:                 cfg.VU,
		Responses:          NewResponses(rch),
		ResponsesChan:      rch,
		labels:             ls,
		currentSegmentMu:   &sync.Mutex{},
		responsesData: &ResponseData{
			okDataMu:        &sync.Mutex{},
			OKData:          NewSliceBuffer[any](cfg.CallResultBufLen),
			okResponsesMu:   &sync.Mutex{},
			OKResponses:     NewSliceBuffer[*Response](cfg.CallResultBufLen),
			failResponsesMu: &sync.Mutex{},
			FailResponses:   NewSliceBuffer[*Response](cfg.CallResultBufLen),
		},
		errsMu:            &sync.Mutex{},
		errs:              NewSliceBuffer[string](cfg.CallResultBufLen),
		stats:             &Stats{},
		Log:               l,
		lokiResponsesChan: make(chan *Response, 50000),
	}
	var err error
	if cfg.LokiConfig != nil {
		g.loki, err = NewLokiClient(cfg.LokiConfig)
		if err != nil {
			return nil, err
		}
	}
	CPUCheckLoop()
	return g, nil
}

// runExecuteLoop initiates the execution loop for the load generator based on the configured load type.
// For RPS (Requests Per Second), it starts a goroutine to handle paced calls controlled by the current RPS.
// For VU (Virtual Users), it initializes and runs the specified number of virtual users concurrently.
// It updates the generator's current segment and statistics accordingly.
func (g *Generator) runExecuteLoop() {
	g.currentSegment = g.scheduleSegments[0]
	g.stats.LastSegment.Store(int64(len(g.scheduleSegments)))
	switch g.Cfg.LoadType {
	case RPS:
		g.ResponsesWaitGroup.Add(1)
		// we run pacedCall controlled by stats.CurrentRPS
		go func() {
			for {
				select {
				case <-g.ResponsesCtx.Done():
					g.ResponsesWaitGroup.Done()
					g.Log.Info().Msg("RPS generator has stopped")
					return
				default:
					g.pacedCall()
				}
			}
		}()
	case VU:
		g.currentSegmentMu.Lock()
		g.stats.CurrentVUs.Store(g.currentSegment.From)
		g.currentSegmentMu.Unlock()
		// we start all vus once
		vus := g.stats.CurrentVUs.Load()
		for i := 0; i < int(vus); i++ {
			inst := g.vu.Clone(g)
			g.runVU(inst)
			g.vus = append(g.vus, inst)
		}
	}
}

// runSetupWithTimeout initiates the setup process for a VirtualUser with a specified timeout.
// It returns true if the setup completes successfully within the timeout period, otherwise false.
// If the setup fails or times out, an error response is sent to the ResponsesChan.
func (g *Generator) runSetupWithTimeout(vu VirtualUser) bool {
	startedAt := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), g.Cfg.SetupTimeout)
	defer cancel()
	setupChan := make(chan bool)
	go func() {
		if err := vu.Setup(g); err != nil {
			g.ResponsesChan <- &Response{StartedAt: &startedAt, Error: errors.Wrap(err, ErrSetup.Error()).Error(), Failed: true}
			setupChan <- false
		} else {
			setupChan <- true
		}
	}()
	select {
	case <-ctx.Done():
		g.ResponsesChan <- &Response{StartedAt: &startedAt, Error: ErrSetupTimeout.Error(), Timeout: true}
		return false
	case success := <-setupChan:
		return success
	}
}

// runTeardownWithTimeout executes the Teardown method of a VirtualUser with a specified timeout.
// It returns true if the teardown completes successfully within the timeout period, otherwise false.
// If the teardown fails or times out, an error response is sent to the ResponsesChan.
func (g *Generator) runTeardownWithTimeout(vu VirtualUser) bool {
	startedAt := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), g.Cfg.TeardownTimeout)
	defer cancel()
	setupChan := make(chan bool)
	go func() {
		if err := vu.Teardown(g); err != nil {
			g.ResponsesChan <- &Response{StartedAt: &startedAt, Error: errors.Wrap(err, ErrTeardown.Error()).Error(), Failed: true}
			setupChan <- false
		} else {
			setupChan <- true
		}
	}()
	select {
	case <-ctx.Done():
		g.ResponsesChan <- &Response{StartedAt: &startedAt, Error: ErrTeardownTimeout.Error(), Timeout: true}
		return false
	case success := <-setupChan:
		return success
	}
}

// runVU initiates and manages the execution of a VirtualUser (VU) in a separate goroutine.
// It adds to the ResponsesWaitGroup to track active VUs and ensures proper teardown upon completion or timeout.
// The function handles setup, execution, and teardown phases, monitoring for pause signals and context cancellations.
// It sends responses to the ResponsesChan upon timeout or successful execution.
func (g *Generator) runVU(vu VirtualUser) {
	g.ResponsesWaitGroup.Add(1)
	go func() {
		defer g.ResponsesWaitGroup.Done()
		if !g.runSetupWithTimeout(vu) {
			return
		}
		for {
			if g.stats.RunPaused.Load() {
				continue
			}
			startedAt := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), g.Cfg.CallTimeout)
			vuChan := make(chan struct{})
			go func() {
				vu.Call(g)
				select {
				case <-ctx.Done():
					return
				default:
					vuChan <- struct{}{}
				}
			}()
			select {
			case <-g.ResponsesCtx.Done():
				cancel()
				return
			case <-vu.StopChan():
				g.runTeardownWithTimeout(vu)
				cancel()
				return
			case <-ctx.Done():
				g.ResponsesChan <- &Response{StartedAt: &startedAt, Error: ErrCallTimeout.Error(), Timeout: true}
				cancel()
			case <-vuChan:
			}
		}
	}()
}

// processSegment updates the current segment of the generator's schedule.
// It adjusts the rate limit or virtual users (VUs) based on the load type
// configuration. It returns true if the current segment is the last segment,
// indicating the schedule is complete, otherwise returns false.
func (g *Generator) processSegment() bool {
	defer func() {
		g.stats.RunStarted.Store(true)
		g.Log.Info().
			Int64("Segment", g.stats.CurrentSegment.Load()).
			Int64("VUs", g.stats.CurrentVUs.Load()).
			Int64("RPS", g.stats.CurrentRPS.Load()).
			Msg("Schedule segment")
	}()
	if g.stats.CurrentSegment.Load() == g.stats.LastSegment.Load() {
		return true
	}
	g.currentSegmentMu.Lock()
	g.currentSegment = g.scheduleSegments[g.stats.CurrentSegment.Load()]
	g.currentSegmentMu.Unlock()
	g.stats.CurrentSegment.Add(1)
	switch g.Cfg.LoadType {
	case RPS:
		newRateLimit := ratelimit.New(int(g.currentSegment.From), ratelimit.Per(g.Cfg.RateLimitUnitDuration), ratelimit.WithoutSlack)
		g.rl.Store(&newRateLimit)
		g.stats.CurrentRPS.Store(g.currentSegment.From)
	case VU:
		oldVUs := g.stats.CurrentVUs.Load()
		newVUs := g.currentSegment.From
		g.stats.CurrentVUs.Store(newVUs)

		vusToSpawn := int(math.Abs(float64(max(oldVUs, g.currentSegment.From) - min(oldVUs, g.currentSegment.From))))
		log.Debug().Int64("OldVUs", oldVUs).Int64("NewVUs", newVUs).Int("VUsDelta", vusToSpawn).Msg("Changing VUs")
		if oldVUs == newVUs {
			return false
		}
		if oldVUs > g.currentSegment.From {
			for i := 0; i < vusToSpawn; i++ {
				g.vus[i].Stop(g)
			}
			g.vus = g.vus[vusToSpawn:]
		} else {
			for i := 0; i < vusToSpawn; i++ {
				inst := g.vu.Clone(g)
				g.runVU(inst)
				g.vus = append(g.vus, inst)
			}
		}
	}
	return false
}

// runScheduleLoop starts a goroutine that continuously processes schedule segments.
// It listens for a cancellation signal from ResponsesCtx to terminate the loop and log completion.
// The loop pauses between segment processing based on the current segment's duration.
func (g *Generator) runScheduleLoop() {
	go func() {
		for {
			select {
			case <-g.ResponsesCtx.Done():
				g.Log.Info().
					Int64("Segment", g.stats.CurrentSegment.Load()).
					Int64("VUs", g.stats.CurrentVUs.Load()).
					Int64("RPS", g.stats.CurrentRPS.Load()).
					Msg("Finished all schedule segments")
				g.Log.Info().Msg("Scheduler exited")
				return
			default:
				if g.processSegment() {
					return
				}
				time.Sleep(g.currentSegment.Duration)
			}
		}
	}()
}

// storeResponses processes and stores the given response based on its status.
// It checks if the response should be recorded and updates statistics accordingly.
// If the response indicates a failure or timeout, it logs the error and updates failure statistics.
// Successful responses are stored separately. If configured, it may stop the generator on errors.
func (g *Generator) storeResponses(res *Response) {
	if g.Cfg.CallTimeout > 0 && res.Duration > g.Cfg.CallTimeout && !res.Timeout {
		return
	}
	if !g.sampler.ShouldRecord(res, g.stats) {
		return
	}
	if g.Cfg.LokiConfig != nil {
		g.lokiResponsesChan <- res
	}
	g.responsesData.okDataMu.Lock()
	g.responsesData.failResponsesMu.Lock()
	g.errsMu.Lock()
	if res.Failed {
		g.stats.RunFailed.Store(true)
		g.stats.Failed.Add(1)
		g.errs.Append(res.Error)
		g.responsesData.FailResponses.Append(res)
		g.Log.Error().Str("Err", res.Error).Msg("load generator request failed")
	} else if res.Timeout {
		g.stats.RunFailed.Store(true)
		g.stats.CallTimeout.Add(1)
		g.stats.Failed.Add(1)
		g.errs.Append(res.Error)
		g.responsesData.FailResponses.Append(res)
		g.Log.Error().Str("Err", res.Error).Msg("load generator request timed out")
	} else {
		g.stats.Success.Add(1)
		g.responsesData.OKData.Append(res.Data)
		g.responsesData.OKResponses.Append(res)
	}
	g.responsesData.okDataMu.Unlock()
	g.responsesData.failResponsesMu.Unlock()
	g.errsMu.Unlock()
	if (g.stats.Failed.Load() > 0 || g.stats.CallTimeout.Load() > 0) && g.Cfg.FailOnErr {
		g.Log.Warn().Msg("Generator has stopped on first error")
		g.responsesCancel()
	}
}

// collectVUResults collects virtual user results from the ResponsesChan channel.
// It calculates the duration and finished time for each response and stores them.
// The function runs in a separate goroutine and exits when the data context is done.
func (g *Generator) collectVUResults() {
	if g.Cfg.LoadType == RPS {
		return
	}
	g.dataWaitGroup.Add(1)
	go func() {
		defer g.dataWaitGroup.Done()
		for {
			select {
			case <-g.dataCtx.Done():
				g.Log.Info().Msg("Collect data exited")
				return
			case res := <-g.ResponsesChan:
				if res.StartedAt != nil {
					res.Duration = time.Since(*res.StartedAt)
				}
				tn := time.Now()
				res.FinishedAt = &tn
				g.storeResponses(res)
			}
		}
	}()
}

// pacedCall executes a controlled call to a target service, respecting the configured rate limits and timeouts.
// It checks the generator's run state and manages the lifecycle of the call, including handling timeouts and storing responses.
func (g *Generator) pacedCall() {
	if !g.Stats().RunStarted.Load() {
		return
	}
	l := *g.rl.Load()
	l.Take()
	if g.stats.RunPaused.Load() {
		return
	}
	if g.stats.RunStopped.Load() {
		return
	}
	result := make(chan *Response)
	requestCtx, cancel := context.WithTimeout(context.Background(), g.Cfg.CallTimeout)
	callStartTS := time.Now()
	go func() {
		result <- g.gun.Call(g)
	}()
	g.ResponsesWaitGroup.Add(1)
	go func() {
		defer g.ResponsesWaitGroup.Done()
		select {
		case <-requestCtx.Done():
			ts := time.Now()
			cr := &Response{Duration: time.Since(callStartTS), FinishedAt: &ts, Timeout: true, Error: ErrCallTimeout.Error()}
			g.storeResponses(cr)
		case res := <-result:
			defer close(result)
			res.Duration = time.Since(callStartTS)
			ts := time.Now()
			res.FinishedAt = &ts
			g.storeResponses(res)
		}
		cancel()
	}()
}

// Run initiates the load generation process for the Generator instance.
// It logs the start of the process, manages statistics, and optionally sends
// data to Loki if configured. The function executes scheduled and execution
// loops, collects virtual user results, and waits for completion if specified.
// It returns the result of the wait operation or nil if not waiting.
func (g *Generator) Run(wait bool) (interface{}, bool) {
	g.Log.Info().Msg("Load generator started")
	g.printStatsLoop()
	if g.Cfg.LokiConfig != nil {
		g.sendResponsesToLoki()
		g.sendStatsToLoki()
	}
	g.runScheduleLoop()
	g.runExecuteLoop()
	g.collectVUResults()
	if wait {
		return g.Wait()
	}
	return nil, false
}

// Pause logs a warning message indicating that the generator has been paused.
// It updates the generator's statistics to reflect the paused state.
func (g *Generator) Pause() {
	g.Log.Warn().Msg("Generator was paused")
	g.stats.RunPaused.Store(true)
}

// Resume resumes the operation of the Generator by updating its state to indicate it is no longer paused. It logs a warning message indicating the resumption.
func (g *Generator) Resume() {
	g.Log.Warn().Msg("Generator was resumed")
	g.stats.RunPaused.Store(false)
}

// Stop halts the generator's operation if it has not already been stopped.
// It sets the run status to stopped and failed, logs a warning, and cancels
// any ongoing responses. It then waits for all responses to finish and returns
// the result of the wait operation, including any data collected and a boolean
// indicating if the run failed.
func (g *Generator) Stop() (interface{}, bool) {
	if g.stats.RunStopped.Load() {
		return nil, true
	}
	g.stats.RunStarted.Store(false)
	g.stats.RunStopped.Store(true)
	g.stats.RunFailed.Store(true)
	g.Log.Warn().Msg("Graceful stop")
	g.responsesCancel()
	return g.Wait()
}

// Wait blocks until all responses and data processing are complete.
// It logs the waiting process, updates statistics, and handles Loki streaming if configured.
// It returns the collected data and a boolean indicating if the run failed.
func (g *Generator) Wait() (interface{}, bool) {
	g.Log.Info().Msg("Waiting for all responses to finish")
	g.ResponsesWaitGroup.Wait()
	g.stats.Duration = g.Cfg.duration.Nanoseconds()
	g.stats.CurrentTimeUnit = g.Cfg.RateLimitUnitDuration.Nanoseconds()
	if g.Cfg.LokiConfig != nil {
		g.dataCancel()
		g.dataWaitGroup.Wait()
		g.stopLokiStream()
	}
	return g.GetData(), g.stats.RunFailed.Load()
}

// InputSharedData retrieves the shared data from the generator's configuration.
// It returns the shared data as an interface{}.
func (g *Generator) InputSharedData() interface{} {
	return g.Cfg.SharedData
}

// Errors returns a slice of strings containing the error messages collected by the Generator.
// It retrieves the error data stored in the Generator's internal structure.
func (g *Generator) Errors() []string {
	return g.errs.Data
}

// GetData returns the response data collected by the Generator.
// It provides access to the stored ResponseData after all responses have been processed.
func (g *Generator) GetData() *ResponseData {
	return g.responsesData
}

// Stats returns the current statistics of the Generator.
// It provides access to the Generator's runtime state, such as whether the run has started, paused, or stopped.
func (g *Generator) Stats() *Stats {
	return g.stats
}

/* Loki's methods to handle CallResult/Stats and stream it to Loki */

// stopLokiStream stops the Loki stream if the Loki configuration is present and valid.
// It logs the stopping and exiting of the Loki process.
func (g *Generator) stopLokiStream() {
	if g.Cfg.LokiConfig != nil && g.Cfg.LokiConfig.URL != "" {
		g.Log.Info().Msg("Stopping Loki")
		g.loki.StopNow()
		g.Log.Info().Msg("Loki exited")
	}
}

// handleLokiResponsePayload processes a Loki response payload by merging labels,
// removing unnecessary timestamps, and sending the structured data to Loki.
// It logs any errors encountered and stops the generator if an error occurs.
func (g *Generator) handleLokiResponsePayload(r *Response) {
	labels := g.labels.Merge(model.LabelSet{
		"test_data_type": "responses",
		CallGroupLabel:   model.LabelValue(r.Group),
	})
	// we are removing time.Time{} because when it marshalled to string it creates N responses for some Loki queries
	// and to minimize the payload, duration is already calculated at that point
	ts := r.FinishedAt
	r.StartedAt = nil
	r.FinishedAt = nil
	err := g.loki.HandleStruct(labels, *ts, r)
	if err != nil {
		g.Log.Err(err).Send()
		g.Stop()
	}
}

// handleLokiStatsPayload merges the generator's labels with a predefined label set
// and sends the resulting structured data to Loki. If an error occurs during this
// process, it logs the error and stops the generator. This function is typically
// used in a loop to periodically send statistics to Loki.
func (g *Generator) handleLokiStatsPayload() {
	ls := g.labels.Merge(model.LabelSet{
		"test_data_type": "stats",
	})
	err := g.loki.HandleStruct(ls, time.Now(), g.StatsJSON())
	if err != nil {
		g.Log.Err(err).Send()
		g.Stop()
	}
}

// sendResponsesToLoki streams data to Loki using the configured URL and labels.
// It listens for responses on the lokiResponsesChan channel and handles them
// until the data context is done. This function is typically called when the
// Loki configuration is present, as part of the data streaming process.
func (g *Generator) sendResponsesToLoki() {
	g.Log.Info().
		Str("URL", g.Cfg.LokiConfig.URL).
		Interface("DefaultLabels", g.Cfg.Labels).
		Msg("Streaming data to Loki")
	g.dataWaitGroup.Add(1)
	go func() {
		defer g.dataWaitGroup.Done()
		for {
			select {
			case <-g.dataCtx.Done():
				g.Log.Info().Msg("Loki responses exited")
				return
			case r := <-g.lokiResponsesChan:
				g.handleLokiResponsePayload(r)
			}
		}
	}()
}

// sendStatsToLoki initiates a goroutine to periodically send statistics to Loki.
// It adds to the dataWaitGroup to ensure synchronization and listens for context
// cancellation to gracefully exit. The function operates based on the configured
// StatsPollInterval and handles the payload for Loki statistics.
func (g *Generator) sendStatsToLoki() {
	g.dataWaitGroup.Add(1)
	go func() {
		defer g.dataWaitGroup.Done()
		for {
			select {
			case <-g.dataCtx.Done():
				g.Log.Info().Msg("Loki stats exited")
				return
			default:
				time.Sleep(g.Cfg.StatsPollInterval)
				g.handleLokiStatsPayload()
			}
		}
	}()
}

/* Local logging methods */

// StatsJSON returns a map containing various statistics and metrics of the Generator.
// The map includes information such as node ID, current requests per second,
// current instances, samples recorded and skipped, run status, and success or failure counts. 
// This data is typically used for monitoring and logging purposes.
func (g *Generator) StatsJSON() map[string]interface{} {
	return map[string]interface{}{
		"node_id":           g.Cfg.nodeID,
		"current_rps":       g.stats.CurrentRPS.Load(),
		"current_instances": g.stats.CurrentVUs.Load(),
		"samples_recorded":  g.stats.SamplesRecorded.Load(),
		"samples_skipped":   g.stats.SamplesSkipped.Load(),
		"run_stopped":       g.stats.RunStopped.Load(),
		"run_failed":        g.stats.RunFailed.Load(),
		"failed":            g.stats.Failed.Load(),
		"success":           g.stats.Success.Load(),
		"callTimeout":       g.stats.CallTimeout.Load(),
		"load_duration":     g.stats.Duration,
		"current_time_unit": g.stats.CurrentTimeUnit,
	}
}

// printStatsLoop starts a goroutine that periodically logs the current statistics of the generator.
// It continues to log until the context is done, at which point it logs an exit message and returns.
func (g *Generator) printStatsLoop() {
	g.ResponsesWaitGroup.Add(1)
	go func() {
		defer g.ResponsesWaitGroup.Done()
		for {
			select {
			case <-g.ResponsesCtx.Done():
				g.Log.Info().Msg("Stats loop exited")
				return
			default:
				time.Sleep(g.Cfg.StatsPollInterval)
				g.Log.Info().
					Int64("Success", g.stats.Success.Load()).
					Int64("Failed", g.stats.Failed.Load()).
					Int64("CallTimeout", g.stats.CallTimeout.Load()).
					Msg("Load stats")
			}
		}
	}()
}

// LabelsMapToModel converts a map of strings to a model.LabelSet.
// Each key-value pair in the map is transformed into a model.LabelName
// and model.LabelValue, respectively, and added to the LabelSet.
// It returns the constructed model.LabelSet.
func LabelsMapToModel(m map[string]string) model.LabelSet {
	ls := model.LabelSet{}
	for k, v := range m {
		ls[model.LabelName(k)] = model.LabelValue(v)
	}
	return ls
}
