package framework

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPostAnnotationIntegration tests the Annotation method against a real Grafana instance.
func TestPostAnnotationIntegration(t *testing.T) {
	t.Skip("manual grafana integration test")
	grafanaURL := os.Getenv("GRAFANA_URL")
	apiKey := os.Getenv("GRAFANA_TOKEN")
	if grafanaURL == "" || apiKey == "" {
		t.Skip("Skipping integration test: GRAFANA_URL or GRAFANA_API_KEY environment variables not set")
	}
	client := NewGrafanaClient(grafanaURL, apiKey)

	annotation := Annotation{
		Text:      "CTFv2 test annotation",
		StartTime: Ptr(time.Now().Add(-11 * time.Minute)),
		//EndTime:      Ptr(time.Now()),
		DashboardUID: []string{"WaspDebug", "e98b5451-12dc-4a8b-9576-2c0b67ddbd0c"},
		Tags:         []string{"tag-3", "tag-4"},
	}
	response, resp, err := client.Annotate(annotation)
	assert.NoError(t, err, "Annotate should not return an error")
	assert.Equal(t, 200, resp[0].StatusCode(), "Expected HTTP status code 200")
	assert.NotEmpty(t, response[0].ID, "Annotation ID should not be empty")
	assert.NotEmpty(t, response[0].Message, "Annotation message should not be empty")
	assert.Equal(t, 200, resp[1].StatusCode(), "Expected HTTP status code 200")
	assert.NotEmpty(t, response[1].ID, "Annotation ID should not be empty")
	assert.NotEmpty(t, response[1].Message, "Annotation message should not be empty")
}

func Ptr[T any](value T) *T {
	return &value
}
