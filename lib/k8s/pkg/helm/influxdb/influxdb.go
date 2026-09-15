package influxdb

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
)

// influxdataChartURL pins the chart version; helm ignores --version for URL refs.
const (
	influxdataChartURL   = "https://github.com/influxdata/helm-charts/releases/download/influxdb-4.12.5/influxdb-4.12.5.tgz"
	defaultImageRegistry = "804282218731.dkr.ecr.us-west-2.amazonaws.com"
)

type Props struct {
}

type Chart struct {
	Name    string
	Path    string
	Version string
	Props   *Props
	Values  *map[string]interface{}
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return m.Name
}

func (m Chart) GetPath() string {
	return m.Path
}

func (m Chart) GetVersion() string {
	return m.Version
}

func (m Chart) GetProps() interface{} {
	return m.Props
}

func (m Chart) GetValues() *map[string]interface{} {
	return m.Values
}

func (m Chart) GetLabels() map[string]string {
	return map[string]string{
		"chain.link/component": "influxdb",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	return nil
}

func defaultProps(reg string) map[string]interface{} {
	return map[string]interface{}{
		"image": map[string]interface{}{
			"repository": fmt.Sprintf("%s/docker-io/library/influxdb", reg),
			"tag":        "1.8.10-alpine",
		},
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"memory": "19000Mi",
				"cpu":    "6",
			},
			"requests": map[string]interface{}{
				"memory": "16000Mi",
				"cpu":    "5",
			},
		},
	}
}

func registry() string {
	if reg := os.Getenv(config.EnvVarInfluxdbImageRegistry); reg != "" {
		return reg
	}
	return defaultImageRegistry
}

func New(props map[string]interface{}) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned keeps its signature for API compatibility; the chart version is pinned in influxdataChartURL and the version argument is ignored.
func NewVersioned(helmVersion string, props map[string]interface{}) environment.ConnectedChart {
	dp := defaultProps(registry())
	config.MustMerge(&dp, props)
	return Chart{
		Name:    "influxdb",
		Path:    influxdataChartURL,
		Values:  &dp,
		Version: "",
	}
}
