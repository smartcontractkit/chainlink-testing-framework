package test_env

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"

	"math/big"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/types/envcommon"
	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/types/node"
	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/utils/templates"
)

type ClNode struct {
	envcommon.EnvComponent
	API            *client.ChainlinkClient
	NodeConfigOpts node.ConfigOpts
	DbC            *tc.Container
	DbCName        string
	DbOpts         envcommon.PgOpts
}

func NewClNode(compOpts envcommon.EnvComponentOpts, opts node.ConfigOpts, dbContainerName string) *ClNode {
	return &ClNode{
		EnvComponent:   envcommon.NewEnvComponent("cl-node", compOpts),
		DbCName:        dbContainerName,
		NodeConfigOpts: opts,
		DbOpts:         envcommon.NewDefaultPgOpts("cl-node", compOpts.Networks),
	}
}

func (m *ClNode) AddBootstrapJob(verifierAddr common.Address, fromBlock uint64, chainId int64,
	feedId [32]byte) (*client.Job, error) {
	spec := utils.BuildBootstrapSpec(verifierAddr, chainId, fromBlock, feedId)
	return m.API.MustCreateJob(spec)
}

func (m *ClNode) GetContainerName() string {
	name, err := m.EnvComponent.Container.Name(context.Background())
	if err != nil {
		return ""
	}
	return strings.Replace(name, "/", "", -1)
}

func (m *ClNode) GetPeerUrl() (string, error) {
	p2pKeys, err := m.API.MustReadP2PKeys()
	if err != nil {
		return "", err
	}
	p2pId := p2pKeys.Data[0].Attributes.PeerID

	return fmt.Sprintf("%s@%s:%d", p2pId, m.GetContainerName(), 6690), nil
}

func (m *ClNode) GetNodeCSAKeys() (*client.CSAKeys, error) {
	csaKeys, _, err := m.API.ReadCSAKeys()
	if err != nil {
		return nil, err
	}
	return csaKeys, err
}

func (m *ClNode) ChainlinkNodeAddress() (common.Address, error) {
	addr, err := m.API.PrimaryEthAddress()
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(addr), nil
}

func (m *ClNode) Fund(g *Geth, amount *big.Float) error {
	toAddress, err := m.API.PrimaryEthAddress()
	if err != nil {
		return err
	}
	gasEstimates, err := g.EthClient.EstimateGas(ethereum.CallMsg{})
	if err != nil {
		return err
	}
	return g.EthClient.Fund(toAddress, amount, gasEstimates)
}

func (m *ClNode) StartContainer(lw *logwatch.LogWatch) error {
	pgReq := envcommon.GetPgContainerRequest(m.DbCName, m.DbOpts)
	pgC, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *pgReq,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return err
	}

	nodeSecrets, err := templates.ExecuteNodeSecretsTemplate(pgReq.Name, "5432")
	if err != nil {
		return err
	}
	clReq, err := m.getContainerRequest(nodeSecrets)
	if err != nil {
		return err
	}
	clC, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *clReq,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return errors.Wrapf(err, "could not start chainlink node container")
	}
	if lw != nil {
		if err := lw.ConnectContainer(context.Background(), clC, "chainlink", true); err != nil {
			return err
		}
	}
	ctName, err := clC.Name(context.Background())
	if err != nil {
		return err
	}
	ctName = strings.Replace(ctName, "/", "", -1)
	clEndpoint, err := clC.Endpoint(context.Background(), "http")
	if err != nil {
		return err
	}

	log.Info().Str("containerName", ctName).
		Str("clEndpoint", clEndpoint).
		Msg("Started Chainlink Node container")

	clClient, err := client.NewChainlinkClient(&client.ChainlinkConfig{
		URL:      clEndpoint,
		Email:    "local@local.com",
		Password: "localdevpassword",
	})
	if err != nil {
		return errors.Wrapf(err, "could not connect Node HTTP Client")
	}

	m.EnvComponent.Container = clC
	m.DbC = &pgC
	m.API = clClient

	return nil
}

func (m *ClNode) getContainerRequest(secrets string) (
	*tc.ContainerRequest, error) {
	configFile, err := os.CreateTemp("", "node_config")
	if err != nil {
		return nil, err
	}
	config, err := node.ExecuteNodeConfigTemplate(m.NodeConfigOpts)
	if err != nil {
		return nil, err
	}
	_, err = configFile.WriteString(config)
	if err != nil {
		return nil, err
	}

	secretsFile, err := os.CreateTemp("", "node_secrets")
	if err != nil {
		return nil, err
	}
	_, err = secretsFile.WriteString(secrets)
	if err != nil {
		return nil, err
	}

	adminCreds := "local@local.com\nlocaldevpassword"
	adminCredsFile, err := os.CreateTemp("", "admin_creds")
	if err != nil {
		return nil, err
	}
	_, err = adminCredsFile.WriteString(adminCreds)
	if err != nil {
		return nil, err
	}

	apiCreds := "local@local.com\nlocaldevpassword"
	apiCredsFile, err := os.CreateTemp("", "api_creds")
	if err != nil {
		return nil, err
	}
	_, err = apiCredsFile.WriteString(apiCreds)
	if err != nil {
		return nil, err
	}

	configPath := "/home/cl-node-config.toml"
	secPath := "/home/cl-node-secrets.toml"
	adminCrePath := "/home/admin-credentials.txt"
	apiCrePath := "/home/api-credentials.txt"

	image, ok := os.LookupEnv("CHAINLINK_IMAGE")
	if !ok {
		return nil, errors.New("CHAINLINK_IMAGE env must be set")
	}
	tag, ok := os.LookupEnv("CHAINLINK_VERSION")
	if !ok {
		return nil, errors.New("CHAINLINK_VERSION env must be set")
	}

	return &tc.ContainerRequest{
		Name:         m.EnvComponent.ContainerName,
		Image:        fmt.Sprintf("%s:%s", image, tag),
		ExposedPorts: []string{"6688/tcp"},
		Entrypoint: []string{"chainlink",
			"-c", configPath,
			"-s", secPath,
			"node", "start", "-d",
			"-p", adminCrePath,
			"-a", apiCrePath,
		},
		Networks: m.Networks,
		WaitingFor: tcwait.ForHTTP("/health").
			WithPort("6688/tcp").
			WithStartupTimeout(90 * time.Second).
			WithPollInterval(1 * time.Second),
		Files: []tc.ContainerFile{
			{
				HostFilePath:      configFile.Name(),
				ContainerFilePath: configPath,
				FileMode:          0644,
			},
			{
				HostFilePath:      secretsFile.Name(),
				ContainerFilePath: secPath,
				FileMode:          0644,
			},
			{
				HostFilePath:      adminCredsFile.Name(),
				ContainerFilePath: adminCrePath,
				FileMode:          0644,
			},
			{
				HostFilePath:      apiCredsFile.Name(),
				ContainerFilePath: apiCrePath,
				FileMode:          0644,
			},
		},
	}, nil
}
