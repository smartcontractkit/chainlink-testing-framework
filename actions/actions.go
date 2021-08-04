package actions

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"math/big"
	"sort"
	"strings"
	"sync"
)

// FundChainlinkNodes will fund all of the Chainlink nodes with a given amount of ETH in wei
func FundChainlinkNodes(
	nodes []client.Chainlink,
	blockchain client.BlockchainClient,
	fromWallet client.BlockchainWallet,
	nativeAmount,
	linkAmount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = blockchain.Fund(fromWallet, toAddress, nativeAmount, linkAmount)
		if err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddresses(nodes []client.Chainlink) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// GetRoundStartTimesAcrossNodes gets round start times across nodes
func GetRoundStartTimesAcrossNodes(
	contractAddr string,
	jobs []contracts.JobByInstance,
	nodesByHostPort map[string]client.Chainlink,
	mu *sync.Mutex,
	contractsStartTimes map[string]int64,
) error {
	startTimesAcrossNodes := make([]int64, 0)
	for _, j := range jobs {
		// get node for a job
		node := nodesByHostPort[j.Instance]
		runs, err := node.ReadRunsByJob(j.ID)
		if err != nil {
			return err
		}
		runsStartTimes := make([]int64, 0)
		for _, r := range runs.Data {
			runsStartTimes = append(runsStartTimes, r.Attributes.CreatedAt.UnixNano())
		}
		sort.SliceStable(runsStartTimes, func(i, j int) bool {
			return runsStartTimes[i] > runsStartTimes[j]
		})
		// debug
		stHr := make([]int64, 0)
		for _, st := range runsStartTimes {
			stHr = append(stHr, st/1e9)
		}
		log.Debug().Interface("Start times array", stHr).Send()
		//
		lastRunStartTime := runsStartTimes[0]
		log.Debug().Int64("Last run start time", lastRunStartTime/1e9).Send()
		startTimesAcrossNodes = append(startTimesAcrossNodes, lastRunStartTime)
	}
	mu.Lock()
	defer mu.Unlock()
	// earliest start across nodes for contract
	sort.SliceStable(startTimesAcrossNodes, func(i, j int) bool {
		return startTimesAcrossNodes[i] < startTimesAcrossNodes[j]
	})
	contractsStartTimes[contractAddr] = startTimesAcrossNodes[0]
	return nil
}

// GetRoundCompleteTimestamps gets round finish timestamp from last event seen on-chain
func GetRoundCompleteTimestamps(
	fi contracts.FluxAggregator,
	roundID int,
	submissions int,
	submissionVal *big.Int,
	mu *sync.Mutex,
	endTimes map[string]int64,
	client client.BlockchainClient,
) error {
	events, err := fi.FilterRoundSubmissions(context.Background(), submissionVal, roundID)
	if err != nil {
		return err
	}
	if len(events) == submissions {
		lastEvent := events[len(events)-1]
		hTime, err := client.HeaderTimestampByNumber(context.Background(), big.NewInt(int64(lastEvent.BlockNumber)))
		if err != nil {
			return err
		}
		log.Debug().
			Str("Contract", fi.Address()).
			Uint64("Header timestamp", hTime).
			Msg("All submissions found")
		mu.Lock()
		defer mu.Unlock()
		// nanoseconds
		endTimes[fi.Address()] = int64(hTime) * 1e9
		return nil
	}
	return fmt.Errorf("not all submissions found for contract: %s", fi.Address())
}

// EncodeOnChainExternalJobID encodes external job uuid to on-chain representation
func EncodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}

// EncodeOnChainVRFProvingKey encodes uncompressed public VRF key to on-chain representation
func EncodeOnChainVRFProvingKey(vrfKey client.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var set1 bool
	var set2 bool
	// strip 0x to convert to int
	provingKey[0], set1 = new(big.Int).SetString(uncompressed[2:66], 16)
	if !set1 {
		return [2]*big.Int{}, errors.New("can not convert VRF key to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(uncompressed[66:], 16)
	if !set2 {
		return [2]*big.Int{}, errors.New("can not convert VRF key to *big.Int")
	}
	return provingKey, nil
}
