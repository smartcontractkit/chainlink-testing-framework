package environment

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/imports/k8s"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg"
	a "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const (
	COVERAGE_DIR       string = "cover"
	FAILED_FUND_RETURN string = "FAILED_FUND_RETURN"
	TEST_FAILED        string = "TEST_FAILED"
)

const (
	ErrInvalidOCI string = "OCI chart url should be in format oci://$ECR_URL/$ECR_REGISTRY_NAME/$CHART_NAME:[?$CHART_VERSION], was %s"
	ErrOCIPull    string = "failed to pull OCI repo: %s"
)

var (
	defaultNamespaceAnnotations = map[string]*string{
		"prometheus.io/scrape":                             ptr.Ptr("true"),
		"backyards.banzaicloud.io/image-registry-access":   ptr.Ptr("true"),
		"backyards.banzaicloud.io/public-dockerhub-access": ptr.Ptr("true"),
	}
)

// ConnectedChart interface to interact both with cdk8s apps and helm charts
type ConnectedChart interface {
	// IsDeploymentNeeded
	// true - we deploy/connect and expose environment data
	// false - we are using external environment, but still exposing data
	IsDeploymentNeeded() bool
	// GetName name of the deployed part
	GetName() string
	// GetPath get Helm chart path, repo or local path
	GetPath() string
	// GetVersion gets the chart's version, empty string if none is specified
	GetVersion() string
	// GetProps get code props if it's typed environment
	GetProps() any
	// GetValues get values.yml props as map, if it's Helm
	GetValues() *map[string]any
	// ExportData export deployment part data in the env
	ExportData(e *Environment) error
	// GetLabels returns labels for the component
	GetLabels() map[string]string
}

// Config is an environment common configuration, labels, annotations, connection types, readiness check, etc.
type Config struct {
	// TTL is time to live for the environment, used with kyverno
	TTL time.Duration
	// NamespacePrefix is a static namespace prefix
	NamespacePrefix string
	// Namespace is full namespace name
	Namespace string
	// Labels is a set of labels applied to the namespace in a format of "key=value"
	Labels []string
	// PodLabels is a set of labels applied to every pod in the namespace
	PodLabels map[string]string
	// WorkloadLabels is a set of labels applied to every workload in the namespace
	WorkloadLabels map[string]string
	// PreventPodEviction if true sets a k8s annotation safe-to-evict=false to prevent pods from being evicted
	// Note: This should only be used if your test is completely incapable of handling things like K8s rebalances without failing.
	// If that is the case, it's worth the effort to make your test fault-tolerant soon. The alternative is expensive and infuriating.
	PreventPodEviction bool
	// Allow deployment to nodes with these tolerances
	Tolerations []map[string]string
	// Restrict deployment to only nodes matching a particular node role
	NodeSelector map[string]string
	// ReadyCheckData is settings for readiness probes checks for all deployment components
	// checking that all pods are ready by default with 8 minutes timeout
	//	&client.ReadyCheckData{
	//		ReadinessProbeCheckSelector: "",
	//		Timeout:                     15 * time.Minute,
	//	}
	ReadyCheckData *client.ReadyCheckData
	// DryRun if true, app will just generate a manifest in local dir
	DryRun bool
	// InsideK8s used for long-running soak tests where you connect to env from the inside
	InsideK8s bool
	// SkipManifestUpdate will skip updating the manifest upon connecting to the environment. Should be true if you wish to update the manifest (e.g. upgrade pods)
	SkipManifestUpdate bool
	// KeepConnection keeps connection until interrupted with a signal, useful when prototyping and debugging a new env
	KeepConnection bool
	// RemoveOnInterrupt automatically removes an environment on interrupt
	RemoveOnInterrupt bool
	// UpdateWaitInterval an interval to wait for deployment update started
	UpdateWaitInterval time.Duration

	// Remote Runner Specific Variables //
	// JobImage an image to run environment as a job inside k8s
	JobImage string
	// Specify only if you want remote-runner to start with a specific name
	RunnerName string
	// Specify only if you want to mount reports from test run in remote runner
	ReportPath string
	// JobLogFunction a function that will be run on each log
	JobLogFunction func(*Environment, string)
	// Test the testing library current Test struct
	Test *testing.T
	// jobDeployed used to limit us to 1 remote runner deploy
	jobDeployed bool
	// detachRunner should we detach the remote runner after starting the test
	detachRunner bool
	// fundReturnFailed the status of a fund return
	fundReturnFailed bool
	// Skip validating that all required chain.link labels are present in the final manifest
	SkipRequiredChainLinkLabelsValidation bool
}

func defaultEnvConfig() *Config {
	return &Config{
		TTL:                20 * time.Minute,
		NamespacePrefix:    "chainlink-test-env",
		UpdateWaitInterval: 1 * time.Second,
		ReadyCheckData: &client.ReadyCheckData{
			ReadinessProbeCheckSelector: "",
			Timeout:                     15 * time.Minute,
		},
	}
}

// Environment describes a launched test environment
type Environment struct {
	App                  cdk8s.App
	CurrentManifest      string
	root                 cdk8s.Chart
	Charts               []ConnectedChart  // All connected charts in the
	Cfg                  *Config           // The environment specific config
	Client               *client.K8sClient // Client connecting to the K8s cluster
	Fwd                  *client.Forwarder // Used to forward ports from local machine to the K8s cluster
	Artifacts            *Artifacts
	Chaos                *client.Chaos
	httpClient           *resty.Client
	URLs                 map[string][]string    // General URLs of launched resources. Uses '_local' to delineate forwarded ports
	ChainlinkNodeDetails []*ChainlinkNodeDetail // ChainlinkNodeDetails has convenient details for connecting to chainlink deployments
	err                  error
}

// ChainlinkNodeDetail contains details about a chainlink node deployment
type ChainlinkNodeDetail struct {
	// ChartName details the name of the Helm chart this node uses, handy for modifying deployment values
	// Note: if you are using replicas of the same chart, this will be the same for all nodes
	// Use NewDeployment function for Chainlink nodes to make use of this
	ChartName string
	// PodName is the name of the pod running the chainlink node
	PodName string
	// LocalIP is the URL to connect to the node from the local machine
	LocalIP string
	// InternalIP is the URL to connect to the node from inside the K8s cluster
	InternalIP string
	// DBLocalIP is the URL to connect to the node's database from the local machine
	DBLocalIP string
}

