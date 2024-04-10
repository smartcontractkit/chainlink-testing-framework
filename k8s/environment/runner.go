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
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

const REMOTE_RUNNER_NAME = "remote-test-runner"

type Chart struct {
	Props *Props
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return m.Props.BaseName
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
	ReportPath         string
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
		ptr.Ptr(fmt.Sprintf("%s-role", props.BaseName)),
		&k8s.KubeRoleProps{
			Metadata: &k8s.ObjectMeta{
				Name: ptr.Ptr(props.BaseName),
			},
			Rules: &[]*k8s.PolicyRule{
				{
					ApiGroups: &[]*string{
						ptr.Ptr(""), // this empty line is needed or k8s get really angry
						ptr.Ptr("apps"),
						ptr.Ptr("batch"),
						ptr.Ptr("core"),
						ptr.Ptr("networking.k8s.io"),
						ptr.Ptr("storage.k8s.io"),
						ptr.Ptr("policy"),
						ptr.Ptr("chaos-mesh.org"),
						ptr.Ptr("monitoring.coreos.com"),
						ptr.Ptr("rbac.authorization.k8s.io"),
					},
					Resources: &[]*string{
						ptr.Ptr("*"),
					},
					Verbs: &[]*string{
						ptr.Ptr("*"),
					},
				},
			},
		})
	k8s.NewKubeRoleBinding(
		chart,
		ptr.Ptr(fmt.Sprintf("%s-role-binding", props.BaseName)),
		&k8s.KubeRoleBindingProps{
			RoleRef: &k8s.RoleRef{
				ApiGroup: ptr.Ptr("rbac.authorization.k8s.io"),
				Kind:     ptr.Ptr("Role"),
				Name:     ptr.Ptr("remote-test-runner"),
			},
			Metadata: nil,
			Subjects: &[]*k8s.Subject{
				{
					Kind:      ptr.Ptr("ServiceAccount"),
					Name:      ptr.Ptr("default"),
					Namespace: ptr.Ptr(props.TargetNamespace),
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
		ptr.Ptr(fmt.Sprintf("%s-job", props.BaseName)),
		&k8s.KubeJobProps{
			Metadata: &k8s.ObjectMeta{
				Name: ptr.Ptr(props.BaseName),
			},
			Spec: &k8s.JobSpec{
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: a.ConvertAnnotations(defaultRunnerPodAnnotations),
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: ptr.Ptr("default"),
						Containers:         container(props),
						RestartPolicy:      ptr.Ptr(restartPolicy),
						Volumes: &[]*k8s.Volume{
							{
								Name:     ptr.Ptr("persistence"),
								EmptyDir: &k8s.EmptyDirVolumeSource{},
							},
							{
								Name:     ptr.Ptr("reports"),
								EmptyDir: &k8s.EmptyDirVolumeSource{},
							},
						},
					},
				},
				ActiveDeadlineSeconds: nil,
				BackoffLimit:          ptr.Ptr(backOffLimit),
			},
		})
}

func container(props *Props) *[]*k8s.Container {
	cpu := os.Getenv(config.EnvVarRemoteRunnerCpu)
	if cpu == "" {
		cpu = "1000m"
	}
	mem := os.Getenv(config.EnvVarRemoteRunnerMem)
	if mem == "" {
		mem = "1024Mi"
	}
	return ptr.Ptr([]*k8s.Container{
		{
			Name:            ptr.Ptr(fmt.Sprintf("%s-node", props.BaseName)),
			Image:           ptr.Ptr(props.Image),
			ImagePullPolicy: ptr.Ptr("Always"),
			Env:             jobEnvVars(props),
			Resources:       a.ContainerResources(cpu, mem, cpu, mem),
			VolumeMounts: &[]*k8s.VolumeMount{
				{
					Name:      ptr.Ptr("persistence"),
					MountPath: ptr.Ptr("/persistence"),
				},
				{
					Name:      ptr.Ptr("reports"),
					MountPath: ptr.Ptr(fmt.Sprintf("/go/testdir/integration-tests/%s", props.ReportPath)),
					SubPath:   ptr.Ptr(props.ReportPath),
				},
			},
		},
		// we create this container to share same volume as remote-runner-node container. This container
		// keeps on running and stays alive after the remote-runner-node gets completed, so that
		// the calling test can access all files generated by remote runner.
		{
			Name:            ptr.Ptr(fmt.Sprintf("%s-data-files", props.BaseName)),
			Image:           ptr.Ptr("busybox:stable"),
			ImagePullPolicy: ptr.Ptr("Always"),
			Command:         ptr.Ptr([]*string{ptr.Ptr("/bin/sh"), ptr.Ptr("-ec"), ptr.Ptr("while :; do echo '.'; sleep 5 ; done")}),
			Ports: ptr.Ptr([]*k8s.ContainerPort{
				{
					ContainerPort: ptr.Ptr(float64(80)),
				},
			}),
			VolumeMounts: &[]*k8s.VolumeMount{
				{
					Name:      ptr.Ptr("reports"),
					MountPath: ptr.Ptr("reports"),
				},
			},
		},
	})
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
		config.EnvVarUser,
		config.EnvVarNodeSelector,
		config.EnvVarDBURL,
		config.EnvVarInternalDockerRepo,
		config.EnvVarLocalCharts,
		config.EnvBase64ConfigOverride,
		config.EnvBase64NetworkConfig,
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
				log.Debug().Str(e[:i], e[i+1:]).Msg("Forwarding generic Env Var")
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
