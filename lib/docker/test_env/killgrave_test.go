// nolint
package test_env

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

type kgTest struct {
	Name          string
	AdapterResult interface{}
	Expected      string
	Path          string
	Headers       map[string]string
}

func TestKillgraveNoUserImposters(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	k := NewKillgrave([]string{network.Name}, "").
		WithTestInstance(t)
	err = k.StartContainer()
	require.NoError(t, err)

	runTestWithExpectations(t, k, []kgTest{})
}

func TestKillgraveMocks(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	k := NewKillgrave([]string{network.Name}, "./killgrave_imposters").
		WithTestInstance(t)
	err = k.StartContainer()
	require.NoError(t, err)

	expectations := []kgTest{
		{
			Name:     "LoadedSix",
			Expected: "{\"id\":\"\",\"error\":null,\"data\":{\"result\":6}}",
			Path:     "/six",
			Headers:  map[string]string{"Content-Type": "text/plain"},
		},
		{
			Name:     "SetStringValuePath",
			Expected: "bar",
			Path:     "/stringany",
			Headers:  map[string]string{"Content-Type": "text/plain"},
		},
		{
			Name:          "SetAdapterBasedAnyValuePath",
			AdapterResult: "bar",
			Expected:      "{\"id\":\"\",\"data\":{\"result\":\"bar\"},\"error\":null}",
			Path:          "/adapterany",
			Headers:       map[string]string{"Content-Type": "application/json"},
		},
		{
			Name:          "SetAdapterBasedAnyValuePathObject",
			AdapterResult: map[string]string{"foo": "bar"},
			Expected:      "{\"id\":\"\",\"data\":{\"result\":{\"foo\":\"bar\"}},\"error\":null}",
			Path:          "/adapteranyobject",
			Headers:       map[string]string{"Content-Type": "application/json"},
		},
		{
			Name:          "SetAdapterBasedIntValuePath",
			AdapterResult: 5,
			Expected:      "{\"id\":\"\",\"data\":{\"result\":5},\"error\":null}",
			Path:          "/adapterint",
			Headers:       map[string]string{"Content-Type": "application/json"},
		},
		{
			Name:          "LongPathForAdapterInt",
			AdapterResult: 5,
			Expected:      "{\"id\":\"\",\"data\":{\"result\":5},\"error\":null}",
			Path:          "/long/adapter/path",
			Headers:       map[string]string{"Content-Type": "application/json"},
		},
		{
			Name:          "MissingLeadingSlash",
			AdapterResult: 5,
			Expected:      "{\"id\":\"\",\"data\":{\"result\":5},\"error\":null}",
			Path:          "noleadingslash",
			Headers:       map[string]string{"Content-Type": "application/json"},
		},
	}

	runTestWithExpectations(t, k, expectations)
}

func runTestWithExpectations(t *testing.T, k *Killgrave, expectations []kgTest) {
	n := t.Name()
	expectations = append(expectations, kgTest{
		Name:     "DefaultFive",
		Expected: "{\"id\":\"\",\"data\":{\"result\":5},\"error\":null}",
		Path:     "/five",
		Headers:  map[string]string{"Content-Type": "text/plain"},
	})
	var err error
	// Check the different kinds of responses
	for _, test := range expectations {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			m := []string{http.MethodGet}
			switch t.Name() {
			case fmt.Sprintf("%s/DefaultFive", n):
				// do nothing, it is provided by default
			case fmt.Sprintf("%s/LoadedSix", n):
				// do nothing, it is loaded from the imposters directory
			case fmt.Sprintf("%s/SetStringValuePath", n):
				err = k.SetStringValuePath(test.Path, m, test.Headers, test.Expected)
			case fmt.Sprintf("%s/SetAdapterBasedAnyValuePath", n):
				err = k.SetAdapterBasedAnyValuePath(test.Path, m, test.AdapterResult)
			case fmt.Sprintf("%s/SetAdapterBasedAnyValuePathObject", n):
				err = k.SetAdapterBasedAnyValuePath(test.Path, m, test.AdapterResult)
			case fmt.Sprintf("%s/SetAdapterBasedIntValuePath", n):
				err = k.SetAdapterBasedIntValuePath(test.Path, m, test.AdapterResult.(int))
			case fmt.Sprintf("%s/LongPathForAdapterInt", n):
				err = k.SetAdapterBasedIntValuePath(test.Path, m, test.AdapterResult.(int))
			case fmt.Sprintf("%s/MissingLeadingSlash", n):
				err = k.SetAdapterBasedIntValuePath(test.Path, m, test.AdapterResult.(int))
			default:
				require.Fail(t, fmt.Sprintf("unknown test name %s", t.Name()))
			}
			require.NoError(t, err)

			var url string
			if strings.HasPrefix(test.Path, "/") {
				url = fmt.Sprintf("%s%s", k.ExternalEndpoint, test.Path)
			} else {
				url = fmt.Sprintf("%s/%s", k.ExternalEndpoint, test.Path)
			}
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			req, err := http.NewRequest(m[0], url, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("url: %s", url))

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			require.NoError(t, err)

			responseString := buf.String()
			require.Equal(t, test.Expected, responseString)
		})
	}
}

func TestKillgraveRequestDump(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	k := NewKillgrave([]string{network.Name}, "./killgrave_imposters").
		WithTestInstance(t)
	err = k.StartContainer()
	require.NoError(t, err)

	path := "/stringany"
	m := []string{http.MethodGet}
	headers := map[string]string{"Content-Type": "text/plain"}
	err = k.SetStringValuePath("/stringany", m, headers, "{\"id\":\"\",\"data\":{\"result\":5},\"error\":null}")
	require.NoError(t, err)
	var url string
	if strings.HasPrefix(path, "/") {
		url = fmt.Sprintf("%s%s", k.ExternalEndpoint, path)
	} else {
		url = fmt.Sprintf("%s/%s", k.ExternalEndpoint, path)
	}
	bodyRequest := []byte("{\n\"a\":5,\n\"b\":6\n}")
	req1, err := http.NewRequest(m[0], url, bytes.NewBuffer(bodyRequest))
	require.NoError(t, err)
	req1.Header.Set("Content-Type", "application/json")
	req2, err := http.NewRequest(m[0], url, bytes.NewBuffer(bodyRequest))
	require.NoError(t, err)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp1, err := client.Do(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()
	require.Equal(t, http.StatusOK, resp1.StatusCode, fmt.Sprintf("url: %s", url))
	resp2, err := client.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()
	require.Equal(t, http.StatusOK, resp2.StatusCode, fmt.Sprintf("url: %s", url))

	requests, err := k.GetReceivedRequests()
	require.NoError(t, err)
	fmt.Printf("Requests: %+v\n", requests)
	require.Equal(t, 2, len(requests))
	require.Equal(t, string(bodyRequest), requests[0].Body)
	require.Equal(t, string(bodyRequest), requests[1].Body)
}
