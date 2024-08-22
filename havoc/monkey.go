package havoc

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	MonkeyModeSeq    = "seq"
	MonkeyModeRandom = "rand"

	ErrInvalidMode = "monkey mode is invalid, should be either \"seq\" or \"rand\""
)

type ExperimentAction struct {
	Name           string
	ExperimentKind string
	ExperimentSpec string
	TimeStart      int64
	TimeEnd        int64
}

type ExperimentAnnotationBody struct {
	DashboardUID string   `json:"dashboardUID"`
	Time         int64    `json:"time"`
	TimeEnd      int64    `json:"timeEnd"`
	Tags         []string `json:"tags"`
	Text         string   `json:"text"`
}

type Controller struct {
	cfg               *Config
	client            *resty.Client
	ctx               context.Context
	cancel            context.CancelFunc
	wg                *sync.WaitGroup
	errors            []error
	experimentActions []*ExperimentAction
}

func NewController(cfg *Config) (*Controller, error) {
	InitDefaultLogging()
	if cfg == nil {
		cfg = DefaultConfig()
		dumpConfig(cfg)
	}
	c := resty.New()
	c.SetBaseURL(cfg.Havoc.Grafana.URL)
	c.SetAuthScheme("Bearer")
	c.SetAuthToken(cfg.Havoc.Grafana.Token)
	return &Controller{
		client:            c,
		cfg:               cfg,
		wg:                &sync.WaitGroup{},
		errors:            make([]error, 0),
		experimentActions: make([]*ExperimentAction, 0),
	}, nil
}

// AnnotateExperiment sends annotation marker to Grafana dashboard
func (m *Controller) AnnotateExperiment(a *ExperimentAction) error {
	if m.cfg.Havoc.Grafana.URL == "" || m.cfg.Havoc.Grafana.Token == "" {
		L.Warn().Msg("Dashboards are not selected, experiment time wasn't annotated, please check README to enable Grafana integration")
		return nil
	}
	for _, dashboardUID := range m.cfg.Havoc.Grafana.DashboardUIDs {
		start := a.TimeStart * 1e3
		end := a.TimeEnd * 1e3
		specBody := fmt.Sprintf("<pre>%s</pre>", a.ExperimentSpec)
		aa := &ExperimentAnnotationBody{
			DashboardUID: dashboardUID,
			Time:         start,
			TimeEnd:      end,
			Tags:         []string{"havoc", a.ExperimentKind},
			Text: fmt.Sprintf(
				"File: %s\n%s",
				a.Name,
				specBody,
			),
		}
		_, err := m.client.R().
			SetBody(aa).
			Post(fmt.Sprintf("%s/api/annotations", m.cfg.Havoc.Grafana.URL))
		if err != nil {
			return err
		}
		L.Info().
			Str("DashboardUID", dashboardUID).
			Str("Name", a.Name).
			Int64("Start", a.TimeStart).
			Int64("End", a.TimeEnd).
			Msg("Annotated experiment")
	}
	return nil
}

func (m *Controller) ApplyAndAnnotate(exp *NamedExperiment) error {
	ea := &ExperimentAction{
		Name:           exp.Name,
		ExperimentKind: exp.Kind,
		ExperimentSpec: string(exp.CRDBytes),
		TimeStart:      time.Now().Unix(),
	}
	if err := m.ApplyExperiment(exp, true); err != nil {
		return err
	}
	ea.TimeEnd = time.Now().Unix()
	return m.AnnotateExperiment(ea)
}

func (m *Controller) Run() error {
	L.Info().Msg("Starting chaos monkey")
	dur, err := time.ParseDuration(m.cfg.Havoc.Monkey.Duration)
	if err != nil {
		return err
	}
	m.ctx, m.cancel = context.WithTimeout(context.Background(), dur)
	defer m.cancel()
	existingExperimentTypes, err := m.readExistingExperimentTypes(m.cfg.Havoc.Dir)
	if err != nil {
		m.errors = append(m.errors, err)
		return err
	}

	m.wg.Add(1)
	switch m.cfg.Havoc.Monkey.Mode {
	case MonkeyModeSeq:
		for _, expType := range existingExperimentTypes {
			experiments, err := m.ReadExperimentsFromDir([]string{expType}, m.cfg.Havoc.Dir)
			if err != nil {
				m.errors = append(m.errors, err)
				return err
			}
			for _, exp := range experiments {
				if err := m.ApplyAndAnnotate(exp); err != nil {
					m.errors = append(m.errors, err)
					return err
				}
				cdDuration, err := time.ParseDuration(m.cfg.Havoc.Monkey.Cooldown)
				if err != nil {
					m.errors = append(m.errors, err)
					return err
				}
				select {
				case <-m.ctx.Done():
					m.wg.Done()
					L.Info().Msg("Monkey has finished by timeout")
					return nil
				default:
				}
				L.Info().
					Dur("Duration", cdDuration).
					Msg("Cooldown between experiments")
				time.Sleep(cdDuration)
			}
		}
		L.Info().Msg("Monkey has finished all scheduled experiments")
		m.wg.Done()
	case MonkeyModeRandom:
		allExperiments := make([]*NamedExperiment, 0)
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for _, expType := range existingExperimentTypes {
			experiments, err := m.ReadExperimentsFromDir([]string{expType}, m.cfg.Havoc.Dir)
			if err != nil {
				m.errors = append(m.errors, err)
				return err
			}
			allExperiments = append(allExperiments, experiments...)
		}
		for {
			select {
			case <-m.ctx.Done():
				m.wg.Done()
				L.Info().Msg("Monkey has finished by timeout")
				return nil
			default:
				exp := pickExperiment(r, allExperiments)
				if err := m.ApplyAndAnnotate(exp); err != nil {
					m.errors = append(m.errors, err)
					return err
				}
				cdDuration, err := time.ParseDuration(m.cfg.Havoc.Monkey.Cooldown)
				if err != nil {
					m.errors = append(m.errors, err)
					return err
				}
				L.Info().
					Dur("Duration", cdDuration).
					Msg("Cooldown between experiments")
				time.Sleep(cdDuration)
			}
		}
	default:
		return errors.New(ErrInvalidMode)
	}
	return nil
}

func (m *Controller) Stop() []error {
	L.Info().Msg("Stopping chaos monkey")
	m.cancel()
	m.wg.Wait()
	L.Info().Errs("Errors", m.errors).Msg("Chaos monkey stopped")
	return m.errors
}

func (m *Controller) Wait() []error {
	L.Info().Msg("Waiting for chaos monkey to finish")
	m.wg.Wait()
	L.Info().Errs("Errors", m.errors).Msg("Chaos monkey finished")
	return m.errors
}

func pickExperiment(r *rand.Rand, s []*NamedExperiment) *NamedExperiment {
	return s[r.Intn(len(s))]
}
