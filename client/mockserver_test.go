package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/stretchr/testify/require"
)

func TestSetValuePath(t *testing.T) {
	t.Parallel()

	server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			switch req.URL.Path {
			case "/expectation":
				writeResponse(t, rw, http.StatusCreated, nil)
			default:
				require.Fail(t, "Path '%s' not supported", req.URL.Path)
			}
		} else {
			require.Fail(t, "Method '%s' not supported", req.Method)
		}
	})
	defer server.Close()

	mockServerClient := newDefaultClient(server.URL)
	err := mockServerClient.SetValuePath("variable", 5)
	require.NoError(t, err)
}

func TestPutExpectations(t *testing.T) {
	t.Parallel()

	server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			switch req.URL.Path {
			case "/expectation":
				writeResponse(t, rw, http.StatusCreated, nil)
			default:
				require.Fail(t, "Path '%s' not supported", req.URL.Path)
			}
		} else {
			require.Fail(t, "Method '%s' not supported", req.Method)
		}
	})
	defer server.Close()

	mockServerClient := newDefaultClient(server.URL)
	var nodesInfo []client.NodeInfoJSON

	nodesInitializer := client.HttpInitializer{
		Request:  client.HttpRequest{Path: "/nodes.json"},
		Response: client.HttpResponse{Body: nodesInfo},
	}
	initializers := []client.HttpInitializer{nodesInitializer}

	err := mockServerClient.PutExpectations(initializers)
	require.NoError(t, err)
}

func TestClearExpectations(t *testing.T) {
	t.Parallel()

	server := mockedServer(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			switch req.URL.Path {
			case "/clear":
				writeResponse(t, rw, http.StatusOK, nil)
			default:
				require.Fail(t, "Path '%s' not supported", req.URL.Path)
			}
		} else {
			require.Fail(t, "Method '%s' not supported", req.Method)
		}
	})
	defer server.Close()

	mockServerClient := newDefaultClient(server.URL)
	err := mockServerClient.ClearExpectation(client.PathSelector{Path: "/nodes.json"})
	require.NoError(t, err)
}

func newDefaultClient(url string) *client.MockserverClient {
	ms := client.NewMockserverClient(&client.MockserverConfig{
		LocalURL:   url,
		ClusterURL: url,
	})
	return ms
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
	require.NoError(t, err)
	_, err = rw.Write(b)
	require.NoError(t, err)
}
