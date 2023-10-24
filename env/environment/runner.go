package environment

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/imports/k8s"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
)

const REMOTE_RUNNER_NAME = "remote-test-runner"

type Chart struct {
	Props *Props
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return REMOTE_RUNNER_NAME
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

func (m Chart) ExportData(e *Environment) error {
	return nil
}

func NewRunner(props *Props) func(root cdk8s.Chart) ConnectedChart {
	return func(root cdk8s.Chart) ConnectedChart {
		c := &Chart{
			Props: props,
		}
		role(root, props)
		job(root, props)
		return c
	}
}

type Props struct {
	BaseName           string
	TargetNamespace    string
	Labels             *map[string]*string
	Image              string
	TestName           string
	NoManifestUpdate   bool
	PreventPodEviction bool
}

func role(chart cdk8s.Chart, props *Props) {
	k8s.NewKubeRole(
		chart,
		a.Str(fmt.Sprintf("%s-role", props.BaseName)),
		&k8s.KubeRoleProps{
			Metadata: &k8s.ObjectMeta{
				Name: a.Str(props.BaseName),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: &[]*string{
						a.Str(""), // this empty line is needed or k8s get really angry
						a.Str("apps"),
						a.Str("batch"),
						a.Str("core"),
						a.Str("networking.k8s.io"),
						a.Str("storage.k8s.io"),
						a.Str("policy"),
						a.Str("chaos-mesh.org"),
						a.Str("monitoring.coreos.com"),
						a.Str("rbac.authorization.k8s.io"),
					},
					Resources: &[]*string{
						a.Str("*"),
					},
					Verbs: &[]*string{
						a.Str("*"),
					},
				},
			},
		})
	k8s.NewKubeRoleBinding(
		chart,
		a.Str(fmt.Sprintf("%s-role-binding", props.BaseName)),
		&k8s.KubeRoleBindingProps{
			RoleRef: &k8s.RoleRef{
				ApiGroup: a.Str("rbac.authorization.k8s.io"),
				Kind:     a.Str("Role"),
				Name:     a.Str("remote-test-runner"),
			},
			Metadata: nil,
			Subjects: &[]*k8s.Subject{
				{
					Kind:      a.Str("ServiceAccount"),
					Name:      a.Str("default"),
					Namespace: a.Str(props.TargetNamespace),
				},
			},
		},
	)
}

func job(chart cdk8s.Chart, props *Props) {
	defaultRunnerPodAnnotations := markNotSafeToEvict(props.PreventPodEviction, nil)
	restartPolicy := "Never"
	backOffLimit := float64(0)
	if os.Getenv(config.EnvVarDetachRunner) == "true" { // If we're running detached, we're likely running a long-form test
		restartPolicy = "OnFailure"
		backOffLimit = 100000 // effectively infinite (I hope)
	}
	k8s.NewKubeJob(
		chart,
		a.Str(fmt.Sprintf("%s-job", props.BaseName)),
		&k8s.KubeJobProps{
			Metadata: &k8s.ObjectMeta{
				Name: a.Str(props.BaseName),
			},
			Spec: &k8s.JobSpec{
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: a.ConvertAnnotations(defaultRunnerPodAnnotations),
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: a.Str("default"),
						Containers: &[]*k8s.Container{
							container(props),
						},
						RestartPolicy: a.Str(restartPolicy),
						Volumes: &[]*k8s.Volume{
							{
								Name:     a.Str("persistence"),
								EmptyDir: &k8s.EmptyDirVolumeSource{},
							},
						},
					},
				},
				ActiveDeadlineSeconds: nil,
				BackoffLimit:          a.Num(backOffLimit),
			},
		})
}

func container(props *Props) *k8s.Container {
	cpu := os.Getenv(config.EnvVarRemoteRunnerCpu)
	if cpu == "" {
		cpu = "1000m"
	}
	mem := os.Getenv(config.EnvVarRemoteRunnerMem)
	if mem == "" {
		mem = "1024Mi"
	}
	return &k8s.Container{
		Name:            a.Str(fmt.Sprintf("%s-node", props.BaseName)),
		Image:           a.Str(props.Image),
		ImagePullPolicy: a.Str("Always"),
		Env:             jobEnvVars(props),
		Resources:       a.ContainerResources(cpu, mem, cpu, mem),
		VolumeMounts: &[]*k8s.VolumeMount{
			{
				Name:      a.Str("persistence"),
				MountPath: a.Str("/persistence"),
			},
		},
	}
}

func jobEnvVars(props *Props) *[]*k8s.EnvVar {
	// Use a map to set values so we can easily overwrite duplicate values
	env := make(map[string]string)

	// Propagate common environment variables to the runner
	lookups := []string{
		config.EnvVarCLImage,
		config.EnvVarCLTag,
		config.EnvVarCLCommitSha,
		config.EnvVarLogLevel,
		config.EnvVarTestTrigger,
		config.EnvVarToleration,
		config.EnvVarSlackChannel,
		config.EnvVarSlackKey,
		config.EnvVarSlackUser,
		config.EnvVarPyroscopeKey,
		config.EnvVarPyroscopeEnvironment,
		config.EnvVarPyroscopeServer,
		config.EnvVarUser,
		config.EnvVarNodeSelector,
		config.EnvVarSelectedNetworks,
		config.EnvVarDBURL,
		config.EnvVarEVMKeys,
		config.EnvVarInternalDockerRepo,
		config.EnvVarEVMUrls,
		config.EnvVarEVMHttpUrls,
	}
	for _, k := range lookups {
		v, success := os.LookupEnv(k)
		if success && len(v) > 0 {
			log.Debug().Str(k, v).Msg("Forwarding Env Var")
			env[k] = v
		}
	}

	// Propagate prefixed variables to the runner
	// These should overwrite anything that was unprefixed if they match up
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			if strings.HasPrefix(e[:i], config.EnvVarPrefix) {
				withoutPrefix := strings.Replace(e[:i], config.EnvVarPrefix, "", 1)
				env[withoutPrefix] = e[i+1:]
			}
		}
	}

	// Add variables that should need specific values for thre remote runner
	env[config.EnvVarNamespace] = props.TargetNamespace
	env["TEST_NAME"] = props.TestName
	env[config.EnvVarInsideK8s] = "true"
	env[config.EnvVarNoManifestUpdate] = strconv.FormatBool(props.NoManifestUpdate)

	// convert from map to the expected array
	cdk8sVars := make([]*k8s.EnvVar, 0)
	for k, v := range env {
		cdk8sVars = append(cdk8sVars, a.EnvVarStr(k, v))
	}
	return &cdk8sVars
}
