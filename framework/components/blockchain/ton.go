package blockchain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultTonSimpleServerPort     = 8000
	DefaultTonHTTPAPIPort          = 8081
	DefaultTonIndexAPIPort         = 8082
	DefaultTonTONExplorerPort      = 8080
	DefaultTonFaucetPort           = 88
	DefaultTonLiteServerPort       = 40004
	DefaultTonDhtServerPort        = 40003
	DefaultTonValidatorConsolePort = 40002
	DefaultTonValidatorPublicPort  = 40001
	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"
)

var (
	commonDBVars = map[string]string{
		"POSTGRES_DIALECT":                  "postgresql+asyncpg",
		"POSTGRES_HOST":                     "index-postgres",
		"POSTGRES_PORT":                     "5432",
		"POSTGRES_USER":                     "postgres",
		"POSTGRES_DB":                       "ton_index",
		"POSTGRES_PASSWORD":                 "PostgreSQL1234",
		"POSTGRES_DBNAME":                   "ton_index",
		"TON_INDEXER_TON_HTTP_API_ENDPOINT": "http://tonhttpapi:8081/",
		"TON_INDEXER_IS_TESTNET":            "0",
		"TON_INDEXER_REDIS_DSN":             "redis://redis:6379",
		"TON_WORKER_FROM":                   "1",
		"TON_WORKER_DBROOT":                 "/tondb",
		"TON_WORKER_BINARY":                 "ton-index-postgres-v2",
		"TON_WORKER_ADDITIONAL_ARGS":        "",
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

type hostPortMapping struct {
	SimpleServer string
	LiteServer   string
	DHTServer    string
	Console      string
	ValidatorUDP string
	HTTPAPIPort  string
	ExplorerPort string
	FaucetPort   string
	IndexAPIPort string
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

func randomID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano()) // fallback
	}
	return hex.EncodeToString(b)
}

func generateUniquePortsFromBase(basePort string) (*hostPortMapping, error) {
	base, err := strconv.Atoi(basePort)
	if err != nil {
		return nil, fmt.Errorf("invalid base port %s: %w", basePort, err)
	}

	// use larger multipliers to ensure no overlap between different base ports
	// 8000->0, 8001->100, 8002->200, etc.
	portMultiplier := (base - DefaultTonSimpleServerPort) * 100

	mapping := &hostPortMapping{
		SimpleServer: basePort,
		HTTPAPIPort:  strconv.Itoa(DefaultTonHTTPAPIPort + portMultiplier),
		ExplorerPort: strconv.Itoa(DefaultTonTONExplorerPort + portMultiplier),
		IndexAPIPort: strconv.Itoa(DefaultTonIndexAPIPort + portMultiplier),
		FaucetPort:   strconv.Itoa(DefaultTonFaucetPort + portMultiplier),
		LiteServer:   strconv.Itoa(DefaultTonLiteServerPort + portMultiplier),
		DHTServer:    strconv.Itoa(DefaultTonDhtServerPort + portMultiplier),
		Console:      strconv.Itoa(DefaultTonValidatorConsolePort + portMultiplier),
		ValidatorUDP: strconv.Itoa(DefaultTonValidatorPublicPort + portMultiplier),
	}

	return mapping, nil
}

func defaultTon(in *Input) {
	if in.Image == "" {
		in.Image = "ghcr.io/neodix42/mylocalton-docker:latest"
	}
	if in.Port == "" {
		in.Port = strconv.Itoa(DefaultTonSimpleServerPort)
	}
}

