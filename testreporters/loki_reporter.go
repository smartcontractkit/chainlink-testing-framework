package testreporters

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Stats basic generator load stats
type Stats interface {
	ToJSON() map[string]interface{}
}

type LokiReporter struct {
	t                 *testing.T
	Log               zerolog.Logger
	dataCtx           context.Context
	dataCancel        context.CancelFunc
	LokiClient        client.ExtendedLokiClient
	ResponseChannel   chan loadgen.CallResult
	LokiCfg           *client.LokiConfig
	dataWaitGroup     *sync.WaitGroup
	Labels            map[string]string
	labels            model.LabelSet
	Summary           Stats
	StatsPollInterval time.Duration
}

func (l *LokiReporter) SetSummary(s Stats) {
	l.Summary = s
}

func (l *LokiReporter) Start() {
	l.dataWaitGroup.Add(1)
	go func() {
		defer l.dataWaitGroup.Done()
		for {
			select {
			case <-l.dataCtx.Done():
				return
			case r := <-l.ResponseChannel:
				l.handleLokiResponsePayload(r)
			}
		}
	}()
}

// handleLokiResponsePayload handles the response payload
func (l *LokiReporter) handleLokiResponsePayload(cr loadgen.CallResult) {
	ls := l.labels.Merge(model.LabelSet{
		"go_test_name":   model.LabelValue(l.t.Name()),
		"test_data_type": "responses",
	})
	// we are removing time.Time{} because when it marshalled to string it creates N responses for some Loki queries
	// and to minimize the payload, duration is already calculated at that point
	ts := cr.FinishedAt
	cr.StartedAt = nil
	cr.FinishedAt = nil
	err := l.LokiClient.HandleStruct(ls, *ts, cr)
	if err != nil {
		l.Log.Err(err).Send()
	}
}

// runLokiPromtailStats pushes Summary payloads to Loki
func (l *LokiReporter) runLokiPromtailStats() {
	l.dataWaitGroup.Add(1)
	go func() {
		defer l.dataWaitGroup.Done()
		for {
			select {
			case <-l.dataCtx.Done():
				l.Log.Info().Msg("Loki stats exited")
				return
			default:
				time.Sleep(l.StatsPollInterval)
				l.handleLokiStatsPayload()
			}
		}
	}()
}

// handleLokiStatsPayload handles the summary payload
func (l *LokiReporter) handleLokiStatsPayload() {
	ls := l.labels.Merge(model.LabelSet{
		"go_test_name":   model.LabelValue(l.t.Name()),
		"test_data_type": "stats",
	})
	err := l.LokiClient.HandleStruct(ls, time.Now(), l.Summary.ToJSON())
	if err != nil {
		l.Log.Err(err).Send()
	}
}

// Stop stops the Loki reporter
func (l *LokiReporter) Stop() {
	l.handleLokiStatsPayload()
	l.dataCancel()
	l.dataWaitGroup.Wait()
	l.stopLokiStream()
}

// stopLokiStream stops the Loki stream client
func (l *LokiReporter) stopLokiStream() {
	l.Log.Info().Msg("Stopping Loki")
	l.LokiClient.Stop()
	l.Log.Info().Msg("Loki exited")
}

// NewLokiReporter creates a new LokiReporter
func NewLokiReporter(t *testing.T, ctx context.Context, cfg *client.LokiConfig, labels map[string]string) (*LokiReporter, error) {
	c, err := client.NewLokiClient(cfg)
	if err != nil {
		return nil, err
	}
	// creating logger from *testing.T context or using a global logger
	var l zerolog.Logger
	if t != nil {
		l = zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t))).With().Timestamp().Logger()
	} else {
		l = log.Logger
	}

	dataCtx, dataCancel := context.WithCancel(ctx)
	return &LokiReporter{
		t:                 t,
		Log:               l,
		dataCtx:           dataCtx,
		dataCancel:        dataCancel,
		LokiClient:        c,
		ResponseChannel:   make(chan loadgen.CallResult),
		LokiCfg:           cfg,
		dataWaitGroup:     &sync.WaitGroup{},
		labels:            utils.LabelsMapToModel(labels),
		StatsPollInterval: 30 * time.Second,
	}, nil
}
