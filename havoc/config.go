package havoc

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

const (
	ErrReadSethConfig      = "failed to read TOML config for havoc"
	ErrUnmarshalSethConfig = "failed to unmarshal TOML config for havoc"

	ErrFailureGroupIsNil      = "failure group must be specified in config"
	ErrLatencyGroupIsNil      = "latency group must be specified in config"
	ErrStressCPUGroupIsNil    = "stress cpu group must be specified in config"
	ErrStressMemoryGroupIsNil = "stress memory group must be specified in config"
	ErrFormat                 = "format error"
)

const (
	DefaultExperimentsDir           = "havoc-experiments"
	DefaultPodFailureDuration       = "1m"
	DefaultNetworkLatencyDuration   = "1m"
	DefaultNetworkPartitionDuration = "1m"
	DefaultHTTPDuration             = "1m"
	DefaultNetworkPartitionLabel    = "havoc-network-group"
	DefaultComponentGroupLabelKey   = "havoc-component-group"
	DefaultStressMemoryDuration     = "1m"
	DefaultStressMemoryWorkers      = 1
	DefaultStressMemoryAmount       = "512MB"
	DefaultStressCPUDuration        = "1m"
	DefaultStressCPUWorkers         = 1
	DefaultStressCPULoad            = 100
	DefaultNetworkLatency           = "300ms"
	DefaultMonkeyDuration           = "24h"
	DefaultMonkeyMode               = "seq"
	DefaultMonkeyCooldown           = "30s"
)

var (
	DefaultGroupPercentage                 = []string{"10", "20", "30"}
	DefaultGroupFixed                      = []string{"1", "2", "3"}
	DefaultNetworkPartitionGroupPercentage = []string{"100"}
)

var (
	DefaultIgnoreGroupLabels = []string{
		"mainnet",
		"release",
		"intents.otterize.com",
		"pod-template-hash",
		"rollouts-pod-template-hash",
		"chain.link/app",
		"chain.link/cost-center",
		"chain.link/env",
		"chain.link/project",
		"chain.link/team",
		"app.kubernetes.io/part-of",
		"app.kubernetes.io/managed-by",
		"app.chain.link/product",
		"app.kubernetes.io/version",
		"app.chain.link/blockchain",
		"app.kubernetes.io/instance",
		"app.kubernetes.io/name",
	}
)

type Config struct {
	Havoc *Havoc `toml:"havoc"`
}

type Havoc struct {
	Dir                  string                `toml:"dir"`
	ExperimentTypes      []string              `toml:"experiment_types"`
	NamespaceLabelFilter string                `toml:"namespace_label_filter"`
	ComponentLabelKey    string                `toml:"component_label_key"`
	IgnoredPods          []string              `toml:"ignore_pods"`
	IgnoreGroupLabels    []string              `toml:"ignore_group_labels"`
	Failure              *Failure              `toml:"failure"`
	Latency              *Latency              `toml:"latency"`
	NetworkPartition     *NetworkPartition     `toml:"network_partition"`
	StressMemory         *StressMemory         `toml:"stress_memory"`
	StressCPU            *StressCPU            `toml:"stress_cpu"`
	ExternalTargets      *ExternalTargets      `toml:"external_targets"`
	BlockchainRewindHead *BlockchainRewindHead `toml:"blockchain_rewind_head"`
	OpenAPI              *OpenAPI              `toml:"openapi"`
	Monkey               *Monkey               `toml:"monkey"`
	Grafana              *Grafana              `toml:"grafana"`
}

func dumpConfig(cfg *Config) {
	if L.GetLevel() == zerolog.DebugLevel {
		d, _ := toml.Marshal(cfg)
		_ = os.WriteFile("config_dump.toml", d, os.ModePerm)
	}
}

