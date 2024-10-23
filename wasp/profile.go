package wasp

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/grafana"
)

// Profile is a set of concurrent generators forming some workload profile
type Profile struct {
	ProfileID    string // Unique identifier for the profile
	Generators   []*Generator
	testEndedWg  *sync.WaitGroup
	bootstrapErr error
	grafanaAPI   *grafana.Client
	grafanaOpts  GrafanaOpts
	startTime    time.Time
	endTime      time.Time
}

// Run runs all generators and wait until they finish
func (m *Profile) Run(wait bool) (*Profile, error) {
	if m.bootstrapErr != nil {
		return m, m.bootstrapErr
	}
	if err := waitSyncGroupReady(); err != nil {
		return m, err
	}
	m.startTime = time.Now()
	if len(m.grafanaOpts.AnnotateDashboardUID) > 0 {
		m.annotateRunStartOnGrafana()
	}
	for _, g := range m.Generators {
		g.Run(false)
	}
	if wait {
		m.Wait()
	}
	if m.grafanaOpts.WaitBeforeAlertCheck > 0 {
		log.Info().Msgf("Waiting %s before checking for alerts..", m.grafanaOpts.WaitBeforeAlertCheck)
		time.Sleep(m.grafanaOpts.WaitBeforeAlertCheck)
	}
	m.endTime = time.Now()
	if len(m.grafanaOpts.AnnotateDashboardUID) > 0 {
		m.annotateRunEndOnGrafana()
	}
	if m.grafanaOpts.CheckDashboardAlertsAfterRun != "" {
		m.printDashboardLink()
		alerts, err := CheckDashboardAlerts(m.grafanaAPI, m.startTime, time.Now(), m.grafanaOpts.CheckDashboardAlertsAfterRun)
		if len(alerts) > 0 {
			log.Info().Msgf("Alerts found\n%s", grafana.FormatAlertsTable(alerts))
		}
		if err != nil {
			return m, err
		}
	} else {
		m.printDashboardLink()
	}

	return m, nil
}

func (m *Profile) printDashboardLink() {
	if m.grafanaAPI == nil {
		log.Warn().Msg("Grafana API not set, skipping dashboard link print")
		return
	}
	d, _, err := m.grafanaAPI.GetDashboard(m.grafanaOpts.AnnotateDashboardUID)
	if err != nil {
		log.Warn().Msgf("could not get dashboard link: %s", err)
	}
	from := m.startTime.Add(-time.Second * 10).UnixMilli()
	to := m.endTime.Add(time.Second * 10).Add(m.grafanaOpts.WaitBeforeAlertCheck).UnixMilli()
	url := fmt.Sprintf("%s%s?from=%d&to=%d", m.grafanaOpts.GrafanaURL, d.Meta["url"].(string), from, to)

	if err != nil {
		log.Warn().Msgf("could not get dashboard link: %s", err)
	} else {
		log.Info().Msgf("Dashboard URL: %s", url)
	}
}

func (m *Profile) annotateRunStartOnGrafana() {
	if m.grafanaAPI == nil {
		log.Warn().Msg("Grafana API not set, skipping annotations")
		return
	}
	var sb strings.Builder
	sb.WriteString("<body>")
	sb.WriteString("<h4>Test Started</h4>")
	sb.WriteString(fmt.Sprintf("<div>WASP profileId: %s</div>", m.ProfileID))
	sb.WriteString(fmt.Sprintf("<div>Start time: %s</div>", m.startTime.Format(time.RFC3339)))
	sb.WriteString("<br>")
	sb.WriteString("<h5>Generators:</h5>")
	sb.WriteString("<ul>")
	for _, g := range m.Generators {
		sb.WriteString(fmt.Sprintf("<li>%s</li>", g.Cfg.GenName))
	}
	sb.WriteString("</ul>")
	sb.WriteString("</body>")

	a := grafana.PostAnnotation{
		DashboardUID: m.grafanaOpts.AnnotateDashboardUID,
		Time:         &m.startTime,
		Text:         sb.String(),
	}
	_, _, err := m.grafanaAPI.PostAnnotation(a)
	if err != nil {
		log.Warn().Msgf("could not annotate on Grafana: %s", err)
	}
}

