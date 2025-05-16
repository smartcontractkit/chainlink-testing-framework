package blockchain

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	// default ports from mylocalton-docker
	DefaultTonHTTPAPIPort      = "8081"
	DefaultTonSimpleServerPort = "8000"
	DefaultTonTONExplorerPort  = "8080"
	DefaultTonLiteServerPort   = "40004"

	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"
)

var (
	CommonDBVars = map[string]string{
		"POSTGRES_DIALECT":                  "postgresql+asyncpg",
		"POSTGRES_HOST":                     "index-postgres",
		"POSTGRES_PORT":                     "5432",
		"POSTGRES_USER":                     "postgres",
		"POSTGRES_DB":                       "ton_index",
		"POSTGRES_PASSWORD":                 "PostgreSQL1234",
		"POSTGRES_DBNAME":                   "ton_index",
		"TON_INDEXER_TON_HTTP_API_ENDPOINT": "http://tonhttpapi:8080/",
		"TON_INDEXER_IS_TESTNET":            "0",
		"TON_INDEXER_REDIS_DSN":             "redis://redis:6379",

		"TON_WORKER_FROM":            "1",
		"TON_WORKER_DBROOT":          "/tondb",
		"TON_WORKER_BINARY":          "ton-index-postgres-v2",
		"TON_WORKER_ADDITIONAL_ARGS": "",
	}
)

type containerTemplate struct {
	Name    string
	Image   string
	Env     map[string]string
	Mounts  []testcontainers.ContainerMount
	Ports   []string
	WaitFor wait.Strategy
	Command []string
	Network string
	Alias   string
}

func commonContainer(
	ctx context.Context,
	name string,
	image string,
	env map[string]string,
	mounts []testcontainers.ContainerMount,
	exposedPorts []string,
	waitStrategy wait.Strategy,
	command []string,
	network string,
	alias string,
) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:         name,
		Labels:       framework.DefaultTCLabels(),
		Image:        image,
		Env:          env,
		Mounts:       mounts,
		ExposedPorts: exposedPorts,
		Networks:     []string{network},
		NetworkAliases: map[string][]string{
			network: {alias},
		},
		WaitingFor: waitStrategy,
		Cmd:        command,
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func defaultTon(in *Input) {
	if in.Image == "" {
		in.Image = "ghcr.io/neodix42/mylocalton-docker:latest"
	}
}

