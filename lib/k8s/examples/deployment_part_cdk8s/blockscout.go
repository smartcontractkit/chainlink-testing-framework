package deployment_part_cdk8s

import (
	"fmt"

	cdk8s "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/imports/k8s"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg"
	a "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const (
	URLsKey = "blockscout"
)

type Chart struct {
	Props *Props
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return "blockscout"
}

func (m Chart) GetProps() interface{} {
	return m.Props
}

func (m Chart) GetPath() string {
	return ""
}

func (m Chart) GetVersion() string {
	return ""
}

func (m Chart) GetValues() *map[string]interface{} {
	return nil
}

func (m Chart) GetLabels() map[string]string {
	return map[string]string{
		"chain.link/component": "blockscout",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	bsURL, err := e.Fwd.FindPort("blockscout:0", "blockscout-node", "explorer").As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	log.Info().Str("URL", bsURL).Msg("Blockscout explorer")
	e.URLs[URLsKey] = []string{bsURL}
	return nil
}

func New(props *Props) func(root cdk8s.Chart) environment.ConnectedChart {
	return func(root cdk8s.Chart) environment.ConnectedChart {
		dp := defaultProps()
		config.MustMerge(dp, props)
		c := &Chart{
			Props: dp,
		}
		vars := vars{
			Labels: &map[string]*string{
				"app": ptr.Ptr(c.GetName()),
			},
			ConfigMapName: fmt.Sprintf("%s-cm", c.GetName()),
			BaseName:      c.GetName(),
			Port:          4000,
			Props:         dp,
		}
		service(root, vars)
		deployment(root, vars)
		return c
	}
}

type Props struct {
	HttpURL string `envconfig:"http_url"`
	WsURL   string `envconfig:"ws_url"`
}

func defaultProps() *Props {
	return &Props{
		HttpURL: "http://geth:8544",
		WsURL:   "ws://geth:8546",
	}
}

// vars some shared labels/selectors and names that must match in resources
type vars struct {
	Labels        *map[string]*string
	BaseName      string
	ConfigMapName string
	Port          float64
	Props         *Props
}

func service(chart cdk8s.Chart, vars vars) {
	k8s.NewKubeService(chart, ptr.Ptr(fmt.Sprintf("%s-service", vars.BaseName)), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name: ptr.Ptr(vars.BaseName),
		},
		Spec: &k8s.ServiceSpec{
			Ports: &[]*k8s.ServicePort{
				{
					Name:       ptr.Ptr("explorer"),
					Port:       ptr.Ptr(vars.Port),
					TargetPort: k8s.IntOrString_FromNumber(ptr.Ptr[float64](4000)),
				},
			},
			Selector: vars.Labels,
		},
	})
}

func postgresContainer(p vars) *k8s.Container {
	return &k8s.Container{
		Name:  ptr.Ptr(fmt.Sprintf("%s-db", p.BaseName)),
		Image: ptr.Ptr("postgres:13.6"),
		Ports: &[]*k8s.ContainerPort{
			{
				Name:          ptr.Ptr("postgres"),
				ContainerPort: ptr.Ptr[float64](5432),
			},
		},
		Env: &[]*k8s.EnvVar{
			a.EnvVarStr("POSTGRES_PASSWORD", "postgres"),
			a.EnvVarStr("POSTGRES_DB", "blockscout"),
		},
		LivenessProbe: &k8s.Probe{
			Exec: &k8s.ExecAction{
				Command: pkg.PGIsReadyCheck()},
			InitialDelaySeconds: ptr.Ptr[float64](60),
			PeriodSeconds:       ptr.Ptr[float64](60),
		},
		ReadinessProbe: &k8s.Probe{
			Exec: &k8s.ExecAction{
				Command: pkg.PGIsReadyCheck()},
			InitialDelaySeconds: ptr.Ptr[float64](2),
			PeriodSeconds:       ptr.Ptr[float64](2),
		},
		Resources: a.ContainerResources("1000m", "2048Mi", "1000m", "2048Mi"),
	}
}

func deployment(chart cdk8s.Chart, vars vars) {
	k8s.NewKubeDeployment(
		chart,
		ptr.Ptr(fmt.Sprintf("%s-deployment", vars.BaseName)),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name: ptr.Ptr(vars.BaseName),
			},
			Spec: &k8s.DeploymentSpec{
				Selector: &k8s.LabelSelector{
					MatchLabels: vars.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: vars.Labels,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: ptr.Ptr("default"),
						Containers: &[]*k8s.Container{
							container(vars),
							postgresContainer(vars),
						},
					},
				},
			},
		})
}

func container(vars vars) *k8s.Container {
	return &k8s.Container{
		Name:            ptr.Ptr(fmt.Sprintf("%s-node", vars.BaseName)),
		Image:           ptr.Ptr("f4hrenh9it/blockscout:v1"),
		ImagePullPolicy: ptr.Ptr("Always"),
		Command:         &[]*string{ptr.Ptr(`/bin/bash`)},
		Args: &[]*string{
			ptr.Ptr("-c"),
			ptr.Ptr("mix ecto.create && mix ecto.migrate && mix phx.server"),
		},
		Ports: &[]*k8s.ContainerPort{
			{
				Name:          ptr.Ptr("explorer"),
				ContainerPort: ptr.Ptr(vars.Port),
			},
		},
		ReadinessProbe: &k8s.Probe{
			HttpGet: &k8s.HttpGetAction{
				Port: k8s.IntOrString_FromNumber(ptr.Ptr(vars.Port)),
				Path: ptr.Ptr("/"),
			},
			InitialDelaySeconds: ptr.Ptr[float64](20),
			PeriodSeconds:       ptr.Ptr[float64](5),
		},
		Env: &[]*k8s.EnvVar{
			a.EnvVarStr("MIX_ENV", "prod"),
			a.EnvVarStr("ECTO_USE_SSL", "'false'"),
			a.EnvVarStr("COIN", "DAI"),
			a.EnvVarStr("ETHEREUM_JSONRPC_VARIANT", "geth"),
			a.EnvVarStr("ETHEREUM_JSONRPC_HTTP_URL", vars.Props.HttpURL),
			a.EnvVarStr("ETHEREUM_JSONRPC_WS_URL", vars.Props.WsURL),
			a.EnvVarStr("DATABASE_URL", "postgresql://postgres:@localhost:5432/blockscout?ssl=false"),
		},
		Resources: a.ContainerResources("300m", "2048Mi", "300m", "2048Mi"),
	}
}
