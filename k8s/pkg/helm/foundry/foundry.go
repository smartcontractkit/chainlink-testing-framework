package foundry

import (
	"fmt"
	"net/url"
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
	ServiceName string
	AppLabel    string
	Path        string
	Version     string
	Props       *Props
	Values      *map[string]any
	WSURL       []string
	HTTPURL     []string
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
	forwardedHttp, err := e.Fwd.FindPort(appInstance, ChartName, "http").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	internalHttp, err := e.Fwd.FindPort(appInstance, ChartName, "http").As(client.RemoteConnection, client.HTTP)
	if err != nil {
		return err
	}
	parsed, err := url.Parse(internalHttp)
	if err != nil {
		return err
	}
	port := parsed.Port()
	forwardedWs, err := e.Fwd.FindPort(appInstance, ChartName, "http").As(client.LocalConnection, client.WS)
	if err != nil {
		return err
	}
	if e.Cfg.InsideK8s {
		services, err := e.Client.ListServices(e.Cfg.Namespace, m.AppLabel)
		if err != nil {
			return err
		}
		internalWs := fmt.Sprintf("ws://%s:%s", services.Items[0].Name, port)
		internalHttp = fmt.Sprintf("http://%s:%s", services.Items[0].Name, port)
		m.WSURL = []string{internalWs}
		m.HTTPURL = []string{internalHttp}
		e.URLs[appInstance] = []string{internalWs}
		e.URLs[appInstance+"_http"] = []string{internalHttp}
	} else {
		m.WSURL = []string{forwardedWs}
		m.HTTPURL = []string{forwardedHttp}
		e.URLs[appInstance] = []string{forwardedWs}
		e.URLs[appInstance+"_http"] = []string{forwardedHttp}
	}

	for k, v := range e.URLs {
		if strings.Contains(k, appInstance) {
			log.Info().Str("Name", k).Strs("URLs", v).Msg("Forked network")
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
		serviceName = props.Values["fullnameOverride"].(string)
		appLabel = fmt.Sprintf("app=%s", props.Values["fullnameOverride"].(string))
	} else {
		serviceName = ChartName
		appLabel = fmt.Sprintf("app=%s", ChartName)
	}
	return &Chart{
		ServiceName: serviceName,
		AppLabel:    appLabel,
		// Path:    "chainlink-qa/foundry",
		Path:    "/Users/lukasz/Documents/smartcontractkit/chainlink-testing-framework/charts/foundry",
		Values:  &dp.Values,
		Props:   dp,
		Version: helmVersion,
	}
}
