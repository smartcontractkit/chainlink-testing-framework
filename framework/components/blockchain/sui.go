package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/docker/docker/api/types/container"
	"github.com/go-resty/resty/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultFaucetPort    = "9123/tcp"
	DefaultFaucetPortNum = "9123"
	DefaultSuiNodePort   = "9000"
)

// SuiWalletInfo info about Sui account/wallet
type SuiWalletInfo struct {
	Alias           *string `json:"alias"`           // Alias key name, usually "null"
	Flag            int     `json:"flag"`            // Flag is an integer
	KeyScheme       string  `json:"keyScheme"`       // Key scheme is a string
	Mnemonic        string  `json:"mnemonic"`        // Mnemonic is a string
	PeerId          string  `json:"peerId"`          // Peer ID is a string
	PublicBase64Key string  `json:"publicBase64Key"` // Public key in Base64 format
	SuiAddress      string  `json:"suiAddress"`      // Sui address is a 0x prefixed hex string
}

// funds provided key using local faucet
// we can't use the best client available - block-vision/sui-go-sdk for that, since some versions have old API and it is hardcoded
// https://github.com/block-vision/sui-go-sdk/blob/main/sui/faucet_api.go#L16
func fundAccount(url string, address string) error {
	r := resty.New().SetBaseURL(url)
	b := &models.FaucetRequest{
		FixedAmountRequest: &models.FaucetFixedAmountRequest{
			Recipient: address,
		},
	}
	resp, err := r.R().SetBody(b).SetHeader("Content-Type", "application/json").Post("/gas")
	if err != nil {
		return err
	}
	framework.L.Info().Any("Resp", resp).Msg("Address is funded!")
	return nil
}

// generateKeyData generates a wallet and returns all the data
func generateKeyData(containerName string, keyCipherType string) (*SuiWalletInfo, error) {
	cmdStr := []string{"sui", "keytool", "generate", keyCipherType, "--json"}
	dc, err := framework.NewDockerClient()
	if err != nil {
		return nil, err
	}
	keyOut, err := dc.ExecContainer(containerName, cmdStr)
	if err != nil {
		return nil, err
	}
	// formatted JSON with, no plain --json version, remove special symbols
	cleanKey := strings.ReplaceAll(keyOut, "\x00", "")
	cleanKey = strings.ReplaceAll(cleanKey, "\x01", "")
	cleanKey = strings.ReplaceAll(cleanKey, "\x02", "")
	cleanKey = strings.ReplaceAll(cleanKey, "\n", "")
	cleanKey = "{" + cleanKey[2:]
	var key *SuiWalletInfo
	if err := json.Unmarshal([]byte(cleanKey), &key); err != nil {
		return nil, err
	}
	framework.L.Info().Interface("Key", key).Msg("Test key")
	return key, nil
}

func defaultSui(in *Input) {
	if in.Image == "" {
		in.Image = "mysten/sui-tools:devnet"
	}
	if in.Port != "" {
		framework.L.Warn().Msgf("'port' field is set but only default port can be used: %s", DefaultSuiNodePort)
	}
	in.Port = DefaultSuiNodePort
}

func newSui(in *Input) (*Output, error) {
	defaultSui(in)
	ctx := context.Background()
	containerName := framework.DefaultTCName("blockchain-node")

	absPath, err := filepath.Abs(in.ContractsDir)
	if err != nil {
		return nil, err
	}

	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	req := testcontainers.ContainerRequest{
		Image:        in.Image,
		ExposedPorts: []string{in.Port, DefaultFaucetPort},
		Name:         containerName,
		Labels:       framework.DefaultTCLabels(),
		Networks:     []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort, DefaultFaucetPort)
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		ImagePlatform: "linux/amd64",
		Env: map[string]string{
			"RUST_LOG": "off,sui_node=info",
		},
		Cmd: []string{
			"sui",
			"start",
			"--force-regenesis",
			"--with-faucet",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      absPath,
				ContainerFilePath: "/",
			},
		},
		// we need faucet for funding
		WaitingFor: wait.ForListeningPort(DefaultFaucetPort).WithStartupTimeout(10 * time.Second).WithPollInterval(200 * time.Millisecond),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}
	suiAccount, err := generateKeyData(containerName, "ed25519")
	if err != nil {
		return nil, err
	}
	if err := fundAccount(fmt.Sprintf("http://%s:%s", "127.0.0.1", DefaultFaucetPortNum), suiAccount.SuiAddress); err != nil {
		return nil, err
	}
	return &Output{
		UseCache:            true,
		Type:                in.Type,
		Family:              FamilySui,
		ContainerName:       containerName,
		NetworkSpecificData: &NetworkSpecificData{SuiAccount: suiAccount},
		Nodes: []*Node{
			{
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s", host, in.Port),
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
