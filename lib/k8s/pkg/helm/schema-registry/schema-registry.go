package schema_registry

import (
	"fmt"
	"os"
	"strconv"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/projectpath"
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
		"chain.link/component": "schema-registry",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	return nil
}

func defaultProps() map[string]interface{} {
	return map[string]interface{}{}
}

func New(props map[string]interface{}) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props map[string]interface{}) environment.ConnectedChart {
	dp := defaultProps()
	config.MustMerge(&dp, props)
	chartPath := "chainlink-qa/schema-registry"
	if b, err := strconv.ParseBool(os.Getenv(config.EnvVarLocalCharts)); err == nil && b {
		chartPath = fmt.Sprintf("%s/schema-registry", projectpath.ChartsRoot)
	}
	return Chart{
		Name:    "cp-schema-registry",
		Path:    chartPath,
		Values:  &dp,
		Version: helmVersion,
	}
}
