package polygon_edge

import (
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
)

type Props struct {
	NetworkName string   `envconfig:"network_name"`
	Simulated   bool     `envconfig:"network_simulated"`
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
	return m.Props.Simulated
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
	if m.Props.Simulated {
		gethLocalWs, err := e.Fwd.FindPort("polygon-edge:0", "polygon-edge", "http").As(client.LocalConnection, client.WSSUFFIX)
		if err != nil {
			return err
		}
		gethInternalWs, err := e.Fwd.FindPort("polygon-edge:0", "polygon-edge", "http").As(client.RemoteConnection, client.WS)
		if err != nil {
			return err
		}
		if e.Cfg.InsideK8s {
			e.URLs[m.Props.NetworkName] = []string{gethInternalWs}
		} else {
			e.URLs[m.Props.NetworkName] = []string{gethLocalWs}
		}

		// For cases like starknet we need the internalHttp address to set up the L1<>L2 interaction
		e.URLs[m.Props.NetworkName+"_internal"] = []string{gethInternalWs}

		log.Info().Str("Name", "Edge").Str("URLs", gethLocalWs).Msg("Edge network")
	} else {
		e.URLs[m.Props.NetworkName] = m.Props.WsURLs
		log.Info().Str("Name", m.Props.NetworkName).Strs("URLs", m.Props.WsURLs).Msg("Edge network")
	}
	return nil
}

func defaultProps() *Props {
	return &Props{
		NetworkName: "edge",
		Simulated:   true,
		Values:      map[string]interface{}{},
	}
}

func New(props *Props) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props *Props) environment.ConnectedChart {
	targetProps := defaultProps()
	if props == nil {
		props = targetProps
	}
	config.MustMerge(targetProps, props)
	config.MustMerge(&targetProps.Values, props.Values)
	targetProps.Simulated = props.Simulated // Mergo has issues with boolean merging for simulated networks
	if targetProps.Simulated {
		return Chart{
			HelmProps: &HelmProps{
				Name:    "edge",
				Path:    "chainlink-qa/polygon-edge",
				Values:  &targetProps.Values,
				Version: helmVersion,
			},
			Props: targetProps,
		}
	}
	return Chart{
		Props: targetProps,
	}
}