// New creates new environment
func New(cfg *Config) *Environment {
	logging.Init()
	if cfg == nil {
		cfg = &Config{}
	}
	targetCfg := defaultEnvConfig()
	config.MustMerge(targetCfg, cfg)
	ns := os.Getenv(config.EnvVarNamespace)
	if ns != "" {
		cfg.Namespace = ns
	}
	if cfg.Namespace != "" {
		log.Info().Str("Namespace", cfg.Namespace).Msg("Namespace selected")
		targetCfg.Namespace = cfg.Namespace
	} else {
		targetCfg.Namespace = fmt.Sprintf("%s-%s", targetCfg.NamespacePrefix, uuid.NewString()[0:5])
		log.Info().Str("Namespace", targetCfg.Namespace).Msg("Creating new namespace")
	}
	targetCfg.JobImage = os.Getenv(config.EnvVarJobImage)
	targetCfg.detachRunner, _ = strconv.ParseBool(os.Getenv(config.EnvVarDetachRunner))
	targetCfg.InsideK8s, _ = strconv.ParseBool(os.Getenv(config.EnvVarInsideK8s))

	c, err := client.NewK8sClient()
	if err != nil {
		return &Environment{err: err}
	}
	e := &Environment{
		URLs:   make(map[string][]string),
		Charts: make([]ConnectedChart, 0),
		Client: c,
		Cfg:    targetCfg,
		Fwd:    client.NewForwarder(c, targetCfg.KeepConnection),
	}
	arts, err := NewArtifacts(e.Client, e.Cfg.Namespace)
	if err != nil {
		log.Error().Err(err).Msg("failed to create artifacts client")
		return &Environment{err: err}
	}
	e.Artifacts = arts

	config.JSIIGlobalMu.Lock()
	defer config.JSIIGlobalMu.Unlock()
	if err := e.initApp(); err != nil {
		log.Error().Err(err).Msg("failed to apply the initial manifest to create the namespace")
		return &Environment{err: err}
	}
	e.Chaos = client.NewChaos(c, e.Cfg.Namespace)

	// setup test cleanup if this is using a remote runner
	// and not in detached mode
	// and not using an existing environment
	if targetCfg.JobImage != "" && !targetCfg.detachRunner && !targetCfg.SkipManifestUpdate {
		targetCfg.fundReturnFailed = false
		if targetCfg.Test != nil {
			targetCfg.Test.Cleanup(func() {
				err := e.Shutdown()
				require.NoError(targetCfg.Test, err)
			})
		}
	}
	return e
}

var requiredChainLinkNsLabels = []string{"chain.link/team", "chain.link/cost-center", "chain.link/product"}
var requiredChainLinkWorkloadAndPodLabels = append([]string{}, append(requiredChainLinkNsLabels, "chain.link/component")...)

// validateRequiredChainLinkLabels validates whether the namespace, workloads and pods have the required chain.link labels
// and returns an error with a list of missing labels if any
func (m *Environment) validateRequiredChainLinkLabels() error {
	if m.root.Labels() == nil {
		return fmt.Errorf("namespace labels are nil, but it should contain at least '%s' labels. Please add them to your environment config under 'Labels' key", strings.Join(requiredChainLinkNsLabels, ", "))
	}

	var missingNsLabels []string
	for _, l := range requiredChainLinkNsLabels {
		if _, ok := (*m.root.Labels())[l]; !ok {
			missingNsLabels = append(missingNsLabels, l)
		}
	}

	children := m.root.Node().Children()
	// map[workflow name][missing labels]
	missingWorkloadLabels := make(map[string][]string)
	// map[workflow name][missing labels]
	missingPodLabels := make(map[string][]string)

	if children == nil {
		return nil
	}

	var podHasLabel = func(labelName string, podLabels map[string]string) bool {
		if len(podLabels) == 0 {
			return false
		}
		for label, _ := range podLabels {
			if label == labelName {
				return true
			}
		}

		return false
	}

	for _, child := range *children {
		// most of our workloads are Helm charts
		if h, ok := child.(cdk8s.Helm); ok {
			for _, ao := range *h.ApiObjects() {
				kind := *ao.Kind()
				chartName := *ao.Name()
				// map[label]value
				var podLabels map[string]string
				nodeMightHavePods := mightHavePods(kind)
				shouldHavePodLabels := false
				if nodeMightHavePods && hasPods(kind, *ao.Chart().ToJson()) {
					shouldHavePodLabels = true
					podLabels = getJsonPodLabels(kind, *ao.Chart().ToJson())
				}
				for _, requiredLabel := range requiredChainLinkWorkloadAndPodLabels {
					maybeLabel := ao.Metadata().GetLabel(&requiredLabel)
					if maybeLabel == nil {
						missingWorkloadLabels[chartName] = append(missingWorkloadLabels[chartName], requiredLabel)
					}
					if shouldHavePodLabels {
						labelFound := podHasLabel(requiredLabel, podLabels)
						if !labelFound {
							missingPodLabels[chartName] = append(missingPodLabels[chartName], requiredLabel)
						}
					}
				}
			}
		}
		// but legacy runners have no Helm charts, but are programmatically defined as KubeJobs
		if j, ok := child.(k8s.KubeJob); ok {
			// we have already checked the Namespace
			if j.Kind() == nil || *j.Kind() == "Namespace" {
				continue
			}

			kind := *j.Kind()
			name := *j.Name()

			var podLabels map[string]string
			nodeMightHavePods := mightHavePods(kind)
			shouldHavePodLabels := false
			if nodeMightHavePods && hasPods(kind, []interface{}{j.ToJson()}) {
				shouldHavePodLabels = true
				podLabels = getJsonPodLabels(kind, []interface{}{j.ToJson()})
			}

			for _, requiredLabel := range requiredChainLinkWorkloadAndPodLabels {
				maybeLabel := j.Metadata().GetLabel(&requiredLabel)
				if maybeLabel == nil {
					missingWorkloadLabels[name] = append(missingWorkloadLabels[name], requiredLabel)
				}
				if shouldHavePodLabels {
					labelFound := podHasLabel(requiredLabel, podLabels)
					if !labelFound {
						missingPodLabels[name] = append(missingPodLabels[name], requiredLabel)
					}
				}
			}
		}
	}

	if len(missingWorkloadLabels) > 0 {
		sb := strings.Builder{}
		sb.WriteString("missing required labels for workloads:\n")
		for chart, missingLabels := range missingWorkloadLabels {
			for _, label := range missingLabels {
				sb.WriteString(fmt.Sprintf("\t'%s': '%s'\n", chart, label))
			}
		}
		sb.WriteString("Please add them to your environment configuration under 'WorkloadLabels' key. And check whether every chart has 'chain.link/component' label defined.")
		return errors.New(sb.String())
	}

	if len(missingPodLabels) > 0 {
		sb := strings.Builder{}
		sb.WriteString("missing required labels for pods:\n")
		for chart, missingLabels := range missingPodLabels {
			for _, label := range missingLabels {
				sb.WriteString(fmt.Sprintf("\t'%s': '%s'\n", chart, label))
			}
		}
		sb.WriteString("Please add them to your environment configuration under 'WorkloadLabels' key. And check whether every pod in the chart has 'chain.link/component' label defined.")
		return errors.New(sb.String())
	}

	return nil
}

