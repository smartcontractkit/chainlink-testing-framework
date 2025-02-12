// File: sentinel_coordinator.go
package sentinel

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
)

type SentinelCoordinator struct {
	Sentinel *Sentinel
	wg       *sync.WaitGroup
	Ctx      context.Context
	cancel   context.CancelFunc
	log      zerolog.Logger
}

func NewSentinelCoordinator(log zerolog.Logger) *SentinelCoordinator {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	s := NewSentinel(SentinelConfig{})

	return &SentinelCoordinator{
		Sentinel: s,
		wg:       wg,
		Ctx:      ctx,
		cancel:   cancel,
		log:      log,
	}
}
