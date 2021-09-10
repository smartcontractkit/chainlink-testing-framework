package chaos

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
	"sort"
	"sync"
	"time"
)

const (
	NextBlockTimeout    = 180 * time.Second
	NetworkJoinAttempts = 30
	NetworkJoinInterval = 5 * time.Second
	VerifyAttempts      = 20
	VerifyInterval      = 5 * time.Second
)

// NodeBlock block received from particular node
type NodeBlock struct {
	NodeID int
	*types.Block
}

// ReorgConfirmer reorg stats collecting struct
type ReorgConfirmer struct {
	env                 environment.Environment
	c                   client.BlockchainClient
	joinBlockNumber     uint64
	blockReceivers      []*AggregatingBlockReceiver
	blocksByNodeMu      *sync.Mutex
	reorgedBlocksByNode map[int][]*types.Block
	blocksByNode        map[uint64]map[int]NodeBlock
	blocksChan          chan NodeBlock
}

// NewReorgConfirmer creates a type that can create reorg chaos and confirm reorg has happened
func NewReorgConfirmer(c client.BlockchainClient, env environment.Environment) (*ReorgConfirmer, error) {
	rc := &ReorgConfirmer{
		env:                 env,
		c:                   c,
		blockReceivers:      make([]*AggregatingBlockReceiver, 0),
		blocksByNodeMu:      &sync.Mutex{},
		reorgedBlocksByNode: make(map[int][]*types.Block),
		blocksByNode:        make(map[uint64]map[int]NodeBlock),
		blocksChan:          make(chan NodeBlock),
	}
	rc.subscribe()
	go rc.aggregateBlocks()
	if err := rc.awaitNetworkJoined(); err != nil {
		return nil, err
	}
	return rc, nil
}

// Verify verifies that reorg has happened, unsubscribing for blocks
func (rc *ReorgConfirmer) Verify(nodeID int, depth int) error {
	defer rc.shutdown()
	err := retry.Do(func() error {
		log.Debug().Msg("Verifying reorged blocks")
		rc.blocksByNodeMu.Lock()
		defer rc.blocksByNodeMu.Unlock()
		if len(rc.reorgedBlocksByNode[nodeID]) >= depth {
			log.Debug().Msg("Reorg verified")
			return nil
		}
		return fmt.Errorf("no reorg blocks found")
	}, retry.DelayType(retry.FixedDelay), retry.Attempts(VerifyAttempts), retry.Delay(VerifyInterval))
	if err != nil {
		return err
	}
	return nil
}

func (rc *ReorgConfirmer) shutdown() {
	for _, br := range rc.blockReceivers {
		br.cancel()
	}
	rc.unsubscribe()
}

func (rc *ReorgConfirmer) Fork(dur time.Duration) error {
	log.Info().Msg("Forking network")
	exp, err := rc.env.ApplyChaos(&experiments.NetworkPartition{
		FromMode:       "one",
		FromLabelKey:   "app",
		FromLabelValue: "ethereum-geth-tx",
		ToMode:         "all",
		ToLabelKey:     "app",
		ToLabelValue:   "ethereum-geth-miner",
		Duration:       dur,
	})
	if err != nil {
		return err
	}
	log.Info().Msg("Network forked")
	time.Sleep(dur)
	log.Debug().Str("Experiment", exp).Msg("Joining network")
	return nil
}

func (rc *ReorgConfirmer) unsubscribe() {
	for idx, c := range rc.c.GetClients() {
		key := fmt.Sprintf("%s_%d", client.BlocksSubPrefix, idx)
		c.DeleteHeaderEventSubscription(key)
	}
	rc.blockReceivers = make([]*AggregatingBlockReceiver, 0)
}

func (rc *ReorgConfirmer) subscribe() {
	for idx, c := range rc.c.(*client.EthereumClients).Clients {
		key := fmt.Sprintf("%s_%d", client.BlocksSubPrefix, idx)
		br := NewAggregatingBlockReceiver(idx, rc.blocksChan)
		c.AddHeaderEventSubscription(key, br)
		rc.blockReceivers = append(rc.blockReceivers, br)
	}
}

