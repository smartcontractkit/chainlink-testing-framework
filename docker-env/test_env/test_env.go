package test_env

import (
	"sync"

	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-testing-framework/docker-env"
	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/types/envcommon"
	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/types/node"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
)

type CLClusterTestEnv struct {
	cfg        *TestEnvConfig
	Network    *tc.DockerNetwork
	LogWatch   *logwatch.LogWatch
	CLNodes    []*ClNode
	Geth       *Geth
	MockServer *MockServer
}

func NewTestEnv() (*CLClusterTestEnv, error) {
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	return &CLClusterTestEnv{
		Network: network,
		Geth: NewGeth(envcommon.EnvComponentOpts{
			Networks: networks,
		}),
		MockServer: NewMockServer(envcommon.EnvComponentOpts{
			Networks: networks,
		}),
	}, nil
}

func NewTestEnvFromCfg(cfg *TestEnvConfig) (*CLClusterTestEnv, error) {
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	log.Info().Interface("Cfg", cfg).Send()
	return &CLClusterTestEnv{
		cfg:     cfg,
		Network: network,
		Geth: NewGeth(envcommon.EnvComponentOpts{
			ReuseContainerName: cfg.Geth.ContainerName,
			Networks:           networks,
		}),
		MockServer: NewMockServer(envcommon.EnvComponentOpts{
			ReuseContainerName: cfg.MockServer.ContainerName,
			Networks:           networks,
		}),
	}, nil
}

func (m *CLClusterTestEnv) StartGeth() error {
	return m.Geth.StartContainer(m.LogWatch)
}

func (m *CLClusterTestEnv) StartMockServer() error {
	return m.MockServer.StartContainer(m.LogWatch)
}

// StartClNodes start one bootstrap node and {count} OCR nodes
func (m *CLClusterTestEnv) StartClNodes(nodeConfigOpts node.ConfigOpts, count int) error {
	var wg sync.WaitGroup
	var errs = []error{}
	var mu sync.Mutex

	// Start nodes
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			dbContainerName := fmt.Sprintf("cl-db-%s", uuid.NewString())
			opts := envcommon.EnvComponentOpts{
				Networks: []string{m.Network.Name},
			}
			if m.cfg != nil {
				opts.ReuseContainerName = m.cfg.Nodes[i].NodeContainerName
				dbContainerName = m.cfg.Nodes[i].DbContainerName
			}
			n := NewClNode(opts, nodeConfigOpts, dbContainerName)
			err := n.StartContainer(m.LogWatch)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			} else {
				mu.Lock()
				m.CLNodes = append(m.CLNodes, n)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return multierr.Combine(errs...)
	}
	return nil
}

func (m *CLClusterTestEnv) GetDefaultNodeConfigOpts() node.ConfigOpts {
	return node.ConfigOpts{
		EVM: struct {
			HttpUrl string
			WsUrl   string
		}{
			HttpUrl: m.Geth.InternalHttpUrl,
			WsUrl:   m.Geth.InternalWsUrl,
		},
	}
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func (m *CLClusterTestEnv) ChainlinkNodeAddresses() ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, n := range m.CLNodes {
		primaryAddress, err := n.ChainlinkNodeAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, primaryAddress)
	}
	return addresses, nil
}

// FundChainlinkNodes will fund all the provided Chainlink nodes with a set amount of native currency
func (m *CLClusterTestEnv) FundChainlinkNodes(amount *big.Float) error {
	for _, cl := range m.CLNodes {
		if err := cl.Fund(m.Geth, amount); err != nil {
			return err
		}
	}
	return m.Geth.EthClient.WaitForEvents()
}

func (m *CLClusterTestEnv) GetNodeCSAKeys() ([]string, error) {
	var keys []string
	for _, n := range m.CLNodes {
		csaKeys, err := n.GetNodeCSAKeys()
		if err != nil {
			return nil, err
		}
		keys = append(keys, csaKeys.Data[0].ID)
	}
	return keys, nil
}

func (m *CLClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by default so ryuk will remove all
	// the containers and the network
	return nil
}