func (m *Environment) initApp() error {
	var err error
	m.App = cdk8s.NewApp(&cdk8s.AppProps{
		YamlOutputType: cdk8s.YamlOutputType_FILE_PER_APP,
	})
	m.Cfg.Labels = append(m.Cfg.Labels, "app.kubernetes.io/managed-by=cdk8s")
	owner := os.Getenv(config.EnvVarUser)
	if owner == "" {
		return fmt.Errorf("missing owner environment variable, please set %s to your name or if you are seeing this in CI please set it to ${{ github.actor }}", config.EnvVarUser)
	}
	m.Cfg.Labels = append(m.Cfg.Labels, fmt.Sprintf("owner=%s", owner))

	if os.Getenv(config.EnvVarCLCommitSha) != "" {
		m.Cfg.Labels = append(m.Cfg.Labels, fmt.Sprintf("commit=%s", os.Getenv(config.EnvVarCLCommitSha)))
	}
	testTrigger := os.Getenv(config.EnvVarTestTrigger)
	if testTrigger == "" {
		testTrigger = "manual"
	}
	m.Cfg.Labels = append(m.Cfg.Labels, fmt.Sprintf("triggered-by=%s", testTrigger))

	if tolerationRole := os.Getenv(config.EnvVarToleration); tolerationRole != "" {
		m.Cfg.Tolerations = []map[string]string{{
			"key":      "node-role",
			"operator": "Equal",
			"value":    tolerationRole,
			"effect":   "NoSchedule",
		}}
	}

	if selectorRole := os.Getenv(config.EnvVarNodeSelector); selectorRole != "" {
		m.Cfg.NodeSelector = map[string]string{
			"node-role": selectorRole,
		}
	}

	m.Cfg.Labels = append(m.Cfg.Labels, fmt.Sprintf("%s=%s", pkg.TTLLabelKey, *a.ShortDur(m.Cfg.TTL)))
	nsLabels, err := a.ConvertLabels(m.Cfg.Labels)
	if err != nil {
		return err
	}

	m.root = cdk8s.NewChart(m.App, ptr.Ptr(fmt.Sprintf("root-chart-%s", m.Cfg.Namespace)), &cdk8s.ChartProps{
		Labels:    nsLabels,
		Namespace: ptr.Ptr(m.Cfg.Namespace),
	})

	k8s.NewKubeNamespace(m.root, ptr.Ptr("namespace"), &k8s.KubeNamespaceProps{
		Metadata: &k8s.ObjectMeta{
			Name:        ptr.Ptr(m.Cfg.Namespace),
			Labels:      nsLabels,
			Annotations: &defaultNamespaceAnnotations,
		},
	})
	m.CurrentManifest = *m.App.SynthYaml()
	// loop retry applying the initial manifest with the namespace and other basics
	ctx, cancel := context.WithTimeout(testcontext.Get(m.Cfg.Test), m.Cfg.ReadyCheckData.Timeout)
	defer cancel()
	startTime := time.Now()
	deadline, _ := ctx.Deadline()
	for {
		err = m.Client.Apply(ctx, m.CurrentManifest, m.Cfg.Namespace, true)
		if err == nil || ctx.Err() != nil {
			break
		}
		elapsed := time.Since(startTime)
		remaining := time.Until(deadline)
		log.Debug().Err(err).Msgf("Failed to apply initial manifest, will continue to retry. Time elapsed: %s, Time until timeout %s\n", elapsed, remaining)
		time.Sleep(5 * time.Second)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("failed to apply manifest within %s", m.Cfg.ReadyCheckData.Timeout)
	}
	if m.Cfg.PodLabels == nil {
		m.Cfg.PodLabels = map[string]string{}
	}
	m.Cfg.PodLabels[pkg.NamespaceLabelKey] = m.Cfg.Namespace
	return err
}

// AddChart adds a chart to the deployment
func (m *Environment) AddChart(f func(root cdk8s.Chart) ConnectedChart) *Environment {
	if m.err != nil {
		return m
	}
	config.JSIIGlobalMu.Lock()
	defer config.JSIIGlobalMu.Unlock()
	m.Charts = append(m.Charts, f(m.root))
	return m
}

func (m *Environment) removeChart(name string) error {
	chartIndex, _, err := m.findChart(name)
	if err != nil {
		return err
	}
	m.Charts = append(m.Charts[:chartIndex], m.Charts[chartIndex+1:]...)
	m.root.Node().TryRemoveChild(ptr.Ptr(name))
	return nil
}

// findChart finds a chart by name, returning the index of it in the Charts slice, and the chart itself
func (m *Environment) findChart(name string) (index int, chart ConnectedChart, err error) {
	for i, c := range m.Charts {
		if c.GetName() == name {
			return i, c, nil
		}
	}
	return -1, nil, fmt.Errorf("chart %s not found", name)
}

// ReplaceHelm entirely replaces an existing helm chart with a new one
// Note: you need to call Run() after this to apply the changes. If you're modifying ConfigMap values, you'll probably
// need to use RollOutStatefulSets to apply the changes to the pods. https://stackoverflow.com/questions/57356521/rollingupdate-for-stateful-set-doesnt-restart-pods-and-changes-from-updated-con
func (m *Environment) ReplaceHelm(name string, chart ConnectedChart) (*Environment, error) {
	if m.err != nil {
		return nil, m.err
	}
	config.JSIIGlobalMu.Lock()
	defer config.JSIIGlobalMu.Unlock()
	if err := m.removeChart(name); err != nil {
		return nil, err
	}
	if m.Cfg.JobImage != "" || !chart.IsDeploymentNeeded() {
		return m, fmt.Errorf("cannot modify helm chart '%s' that does not need deployment, it may be in a remote runner or detached mode", name)
	}
	log.Trace().
		Str("Chart", chart.GetName()).
		Str("Path", chart.GetPath()).
		Interface("Props", chart.GetProps()).
		Interface("Values", chart.GetValues()).
		Msg("Chart deployment values")
	h := cdk8s.NewHelm(m.root, ptr.Ptr(chart.GetName()), &cdk8s.HelmProps{
		Chart: ptr.Ptr(chart.GetPath()),
		HelmFlags: &[]*string{
			ptr.Ptr("--namespace"),
			ptr.Ptr(m.Cfg.Namespace),
		},
		ReleaseName: ptr.Ptr(chart.GetName()),
		Values:      chart.GetValues(),
	})

	workloadLabels, err := getComponentLabels(m.Cfg.WorkloadLabels, chart.GetLabels())
	if err != nil {
		m.err = err
		return nil, err
	}

	podLabels, err := getComponentLabels(m.Cfg.PodLabels, chart.GetLabels())
	if err != nil {
		m.err = err
		return nil, err
	}

	addRequiredChainLinkLabelsToWorkloads(h, workloadLabels)
	addDefaultPodAnnotationsAndLabels(h, markNotSafeToEvict(m.Cfg.PreventPodEviction, nil), podLabels)
	m.Charts = append(m.Charts, chart)
	return m, nil
}

