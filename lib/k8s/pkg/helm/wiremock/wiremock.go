package wiremock

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
	LocalURLsKey    = "wiremock_local"
	InternalURLsKey = "wiremock_internal"
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
		"chain.link/component": "wiremock",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	mockLocal, err := e.Fwd.FindPort("wiremock:0", "wiremock", "serviceport").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	services, err := e.Client.ListServices(e.Cfg.Namespace, fmt.Sprintf("app=%s", m.Name))
	if err != nil {
		return err
	}
	var mockInternal string
	if services != nil && len(services.Items) != 0 {
		mockInternal = fmt.Sprintf("http://%s:80", services.Items[0].Name)
	} else {
		mockInternal, err = e.Fwd.FindPort("wiremock:0", "wiremock", "serviceport").As(client.RemoteConnection, client.HTTP)
		if err != nil {
			return err
		}
	}
	if e.Cfg.InsideK8s {
		mockLocal = mockInternal
	}

	e.URLs[LocalURLsKey] = []string{mockLocal}
	e.URLs[InternalURLsKey] = []string{mockInternal}
	log.Info().Str("Local Connection", mockLocal).Str("Internal Connection", mockInternal).Msg("Wiremock")
	return nil
}

func defaultProps() map[string]interface{} {
	internalRepo := os.Getenv(config.EnvVarInternalDockerRepo)
	wiremockRepo := "wiremock/wiremock"
	if internalRepo != "" {
		wiremockRepo = fmt.Sprintf("%s/wiremock/wiremock", internalRepo)
	}

	return map[string]interface{}{
		"replicaCount": "1",
		"image": map[string]interface{}{
			"repository": wiremockRepo,
			"tag":        "3.4.2",
			"pullPolicy": "IfNotPresent",
		},
		"nameOverride":     "",
		"fullnameOverride": "",
		"service": map[string]interface{}{
			"name":         "wiremock",
			"type":         "ClusterIP",
			"externalPort": "80",
			"internalPort": "9021",
		},
		"ingress": map[string]interface{}{
			"enabled": "false",
		},
		"env": map[string]interface{}{
			"WIREMOCK_OPTIONS": "--port=9021 --async-response-enabled=true --async-response-threads=100 --max-request-journal=1000 --local-response-templating --root-dir=/home/wiremock/storage",
		},
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2Gi",
			},
		},
		"storage": map[string]interface{}{
			"capacity": "100Mi",
		},
		"scheme": "HTTP",
	}
}

func New(props map[string]interface{}) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props map[string]interface{}) environment.ConnectedChart {
	dp := defaultProps()
	config.MustMerge(&dp, props)
	chartPath := "chainlink-qa/wiremock"
	if b, err := strconv.ParseBool(os.Getenv(config.EnvVarLocalCharts)); err == nil && b {
		chartPath = fmt.Sprintf("%s/wiremock", projectpath.ChartsRoot)
	}
	return Chart{
		Name:    "wiremock",
		Path:    chartPath,
		Values:  &dp,
		Version: helmVersion,
	}
}
