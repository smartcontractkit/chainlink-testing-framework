package framework

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	resty *resty.Client
}

// NewGrafanaClient initializes a new Grafana client with the specified URL and API key.
func NewGrafanaClient(url, bearerToken string) *Client {
	return &Client{
		resty: resty.New().
			SetBaseURL(url).
			SetHeader("Authorization", "Bearer "+bearerToken),
	}
}

type Annotation struct {
	PanelID      *int
	DashboardUID []string
	StartTime    *time.Time
	EndTime      *time.Time
	Tags         []string
	Text         string
}

type PostAnnotationResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

// Annotate adds annotation to all the dashboards, works for both single point annotation with just StartTime and for ranges with StartTime/EndTime
func (c *Client) Annotate(annotation Annotation) ([]PostAnnotationResponse, []*resty.Response, error) {
	var results []PostAnnotationResponse
	var responses []*resty.Response

	for _, uid := range annotation.DashboardUID {
		a := map[string]interface{}{
			"dashboardUID": uid,
			"tags":         annotation.Tags,
			"text":         annotation.Text,
		}
		if annotation.PanelID != nil {
			a["panelId"] = *annotation.PanelID
		}
		if annotation.StartTime != nil {
			a["time"] = annotation.StartTime.UnixMilli()
		}
		if annotation.EndTime != nil {
			a["timeEnd"] = annotation.EndTime.UnixMilli()
		}

		var result PostAnnotationResponse
		r, err := c.resty.R().
			SetBody(a).
			SetResult(&result).
			Post("/api/annotations")
		if err != nil {
			return nil, nil, err // Return early if any request fails
		}

		results = append(results, result)
		responses = append(responses, r)
	}

	return results, responses, nil
}
