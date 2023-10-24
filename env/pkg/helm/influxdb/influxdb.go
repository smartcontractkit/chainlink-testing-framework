package influxdb

import (
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
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

func (m Chart) ExportData(e *environment.Environment) error {
	return nil
}

func defaultProps() map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"enabled": "false",
		},
		"image": map[string]interface{}{
			"tag": "1.7.10",
		},
		"influxdb": map[string]interface{}{
			"readinessProbe": map[string]interface{}{
				"enabled": false,
			},
			"livenessProbe": map[string]interface{}{
				"enabled": false,
			},
			"startupProbe": map[string]interface{}{
				"enabled": false,
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
		},
	}
}

func New(props map[string]interface{}) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props map[string]interface{}) environment.ConnectedChart {
	dp := defaultProps()
	config.MustMerge(&dp, props)
	return Chart{
		Name:    "influxdb",
		Path:    "bitnami/influxdb",
		Values:  &dp,
		Version: helmVersion,
	}
}
