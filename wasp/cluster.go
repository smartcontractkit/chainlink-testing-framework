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

// Defaults initializes default values for the Helm configuration and related paths in the ClusterConfig. 
// It sets the namespace, sync identifier, deployment timeout, resource requests and limits, 
// and ensures that default chart, Dockerfile, Docker ignore file, and build script are created 
// if their respective paths are not provided. 
// If any errors occur during file writing or path resolution, it returns the error encountered. 
// This function is typically called during the creation of a new ClusterProfile to ensure 
// that the configuration is complete and valid before proceeding with further operations.
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

// Validate checks the ClusterConfig for required fields. 
// It ensures that the Namespace is not empty and that the "jobs" key in HelmValues is set. 
// If any of these validations fail, it returns an error detailing the missing fields. 
// If all validations pass, it returns nil, indicating that the configuration is valid.
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
// It expects the URI to be in the format ${registry}/${repo}:${tag}. 
// If the format is invalid, it returns an error indicating the expected format. 
// On success, it returns the extracted registry, repository, and tag as strings, along with a nil error.
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

// NewClusterProfile creates a new ClusterProfile instance based on the provided ClusterConfig. 
// It validates the configuration and sets default values before initializing the ClusterProfile. 
// If the configuration includes an update for the image, it builds and pushes the image as part of the initialization process. 
// The function returns a pointer to the newly created ClusterProfile and an error if any issues occur during the process. 
// If the configuration is invalid or if there are errors in parsing the timeout duration, an appropriate error is returned.
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

// buildAndPushImage builds a Docker image using the specified configuration and pushes it to a container registry. 
// It constructs the command to execute based on the provided build script, Dockerfile, and context path, 
// along with the image tag, registry, and repository details parsed from the image URI. 
// If any errors occur during the parsing of the image URI or the execution of the build command, 
// the function returns the corresponding error. 
// On successful execution, it returns nil.
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

// deployHelm installs a Helm chart using the specified test name and configuration settings. 
// It constructs a command to execute the Helm install operation, incorporating any Helm values 
// provided in the configuration. The function also ensures that the temporary Helm file is 
// removed after execution. If the command execution fails, it returns an error indicating 
// the failure of the deployment process.
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

// Run executes the deployment of a Helm chart and tracks the associated jobs. 
// It generates a unique test name, ensuring it starts with a letter, and then 
// calls the deployHelm method to initiate the deployment. If the deployment 
// is successful, it retrieves the number of jobs from the configuration and 
// invokes the TrackJobs method to monitor the job status. 
// The function returns an error if any step in the process fails.
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
