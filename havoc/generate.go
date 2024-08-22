package havoc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	ErrParsingTemplate = "failed to parse Go text template"

	ErrExperimentTimeout = "waiting for experiment to finish timed out"
	ErrExperimentApply   = "error applying experiment manifest"
	ErrInvalidCustomKind = "invalid custom Kind of experiment"
)

const (
	DebugContainerImage = "curlimages/curl:latest"
)

var (
	RecommendedExperimentTypes = []string{
		ChaosTypeFailure,
		ChaosTypeLatency,
		ChaosTypeGroupFailure,
		ChaosTypeGroupLatency,
		ChaosTypeStressMemory,
		ChaosTypeStressGroupMemory,
		ChaosTypeStressCPU,
		ChaosTypeStressGroupCPU,
		ChaosTypePartitionGroup,
		ChaosTypeHTTP,
		//ChaosTypePartitionExternal,
	}
)

// MarshalTemplate Helper to marshal templates
func MarshalTemplate(jobSpec interface{}, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", errors.Wrap(err, ErrParsingTemplate)
	}
	err = tmpl.Execute(&buf, jobSpec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

type CommonExperimentMeta struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

type HTTPExperiment struct {
	ExperimentName string
	Metadata       *Metadata
	Namespace      string
	Mode           string
	ModeValue      string
	Selector       string
	PodName        string
	Port           int64
	Target         string
	Path           string
	Method         string
	Abort          bool
	Duration       string
}

func (m HTTPExperiment) String() (string, error) {
	tpl := `
kind: HTTPChaos
apiVersion: chaos-mesh.org/v1alpha1
metadata:
  name: {{ .ExperimentName }}
spec:
  mode: {{ .Mode }}
  {{- if .ModeValue }}
  value: '{{ .ModeValue }}'
  {{- end }}
  selector:
    namespaces:
      - {{ .Namespace }}
    {{- if .Selector}}
    labelSelectors:
      {{ .Selector }}
	{{- else}}
    fieldSelectors:
      metadata.name: {{ .PodName }}
	{{- end}}
  target: Request
  port: {{ .Port }}
  method: {{ .Method }}
  path: {{ .Path }}
  abort: {{ .Abort }}
  duration: {{ .Duration }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type BlockchainRewindHeadExperiment struct {
	ExperimentName        string    `yaml:"experimentName"`
	Metadata              *Metadata `yaml:"metadata"`
	Namespace             string    `yaml:"namespace"`
	PodName               string    `yaml:"podName"`
	ExecutorPodPrefix     string    `yaml:"executorPodPrefix"`
	ExecutorContainerName string    `yaml:"executorContainerName"`
	NodeInternalHTTPURL   string    `yaml:"nodeInternalHTTPURL"`
	Blocks                int64     `yaml:"blocks"`
}

type Metadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

func (m BlockchainRewindHeadExperiment) String() (string, error) {
	tpl := `
kind: blockchain_rewind_head
name: {{ .ExperimentName }}
metadata:
  name: {{ .Metadata.Name }}
podName: {{ .PodName }}
executorContainerName: {{ .ExecutorContainerName }}
nodeInternalHTTPURL: {{ .NodeInternalHTTPURL }}
namespace: {{ .Namespace }}
blocks: {{ .Blocks }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type NetworkChaosExperiment struct {
	ExperimentName string
	Mode           string
	ModeValue      string
	Namespace      string
	Duration       string
	Latency        string
	PodName        string
	Selector       string
}

func (m NetworkChaosExperiment) String() (string, error) {
	tpl := `
kind: NetworkChaos
apiVersion: chaos-mesh.org/v1alpha1
metadata:
  name: {{ .ExperimentName }}
  namespace: {{ .Namespace }}
spec:
  selector:
    namespaces:
      - {{ .Namespace }}
    {{- if .Selector}}
    labelSelectors:
      {{ .Selector }}
	{{- else}}
    fieldSelectors:
      metadata.name: {{ .PodName }}
	{{- end}}
  mode: {{ .Mode }}
  {{- if .ModeValue }}
  value: '{{ .ModeValue }}'
  {{- end }}
  action: delay
  duration: {{ .Duration }}
  delay:
    latency: {{ .Latency }}
  direction: from
  target:
    selector:
      namespaces:
        - {{ .Namespace }}
      {{- if .Selector}}
      labelSelectors:
        {{ .Selector }}
	  {{- else}}
      fieldSelectors:
        metadata.name: {{ .PodName }}
	  {{- end}}
    mode: {{ .Mode }}
    {{- if .ModeValue }}
    value: '{{ .ModeValue }}'
    {{- end }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type NetworkChaosGroupPartitionExperiment struct {
	ExperimentName string
	ModeTo         string
	ModeToValue    string
	ModeFrom       string
	ModeFromValue  string
	Direction      string
	Namespace      string
	Duration       string
	SelectorFrom   string
	SelectorTo     string
}

func (m NetworkChaosGroupPartitionExperiment) String() (string, error) {
	tpl := `
kind: NetworkChaos
apiVersion: chaos-mesh.org/v1alpha1
metadata:
  name: {{ .ExperimentName }}
  namespace: {{ .Namespace }}
spec:
  selector:
    namespaces:
      - {{ .Namespace }}
    labelSelectors:
      {{ .SelectorFrom }}
  action: partition
  mode: {{ .ModeFrom }}
  {{- if .ModeFromValue }}
  value: '{{ .ModeFromValue }}'
  {{- end }}
  duration: {{ .Duration }}
  direction: {{ .Direction }}
  target:
    mode: {{ .ModeTo }}
    {{- if .ModeToValue }}
    value: '{{ .ModeToValue }}'
    {{- end }}
    selector:
      namespaces:
        - {{ .Namespace }}
      labelSelectors:
        {{ .SelectorTo }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type NetworkChaosExternalPartitionExperiment struct {
	ExperimentName string
	Namespace      string
	Duration       string
	PodName        string
	ExternalURL    string
}

func (m NetworkChaosExternalPartitionExperiment) String() (string, error) {
	tpl := `
kind: NetworkChaos
apiVersion: chaos-mesh.org/v1alpha1
metadata:
  name: {{ .ExperimentName }}
  namespace: {{ .Namespace }}
spec:
  selector:
    namespaces:
      - {{ .Namespace }}
  mode: all
  action: partition
  duration: {{ .Duration }}
  direction: to
  target:
    selector:
      namespaces:
        - {{ .Namespace }}
    mode: all
  externalTargets:
    - {{ .ExternalURL }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type PodFailureExperiment struct {
	ExperimentName string
	Mode           string
	ModeValue      string
	Namespace      string
	Duration       string
	PodName        string
	Selector       string
}

func (m PodFailureExperiment) String() (string, error) {
	tpl := `
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: {{ .ExperimentName }}
  namespace: {{ .Namespace }}
spec:
  action: pod-failure
  mode: {{ .Mode }}
  {{- if .ModeValue }}
  value: '{{ .ModeValue }}'
  {{- end }}
  duration: {{ .Duration }}
  selector:
    {{- if .Selector}}
    labelSelectors:
      {{ .Selector }}
	{{- else}}
    fieldSelectors:
      metadata.name: {{ .PodName }}
	{{- end}}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type PodStressCPUExperiment struct {
	ExperimentName string
	Mode           string
	ModeValue      string
	Namespace      string
	Workers        int
	Load           int
	Duration       string
	PodName        string
	Selector       string
}

func (m PodStressCPUExperiment) String() (string, error) {
	tpl := `
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: {{ .ExperimentName }}
  namespace: {{ .Namespace }}
spec:
  mode: {{ .Mode }}
  {{- if .ModeValue }}
  value: '{{ .ModeValue }}'
  {{- end }}
  duration: {{ .Duration }}
  selector:
    {{- if .Selector}}
    labelSelectors:
      {{ .Selector }}
	{{- else}}
    fieldSelectors:
      metadata.name: {{ .PodName }}
	{{- end}}
  stressors:
    cpu:
      workers: {{ .Workers }}
      load: {{ .Load }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type PodStressMemoryExperiment struct {
	ExperimentName string
	Mode           string
	ModeValue      string
	Namespace      string
	Workers        int
	Memory         string
	Duration       string
	PodName        string
	Selector       string
}

func (m PodStressMemoryExperiment) String() (string, error) {
	tpl := `
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: {{ .ExperimentName }}
  namespace: {{ .Namespace }}
spec:
  mode: {{ .Mode }}
  {{- if .ModeValue }}
  value: '{{ .ModeValue }}'
  {{- end }}
  duration: {{ .Duration }}
  selector:
    {{- if .Selector}}
    labelSelectors:
      {{ .Selector }}
	{{- else}}
    fieldSelectors:
      metadata.name: {{ .PodName }}
	{{- end}}
  stressors:
    memory:
      workers: {{ .Workers }}
      size: {{ .Memory }}
`
	return MarshalTemplate(
		m,
		uuid.NewString(),
		tpl,
	)
}

type CRD struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec interface{} `yaml:"spec"` // Use interface{} if the spec can have various structures
}

type NamedExperiment struct {
	CRD
	Name     string
	Path     string
	CRDBytes []byte
}

func NewNamedExperiment(expPath string) (*NamedExperiment, error) {
	data, err := os.ReadFile(expPath)
	if err != nil {
		return nil, err
	}

	var exp CRD
	err = yaml.Unmarshal(data, &exp)
	if err != nil {
		return nil, err
	}
	expName := exp.Metadata.Name
	if expName == "" {
		return nil, errors.Errorf("experiment metadata.name is empty")
	}

	return &NamedExperiment{
		CRD:      exp,
		Name:     expName,
		Path:     expPath,
		CRDBytes: data,
	}, nil
}

func (m *Controller) readExistingExperimentTypes(dir string) ([]string, error) {
	expTypes := make([]string, 0)
	err := filepath.Walk(
		dir,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() != dir {
				expTypes = append(expTypes, info.Name())
				return nil
			}
			return err
		})
	if err != nil {
		return nil, err
	}
	sort.Slice(expTypes, func(i, j int) bool {
		return expTypes[i] < expTypes[j]
	})
	L.Info().Strs("Order", expTypes).Msg("Order of experiment dirs execution")
	return expTypes, nil
}

func (m *Controller) ReadExperimentsFromDir(expTypes []string, dir string) ([]*NamedExperiment, error) {
	expData := make([]*NamedExperiment, 0)
	for _, expType := range expTypes {
		targetDir := fmt.Sprintf("%s/%s", dir, expType)
		if _, err := os.Stat(targetDir); err != nil {
			// it's okay, some experiments may be skipped due configuration
			continue
		}
		err := filepath.Walk(
			fmt.Sprintf("%s/%s", dir, expType),
			func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				exp, err := NewNamedExperiment(path)
				if err != nil {
					return err
				}
				expData = append(expData, exp)
				return err
			})
		if err != nil {
			return nil, err
		}
	}
	return expData, nil
}

// maybeFailAll is a special case where we've labelled component properly but
// there is only one component, no need to apply experiment to part of a group, so we set 100%
func maybeFailAll(e lo.Entry[string, int], origValue string) string {
	if e.Value == 1 {
		return "100"
	}
	return origValue
}

func (m *Controller) generate(
	namespace string,
	oapiSpecs []*OAPISpecData,
	allPodsInfo map[string][]*PodResponse,
	podsInfo []*PodResponse,
	groupLabels []lo.Entry[string, int],
	netLabels [][]string,
) (*ChaosSpecs, error) {
	allExperimentsByType := make(map[string]map[string]string)
	for _, expType := range m.cfg.Havoc.ExperimentTypes {
		experiments := make(map[string]string)
		switch expType {
		case ChaosTypeHTTP:
			for _, entry := range groupLabels {
				if _, ok := m.cfg.Havoc.OpenAPI.Mapping[m.groupValueFromLabelSelector(entry.Key)]; ok {
					if err := m.generateOAPIExperiments(experiments, namespace, entry, oapiSpecs); err != nil {
						return nil, err
					}
				}
			}
		case ChaosTypeBlockchainSetHead:
			for _, p := range allPodsInfo {
				for _, pi := range p {
					for _, nodeCfg := range m.cfg.Havoc.BlockchainRewindHead.NodesConfig {
						if strings.Contains(pi.Metadata.Name, nodeCfg.ExecutorPodPrefix) {
							for _, b := range nodeCfg.Blocks {
								name := fmt.Sprintf("%s-%s-%d", ChaosTypeBlockchainSetHead, pi.Metadata.Name, b)
								experiment, err := BlockchainRewindHeadExperiment{
									ExperimentName:        name,
									Metadata:              &Metadata{Name: name},
									Namespace:             namespace,
									NodeInternalHTTPURL:   nodeCfg.NodeInternalHTTPURL,
									PodName:               pi.Metadata.Name,
									ExecutorContainerName: nodeCfg.ExecutorContainerName,
									Blocks:                b,
								}.String()
								if err != nil {
									return nil, err
								}
								shortName := fmt.Sprintf("%s-%d", pi.Metadata.Name, b)
								experiments[shortName] = experiment
							}
						}
					}
				}
			}
		case ChaosTypePartitionExternal:
			if m.cfg.Havoc.ExternalTargets == nil {
				continue
			}
			for _, u := range m.cfg.Havoc.ExternalTargets.URLs {
				nsAndURLHash := fmt.Sprintf("%s-%s", namespace, urlHash(u))
				experiment, err := NetworkChaosExternalPartitionExperiment{
					Namespace:      namespace,
					ExperimentName: fmt.Sprintf("%s-%s", ChaosTypePartitionExternal, nsAndURLHash),
					Duration:       m.cfg.Havoc.ExternalTargets.Duration,
					ExternalURL:    fmt.Sprintf("'%s'", u),
				}.String()
				if err != nil {
					return nil, err
				}
				experiments[nsAndURLHash] = experiment
			}
		case ChaosTypePartitionGroup:
			for _, pair := range netLabels {
				for _, groupModeValue := range m.cfg.Havoc.NetworkPartition.GroupPercentage {
					sanitizedLabel := sanitizeLabel(fmt.Sprintf("%s-to-%s", pair[0], pair[1]))
					sanitizedLabel = fmt.Sprintf("%s-%s-perc", sanitizedLabel, groupModeValue)
					experiment, err := NetworkChaosGroupPartitionExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypePartitionGroup, sanitizedLabel),
						Duration:       m.cfg.Havoc.NetworkPartition.Duration,
						ModeFrom:       "fixed-percent",
						ModeFromValue:  groupModeValue,
						ModeTo:         "fixed-percent",
						ModeToValue:    groupModeValue,
						Direction:      "from",
						SelectorFrom:   pair[0],
						SelectorTo:     pair[1],
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
				for _, groupModeValue := range m.cfg.Havoc.NetworkPartition.GroupFixed {
					sanitizedLabel := sanitizeLabel(fmt.Sprintf("%s-to-%s", pair[0], pair[1]))
					sanitizedLabel = fmt.Sprintf("%s-%s-fixed", sanitizedLabel, groupModeValue)
					experiment, err := NetworkChaosGroupPartitionExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypePartitionGroup, sanitizedLabel),
						Duration:       m.cfg.Havoc.NetworkPartition.Duration,
						ModeFrom:       "fixed-percent",
						ModeFromValue:  groupModeValue,
						ModeTo:         "fixed-percent",
						ModeToValue:    groupModeValue,
						Direction:      "from",
						SelectorFrom:   pair[0],
						SelectorTo:     pair[1],
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
			}
		case ChaosTypeFailure:
			for _, pi := range podsInfo {
				experiment, err := PodFailureExperiment{
					Namespace:      namespace,
					ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeFailure, pi.Metadata.Name),
					Mode:           "one",
					Duration:       m.cfg.Havoc.Failure.Duration,
					PodName:        pi.Metadata.Name,
				}.String()
				if err != nil {
					return nil, err
				}
				experiments[pi.Metadata.Name] = experiment
			}
		case ChaosTypeLatency:
			for _, pi := range podsInfo {
				experiment, err := NetworkChaosExperiment{
					Namespace:      namespace,
					ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeLatency, pi.Metadata.Name),
					Mode:           "one",
					Duration:       m.cfg.Havoc.Latency.Duration,
					Latency:        m.cfg.Havoc.Latency.Latency,
					PodName:        pi.Metadata.Name,
				}.String()
				if err != nil {
					return nil, err
				}
				experiments[pi.Metadata.Name] = experiment
			}
		case ChaosTypeStressCPU:
			for _, pi := range podsInfo {
				experiment, err := PodStressCPUExperiment{
					Namespace:      namespace,
					ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeStressCPU, pi.Metadata.Name),
					Duration:       m.cfg.Havoc.StressCPU.Duration,
					Workers:        m.cfg.Havoc.StressCPU.Workers,
					Load:           m.cfg.Havoc.StressCPU.Load,
					Mode:           "one",
					PodName:        pi.Metadata.Name,
				}.String()
				if err != nil {
					return nil, err
				}
				experiments[pi.Metadata.Name] = experiment
			}
		case ChaosTypeStressMemory:
			for _, pi := range podsInfo {
				experiment, err := PodStressMemoryExperiment{
					Namespace:      namespace,
					ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeStressMemory, pi.Metadata.Name),
					Duration:       m.cfg.Havoc.StressMemory.Duration,
					Workers:        m.cfg.Havoc.StressMemory.Workers,
					Memory:         m.cfg.Havoc.StressMemory.Memory,
					Mode:           "one",
					PodName:        pi.Metadata.Name,
				}.String()
				if err != nil {
					return nil, err
				}
				experiments[pi.Metadata.Name] = experiment
			}
		case ChaosTypeStressGroupMemory:
			for _, entry := range groupLabels {
				for _, groupModeValue := range m.cfg.Havoc.StressMemory.GroupPercentage {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-perc", sanitizedLabel, groupModeValue)
					experiment, err := PodStressMemoryExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeStressGroupMemory, sanitizedLabel),
						Duration:       m.cfg.Havoc.StressMemory.Duration,
						Workers:        m.cfg.Havoc.StressMemory.Workers,
						Memory:         m.cfg.Havoc.StressMemory.Memory,
						Mode:           "fixed-percent",
						ModeValue:      groupModeValue,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
				for _, groupModeValue := range m.cfg.Havoc.StressMemory.GroupFixed {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-fixed", sanitizedLabel, groupModeValue)
					experiment, err := PodStressMemoryExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeStressGroupMemory, sanitizedLabel),
						Duration:       m.cfg.Havoc.StressMemory.Duration,
						Workers:        m.cfg.Havoc.StressMemory.Workers,
						Memory:         m.cfg.Havoc.StressMemory.Memory,
						Mode:           "fixed",
						ModeValue:      groupModeValue,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
			}
		case ChaosTypeStressGroupCPU:
			for _, entry := range groupLabels {
				for _, groupModeValue := range m.cfg.Havoc.StressCPU.GroupPercentage {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-perc", sanitizedLabel, groupModeValue)
					experiment, err := PodStressCPUExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeStressGroupCPU, sanitizedLabel),
						Duration:       m.cfg.Havoc.StressCPU.Duration,
						Workers:        m.cfg.Havoc.StressCPU.Workers,
						Load:           m.cfg.Havoc.StressCPU.Load,
						Mode:           "fixed-percent",
						ModeValue:      groupModeValue,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
				for _, groupModeValue := range m.cfg.Havoc.StressCPU.GroupFixed {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-fixed", sanitizedLabel, groupModeValue)
					experiment, err := PodStressCPUExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeStressGroupCPU, sanitizedLabel),
						Duration:       m.cfg.Havoc.StressCPU.Duration,
						Workers:        m.cfg.Havoc.StressCPU.Workers,
						Load:           m.cfg.Havoc.StressCPU.Load,
						Mode:           "fixed",
						ModeValue:      groupModeValue,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
			}
		case ChaosTypeGroupFailure:
			for _, entry := range groupLabels {
				for _, groupModeValue := range m.cfg.Havoc.Failure.GroupPercentage {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-perc", sanitizedLabel, groupModeValue)
					experiment, err := PodFailureExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeGroupFailure, sanitizedLabel),
						Duration:       m.cfg.Havoc.Failure.Duration,
						Mode:           "fixed-percent",
						ModeValue:      groupModeValue,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
				for _, groupModeValue := range m.cfg.Havoc.Failure.GroupFixed {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-fixed", sanitizedLabel, groupModeValue)
					experiment, err := PodFailureExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeGroupFailure, sanitizedLabel),
						Duration:       m.cfg.Havoc.Failure.Duration,
						Mode:           "fixed",
						ModeValue:      groupModeValue,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
			}
		case ChaosTypeGroupLatency:
			for _, entry := range groupLabels {
				for _, groupModeValue := range m.cfg.Havoc.Latency.GroupPercentage {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-perc", sanitizedLabel, groupModeValue)
					experiment, err := NetworkChaosExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeGroupLatency, sanitizedLabel),
						Mode:           "fixed-percent",
						ModeValue:      groupModeValue,
						Duration:       m.cfg.Havoc.Latency.Duration,
						Latency:        m.cfg.Havoc.Latency.Latency,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
				for _, groupModeValue := range m.cfg.Havoc.Latency.GroupFixed {
					groupModeValue = maybeFailAll(entry, groupModeValue)
					sanitizedLabel := sanitizeLabel(entry.Key)
					sanitizedLabel = fmt.Sprintf("%s-%s-fixed", sanitizedLabel, groupModeValue)
					experiment, err := NetworkChaosExperiment{
						Namespace:      namespace,
						ExperimentName: fmt.Sprintf("%s-%s", ChaosTypeGroupLatency, sanitizedLabel),
						Mode:           "fixed",
						ModeValue:      groupModeValue,
						Duration:       m.cfg.Havoc.Latency.Duration,
						Latency:        m.cfg.Havoc.Latency.Latency,
						Selector:       entry.Key,
					}.String()
					if err != nil {
						return nil, err
					}
					experiments[sanitizedLabel] = experiment
				}
			}
		}
		allExperimentsByType[expType] = experiments
	}
	return &ChaosSpecs{
		ExperimentsByType: allExperimentsByType,
	}, nil
}

func urlHash(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func sanitizeLabel(label string) string {
	sanitizedLabel := strings.Replace(label, "'", "", -1)
	sanitizedLabel = strings.Replace(sanitizedLabel, ": ", "-", -1)
	sanitizedLabel = strings.Replace(sanitizedLabel, ".", "-", -1)
	sanitizedLabel = strings.Replace(sanitizedLabel, "/", "-", -1)
	return sanitizedLabel
}

type EventJSONItemResponse struct {
	APIVersion     string    `json:"apiVersion"`
	Count          int       `json:"count"`
	EventTime      any       `json:"eventTime"`
	FirstTimestamp time.Time `json:"firstTimestamp"`
	InvolvedObject struct {
		APIVersion      string `json:"apiVersion"`
		Kind            string `json:"kind"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"involvedObject"`
	Kind          string    `json:"kind"`
	LastTimestamp time.Time `json:"lastTimestamp"`
	Message       string    `json:"message"`
	Metadata      struct {
		Annotations struct {
			ChaosMeshOrgType string `json:"chaos-mesh.org/type"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		ResourceVersion   string    `json:"resourceVersion"`
		UID               string    `json:"uid"`
	} `json:"metadata"`
	Reason             string `json:"reason"`
	ReportingComponent string `json:"reportingComponent"`
	ReportingInstance  string `json:"reportingInstance"`
	Source             struct {
		Component string `json:"component"`
	} `json:"source"`
	Type string `json:"type"`
}

type EventsJSONResponse struct {
	APIVersion string                   `json:"apiVersion"`
	Items      []*EventJSONItemResponse `json:"items"`
	Kind       string                   `json:"kind"`
	Metadata   struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

func eventsForLastMinutes(out string, timeOfApplication time.Time) error {
	var d *EventsJSONResponse
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		return err
	}
	L.Debug().Msg("Listing all experiment events")
	for _, i := range d.Items {
		if i.LastTimestamp.After(timeOfApplication) {
			L.Info().
				Time("Time", i.LastTimestamp).
				Str("Reason", i.Reason).
				Str("Message", i.Message).
				Send()
		}
	}
	return nil
}

func (m *Controller) ApplyExperiment(exp *NamedExperiment, wait bool) error {
	timeOfApplication := time.Now()
	var errDefer error
	if exp.Kind == ChaosTypeBlockchainSetHead {
		return m.ApplyCustomKindChaosFile(exp, ChaosTypeBlockchainSetHead, wait)
	}
	L.Info().
		Str("Dir", m.cfg.Havoc.Dir).
		Str("Type", exp.Kind).
		Str("Name", exp.Metadata.Name).
		Msg("Applying experiment manifest")
	fmt.Println(string(exp.CRDBytes))
	_, err := ExecCmd(fmt.Sprintf("kubectl apply -f %s", exp.Path))
	if err != nil {
		return errors.Wrap(err, ErrExperimentApply)
	}
	if wait {
		resourceType := ExperimentTypesToCRDNames[exp.Kind]
		if resourceType == "" {
			return errors.Errorf("%s resource not present in %+v list", exp.Kind, ExperimentTypesToCRDNames)
		}
		_, err = ExecCmd(
			fmt.Sprintf("kubectl wait -n %s %s --field-selector=metadata.name=%s --for condition=AllRecovered=True --timeout %s",
				exp.Metadata.Namespace,
				resourceType,
				exp.Metadata.Name,
				DefaultCMDTimeout,
			))
		if err != nil {
			return errors.Wrap(err, ErrExperimentTimeout)
		}
		out, err := ExecCmd(
			fmt.Sprintf("kubectl get -n %s events --field-selector involvedObject.name=%s -o json",
				exp.Metadata.Namespace,
				exp.Name,
			))
		if err != nil {
			return err
		}
		if err = eventsForLastMinutes(out, timeOfApplication); err != nil {
			return err
		}
		_, err = ExecCmd(fmt.Sprintf("kubectl -n %s delete %s %s", exp.Metadata.Namespace, resourceType, exp.Name))
		if err != nil {
			return err
		}
		L.Info().Msg("Chaos experiment successfully recovered")
	}
	return errDefer
}

type CurrentBlockResponse struct {
	Result string `json:"result"`
}

func (m *Controller) ApplyCustomKindChaosFile(exp *NamedExperiment, chaosType string, wait bool) error {
	switch chaosType {
	case ChaosTypeBlockchainSetHead:
		var rewind *BlockchainRewindHeadExperiment
		data, err := os.ReadFile(exp.Path)
		if err != nil {
			return nil
		}
		if err := yaml.Unmarshal(data, &rewind); err != nil {
			return err
		}
		L.Info().
			Str("Dir", m.cfg.Havoc.Dir).
			Str("Type", chaosType).
			Str("Name", exp.Name).
			Msg("Applying custom experiment")
		fmt.Println(string(exp.CRDBytes))
		lastBlkCommand := fmt.Sprintf(`kubectl -n %s -it debug %s --image=%s --target=%s -- curl -s -X POST -H Content-Type:application/json --data {"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":88} %s`,
			rewind.Namespace,
			rewind.PodName,
			DebugContainerImage,
			rewind.ExecutorContainerName,
			rewind.NodeInternalHTTPURL,
		)
		out, err := ExecCmd(lastBlkCommand)
		if err != nil {
			return err
		}
		msg, err := findJSONMsg(out)
		if err != nil {
			return err
		}
		var res *CurrentBlockResponse
		if err := json.Unmarshal([]byte(msg), &res); err != nil {
			return err
		}
		decimalLastBlock, err := strconv.ParseInt(res.Result[2:], 16, 64)
		if err != nil {
			return err
		}
		moveToBlock := decimalLastBlock - rewind.Blocks
		moveToBlockHex := strconv.FormatInt(moveToBlock, 16)
		setHeadCommand := fmt.Sprintf(`kubectl -n %s -it debug %s --image=%s --target=%s -- curl -s -X POST -H Content-Type:application/json --data {"jsonrpc":"2.0","method":"debug_setHead","params":["0x%s"],"id":5} %s`,
			rewind.Namespace,
			rewind.PodName,
			DebugContainerImage,
			rewind.ExecutorContainerName,
			moveToBlockHex,
			rewind.NodeInternalHTTPURL,
		)
		_, err = ExecCmd(setHeadCommand)
		if err != nil {
			return err
		}
	default:
		return errors.New(ErrInvalidCustomKind)
	}
	return nil
}

func findJSONMsg(s string) (string, error) {
	startIndex := strings.Index(s, "{")
	endIndex := strings.LastIndex(s, "}")
	if startIndex != -1 && endIndex != -1 {
		substring := s[startIndex : endIndex+1]
		L.Debug().
			Str("Message", substring).
			Msg("JSON substring response")
		return substring, nil
	} else {
		return "", errors.New("no JSON substring found in response")
	}
}

// GenerateSpecs generates specs from namespace, should be used programmatically in tests
func (m *Controller) GenerateSpecs(ns string) error {
	podsInfo, err := m.GetPodsInfo(ns)
	if err != nil {
		return err
	}
	_, _, err = m.generateSpecs(ns, podsInfo)
	return err
}

func (m *Controller) generateSpecs(namespace string, podListResponse *PodsListResponse) (*ChaosSpecs, []*PodResponse, error) {
	L.Trace().
		Interface("PodListResponse", podListResponse).
		Msg("Found pods")
	all, noGroup, componentLabels, networkLabels, err := m.processPodInfoLo(podListResponse)
	if err != nil {
		return nil, nil, err
	}
	L.Info().Msg("Processing OpenAPI specs")
	specs, err := m.ParseOpenAPISpecs()
	if err != nil {
		return nil, nil, err
	}
	L.Info().Msg("Generating chaos experiments")
	csp, err := m.generate(namespace, specs, all, noGroup, componentLabels, networkLabels)
	if err != nil {
		return nil, nil, err
	}
	return csp, noGroup, csp.Dump(m.cfg.Havoc.Dir)
}
