package foundry

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
)

type Props struct {
	Values map[string]interface{}
}

type Chart struct {
	Name             string
	AppLabel         string
	Path             string
	Version          string
	Props            *Props
	Values           *map[string]any
	ClusterWSURL     string
	ClusterHTTPURL   string
	ForwardedWSURL   string
	ForwardedHTTPURL string
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

func (m *Chart) ExportData(e *environment.Environment) error {
	appInstance := fmt.Sprintf("%s:0", m.Name) // uniquely identifies an instance of an anvil service running in a pod
	var err error
	m.ForwardedHTTPURL, err = e.Fwd.FindPort(appInstance, m.GetName(), "http").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	m.ForwardedWSURL, err = e.Fwd.FindPort(appInstance, m.GetName(), "http").As(client.LocalConnection, client.WS)
	if err != nil {
		return err
	}

	e.URLs[appInstance+"_cluster_ws"] = []string{m.ClusterWSURL}
	e.URLs[appInstance+"_cluster_http"] = []string{m.ClusterHTTPURL}
	e.URLs[appInstance+"_forwarded_ws"] = []string{m.ForwardedWSURL}
	e.URLs[appInstance+"_forwarded_http"] = []string{m.ForwardedHTTPURL}

	for k, v := range e.URLs {
		if strings.Contains(k, appInstance) {
			log.Info().Str("Name", k).Strs("URLs", v).Msg("Anvil URLs")
		}
	}

	return nil
}

func defaultProps() *Props {
	return &Props{
		Values: map[string]any{
			"replicaCount": "1",
			"anvil": map[string]any{
				"host":                      "0.0.0.0",
				"port":                      "8545",
				"blockTime":                 1,
				"forkRetries":               "5",
				"forkTimeout":               "45000",
				"forkComputeUnitsPerSecond": "330",
				"chainId":                   "1337",
			},
		},
	}
}

func New(props *Props) *Chart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props *Props) *Chart {
	dp := defaultProps()
	config.MustMerge(dp, props)
	config.MustMerge(&dp.Values, props.Values)
	var name string
	if props.Values["fullnameOverride"] != nil {
		// If fullnameOverride is set it is used as the service name and app label
		name = dp.Values["fullnameOverride"].(string)
	} else {
		// Use default name with random suffix to allow multiple charts in the same namespace
		name = fmt.Sprintf("anvil-%s", uuid.New().String()[0:5])
	}
	anvilValues := dp.Values["anvil"].(map[string]any)
	return &Chart{
		Name:           name,
		AppLabel:       fmt.Sprintf("app=%s", name),
		ClusterWSURL:   fmt.Sprintf("ws://%s:%s", name, anvilValues["port"].(string)),
		ClusterHTTPURL: fmt.Sprintf("http://%s:%s", name, anvilValues["port"].(string)),
		Path:           "chainlink-qa/foundry",
		Values:         &dp.Values,
		Props:          dp,
		Version:        helmVersion,
	}
}
