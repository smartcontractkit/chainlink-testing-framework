package foundry

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
)

const (
	ChartName = "foundry"
)

type Props struct {
	Values map[string]interface{}
}

type Chart struct {
	ServiceName      string
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
	return ChartName
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
	appInstance := fmt.Sprintf("%s:0", m.ServiceName) // uniquely identifies an instance of an anvil service running in a pod
	var err error
	m.ForwardedHTTPURL, err = e.Fwd.FindPort(appInstance, ChartName, "http").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	m.ForwardedWSURL, err = e.Fwd.FindPort(appInstance, ChartName, "http").As(client.LocalConnection, client.WS)
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
	var serviceName, appLabel string
	// If fullnameOverride is set it is used as the service name and app label
	if props.Values["fullnameOverride"] != nil {
		serviceName = dp.Values["fullnameOverride"].(string)
		appLabel = fmt.Sprintf("app=%s", dp.Values["fullnameOverride"].(string))
	} else {
		serviceName = ChartName
		appLabel = fmt.Sprintf("app=%s", ChartName)
	}
	anvilValues := dp.Values["anvil"].(map[string]any)
	return &Chart{
		ServiceName:    serviceName,
		ClusterWSURL:   fmt.Sprintf("ws://%s:%s", serviceName, anvilValues["port"].(string)),
		ClusterHTTPURL: fmt.Sprintf("http://%s:%s", serviceName, anvilValues["port"].(string)),
		AppLabel:       appLabel,
		Path:           "chainlink-qa/foundry",
		Values:         &dp.Values,
		Props:          dp,
		Version:        helmVersion,
	}
}
