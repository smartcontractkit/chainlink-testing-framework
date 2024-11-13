package wasp

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/grafana"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// AlertChecker is checking alerts according to dashboardUUID and requirements labels
type AlertChecker struct {
	RequirementLabelKey string
	T                   *testing.T
	l                   zerolog.Logger
	grafanaClient       *grafana.Client
}

// NewAlertChecker initializes a new AlertChecker instance by retrieving the Grafana URL and API token from the environment variables. 
// It panics if either the GRAFANA_URL or GRAFANA_TOKEN environment variable is not set, ensuring that the necessary configuration is available. 
// The function creates a Grafana client using the provided URL and API token, and sets up the AlertChecker with a default requirement label key and a logger. 
// It returns a pointer to the newly created AlertChecker instance, which can be used for checking alerts in Grafana.
func NewAlertChecker(t *testing.T) *AlertChecker {
	url := os.Getenv("GRAFANA_URL")
	if url == "" {
		panic(fmt.Errorf("GRAFANA_URL env var must be defined"))
	}
	apiKey := os.Getenv("GRAFANA_TOKEN")
	if apiKey == "" {
		panic(fmt.Errorf("GRAFANA_TOKEN env var must be defined"))
	}

	grafanaClient := grafana.NewGrafanaClient(url, apiKey)

	return &AlertChecker{
		RequirementLabelKey: "requirement_name",
		T:                   t,
		grafanaClient:       grafanaClient,
		l:                   GetLogger(t, "AlertChecker"),
	}
}

// AnyAlerts checks for alerts associated with a specific dashboard and requirement label. 
// It retrieves alert groups from the Grafana Alert Manager and scans through the alerts to determine 
// if any alert matches the provided dashboard UUID and requirement label value. 
// If an alert is found, it logs the alert details and marks that an alert has been raised. 
// The function returns the list of alert groups and any error encountered during the retrieval process. 
// If the test context is active and an alert was raised, it will mark the test as failed.
func (m *AlertChecker) AnyAlerts(dashboardUUID, requirementLabelValue string) ([]grafana.AlertGroupsResponse, error) {
	raised := false
	defer func() {
		if m.T != nil && raised {
			m.T.Fail()
		}
	}()
	alertGroups, _, err := m.grafanaClient.AlertManager.GetAlertGroups()
	if err != nil {
		return alertGroups, fmt.Errorf("failed to get alert groups: %s", err)
	}
	for _, a := range alertGroups {
		for _, aa := range a.Alerts {
			log.Debug().Interface("Alert", aa).Msg("Scanning alert")
			if aa.Annotations.DashboardUID == dashboardUUID && aa.Labels[m.RequirementLabelKey] == requirementLabelValue {
				log.Warn().
					Str("Summary", aa.Annotations.Summary).
					Str("Description", aa.Annotations.Description).
					Str("URL", aa.GeneratorURL).
					Interface("Labels", aa.Labels).
					Time("StartsAt", aa.StartsAt).
					Time("UpdatedAt", aa.UpdatedAt).
					Str("State", aa.Status.State).
					Msg("Alert fired")
				raised = true
			}
		}
	}
	return alertGroups, nil
}

// CheckDashboardAlerts retrieves annotations of type "alert" from a specified Grafana dashboard within a given time range. 
// It returns a slice of annotations and an error if any occurred during the retrieval process. 
// If the retrieval is successful, the function sorts the annotations by time from oldest to newest and checks if any alerts are in an "alerting" state. 
// If at least one alert is found to be firing, it returns the annotations along with an error indicating that an alert was firing. 
// If no alerts are firing, it returns the annotations with a nil error.
func CheckDashboardAlerts(grafanaClient *grafana.Client, from, to time.Time, dashboardUID string) ([]grafana.Annotation, error) {
	annotationType := "alert"
	alerts, _, err := grafanaClient.GetAnnotations(grafana.AnnotationsQueryParams{
		DashboardUID: &dashboardUID,
		From:         &from,
		To:           &to,
		Type:         &annotationType,
	})
	if err != nil {
		return alerts, fmt.Errorf("could not check for alerts: %s", err)
	}

	// Sort the annotations by time oldest to newest
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].Time.Before(alerts[j].Time.Time)
	})

	// Check if any alerts are in alerting state
	for _, a := range alerts {
		if strings.ToLower(a.NewState) == "alerting" {
			return alerts, errors.New("at least one alert was firing")
		}
	}

	return alerts, nil
}
