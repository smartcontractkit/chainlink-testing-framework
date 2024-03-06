package runner

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/k8s_client"
)

const (
	defaultHelmDeployTimeoutSec = "10m"
)

// k8s pods resources
const (
	DefaultRequestsCPU    = "1000m"
	DefaultRequestsMemory = "512Mi"
	DefaultLimitsCPU      = "1000m"
	DefaultLimitsMemory   = "512Mi"
)

var DefaultDockerfile []byte

var DefaultDockerIgnorefile []byte

var (
	ErrNoNamespace = errors.New("namespace is empty")
	ErrNoJobs      = errors.New("HelmValues should contain \"jobs\" field used to scale your cluster jobs, jobs must be > 0")
)

// Config defines k8s jobs settings
type Config struct {
	ChartPath            string
	Namespace            string
	KeepJobs             bool
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

func (m *Config) Defaults() error {
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
	return nil
}

func (m *Config) Validate() (err error) {
	if m.Namespace == "" {
		err = errors.Join(err, ErrNoNamespace)
	}
	if m.HelmValues["jobs"] == "" {
		err = errors.Join(err, ErrNoJobs)
	}
	return
}

type K8sTestRun struct {
	cfg    *Config
	c      *k8s_client.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewK8sTestRun(cfg *Config) (*K8sTestRun, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.Defaults(); err != nil {
		return nil, err
	}
	log.Info().Interface("Config", cfg).Msg("Cluster configuration")
	dur, err := time.ParseDuration(cfg.HelmValues["test.timeout"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse test timeout duration")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), dur)
	cp := &K8sTestRun{
		cfg:    cfg,
		c:      k8s_client.NewClient(),
		Ctx:    ctx,
		Cancel: cancelFunc,
	}
	return cp, nil
}

func (m *K8sTestRun) getChartPath() string {
	if m.cfg.ChartPath == "" {
		// Use default chart
		_, f, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(f))
		chartPath := path.Join(filepath.Dir(dir), "../chart")
		return chartPath
	}
	return m.cfg.ChartPath
}

func (m *K8sTestRun) deployHelm(testName string) error {
	//nolint
	defer os.Remove(m.cfg.tmpHelmFilePath)
	var cmd strings.Builder
	cmd.WriteString(fmt.Sprintf("helm install %s %s", testName, m.getChartPath()))
	for k, v := range m.cfg.HelmValues {
		cmd.WriteString(fmt.Sprintf(" --set %s=%s", k, v))
	}
	cmd.WriteString(fmt.Sprintf(" -n %s", m.cfg.Namespace))
	cmd.WriteString(fmt.Sprintf(" --timeout %s", m.cfg.HelmDeployTimeoutSec))
	log.Info().Str("Cmd", cmd.String()).Msg("Deploying jobs")
	return ExecCmd(cmd.String())
}

// Run starts a new test
func (m *K8sTestRun) Run() error {
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

// ExecCmd executes os command, logging both streams
func ExecCmd(command string) error {
	return ExecCmdWithStreamFunc(command, func(m string) {
		log.Info().Str("Text", m).Msg("Command output")
	})
}

// readStdPipe continuously read a pipe from the command
func readStdPipe(pipe io.ReadCloser, streamFunc func(string)) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if streamFunc != nil {
			streamFunc(m)
		}
	}
}

// ExecCmdWithStreamFunc executes command with stream function
func ExecCmdWithStreamFunc(command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
	cmd := exec.Command(c[0], c[1:]...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go readStdPipe(stderr, outputFunction)
	go readStdPipe(stdout, outputFunction)
	return cmd.Wait()
}
