package grafana

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-resty/resty/v2"
)

type GrafanaClient struct {
	AlertManager *AlertManagerClient
	resty        *resty.Client
}

type AlertManagerClient struct {
	resty *resty.Client
}

func (g *AlertManagerClient) GetAlertGroups() ([]AlertGroupsResponse, *resty.Response, error) {
	var result []AlertGroupsResponse
	r, err := g.resty.R().SetResult(&result).Get("/api/alertmanager/grafana/api/v2/alerts/groups")
	return result, r, err
}

func (g *AlertManagerClient) GetAlterManagerAlerts() ([]interface{}, *resty.Response, error) {
	var result []interface{}
	r, err := g.resty.R().SetResult(&result).Get("/api/alertmanager/grafana/api/v2/alerts")
	return result, r, err
}

func NewGrafanaClient(url, apiKey string) *GrafanaClient {
	isDebug := os.Getenv("DEBUG_RESTY") == "true"
	resty := resty.New().SetDebug(isDebug).SetBaseURL(url).SetHeader("Authorization", "Bearer "+apiKey)

	return &GrafanaClient{
		resty:        resty,
		AlertManager: &AlertManagerClient{resty: resty},
	}
}

func (g *GrafanaClient) GetAlertsRules() ([]ProvisionedAlertRule, *resty.Response, error) {
	var result []ProvisionedAlertRule
	r, err := g.resty.R().SetResult(&result).Get("/api/v1/provisioning/alert-rules")
	return result, r, err
}

func (g *GrafanaClient) GetAlertRulesForDashboardID(dashboardID string) ([]ProvisionedAlertRule, error) {
	rules, _, err := g.GetAlertsRules()
	if err != nil {
		return nil, err
	}
	var results []ProvisionedAlertRule
	for _, rule := range rules {
		ruleDashboardUID := rule.Annotations["__dashboardUid__"]
		if ruleDashboardUID == dashboardID {
			results = append(results, rule)
		}
	}
	return results, nil
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Convert bytes to string and parse the number
	var ms int64
	err := json.Unmarshal(b, &ms)
	if err != nil {
		return err
	}
	// Set the time using milliseconds
	ct.Time = time.Unix(0, ms*int64(time.Millisecond))
	return nil
}

type Annotation struct {
	ID           int64         `json:"id"`
	AlertID      int64         `json:"alertId"`
	DashboardID  int64         `json:"dashboardId"`
	DashboardUID string        `json:"dashboardUID"`
	PanelID      int64         `json:"panelId"`
	PrevState    string        `json:"prevState"`
	NewState     string        `json:"newState"`
	Text         string        `json:"text"`
	Time         CustomTime    `json:"time"`
	TimeEnd      CustomTime    `json:"timeEnd"`
	Created      CustomTime    `json:"created"`
	Updated      CustomTime    `json:"updated"`
	Tags         []interface{} `json:"tags"`
	Data         interface{}   `json:"data"`
}

func FormatAlertsTable(alerts []Annotation) string {
	var b strings.Builder
	// Initialize a new tab writer with output directed to b, a min width of 8, a tab width of 8, a padding of 2, and tabs replaced by spaces
	w := tabwriter.NewWriter(&b, 8, 8, 2, ' ', 0)

	// Write the table header
	fmt.Fprintln(w, "Time\tPrevState\tNewState\tAlertID\tAlert")

	// Write each alert as a row in the table
	for _, a := range alerts {
		fmt.Fprintf(w,
			"%s\t%s\t%s\t%d\t%s\n",
			a.Time.UTC().Format(time.RFC3339), a.PrevState, a.NewState, a.ID, a.Text,
		)
	}

	// Flush writes to the underlying writer
	w.Flush()

	return b.String()
}

type AnnotationsQueryParams struct {
	Limit        *int
	AlertID      *int
	DashboardID  *int
	DashboardUID *string // when dashboardUID presents, dashboardId would be ignored by /annotations API
	Type         *string
	From         *time.Time
	To           *time.Time
}

type PostAnnotation struct {
	DashboardID  *int
	PanelID      *int
	DashboardUID string
	Time         *time.Time
	TimeEnd      *time.Time
	Tags         []string
	Text         string
}

func (g *GrafanaClient) GetAnnotations(params AnnotationsQueryParams) ([]Annotation, *resty.Response, error) {
	query := make(url.Values)
	if params.Limit != nil {
		query.Set("limit", fmt.Sprintf("%d", *params.Limit))
	}
	if params.AlertID != nil {
		query.Set("alertId", fmt.Sprintf("%d", *params.AlertID))
	}
	if params.DashboardID != nil {
		query.Set("dashboardId", fmt.Sprintf("%d", *params.DashboardID))
	}
	if params.DashboardUID != nil {
		query.Set("dashboardUID", *params.DashboardUID)
	}
	// type can be: alert or annotation
	if params.Type != nil {
		query.Set("type", *params.Type)
	}

	// Grafana issue https://github.com/grafana/grafana/issues/63130
	if (params.From != nil && params.To == nil) || (params.To != nil && params.From == nil) {
		return nil, nil, fmt.Errorf("both From and To must be set")
	}

	if params.From != nil {
		query.Set("from", fmt.Sprintf("%d", params.From.UnixMilli()))
	}
	if params.To != nil {
		query.Set("to", fmt.Sprintf("%d", params.To.UnixMilli()))
	}

	var result []Annotation
	r, err := g.resty.R().
		SetResult(&result).
		SetQueryString(query.Encode()).
		Get("/api/annotations")
	return result, r, err
}

