package chainlink

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
)

const (
	AppName              = "chainlink"
	NodesLocalURLsKey    = "chainlink_local"
	NodesInternalURLsKey = "chainlink_internal"
	DBsLocalURLsKey      = "chainlink_db"
)

type Props struct {
	HasReplicas bool
}

type Chart struct {
	Name    string
	Index   int
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

func (m Chart) GetProps() any {
	return m.Props
}

func (m Chart) GetValues() *map[string]any {
	return m.Values
}

func (m Chart) ExportData(e *environment.Environment) error {
	// fetching all apps with label app=chainlink-${deploymentIndex}:${instanceIndex}
	pods, err := e.Fwd.Client.ListPods(e.Cfg.Namespace, fmt.Sprintf("app=%s", m.Name))
	if err != nil {
		return err
	}
	for i := 0; i < len(pods.Items); i++ {
		chainlinkNode := fmt.Sprintf("%s:node-%d", m.Name, i+1)
		pgNode := fmt.Sprintf("%s-postgres:node-%d", m.Name, i+1)
		internalConnection, err := e.Fwd.FindPort(chainlinkNode, "node", "access").As(client.RemoteConnection, client.HTTP)
		if err != nil {
			return err
		}
		if !m.Props.HasReplicas {
			services, err := e.Client.ListServices(e.Cfg.Namespace, fmt.Sprintf("app=%s", m.Name))
			if err != nil {
				return err
			}
			if services != nil && len(services.Items) > 0 {
				internalConnection = fmt.Sprintf("http://%s:6688", services.Items[0].Name)
			}
		}

		var localConnection string
		if e.Cfg.InsideK8s {
			localConnection = internalConnection
		} else {
			localConnection, err = e.Fwd.FindPort(chainlinkNode, "node", "access").As(client.LocalConnection, client.HTTP)
			if err != nil {
				return err
			}
		}
		e.URLs[NodesInternalURLsKey] = append(e.URLs[NodesInternalURLsKey], internalConnection)
		e.URLs[NodesLocalURLsKey] = append(e.URLs[NodesLocalURLsKey], localConnection)

		dbLocalConnection, err := e.Fwd.FindPort(pgNode, "chainlink-db", "postgres").
			As(client.LocalConnection, client.HTTP)
		if err != nil {
			return err
		}
		e.URLs[DBsLocalURLsKey] = append(e.URLs[DBsLocalURLsKey], dbLocalConnection)
		log.Debug().
			Str("Chart Name", m.Name).
			Str("Local IP", localConnection).
			Str("Local DB IP", dbLocalConnection).
			Str("K8s Internal Connection", internalConnection).
			Msg("Chainlink Node Details")

		nodeDetails := &environment.ChainlinkNodeDetail{
			ChartName:  m.Name,
			PodName:    pods.Items[i].Name,
			LocalIP:    localConnection,
			InternalIP: internalConnection,
			DBLocalIP:  dbLocalConnection,
		}
		e.ChainlinkNodeDetails = append(e.ChainlinkNodeDetails, nodeDetails)
	}
	return nil
}

func defaultProps() map[string]any {
	internalRepo := os.Getenv(config.EnvVarInternalDockerRepo)
	chainlinkImage := "public.ecr.aws/chainlink/chainlink"
	postgresImage := "postgres"
	if internalRepo != "" {
		chainlinkImage = fmt.Sprintf("%s/chainlink", internalRepo)
		postgresImage = fmt.Sprintf("%s/postgres", internalRepo)
	}
	env := make(map[string]string)
	pyroscopeAuth := os.Getenv(config.EnvVarPyroscopeKey)
	if pyroscopeAuth != "" {
		env["CL_PYROSCOPE_AUTH_TOKEN"] = pyroscopeAuth
	}
	return map[string]any{
		"replicas": "1",
		"env":      env,
		"chainlink": map[string]any{
			"image": map[string]any{
				"image":   chainlinkImage,
				"version": "develop",
			},
			"web_port": "6688",
			"p2p_port": "6690",
			"resources": map[string]any{
				"requests": map[string]any{
					"cpu":    "350m",
					"memory": "1024Mi",
				},
				"limits": map[string]any{
					"cpu":    "350m",
					"memory": "1024Mi",
				},
			},
		},
		"db": map[string]any{
			"image": map[string]any{
				"image":   postgresImage,
				"version": "11.15",
			},
			"stateful": false,
			"capacity": "1Gi",
			"resources": map[string]any{
				"requests": map[string]any{
					"cpu":    "250m",
					"memory": "256Mi",
				},
				"limits": map[string]any{
					"cpu":    "250m",
					"memory": "256Mi",
				},
			},
		},
	}
}

func New(index int, props map[string]any) environment.ConnectedChart {
	return NewVersioned(index, "", props)
}

// NewVersioned enables you to select a specific helm chart version
func NewVersioned(index int, helmVersion string, props map[string]any) environment.ConnectedChart {
	dp := defaultProps()
	config.MustEnvOverrideVersion(&dp)
	config.MustMerge(&dp, props)
	p := &Props{
		HasReplicas: false,
	}
	if props["replicas"] != nil && props["replicas"] != "1" {
		p.HasReplicas = true
		replicas := props["replicas"].(int)
		var nodesMap []map[string]any
		for i := 0; i < replicas; i++ {
			nodesMap = append(nodesMap, map[string]any{
				"name": fmt.Sprintf("node-%d", i+1),
			})
		}
		dp["nodes"] = nodesMap
	}
	return Chart{
		Index:   index,
		Name:    fmt.Sprintf("%s-%d", AppName, index),
		Path:    "chainlink-qa/chainlink",
		Version: helmVersion,
		Values:  &dp,
		Props:   p,
	}
}