func DefaultConfig() *Config {
	return &Config{
		Havoc: &Havoc{
			Dir:               DefaultExperimentsDir,
			ExperimentTypes:   RecommendedExperimentTypes,
			ComponentLabelKey: DefaultComponentGroupLabelKey,
			IgnoreGroupLabels: DefaultIgnoreGroupLabels,
			Failure: &Failure{
				Duration:   DefaultPodFailureDuration,
				GroupFixed: DefaultGroupFixed,
			},
			Latency: &Latency{
				Duration:   DefaultNetworkLatencyDuration,
				Latency:    DefaultNetworkLatency,
				GroupFixed: DefaultGroupFixed,
			},
			StressMemory: &StressMemory{
				Duration:   DefaultStressMemoryDuration,
				Workers:    DefaultStressMemoryWorkers,
				Memory:     DefaultStressMemoryAmount,
				GroupFixed: DefaultGroupFixed,
			},
			StressCPU: &StressCPU{
				Duration:   DefaultStressCPUDuration,
				Workers:    DefaultStressCPUWorkers,
				Load:       DefaultStressCPULoad,
				GroupFixed: DefaultGroupFixed,
			},
			NetworkPartition: &NetworkPartition{
				Duration:        DefaultNetworkPartitionDuration,
				Label:           DefaultNetworkPartitionLabel,
				GroupPercentage: DefaultNetworkPartitionGroupPercentage,
			},
			OpenAPI: &OpenAPI{
				Duration:   DefaultHTTPDuration,
				GroupFixed: DefaultGroupFixed,
			},
			Monkey: &Monkey{
				Duration: DefaultMonkeyDuration,
				Mode:     DefaultMonkeyMode,
				Cooldown: DefaultMonkeyCooldown,
			},
			Grafana: &Grafana{
				URL:   os.Getenv("GRAFANA_URL"),
				Token: os.Getenv("GRAFANA_TOKEN"),
			},
		},
	}
}

func (c *Config) Validate() []error {
	errs := make([]error, 0)
	if c.Havoc.Dir == "" {
		errs = append(errs, errors.Wrap(errors.New(ErrFormat), "monkey.dir must not be empty"))
	}
	if c.Havoc.Failure == nil {
		errs = append(errs, errors.New(ErrFailureGroupIsNil))
	}
	if c.Havoc.Latency == nil {
		errs = append(errs, errors.New(ErrLatencyGroupIsNil))
	}
	if c.Havoc.StressCPU == nil {
		errs = append(errs, errors.New(ErrStressCPUGroupIsNil))
	}
	if c.Havoc.StressMemory == nil {
		errs = append(errs, errors.New(ErrStressMemoryGroupIsNil))
	}
	if c.Havoc.Failure != nil {
		if c.Havoc.Failure.Duration == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "failure.duration must be in Go duration format, 1d2h3m0s"))
		}
	}
	if c.Havoc.Latency != nil {
		if c.Havoc.Latency.Duration == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "latency.duration must be in Go duration format, 1d2h3m0s"))
		}
		if c.Havoc.Latency.Latency == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "latency.latency must be in milliseconds format, ex.: 300ms"))
		}
	}
	if c.Havoc.StressMemory != nil {
		if c.Havoc.StressMemory.Workers <= 0 {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "stress_memory.workers must be set, ex.: \"4\""))
		}
		if c.Havoc.StressMemory.Memory == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "stress_memory.memory must be set, ex.: \"256MB\" or \"25%\""))
		}
	}
	if c.Havoc.StressCPU != nil {
		if c.Havoc.StressCPU.Workers <= 0 {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "stress_cpu.workers must be set, ex.: \"1\""))
		}
		if c.Havoc.StressCPU.Load <= 0 {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "stress_cpu.load must be set, ex.: \"100\""))
		}
	}
	if c.Havoc.BlockchainRewindHead != nil {
		if c.Havoc.BlockchainRewindHead.Duration == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "havoc.blockchain_rewind_head.duration must be set, ex.: \"30s\""))
		}
		for _, bn := range c.Havoc.BlockchainRewindHead.NodesConfig {
			if bn.ExecutorPodPrefix == "" {
				errs = append(errs, errors.Wrap(errors.New(ErrFormat), "havoc.blockchain_rewind_head.nodes.executor_pod_prefix must be set, ex.: \"geth\""))
			}
			if bn.ExecutorContainerName == "" {
				errs = append(errs, errors.Wrap(errors.New(ErrFormat), "havoc.blockchain_rewind_head.nodes.executor_container_name must be set, ex.: \"geth-network\""))
			}
			if bn.NodeInternalHTTPURL == "" {
				errs = append(errs, errors.Wrap(errors.New(ErrFormat), "havoc.blockchain_rewind_head.nodes.node_internal_http_url must be set, ex.: \"geth-1337:8544\""))
			}
			if len(bn.Blocks) == 0 {
				errs = append(errs, errors.Wrap(errors.New(ErrFormat), "havoc.blockchain_rewind_head.nodes.blocks must be set, ex.: \"10\""))
			}
		}
	}
	if c.Havoc.Monkey != nil {
		if c.Havoc.Monkey.Mode == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "monkey.mode must be either \"seq\" or \"rand\""))
		}
		if c.Havoc.Monkey.Duration == "" {
			errs = append(errs, errors.Wrap(errors.New(ErrFormat), "monkey.duration must be in Go duration format, 1d2h3m0s"))
		}
	}
	return errs
}

