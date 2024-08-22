package havoc

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	ChaosTypeBlockchainSetHead = "blockchain_rewind_head"
	ChaosTypeFailure           = "failure"
	ChaosTypeGroupFailure      = "group-failure"
	ChaosTypeLatency           = "latency"
	ChaosTypeGroupLatency      = "group-latency"
	ChaosTypeStressMemory      = "memory"
	ChaosTypeStressGroupMemory = "group-memory"
	ChaosTypeStressCPU         = "cpu"
	ChaosTypeStressGroupCPU    = "group-cpu"
	ChaosTypePartitionExternal = "external"
	ChaosTypePartitionGroup    = "group-partition"
	ChaosTypeHTTP              = "http"
)

var (
	ExperimentTypesToCRDNames = map[string]string{
		"PodChaos":     "podchaos.chaos-mesh.org",
		"StressChaos":  "stresschaos.chaos-mesh.org",
		"NetworkChaos": "networkchaos.chaos-mesh.org",
		"HTTPChaos":    "httpchaos.chaos-mesh.org",
	}
)

var L zerolog.Logger

func SetGlobalLogger(l zerolog.Logger) {
	L = l.With().Str("Component", "havoc").Logger()
}

func InitDefaultLogging() {
	lvl, err := zerolog.ParseLevel(os.Getenv("HAVOC_LOG_LEVEL"))
	if err != nil {
		panic(err)
	}
	if lvl.String() == "" {
		lvl = zerolog.InfoLevel
	}
	L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl)
}

type ChaosSpecs struct {
	ExperimentsByType map[string]map[string]string
}

func (m *ChaosSpecs) Dump(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}
	L.Info().Str("Dir", dir).Msg("Writing experiments to a dir")
	for expType := range m.ExperimentsByType {
		if len(m.ExperimentsByType[expType]) == 0 {
			continue
		}
		if err := os.Mkdir(fmt.Sprintf("%s/%s", dir, expType), os.ModePerm); err != nil {
			return err
		}
		for expName, expBody := range m.ExperimentsByType[expType] {
			fname := strings.ToLower(fmt.Sprintf("%s/%s/%s-%s.yaml", dir, expType, expType, expName))
			if err := os.WriteFile(fname, []byte(expBody), os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
