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
	"github.com/docker/go-connections/nat"
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
	Alias           *string `toml:"alias" json:"alias" comment:"Alias key name, usually null"`                   // Alias key name, usually "null"
	Flag            int     `toml:"flag" json:"flag" comment:"-"`                                                // Flag is an integer
	KeyScheme       string  `toml:"key_scheme" json:"keyScheme" comment:"Sui key scheme"`                        // Key scheme is a string
	Mnemonic        string  `toml:"mnemonic" json:"mnemonic" comment:"Sui key mnemonic"`                         // Mnemonic is a string
	PeerId          string  `toml:"peer_id" json:"peerId" comment:"Sui key peer ID"`                             // Peer ID is a string
	PublicBase64Key string  `toml:"public_base64_key" json:"publicBase64Key" comment:"Sui key in base64 format"` // Public key in Base64 format
	SuiAddress      string  `toml:"sui_address" json:"suiAddress" comment:"Sui key address"`                     // Sui address is a 0x prefixed hex string
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
func generateKeyData(ctx context.Context, containerName string, keyCipherType string) (*SuiWalletInfo, error) {
	cmdStr := []string{"sui", "keytool", "generate", keyCipherType, "--json"}
	dc, err := framework.NewDockerClient()
	if err != nil {
		return nil, err
	}
	keyOut, err := dc.ExecContainerWithContext(ctx, containerName, cmdStr)
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
		in.Image = "mysten/sui-tools:devnet-v1.61.0"
	}
	if in.Port == "" {
		in.Port = DefaultSuiNodePort
	}
	if in.FaucetPort == "" {
		in.FaucetPort = DefaultFaucetPortNum
	}
}

func newSui(ctx context.Context, in *Input) (*Output, error) {
	defaultSui(in)
	containerName := framework.DefaultTCName("blockchain-node")

	absPath, err := filepath.Abs(in.ContractsDir)
	if err != nil {
		return nil, err
	}

	// Sui container always listens on port 9000 internally
	containerPort := fmt.Sprintf("%s/tcp", DefaultSuiNodePort)

	// default to amd64, unless otherwise specified
	imagePlatform := "linux/amd64"
	if in.ImagePlatform != nil {
		imagePlatform = *in.ImagePlatform
	}

	req := testcontainers.ContainerRequest{
		Image:        in.Image,
		ExposedPorts: []string{containerPort, DefaultFaucetPort},
		Name:         containerName,
		Labels:       framework.DefaultTCLabels(),
		Networks:     []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		HostConfigModifier: func(h *container.HostConfig) {
			// Map user-provided host port to container's default port (9000)
			h.PortBindings = nat.PortMap{
				nat.Port(containerPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: in.Port,
					},
				},
				nat.Port(DefaultFaucetPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: in.FaucetPort,
					},
				},
			}
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		ImagePlatform: imagePlatform,
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
		WaitingFor: wait.ForListeningPort(DefaultFaucetPort).WithStartupTimeout(1 * time.Minute).WithPollInterval(200 * time.Millisecond),
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
	suiAccount, err := generateKeyData(ctx, containerName, "ed25519")
	if err != nil {
		return nil, err
	}
	if err := fundAccount(fmt.Sprintf("http://%s:%s", "127.0.0.1", in.FaucetPort), suiAccount.SuiAddress); err != nil {
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
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, DefaultSuiNodePort),
			},
		},
	}, nil
}
