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
	NetworkName string
	Values      map[string]interface{}
}

type Chart struct {
	Name    string
	Path    string
	Version string
	Props   *Props
	Values  *map[string]any
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
	podName := fmt.Sprintf("%s-%s:0", m.Props.NetworkName, ChartName)
	localHttp, err := e.Fwd.FindPort(podName, ChartName, "http").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	internalHttp, err := e.Fwd.FindPort(podName, ChartName, "http").As(client.RemoteConnection, client.HTTP)
	if err != nil {
		return err
	}
	parsed, err := url.Parse(internalHttp)
	if err != nil {
		return err
	}
	port := parsed.Port()
	localWs, err := e.Fwd.FindPort(podName, ChartName, "http").As(client.LocalConnection, client.WS)
	if err != nil {
		return err
	}
	if e.Cfg.InsideK8s {
		services, err := e.Client.ListServices(e.Cfg.Namespace, fmt.Sprintf("app=%s-%s", m.Props.NetworkName, ChartName))
		if err != nil {
			return err
		}
		internalWs := fmt.Sprintf("ws://%s:%s", services.Items[0].Name, port)
		internalHttp = fmt.Sprintf("http://%s:%s", services.Items[0].Name, port)
		e.URLs[m.Props.NetworkName] = []string{internalWs}
		e.URLs[m.Props.NetworkName+"_http"] = []string{internalHttp}
	} else {
		e.URLs[m.Props.NetworkName] = []string{localWs}
		e.URLs[m.Props.NetworkName+"_http"] = []string{localHttp}
	}

	for k, v := range e.URLs {
		if strings.Contains(k, m.Props.NetworkName) {
			log.Info().Str("Name", k).Strs("URLs", v).Msg("Forked network")
		}
	}

	return nil
}

func defaultProps() *Props {
	return &Props{
		NetworkName: "fork",
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

func New(props *Props) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props *Props) environment.ConnectedChart {
	dp := defaultProps()
	config.MustMerge(dp, props)
	config.MustMerge(&dp.Values, props.Values)
	dp.NetworkName = strings.ReplaceAll(strings.ToLower(dp.NetworkName), " ", "-")
	return Chart{
		Name:    dp.NetworkName,
		Path:    "chainlink-qa/foundry",
		Values:  &dp.Values,
		Props:   dp,
		Version: helmVersion,
	}
}