func addDefaultPodAnnotationsAndLabels(h cdk8s.Helm, annotations, labels map[string]string) {
	annotationsCopy := map[string]string{}
	for k, v := range annotations {
		annotationsCopy[k] = v
	}
	for _, ao := range *h.ApiObjects() {
		if ao.Kind() == nil {
			continue
		}
		kind := *ao.Kind()
		if mightHavePods(kind) {
			// we aren't guaranteed to have annotations in existence so we have to dig down to see if they exist
			// and add any to our current list we want to add
			aj := *ao.Chart().ToJson()
			// loop over the json array until we get the expected kind and look for existing annotations
			for _, dep := range aj {
				l := fmt.Sprint(dep)
				if !strings.Contains(l, fmt.Sprintf("kind:%s", kind)) {
					continue
				}
				depM := dep.(map[string]interface{})
				spec, ok := depM["spec"].(map[string]interface{})
				if !ok {
					continue
				}
				template, ok := spec["template"].(map[string]interface{})
				if !ok {
					continue
				}
				metadata, ok := template["metadata"].(map[string]interface{})
				if !ok {
					continue
				}
				annot, ok := metadata["annotations"].(map[string]interface{})
				if !ok {
					continue
				}
				for k, v := range annot {
					annotationsCopy[k] = v.(string)
				}
			}
			annotationPath := "/spec/template/metadata/annotations"
			if strings.EqualFold("cronjob", kind) {
				annotationPath = "/spec/jobTemplate/spec/template/metadata/annotations"
			}
			ao.AddJsonPatch(cdk8s.JsonPatch_Add(ptr.Ptr(annotationPath), annotationsCopy))

			// loop over the labels and apply them to both the labels and selectors
			// these should in theory always have at least one label/selector combo in existence so we don't
			// have to do the existence check like we do for the annotations
			for k, v := range labels {
				// Escape the keys according to JSON Pointer syntax in RFC 6901
				escapedKey := strings.ReplaceAll(strings.ReplaceAll(k, "~", "~0"), "/", "~1")
				ao.AddJsonPatch(cdk8s.JsonPatch_Add(ptr.Ptr(fmt.Sprintf("/spec/template/metadata/labels/%s", escapedKey)), v))
				// CronJob doesn't have a selector, so we don't need to add it
				if !strings.EqualFold("cronjob", kind) {
					ao.AddJsonPatch(cdk8s.JsonPatch_Add(ptr.Ptr(fmt.Sprintf("/spec/selector/matchLabels/%s", escapedKey)), v))
				}
			}
		}
	}
}

func addRequiredChainLinkLabelsToWorkloads(h cdk8s.Helm, labels map[string]string) {
	for _, ao := range *h.ApiObjects() {
		for k, v := range labels {
			ao.Metadata().AddLabel(&k, &v)
		}
	}
}

// UpdateHelm update a helm chart with new values. The pod will launch with an `updated=true` label if it's a Chainlink node.
// Note: If you're modifying ConfigMap values, you'll probably need to use RollOutStatefulSets to apply the changes to the pods.
// https://stackoverflow.com/questions/57356521/rollingupdate-for-stateful-set-doesnt-restart-pods-and-changes-from-updated-con
func (m *Environment) UpdateHelm(name string, values map[string]any) (*Environment, error) {
	if m.err != nil {
		return nil, m.err
	}
	_, chart, err := m.findChart(name)
	if err != nil {
		return nil, err
	}
	if _, labelsExist := values["labels"]; !labelsExist {
		values["labels"] = make(map[string]*string)
	}
	values["labels"].(map[string]*string)["updated"] = ptr.Ptr("true")
	if err = mergo.Merge(chart.GetValues(), values, mergo.WithOverride); err != nil {
		return nil, err
	}
	return m.ReplaceHelm(name, chart)
}

// Charts adds multiple helm charts to the testing environment
func (m *Environment) AddHelmCharts(charts []ConnectedChart) *Environment {
	if m.err != nil {
		return m
	}
	for _, c := range charts {
		m.AddHelm(c)
	}
	return m
}

// AddHelm adds a helm chart to the testing environment
func (m *Environment) AddHelm(chart ConnectedChart) *Environment {
	if m.err != nil {
		return m
	}
	if m.Cfg.JobImage != "" || !chart.IsDeploymentNeeded() {
		return m
	}
	config.JSIIGlobalMu.Lock()
	defer config.JSIIGlobalMu.Unlock()

	values := &map[string]any{
		"tolerations":  m.Cfg.Tolerations,
		"nodeSelector": m.Cfg.NodeSelector,
	}
	config.MustMerge(values, chart.GetValues())
	log.Trace().
		Str("Chart", chart.GetName()).
		Str("Path", chart.GetPath()).
		Interface("Props", chart.GetProps()).
		Interface("Values", values).
		Msg("Chart deployment values")
	helmFlags := []*string{
		ptr.Ptr("--namespace"),
		ptr.Ptr(m.Cfg.Namespace),
		ptr.Ptr("--skip-tests"),
	}
	if chart.GetVersion() != "" {
		helmFlags = append(helmFlags, ptr.Ptr("--version"), ptr.Ptr(chart.GetVersion()))
	}
	chartPath, err := m.PullOCIChart(chart)
	if err != nil {
		m.err = err
		return m
	}
	h := cdk8s.NewHelm(m.root, ptr.Ptr(chart.GetName()), &cdk8s.HelmProps{
		Chart:       ptr.Ptr(chartPath),
		HelmFlags:   &helmFlags,
		ReleaseName: ptr.Ptr(chart.GetName()),
		Values:      values,
	})

	workloadLabels, err := getComponentLabels(m.Cfg.WorkloadLabels, chart.GetLabels())
	if err != nil {
		m.err = err
		return m
	}

	podLabels, err := getComponentLabels(m.Cfg.PodLabels, chart.GetLabels())
	if err != nil {
		m.err = err
		return m
	}
	addRequiredChainLinkLabelsToWorkloads(h, workloadLabels)
	addDefaultPodAnnotationsAndLabels(h, markNotSafeToEvict(m.Cfg.PreventPodEviction, nil), podLabels)
	m.Charts = append(m.Charts, chart)
	return m
}

