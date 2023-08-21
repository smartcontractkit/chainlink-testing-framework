package test_env

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

const (
	BLKSCOUT_IMAGE     = "blockscout/blockscout:5.2.1.commit.14c2a7cc"
	BLKSCOUT_HTTP_PORT = "4000"
	BLKSCOUT_PASS      = "blockscoutpass"
	BLCSCOUT_DB_NAME   = "blockscout"
)

type Blockscout struct {
	EnvComponent
	NetworkConfig blockchain.EVMNetwork
	HTTPURL       string
	Image         string
	pgContainer   *PostgresDb
}

func NewBlockscout(networks []string, networkConfig blockchain.EVMNetwork) *Blockscout {
	return &Blockscout{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "blockscout", networkConfig.Name),
			Networks:      networks,
		},
		NetworkConfig: networkConfig,
		Image:         BLKSCOUT_IMAGE,
	}
}

type BlockscoutOption = func(c *Blockscout)

func WithImage(image string) BlockscoutOption {
	return func(b *Blockscout) {
		if image == "" {
			return
		}
		b.Image = image
	}
}

func (b *Blockscout) Start(opts ...BlockscoutOption) error {
	for _, opt := range opts {
		opt(b)
	}
	b.pgContainer = NewPostgresDb(b.Networks,
		WithPostgresDbContainerName(fmt.Sprintf("pg-%s", b.ContainerName)),
		WithPostgresDbDatabaseName(BLCSCOUT_DB_NAME),
		WithPostgresDbPassword(BLKSCOUT_PASS),
	)
	err := b.pgContainer.StartContainer()
	if err != nil {
		return err
	}
	req := b.getBlockscoutContainerRequest()
	blockscout, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return err
	}
	host, err := blockscout.Host(context.Background())
	if err != nil {
		return err
	}
	port, err := blockscout.MappedPort(context.Background(), BLKSCOUT_HTTP_PORT)
	if err != nil {
		return err
	}
	b.HTTPURL = fmt.Sprintf("http://%s:%s", host, port.Port())
	log.Info().
		Str("Network", b.NetworkConfig.Name).
		Str("url", b.HTTPURL).
		Msg("Started Blockscout")
	return nil
}

func (b *Blockscout) getBlockscoutContainerRequest() tc.ContainerRequest {
	return tc.ContainerRequest{
		Name:  b.ContainerName,
		Image: b.Image,
		ExposedPorts: []string{
			fmt.Sprintf("%s/tcp", BLKSCOUT_HTTP_PORT),
		},
		Networks: b.Networks,
		WaitingFor: tcwait.ForHTTP("/").
			WithPort(nat.Port(fmt.Sprintf("%s/tcp", BLKSCOUT_HTTP_PORT))),
		Cmd: []string{
			"bin/blockscout", "eval", `"Elixir.Explorer.ReleaseTasks.create_and_migrate()\"`, "&&",
			"bin/blockscout", "start",
		},
		Env: map[string]string{
			"MIX_ENV":                   "prod",
			"ECTO_USE_SSL":              "'false'",
			"COIN":                      "DAI",
			"ETHEREUM_JSONRPC_VARIANT":  "geth",
			"ETHEREUM_JSONRPC_HTTP_URL": b.NetworkConfig.HTTPURLs[0],
			"ETHEREUM_JSONRPC_WS_URL":   b.NetworkConfig.URLs[0],
			"DATABASE_URL": fmt.Sprintf(
				"postgresql://%s:%s@%s:5432/%s?ssl=false",
				b.pgContainer.User,
				b.pgContainer.Password,
				b.pgContainer.ContainerName,
				b.pgContainer.DbName),
		},
	}
}
