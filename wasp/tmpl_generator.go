package wasp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// TestCodeGenerator defines the contract for all code generators
// it is not mandatory but it expresses an idea how we structure code generation
// A single code generator may:
// - Generate a new Go module (go.mod) or just new files for existing module
// - Use static params or read dynamic data (k8s namespace, for example)
// - Generate multiple test cases
// - Generate only a single table test (need more - let's have another generator)
// - Generate boilerplate (one or multiple files)
type TestCodeGenerator interface {
	// Read reads dynamic data from any source, ex. K8s namespace, Docker containers, etc
	Read() error
	// GenerateGoMod generated a new "go.mod" and adds all required dependencies
	GenerateGoMod() (string, error)
	// GenerateTestCases generates one or more test cases for a table test
	GenerateTestCases() ([]TestCaseParams, error)
	// GenerateTableTest generates a table test
	GenerateTableTest() (string, error)
	// Write uses all the data from the methods above and writes files
	Write() error
}

var _ TestCodeGenerator = (*LoadTestCodegen)(nil)

const (
	DefaultTestSuiteName = "TestGeneratedLoadChaos"
	DefaultUniqLabel     = "app.kubernetes.io/instance"
	DefaultTestModule    = "main"
)

/* Templates */

const (
	GoModTemplate = `module {{.ModuleName}}

go 1.25

replace github.com/smartcontractkit/chainlink-testing-framework/wasp => ../../../wasp/
`

	// TableTestTmpl is a load/chaos table test template
	TableTestTmpl = `package {{.Package}}

import (
		"context"
		"os"
		"testing"
		"time"

		f "github.com/smartcontractkit/chainlink-testing-framework/framework"
		"github.com/go-resty/resty/v2"
		"github.com/smartcontractkit/chainlink-testing-framework/wasp"
		havoc "github.com/smartcontractkit/chainlink-testing-framework/wasp/havoc"
		"github.com/stretchr/testify/require"
)

{{.GunCode}}

type GenK8sChaos struct {
		WaitBeforeStart             string
		Namespace                   string
		DashboardUUIDs              []string
		ExperimentDuration          string
		ExperimentInjectionDuration string
		RemoveK8sChaos              bool
}

type GenChaosCfg struct {
		Chaos *GenK8sChaos
}

func {{.TableTestName}}(t *testing.T) {
		{{.WorkloadCode}}

		cfg := &GenChaosCfg{
			Chaos: &GenK8sChaos{
				WaitBeforeStart:             "30s",
				Namespace:                   "{{.Namespace}}",
				// TODO: find and insert the dashboard UUIDs on which you need annotations!
				DashboardUUIDs:              []string{},
				// TODO: choose each experiment duration
				ExperimentDuration:          "10s",
				ExperimentInjectionDuration: "5s",
				// TODO: whether to remove the chaos instance CRDs or not, useful for debugging in corner cases where logs can't help
				RemoveK8sChaos:              true,
			},
		}

		c, err := havoc.NewChaosMeshClient()
		require.NoError(t, err)
		cr := havoc.NewNamespaceRunner(f.L, c, cfg.Chaos.RemoveK8sChaos)
		gc := f.NewGrafanaClient(os.Getenv("GRAFANA_URL"), os.Getenv("GRAFANA_TOKEN"))

		testCases := []struct {
			name     string
			run      func(t *testing.T)
			validate func(t *testing.T)
		}{ {{range .TestCases}}
			{
				name: "{{.Name}}",
				run: func(t *testing.T) {
					{{.RunFunc}}
				},
				validate: func(t *testing.T) {
					// TODO: add post-experiment validation here if needed
				},
			},{{end}}
		}

		startsIn := f.MustParseDuration(cfg.Chaos.WaitBeforeStart)
		f.L.Info().Msgf("Starting chaos tests in %s", startsIn)
		time.Sleep(startsIn)

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				n := time.Now()
				testCase.run(t)
				time.Sleep(f.MustParseDuration(cfg.Chaos.ExperimentDuration))
				if os.Getenv("GRAFANA_URL") != "" {
					_, _, err := gc.Annotate(f.A(cfg.Chaos.Namespace, testCase.name, cfg.Chaos.DashboardUUIDs, havoc.Ptr(n), havoc.Ptr(time.Now())))
					require.NoError(t, err)
				}
				testCase.validate(t)
			})
		}
}`

	// GunHTTPCode is the most simple HTTP Gun example for WASP load generator
	GunHTTPCode = `type ExampleGun struct {
		target string
		client *resty.Client
		Data   []string
}

func NewExampleHTTPGun(target string) *ExampleGun {
		return &ExampleGun{
			client: resty.New(),
			target: target,
			Data:   make([]string, 0),
		}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *ExampleGun) Call(l *wasp.Generator) *wasp.Response {
		var result map[string]any
		r, err := m.client.R().
			SetResult(&result).
			Get(m.target)
		if err != nil {
			return &wasp.Response{Data: result, Error: err.Error()}
		}
		if r.Status() != "200 OK" {
			return &wasp.Response{Data: result, Error: "not 200", Failed: true}
		}
		return &wasp.Response{Data: result}
}
`

	// WorkloadHTTPCode is the most simple test part that uses WASP's HTTP load generator
	WorkloadHTTPCode = `
		// TODO: add service URL for http/grpc or other web services
		// remove in case we'll be generating blockchain transactions
		serviceUnderLoadURL := ""

		labels := map[string]string{
			// TODO: set "go_test_name" and "gen_name" to be able to filter things in WASP dashboard
			"go_test_name": "generator_healthcheck",
			"gen_name":     "generator_healthcheck",
			// TODO: in case we need to compare between commits/tags these parameters can be filled in runtime
			"branch":       "test",
			"commit":       "test",
		}

		gen, err := wasp.NewGenerator(&wasp.Config{
			LoadType: wasp.RPS,
			T:        t,
			// TODO: set required schedule, see wasp.Steps and wasp.CombineAndRepeat
			Schedule:   wasp.Plain(5, 60*time.Second),
			Gun:        NewExampleHTTPGun(serviceUnderLoadURL),
			Labels:     labels,
			// TODO: by default we are sending data to CTF obs stack, use wasp.
			LokiConfig: wasp.LocalCTFObsConfig(),
		})
		require.NoError(t, err)

		gen.Run(false)
`

	// Kubernetes ChaosMesh templates for our wrapper Havoc, majority of them relies on assumption that your
	// K8s deployments follow best labelling practices
	// described here: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels

	// PodFailTmpl fails one pod by unique label
	PodFailTmpl = `_, err := cr.RunPodFail(context.Background(),
				havoc.PodFailCfg{
					Namespace:         cfg.Chaos.Namespace,
					LabelKey:          "{{.LabelKey}}",
					LabelValues:       []string{"{{.LabelValue}}"},
					InjectionDuration: f.MustParseDuration(cfg.Chaos.ExperimentInjectionDuration),
				})
require.NoError(t, err)`

	// PodDelayTmpl simulates network delay for one pud by unique label
	PodDelayTmpl = `_, err := cr.RunPodDelay(context.Background(),
    havoc.PodDelayCfg{
        Namespace:         cfg.Chaos.Namespace,
        LabelKey:          "{{.LabelKey}}",
        LabelValues:       []string{"{{.LabelValue}}"},
        Latency:           {{.LatencyMs}} * time.Millisecond,
        Jitter:            {{.JitterMs}} * time.Millisecond,
        Correlation:       "0",
        InjectionDuration: f.MustParseDuration(cfg.Chaos.ExperimentInjectionDuration),
    })
require.NoError(t, err)
`
)

