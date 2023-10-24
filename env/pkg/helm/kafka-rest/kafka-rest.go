package kafka_rest

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
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

func (m Chart) ExportData(e *environment.Environment) error {
	urls := make([]string, 0)
	local, err := e.Fwd.FindPort("cp-kafka-rest:0", "kafka-rest", "http").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	remote, err := e.Fwd.FindPort("cp-kafka-rest:0", "kafka-rest", "http").As(client.RemoteConnection, client.HTTP)
	if err != nil {
		return err
	}
	if e.Cfg.InsideK8s {
		urls = append(urls, remote, remote)
	} else {
		urls = append(urls, local, remote)
	}
	e.URLs["cp-kafka-rest"] = urls
	log.Info().Str("URL", local).Msg("KafkaRest local connection")
	log.Info().Str("URL", remote).Msg("KafkaRest remote connection")
	return nil
}

func defaultProps() map[string]interface{} {
	return map[string]interface{}{}
}

func New(props map[string]interface{}) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props map[string]interface{}) environment.ConnectedChart {
	dp := defaultProps()
	config.MustMerge(&dp, props)
	return Chart{
		Name:    "cp-kafka-rest",
		Path:    "chainlink-qa/kafka-rest",
		Values:  &dp,
		Version: helmVersion,
	}
}
