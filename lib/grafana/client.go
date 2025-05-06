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
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

type Client struct {
	AlertManager *AlertManagerClient
	AlertRuler   *AlertRulerClient
	resty        *resty.Client
}

type AlertManagerClient struct {
	resty *resty.Client
}

type AlertRulerClient struct {
	resty *resty.Client
}

// GetAlertGroups retrieves the alert groups from the AlertManager API.
// It returns a slice of AlertGroupsResponse, the HTTP response, and any error encountered.
// This function is useful for monitoring and managing alert configurations.
func (g *AlertManagerClient) GetAlertGroups() ([]AlertGroupsResponse, *resty.Response, error) {
	var result []AlertGroupsResponse
	r, err := g.resty.R().SetResult(&result).Get("/api/alertmanager/grafana/api/v2/alerts/groups")
	return result, r, err
}

// GetAlterManagerAlerts retrieves alerts from the AlertManager API.
// It returns a slice of alerts, the HTTP response, and any error encountered.
func (g *AlertManagerClient) GetAlterManagerAlerts() ([]interface{}, *resty.Response, error) {
	var result []interface{}
	r, err := g.resty.R().SetResult(&result).Get("/api/alertmanager/grafana/api/v2/alerts")
	return result, r, err
}

// GetDatasources retrieves a map of datasource names to their unique identifiers from the API.
// It also identifies the default datasource, if available, and includes it in the returned map.
// This function is useful for applications needing to interact with various datasources dynamically.
func (c *Client) GetDatasources() (map[string]string, *resty.Response, error) {
	var result []struct {
		UID       string `json:"uid"`
		Name      string `json:"name"`
		IsDefault bool   `json:"isDefault"`
	}

	r, err := c.resty.R().SetResult(&result).Get("/api/datasources")
	if err != nil {
		return nil, r, fmt.Errorf("error making API request: %w", err)
	}

	datasourcesMap := make(map[string]string, len(result))
	for _, ds := range result {
		datasourcesMap[ds.Name] = ds.UID

		if ds.IsDefault {
			datasourcesMap["$grabana_default_datasource_key$"] = ds.UID
		}
	}

	return datasourcesMap, r, err
}

// NewGrafanaClient initializes a new Grafana client with the specified URL and API key.
// It sets up the necessary headers and configurations for making API requests to Grafana,
// providing access to AlertManager and AlertRuler functionalities.
func NewGrafanaClient(url, apiKey string) *Client {
	isDebug := os.Getenv("RESTY_DEBUG") == "true"
	resty := resty.New().
		SetDebug(isDebug).
		SetBaseURL(url).
		SetHeader("Authorization", "Bearer "+apiKey)
	return &Client{
		resty:        resty,
		AlertManager: &AlertManagerClient{resty: resty},
		AlertRuler:   &AlertRulerClient{resty: resty},
	}
}

type GetDashboardResponse struct {
	Meta      map[string]interface{} `json:"meta"`
	Dashboard *dashboard.Dashboard   `json:"dashboard"`
}

// GetDashboard retrieves the dashboard associated with the given unique identifier (uid).
// It returns the dashboard data, the HTTP response, and any error encountered during the request.
func (c *Client) GetDashboard(uid string) (GetDashboardResponse, *resty.Response, error) {
	var result GetDashboardResponse
	r, err := c.resty.R().SetResult(&result).Get("/api/dashboards/uid/" + uid)
	return result, r, err
}

type PostDashboardRequest struct {
	Dashboard interface{} `json:"dashboard"`
	FolderID  int         `json:"folderId"`
	Overwrite bool        `json:"overwrite"`
}

// nolint:revive
type GrafanaResponse struct {
	ID      *uint   `json:"id"`
	OrgID   *uint   `json:"orgId"`
	Message *string `json:"message"`
	Slug    *string `json:"slug"`
	Version *int    `json:"version"`
	Status  *string `json:"status"`
	UID     *string `json:"uid"`
	URL     *string `json:"url"`
}

// PostDashboard sends a request to create or update a Grafana dashboard.
// It returns the response containing the dashboard details, the HTTP response,
// and any error encountered during the request. This function is useful for
// programmatically managing Grafana dashboards.
func (c *Client) PostDashboard(dashboard PostDashboardRequest) (GrafanaResponse, *resty.Response, error) {
	var grafanaResp GrafanaResponse

	resp, err := c.resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(dashboard).
		SetResult(&grafanaResp). // SetResult will automatically unmarshal the response into grafanaResp
		Post("/api/dashboards/db")

	if err != nil {
		return GrafanaResponse{}, resp, fmt.Errorf("error making API request: %w", err)
	}

	statusCode := resp.StatusCode()
	if statusCode != 200 && statusCode != 201 {
		return GrafanaResponse{}, resp, fmt.Errorf("error creating/updating dashboard, received unexpected status code %d: %s", statusCode, resp.String())
	}

	return grafanaResp, resp, nil
}

