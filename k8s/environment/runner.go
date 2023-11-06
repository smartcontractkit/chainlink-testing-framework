package environment

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/imports/k8s"
	a "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
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
		utils.Ptr(fmt.Sprintf("%s-role", props.BaseName)),
		&k8s.KubeRoleProps{
			Metadata: &k8s.ObjectMeta{
				Name: utils.Ptr(props.BaseName),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: &[]*string{
						utils.Ptr(""), // this empty line is needed or k8s get really angry
						utils.Ptr("apps"),
						utils.Ptr("batch"),
						utils.Ptr("core"),
						utils.Ptr("networking.k8s.io"),
						utils.Ptr("storage.k8s.io"),
						utils.Ptr("policy"),
						utils.Ptr("chaos-mesh.org"),
						utils.Ptr("monitoring.coreos.com"),
						utils.Ptr("rbac.authorization.k8s.io"),
					},
					Resources: &[]*string{
						utils.Ptr("*"),
					},
					Verbs: &[]*string{
						utils.Ptr("*"),
					},
				},
			},
		})
	k8s.NewKubeRoleBinding(
		chart,
		utils.Ptr(fmt.Sprintf("%s-role-binding", props.BaseName)),
		&k8s.KubeRoleBindingProps{
			RoleRef: &k8s.RoleRef{
				ApiGroup: utils.Ptr("rbac.authorization.k8s.io"),
				Kind:     utils.Ptr("Role"),
				Name:     utils.Ptr("remote-test-runner"),
			},
			Metadata: nil,
			Subjects: &[]*k8s.Subject{
				{
					Kind:      utils.Ptr("ServiceAccount"),
					Name:      utils.Ptr("default"),
					Namespace: utils.Ptr(props.TargetNamespace),
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
		utils.Ptr(fmt.Sprintf("%s-job", props.BaseName)),
		&k8s.KubeJobProps{
			Metadata: &k8s.ObjectMeta{
				Name: utils.Ptr(props.BaseName),
			},
			Spec: &k8s.JobSpec{
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: a.ConvertAnnotations(defaultRunnerPodAnnotations),
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: utils.Ptr("default"),
						Containers: &[]*k8s.Container{
							container(props),
						},
						RestartPolicy: utils.Ptr(restartPolicy),
						Volumes: &[]*k8s.Volume{
							{
								Name:     utils.Ptr("persistence"),
								EmptyDir: &k8s.EmptyDirVolumeSource{},
							},
						},
					},
				},
				ActiveDeadlineSeconds: nil,
				BackoffLimit:          utils.Ptr(backOffLimit),
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
		Name:            utils.Ptr(fmt.Sprintf("%s-node", props.BaseName)),
		Image:           utils.Ptr(props.Image),
		ImagePullPolicy: utils.Ptr("Always"),
		Env:             jobEnvVars(props),
		Resources:       a.ContainerResources(cpu, mem, cpu, mem),
		VolumeMounts: &[]*k8s.VolumeMount{
			{
				Name:      utils.Ptr("persistence"),
				MountPath: utils.Ptr("/persistence"),
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
		config.EnvVarLocalCharts,
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
