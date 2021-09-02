package contracts

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"time"
)

const (
	DAGAwaitTimeout = 120 * time.Second
)

// AwaitMining awaits first block after DAG generation on multi-node geth networks
func AwaitMining(c client.BlockchainClient) error {
	log.Info().Msg("Awaiting first block to be mined")
	key := "next_block"
	cf := NewNextBlockConfirmer()
	c.AddHeaderEventSubscription(key, cf)
	if err := c.WaitForEvents(); err != nil {
		return err
	}
	return nil
}

// NextBlockConfirmer await for the next block
type NextBlockConfirmer struct {
	doneChan chan struct{}
	done     bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewNextBlockConfirmer generic next block confirmer
func NewNextBlockConfirmer() *NextBlockConfirmer {
	ctx, cancel := context.WithTimeout(context.Background(), DAGAwaitTimeout)
	return &NextBlockConfirmer{
		done:     false,
		doneChan: make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (f *NextBlockConfirmer) ReceiveBlock(_ *types.Block) error {
	if f.done {
		return nil
	}
	log.Info().Msg("First block received")
	f.done = true
	f.doneChan <- struct{}{}
	return nil
}

func (f *NextBlockConfirmer) Wait() error {
	for {
		select {
		case <-f.doneChan:
			f.cancel()
			return nil
		case <-f.ctx.Done():
			return errors.New("timeout waiting for the next block to confirm")
		}
	}
}
