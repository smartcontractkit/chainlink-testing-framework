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

// NewAlertChecker initializes and returns a new AlertChecker instance.
// It retrieves the Grafana URL and API token from environment variables
// GRAFANA_URL and GRAFANA_TOKEN, respectively. If either is not set, it
// panics. The function creates a Grafana client using these credentials
// and sets up the AlertChecker with a logger and a requirement label key.
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

// AnyAlerts checks for alerts in Grafana that match the specified dashboard UUID and requirement label value.
// It returns a slice of AlertGroupsResponse and any error encountered during the retrieval of alert groups.
// If an alert is found, it logs the alert details and marks the test as failed if a testing object is present.
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

// CheckDashboardAlerts retrieves and checks alerts from a Grafana dashboard within a specified time range.
// It returns a slice of annotations representing the alerts and an error if any alert is in the "alerting" state
// or if there is an issue retrieving the annotations. The alerts are sorted by time from oldest to newest.
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
