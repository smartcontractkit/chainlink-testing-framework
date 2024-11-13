package wasp

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHelmDeployTimeoutSec = "10m"
	defaultArchiveName          = "wasp-0.1.8.tgz"
	defaultDockerfilePath       = "DockerfileWasp"
	defaultDockerfileIgnorePath = "DockerfileWasp.dockerignore"
	defaultBuildScriptPath      = "./build.sh"
)

// k8s pods resources
const (
	DefaultRequestsCPU    = "1000m"
	DefaultRequestsMemory = "512Mi"
	DefaultLimitsCPU      = "1000m"
	DefaultLimitsMemory   = "512Mi"
)

//go:embed charts/wasp/wasp-0.1.8.tgz
var defaultChart []byte

//go:embed DockerfileWasp
var DefaultDockerfile []byte

//go:embed DockerfileWasp.dockerignore
var DefaultDockerIgnorefile []byte

// TODO: remove the whole stuff in favor of remote-runner
var DefaultBuildScript []byte

var (
	ErrNoNamespace = errors.New("namespace is empty")
	ErrNoJobs      = errors.New("HelmValues should contain \"jobs\" field used to scale your cluster jobs, jobs must be > 0")
)

// ClusterConfig defines k8s jobs settings
type ClusterConfig struct {
	ChartPath            string
	Namespace            string
	KeepJobs             bool
	UpdateImage          bool
	DockerCmdExecPath    string
	DockerfilePath       string
	DockerIgnoreFilePath string
	BuildScriptPath      string
	BuildCtxPath         string
	ImageTag             string
	RegistryName         string
	RepoName             string
	HelmDeployTimeoutSec string
	HelmValues           map[string]string
	// generated values
	tmpHelmFilePath string
}

// Defaults sets default values for the ClusterConfig fields if they are not already set.
// It initializes Helm values, file paths, and resource requests and limits with default values.
// It returns an error if any file operations fail during the setup process.
func (m *ClusterConfig) Defaults(a int) error {
	// TODO: will it be more clear if we move Helm values to a struct
	// TODO: or should it be like that for extensibility of a chart without reflection?
	m.HelmValues["namespace"] = m.Namespace
	// nolint
	m.HelmValues["sync"] = fmt.Sprintf("a%s", uuid.NewString()[0:5])
	if m.HelmDeployTimeoutSec == "" {
		m.HelmDeployTimeoutSec = defaultHelmDeployTimeoutSec
	}
	if m.HelmValues["test.timeout"] == "" {
		m.HelmValues["test.timeout"] = "12h"
	}
	if m.HelmValues["resources.requests.cpu"] == "" {
		m.HelmValues["resources.requests.cpu"] = DefaultRequestsCPU
	}
	if m.HelmValues["resources.requests.memory"] == "" {
		m.HelmValues["resources.requests.memory"] = DefaultRequestsMemory
	}
	if m.HelmValues["resources.limits.cpu"] == "" {
		m.HelmValues["resources.limits.cpu"] = DefaultLimitsCPU
	}
	if m.HelmValues["resources.limits.memory"] == "" {
		m.HelmValues["resources.limits.memory"] = DefaultLimitsMemory
	}
	if m.ChartPath == "" {
		log.Info().Msg("Using default embedded chart")
		if err := os.WriteFile(defaultArchiveName, defaultChart, os.ModePerm); err != nil {
			return err
		}
		m.tmpHelmFilePath, m.ChartPath = defaultArchiveName, defaultArchiveName
	}
	if m.DockerfilePath == "" {
		log.Info().Msg("Using default embedded DockerfileWasp")
		if err := os.WriteFile(defaultDockerfilePath, DefaultDockerfile, os.ModePerm); err != nil {
			return err
		}
		p, err := filepath.Abs(defaultDockerfilePath)
		if err != nil {
			return err
		}
		m.DockerfilePath = p
	}
	if m.DockerIgnoreFilePath == "" {
		log.Info().Msg("Using default embedded DockerfileWasp.dockerignore")
		if err := os.WriteFile(defaultDockerfileIgnorePath, DefaultDockerIgnorefile, os.ModePerm); err != nil {
			return err
		}
		p, err := filepath.Abs(defaultDockerfileIgnorePath)
		if err != nil {
			return err
		}
		m.DockerIgnoreFilePath = p
	}
	if m.BuildScriptPath == "" {
		log.Info().Msg("Using default build script")
		fname := strings.Replace(defaultBuildScriptPath, "./", "", -1)
		if err := os.WriteFile(fname, DefaultBuildScript, os.ModePerm); err != nil {
			return err
		}
		m.BuildScriptPath = defaultBuildScriptPath
	}
	return nil
}