func getComponentLabels(podLabels, chartLabels map[string]string) (map[string]string, error) {
	componentLabels := make(map[string]string)
	err := mergo.Merge(&componentLabels, podLabels, mergo.WithOverride)
	if err != nil {
		return nil, err
	}
	err = mergo.Merge(&componentLabels, chartLabels, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	return componentLabels, nil
}

// PullOCIChart handles working with OCI format repositories
// https://helm.sh/docs/topics/registries/
// API is not compatible between helm repos and OCI repos, so we download and untar the chart
func (m *Environment) PullOCIChart(chart ConnectedChart) (string, error) {
	if !strings.HasPrefix(chart.GetPath(), "oci") {
		return chart.GetPath(), nil
	}
	cp := strings.Split(chart.GetPath(), "/")
	if len(cp) != 5 {
		return "", fmt.Errorf(ErrInvalidOCI, chart.GetPath())
	}
	sp := strings.Split(chart.GetPath(), ":")

	var cmd string
	var chartName string
	chartName = cp[len(cp)-1]
	chartDir := uuid.NewString()
	switch len(sp) {
	case 2:
		cmd = fmt.Sprintf("helm pull %s --untar --untardir %s", chart.GetPath(), chartDir)
	case 3:
		chartName = strings.Split(chartName, ":")[0]
		cmd = fmt.Sprintf("helm pull %s --version %s --untar --untardir %s", fmt.Sprintf("%s:%s", sp[0], sp[1]), sp[2], chartDir)
	default:
		return "", fmt.Errorf(ErrInvalidOCI, chart.GetPath())
	}
	log.Info().Str("CMD", cmd).Msg("Running helm cmd")
	if err := client.ExecCmd(cmd); err != nil {
		return "", fmt.Errorf(ErrOCIPull, chart.GetPath())
	}
	localChartPath := fmt.Sprintf("%s/%s/", chartDir, chartName)
	log.Info().Str("Path", localChartPath).Msg("Local chart path")
	return localChartPath, nil
}

// PrintExportData prints export data
func (m *Environment) PrintExportData() error {
	m.URLs = make(map[string][]string)
	for _, c := range m.Charts {
		err := c.ExportData(m)
		if err != nil {
			return err
		}
	}
	log.Debug().Interface("URLs", m.URLs).Msg("Connection URLs")
	return nil
}

// DumpLogs dumps all logs into a file
func (m *Environment) DumpLogs(path string) error {
	arts, err := NewArtifacts(m.Client, m.Cfg.Namespace)
	if err != nil {
		return err
	}
	if path == "" {
		path = fmt.Sprintf("logs/%s-%d", m.Cfg.Namespace, time.Now().Unix())
	}
	return arts.DumpTestResult(path, "chainlink")
}

// ResourcesSummary returns resources summary for selected pods as a map, used in reports
func (m *Environment) ResourcesSummary(selector string) (map[string]map[string]string, error) {
	pl, err := m.Client.ListPods(m.Cfg.Namespace, selector)
	if err != nil {
		return nil, err
	}
	if len(pl.Items) == 0 {
		return nil, fmt.Errorf("no pods found for selector: %s", selector)
	}
	resources := make(map[string]map[string]string)
	for _, p := range pl.Items {
		for _, c := range p.Spec.Containers {
			if resources[c.Name] == nil {
				resources[c.Name] = make(map[string]string)
			}
			cpuRes := c.Resources.Requests["cpu"]
			resources[c.Name]["cpu"] = cpuRes.String()
			memRes := c.Resources.Requests["memory"]
			resources[c.Name]["memory"] = memRes.String()
		}
	}
	return resources, nil
}

// ClearCharts recreates cdk8s app
func (m *Environment) ClearCharts() error {
	m.Charts = make([]ConnectedChart, 0)
	if err := m.initApp(); err != nil {
		log.Error().Err(err).Msg("failed to apply the initial manifest to create the namespace")
		return err
	}
	return nil
}

func (m *Environment) Manifest() string {
	return m.CurrentManifest
}

// Update current manifest based on the cdk8s app state
func (m *Environment) UpdateManifest() {
	config.JSIIGlobalMu.Lock()
	m.CurrentManifest = *m.App.SynthYaml()
	config.JSIIGlobalMu.Unlock()
}

// RunCustomReadyConditions Runs the environment with custom ready conditions for a supplied pod count
func (m *Environment) RunCustomReadyConditions(customCheck *client.ReadyCheckData, podCount int) error {
	if m.err != nil {
		return m.err
	}
	if m.Cfg.jobDeployed {
		return nil
	}
	if m.Cfg.JobImage != "" {
		if m.Cfg.Test == nil {
			return fmt.Errorf("Test must be configured in the environment when using the remote runner")
		}
		remoteRunnerLabels := map[string]*string{pkg.NamespaceLabelKey: ptr.Ptr(m.Cfg.Namespace)}
		for l, v := range m.Cfg.WorkloadLabels {
			remoteRunnerLabels[l] = ptr.Ptr(v)
		}
		// if no runner name is specified use constant
		if m.Cfg.RunnerName == "" {
			m.Cfg.RunnerName = REMOTE_RUNNER_NAME
		}
		m.AddChart(NewRunner(&Props{
			BaseName:           m.Cfg.RunnerName,
			ReportPath:         m.Cfg.ReportPath,
			TargetNamespace:    m.Cfg.Namespace,
			Labels:             &remoteRunnerLabels,
			Image:              m.Cfg.JobImage,
			TestName:           m.Cfg.Test.Name(),
			SkipManifestUpdate: m.Cfg.SkipManifestUpdate,
			PreventPodEviction: m.Cfg.PreventPodEviction,
		}))
		// add a pod to access reports generated by remote-runner, even after remote-runner's job execution completion
		if m.Cfg.ReportPath != "" {
			m.AddChart(DataFromRunner(&Props{
				BaseName:           m.Cfg.RunnerName,
				ReportPath:         m.Cfg.ReportPath,
				TargetNamespace:    m.Cfg.Namespace,
				Labels:             &remoteRunnerLabels,
				Image:              m.Cfg.JobImage,
				TestName:           m.Cfg.Test.Name(),
				SkipManifestUpdate: m.Cfg.SkipManifestUpdate,
				PreventPodEviction: m.Cfg.PreventPodEviction,
			}))
		}
	}
	m.UpdateManifest()
	m.ChainlinkNodeDetails = []*ChainlinkNodeDetail{} // Resets potentially old details if re-deploying
	if m.Cfg.DryRun {
		log.Info().Msg("Dry-run mode, manifest synthesized and saved as tmp-manifest.yaml")
		return nil
	}
	manifestUpdate := os.Getenv(config.EnvVarSkipManifestUpdate)
	if manifestUpdate != "" {
		mu, err := strconv.ParseBool(manifestUpdate)
		if err != nil {
			return fmt.Errorf("manifest update should be bool: true, false")
		}
		m.Cfg.SkipManifestUpdate = mu
	}
	log.Debug().Bool("ManifestUpdate", m.Cfg.SkipManifestUpdate).Msg("Update mode")

	if !m.Cfg.SkipRequiredChainLinkLabelsValidation {
		// make sure all required chain.link labels are present in the final manifest
		if err := m.validateRequiredChainLinkLabels(); err != nil {
			return err
		}
	}

	if !m.Cfg.SkipManifestUpdate || m.Cfg.JobImage != "" {
		if err := m.DeployCustomReadyConditions(customCheck, podCount); err != nil {
			log.Error().Err(err).Msg("Error deploying environment")
			_ = m.Shutdown()
			return err
		}
	}
	if m.Cfg.JobImage != "" {
		log.Info().Msg("Waiting for remote runner to complete")
		// Do not wait for the job to complete if we are running something like a soak test in the remote runner
		if m.Cfg.detachRunner {
			return nil
		}
		if err := m.Client.WaitForJob(m.Cfg.Namespace, m.Cfg.RunnerName, func(message string) {
			if m.Cfg.JobLogFunction != nil {
				m.Cfg.JobLogFunction(m, message)
			} else {
				DefaultJobLogFunction(m, message)
			}
		}); err != nil {
			return err
		}
		if m.Cfg.fundReturnFailed {
			return fmt.Errorf("failed to return funds in remote runner")
		}
		m.Cfg.jobDeployed = true
	} else {
		if err := m.Fwd.Connect(m.Cfg.Namespace, "", m.Cfg.InsideK8s); err != nil {
			return err
		}
		log.Debug().Interface("Ports", m.Fwd.Info).Msg("Forwarded ports")
		m.Fwd.PrintLocalPorts()
		if err := m.PrintExportData(); err != nil {
			return err
		}
		arts, err := NewArtifacts(m.Client, m.Cfg.Namespace)
		if err != nil {
			log.Error().Err(err).Msg("failed to create artifacts client")
			return err
		}
		m.Artifacts = arts
		if len(m.URLs["goc"]) != 0 {
			m.httpClient = resty.New().SetBaseURL(m.URLs["goc"][0])
		}
		if m.Cfg.KeepConnection {
			log.Info().Msg("Keeping forwarder connections, press Ctrl+C to interrupt")
			if m.Cfg.RemoveOnInterrupt {
				log.Warn().Msg("Environment will be removed on interrupt")
			}
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
			<-ch
			log.Warn().Msg("Interrupted")
			if m.Cfg.RemoveOnInterrupt {
				return m.Client.RemoveNamespace(m.Cfg.Namespace)
			}
		}
	}
	return nil
}

// RunUpdated runs the environment and checks for pods with `updated=true` label
func (m *Environment) RunUpdated(podCount int) error {
	if m.err != nil {
		return m.err
	}
	conds := &client.ReadyCheckData{
		ReadinessProbeCheckSelector: "updated=true",
		Timeout:                     10 * time.Minute,
	}
	return m.RunCustomReadyConditions(conds, podCount)
}

// Run deploys or connects to already created environment
func (m *Environment) Run() error {
	if m.err != nil {
		return m.err
	}
	return m.RunCustomReadyConditions(nil, 0)
}

func (m *Environment) enumerateApps() error {
	apps, err := m.Client.UniqueLabels(m.Cfg.Namespace, client.AppLabel)
	if err != nil {
		return err
	}
	for _, app := range apps {
		if err := m.Client.EnumerateInstances(m.Cfg.Namespace, fmt.Sprintf("app=%s", app)); err != nil {
			return err
		}
	}
	return nil
}

// DeployCustomReadyConditions deploy current manifest with added custom readiness checks
func (m *Environment) DeployCustomReadyConditions(customCheck *client.ReadyCheckData, customPodCount int) error {
	if m.err != nil {
		return m.err
	}
	log.Info().Str("Namespace", m.Cfg.Namespace).Msg("Deploying namespace")

	if m.Cfg.DryRun {
		return m.Client.DryRun(m.CurrentManifest)
	}
	ctx, cancel := context.WithTimeout(testcontext.Get(m.Cfg.Test), m.Cfg.ReadyCheckData.Timeout)
	defer cancel()
	err := m.Client.Apply(ctx, m.CurrentManifest, m.Cfg.Namespace, true)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("timeout waiting for environment to be ready")
	}
	if err != nil {
		return err
	}
	if int64(m.Cfg.UpdateWaitInterval) != 0 {
		time.Sleep(m.Cfg.UpdateWaitInterval)
	}

	expectedPodCount := m.findPodCountInDeploymentManifest()

	if err := m.Client.WaitPodsReady(m.Cfg.Namespace, m.Cfg.ReadyCheckData, expectedPodCount); err != nil {
		return err
	}
	if customCheck != nil {
		if err := m.Client.WaitPodsReady(m.Cfg.Namespace, customCheck, customPodCount); err != nil {
			return err
		}
	}
	return m.enumerateApps()
}

