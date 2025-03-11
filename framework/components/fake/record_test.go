package fake_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

func TestSmokeComponentFake(t *testing.T) {
	cfg := &fake.Input{
		Port: 9111,
	}
	out, err := fake.NewFakeDataProvider(cfg)
	require.NoError(t, err)
	r := resty.New().SetBaseURL(out.BaseURLHost)

	t.Run("can mock a response with Func", func(t *testing.T) {
		apiPath := "/fake/api/1"
		err = fake.Func(
			"GET",
			apiPath,
			func(c *gin.Context) {
				c.JSON(200, gin.H{
					"status": "ok",
				})
			},
		)
		require.NoError(t, err)
		var respBody struct {
			Status string `json:"status"`
		}
		resp, err := r.R().SetResult(&respBody).Get(apiPath)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode())
		require.Equal(t, "ok", respBody.Status)
	})

	t.Run("can mock a response with JSON", func(t *testing.T) {
		apiPath := "/fake/api/2"
		err = fake.JSON(
			"GET",
			apiPath,
			map[string]any{
				"status": "ok",
			}, 200,
		)
		require.NoError(t, err)
		var respBody struct {
			Status string `json:"status"`
		}
		resp, err := r.R().SetResult(&respBody).Get(apiPath)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode())
		require.Equal(t, "ok", respBody.Status)
	})

	t.Run("can record request/response and access it", func(t *testing.T) {
		method := "POST"
		apiPath := "/fake/api/3"
		err = fake.JSON(
			method,
			apiPath,
			map[string]any{
				"status": "ok",
			}, 200,
		)
		require.NoError(t, err)
		reqBody := struct {
			SomeData string `json:"some_data"`
		}{
			SomeData: "some_data",
		}
		var respBody struct {
			Status string `json:"status"`
		}
		_, err := r.R().SetBody(reqBody).SetResult(&respBody).Post(apiPath)
		require.NoError(t, err)

		// get request and response
		recordedData, err := fake.R.Get(method, apiPath)
		require.NoError(t, err)
		require.Equal(t, []*fake.Record{
			{
				Method: "POST",
				Path:   apiPath,
				Headers: http.Header{
					"Accept-Encoding": []string{"gzip"},
					"Content-Type":    []string{"application/json"},
					"Content-Length":  []string{"25"},
					"User-Agent":      []string{"go-resty/2.15.3 (https://github.com/go-resty/resty)"},
				},
				ReqBody: `{"some_data":"some_data"}`,
				ResBody: `{"status":"ok"}`,
				Status:  200,
			},
		}, recordedData)
	})
}