// Validate checks the ClusterConfig for required fields and returns an error if any are missing.
// Specifically, it ensures that the Namespace and HelmValues["jobs"] fields are not empty.
func (m *ClusterConfig) Validate() (err error) {
	if m.Namespace == "" {
		err = errors.Join(err, ErrNoNamespace)
	}
	if m.HelmValues["jobs"] == "" {
		err = errors.Join(err, ErrNoJobs)
	}
	return
}

// parseECRImageURI parses an Amazon ECR image URI into its constituent parts: registry, repository, and tag.
// It returns an error if the URI does not match the expected format `${registry}/${repo}:${tag}`.
func parseECRImageURI(uri string) (registry, repo, tag string, err error) {
	re := regexp.MustCompile(`^([^/]+)/([^:]+):(.+)$`)
	matches := re.FindStringSubmatch(uri)
	if len(matches) != 4 {
		return "", "", "", fmt.Errorf("invalid ECR image URI format, must be ${registry}/${repo}:${tag}")
	}
	return matches[1], matches[2], matches[3], nil
}

// ClusterProfile is a k8s cluster test for some workload profile
type ClusterProfile struct {
	cfg    *ClusterConfig
	c      *K8sClient
	Ctx    context.Context
	Cancel context.CancelFunc
}

// NewClusterProfile creates a new ClusterProfile using the provided ClusterConfig.
// It validates and applies default settings to the configuration. It logs the configuration
// and sets a context with a timeout based on the test timeout duration specified in the config.
// If UpdateImage is true, it builds and pushes the image. It returns the ClusterProfile and
// any error encountered during the process.
func NewClusterProfile(cfg *ClusterConfig) (*ClusterProfile, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.Defaults(0); err != nil {
		return nil, err
	}
	log.Info().Interface("Config", cfg).Msg("Cluster configuration")
	dur, err := time.ParseDuration(cfg.HelmValues["test.timeout"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse test timeout duration")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), dur)
	cp := &ClusterProfile{
		cfg:    cfg,
		c:      NewK8sClient(),
		Ctx:    ctx,
		Cancel: cancelFunc,
	}
	if cp.cfg.UpdateImage {
		return cp, cp.buildAndPushImage()
	}
	return cp, nil
}

// buildAndPushImage parses the ECR image URI from the configuration and constructs a command to build and push a Docker image. It returns an error if the URI parsing or command execution fails.
func (m *ClusterProfile) buildAndPushImage() error {
	registry, repo, tag, err := parseECRImageURI(m.cfg.HelmValues["image"])
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("%s %s %s %s %s %s %s",
		m.cfg.BuildScriptPath,
		m.cfg.DockerfilePath,
		m.cfg.BuildCtxPath,
		tag,
		registry,
		repo,
		m.cfg.DockerCmdExecPath,
	)
	log.Info().Str("Cmd", cmd).Msg("Building docker")
	return ExecCmd(cmd)
}

// deployHelm installs a Helm chart using the specified testName as the release name.
// It constructs the Helm command with the chart path, values, namespace, and timeout
// from the ClusterProfile configuration. It returns any error encountered during the
// execution of the Helm command.
func (m *ClusterProfile) deployHelm(testName string) error {
	//nolint
	defer os.Remove(m.cfg.tmpHelmFilePath)
	var cmd strings.Builder
	cmd.WriteString(fmt.Sprintf("helm install %s %s", testName, m.cfg.ChartPath))
	for k, v := range m.cfg.HelmValues {
		cmd.WriteString(fmt.Sprintf(" --set %s=%s", k, v))
	}
	cmd.WriteString(fmt.Sprintf(" -n %s", m.cfg.Namespace))
	cmd.WriteString(fmt.Sprintf(" --timeout %s", m.cfg.HelmDeployTimeoutSec))
	log.Info().Str("Cmd", cmd.String()).Msg("Deploying jobs")
	return ExecCmd(cmd.String())
}

// Run generates a unique test name, modifies it to comply with Helm naming rules,
// and deploys a Helm chart using this name. It retrieves the number of jobs from
// the Helm values and tracks these jobs within the specified namespace. It returns
// an error if any step in the process fails.
func (m *ClusterProfile) Run() error {
	testName := uuid.NewString()[0:8]
	tn := []rune(testName)
	// replace first letter, since helm does not allow it to start with numbers
	tn[0] = 'a'
	if err := m.deployHelm(string(tn)); err != nil {
		return err
	}
	jobNum, err := strconv.Atoi(m.cfg.HelmValues["jobs"])
	if err != nil {
		return err
	}
	return m.c.TrackJobs(m.Ctx, m.cfg.Namespace, m.cfg.HelmValues["sync"], jobNum, m.cfg.KeepJobs)
}