/* Template params in heirarchical order, module -> file(table test) -> test */

// GoModParams params for generating go.mod file
type GoModParams struct {
	ModuleName string
}

// TableTestParams params for generating a table test
type TableTestParams struct {
	Package       string
	Namespace     string
	TableTestName string
	TestCases     []TestCaseParams
	WorkloadCode  string
	GunCode       string
}

// TestCaseParams params for generating a test case
type TestCaseParams struct {
	Name    string
	RunFunc string
}

// PodFail params for pod delay test case
type PodFailParams struct {
	LabelKey   string
	LabelValue string
}

// PodDelayParams params for pod delay test case
type PodDelayParams struct {
	LabelKey   string
	LabelValue string
	LatencyMs  int
	JitterMs   int
}

// K8sPodsInterface defines the interface for Kubernetes interactions
// we are not using mockery lib here because we don't need more than 2-3 methods
type K8sPodsInterface interface {
	GetPods(ctx context.Context, namespace string) (*corev1.PodList, error)
}

// K8s implements K8sClient using the actual Kubernetes client
type K8s struct {
	clientset *kubernetes.Clientset
}

// NewK8s creates new K8s client
func NewK8s() (*K8s, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}
	return &K8s{clientset: clientset}, nil
}

// GetPods returns K8s Pods list
func (r *K8s) GetPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	return r.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

