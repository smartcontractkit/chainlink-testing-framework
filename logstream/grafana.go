package logstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

	req, err := http.NewRequest(http.MethodPost, grafanaUrl, bodyReader)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ShorteningFailedErr, err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ShorteningFailedErr, err)
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&responseObject)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ShorteningFailedErr, err)
	}

	return responseObject.Url, nil
}