func newTon(in *Input) (*Output, error) {
	defaultTon(in)

	hostPorts, err := generateUniquePortsFromBase(in.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique ports: %w", err)
	}

	ctx := context.Background()

	instanceID := randomID()
	networkName := fmt.Sprintf("ton-%s", instanceID)

	_, err = testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
			Attachable:     true,
			Labels:         framework.DefaultTCLabels(),
		},
	})
	framework.L.Info().Str("output", string(networkName)).Msg("TON Docker network created")

	tonServices := []containerTemplate{
		{
			Name:  fmt.Sprintf("TON-genesis-%s", instanceID),
			Image: "ghcr.io/neodix42/mylocalton-docker:latest",
			Ports: []string{
				fmt.Sprintf("%s:8000/tcp", hostPorts.SimpleServer),
				// Note: LITE_PORT port is used by the lite-client to connect to the genesis node in config
				fmt.Sprintf("%s:%s/tcp", hostPorts.LiteServer, hostPorts.LiteServer),
				fmt.Sprintf("%s:40003/udp", hostPorts.DHTServer),
				fmt.Sprintf("%s:40002/tcp", hostPorts.Console),
				fmt.Sprintf("%s:40001/udp", hostPorts.ValidatorUDP),
			},
			Env: map[string]string{
				"GENESIS": "true",
				"NAME":    "genesis",
				// Note: LITE_PORT port is used by the lite-client to connect to the genesis node in config
				"LITE_PORT":         hostPorts.LiteServer,
				"CUSTOM_PARAMETERS": "--state-ttl 315360000 --archive-ttl 315360000",
			},
			WaitFor: wait.ForExec([]string{
				"/usr/local/bin/lite-client", "-a", fmt.Sprintf("127.0.0.1:%s", hostPorts.LiteServer), "-b",
				"E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY=", "-t", "3", "-c", "last",
			}).WithStartupTimeout(2 * time.Minute),
			Network: networkName,
			Alias:   "genesis",
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("shared-data-%s", networkName)},
					Target: "/usr/share/data",
				},
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("ton-db-%s", networkName)},
					Target: "/var/ton-work/db",
				},
			},
		},
		{
			Image:   "redis:latest",
			Name:    fmt.Sprintf("TON-redis-%s", instanceID),
			Network: networkName,
			Alias:   "redis",
		},
		{
			Image:   "postgres:17",
			Name:    fmt.Sprintf("TON-index-postgres-%s", instanceID),
			Network: networkName,
			Alias:   "index-postgres",
			Env:     commonDBVars,
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("pg-%s", networkName)},
					Target: "/var/lib/postgresql/data",
				},
			},
		},
		{
			Name:  fmt.Sprintf("TON-tonhttpapi-%s", instanceID),
			Image: "ghcr.io/neodix42/ton-http-api:latest",
			Env: map[string]string{
				"TON_API_LOGS_JSONIFY":             "0",
				"TON_API_LOGS_LEVEL":               "ERROR",
				"TON_API_TONLIB_LITESERVER_CONFIG": "/usr/share/data/global.config.json",
				"TON_API_TONLIB_CDLL_PATH":         "/usr/share/data/libtonlibjson.so",
				"TON_API_GET_METHODS_ENABLED":      "1",
				"TON_API_JSON_RPC_ENABLED":         "1",

				"POSTGRES_DIALECT":       commonDBVars["POSTGRES_DIALECT"],
				"POSTGRES_HOST":          commonDBVars["POSTGRES_HOST"],
				"POSTGRES_PORT":          commonDBVars["POSTGRES_PORT"],
				"POSTGRES_USER":          commonDBVars["POSTGRES_USER"],
				"POSTGRES_PASSWORD":      commonDBVars["POSTGRES_PASSWORD"],
				"POSTGRES_DBNAME":        commonDBVars["POSTGRES_DBNAME"],
				"TON_INDEXER_IS_TESTNET": commonDBVars["TON_INDEXER_IS_TESTNET"],
				"TON_INDEXER_REDIS_DSN":  commonDBVars["TON_INDEXER_REDIS_DSN"],
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("shared-data-%s", networkName)},
					Target: "/usr/share/data",
				},
			},
			Ports:   []string{fmt.Sprintf("%s:8081/tcp", hostPorts.HTTPAPIPort)},
			WaitFor: wait.ForHTTP("/healthcheck").WithStartupTimeout(90 * time.Second),
			Command: []string{"-c", "gunicorn -k uvicorn.workers.UvicornWorker -w 1 --bind 0.0.0.0:8081 pyTON.main:app"},
			Network: networkName,
			Alias:   "tonhttpapi",
		},
		{
			Name:  fmt.Sprintf("TON-faucet-%s", instanceID),
			Image: "ghcr.io/neodix42/mylocalton-docker-faucet:latest",
			Env: map[string]string{
				"FAUCET_USE_RECAPTCHA": "false",
				"RECAPTCHA_SITE_KEY":   "",
				"RECAPTCHA_SECRET":     "",
				"SERVER_PORT":          "88",
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("shared-data-%s", networkName)},
					Target: "/usr/share/data",
				},
			},
			Ports:   []string{fmt.Sprintf("%s:88/tcp", hostPorts.FaucetPort)},
			WaitFor: wait.ForHTTP("/").WithStartupTimeout(90 * time.Second),
			Network: networkName,
			Alias:   "faucet",
		},
	}

	tonIndexingAndObservability := []containerTemplate{
		{
			Name:  fmt.Sprintf("TON-explorer-%s", instanceID),
			Image: "ghcr.io/neodix42/mylocalton-docker-explorer:latest",
			Env: map[string]string{
				"SERVER_PORT": "8080",
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("shared-data-%s", networkName)},
					Target: "/usr/share/data",
				},
			},
			Ports:   []string{fmt.Sprintf("%s:8080/tcp", hostPorts.ExplorerPort)},
			Network: networkName,
			Alias:   "explorer",
		},
		{
			Name:  fmt.Sprintf("TON-index-worker-%s", instanceID),
			Image: "toncenter/ton-indexer-worker:v1.2.0-test",
			Env:   commonDBVars,
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("ton-db-%s", networkName)},
					Target: "/tondb",
				},
				{
					Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("index-workdir-%s", networkName)},
					Target: "/workdir",
				},
			},
			WaitFor: wait.ForLog("Starting indexing from seqno").WithStartupTimeout(90 * time.Second),
			Command: []string{"--working-dir", "/workdir", "--from", "1", "--threads", "8"},
			Network: networkName,
			Alias:   "index-worker",
		},
		{
			Name:  fmt.Sprintf("TON-index-api-%s", instanceID),
			Image: "toncenter/ton-indexer-api:v1.2.0-test",
			Env:   commonDBVars,
			Ports: []string{fmt.Sprintf("%s:8082/tcp", hostPorts.IndexAPIPort)},
			WaitFor: wait.ForHTTP("/").
				WithStartupTimeout(90 * time.Second),
			Command: []string{"-bind", ":8082", "-prefork", "-threads", "4", "-v2", "http://tonhttpapi:8081/"},
			Network: networkName,
			Alias:   "index-api",
		},
		{
			Name:  fmt.Sprintf("TON-event-classifier-%s", instanceID),
			Image: "toncenter/ton-indexer-classifier:v1.2.0-test",
			Env:   commonDBVars,
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
	if os.Getenv("CI") != "true" {
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
		Nodes: []*Node{{
			// Note: define if we need more access other than the global config(tonutils-go only uses liteclients defined in the config)
			ExternalHTTPUrl: fmt.Sprintf("%s:%s", "localhost", hostPorts.SimpleServer),
			InternalHTTPUrl: fmt.Sprintf("%s:%s", name, "8000"),
		}},
	}, nil
}