// Deploy deploys current manifest and check logs for readiness
func (m *Environment) Deploy() error {
	return m.DeployCustomReadyConditions(nil, 0)
}

// RolloutStatefulSets applies "rollout statefulset" to all existing statefulsets in our namespace
func (m *Environment) RolloutStatefulSets() error {
	if m.err != nil {
		return m.err
	}
	ctx, cancel := context.WithTimeout(testcontext.Get(m.Cfg.Test), m.Cfg.ReadyCheckData.Timeout)
	defer cancel()
	err := m.Client.RolloutStatefulSets(ctx, m.Cfg.Namespace)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("timeout waiting for rollout statefulset to complete")
	}
	return err
}

// CopyFromPod lists pods with given selector, it copies files from local to destPath at pods filtered by given selector
func (m *Environment) CopyFromPod(selector, containerName, srcPath, destPath string) error {
	pl, err := m.Client.ListPods(m.Cfg.Namespace, selector)
	if err != nil {
		return err
	}
	if len(pl.Items) == 0 {
		return fmt.Errorf("no pods found for selector: %s", selector)
	}
	for _, p := range pl.Items {
		err := m.Client.CopyFromPod(context.Background(), m.Cfg.Namespace, p.Name, containerName, srcPath, destPath)
		if err != nil {
			return fmt.Errorf("%w error copying from %s:%s to destination path %s", err, p.Name, srcPath, destPath)
		}
	}
	return nil
}

