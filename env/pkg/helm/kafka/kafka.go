package kafka

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
			"clientProtocol":      "plaintext",
			"interBrokerProtocol": "plaintext",
		},
		"image": map[string]interface{}{
			"debug": true,
		},
		"provisioning": map[string]interface{}{
			"enabled": true,
			"resources": map[string]interface{}{
				"limits": map[string]interface{}{
					"cpu":    "0.1",
					"memory": "500M",
				},
			},
		},
		"zookeeper": map[string]interface{}{
			"persistence": map[string]interface{}{
				"enabled": true,
			},
		},
		"podSecurityContext": map[string]interface{}{
			"enabled": false,
		},
		"containerSecurityContext": map[string]interface{}{
			"enabled": false,
		},
		"persistence": map[string]interface{}{
			"enabled": false,
		},
		"livenessProbe": map[string]interface{}{
			"enabled":             true,
			"initialDelaySeconds": 10,
			"timeoutSeconds":      5,
			"failureThreshold":    3,
			"periodSeconds":       10,
			"successThreshold":    1,
		},
		"readinessProbe": map[string]interface{}{
			"enabled":             true,
			"initialDelaySeconds": 5,
			"failureThreshold":    6,
			"timeoutSeconds":      5,
			"periodSeconds":       10,
			"successThreshold":    1,
		},
		"startupProbe": map[string]interface{}{
			"enabled":             true,
			"initialDelaySeconds": 30,
			"periodSeconds":       10,
			"timeoutSeconds":      1,
			"failureThreshold":    15,
			"successThreshold":    1,
		},
		"commonLabels": map[string]interface{}{
			"app": "kafka",
		},
		"commonAnnotations": map[string]interface{}{
			"app": "kafka",
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
		Name:    "kafka",
		Path:    "bitnami/kafka",
		Values:  &dp,
		Version: helmVersion,
	}
}
