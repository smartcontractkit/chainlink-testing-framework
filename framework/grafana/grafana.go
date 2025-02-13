package framework

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

// NewGrafanaClient initializes a new Grafana client with the specified URL and API key.
func NewGrafanaClient(url, apiKey string) *Client {
	return &Client{
		resty: resty.New().
			SetBaseURL(url).
			SetHeader("Authorization", "Bearer "+apiKey),
	}
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
	Time         time.Time     `json:"time"`
	TimeEnd      time.Time     `json:"timeEnd"`
	Created      time.Time     `json:"created"`
	Updated      time.Time     `json:"updated"`
	Tags         []interface{} `json:"tags"`
	Data         interface{}   `json:"data"`
}

type AnnotationsQueryParams struct {
	Limit        *int
	AlertID      *int
	DashboardID  *int
	DashboardUID *string
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
	if params.Type != nil {
		query.Set("type", *params.Type)
	}

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

// PostAnnotation sends a new annotation to a specified dashboard.
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
