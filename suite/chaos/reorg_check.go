package chaos

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/environment"
	"golang.org/x/sync/errgroup"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	MinerIDTemplate    = "miner-%d"
	ReorgCheckAttempts = 15
	ReorgCheckInterval = 10 * time.Second
)

// ReorgChecker collects data to verify reorged blocks
type ReorgChecker struct {
	// main client that is connected to the Chainlink node on which we perform reorg
	Client client.BlockchainClient
	Env    environment.Environment
	// Other miners clients
	AltClients []client.BlockchainClient
	FromBlock  uint64
	HeadersMu  *sync.Mutex
	Headers    map[uint64]map[string]*types.Header
	TXHashes   []common.Hash
}

// NewReorgChecker created new re-org checker connected to all miner nodes
func NewReorgChecker(connectedClient client.BlockchainClient, env environment.Environment, cfg *config.Config, altMinersPort uint16) (*ReorgChecker, error) {
	sd, err := env.GetAllServiceDetails(altMinersPort)
	if err != nil {
		return nil, err
	}
	altClients := make([]client.BlockchainClient, 0)
	for idx, d := range sd {
		log.Debug().Str("Remote", d.RemoteURL.String()).Str("Local", d.LocalURL.String()).Msg("Miners RPCs")
		n, err := client.NewNetworkFromConfig(cfg)
		if err != nil {
			return nil, err
		}
		id := fmt.Sprintf(MinerIDTemplate, idx)
		n.SetURL(strings.Replace(d.LocalURL.String(), "http", "ws", -1))
		n.SetID(id)
		connectedClient.SetID("tx")
		cl, err := client.NewBlockchainClient(n)
		if err != nil {
			return nil, err
		}
		cl.SetID(id)
		altClients = append(altClients, cl)
	}
	fromBlock, err := connectedClient.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	m := &ReorgChecker{
		Env:        env,
		Client:     connectedClient,
		AltClients: altClients,
		FromBlock:  fromBlock,
		HeadersMu:  &sync.Mutex{},
		Headers:    make(map[uint64]map[string]*types.Header),
		TXHashes:   make([]common.Hash, 0),
	}
	return m, nil
}

// Fork separates reorg node from other miners
func (rc *ReorgChecker) Fork(dur time.Duration) error {
	log.Info().Msg("Forking network")
	exp, err := rc.Env.ApplyChaos(&experiments.NetworkPartition{
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

// getReorgedHeaders gets current version of headers for all nodes
func (rc *ReorgChecker) getReorgedHeaders() error {
	allClients := append([]client.BlockchainClient{rc.Client}, rc.AltClients...)
	g := errgroup.Group{}
	for _, ac := range allClients {
		ac := ac
		latestBn, err := ac.BlockNumber(context.Background())
		if err != nil {
			return err
		}
		for i := rc.FromBlock; i < latestBn; i++ {
			i := i
			g.Go(func() error {
				h, err := ac.(*client.EthereumClient).Client.HeaderByNumber(context.Background(), big.NewInt(int64(i)))
				if err != nil {
					return err
				}
				hn := h.Number.Uint64()
				rc.HeadersMu.Lock()
				if rc.Headers[hn] == nil {
					rc.Headers[hn] = make(map[string]*types.Header)
				}
				rc.Headers[hn][ac.GetID()] = h
				rc.HeadersMu.Unlock()
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// Verify verifies that target node has reorged blocks
func (rc *ReorgChecker) Verify() error {
	log.Info().Msg("Verifying reorg")
	err := retry.Do(func() error {
		log.Debug().Msg("Checking blocks for reorg")
		err := rc.getReorgedHeaders()
		if err != nil {
			return err
		}
		err = rc.printBlockVersions()
		if err != nil {
			return err
		}
		if err := rc.FindReorgedBlocks(); err != nil {
			return err
		}
		rc.Headers = make(map[uint64]map[string]*types.Header)
		return nil
	}, retry.DelayType(retry.FixedDelay), retry.Attempts(ReorgCheckAttempts), retry.Delay(ReorgCheckInterval))
	if err != nil {
		return err
	}
	return nil
}

// FindReorgedBlocks finds reorged blocks, returns error if no blocks found
func (rc *ReorgChecker) FindReorgedBlocks() error {
	seenHeaders := rc.Client.(*client.EthereumClient).Headers
	reorged := 0
	for _, seenHeader := range seenHeaders {
		for _, headerNow := range rc.Headers {
			if rc.isBlockReorged(seenHeader, headerNow["tx"]) {
				log.Info().
					Uint64("Number", seenHeader.Number.Uint64()).
					Str("Hash before", seenHeader.Hash().Hex()).
					Str("Hash after", headerNow["tx"].Hash().Hex()).
					Msg("Block was reorged")
				reorged += 1
			}
		}
	}
	log.Info().Int("Blocks", reorged).Msg("Total blocks reorged")
	if reorged == 0 {
		return errors.New("reorg failed, no blocks were replaced")
	}
	return nil
}

func (rc *ReorgChecker) printBlockVersions() error {
	bnSorted := make([]uint64, 0)
	for bn := range rc.Headers {
		bnSorted = append(bnSorted, bn)
	}
	sort.Slice(bnSorted, func(i, j int) bool { return bnSorted[i] < bnSorted[j] })
	for _, bn := range bnSorted {
		log.Debug().Uint64("Number", bn).Msg("Blocks info")
		for nodeID, b := range rc.Headers[bn] {
			log.Debug().Str("Node", nodeID).Str("Hash", b.Hash().Hex()).Msg("Block hash")
		}
	}
	return nil
}

func (rc *ReorgChecker) isBlockReorged(h *types.Header, h2 *types.Header) bool {
	if h == nil || h2 == nil {
		return false
	}
	return h.Number.Uint64() == h2.Number.Uint64() && h.Hash().Hex() != h2.Hash().Hex()
}
