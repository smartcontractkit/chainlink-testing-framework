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

// Run starts the Profile execution, handling bootstrap errors and ensuring synchronization readiness.
// It records the start time, optionally annotates the run start on Grafana, and executes all associated generators.
// If the wait parameter is true, it waits for all generators to complete. Before alert checking,
// it may pause for a configured duration. Upon completion, it records the end time, optionally annotates
// the run end on Grafana, checks for dashboard alerts if configured, and returns the updated Profile
// or any encountered error.
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

// printDashboardLink retrieves the Grafana dashboard URL for the current profile run and logs it.  
// If the Grafana API is not configured or an error occurs while fetching the dashboard,  
// appropriate warnings are logged instead. This function is typically called after a profiling  
// run to provide a direct link to the associated Grafana dashboard for analysis.
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

// annotateRunStartOnGrafana creates and posts a "Test Started" annotation to Grafana with profile details.
// It includes the profile ID, start time, and a list of generators. If the Grafana API
// is not configured, it logs a warning and skips the annotation. Any errors encountered
// while posting the annotation are also logged.
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

// annotateRunEndOnGrafana records the completion of a profile run by posting an annotation to Grafana.
// It includes the profile ID, end time, and a list of generators used.
// If the Grafana API is not configured or the annotation fails, a warning is logged.
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

// Pause pauses all generators associated with the profile.
// It iterates through each generator in the Profile and invokes their Pause method,
// which logs a warning and updates their paused status.
func (m *Profile) Pause() {
	for _, g := range m.Generators {
		g.Pause()
	}
}

// Resume resumes all generators associated with the Profile.
func (m *Profile) Resume() {
	for _, g := range m.Generators {
		g.Resume()
	}
}

// Wait blocks until all generators associated with the Profile have completed their execution.
// It launches a goroutine for each generator's Wait method and waits for all to finish.
// This ensures that the Profile only proceeds once all generator processes are done.
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

// NewProfile creates and returns a new Profile with a unique ProfileID, an empty slice of Generators, and an initialized sync.WaitGroup for managing test completion.
func NewProfile() *Profile {
	return &Profile{
		ProfileID:   uuid.NewString()[0:5],
		Generators:  make([]*Generator, 0),
		testEndedWg: &sync.WaitGroup{},
	}
}

// Add appends the provided Generator to the Profile's Generators slice.
// If an error is supplied, it sets the Profile's bootstrapErr field instead.
// It returns the updated Profile.
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

// WithGrafana initializes the Profile with a Grafana client using the provided GrafanaOpts.
// It sets the grafanaAPI and grafanaOpts fields and returns the updated Profile.
func (m *Profile) WithGrafana(opts *GrafanaOpts) *Profile {
	m.grafanaAPI = grafana.NewGrafanaClient(opts.GrafanaURL, opts.GrafanaToken)
	m.grafanaOpts = *opts
	return m
}

// waitSyncGroupReady waits for the synchronization group to be ready.
// It checks the required environment variables and ensures the synchronization
// process completes successfully. An error is returned if the synchronization fails.
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
