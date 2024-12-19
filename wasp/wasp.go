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
	ErrMissingSegmentType     = errors.New("Segment Type myst be set")
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

// NewVUControl creates a new VUControl instance used to manage the lifecycle and control of a virtual user.
func NewVUControl() *VUControl {
	return &VUControl{stop: make(chan struct{}, 1)}
}

// VUControl is a base VU that allows us to control the schedule and bring VUs up and down
type VUControl struct {
	stop chan struct{}
}

// Stop signals VUControl to cease operations by sending a stop signal through the stop channel.
func (m *VUControl) Stop(_ *Generator) {
	m.stop <- struct{}{}
}

// StopChan returns the channel used to signal when the VUControl is stopped.
// It allows consumers to listen for termination events and handle cleanup accordingly.
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

type SegmentType string

const (
	SegmentType_Plain SegmentType = "plain"
	SegmentType_Steps SegmentType = "steps"
)

// Segment load test schedule segment
type Segment struct {
	From      int64         `json:"from"`
	Duration  time.Duration `json:"duration"`
	Type      SegmentType   `json:"type"`
	StartTime time.Time     `json:"time_start"`
	EndTime   time.Time     `json:"time_end"`
}

// Validate checks that the Segment has a valid starting point and duration.
// It returns an error if the starting point is non-positive or the duration is zero.
// Use it to ensure the Segment is properly configured before processing.
func (ls *Segment) Validate() error {
	if ls.From <= 0 {
		return ErrStartFrom
	}
	if ls.Duration == 0 {
		return ErrInvalidSegmentDuration
	}
	if ls.Type == "" {
		return ErrMissingSegmentType
	}

	return nil
}

// Config is for shared load test data and configuration
type Config struct {
	T                     *testing.T        `json:"-"`
	GenName               string            `json:"generator_name"`
	LoadType              ScheduleType      `json:"load_type"`
	Labels                map[string]string `json:"-"`
	LokiConfig            *LokiConfig       `json:"-"`
	Schedule              []*Segment        `json:"schedule"`
	RateLimitUnitDuration time.Duration     `json:"rate_limit_unit_duration"`
	CallResultBufLen      int               `json:"-"`
	StatsPollInterval     time.Duration     `json:"-"`
	CallTimeout           time.Duration     `json:"call_timeout"`
	SetupTimeout          time.Duration     `json:"-"`
	TeardownTimeout       time.Duration     `json:"-"`
	FailOnErr             bool              `json:"-"`
	Gun                   Gun               `json:"-"`
	VU                    VirtualUser       `json:"-"`
	Logger                zerolog.Logger    `json:"-"`
	SharedData            interface{}       `json:"-"`
	SamplerConfig         *SamplerConfig    `json:"-"`
	// calculated fields
	duration time.Duration
	// only available in cluster mode
	nodeID string
}

// Validate checks the Config fields for correctness, sets default values for unset parameters,
// and ensures required configurations are provided. It returns an error if the configuration
// is incomplete or invalid, ensuring the Config is ready for use.
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
	rpsLoopOnce        *sync.Once
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

// NewGenerator initializes a Generator with the provided configuration.
// It validates the config, sets up contexts, logging, and labels.
// Use it to create a Generator for managing service schedules and data collection.
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
		rpsLoopOnce:        &sync.Once{},
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
	return g, nil
}

// runGunLoop runs the generator's Gun loop
// It manages request pacing for RPS after the first segment is loaded.
func (g *Generator) runGunLoop() {
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
}

// runSetupWithTimeout executes the VirtualUser's setup within the configured timeout.
// It returns true if the setup completes successfully before the timeout, otherwise false.
// Use it to ensure that setup processes do not exceed the allowed time.
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

// runTeardownWithTimeout attempts to teardown the given VirtualUser within the configured timeout.
// It returns true if successful, or false if a timeout or error occurs.
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

// runVU starts and manages the execution cycle for a VirtualUser. It handles setup, executes user calls with timeout control, processes responses, and ensures proper teardown. Use it to simulate and manage individual virtual user behavior within the Generator.
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

// processSegment processes the next schedule segment, updating rate limits or virtual users based on configuration.
// It returns true when all segments have been handled, signaling the scheduler to terminate.
func (g *Generator) processSegment() bool {
	defer func() {
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
	g.currentSegment.StartTime = time.Now()
	switch g.Cfg.LoadType {
	case RPS:
		newRateLimit := ratelimit.New(int(g.currentSegment.From), ratelimit.Per(g.Cfg.RateLimitUnitDuration), ratelimit.WithoutSlack)
		g.rl.Store(&newRateLimit)
		g.stats.CurrentRPS.Store(g.currentSegment.From)
		// start Gun loop once, in next segments we control it using g.rl ratelimiter
		g.rpsLoopOnce.Do(func() {
			g.runGunLoop()
		})
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

// runScheduleLoop initiates an asynchronous loop that processes scheduling segments and monitors for completion signals.
// It enables the generator to handle load distribution seamlessly in the background.
func (g *Generator) runScheduleLoop() {
	g.currentSegment = g.scheduleSegments[0]
	g.stats.LastSegment.Store(int64(len(g.scheduleSegments)))
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
				g.currentSegment.EndTime = time.Now()
			}
		}
	}()
}

