package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/go-resty/resty/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
)

const (
	DefaultFaucetPort    = "9123/tcp"
	DefaultFaucetPortNum = "9123"
	DefaultSuiNodePort   = "9000"
	// DefaultSuiImage is the mysten/sui-tools image when Input.Image is empty on non-arm64 hosts.
	DefaultSuiImage = "mysten/sui-tools:devnet-v1.69.0"
	// DefaultSuiImageARM64 is used when Input.Image is empty on arm64 (e.g. Apple Silicon).
	DefaultSuiImageARM64 = "mysten/sui-tools:ci-arm64"
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

// demuxDockerExecOutput converts Docker exec attach output to plain text when it uses the
// multiplexed stream format (first byte 1=stdout / 2=stderr). Must run before stripping 0x01,
// which appears in stream headers and would corrupt the stream if removed globally.
func demuxDockerExecOutput(raw string) string {
	if len(raw) == 0 {
		return raw
	}
	if raw[0] != 1 && raw[0] != 2 {
		return raw
	}
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, strings.NewReader(raw)); err != nil {
		return raw
	}
	out := stdout.String() + stderr.String()
	// Invalid or partial multiplex streams can make StdCopy succeed with empty output; keep raw so
	// parseSuiKeytoolGenerateJSON can still find JSON after a single-byte preamble (e.g. 0x01).
	if out == "" {
		return raw
	}

	return out
}

// parseSuiKeytoolGenerateJSON extracts a SuiWalletInfo from `sui keytool generate --json` output.
// The CLI may print a preamble, and v1.69+ may emit compact one-line JSON; older parsers assumed a
// legacy layout (newline after '{') and corrupt compact output.
func parseSuiKeytoolGenerateJSON(keyOut string) (*SuiWalletInfo, error) {
	text := demuxDockerExecOutput(keyOut)
	s := strings.ReplaceAll(text, "\x00", "")
	for i := range s {
		if s[i] != '{' {
			continue
		}
		var key SuiWalletInfo
		dec := json.NewDecoder(bytes.NewReader([]byte(s[i:])))
		if err := dec.Decode(&key); err != nil {
			continue
		}
		if key.SuiAddress != "" {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("failed to parse SuiWalletInfo from keytool output: %.200q", keyOut)
}

// generateKeyData generates a wallet and returns all the data
func generateKeyData(ctx context.Context, containerName string, keyCipherType string) (*SuiWalletInfo, error) {
	dc, err := framework.NewDockerClient()
	if err != nil {
		return nil, err
	}

	// Ensure a valid Sui client config exists. `sui start --force-regenesis`
	// creates its config under /root/.sui/sui_config/ but the client.yaml it
	// generates may not exist yet when this runs, so we use `sui client --yes`
	// with an explicit config flag to force creation.
	initCmd := []string{"sui", "client", "--client.config", "/root/.sui/sui_config/client.yaml", "--yes", "envs"}
	if initOut, initErr := dc.ExecContainerWithContext(ctx, containerName, initCmd); initErr != nil {
		framework.L.Warn().Err(initErr).Str("out", initOut).Msg("sui client init returned error (may be harmless)")
	}

	cmdStr := []string{"sui", "keytool", "generate", keyCipherType, "--json"}
	keyOut, err := dc.ExecContainerWithContext(ctx, containerName, cmdStr)
	if err != nil {
		return nil, err
	}
	key, err := parseSuiKeytoolGenerateJSON(keyOut)
	if err != nil {
		return nil, fmt.Errorf("%w (raw output: %.300q)", err, keyOut)
	}

	framework.L.Info().Str("suiAddress", key.SuiAddress).Msg("CTF test key generated")

	return key, nil
}

func defaultSui(in *Input) {
	if in.Image == "" {
		if runtime.GOARCH == "arm64" {
			in.Image = DefaultSuiImageARM64
			if in.ImagePlatform == nil {
				arm := "linux/arm64"
				in.ImagePlatform = &arm
			}
		} else {
			in.Image = DefaultSuiImage
		}
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

	var files []testcontainers.ContainerFile
	if in.ContractsDir != "" {
		absPath, err := filepath.Abs(in.ContractsDir)
		if err != nil {
			return nil, err
		}
		files = []testcontainers.ContainerFile{
			{
				HostFilePath:      absPath,
				ContainerFilePath: "/",
			},
		}
	}

	// Sui container always listens on port 9000 internally
	containerPort := fmt.Sprintf("%s/tcp", DefaultSuiNodePort)

	imagePlatform := "linux/amd64"
	if in.ImagePlatform != nil && *in.ImagePlatform != "" {
		imagePlatform = *in.ImagePlatform
	}

	if pods.K8sEnabled() {
		return nil, fmt.Errorf("K8s support is not yet implemented")
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
<<<<<<< HEAD
		Files:      files,
=======
		Files: files,
		// we need faucet for funding
>>>>>>> main
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
		Container:           c,
		NetworkSpecificData: &NetworkSpecificData{SuiAccount: suiAccount},
		Nodes: []*Node{
			{
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s", host, in.Port),
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, DefaultSuiNodePort),
			},
		},
	}, nil
}
