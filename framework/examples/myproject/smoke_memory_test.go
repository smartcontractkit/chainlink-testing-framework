package examples

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type MemoryStats struct {
	HeapInUseMB   float64
	HeapAllocMB   float64
	SysMB         float64
	Goroutines    float64
	LastGCTime    float64
	GCDurationSum float64
	GCCount       float64
	AvgGCTimeMS   float64
}

type CfgMemory struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestMemoryMetrics(t *testing.T) {
	in, err := framework.Load[CfgMemory](t)
	require.NoError(t, err)

	bcOut, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	nsOut, err := ns.NewSharedDBNodeSet(in.NodeSet, bcOut)
	require.NoError(t, err)

	t.Run("memory_metrics_monitoring", func(t *testing.T) {
		// add some OCR test here
		_ = bcOut
		_ = nsOut
		time.Sleep(5 * time.Second)

		metrics, err := fetchMetrics(fmt.Sprintf("%s/metrics", nsOut.CLNodes[1].Node.HostURL))
		require.NoError(t, err)
		stats, err := parseMetrics(metrics)
		require.NoError(t, err)

		framework.L.Info().
			Float64("heap_in_use_mb", stats.HeapInUseMB).
			Float64("heap_alloc_mb", stats.HeapAllocMB).
			Float64("system_mb", stats.SysMB).
			Float64("goroutines", stats.Goroutines).
			Float64("avg_gc_time_ms", stats.AvgGCTimeMS).
			Msg("Checking memory metrics")

		// memory limits
		require.LessOrEqual(t, stats.SysMB, 1000.0)
		require.LessOrEqual(t, stats.HeapInUseMB, 300.0)
		require.LessOrEqual(t, stats.HeapAllocMB, 300.0)
		// goroutines
		require.LessOrEqual(t, stats.Goroutines, 1000.0)
		// gc pauses
		require.Less(t, stats.AvgGCTimeMS, 10.0)
	})
}

func fetchMetrics(url string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "text/plain").
		Get(url)

	if err != nil {
		return "", fmt.Errorf("failed to fetch metrics: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return resp.String(), nil
}

func parseMetrics(metrics string) (MemoryStats, error) {
	var stats MemoryStats
	scanner := bufio.NewScanner(strings.NewReader(metrics))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.Split(line, " ")) < 2 {
			continue
		}

		parts := strings.Split(line, " ")
		metricName := strings.Split(parts[0], "{")[0]
		value, _ := strconv.ParseFloat(parts[1], 64)

		switch metricName {
		case "go_memstats_heap_inuse_bytes":
			stats.HeapInUseMB = value / 1024 / 1024
		case "go_memstats_alloc_bytes":
			stats.HeapAllocMB = value / 1024 / 1024
		case "go_memstats_sys_bytes":
			stats.SysMB = value / 1024 / 1024
		case "go_goroutines":
			stats.Goroutines = value
		case "go_gc_duration_seconds_sum":
			stats.GCDurationSum = value
		case "go_gc_duration_seconds_count":
			stats.GCCount = value
		}
	}

	if stats.GCCount > 0 {
		stats.AvgGCTimeMS = (stats.GCDurationSum / stats.GCCount) * 1000
	}

	return stats, scanner.Err()
}
