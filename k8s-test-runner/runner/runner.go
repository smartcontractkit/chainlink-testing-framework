package runner

import (
	"context"
	// "embed"
	_ "embed"
	"fmt"
	"os"

	"path"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/exec"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/k8s_client"
)

type K8sTestRun struct {
	cfg            *config.Runner
	c              *k8s_client.Client
	Ctx            context.Context
	Cancel         context.CancelFunc
	ChartOverrides map[string]interface{}
}

func NewK8sTestRun(cfg *config.Runner, chartOverrides map[string]interface{}) (*K8sTestRun, error) {
	log.Info().Interface("Config", cfg).Msg("Cluster configuration")
	runTimeout := cfg.TestTimeout + time.Minute*10
	ctx, cancelFunc := context.WithTimeout(context.Background(), runTimeout)

	return &K8sTestRun{
		cfg:            cfg,
		c:              k8s_client.NewClient(),
		Ctx:            ctx,
		Cancel:         cancelFunc,
		ChartOverrides: chartOverrides,
	}, nil
}

func (m *K8sTestRun) getChartPath() string {
	if m.cfg.ChartPath == "" {
		// Use default chart
		_, f, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(f))
		chartPath := path.Join(dir, "../chart")
		return chartPath
	}
	return m.cfg.ChartPath
}

func (m *K8sTestRun) deployHelm(testName string) error {
	overridesYAML, err := yaml.Marshal(m.ChartOverrides)
	if err != nil {
		return fmt.Errorf("failed to convert overrides to YAML: %w", err)
	}

	// Create a temporary file to save the overrides
	tmpFile, err := os.CreateTemp("", "helm-overrides-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Write the YAML data to the temp file
	if _, err = tmpFile.Write(overridesYAML); err != nil {
		return fmt.Errorf("failed to write overrides to temp file: %w", err)
	}

	// Ensure to flush writes to storage
	if err = tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to flush writes to temp file: %w", err)
	}

	// Run helm install command
	cmd := fmt.Sprintf("helm install %s %s --namespace %s --values %s --timeout 30s", testName, m.getChartPath(), m.cfg.Namespace, tmpFile.Name())
	if m.cfg.Debug {
		cmd += " --debug"
	}
	log.Info().Str("cmd", cmd).Msg("Running helm install...")
	return exec.CmdWithStreamFunc(cmd, func(m string) {
		fmt.Printf("%s\n", m)
	})
}

func (m *K8sTestRun) Run() error {
	testName := uuid.NewString()[0:8]
	tn := []rune(testName)
	// replace first letter, since helm does not allow it to start with numbers
	tn[0] = 'a'
	if err := m.deployHelm(string(tn)); err != nil {
		return err
	}
	jobs, err := m.c.ListJobs(m.Ctx, m.cfg.Namespace, m.cfg.SyncValue)
	if err == nil {
		for _, j := range jobs.Items {
			log.Info().Str("job", j.Name).Str("namespace", m.cfg.Namespace).Msg("Job created")
		}
	}
	// Exit early in detached mode
	if m.cfg.DetachedMode {
		log.Info().Msg("Running in detached mode, exiting early")
		return nil
	}
	err = m.c.WaitUntilJobsComplete(m.Ctx, m.cfg.Namespace, m.cfg.SyncValue, m.cfg.JobCount)
	m.c.PrintPodLogs(m.Ctx, m.cfg.Namespace, m.cfg.SyncValue)
	return err
}