// MockK8s implements K8sClient for testing
type MockK8s struct {
	Pods *corev1.PodList
	Err  error
}

// GetPods returns K8s Pods list
func (m *MockK8s) GetPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	return m.Pods, m.Err
}

/* Codegen logic */

// LoadTestBuilder builder for load test codegen
type LoadTestBuilder struct {
	namespace       string
	testSuiteName   string
	pods            []corev1.Pod
	uniqPodLabelKey string
	latencyMs       int
	jitterMs        int
	includeWorkload bool
	outputDir       string
	moduleName      string
	k8sClient       K8sPodsInterface
}

// LoadTestCodegen is a load test code generator that creates workload and chaos experiments
type LoadTestCodegen struct {
	cfg *LoadTestBuilder
}

// NewLoadTestGenBuilder creates a new test generator builder with a real K8s client
func NewLoadTestGenBuilder(client K8sPodsInterface, namespace string) *LoadTestBuilder {
	return &LoadTestBuilder{
		namespace:       namespace,
		testSuiteName:   DefaultTestSuiteName,
		uniqPodLabelKey: DefaultUniqLabel,
		k8sClient:       client,
		includeWorkload: false,
		outputDir:       ".",
		moduleName:      "chaos-tests",
	}
}

// WithWorkload enables workload generation in the test
func (g *LoadTestBuilder) TestSuiteName(n string) *LoadTestBuilder {
	if !strings.HasPrefix(n, "Test") {
		n = fmt.Sprintf("Test%s", n)
	}
	g.testSuiteName = n
	return g
}

// Latency sets default latency for delay experiments
func (g *LoadTestBuilder) Latency(l int) *LoadTestBuilder {
	g.latencyMs = l
	return g
}

// Jitter sets default jitter for delay experiments
func (g *LoadTestBuilder) Jitter(j int) *LoadTestBuilder {
	g.jitterMs = j
	return g
}

// UniqPodLabelKey K8s Pod label key that uniqely identifies Pod for chaos experiment
// read more here https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func (g *LoadTestBuilder) UniqPodLabelKey(k string) *LoadTestBuilder {
	g.uniqPodLabelKey = k
	return g
}

// Workload enables workload generation in the test
func (g *LoadTestBuilder) Workload(i bool) *LoadTestBuilder {
	g.includeWorkload = i
	return g
}

// OutputDir sets the output directory for generated files
func (g *LoadTestBuilder) OutputDir(dir string) *LoadTestBuilder {
	g.outputDir = dir
	return g
}

// GoModName sets the Go module name for the generated project
func (g *LoadTestBuilder) GoModName(name string) *LoadTestBuilder {
	g.moduleName = name
	return g
}

// validate verifier that we can build codegen with provided params
// empty for now, add validation here if later it'd become more complex
func (g *LoadTestBuilder) validate() error {
	return nil
}

// Validate validate generation params
// for now it's empty but for more complex mutually exclusive cases we should
// add validation here
func (g *LoadTestBuilder) Build() (*LoadTestCodegen, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}
	return &LoadTestCodegen{g}, nil
}

// Read read K8s namespace and find all the pods
// some pods may be crashing but it doesn't matter for code generation
func (g *LoadTestCodegen) Read() error {
	log.Info().Str("", g.cfg.namespace).Msg("Scanning namespace for pods")

	pods, err := g.cfg.k8sClient.GetPods(context.Background(), g.cfg.namespace)
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	g.cfg.pods = pods.Items
	log.Info().Int("Pods", len(g.cfg.pods)).Msg("Found pods in namespace")

	for _, pod := range g.cfg.pods {
		log.Debug().
			Str("pod", pod.Name).
			Str("app", pod.Labels["app"]).
			Interface("labels", pod.Labels).
			Msg("Pod details")
	}

	return nil
}

