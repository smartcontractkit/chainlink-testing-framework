package parrot

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWriterRecorder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		responseFunc       http.HandlerFunc
		expectedRespCode   int
		expectedRespBody   string
		expectedRespHeader http.Header
	}{
		{
			name: "good response",
			responseFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("Squawk"))
				require.NoError(t, err, "error writing response")
			},
			expectedRespCode: http.StatusOK,
			expectedRespBody: "Squawk",
			expectedRespHeader: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
		{
			name: "error response",
			responseFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Squawk", http.StatusInternalServerError)
			},
			expectedRespCode: http.StatusInternalServerError,
			expectedRespBody: "Squawk\n", // http.Error adds a newline
			expectedRespHeader: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			writerRecorder := newResponseWriterRecorder(recorder)
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			handler := http.HandlerFunc(tc.responseFunc)
			handler.ServeHTTP(writerRecorder, req)

			actualResp := recorder.Result()
			recordedResp := writerRecorder.Result()
			t.Cleanup(func() {
				_ = actualResp.Body.Close()
				_ = recordedResp.Body.Close()
			})

			actualBody, err := io.ReadAll(actualResp.Body)
			require.NoError(t, err, "error reading actual response body")
			recordedBody, err := io.ReadAll(recordedResp.Body)
			require.NoError(t, err, "error reading recorded response body")

			assert.Equal(t, tc.expectedRespCode, actualResp.StatusCode, "actual response has unexpected status code")
			assert.Equal(t, tc.expectedRespCode, recordedResp.StatusCode, "recorded response has unexpected status code")
			assert.Equal(t, tc.expectedRespBody, string(actualBody), "actual response has unexpected body")
			assert.Equal(t, tc.expectedRespBody, string(recordedBody), "recorded response has unexpected body")
		})
	}
}

func TestRecorder(t *testing.T) {

}
