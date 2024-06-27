package logstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

const (
	ShorteningFailedErr = "failed to shorten Grafana URL"
)

func ShortenUrl(grafanaUrl, urlToShorten, bearerToken string) (string, error) {
	jsonBody := []byte(`{"path":"` + urlToShorten + `"}`)
	bodyReader := bytes.NewReader(jsonBody)

	var responseObject struct {
		Uid string `json:"uid"`
		Url string `json:"url"`
	}

	var res *http.Response

	if err := retry.Do(
		func() error {
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sapi/short-urls", grafanaUrl), bodyReader)
			if err != nil {
				return err
			}

			req.Header.Add("Authorization", "Bearer "+bearerToken)
			req.Header.Add("Content-Type", "application/json")

			res, err = http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			if res.StatusCode != http.StatusOK {
				return err
			}

			return nil
		},
		retry.DelayType(retry.FixedDelay),
		retry.Attempts(10),
		retry.Delay(time.Duration(1)*time.Second),
	); err != nil {
		return "", fmt.Errorf("%s: %w", ShorteningFailedErr, err)
	}

	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&responseObject)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ShorteningFailedErr, err)
	}

	return responseObject.Url, nil
}