// aggregateBlocks aggregates blocks, if hash was already seen adds to reorged blocks
func (rc *ReorgConfirmer) aggregateBlocks() {
	for b := range rc.blocksChan {
		rc.blocksByNodeMu.Lock()
		if _, ok := rc.blocksByNode[b.NumberU64()]; !ok {
			rc.blocksByNode[b.NumberU64()] = make(map[int]NodeBlock)
		}
		seenBlock, seen := rc.blocksByNode[b.NumberU64()][b.NodeID]
		if seen && seenBlock.Hash().Hex() != b.Hash().Hex() && b.NumberU64() >= rc.joinBlockNumber {
			log.Info().Int("Node", b.NodeID).
				Uint64("Number", b.NumberU64()).
				Str("Old hash", seenBlock.Hash().Hex()).
				Str("New Hash", b.Hash().Hex()).
				Msg("Block hash was updated")
			if rc.reorgedBlocksByNode[b.NodeID] == nil {
				rc.reorgedBlocksByNode[b.NodeID] = make([]*types.Block, 0)
			}
			rc.reorgedBlocksByNode[b.NodeID] = append(rc.reorgedBlocksByNode[b.NodeID], b.Block)
		}
		rc.blocksByNode[b.NumberU64()][b.NodeID] = b
		rc.blocksByNodeMu.Unlock()
	}
}

// isJoinBlock checks that we have a block with all versions from different client and block hashes are equal
func (rc *ReorgConfirmer) isJoinBlock(m map[int]NodeBlock) bool {
	values := make([]string, 0)
	for _, v := range m {
		values = append(values, v.Hash().Hex())
	}
	if len(values) < len(rc.c.GetClients()) {
		return false
	}
	for i := 1; i < len(values); i++ {
		if values[i] != values[0] {
			return false
		}
	}
	return true
}

// awaitNetworkJoined awaits common block seen by all nodes
func (rc *ReorgConfirmer) awaitNetworkJoined() error {
	err := retry.Do(func() error {
		log.Debug().Msg("Checking for a join block")
		rc.blocksByNodeMu.Lock()
		defer rc.blocksByNodeMu.Unlock()
		revBlocks := make([]uint64, 0)
		for bn := range rc.blocksByNode {
			revBlocks = append(revBlocks, bn)
		}
		sort.Slice(revBlocks, func(i, j int) bool {
			return revBlocks[i] > revBlocks[j]
		})
		for _, bn := range revBlocks {
			if rc.isJoinBlock(rc.blocksByNode[bn]) {
				rc.joinBlockNumber = bn
				log.Info().Uint64("Number", bn).Msg("Join block found")
				return nil
			}
		}
		return fmt.Errorf("network is still joining")
	}, retry.DelayType(retry.FixedDelay), retry.Attempts(NetworkJoinAttempts), retry.Delay(NetworkJoinInterval))
	if err != nil {
		return err
	}
	return nil
}

// AggregatingBlockReceiver receives next block, mark it by node id
type AggregatingBlockReceiver struct {
	id         int
	doneChan   chan struct{}
	done       bool
	ctx        context.Context
	cancel     context.CancelFunc
	blocksChan chan NodeBlock
}

// NewAggregatingBlockReceiver generic next block receiver that aggregates blocks per blockchain node
func NewAggregatingBlockReceiver(id int, blocksChan chan NodeBlock) *AggregatingBlockReceiver {
	ctx, cancel := context.WithTimeout(context.Background(), NextBlockTimeout)
	return &AggregatingBlockReceiver{
		id:         id,
		done:       false,
		doneChan:   make(chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
		blocksChan: blocksChan,
	}
}

func (f *AggregatingBlockReceiver) ReceiveBlock(b *types.Block) error {
	if b == nil {
		return nil
	}
	select {
	case f.blocksChan <- NodeBlock{NodeID: f.id, Block: b}:
	case <-f.ctx.Done():
		return nil
	}
	return nil
}

func (f *AggregatingBlockReceiver) Wait() error {
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
