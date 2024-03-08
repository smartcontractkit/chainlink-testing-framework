package runner

import (
	"context"
	_ "embed"
	"fmt"
	"path"
	"runtime"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/config"
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
	settings := cli.New() // Initialize the Helm environment settings
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), m.cfg.Namespace, "", log.Printf); err != nil {
		return fmt.Errorf("failed to initialize Helm action configuration: %w", err)
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = testName
	install.Namespace = m.cfg.Namespace
	install.Timeout = time.Minute * 3

	chart, err := loader.Load(m.getChartPath())
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	// Install helm chart with overrides
	log.Info().Interface("Overrides", m.ChartOverrides).Msg("Installing chart with overrides")
	release, err := install.Run(chart, m.ChartOverrides)
	if err != nil {
		return fmt.Errorf("failed to install chart with overrides: %w", err)
	}
	log.Info().Str("releaseName", release.Name).Str("namespace", release.Namespace).Msg("Chart installed successfully")

	return nil
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
	err := m.c.WaitUntilJobsComplete(m.Ctx, m.cfg.Namespace, m.cfg.SyncValue, m.cfg.JobCount)
	m.c.PrintPodLogs(m.Ctx, m.cfg.Namespace, m.cfg.SyncValue)
	if !m.cfg.KeepJobs {
		err = m.c.RemoveJobs(m.Ctx, m.cfg.Namespace, m.cfg.SyncValue)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove jobs")
		}
	}
	return err
}
