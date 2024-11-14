package main

import (
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/dashboard"
)

// main creates and deploys a default dashboard using environment variables for configuration.
// It sets up the dashboard without extensions or non-functional requirements, enabling straightforward deployment.
func main() {
	// just default dashboard, no NFRs, no dashboard extensions
	// see examples/alerts.go for an example extension
	d, err := dashboard.NewDashboard(nil, nil)
	if err != nil {
		panic(err)
	}
	// set env vars
	//export GRAFANA_URL=...
	//export GRAFANA_TOKEN=...
	//export DATA_SOURCE_NAME=Loki
	//export DASHBOARD_FOLDER=LoadTests
	//export DASHBOARD_NAME=Wasp
	if _, err := d.Deploy(); err != nil {
		panic(err)
	}
}