// CopyToPod lists pods with given selector, it copies files from srcPath at pods filtered by given selector to
// local destPath
func (m *Environment) CopyToPod(selector, containerName, srcPath, destPath string) error {
	pl, err := m.Client.ListPods(m.Cfg.Namespace, selector)
	if err != nil {
		return err
	}
	if len(pl.Items) == 0 {
		return fmt.Errorf("no pods found for selector: %s", selector)
	}
	for _, p := range pl.Items {
		destPath = fmt.Sprintf("%s/%s:/%s", m.Cfg.Namespace, p.Name, destPath)
		_, _, _, err := m.Client.CopyToPod(m.Cfg.Namespace, srcPath, destPath, containerName)
		if err != nil {
			return fmt.Errorf("%w error copying from %s to destination path %s", err, srcPath, destPath)
		}
	}
	return nil
}

// RolloutRestartBySelector applies "rollout restart" to the selected resources
func (m *Environment) RolloutRestartBySelector(resource string, selector string) error {
	if m.err != nil {
		return m.err
	}
	ctx, cancel := context.WithTimeout(testcontext.Get(m.Cfg.Test), m.Cfg.ReadyCheckData.Timeout)
	defer cancel()
	err := m.Client.RolloutRestartBySelector(ctx, m.Cfg.Namespace, resource, selector)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("timeout waiting for rollout restart to complete")
	}
	return err
}

// findPodsInDeploymentManifest finds all the pods we will be deploying
func (m *Environment) findPodCountInDeploymentManifest() int {
	config.JSIIGlobalMu.Lock()
	defer config.JSIIGlobalMu.Unlock()
	podCount := 0
	charts := m.App.Charts()
	for _, chart := range *charts {
		json := chart.ToJson()
		if json == nil {
			continue
		}
		for _, j := range *json {
			m := j.(map[string]any)
			// if the kind is a deployment then we want to see if it has replicas to count towards the app count
			if _, ok := m["kind"]; !ok {
				continue
			}
			kind := m["kind"].(string)
			if kind == "Deployment" || kind == "StatefulSet" {
				if _, ok := m["spec"]; !ok {
					continue
				}
				podCount += getReplicaCount(m["spec"].(map[string]any))
			}
		}

	}
	return podCount
}

func getReplicaCount(spec map[string]any) int {
	if spec == nil {
		return 0
	}
	if _, ok := spec["selector"]; !ok {
		return 0
	}
	s := spec["selector"].(map[string]any)
	if s == nil {
		return 0
	}
	if _, ok := s["matchLabels"]; !ok {
		return 0
	}
	m := s["matchLabels"].(map[string]any)
	if m == nil {
		return 0
	}
	if _, ok := m[client.AppLabel]; !ok {
		return 0
	}
	l := m[client.AppLabel]
	if l == nil {
		return 0
	}

	replicaCount := 0
	var replicas any
	replicas, ok := spec["replicas"]
	if ok {
		replicaCount += int(replicas.(float64))
	} else {
		replicaCount++
	}

	return replicaCount
}

type CoverageProfileParams struct {
	Force             bool     `form:"force" json:"force"`
	Service           []string `form:"service" json:"service"`
	Address           []string `form:"address" json:"address"`
	CoverFilePatterns []string `form:"coverfile" json:"coverfile"`
	SkipFilePatterns  []string `form:"skipfile" json:"skipfile"`
}

func (m *Environment) getCoverageList() (map[string]any, error) {
	var servicesMap map[string]any
	resp, err := m.httpClient.R().
		SetResult(&servicesMap).
		Get("v1/cover/list")
	if err != nil {
		return nil, err
	}
	if resp.Status() != "200 OK" {
		return nil, fmt.Errorf("coverage service list request is not 200")
	}
	return servicesMap, nil
}

func (m *Environment) ClearCoverage() error {
	servicesMap, err := m.getCoverageList()
	if err != nil {
		return err
	}
	for serviceName := range servicesMap {
		r, err := m.httpClient.R().
			SetBody(CoverageProfileParams{Service: []string{serviceName}}).
			Post("v1/cover/clear")
		if err != nil {
			return err
		}
		if r.Status() != "200 OK" {
			return fmt.Errorf("coverage service list request is not 200")
		}
		log.Debug().Str("Service", serviceName).Msg("Coverage cleared")
	}
	return nil
}

