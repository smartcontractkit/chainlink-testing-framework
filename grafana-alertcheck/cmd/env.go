package main

import (
	"fmt"
	"os"
)

// grafanaEnv reads the connection details from the environment only, never
// from a flag — a flag value lands in the process argv and in CI logs, and
// the token must never be logged or otherwise surface in an error string.
func grafanaEnv() (url, token string, err error) {
	url = os.Getenv("GRAFANA_URL")
	if url == "" {
		return "", "", fmt.Errorf("GRAFANA_URL is not set")
	}
	token = os.Getenv("GRAFANA_TOKEN")
	if token == "" {
		return "", "", fmt.Errorf("GRAFANA_TOKEN is not set")
	}
	return url, token, nil
}
