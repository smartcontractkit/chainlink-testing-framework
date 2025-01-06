package dummy

import (
	"fmt"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/imports/k8s"
	a "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const (
	URLsKey = "goc"
)

type Chart struct {
	Props *Props
}

func (m *Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return "goc"
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
		"chain.link/component": "http_dummy",
	}
}

func (m Chart) ExportData(e *environment.Environment) error {
	url, err := e.Fwd.FindPort(
		fmt.Sprintf("%s:0", m.GetName()),
		m.GetName(), "http").
		As(client.LocalConnection, client.HTTP)
	if err != nil {
		return err
	}
	log.Info().Str("URL", url).Msg("goc coverage service")
	e.URLs[URLsKey] = []string{url}
	return nil
}

func New() func(root cdk8s.Chart) environment.ConnectedChart {
	return func(root cdk8s.Chart) environment.ConnectedChart {
		c := &Chart{}
		vars := vars{
			Labels: &map[string]*string{
				"app":                  ptr.Ptr(c.GetName()),
				"chain.link/component": ptr.Ptr("http_dummy"),
			},
			ConfigMapName: fmt.Sprintf("%s-cm", c.GetName()),
			BaseName:      c.GetName(),
			Port:          7777,
		}
		service(root, vars)
		deployment(root, vars)
		return c
	}
}

type Props struct {
	Name string
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
			Name:   ptr.Ptr(vars.BaseName),
			Labels: vars.Labels,
		},
		Spec: &k8s.ServiceSpec{
			Ports: &[]*k8s.ServicePort{
				{
					Name:       ptr.Ptr("http"),
					Port:       ptr.Ptr(vars.Port),
					TargetPort: k8s.IntOrString_FromNumber(ptr.Ptr[float64](7777)),
				},
			},
			Selector: vars.Labels,
		},
	})
}

func deployment(chart cdk8s.Chart, vars vars) {
	k8s.NewKubeDeployment(
		chart,
		ptr.Ptr(fmt.Sprintf("%s-deployment", vars.BaseName)),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name:   ptr.Ptr(vars.BaseName),
				Labels: vars.Labels,
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
						},
					},
				},
			},
		})
}

func container(vars vars) *k8s.Container {
	return &k8s.Container{
		Name:            ptr.Ptr(vars.BaseName),
		Image:           ptr.Ptr("public.ecr.aws/chainlink/goc:latest"),
		ImagePullPolicy: ptr.Ptr("Always"),
		Ports: &[]*k8s.ContainerPort{
			{
				Name:          ptr.Ptr("http"),
				ContainerPort: ptr.Ptr(vars.Port),
			},
		},
		ReadinessProbe: &k8s.Probe{
			HttpGet: &k8s.HttpGetAction{
				Port: k8s.IntOrString_FromNumber(ptr.Ptr(vars.Port)),
				Path: ptr.Ptr("/v1/cover/list"),
			},
			InitialDelaySeconds: ptr.Ptr[float64](20),
			PeriodSeconds:       ptr.Ptr[float64](5),
		},
		Env:       &[]*k8s.EnvVar{},
		Resources: a.ContainerResources("200m", "512Mi", "200m", "512Mi"),
	}
}
