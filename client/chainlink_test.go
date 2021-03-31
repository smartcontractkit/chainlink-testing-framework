package client

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var spec = `{
  "initiators": [
    {
      "type": "runlog"
    }
  ],
  "tasks": [
    {
      "type": "httpget"
    },
    {
      "type": "jsonparse"
    },
    {
      "type": "multiply"
    },
    {
      "type": "ethuint256"
    },
    {
      "type": "ethtx"
    }
  ]
}`

func TestNodeClient_CreateReadDeleteJob(t *testing.T) {
	server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			assert.Equal(t, "/v2/jobs", req.URL.Path)
			writeResponse(t, rw, http.StatusOK, Job{
				Data: JobData{
					ID: "1",
				},
			})
		case http.MethodGet:
			assert.Equal(t, "/v2/jobs/1", req.URL.Path)
			writeResponse(t, rw, http.StatusOK, nil)
		case http.MethodDelete:
			assert.Equal(t, "/v2/jobs/1", req.URL.Path)
			writeResponse(t, rw, http.StatusNoContent, nil)
		}
	})
	defer server.Close()

	c := newDefaultClient(server.URL)
	c.SetClient(server.Client())

	s, err := c.CreateJob("schemaVersion = 1")
	assert.NoError(t, err)

	err = c.ReadJob(s.Data.ID)
	assert.NoError(t, err)

	err = c.DeleteJob(s.Data.ID)
	assert.NoError(t, err)
}

func TestNodeClient_CreateReadDeleteSpec(t *testing.T) {
	specID := "c142042149f64911bb4698fb08572040"

	server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			assert.Equal(t, "/v2/specs", req.URL.Path)
			writeResponse(t, rw, http.StatusOK, Spec{
				Data: SpecData{ID: specID},
			})
		case http.MethodGet:
			assert.Equal(t, fmt.Sprintf("/v2/specs/%s", specID), req.URL.Path)
			writeResponse(t, rw, http.StatusOK, Response{
				Data: map[string]interface{}{},
			})
		case http.MethodDelete:
			assert.Equal(t, fmt.Sprintf("/v2/specs/%s", specID), req.URL.Path)
			writeResponse(t, rw, http.StatusNoContent, nil)
		}
	})
	defer server.Close()

	c := newDefaultClient(server.URL)
	c.SetClient(server.Client())

	s, err := c.CreateSpec(spec)
	assert.NoError(t, err)

	_, err = c.ReadSpec(s.Data.ID)
	assert.NoError(t, err)

	err = c.DeleteSpec(s.Data.ID)
	assert.NoError(t, err)
}

func TestNodeClient_CreateReadDeleteBridge(t *testing.T) {
	bta := BridgeTypeAttributes{
		Name: "example",
		URL:  "https://example.com",
	}

	server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			assert.Equal(t, "/v2/bridge_types", req.URL.Path)
			writeResponse(t, rw, http.StatusOK, nil)
		case http.MethodGet:
			assert.Equal(t, "/v2/bridge_types/example", req.URL.Path)
			writeResponse(t, rw, http.StatusOK, BridgeType{
				Data: BridgeTypeData{
					Attributes: bta,
				},
			})
		case http.MethodDelete:
			assert.Equal(t, "/v2/bridge_types/example", req.URL.Path)
			writeResponse(t, rw, http.StatusOK, nil)
		}
	})
	defer server.Close()

	c := newDefaultClient(server.URL)
	c.SetClient(server.Client())

	err := c.CreateBridge(&bta)
	assert.NoError(t, err)

	bt, err := c.ReadBridge(bta.Name)
	assert.NoError(t, err)

	assert.Equal(t, bt.Data.Attributes.Name, bta.Name)
	assert.Equal(t, bt.Data.Attributes.URL, bta.URL)

	err = c.DeleteBridge(bta.Name)
	assert.NoError(t, err)
}

func newDefaultClient(url string) Chainlink {
	cl := NewChainlink(&ChainlinkConfig{
		Email:    "admin@node.local",
		Password: "twochains",
		URL:      url,
	})
	return cl
}

func mockedServer(handlerFunc http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handlerFunc)
}

func writeResponse(t *testing.T, rw http.ResponseWriter, statusCode int, obj interface{}) {
	rw.WriteHeader(statusCode)
	if obj == nil {
		return
	}
	b, err := json.Marshal(obj)
	require.Nil(t, err)
	_, err = rw.Write(b)
	require.Nil(t, err)
}