func newTon(in *Input) (*Output, error) {
	defaultTon(in)

	ctx := context.Background()

	networkName := "ton"
	lightCLientSubNet := "172.28.0.0/16"
	//nolint:gosec
	_ = exec.Command("docker", "network", "create",
		"--driver=bridge",
		"--attachable",
		fmt.Sprintf("--subnet=%s", lightCLientSubNet),
		"--label=framework=ctf",
		networkName,
	)

	tonServices := []containerTemplate{
		{
			Image: "ghcr.io/neodix42/mylocalton-docker:latest",
			Ports: []string{"8000:8000/tcp", "40004:40004/tcp", "40003:40003/udp", "40002:40002/tcp", "40001:40001/udp"},
			Env: map[string]string{
				"GENESIS":           "true",
				"NAME":              "genesis",
				"CUSTOM_PARAMETERS": "--state-ttl 315360000 --archive-ttl 315360000",
			},
			WaitFor: wait.ForExec([]string{
				"/usr/local/bin/lite-client", "-a", "127.0.0.1:40004", "-b",
				"E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY=", "-t", "3", "-c", "last",
			}).WithStartupTimeout(2 * time.Minute),
			Network: networkName,
			Alias:   "genesis",
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "shared-data"},
					Target: "/usr/share/data",
				},
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "ton-db"},
					Target: "/var/ton-work/db",
				},
			},
		},
		{
			Image:   "redis:latest",
			Name:    "redis",
			Network: networkName,
			Alias:   "redis",
		},
		{
			Image:   "postgres:17",
			Name:    "index-postgres",
			Network: networkName,
			Alias:   "index-postgres",
			Env:     CommonDBVars,
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "pg"},
					Target: "/var/lib/postgresql/data",
				},
			},
		},
		{
			Name:  "tonhttpapi",
			Image: "ghcr.io/neodix42/ton-http-api:latest",
			Env: map[string]string{
				"TON_API_LOGS_JSONIFY":             "0",
				"TON_API_LOGS_LEVEL":               "ERROR",
				"TON_API_TONLIB_LITESERVER_CONFIG": "/usr/share/data/global.config.json",
				"TON_API_TONLIB_CDLL_PATH":         "/usr/share/data/libtonlibjson.so",
				"TON_API_GET_METHODS_ENABLED":      "1",
				"TON_API_JSON_RPC_ENABLED":         "1",

				"POSTGRES_DIALECT":       "postgresql+asyncpg",
				"POSTGRES_HOST":          "index-postgres",
				"POSTGRES_PORT":          "5432",
				"POSTGRES_USER":          "postgres",
				"POSTGRES_PASSWORD":      "PostgreSQL1234",
				"POSTGRES_DBNAME":        "ton_index",
				"TON_INDEXER_IS_TESTNET": "0",
				"TON_INDEXER_REDIS_DSN":  "redis://redis:6379",
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "shared-data"},
					Target: "/usr/share/data",
				},
			},
			Ports:   []string{"8081/tcp"},
			WaitFor: wait.ForHTTP("/healthcheck").WithStartupTimeout(90 * time.Second),
			Command: []string{"-c", "gunicorn -k uvicorn.workers.UvicornWorker -w 1 --bind 0.0.0.0:8081 pyTON.main:app"},
			Network: networkName,
			Alias:   "tonhttpapi",
		},
		{
			Name:  "faucet",
			Image: "ghcr.io/neodix42/mylocalton-docker-faucet:latest",
			Env: map[string]string{
				"FAUCET_USE_RECAPTCHA": "false",
				"RECAPTCHA_SITE_KEY":   "",
				"RECAPTCHA_SECRET":     "",
				"SERVER_PORT":          "88",
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "shared-data"},
					Target: "/usr/share/data",
				},
			},
			Ports:   []string{"88/tcp"},
			WaitFor: wait.ForHTTP("/").WithStartupTimeout(90 * time.Second),
			Network: networkName,
			Alias:   "faucet",
		},
	}

	tonIndexingAndObservability := []containerTemplate{
		{
			Name:  "explorer",
			Image: "ghcr.io/neodix42/mylocalton-docker-explorer:latest",
			Env: map[string]string{
				"SERVER_PORT": "8080",
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "shared-data"},
					Target: "/usr/share/data",
				},
			},
			Ports:   []string{"8080:8080/tcp"},
			Network: networkName,
			Alias:   "explorer",
		},
		{
			Name:  "index-worker",
			Image: "toncenter/ton-indexer-worker:v1.2.0-test",
			Env:   CommonDBVars,
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "ton-db"},
					Target: "/tondb",
				},
				{
					Source: testcontainers.GenericVolumeMountSource{Name: "index-workdir"},
					Target: "/workdir",
				},
			},
			WaitFor: wait.ForLog("Starting indexing from seqno").WithStartupTimeout(90 * time.Second),
			Command: []string{"--working-dir", "/workdir", "--from", "1", "--threads", "8"},
			Network: networkName,
			Alias:   "index-worker",
		},
		{
			Name:  "index-api",
			Image: "toncenter/ton-indexer-api:v1.2.0-test",
			Env:   CommonDBVars,
			Ports: []string{"8082/tcp"},
			WaitFor: wait.ForHTTP("/").
				WithStartupTimeout(90 * time.Second),
			Command: []string{"-bind", ":8082", "-prefork", "-threads", "4", "-v2", "http://tonhttpapi:8081/"},
			Network: networkName,
			Alias:   "index-api",
		},
		{
			Name:  "event-classifier",
			Image: "toncenter/ton-indexer-classifier:v1.2.0-test",
			Env:   CommonDBVars,
			WaitFor: wait.ForLog("Reading finished tasks").
				WithStartupTimeout(90 * time.Second),
			Command: []string{"--pool-size", "4", "--prefetch-size", "1000", "--batch-size", "100"},
			Network: networkName,
			Alias:   "event-classifier",
		},
	}

	containers := make([]testcontainers.Container, 0)
	for _, s := range tonServices {
		c, err := commonContainer(ctx, s.Name, s.Image, s.Env, s.Mounts, s.Ports, s.WaitFor, s.Command, s.Network, s.Alias)
		if err != nil {
			return nil, fmt.Errorf("failed to start %s: %v", s.Name, err)
		}
		containers = append(containers, c)
	}
	// no need for indexers and block explorers in CI
	if os.Getenv("CI") != "" {
		for _, s := range tonIndexingAndObservability {
			c, err := commonContainer(ctx, s.Name, s.Image, s.Env, s.Mounts, s.Ports, s.WaitFor, s.Command, s.Network, s.Alias)
			if err != nil {
				return nil, fmt.Errorf("failed to start %s: %v", s.Name, err)
			}
			containers = append(containers, c)
		}
	}

	genesisTonContainer := containers[0]

	name, err := genesisTonContainer.Name(ctx)
	if err != nil {
		return nil, err
	}
	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTon,
		ContainerName: name,
		// Note: in case we need 1+ validators, we need to modify the compose file
		Nodes: []*Node{{
			// Note: define if we need more access other than the global config(tonutils-go only uses liteclients defined in the config)
			ExternalHTTPUrl: fmt.Sprintf("%s:%s", "localhost", DefaultTonSimpleServerPort),
			InternalHTTPUrl: fmt.Sprintf("%s:%s", name, DefaultTonSimpleServerPort),
		}},
	}, nil
}
