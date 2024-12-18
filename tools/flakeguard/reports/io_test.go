package reports

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendSplunkEvents(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// body, err := io.ReadAll(r.Body)
		// assert.NoError(t, err)
		// defer r.Body.Close()

		// var receivedEvents []Event
		// err = json.Unmarshal(body, &receivedEvents)
		// assert.NoError(t, err)

		// assert.Equal(t, expectedEvents, receivedEvents)
		// w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
}
