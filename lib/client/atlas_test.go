// atlas_test.go
package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

// TestAtlasClient_GetTransactionDetails_Success verifies that GetTransactionDetails successfully parses a valid response.
func TestAtlasClient_GetTransactionDetails_Success(t *testing.T) {
	// Create a mock Atlas server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the endpoint contains "/atlas/search"
		assert.Contains(t, r.URL.String(), "/atlas/search")
		w.WriteHeader(http.StatusOK)
		// Return a valid JSON response
		response := `{"transactionHash": [{"messageId": "abc123"}]}`
		_, err := w.Write([]byte(response))
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	// Create the Atlas client using the mock server URL
	client := NewAtlasClient(mockServer.URL)

	// Call GetTransactionDetails
	resp, err := client.GetTransactionDetails("abc123")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.TransactionHash, 1)
	assert.Equal(t, "abc123", resp.TransactionHash[0].MessageID)
}

// TestAtlasClient_GetTransactionDetails_NonOK verifies that a non-200 response results in an error.
func TestAtlasClient_GetTransactionDetails_NonOK(t *testing.T) {
	// Create a mock Atlas server that returns a non-OK status code
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), "/atlas/search")
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer mockServer.Close()

	client := NewAtlasClient(mockServer.URL)
	resp, err := client.GetTransactionDetails("abc123")
	assert.Nil(t, resp)
	assert.Error(t, err)
}

// TestAtlasClient_GetMessageDetails_Success verifies that GetMessageDetails correctly unmarshals a valid response.
func TestAtlasClient_GetMessageDetails_Success(t *testing.T) {
	// Prepare a minimal valid response for TransactionDetails.
	detailsResponse := TransactionDetails{
		MessageId:         ptr.Ptr("msg123"),
		State:             ptr.Ptr(1),
		Votes:             ptr.Ptr(5),
		SourceNetworkName: ptr.Ptr("source"),
		DestNetworkName:   ptr.Ptr("dest"),
	}
	responseBytes, err := json.Marshal(detailsResponse)
	assert.NoError(t, err)

	// Create a mock server returning the above JSON.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/atlas/message/msg123", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(responseBytes)
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewAtlasClient(mockServer.URL)
	details, err := client.GetMessageDetails("msg123")
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, "msg123", *details.MessageId)
	assert.Equal(t, "source", *details.SourceNetworkName)
	assert.Equal(t, "dest", *details.DestNetworkName)
}

// TestAtlasClient_GetMessageDetails_MessageNotFound checks that a "Message not found" response is handled correctly.
func TestAtlasClient_GetMessageDetails_MessageNotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/atlas/message/msg123", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Message not found"))
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewAtlasClient(mockServer.URL)
	details, err := client.GetMessageDetails("msg123")
	assert.Nil(t, details)
	assert.Error(t, err)
	assert.Equal(t, "message not found", err.Error())
}

// TestAtlasClient_GetMessageDetails_NonOK verifies that a non-200 status code leads to an error.
func TestAtlasClient_GetMessageDetails_NonOK(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/atlas/message/msg123", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := NewAtlasClient(mockServer.URL)
	details, err := client.GetMessageDetails("msg123")
	assert.Nil(t, details)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve message details")
}

// TestAtlasClient_DebugMode verifies that the RESTY_DEBUG environment variable enables debug mode in the Resty client.
func TestAtlasClient_DebugMode(t *testing.T) {
	os.Setenv("RESTY_DEBUG", "true")
	defer os.Unsetenv("RESTY_DEBUG")

	client := NewAtlasClient("http://example.com")
	assert.True(t, client.RestClient.Debug)
}
