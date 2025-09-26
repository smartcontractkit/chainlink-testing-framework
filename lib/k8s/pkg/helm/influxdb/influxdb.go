package influxdb

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
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
		"global": map[string]interface{}{
			"security": map[string]interface{}{
				"allowInsecureImages": true,
			},
		},
		"image": map[string]interface{}{
			"registry":   reg,
			"repository": "containers/debian-12",
			"tag":        "3.4.2",
		},
		"auth": map[string]interface{}{
			"enabled": "false",
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
	reg := os.Getenv("BITNAMI_PRIVATE_REGISTRY")
	if reg == "" {
		panic("BITNAMI_PRIVATE_REGISTRY not set, it is required for Helm charts")
	}
	dp := defaultProps(reg)
	config.MustMerge(&dp, props)
	return Chart{
		Name:    "influxdb",
		Path:    fmt.Sprintf("%s/charts/debian-12/influxdb:7.1.47", reg),
		Values:  &dp,
		Version: helmVersion,
	}
}
