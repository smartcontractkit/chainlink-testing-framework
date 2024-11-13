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

// Run executes the profile, optionally waiting for completion. It returns the updated Profile and any error encountered.
// If a bootstrap error exists, it returns immediately. It waits for the sync group to be ready before starting.
// The start time is recorded, and if Grafana annotations are enabled, it annotates the start.
// It runs all generators and waits if specified. After execution, it may wait before checking alerts.
// The end time is recorded, and if Grafana annotations are enabled, it annotates the end.
// It checks for alerts on the Grafana dashboard if configured, printing any found alerts.
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

// printDashboardLink retrieves and logs the Grafana dashboard URL for the profile.
// It constructs the URL using the start and end times of the profile run, and logs
// a warning if the Grafana API is not set or if there is an error retrieving the dashboard.
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

// annotateRunStartOnGrafana posts an annotation to Grafana indicating the start of a test run.
// It includes details such as the profile ID, start time, and generator names.
// If the Grafana API is not set, it logs a warning and skips the annotation.
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

// annotateRunEndOnGrafana posts an annotation to Grafana indicating the end of a test run.
// It includes details such as the profile ID, end time, and a list of generators.
// If the Grafana API is not set, it logs a warning and skips the annotation.
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
// It iterates over each generator and invokes its Pause method,
// which logs a warning message and updates the generator's state to paused.
func (m *Profile) Pause() {
	for _, g := range m.Generators {
		g.Pause()
	}
}

// Resume resumes all generators in the profile by invoking their Resume method.
func (m *Profile) Resume() {
	for _, g := range m.Generators {
		g.Resume()
	}
}

// Wait blocks until all generator routines in the Profile have completed.
// It ensures that all generators have finished their execution by waiting
// on a synchronization mechanism, allowing for concurrent operations to
// complete before proceeding. This function is typically used after initiating
// generator operations to ensure all tasks are finalized.
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

// NewProfile creates and returns a new Profile instance with a unique ProfileID.
// The ProfileID is a 5-character string generated from a UUID. It initializes
// an empty slice of Generators and a sync.WaitGroup for managing concurrent tasks.
func NewProfile() *Profile {
	return &Profile{
		ProfileID:   uuid.NewString()[0:5],
		Generators:  make([]*Generator, 0),
		testEndedWg: &sync.WaitGroup{},
	}
}

// Add appends a Generator to the Profile's Generators slice if no error is encountered.
// If an error is present, it sets the Profile's bootstrapErr to the error and returns the Profile.
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

// WithGrafana configures the Profile with Grafana settings using the provided options.
// It initializes a new Grafana client with the specified URL and token from opts.
// The function returns the updated Profile instance.
func (m *Profile) WithGrafana(opts *GrafanaOpts) *Profile {
	m.grafanaAPI = grafana.NewGrafanaClient(opts.GrafanaURL, opts.GrafanaToken)
	m.grafanaOpts = *opts
	return m
}

// waitSyncGroupReady checks if the environment variable "WASP_NODE_ID" is set.
// If set, it initializes a Kubernetes client and waits for a synchronization group
// to be ready based on the specified namespace, sync group, and job number.
// It returns any error encountered during this process.
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
