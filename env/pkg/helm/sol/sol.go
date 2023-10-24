package sol

import (
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/environment"
)

type Props struct {
	NetworkName string   `envconfig:"network_name"`
	HttpURLs    []string `envconfig:"http_url"`
	WsURLs      []string `envconfig:"ws_url"`
	Values      map[string]interface{}
}

type HelmProps struct {
	Name    string
	Path    string
	Version string
	Values  *map[string]interface{}
}

type Chart struct {
	HelmProps *HelmProps
	Props     *Props
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetProps() interface{} {
	return m.Props
}

func (m Chart) GetName() string {
	return m.HelmProps.Name
}

func (m Chart) GetPath() string {
	return m.HelmProps.Path
}

func (m Chart) GetVersion() string {
	return m.HelmProps.Version
}

func (m Chart) GetValues() *map[string]interface{} {
	return m.HelmProps.Values
}

func (m Chart) ExportData(e *environment.Environment) error {
	netLocal, err := e.Fwd.FindPort("sol:0", "sol-val", "http-rpc").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	netLocalWS, err := e.Fwd.FindPort("sol:0", "sol-val", "ws-rpc").As(client.LocalConnection, client.WS)
	if err != nil {
		return err
	}
	netInternal, err := e.Fwd.FindPort("sol:0", "sol-val", "http-rpc").As(client.RemoteConnection, client.HTTP)
	if err != nil {
		return err
	}
	netInternalWS, err := e.Fwd.FindPort("sol:0", "sol-val", "ws-rpc").As(client.RemoteConnection, client.WS)
	if err != nil {
		return err
	}
	e.URLs[m.Props.NetworkName] = []string{netLocal, netLocalWS}
	if e.Cfg.InsideK8s {
		e.URLs[m.Props.NetworkName] = []string{netInternal, netInternalWS}
	}
	log.Info().Str("Name", m.Props.NetworkName).Str("URLs", netLocal).Msg("Solana network")
	return nil
}

func defaultProps() *Props {
	return &Props{
		NetworkName: "sol",
		Values: map[string]interface{}{
			"replicas": "1",
			"sol": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   "solanalabs/solana",
					"version": "v1.13.3",
				},
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "2000m",
						"memory": "4000Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "2000m",
						"memory": "4000Mi",
					},
				},
			},
		},
	}
}

func New(props *Props) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props *Props) environment.ConnectedChart {
	if props == nil {
		props = defaultProps()
	}
	return Chart{
		HelmProps: &HelmProps{
			Name:    "sol",
			Path:    "chainlink-qa/solana-validator",
			Values:  &props.Values,
			Version: helmVersion,
		},
		Props: props,
	}
}