// storeResponses processes a Response, updating metrics and recording success or failure.
// It is used to handle generator call results for monitoring and error tracking.
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

// collectVUResults launches a background process to receive and store virtual user responses.
// It enables asynchronous collection of performance data during load testing.
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

// pacedCall initiates a rate-limited request to the external service,
// handling timeouts and storing the response.
// It ensures requests adhere to the generator's configuration and execution state.
func (g *Generator) pacedCall() {
	(*g.rl.Load()).Take()
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

// Run starts the Generator’s scheduling and execution workflows, managing logging and metrics.
// If wait is true, it waits for all processes to complete and returns the results.
// Use Run to execute generator tasks either synchronously or asynchronously.
func (g *Generator) Run(wait bool) (interface{}, bool) {
	g.Log.Info().Msg("Load generator started")
	g.printStatsLoop()
	if g.Cfg.LokiConfig != nil {
		g.sendResponsesToLoki()
		g.sendStatsToLoki()
	}
	g.runScheduleLoop()
	g.collectVUResults()
	if wait {
		return g.Wait()
	}
	return nil, false
}

// Pause signals the generator to stop its operations.
// It is used to gracefully halt the generator when pausing activities is required.
func (g *Generator) Pause() {
	g.Log.Warn().Msg("Generator was paused")
	g.stats.RunPaused.Store(true)
}

// Resume resumes the Generator, allowing it to continue operations after being paused.
// It is typically used to restart paused Generators within a Profile or management structure.
func (g *Generator) Resume() {
	g.Log.Warn().Msg("Generator was resumed")
	g.stats.RunPaused.Store(false)
}

// Stop gracefully halts the generator by updating its run state, logging the event, canceling ongoing responses, and waiting for all processes to complete.
// It returns the final data and a boolean indicating whether the run was successfully stopped.
func (g *Generator) Stop() (interface{}, bool) {
	if g.stats.RunStopped.Load() {
		return nil, true
	}
	g.stats.RunStopped.Store(true)
	g.stats.RunFailed.Store(true)
	g.Log.Warn().Msg("Graceful stop")
	g.responsesCancel()
	return g.Wait()
}

// Wait blocks until all generator operations have completed and returns the collected data and a boolean indicating if the run failed.
func (g *Generator) Wait() (interface{}, bool) {
	g.Log.Info().Msg("Waiting for all responses to finish")
	g.ResponsesWaitGroup.Wait()
	g.stats.Duration = g.Cfg.duration.Nanoseconds()
	g.stats.CurrentTimeUnit = g.Cfg.RateLimitUnitDuration.Nanoseconds()
	g.dataCancel()
	g.dataWaitGroup.Wait()
	if g.Cfg.LokiConfig != nil {
		g.stopLokiStream()
	}
	return g.GetData(), g.stats.RunFailed.Load()
}

// InputSharedData retrieves the shared data from the generator's configuration.
// It allows access to common data shared across different components or processes.
func (g *Generator) InputSharedData() interface{} {
	return g.Cfg.SharedData
}

// Errors returns a slice of error messages collected by the Generator.
// Use this to access all errors encountered during the generation process.
func (g *Generator) Errors() []string {
	return g.errs.Data
}

// GetData retrieves the aggregated response data from the Generator.
// Use it to access all collected responses after processing is complete.
func (g *Generator) GetData() *ResponseData {
	return g.responsesData
}

// Stats returns the current statistics of the Generator.
// It allows callers to access and monitor the generator's state.
func (g *Generator) Stats() *Stats {
	return g.stats
}

/* Loki's methods to handle CallResult/Stats and stream it to Loki */

// stopLokiStream gracefully terminates the Loki streaming service if it is configured.
// It ensures that all Loki-related processes are properly stopped.
func (g *Generator) stopLokiStream() {
	if g.Cfg.LokiConfig != nil && g.Cfg.LokiConfig.URL != "" {
		g.Log.Info().Msg("Stopping Loki")
		g.loki.StopNow()
		g.Log.Info().Msg("Loki exited")
	}
}

// handleLokiResponsePayload enriches a Response with additional labels and submits it to Loki for centralized logging.
// It optimizes the payload by removing unnecessary timestamps and handles any errors that occur during submission.
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

// handleLokiStatsPayload transmits the generator’s current statistics to Loki for monitoring.
// It merges relevant labels with the stats data and handles any transmission errors by logging and stopping the generator.
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

// sendResponsesToLoki starts streaming response data to Loki using the generator's configuration.
// It handles incoming responses for monitoring and logging purposes.
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

// sendStatsToLoki starts a background goroutine that periodically sends generator statistics to Loki for monitoring.
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

// StatsJSON returns the generator's current statistics as a JSON-compatible map.
// It is used to capture and transmit real-time metrics for monitoring and analysis.
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

// printStatsLoop starts a background loop that periodically logs generator statistics.
// It runs until the generator's response context is canceled. Use it to monitor
// success, failure, and timeout metrics in real-time.
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

// LabelsMapToModel transforms a map of string key-value pairs into a model.LabelSet.
// This enables user-defined labels to be integrated into the model for downstream processing.
func LabelsMapToModel(m map[string]string) model.LabelSet {
	ls := model.LabelSet{}
	for k, v := range m {
		ls[model.LabelName(k)] = model.LabelValue(v)
	}
	return ls
}
