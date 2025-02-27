package mock_adapter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/projectpath"
)

const (
	LocalURLsKey    = "qa_mock_adapter_local"    // #nosec G101
	InternalURLsKey = "qa_mock_adapter_internal" // #nosec G101
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
		"chain.link/component": "mock-adapter",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	mockLocal, err := e.Fwd.FindPort("qa-mock-adapter:0", "qa-mock-adapter", "serviceport").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	services, err := e.Client.ListServices(e.Cfg.Namespace, fmt.Sprintf("app=%s", m.Name))
	if err != nil {
		return err
	}
	var mockInternal string
	if services != nil && len(services.Items) != 0 {
		mockInternal = fmt.Sprintf("http://%s:6060", services.Items[0].Name)
	} else {
		mockInternal, err = e.Fwd.FindPort("qa-mock-adapter:0", "qa-mock-adapter", "serviceport").As(client.RemoteConnection, client.HTTP)
		if err != nil {
			return err
		}
	}
	if e.Cfg.InsideK8s {
		mockLocal = mockInternal
	}

	e.URLs[LocalURLsKey] = []string{mockLocal}
	e.URLs[InternalURLsKey] = []string{mockInternal}
	log.Info().Str("Local Connection", mockLocal).Str("Internal Connection", mockInternal).Msg("QA Mock Adapter")
	return nil
}

func defaultProps() map[string]interface{} {
	internalRepo := os.Getenv(config.EnvVarInternalDockerRepo)
	qaMockAdapterRepo := "qa-mock-adapter"
	if internalRepo != "" {
		qaMockAdapterRepo = internalRepo
	}

	return map[string]interface{}{
		"replicaCount": "1",
		"service": map[string]interface{}{
			"type": "NodePort",
			"port": "6060",
		},
		"app": map[string]interface{}{
			"serverPort": "6060",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "200m",
					"memory": "256Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "200m",
					"memory": "256Mi",
				},
			},
		},
		"image": map[string]interface{}{
			"repository": qaMockAdapterRepo,
			"snapshot":   false,
			"pullPolicy": "IfNotPresent",
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
	chartPath := "chainlink-qa/qa-mock-adapter"
	if b, err := strconv.ParseBool(os.Getenv(config.EnvVarLocalCharts)); err == nil && b {
		chartPath = fmt.Sprintf("%s/qa-mock-adapter", projectpath.ChartsRoot)
	}
	return Chart{
		Name:    "qa-mock-adapter",
		Path:    chartPath,
		Values:  &dp,
		Version: helmVersion,
	}
}
