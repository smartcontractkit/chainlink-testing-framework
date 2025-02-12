package framework

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPostAnnotationIntegration tests the PostAnnotation method against a real Grafana instance.
func TestPostAnnotationIntegration(t *testing.T) {
	//t.Skip("can only be run manually to debug")
	grafanaURL := os.Getenv("GRAFANA_URL")
	apiKey := os.Getenv("GRAFANA_TOKEN")
	if grafanaURL == "" || apiKey == "" {
		t.Skip("Skipping integration test: GRAFANA_URL or GRAFANA_API_KEY environment variables not set")
	}
	client := NewGrafanaClient(grafanaURL, apiKey)

	annotation := PostAnnotation{
		DashboardUID: "WaspDebug",
		Text:         "CTFv2 test annotation",
		Tags:         []string{"tag-1", "tag-2"},
		Time:         Ptr(time.Now().Add(-1 * time.Minute)),
		TimeEnd:      Ptr(time.Now()),
	}
	response, resp, err := client.PostAnnotation(annotation)
	assert.NoError(t, err, "PostAnnotation should not return an error")
	assert.Equal(t, 200, resp.StatusCode(), "Expected HTTP status code 200")
	assert.NotEmpty(t, response.ID, "Annotation ID should not be empty")
	assert.NotEmpty(t, response.Message, "Annotation message should not be empty")
}

func Ptr[T any](value T) *T {
	return &value
}