// GetAlertsRules retrieves all provisioned alert rules from the API.
// It returns a slice of ProvisionedAlertRule, the HTTP response, and any error encountered.
func (c *Client) GetAlertsRules() ([]ProvisionedAlertRule, *resty.Response, error) {
	var result []ProvisionedAlertRule
	r, err := c.resty.R().SetResult(&result).Get("/api/v1/provisioning/alert-rules")
	return result, r, err
}

// GetAlertRulesForDashboardID retrieves all provisioned alert rules associated with a specific dashboard ID.
// It returns a slice of ProvisionedAlertRule and any error encountered during the process.
func (c *Client) GetAlertRulesForDashboardID(dashboardID string) ([]ProvisionedAlertRule, error) {
	rules, _, err := c.GetAlertsRules()
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

// UnmarshalJSON parses the JSON-encoded data and sets the CustomTime's value.
// It expects the data to represent a timestamp in milliseconds since the epoch.
// This function is useful for decoding JSON data into a CustomTime type.
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

// FormatAlertsTable formats a slice of alerts into a tabular string representation.
// It provides a clear overview of alert states, including timestamps and IDs,
// making it useful for logging or displaying alert information in a structured format.
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

type PostAnnotationResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

// GetAnnotations retrieves a list of annotations based on specified query parameters.
// It allows filtering by alert ID, dashboard ID, type, and time range, ensuring both From and To are set.
// This function is useful for fetching relevant annotations in a Grafana context.
func (c *Client) GetAnnotations(params AnnotationsQueryParams) ([]Annotation, *resty.Response, error) {
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
	r, err := c.resty.R().
		SetResult(&result).
		SetQueryString(query.Encode()).
		Get("/api/annotations")
	return result, r, err
}

// DeleteAnnotation removes an annotation identified by its ID from the server.
// It returns the response from the server and any error encountered during the request.
func (c *Client) DeleteAnnotation(annotationID int64) (*resty.Response, error) {
	urlPath := fmt.Sprintf("/api/annotations/%d", annotationID)

	r, err := c.resty.R().
		Delete(urlPath)

	return r, err
}

// PostAnnotation sends a new annotation to a specified dashboard, allowing users to add notes or comments.
// It returns the response containing the annotation details, the HTTP response, and any error encountered.
func (c *Client) PostAnnotation(annotation PostAnnotation) (PostAnnotationResponse, *resty.Response, error) {
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
	var result PostAnnotationResponse
	r, err := c.resty.R().
		SetBody(a).
		SetResult(&result).
		Post("/api/annotations")
	return result, r, err
}

// PostAlert retrieves alert rules associated with a specified dashboard UID.
// It returns a map of alert rules, the HTTP response, and any error encountered.
func (g *AlertRulerClient) PostAlert(dashboardUID string) (map[string][]interface{}, *resty.Response, error) {
	var result map[string][]interface{}
	r, err := g.resty.R().SetResult(&result).Get("/api/ruler/grafana/api/v1/rules?dashboard_uid=" + url.QueryEscape(dashboardUID))
	return result, r, err
}

// ruler API is deprecated https://github.com/grafana/grafana/issues/74434
func (g *AlertRulerClient) GetAlertsForDashboard(dashboardUID string) (map[string][]interface{}, *resty.Response, error) {
	var result map[string][]interface{}
	r, err := g.resty.R().SetResult(&result).Get("/api/ruler/grafana/api/v1/rules?dashboard_uid=" + url.QueryEscape(dashboardUID))
	return result, r, err
}

type ProvisionedAlertRule struct {
	ID           int64             `json:"id"`
	UID          string            `json:"uid"`
	FolderUID    string            `json:"folderUID"`
	Title        string            `json:"title"`
	Data         []interface{}     `json:"data"`
	ExecErrState string            `json:"execErrState"`
	Labels       map[string]string `json:"labels"`
	RuleGroup    string            `json:"ruleGroup"`
	UpdatedAt    time.Time         `json:"updated"`
	Annotations  map[string]string `json:"annotations"`
}

// String returns a formatted string representation of the ProvisionedAlertRule,
// including its ID, UID, title, labels, annotations, and other relevant details.
// This function is useful for logging and debugging purposes, providing a clear
// overview of the alert rule's properties.
func (p ProvisionedAlertRule) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ID: %d\n", p.ID))
	sb.WriteString(fmt.Sprintf("UID: %s\n", p.UID))
	sb.WriteString(fmt.Sprintf("FolderUID: %s\n", p.FolderUID))
	sb.WriteString(fmt.Sprintf("Title: %s\n", p.Title))
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