func (m *Environment) SaveCoverage() error {
	if err := MkdirIfNotExists(COVERAGE_DIR); err != nil {
		return err
	}
	servicesMap, err := m.getCoverageList()
	if err != nil {
		return err
	}
	log.Debug().Interface("Services", servicesMap).Msg("Services eligible for coverage")
	for serviceName := range servicesMap {
		r, err := m.httpClient.R().
			SetBody(CoverageProfileParams{Service: []string{serviceName}}).
			Post("v1/cover/profile")
		if err != nil {
			return err
		}
		if r.Status() != "200 OK" {
			return fmt.Errorf("coverage service list request is not 200")
		}
		log.Debug().Str("Service", serviceName).Msg("Coverage received")
		if err := os.WriteFile(fmt.Sprintf("%s/%s.cov", COVERAGE_DIR, serviceName), r.Body(), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown environment, remove namespace
func (m *Environment) Shutdown() error {
	// don't shutdown if returning of funds failed
	if m.Cfg.fundReturnFailed {
		return nil
	}

	// don't shutdown if this is a test running remotely
	if m.Cfg.InsideK8s {
		return nil
	}

	keepEnvs := os.Getenv(config.EnvVarKeepEnvironments)
	if keepEnvs == "" {
		keepEnvs = "NEVER"
	}

	shouldShutdown := false
	switch strings.ToUpper(keepEnvs) {
	case "ALWAYS":
		return nil
	case "ONFAIL":
		if m.Cfg.Test != nil {
			if !m.Cfg.Test.Failed() {
				shouldShutdown = true
			}
		}
	case "NEVER":
		shouldShutdown = true
	default:
		log.Warn().Str("Invalid Keep Value", keepEnvs).
			Msg("Invalid 'keep_environments' value, see the KEEP_ENVIRONMENTS env var")
	}

	if shouldShutdown {
		return m.Client.RemoveNamespace(m.Cfg.Namespace)
	}
	return nil
}

// WillUseRemoteRunner determines if we need to start the remote runner
func (m *Environment) WillUseRemoteRunner() bool {
	val, _ := os.LookupEnv(config.EnvVarJobImage)
	return val != "" && m.Cfg != nil && m.Cfg.Test != nil && m.Cfg.Test.Name() != ""
}

func DefaultJobLogFunction(e *Environment, message string) {
	logChunks := logging.SplitStringIntoChunks(message, 50000)
	for _, chunk := range logChunks {
		e.Cfg.Test.Log(chunk)
	}
	if strings.Contains(message, FAILED_FUND_RETURN) {
		e.Cfg.fundReturnFailed = true
	}
	if strings.Contains(message, TEST_FAILED) {
		e.Cfg.Test.Fail()
	}
}

// markNotSafeToEvict adds the safe to evict annotation to the provided map if needed
func markNotSafeToEvict(preventPodEviction bool, m map[string]string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	if preventPodEviction {
		m["karpenter.sh/do-not-evict"] = "true"
		m["karpenter.sh/do-not-disrupt"] = "true"
		m["cluster-autoscaler.kubernetes.io/safe-to-evict"] = "false"
	}

	return m
}

// GetRequiredChainLinkNamespaceLabels returns the required chain.link namespace labels
// if `CHAINLINK_USER_TEAM` env var is not set it will return an error
func GetRequiredChainLinkNamespaceLabels(product, testType string) ([]string, error) {
	var nsLabels []string
	createdLabels, err := createRequiredChainLinkLabels(product, testType)
	if err != nil {
		return nsLabels, err
	}

	for k, v := range createdLabels {
		nsLabels = append(nsLabels, fmt.Sprintf("%s=%s", k, v))
	}

	return nsLabels, nil
}

// GetRequiredChainLinkWorkloadAndPodLabels returns the required chain.link workload and pod labels
// if `CHAINLINK_USER_TEAM` env var is not set it will return an error
func GetRequiredChainLinkWorkloadAndPodLabels(product, testType string) (map[string]string, error) {
	createdLabels, err := createRequiredChainLinkLabels(product, testType)
	if err != nil {
		return nil, err
	}

	return createdLabels, nil
}

func createRequiredChainLinkLabels(product, testType string) (map[string]string, error) {
	team := os.Getenv(config.EnvVarTeam)
	if team == "" {
		return nil, fmt.Errorf("missing team environment variable, please set %s to your team name or if you are seeing this in CI please either add a new input with team name or hardcode it if this jobs is only run by a single team", config.EnvVarTeam)
	}

	return map[string]string{
		"chain.link/product":     product,
		"chain.link/team":        team,
		"chain.link/cost-center": fmt.Sprintf("test-tooling-%s-test", testType),
	}, nil
}

// mightHavePods returns true if the kind of k8s resource might have pods
func mightHavePods(kind string) bool {
	switch kind {
	// only these objects contain pods in their definition
	case "Deployment", "ReplicaSet", "StatefulSet", "Job", "CronJob", "DaemonSet":
		return true
	}

	return false
}

// hasPods checks if the json representing k8s resource has any pods
func hasPods(kind string, maybeJson []interface{}) bool {
	var hasSpecContainers = func(depMap map[string]interface{}) bool {
		spec, ok := depMap["spec"].(map[string]interface{})
		if !ok {
			return false
		}
		containers, ok := spec["containers"].([]interface{})
		if !ok {
			return false
		}

		return len(containers) > 0
	}

	switch kind {
	case "CronJob":
		foundPods := false
		for _, dep := range maybeJson {
			depM := dep.(map[string]interface{})
			spec, ok := depM["spec"].(map[string]interface{})
			if !ok {
				continue
			}
			jobTemplate, ok := spec["jobTemplate"].(map[string]interface{})
			if !ok {
				continue
			}
			spec2, ok := jobTemplate["spec"].(map[string]interface{})
			if !ok {
				continue
			}
			template, ok := spec2["template"].(map[string]interface{})
			if !ok {
				continue
			}
			if hasSpecContainers(template) {
				foundPods = true
				break
			}
		}
		return foundPods
	default:
		foundPods := false
		for _, dep := range maybeJson {
			depM := dep.(map[string]interface{})
			spec, ok := depM["spec"].(map[string]interface{})
			if !ok {
				continue
			}
			template, ok := spec["template"].(map[string]interface{})
			if !ok {
				continue
			}
			if hasSpecContainers(template) {
				foundPods = true
				break
			}
		}
		return foundPods
	}

	return false
}

// getJsonPodLabels returns the labels for the pods in the json, which represents k8s resource
// it returns pod labels of the first resource with these labels
func getJsonPodLabels(kind string, maybeJson []interface{}) map[string]string {
	templateLabels := make(map[string]string)
	for _, dep := range maybeJson {
		l := fmt.Sprint(dep)
		if !strings.Contains(l, fmt.Sprintf("kind:%s", kind)) {
			continue
		}
		depM := dep.(map[string]interface{})
		spec, ok := depM["spec"].(map[string]interface{})
		if !ok {
			continue
		}

		// CronJob has a different structure for the labels
		var specRoot map[string]interface{}
		if strings.EqualFold(kind, "CronJob") {
			jobTemplate, ok := spec["jobTemplate"].(map[string]interface{})
			if !ok {
				continue
			}
			spec2, ok := jobTemplate["spec"].(map[string]interface{})
			if !ok {
				continue
			}
			specRoot = spec2
		} else {
			specRoot = spec
		}

		template, ok := specRoot["template"].(map[string]interface{})
		if !ok {
			continue
		}
		metadata, ok := template["metadata"].(map[string]interface{})
		if !ok {
			continue
		}
		labels, ok := metadata["labels"].(map[string]interface{})
		if !ok {
			continue
		}

		for k, v := range labels {
			templateLabels[k] = v.(string)
		}
		break
	}

	return templateLabels
}
