package tools

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	QueryMemoryUsage          = `100 * (1 - ((avg_over_time(node_memory_MemFree_bytes[%s]) + avg_over_time(node_memory_Cached_bytes[%s]) + avg_over_time(node_memory_Buffers_bytes[%s])) / avg_over_time(node_memory_MemTotal_bytes[%s])))`
	QueryAllCPUBusyPercentage = `100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[%s])) * 100)`
)

type ResourcesSummary struct {
	MemoryUsage   float64
	CPUPercentage float64
}

type PromChecker struct {
	API v1.API
	Url string
}

func NewPrometheusClient(url string) (*PromChecker, error) {
	client, err := api.NewClient(api.Config{
		Address: url,
	})
	if err != nil {
		return nil, err
	}
	return &PromChecker{
		API: v1.NewAPI(client),
		Url: url,
	}, nil
}

func (p *PromChecker) printWarns(warns v1.Warnings) {
	if len(warns) > 0 {
		log.Info().Interface("Warnings", warns).Msg("Warnings found when performing prometheus query")
	}
}

func (p *PromChecker) validateNotEmptyVec(q string, val model.Value) bool {
	if len(val.(model.Vector)) == 0 {
		log.Warn().Str("query", q).Msg("empty response for prometheus query")
		return false
	}
	return true
}

// CPUBusyPercentage host CPU busy percentage
func (p *PromChecker) CPUBusyPercentage() (float64, error) {
	q := fmt.Sprintf(QueryAllCPUBusyPercentage, "2m")
	val, warns, err := p.API.Query(context.Background(), q, time.Now())
	if err != nil {
		return 0, err
	}
	p.printWarns(warns)
	if !p.validateNotEmptyVec(q, val) {
		return 0, nil
	}
	scalarVal := val.(model.Vector)[0].Value
	return float64(scalarVal), nil
}

// MemoryUsage total memory used by interval
func (p *PromChecker) MemoryUsage() (float64, error) {
	q := fmt.Sprintf(QueryMemoryUsage, "2m", "2m", "2m", "2m")
	val, warns, err := p.API.Query(context.Background(), q, time.Now())
	if err != nil {
		return 0, err
	}
	p.printWarns(warns)
	if !p.validateNotEmptyVec(q, val) {
		return 0, nil
	}
	scalarVal := val.(model.Vector)[0].Value
	return float64(scalarVal), nil
}

func (p *PromChecker) ResourcesSummary() (float64, float64, error) {
	cpu, err := p.CPUBusyPercentage()
	if err != nil {
		return 0, 0, err
	}
	mem, err := p.MemoryUsage()
	if err != nil {
		return 0, 0, err
	}
	return cpu, mem, nil
}
