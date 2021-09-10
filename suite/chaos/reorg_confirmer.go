package chaos

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
	"sync"
	"time"
)

// ReorgConfirmer reorg stats collecting struct
type ReorgConfirmer struct {
	env                      environment.Environment
	c                        client.BlockchainClient
	reorgDepth               int
	blockConsensusThreshold  int
	numberOfNodes            int
	currentDivergenceDepth   int
	currentBlockConsensus    int
	awaitingNetworkConsensus bool
	blockHashes              map[int64][]common.Hash
	chaosExperimentName      string
	mutex                    sync.Mutex

	ctx      context.Context
	cancel   context.CancelFunc
	doneChan chan struct{}
	done     bool
}

// NewReorgConfirmer creates a type that can create reorg chaos and confirm reorg has happened
func NewReorgConfirmer(
	c client.BlockchainClient,
	env environment.Environment,
	reorgDepth int,
	blockConsensusThreshold int,
	timeout time.Duration,
) (*ReorgConfirmer, error) {
	if len(c.GetClients()) == 1 {
		return nil, errors.New("Only one node within the blockchain client detected, cannot reorg")
	}
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	rc := &ReorgConfirmer{
		env:                     env,
		c:                       c,
		reorgDepth:              reorgDepth,
		blockConsensusThreshold: blockConsensusThreshold,
		numberOfNodes:           len(c.GetClients()),
		blockHashes:             map[int64][]common.Hash{},
		mutex:                   sync.Mutex{},
		ctx:                     ctx,
		cancel:                  ctxCancel,
	}
	return rc, rc.forkNetwork()
}

func (rc *ReorgConfirmer) ReceiveBlock(header *types.Block) error {
	if header == nil || rc.done {
		return nil
	}
	if rc.awaitingNetworkConsensus {
		if rc.hasNetworkFormedConsensus(header) {
			rc.doneChan <- struct{}{}
		}
	} else {
		if rc.hasNetworkMetReorgDepth(header) {
			rc.awaitingNetworkConsensus = true
			if err := rc.joinNetwork(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rc *ReorgConfirmer) Wait() error {
	for {
		select {
		case <-rc.doneChan:
			rc.cancel()
			rc.done = true
			return nil
		case <-rc.ctx.Done():
			return errors.New("timeout waiting for reorg to complete")
		}
	}
}

func (rc *ReorgConfirmer) forkNetwork() error {
	expName, err := rc.env.ApplyChaos(&experiments.NetworkPartition{
		FromMode:       "one",
		FromLabelKey:   "app",
		FromLabelValue: "ethereum-geth-tx",
		ToMode:         "all",
		ToLabelKey:     "app",
		ToLabelValue:   "ethereum-geth-miner",
	})
	rc.chaosExperimentName = expName
	return err
}

func (rc *ReorgConfirmer) joinNetwork() error {
	return rc.env.StopChaos(rc.chaosExperimentName)
}

func (rc *ReorgConfirmer) hasNetworkFormedConsensus(header *types.Block) bool {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	blockNumber := header.Number().Int64()
	rc.appendBlockHeader(header)

	// If we've received the same block number from all nodes, check hashes to ensure they've reformed consensus
	if len(rc.blockHashes[blockNumber]) == rc.numberOfNodes {
		firstBlockHash := rc.blockHashes[blockNumber][0]
		for _, blockHash := range rc.blockHashes[blockNumber][1:] {
			if blockHash.String() != firstBlockHash.String() {
				log.Info().
					Int64("Blocknumber", blockNumber).
					Msg("Reorg detected for block, awaiting network rejoin")
				return false
			}
		}
		rc.currentBlockConsensus++
	}

	if rc.currentBlockConsensus >= rc.blockConsensusThreshold {
		log.Info().
			Msg("Network has reformed consensus and joined, reorg complete")
		return true
	}
	return false
}

func (rc *ReorgConfirmer) hasNetworkMetReorgDepth(header *types.Block) bool {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	blockNumber := header.Number().Int64()
	rc.appendBlockHeader(header)

	// If we've received the same block number from all nodes, check hashes to verify if a reorg is taking place
	if len(rc.blockHashes[blockNumber]) == rc.numberOfNodes {
		firstBlockHash := rc.blockHashes[blockNumber][0]
		for _, blockHash := range rc.blockHashes[blockNumber][1:] {
			if blockHash.String() == firstBlockHash.String() {
				log.Info().
					Int64("Blocknumber", blockNumber).
					Msg("Reorg not detected for block")
				rc.currentDivergenceDepth = 0
				return false
			}
		}
		rc.currentDivergenceDepth++
		log.Info().
			Int64("Blocknumber", blockNumber).
			Int("Current divergence depth", rc.currentDivergenceDepth).
			Msg("Block reorg detected")
	}

	if rc.currentDivergenceDepth >= rc.reorgDepth {
		log.Info().Int("Reorg Depth", rc.reorgDepth).Msg("Desired reorg depth met, joining the network")
		return true
	}
	return false
}

func (rc *ReorgConfirmer) appendBlockHeader(header *types.Block) {
	blockNumber := header.Number().Int64()
	if _, ok := rc.blockHashes[blockNumber]; !ok {
		rc.blockHashes[blockNumber] = []common.Hash{}
	}
	rc.blockHashes[blockNumber] = append(rc.blockHashes[blockNumber], header.Hash())
}