func (g *GrafanaClient) PostAnnotation(annotation PostAnnotation) (*resty.Response, error) {
	a := map[string]interface{}{
		"dashboardUID": annotation.DashboardUID,
		"tags":         annotation.Tags,
		"text":         annotation.Text,
	}
	if annotation.DashboardID != nil {
		a["dashboardId"] = *annotation.DashboardID
	}
	if annotation.PanelID != nil {
		a["panelId"] = *annotation.PanelID
	}
	if annotation.Time != nil {
		a["time"] = annotation.Time.UnixMilli()
	}
	if annotation.TimeEnd != nil {
		a["timeEnd"] = annotation.TimeEnd.UnixMilli()
	}
	return g.resty.R().
		SetBody(a).
		Post("/api/annotations")
}

// ruler API is deprecated https://github.com/grafana/grafana/issues/74434
func (g *GrafanaClient) GetAlertsForDashboard(dashboardUID string) (map[string][]interface{}, *resty.Response, error) {
	var result map[string][]interface{}
	r, err := g.resty.R().SetResult(&result).Get("/api/ruler/grafana/api/v1/rules?dashboard_uid=" + url.QueryEscape(dashboardUID))
	return result, r, err
}

type AlertQuery interface{}

type ProvisionedAlertRule struct {
	ID           int64             `json:"id"`
	UID          string            `json:"uid"`
	FolderUID    string            `json:"folderUID"`
	Title        string            `json:"title"`
	Data         []AlertQuery      `json:"data"`
	ExecErrState string            `json:"execErrState"`
	Labels       map[string]string `json:"labels"`
	RuleGroup    string            `json:"ruleGroup"`
	UpdatedAt    time.Time         `json:"updated"`
	Annotations  map[string]string `json:"annotations"`
}

func (p ProvisionedAlertRule) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ID: %d\n", p.ID))
	sb.WriteString(fmt.Sprintf("UID: %s\n", p.UID))
	sb.WriteString(fmt.Sprintf("FolderUID: %s\n", p.FolderUID))
	sb.WriteString(fmt.Sprintf("Title: %s\n", p.Title))
	// sb.WriteString("Data: [\n")
	// for _, query := range p.Data {
	// 	sb.WriteString(fmt.Sprintf("  %+v\n", query))
	// }
	sb.WriteString("]\n")
	sb.WriteString(fmt.Sprintf("ExecErrState: %s\n", p.ExecErrState))
	sb.WriteString("Labels: {\n")
	for key, value := range p.Labels {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
	}
	sb.WriteString("}\n")
	sb.WriteString(fmt.Sprintf("RuleGroup: %s\n", p.RuleGroup))
	sb.WriteString(fmt.Sprintf("UpdatedAt: %s\n", p.UpdatedAt.Format(time.RFC3339)))
	sb.WriteString("Annotations: {\n")
	for key, value := range p.Annotations {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
	}
	sb.WriteString("}\n")
	return sb.String()
}

// AlertGroupsResponse is response body for "api/alertmanager/grafana/api/v2/alerts/groups"
type AlertGroupsResponse struct {
	Alerts []Alert `json:"alerts"`
	Labels struct {
		Alertname     string `json:"alertname"`
		GrafanaFolder string `json:"grafana_folder"`
	} `json:"labels"`
	Receiver struct {
		Active       interface{} `json:"active"`
		Integrations interface{} `json:"integrations"`
		Name         string      `json:"name"`
	} `json:"receiver"`
}

type Alert struct {
	Annotations struct {
		DashboardUID string `json:"__dashboardUid__"`
		OrgID        string `json:"__orgId__"`
		PanelID      string `json:"__panelId__"`
		Description  string `json:"description"`
		RunbookURL   string `json:"runbook_url"`
		Summary      string `json:"summary"`
	} `json:"annotations"`
	EndsAt      time.Time `json:"endsAt"`
	Fingerprint string    `json:"fingerprint"`
	Receivers   []struct {
		Active       interface{} `json:"active"`
		Integrations interface{} `json:"integrations"`
		Name         string      `json:"name"`
	} `json:"receivers"`
	StartsAt time.Time `json:"startsAt"`
	Status   struct {
		InhibitedBy []interface{} `json:"inhibitedBy"`
		SilencedBy  []interface{} `json:"silencedBy"`
		State       string        `json:"state"`
	} `json:"status"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
}
