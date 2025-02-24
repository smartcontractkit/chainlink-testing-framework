package parrot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/projectpath"

	"github.com/rs/zerolog/log"
)

const (
	LocalURLsKey    = "parrot_local"
	InternalURLsKey = "parrot_internal"
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
		"chain.link/component": "parrot",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	parrotLocal, err := e.Fwd.FindPort("parrot:0", "parrot", "http").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	services, err := e.Client.ListServices(e.Cfg.Namespace, fmt.Sprintf("app=%s", m.Name))
	if err != nil {
		return err
	}
	var parrotInternal string
	if services != nil && len(services.Items) != 0 {
		parrotInternal = fmt.Sprintf("http://%s:80", services.Items[0].Name)
	} else {
		parrotInternal, err = e.Fwd.FindPort("parrot:0", "parrot", "http").As(client.RemoteConnection, client.HTTP)
		if err != nil {
			return err
		}
	}
	if e.Cfg.InsideK8s {
		parrotLocal = parrotInternal
	}

	e.URLs[LocalURLsKey] = []string{parrotLocal}
	e.URLs[InternalURLsKey] = []string{parrotInternal}
	log.Info().Str("Local Connection", parrotLocal).Str("Internal Connection", parrotInternal).Msg("Parrot")
	return nil
}

func defaultProps() map[string]interface{} {
	internalRepo := os.Getenv(config.EnvVarInternalDockerRepo)
	mockserverRepo := "parrot"
	if internalRepo != "" {
		mockserverRepo = fmt.Sprintf("%s/parrot", internalRepo)
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
	chartPath := "chainlink-qa/parrot"
	if b, err := strconv.ParseBool(os.Getenv(config.EnvVarLocalCharts)); err == nil && b {
		chartPath = fmt.Sprintf("%s/parrot", projectpath.ChartsRoot)
	}
	return Chart{
		Name:    "parrot",
		Path:    chartPath,
		Values:  &dp,
		Version: helmVersion,
	}
}