func (m *Profile) annotateRunEndOnGrafana() {
	if m.grafanaAPI == nil {
		log.Warn().Msg("Grafana API not set, skipping annotations")
		return
	}
	var sb strings.Builder
	sb.WriteString("<body>")
	sb.WriteString("<h4>Test Ended</h4>")
	sb.WriteString(fmt.Sprintf("<div>WASP profileId: %s</div>", m.ProfileID))
	sb.WriteString(fmt.Sprintf("<div>End time: %s</div>", m.endTime.Format(time.RFC3339)))
	sb.WriteString("<br>")
	sb.WriteString("<h5>Generators:</h5>")
	sb.WriteString("<ul>")
	for _, g := range m.Generators {
		sb.WriteString(fmt.Sprintf("<li>%s</li>", g.Cfg.GenName))
	}
	sb.WriteString("</ul>")
	sb.WriteString("</body>")

	a := grafana.PostAnnotation{
		DashboardUID: m.grafanaOpts.AnnotateDashboardUID,
		Time:         &m.endTime,
		Text:         sb.String(),
	}
	_, _, err := m.grafanaAPI.PostAnnotation(a)
	if err != nil {
		log.Warn().Msgf("could not annotate on Grafana: %s", err)
	}
}

// Pause pauses execution of all generators
func (m *Profile) Pause() {
	for _, g := range m.Generators {
		g.Pause()
	}
}

// Resume resumes execution of all generators
func (m *Profile) Resume() {
	for _, g := range m.Generators {
		g.Resume()
	}
}

// Wait waits until all generators have finished the workload
func (m *Profile) Wait() {
	for _, g := range m.Generators {
		g := g
		m.testEndedWg.Add(1)
		go func() {
			defer m.testEndedWg.Done()
			g.Wait()
		}()
	}
	m.testEndedWg.Wait()
}

// NewProfile creates new VU or Gun profile from parts
func NewProfile() *Profile {
	return &Profile{
		ProfileID:   uuid.NewString()[0:5],
		Generators:  make([]*Generator, 0),
		testEndedWg: &sync.WaitGroup{},
	}
}

func (m *Profile) Add(g *Generator, err error) *Profile {
	if err != nil {
		m.bootstrapErr = err
		return m
	}
	m.Generators = append(m.Generators, g)
	return m
}

type GrafanaOpts struct {
	GrafanaURL                   string        `toml:"grafana_url"`
	GrafanaToken                 string        `toml:"grafana_token_secret"`
	WaitBeforeAlertCheck         time.Duration `toml:"grafana_wait_before_alert_check"`                 // Cooldown period to wait before checking for alerts
	AnnotateDashboardUID         string        `toml:"grafana_annotate_dashboard_uid"`                  // Grafana dashboardUID to annotate start and end of the run
	CheckDashboardAlertsAfterRun string        `toml:"grafana_check_alerts_after_run_on_dashboard_uid"` // Grafana dashboardUID to check for alerts after run
}

func (m *Profile) WithGrafana(opts *GrafanaOpts) *Profile {
	m.grafanaAPI = grafana.NewGrafanaClient(opts.GrafanaURL, opts.GrafanaToken)
	m.grafanaOpts = *opts
	return m
}

// waitSyncGroupReady awaits other pods with WASP_SYNC label to start before starting the test
func waitSyncGroupReady() error {
	if os.Getenv("WASP_NODE_ID") != "" {
		kc := NewK8sClient()
		jobNum, err := strconv.Atoi(os.Getenv("WASP_JOBS"))
		if err != nil {
			return err
		}
		if err := kc.waitSyncGroup(context.Background(), os.Getenv("WASP_NAMESPACE"), os.Getenv("WASP_SYNC"), jobNum); err != nil {
			return err
		}
	}
	return nil
}
