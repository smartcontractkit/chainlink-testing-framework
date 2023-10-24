package goc

import (
	"fmt"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/imports/k8s"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
)

type Chart struct {
	Props *Props
}

func (m *Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return "dummy"
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

func (m Chart) ExportData(e *environment.Environment) error {
	return nil
}

func New() func(root cdk8s.Chart) environment.ConnectedChart {
	return func(root cdk8s.Chart) environment.ConnectedChart {
		c := &Chart{}
		vars := vars{
			Labels: &map[string]*string{
				"app": a.Str(c.GetName()),
			},
			ConfigMapName: fmt.Sprintf("%s-cm", c.GetName()),
			BaseName:      c.GetName(),
			Port:          3000,
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
	k8s.NewKubeService(chart, a.Str(fmt.Sprintf("%s-service", vars.BaseName)), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name: a.Str(vars.BaseName),
		},
		Spec: &k8s.ServiceSpec{
			Ports: &[]*k8s.ServicePort{
				{
					Name:       a.Str("http"),
					Port:       a.Num(vars.Port),
					TargetPort: k8s.IntOrString_FromNumber(a.Num(3000)),
				},
			},
			Selector: vars.Labels,
		},
	})
}

func deployment(chart cdk8s.Chart, vars vars) {
	k8s.NewKubeDeployment(
		chart,
		a.Str(fmt.Sprintf("%s-deployment", vars.BaseName)),
		&k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name: a.Str(vars.BaseName),
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
						ServiceAccountName: a.Str("default"),
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
		Name:            a.Str(vars.BaseName),
		Image:           a.Str("public.ecr.aws/chainlink/goc-target:latest"),
		ImagePullPolicy: a.Str("Always"),
		Ports: &[]*k8s.ContainerPort{
			{
				Name:          a.Str("http"),
				ContainerPort: a.Num(vars.Port),
			},
		},
		ReadinessProbe: &k8s.Probe{
			HttpGet: &k8s.HttpGetAction{
				Port: k8s.IntOrString_FromNumber(a.Num(vars.Port)),
				Path: a.Str("/"),
			},
			InitialDelaySeconds: a.Num(20),
			PeriodSeconds:       a.Num(5),
		},
		Env:       &[]*k8s.EnvVar{},
		Resources: a.ContainerResources("100m", "512Mi", "100m", "512Mi"),
	}
}
