package contracts

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"time"
)

const (
	DAGAwaitTimeout = 180 * time.Second
)

// AwaitMining awaits first block after DAG generation on multi-node geth networks
func AwaitMining(c client.BlockchainClient, extraData []byte) error {
	key := "next_block"
	cf := NewNextBlockConfirmer(extraData)
	c.AddHeaderEventSubscription(key, cf)
	if err := c.WaitForEvents(); err != nil {
		return err
	}
	return nil
}

// NextBlockConfirmer await for the next block
type NextBlockConfirmer struct {
	doneChan  chan struct{}
	done      bool
	extraData []byte
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewNextBlockConfirmer generic next block confirmer
func NewNextBlockConfirmer(extraData []byte) *NextBlockConfirmer {
	ctx, cancel := context.WithTimeout(context.Background(), DAGAwaitTimeout)
	return &NextBlockConfirmer{
		done:      false,
		doneChan:  make(chan struct{}),
		extraData: extraData,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (f *NextBlockConfirmer) ReceiveBlock(b *types.Block) error {
	if f.done {
		return nil
	}
	// in case of a multi-node setup we may need a block from the node on which reorg will be performed
	if f.extraData != nil {
		log.Debug().Str("ExtraData", string(b.Extra())).Msg("Block received")
		if !bytes.Contains(b.Extra(), []byte("tx")) {
			return nil
		}
	}
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
