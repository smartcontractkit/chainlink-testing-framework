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

// NewAlertChecker creates a new AlertChecker using Grafana configurations from environment variables.
// It retrieves GRAFANA_URL and GRAFANA_TOKEN, ensuring they are set.
// Use this function to set up alert checking in tests.
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

// AnyAlerts retrieves alert groups from Grafana and checks for alerts matching the specified dashboard UUID and requirement label value.
// It returns the matching alert groups, enabling users to identify and respond to specific alert conditions.
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

// CheckDashboardAlerts retrieves alert annotations from a Grafana dashboard within the specified time range.
// It returns the sorted alerts and an error if any alert is in the alerting state.
// Use it to verify the status of dashboard alerts after a run.
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
