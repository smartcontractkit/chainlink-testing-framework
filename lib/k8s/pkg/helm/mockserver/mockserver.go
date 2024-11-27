package mockserver

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
	LocalURLsKey    = "mockserver_local"
	InternalURLsKey = "mockserver_internal"
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
		"chain.link/component": "mockserver",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	mockLocal, err := e.Fwd.FindPort("mockserver:0", "mockserver", "serviceport").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	services, err := e.Client.ListServices(e.Cfg.Namespace, fmt.Sprintf("app=%s", m.Name))
	if err != nil {
		return err
	}
	var mockInternal string
	if services != nil && len(services.Items) != 0 {
		mockInternal = fmt.Sprintf("http://%s:1080", services.Items[0].Name)
	} else {
		mockInternal, err = e.Fwd.FindPort("mockserver:0", "mockserver", "serviceport").As(client.RemoteConnection, client.HTTP)
		if err != nil {
			return err
		}
	}
	if e.Cfg.InsideK8s {
		mockLocal = mockInternal
	}

	e.URLs[LocalURLsKey] = []string{mockLocal}
	e.URLs[InternalURLsKey] = []string{mockInternal}
	log.Info().Str("Local Connection", mockLocal).Str("Internal Connection", mockInternal).Msg("Mockserver")
	return nil
}

func defaultProps() map[string]interface{} {
	internalRepo := os.Getenv(config.EnvVarInternalDockerRepo)
	mockserverRepo := "mockserver"
	if internalRepo != "" {
		mockserverRepo = fmt.Sprintf("%s/mockserver", internalRepo)
	}

	return map[string]interface{}{
		"replicaCount": "1",
		"service": map[string]interface{}{
			"type": "NodePort",
			"port": "1080",
		},
		"app": map[string]interface{}{
			"logLevel":               "INFO",
			"serverPort":             "1080",
			"mountedConfigMapName":   "mockserver-config",
			"propertiesFileName":     "mockserver.properties",
			"readOnlyRootFilesystem": "false",
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
			"repository": mockserverRepo,
			"version":    "5.15.0",
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
	chartPath := "chainlink-qa/mockserver"
	if b, err := strconv.ParseBool(os.Getenv(config.EnvVarLocalCharts)); err == nil && b {
		chartPath = fmt.Sprintf("%s/mockserver", projectpath.ChartsRoot)
	}
	return Chart{
		Name:    "mockserver",
		Path:    chartPath,
		Values:  &dp,
		Version: helmVersion,
	}
}