// Write generates a complete boilerplate, can be multiple files
func (g *LoadTestCodegen) Write() error {
	// Create output directory
	if err := os.MkdirAll(g.cfg.outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate go.mod
	goModContent, err := g.GenerateGoMod()
	if err != nil {
		return err
	}
	goModPath := filepath.Join(g.cfg.outputDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0o600); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Generate main test file
	testContent, err := g.GenerateTableTest()
	if err != nil {
		return err
	}
	testPath := filepath.Join(g.cfg.outputDir, "chaos_test.go")
	if err := os.WriteFile(testPath, []byte(testContent), 0o600); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	// nolint
	defer os.Chdir(currentDir)
	if err := os.Chdir(g.cfg.outputDir); err != nil {
		return err
	}
	_, err = exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tidy generated module")
	}
	log.Info().
		Str("OutputDir", g.cfg.outputDir).
		Str("Module", g.cfg.moduleName).
		Bool("Workload", g.cfg.includeWorkload).
		Msg("Test module generated")
	return nil
}

// GenerateGoMod generates a go.mod file
func (g *LoadTestCodegen) GenerateGoMod() (string, error) {
	data := GoModParams{
		ModuleName: g.cfg.moduleName,
	}
	return render(GoModTemplate, data)
}

// GenerateTableTest generates all possible experiments for a namespace
// first generate all small pieces then insert into a table test template
func (g *LoadTestCodegen) GenerateTableTest() (string, error) {
	gunCode, workloadCode := "", ""
	if g.cfg.includeWorkload {
		gunCode, workloadCode = g.GenerateLoadTest()
	}
	testCases, err := g.GenerateTestCases()
	if err != nil {
		return "", err
	}

	data := TableTestParams{
		Package:       DefaultTestModule,
		Namespace:     g.cfg.namespace,
		TableTestName: g.cfg.testSuiteName,
		TestCases:     testCases,
		GunCode:       gunCode,
		WorkloadCode:  workloadCode,
	}
	return render(TableTestTmpl, data)
}

// GenerateLoadTest returns the workload generation code
func (g *LoadTestCodegen) GenerateLoadTest() (string, string) {
	return GunHTTPCode, WorkloadHTTPCode
}

// GenerateTestCases generates table test cases
func (g *LoadTestCodegen) GenerateTestCases() ([]TestCaseParams, error) {
	var testCases []TestCaseParams

	for _, pod := range g.cfg.pods {
		if _, ok := pod.Labels[g.cfg.uniqPodLabelKey]; !ok {
			return nil, fmt.Errorf("pod %s doesn't have uniq label key %s", pod.Name, g.cfg.uniqPodLabelKey)
		}
	}

	// Pod failures
	for _, pod := range g.cfg.pods {
		r, err := render(PodFailTmpl, PodFailParams{
			LabelKey:   g.cfg.uniqPodLabelKey,
			LabelValue: pod.Labels[g.cfg.uniqPodLabelKey],
		})
		if err != nil {
			return nil, err
		}
		testCases = append(testCases, TestCaseParams{
			Name:    fmt.Sprintf("Fail pod %s", pod.Name),
			RunFunc: r,
		})
	}

	// Pod latency
	for _, pod := range g.cfg.pods {
		r, err := render(PodDelayTmpl, PodDelayParams{
			LabelKey:   g.cfg.uniqPodLabelKey,
			LabelValue: pod.Labels[g.cfg.uniqPodLabelKey],
			LatencyMs:  g.cfg.latencyMs,
			JitterMs:   g.cfg.jitterMs,
		})
		if err != nil {
			return nil, err
		}

		testCases = append(testCases, TestCaseParams{
			Name:    fmt.Sprintf("Network delay for %s", pod.Name),
			RunFunc: r,
		})
	}

	return testCases, nil
}

// render is just an internal function to parse and render template
func render(tmpl string, data any) (string, error) {
	parsed, err := template.New("table_test").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse table test template: %w", err)
	}
	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to generate table test: %w", err)
	}
	return buf.String(), err
}
