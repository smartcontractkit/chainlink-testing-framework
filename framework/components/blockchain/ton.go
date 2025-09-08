package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultTonSimpleServerPort = "8000"
	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"

	liteServerPortOffset = 100 // internal, arbitrary offset for lite server port
)

// TON config structures (e.g.: ton-blockchain.github.io/testnet-global.config.json)
type tonLiteServer struct {
	IP   int64 `json:"ip"`
	Port int   `json:"port"`
	ID   struct {
		Key  string `json:"key"`
		Type string `json:"@type"`
	} `json:"id"`
}

type tonConfig struct {
	LiteServers []tonLiteServer `json:"liteservers"`
}

// convert int64 IP to string format (matches https://github.com/xssnick/tonutils-go/liteclient/connection.go/intToIP4)
func intToIP4(ip int64) string {
	uip := uint32(ip) //nolint:gosec // IP conversion is safe for TON format
	return fmt.Sprintf("%d.%d.%d.%d",
		(uip>>24)&0xFF,
		(uip>>16)&0xFF,
		(uip>>8)&0xFF,
		uip&0xFF)
}

// fetch and parse TON config to generate liteserver URLs
func fetchTonConfig(configURL string) ([]string, error) {
	parsedURL, err := url.Parse(configURL)
	if err != nil {
		return nil, fmt.Errorf("invalid config URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
	}

	resp, err := http.Get(configURL) //nolint:gosec // URL is validated above
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()
	defer resp.Body.Close()

	var config tonConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	var liteServerURLs []string
	for _, ls := range config.LiteServers {
		ipStr := intToIP4(ls.IP)
		url := fmt.Sprintf("liteserver://%s@%s:%d", ls.ID.Key, ipStr, ls.Port)
		liteServerURLs = append(liteServerURLs, url)
	}

	return liteServerURLs, nil
}

type portMapping struct {
	SimpleServer string
	LiteServer   string
	DHTServer    string
	Console      string
	ValidatorUDP string
}

func defaultTon(in *Input) {
	if in.Image == "" {
		in.Image = "ghcr.io/neodix42/mylocalton-docker:latest"
	}
	if in.Port == "" {
		in.Port = DefaultTonSimpleServerPort
	}
}

func newTon(in *Input) (*Output, error) {
	defaultTon(in)

	base, err := strconv.Atoi(in.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid base port %s: %w", in.Port, err)
	}

	ports := &portMapping{
		SimpleServer: in.Port,
		LiteServer:   strconv.Itoa(base + liteServerPortOffset),
	}

	ctx := context.Background()

	network, err := network.New(ctx,
		network.WithAttachable(),
		network.WithLabels(framework.DefaultTCLabels()),
	)
	if err != nil {
		return nil, err
	}
	networkName := network.Name

	baseEnv := map[string]string{
		"GENESIS":                        "true",
		"NAME":                           "genesis",
		"LITE_PORT":                      ports.LiteServer,
		"CUSTOM_PARAMETERS":              "--state-ttl 315360000 --archive-ttl 315360000",
		"EMBEDDED_FILE_HTTP_SERVER":      "true",
		"EMBEDDED_FILE_HTTP_SERVER_PORT": in.Port,
	}

	// merge with additional environment variables from input
	finalEnv := baseEnv
	if in.CustomEnv != nil {
		for key, value := range in.CustomEnv {
			finalEnv[key] = value
		}
	}

	req := testcontainers.ContainerRequest{
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Name:            framework.DefaultTCName("ton-genesis"),
		ExposedPorts: []string{
			fmt.Sprintf("%s:%s/tcp", ports.SimpleServer, DefaultTonSimpleServerPort),
			fmt.Sprintf("%s:%s/tcp", ports.LiteServer, ports.LiteServer),
			"40003/udp",
			"40002/tcp",
			"40001/udp",
		},
		Networks:       []string{networkName},
		NetworkAliases: map[string][]string{networkName: {"genesis"}},
		Labels:         framework.DefaultTCLabels(),
		Env:            finalEnv,
		WaitingFor: wait.ForExec([]string{
			"/usr/local/bin/lite-client",
			"-a", fmt.Sprintf("127.0.0.1:%s", ports.LiteServer),
			"-b", "E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY=",
			"-t", "3", "-c", "last",
		}).WithStartupTimeout(2 * time.Minute),
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
		HostConfigModifier: func(h *container.HostConfig) {
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	name, err := c.Name(ctx)
	if err != nil {
		return nil, err
	}

	// fetch config and generate liteserver URLs from actual config
	configURL := fmt.Sprintf("http://localhost:%s/localhost.global.config.json", ports.SimpleServer)

	liteServerURLs, err := fetchTonConfig(configURL)
	if err != nil {
		return nil, err
	}

	if len(liteServerURLs) == 0 {
		return nil, fmt.Errorf("no liteservers found in config")
	}

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTon,
		ContainerName: name,
		Container:     c,
		Nodes: []*Node{{
			// URLs now contain liteserver://publickey@host:port
			ExternalHTTPUrl: liteServerURLs[0],
			InternalHTTPUrl: liteServerURLs[0],
		}},
	}, nil
}
