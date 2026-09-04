package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"net/netip"

	"github.com/avast/retry-go/v4"
	"github.com/block-vision/sui-go-sdk/models"

	"github.com/go-resty/resty/v2"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
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
	DefaultSuiImage = "mysten/sui-tools:devnet-v1.73.0"
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

// faucetFundTimeout bounds the total time spent retrying Sui faucet /gas funding.
// The faucet is served by the same container that just started, and the container
// readiness gate only waits for the TCP port to listen — not for the faucet's HTTP
// handler to be ready. So the first /gas request can race the faucet's own readiness
// and fail with a connection reset. Retry with backoff until the faucet accepts it.
const faucetFundTimeout = 2 * time.Minute

// faucetRequestTimeout bounds a single faucet /gas HTTP request so one hung attempt
// does not consume the entire faucetFundTimeout budget, leaving room for retries.
const faucetRequestTimeout = 10 * time.Second

// faucetFundAttempts is the maximum number of /gas funding attempts before giving up.
// The faucetFundTimeout context still bounds the wall-clock total.
const faucetFundAttempts = uint(15)

// funds provided key using local faucet
// we can't use the best client available - block-vision/sui-go-sdk for that, since some versions have old API and it is hardcoded
// https://github.com/block-vision/sui-go-sdk/blob/main/sui/faucet_api.go#L16
func fundAccount(ctx context.Context, url string, address string) error {
	// Bound the overall retry. A child of the caller's context wins on the earlier
	// deadline, so callers can tighten it further.
	ctx, cancel := context.WithTimeout(ctx, faucetFundTimeout)
	defer cancel()

	r := resty.New().
		SetBaseURL(url).
		SetTimeout(faucetRequestTimeout)

	b := &models.FaucetRequest{
		FixedAmountRequest: &models.FaucetFixedAmountRequest{
			Recipient: address,
		},
	}
	_, err := retry.DoWithData(func() (*resty.Response, error) {
		resp, perr := r.R().
			SetContext(ctx).
			SetBody(b).
			SetHeader("Content-Type", "application/json").
			Post("/gas")
		if perr != nil {
			return nil, perr
		}
		if resp.IsError() {
			return nil, &faucetStatusError{status: resp.StatusCode(), body: resp.Body()}
		}
		return resp, nil
	},
		retry.Context(ctx),
		retry.Attempts(faucetFundAttempts),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(isRetryableFaucetErr),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			framework.L.Warn().Err(err).Uint("attempt", n+1).Uint("attempts", faucetFundAttempts).
				Str("recipient", address).Msg("Retrying Sui faucet /gas funding")
		}),
	)
	if err != nil {
		return fmt.Errorf("fund account via Sui faucet: %w", err)
	}
	framework.L.Info().Str("recipient", address).Msg("Address is funded!")
	return nil
}

// faucetStatusError carries the HTTP status and body of a non-2xx faucet response so
// isRetryableFaucetErr can classify it without re-issuing the request.
type faucetStatusError struct {
	status int
	body   []byte
}

func (e *faucetStatusError) Error() string {
	if trimmed := bytes.TrimSpace(e.body); len(trimmed) > 0 {
		return fmt.Sprintf("faucet returned status %d: %s", e.status, trimmed)
	}
	return fmt.Sprintf("faucet returned status %d", e.status)
}

// isRetryableFaucetErr classifies faucet /gas errors. Transient failures (faucet still
// warming up, brief network blips, rate limiting, 5xx) are retried; failures that retrying
// cannot fix (a malformed request, an already-cancelled context, 4xx other than 429) stop
// immediately so they don't burn the retry budget and prolong startup/teardown.
func isRetryableFaucetErr(err error) bool {
	if err == nil {
		return false
	}
	// A cancelled/deadline-exceeded context means the caller (or the total budget) has
	// given up; never retry, just propagate.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var se *faucetStatusError
	if errors.As(err, &se) {
		switch {
		case se.status == http.StatusTooManyRequests:
			// Rate limited; back off and try again.
			return true
		case se.status >= 400 && se.status < 500:
			// Client-side error (bad recipient, unauthorized, not found, ...).
			// Retrying an identical request won't fix it.
			return false
		case se.status >= 500:
			// Server-side error / faucet not yet ready.
			return true
		}
	}
	// Transport-level failures (connection refused/reset, EOF, DNS, timeouts) are
	// treated as transient while the faucet container comes up.
	return true
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
		return nil, fmt.Errorf("failed to parse sui keytool generate output: %w", err)
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
			h.PortBindings = network.PortMap{
				network.MustParsePort(containerPort): []network.PortBinding{
					{
						HostIP:   netip.MustParseAddr("0.0.0.0"),
						HostPort: in.Port,
					},
				},
				network.MustParsePort(DefaultFaucetPort): []network.PortBinding{
					{
						HostIP:   netip.MustParseAddr("0.0.0.0"),
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
		Files: files,
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
	if err := fundAccount(ctx, fmt.Sprintf("http://%s:%s", "127.0.0.1", in.FaucetPort), suiAccount.SuiAddress); err != nil {
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
