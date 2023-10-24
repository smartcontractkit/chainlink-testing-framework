package grafana

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
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"memory": "1000Mi",
				"cpu":    "1.5",
			},
			"requests": map[string]interface{}{
				"memory": "500Mi",
				"cpu":    "1.0",
			},
		},
		"rbac": map[string]interface{}{
			"create": "false",
		},
		"testFramework": map[string]interface{}{
			"enabled": false,
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
		Name:    "grafana",
		Path:    "grafana/grafana",
		Values:  &dp,
		Version: helmVersion,
	}
}
