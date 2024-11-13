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

// Run executes the profile's generators and manages the run lifecycle. 
// If the wait parameter is true, it will block until all generators have completed. 
// It also handles Grafana annotations for the run start and end, and checks for alerts 
// on the specified dashboard after the run if configured. 
// The function returns the profile instance and any error encountered during execution. 
// If there is a bootstrap error or an error while waiting for the sync group, 
// it will return the profile and the corresponding error.
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

// printDashboardLink retrieves the dashboard link from the Grafana API and logs it. 
// If the Grafana API is not set or an error occurs while fetching the dashboard, 
// a warning message is logged. The function constructs a URL using the dashboard's 
// metadata and the specified time range, which is adjusted by a few seconds 
// before and after the start and end times. If successful, the dashboard URL is 
// logged as an info message.
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

// annotateRunStartOnGrafana sends an annotation to Grafana indicating the start of a test run. 
// It constructs an HTML message containing the profile ID, start time, and a list of generators involved in the run. 
// If the Grafana API is not set, it logs a warning and skips the annotation. 
// In case of an error during the annotation posting, it logs a warning with the error details.
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

// annotateRunEndOnGrafana posts an annotation to a Grafana dashboard indicating the end of a test run. 
// It constructs an HTML message containing the profile ID, end time, and a list of generators used during the run. 
// If the Grafana API is not set, it logs a warning and skips the annotation. 
// In case of an error while posting the annotation, it logs a warning with the error details.
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
// It iterates through each generator in the profile's generator list and calls their Pause method. 
// This action will trigger a warning log indicating that the generator has been paused and update the generator's state to reflect that it is currently paused.
func (m *Profile) Pause() {
	for _, g := range m.Generators {
		g.Pause()
	}
}

// Resume resumes all generators associated with the profile. 
// It iterates through each generator in the profile's generator list and calls their Resume method, 
// which logs a warning message indicating that the generator has been resumed and updates its run state to indicate that it is no longer paused.
func (m *Profile) Resume() {
	for _, g := range m.Generators {
		g.Resume()
	}
}

// Wait blocks until all generator processes associated with the profile have completed. 
// It initiates a wait for each generator in the profile's Generators slice, 
// ensuring that the function does not return until all concurrent operations are finished. 
// This is particularly useful when the profile is run with the wait parameter set to true, 
// allowing for synchronization of the completion of all tasks before proceeding.
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

// NewProfile creates and returns a new instance of Profile. 
// It initializes the ProfileID with a unique identifier, 
// sets up an empty slice for Generators, and prepares a 
// WaitGroup to manage synchronization for test completion. 
// The returned Profile is ready for further configuration 
// and use in the application.
func NewProfile() *Profile {
	return &Profile{
		ProfileID:   uuid.NewString()[0:5],
		Generators:  make([]*Generator, 0),
		testEndedWg: &sync.WaitGroup{},
	}
}

// Add appends a new Generator to the Profile's list of Generators. 
// If the provided error is not nil, it sets the bootstrapErr field 
// of the Profile to the given error and returns the Profile without 
// modifying the Generators list. If there is no error, the function 
// adds the Generator to the list and returns the updated Profile.
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

// WithGrafana initializes the Grafana client with the provided options and updates the profile's Grafana settings. 
// It takes a pointer to GrafanaOpts, which includes the Grafana URL and token. 
// The function returns the updated Profile instance, allowing for method chaining.
func (m *Profile) WithGrafana(opts *GrafanaOpts) *Profile {
	m.grafanaAPI = grafana.NewGrafanaClient(opts.GrafanaURL, opts.GrafanaToken)
	m.grafanaOpts = *opts
	return m
}

// waitSyncGroupReady checks if the synchronization group is ready by verifying the environment variable "WASP_NODE_ID". 
// If the variable is set, it creates a Kubernetes client and retrieves the number of jobs from the "WASP_JOBS" environment variable. 
// It then waits for the specified synchronization group in the given namespace to be ready. 
// If any errors occur during this process, they are returned. 
// If the synchronization group is ready or if "WASP_NODE_ID" is not set, it returns nil.
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