type Failure struct {
	Duration        string   `toml:"duration"`
	GroupPercentage []string `toml:"group_percentage"`
	GroupFixed      []string `toml:"group_fixed"`
}

type Latency struct {
	Duration        string   `toml:"duration"`
	Latency         string   `toml:"latency"`
	GroupPercentage []string `toml:"group_percentage"`
	GroupFixed      []string `toml:"group_fixed"`
}

type NetworkPartition struct {
	Duration        string   `toml:"duration"`
	Label           string   `toml:"label"`
	GroupPercentage []string `toml:"group_percentage"`
	GroupFixed      []string `toml:"group_fixed"`
}

type StressMemory struct {
	Duration        string   `toml:"duration"`
	Workers         int      `toml:"workers"`
	Memory          string   `toml:"memory"`
	GroupPercentage []string `toml:"group_percentage"`
	GroupFixed      []string `toml:"group_fixed"`
}

type StressCPU struct {
	Duration        string   `toml:"duration"`
	Workers         int      `toml:"workers"`
	Load            int      `toml:"load"`
	GroupPercentage []string `toml:"group_percentage"`
	GroupFixed      []string `toml:"group_fixed"`
}

type ExternalTargets struct {
	Duration string   `toml:"duration"`
	URLs     []string `toml:"urls"`
}

type BlockchainRewindHead struct {
	Duration    string                  `toml:"duration"`
	NodesConfig []*BlockchainNodeConfig `toml:"nodes"`
}

type BlockchainNodeConfig struct {
	ExecutorPodPrefix     string  `toml:"executor_pod_prefix"`
	ExecutorContainerName string  `toml:"executor_container_name"`
	NodeInternalHTTPURL   string  `toml:"node_internal_http_url"`
	Blocks                []int64 `toml:"blocks"`
}

type OpenAPI struct {
	Mapping         map[string]*OpenApiSpecInfo `toml:"mapping"`
	Duration        string                      `toml:"duration"`
	GroupPercentage []string                    `toml:"group_percentage"`
	GroupFixed      []string                    `toml:"group_fixed"`
}

type OpenApiSpecInfo struct {
	SpecToPortMappings []*SpecToPort `toml:"spec_to_port"`
}

type SpecToPort struct {
	Port int64  `toml:"port"`
	Path string `toml:"path"`
}

type Monkey struct {
	Duration string `toml:"duration"`
	Cooldown string `toml:"cooldown"`
	Mode     string `toml:"mode"`
}

type Grafana struct {
	URL           string   `toml:"grafana_url"`
	Token         string   `toml:"grafana_token"`
	DashboardUIDs []string `toml:"dashboard_uids"`
}

func ReadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	dumpConfig(cfg)
	if path == "" {
		L.Info().Msg("No config specified, using default configuration")
	} else {
		L.Debug().
			Str("Path", path).
			Msg("Reading config from path")
		d, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, ErrReadSethConfig)
		}
		err = toml.Unmarshal(d, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, ErrUnmarshalSethConfig)
		}
	}
	L.Debug().
		Interface("Config", cfg).
		Msg("Configuration loaded")
	cfg.Havoc.Grafana.URL = os.Getenv("GRAFANA_URL")
	cfg.Havoc.Grafana.Token = os.Getenv("GRAFANA_TOKEN")
	return cfg, nil
}

// nolint
func sliceContains(target string, array []string) bool {
	for _, element := range array {
		if element == target {
			return true
		}
	}
	return false
}

func sliceContainsSubString(target string, array []string) bool {
	for _, element := range array {
		if strings.Contains(target, element) {
			return true
		}
	}
	return false
}
