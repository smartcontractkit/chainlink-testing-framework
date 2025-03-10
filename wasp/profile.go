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

// Run executes the profile's generators, manages Grafana annotations, and handles alert checks.
// If wait is true, it waits for all generators to complete before proceeding.
// It returns the updated Profile and any encountered error.
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
	m.printProfileId()
	m.printDashboardLink()
	if m.grafanaOpts.CheckDashboardAlertsAfterRun != "" {
		alerts, err := CheckDashboardAlerts(m.grafanaAPI, m.startTime, time.Now(), m.grafanaOpts.CheckDashboardAlertsAfterRun)
		if len(alerts) > 0 {
			log.Info().Msgf("Alerts found\n%s", grafana.FormatAlertsTable(alerts))
		}
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func (m *Profile) printProfileId() {
	log.Info().Msgf("Profile ID: %s", m.ProfileID)
}

// printDashboardLink retrieves the Grafana dashboard URL for the current run
// and logs it. It provides users with a direct link to monitor metrics and alerts
// related to the profile execution.
func (m *Profile) printDashboardLink() {
	if m.grafanaAPI == nil {
		log.Warn().Msg("Grafana API not set, skipping dashboard link print")
		return
	}
	d, _, err := m.grafanaAPI.GetDashboard(m.grafanaOpts.AnnotateDashboardUID)
	if err != nil {
		log.Warn().Msgf("could not get dashboard link: %s", err)
		return
	}
	if d.Meta == nil || d.Meta["ur"] == nil {
		log.Warn().Msgf("nil dasbhoard metadata returned from Grafana API with uid %s", m.grafanaOpts.AnnotateDashboardUID)
		return
	}
	from := m.startTime.Add(-time.Second * 10).UnixMilli()
	to := m.endTime.Add(time.Second * 10).Add(m.grafanaOpts.WaitBeforeAlertCheck).UnixMilli()
	url := fmt.Sprintf("%s%s?from=%d&to=%d", m.grafanaOpts.GrafanaURL, d.Meta["url"].(string), from, to)

	log.Info().Msgf("Dashboard URL: %s", url)
}

// annotateRunStartOnGrafana posts a run start annotation to the Grafana dashboard.
// It includes run details for monitoring and tracking purposes.
// Logs a warning if the Grafana API is not configured.
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

// annotateRunEndOnGrafana creates and posts an end-of-run annotation to Grafana,
// including profile ID, end time, and generator details.
// It is used to mark the completion of a profile run on the Grafana dashboard.
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

// Pause suspends all generators within the profile.
// It is used to temporarily halt all generator operations managed by the profile.
func (m *Profile) Pause() {
	for _, g := range m.Generators {
		g.Pause()
	}
}

// Resume resumes all generators associated with the profile, allowing them to continue their operations.
func (m *Profile) Resume() {
	for _, g := range m.Generators {
		g.Resume()
	}
}

// Wait blocks until all generators associated with the Profile have finished executing,
// ensuring all operations are complete before proceeding.
func (m *Profile) Wait() {
	for _, g := range m.Generators {

		m.testEndedWg.Add(1)
		go func() {
			defer m.testEndedWg.Done()
			g.Wait()
		}()
	}
	m.testEndedWg.Wait()
}

// NewProfile creates and returns a new Profile instance.
// It initializes the ProfileID with a unique identifier,
// an empty slice of Generators, and a WaitGroup for synchronization.
// Use it to instantiate profiles with default settings.
func NewProfile() *Profile {
	return &Profile{
		ProfileID:   uuid.NewString()[0:5],
		Generators:  make([]*Generator, 0),
		testEndedWg: &sync.WaitGroup{},
	}
}

// Add appends a Generator to the Profile. If an error is provided, it records the bootstrap error and does not add the Generator.
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

// WithGrafana configures the Profile with Grafana settings.
// It initializes the Grafana client using the provided options
// and returns the updated Profile instance.
func (m *Profile) WithGrafana(opts *GrafanaOpts) *Profile {
	m.grafanaAPI = grafana.NewGrafanaClient(opts.GrafanaURL, opts.GrafanaToken)
	m.grafanaOpts = *opts
	return m
}

// waitSyncGroupReady waits for the synchronization group to be ready based on environment variables.
// It ensures dependencies are initialized before proceeding with execution.
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
