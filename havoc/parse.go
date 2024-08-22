package havoc

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"os"
	"sort"
	"strings"
)

const (
	ErrNoNamespace    = "no namespace found"
	ErrEmptyNamespace = "no pods found inside namespace, namespace is empty or check your filter"
)

const (
	NoGroupKey = "no-group"
)

type ManifestPart struct {
	Kind              string
	Name              string
	FlattenedManifest map[string]interface{}
}

// PodsListResponse pod list response from kubectl in JSON
type PodsListResponse struct {
	Items []*PodResponse `json:"items"`
}

// PodResponse pod info response from kubectl in JSON
type PodResponse struct {
	Metadata struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

type GroupInfo struct {
	Label        string
	PodsAffected int
}

// ActionablePodInfo info about pod and labels for which we can generate a chaos experiment
type ActionablePodInfo struct {
	PodName  string
	Labels   []string
	HasGroup bool
}

func uniquePairs(strings []string) [][]string {
	var pairs [][]string
	for i := 0; i < len(strings); i++ {
		for j := i + 1; j < len(strings); j++ {
			pair := []string{strings[i], strings[j]}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

func (m *Controller) processPodInfoLo(plr *PodsListResponse) (map[string][]*PodResponse, []*PodResponse, []lo.Entry[string, int], [][]string, error) {
	L.Info().Msg("Processing pods info")
	// filtering
	filteredPods := lo.Filter(plr.Items, func(item *PodResponse, index int) bool {
		return !sliceContainsSubString(item.Metadata.Name, m.cfg.Havoc.IgnoredPods)
	})
	labelsToAllow := append([]string{}, m.cfg.Havoc.ComponentLabelKey)
	if m.hasNetworkExperiments() {
		labelsToAllow = append(labelsToAllow, m.cfg.Havoc.NetworkPartition.Label)
	}
	for _, p := range filteredPods {
		p.Metadata.Labels = lo.PickByKeys(p.Metadata.Labels, labelsToAllow)
	}
	if len(filteredPods) == 0 {
		return nil, nil, nil, nil, errors.New(ErrEmptyNamespace)
	}
	// grouping
	byComponent := lo.GroupBy(filteredPods, func(item *PodResponse) string {
		key := m.cfg.Havoc.ComponentLabelKey
		return m.labelSelector(key, item.Metadata.Labels[key])
	})
	var byPartition map[string][]*PodResponse
	if m.hasNetworkExperiments() {
		byPartition = lo.GroupBy(filteredPods, func(item *PodResponse) string {
			key := m.cfg.Havoc.NetworkPartition.Label
			return m.labelSelector(key, item.Metadata.Labels[key])
		})
	}
	componentGroupInfo := lo.MapEntries(byComponent, func(key string, value []*PodResponse) (string, int) {
		return key, len(value)
	})
	componentGroupsInfo := lo.Reject(lo.Entries(componentGroupInfo), func(item lo.Entry[string, int], index int) bool {
		return item.Key == NoGroupKey
	})
	byPartition = lo.OmitByKeys(byPartition, []string{NoGroupKey})
	partKeys := lo.Keys(byPartition)
	sort.Strings(partKeys)
	networkGroupsInfo := uniquePairs(partKeys)

	m.printPartitions(byComponent, "Component groups found")
	m.printPartitions(byPartition, "Network groups found")
	return byComponent, byComponent[NoGroupKey], componentGroupsInfo, networkGroupsInfo, nil
}

func (m *Controller) hasNetworkExperiments() bool {
	if m.cfg.Havoc.NetworkPartition != nil && m.cfg.Havoc.NetworkPartition.Label != "" {
		return true
	}
	return false
}

func (m *Controller) printPartitions(parts map[string][]*PodResponse, msg string) {
	for _, p := range parts {
		for _, pp := range p {
			L.Info().
				Str("Name", pp.Metadata.Name).
				Interface("Labels", pp.Metadata.Labels).
				Msg(msg)
		}
	}
}

// labelSelector transforms selector to ChaosMesh CRD format
func (m *Controller) labelSelector(k, v string) string {
	if v == "" {
		return NoGroupKey
	} else {
		return fmt.Sprintf("'%s': '%s'", k, v)
	}
}

// groupValueFromLabelSelector returns just the selector value
func (m *Controller) groupValueFromLabelSelector(selector string) string {
	val := strings.Split(selector, ": ")[1]
	return strings.ReplaceAll(val, "'", "")
}

// GetPodsInfo gets info about all the pods in the namespace
func (m *Controller) GetPodsInfo(namespace string) (*PodsListResponse, error) {
	if _, err := ExecCmd(fmt.Sprintf("kubectl get ns %s", namespace)); err != nil {
		return nil, errors.Wrap(errors.New(ErrNoNamespace), namespace)
	}
	var cmdBuilder strings.Builder
	cmdBuilder.Write([]byte(fmt.Sprintf("kubectl get pods -n %s ", namespace)))
	if m.cfg.Havoc.NamespaceLabelFilter != "" {
		cmdBuilder.Write([]byte(fmt.Sprintf("-l %s ", m.cfg.Havoc.NamespaceLabelFilter)))
	}
	cmdBuilder.Write([]byte("-o json"))
	out, err := ExecCmd(cmdBuilder.String())
	if err != nil {
		return nil, err
	}
	if err := dumpPodInfo(out); err != nil {
		return nil, err
	}
	var pr *PodsListResponse
	if err := json.Unmarshal([]byte(out), &pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func dumpPodInfo(out string) error {
	if L.GetLevel() == zerolog.DebugLevel {
		var plr *PodsListResponse
		if err := json.Unmarshal([]byte(out), &plr); err != nil {
			return err
		}
		d, err := json.Marshal(plr)
		if err != nil {
			return err
		}
		_ = os.WriteFile("pods_dump.json", d, os.ModePerm)
		return nil
	}
	return nil
}
