package seth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	ErrEmptyConfigPath          = "toml config path is empty, set SETH_CONFIG_PATH"
	ErrContractDeploymentFailed = "contract deployment failed"

	// unused by Seth, but used by upstream
	ErrNoKeyLoaded = "failed to load private key"

	ContractMapFilePattern          = "deployed_contracts_%s_%s.toml"
	RevertedTransactionsFilePattern = "reverted_transactions_%s_%s.json"
)

var (
	// Amount of funds that will be left on the root key, when splitting funds between ephemeral addresses
	ZeroInt64 int64

	TracingLevel_None     = "NONE"
	TracingLevel_Reverted = "REVERTED"
	TracingLevel_All      = "ALL"

	TraceOutput_Console = "console"
	TraceOutput_JSON    = "json"
	TraceOutput_DOT     = "dot"
)

// Client is a vanilla go-ethereum client with enhanced debug logging
type Client struct {
	Cfg                      *Config
	Client                   simulated.Client
	Addresses                []common.Address
	PrivateKeys              []*ecdsa.PrivateKey
	ChainID                  int64
	URL                      string
	Context                  context.Context
	CancelFunc               context.CancelFunc
	Errors                   []error
	ContractStore            *ContractStore
	NonceManager             *NonceManager
	Tracer                   *Tracer
	ContractAddressToNameMap ContractMap
	ABIFinder                *ABIFinder
	HeaderCache              *LFUHeaderCache
}

// NewClientWithConfig creates a new seth client with all deps setup from config
func NewClientWithConfig(cfg *Config) (*Client, error) {
	initDefaultLogging()

	if cfg == nil {
		return nil, fmt.Errorf("seth configuration is nil. " +
			"Ensure you're calling NewClientWithConfig() with a valid config, or use NewClient() to load from SETH_CONFIG_PATH environment variable. " +
			"See documentation for configuration examples")
	}
	if cfgErr := cfg.Validate(); cfgErr != nil {
		return nil, cfgErr
	}

	L.Debug().Msgf("Using tracing level: %s", cfg.TracingLevel)

	cfg.setEphemeralAddrs()
	cs, err := NewContractStore(filepath.Join(cfg.ConfigDir, cfg.ABIDir), filepath.Join(cfg.ConfigDir, cfg.BINDir), cfg.GethWrappersDirs)
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI/contract store: %w\n"+
			"Check that:\n"+
			"  1. 'abi_dir' path is correct (current: %s)\n"+
			"  2. 'bin_dir' path is correct (current: %s)\n"+
			"  3. These directories exist and are readable\n"+
			"  4. Or comment out these settings if not using ABI/BIN files",
			err, cfg.ABIDir, cfg.BINDir)
	}
	if cfg.ephemeral {
		// we don't care about any other keys, only the root key
		// you should not use ephemeral mode with more than 1 key
		if len(cfg.Network.PrivateKeys) > 1 {
			L.Warn().Msg("Ephemeral mode is enabled, but more than 1 key is loaded. Only the first key will be used")
		}
		cfg.Network.PrivateKeys = cfg.Network.PrivateKeys[:1]
		pkeys, err := NewEphemeralKeys(*cfg.EphemeralAddrs)
		if err != nil {
			return nil, err
		}
		cfg.Network.PrivateKeys = append(cfg.Network.PrivateKeys, pkeys...)
	}
	addrs, pkeys, err := cfg.ParseKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to parse private keys: %w\n"+
			"Ensure private keys are valid hex strings (64 characters, without 0x prefix). "+
			"Keys should be set in SETH_ROOT_PRIVATE_KEY env var or 'private_keys_secret' in seth.toml",
			err)
	}
	nm, err := NewNonceManager(cfg, addrs, pkeys)
	if err != nil {
		return nil, fmt.Errorf("failed to create nonce manager: %w\n"+
			"This is usually a configuration issue. Check 'nonce_manager' settings in your config (seth.toml or ClientBuilder)",
			err)
	}

	if !cfg.IsSimulatedNetwork() && cfg.SaveDeployedContractsMap && cfg.ContractMapFile == "" {
		cfg.ContractMapFile = cfg.GenerateContractMapFileName()
	}

	// this part is kind of duplicated in NewClientRaw, but we need to create contract map before creating Tracer
	// so that both the tracer and client have references to the same map
	contractAddressToNameMap := NewEmptyContractMap()
	contractAddressToNameMap.addressMap = make(map[string]string)
	if !cfg.IsSimulatedNetwork() {
		contractAddressToNameMap.addressMap, err = LoadDeployedContracts(cfg.ContractMapFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load deployed contracts map from '%s': %w\n"+
				"If this is a fresh deployment or you don't need contract mapping, ignore this error by setting save_deployed_contracts_map = false",
				cfg.ContractMapFile, err)
		}
	} else {
		L.Debug().Msg("Simulated network, contract map won't be read from file")
	}

	abiFinder := NewABIFinder(contractAddressToNameMap, cs)

	var opts []ClientOpt

	// even if the ethclient that was passed supports tracing, we still need the RPC URL, because we cannot get from
	// the instance of ethclient, since it doesn't expose any such method
	if (cfg.ethclient != nil && shouldInitializeTracer(cfg.ethclient, cfg) && len(cfg.Network.URLs) > 0) || cfg.ethclient == nil {
		tr, err := NewTracer(cs, &abiFinder, cfg, contractAddressToNameMap, addrs)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction tracer: %w\n"+
				"Possible causes:\n"+
				"  1. RPC endpoint is not accessible\n"+
				"  2. RPC node doesn't support debug_traceTransaction API\n"+
				"Solutions:\n"+
				"  1. Verify RPC endpoint URL is correct and accessible\n"+
				"  2. Use an RPC node with debug API enabled (archive node recommended)\n"+
				"  3. Set tracing_level = 'NONE' in config to disable tracing",
				err)
		}
		opts = append(opts, WithTracer(tr))
	}

	opts = append(opts, WithContractStore(cs), WithNonceManager(nm), WithContractMap(contractAddressToNameMap), WithABIFinder(&abiFinder))

	return NewClientRaw(
		cfg,
		addrs,
		pkeys,
		opts...,
	)
}

// NewClient creates a new raw seth client with all deps setup from env vars
func NewClient() (*Client, error) {
	cfg, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	return NewClientWithConfig(cfg)
}

// NewClientRaw creates a new raw seth client without dependencies
func NewClientRaw(
	cfg *Config,
	addrs []common.Address,
	pkeys []*ecdsa.PrivateKey,
	opts ...ClientOpt,
) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("seth configuration is nil. " +
			"Provide a valid Config when calling NewClientRaw(). " +
			"Consider using NewClient() or NewClientWithConfig() instead")
	}
	if cfgErr := cfg.Validate(); cfgErr != nil {
		return nil, cfgErr
	}
	if cfg.ReadOnly && (len(addrs) > 0 || len(pkeys) > 0) {
		return nil, fmt.Errorf("configuration conflict: read-only mode is enabled, but private keys were provided. " +
			"Read-only mode is for querying blockchain state only (no transactions).\n" +
			"To fix:\n" +
			"  1. Remove private keys if you only need to read data\n" +
			"  2. Set 'read_only = false' in config if you need to send transactions")
	}

	var firstUrl string
	var client simulated.Client
	if cfg.ethclient == nil {
		L.Info().Msg("Creating new ethereum client")
		if len(cfg.Network.URLs) == 0 {
			return nil, fmt.Errorf("no RPC URLs provided. " +
				"Set RPC URLs in your seth.toml config under 'urls_secret = [\"http://...\"]' or provide via WithRpcUrl() when using ClientBuilder")
		}

		if len(cfg.Network.URLs) > 1 {
			L.Warn().Msg("Multiple RPC URLs provided, only the first one will be used")
		}

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Network.DialTimeout.Duration())
		defer cancel()
		rpcClient, err := rpc.DialOptions(ctx,
			cfg.MustFirstNetworkURL(),
			rpc.WithHeaders(cfg.RPCHeaders),
			rpc.WithHTTPClient(&http.Client{
				Transport: NewLoggingTransport(),
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RPC endpoint '%s': %w\n"+
				"Troubleshooting steps:\n"+
				"  1. Verify the URL is correct and accessible\n"+
				"  2. Check if the RPC node is running\n"+
				"  3. Verify network connectivity and firewall rules\n"+
				"  4. Check if dial_timeout (%s) is sufficient for your network",
				cfg.MustFirstNetworkURL(), err, cfg.Network.DialTimeout.String())
		}
		client = ethclient.NewClient(rpcClient)
		firstUrl = cfg.MustFirstNetworkURL()
	} else {
		L.Info().
			Str("Type", reflect.TypeOf(cfg.ethclient).String()).
			Msg("Using provided ethereum client")
		client = cfg.ethclient
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	c := &Client{
		Client:      client,
		Cfg:         cfg,
		Addresses:   addrs,
		PrivateKeys: pkeys,
		URL:         firstUrl,
		ChainID:     mustSafeInt64(cfg.Network.ChainID),
		Context:     ctx,
		CancelFunc:  cancelFunc,
	}

	for _, o := range opts {
		o(c)
	}

	if cfg.Network.ChainID == 0 {
		chainId, err := c.Client.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID from RPC: %w\n"+
				"Ensure the RPC endpoint is accessible and the network is running. "+
				"You can also set 'chain_id' explicitly in your seth.toml to avoid this check",
				err)
		}
		cfg.Network.ChainID = chainId.Uint64()
		c.ChainID = mustSafeInt64(cfg.Network.ChainID)
	}

	var err error

	if c.ContractAddressToNameMap.addressMap == nil {
		c.ContractAddressToNameMap = NewEmptyContractMap()
		if !cfg.IsSimulatedNetwork() {
			c.ContractAddressToNameMap.addressMap, err = LoadDeployedContracts(cfg.ContractMapFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load deployed contracts map from '%s': %w\n"+
					"If this is a fresh deployment or you don't need contract mapping, ignore this error by setting save_deployed_contracts_map = false",
					cfg.ContractMapFile, err)
			}
			if len(c.ContractAddressToNameMap.addressMap) > 0 {
				L.Info().
					Int("Size", len(c.ContractAddressToNameMap.addressMap)).
					Str("File name", cfg.ContractMapFile).
					Msg("No contract map provided, read it from file")
			} else {
				L.Info().
					Msg("No contract map provided and no file found, created new one")
			}
		} else {
			L.Debug().Msg("Simulated network, contract map won't be read from file")
			L.Info().
				Msg("No contract map provided and no file found, created new one")
		}
	} else {
		L.Info().
			Int("Size", len(c.ContractAddressToNameMap.addressMap)).
			Msg("Contract map was provided")
	}
	if c.NonceManager != nil {
		c.NonceManager.Client = c
		if len(c.Cfg.Network.PrivateKeys) > 0 {
			if err := c.NonceManager.UpdateNonces(); err != nil {
				return nil, err
			}
		}
	}

	if cfg.CheckRpcHealthOnStart {
		if cfg.ReadOnly {
			return nil, fmt.Errorf("RPC health check is not supported in read-only mode because it requires sending transactions. " +
				"Either:\n" +
				"  1. Set 'read_only = false' to enable transaction capabilities\n" +
				"  2. Set 'check_rpc_health_on_start = false' to skip the health check")
		}
		if c.NonceManager == nil {
			L.Debug().Msg("Nonce manager is not set, RPC health check will be skipped. Client will most probably fail on first transaction")
		} else {
			if err := c.checkRPCHealth(); err != nil {
				return nil, err
			}
		}
	}

	if cfg.PendingNonceProtectionEnabled && cfg.ReadOnly {
		return nil, fmt.Errorf("pending nonce protection is not supported in read-only mode because it requires transaction monitoring. " +
			"Either:\n" +
			"  1. Set 'read_only = false' to enable transaction capabilities\n" +
			"  2. Set 'pending_nonce_protection_enabled = false'")
	}

	cfg.setEphemeralAddrs()

	L.Info().
		Str("NetworkName", cfg.Network.Name).
		Interface("Addresses", addrs).
		Str("RPC", firstUrl).
		Uint64("ChainID", cfg.Network.ChainID).
		Int64("Ephemeral keys", *cfg.EphemeralAddrs).
		Msg("Created new client")

	if cfg.ephemeral {
		if len(c.Addresses) == 0 {
			return nil, fmt.Errorf("ephemeral mode requires exactly one root private key to fund ephemeral addresses, but no keys were loaded. " +
				"Load the root private key via:\n" +
				"  1. SETH_ROOT_PRIVATE_KEY environment variable\n" +
				"  2. 'root_private_key' in seth.toml\n" +
				"  3. WithPrivateKeys() when using ClientBuilder")
		}
		if cfg.ReadOnly {
			return nil, fmt.Errorf("ephemeral mode is not supported in read-only mode because it requires funding transactions. " +
				"Set 'read_only = false' or disable ephemeral mode by removing 'ephemeral_addresses_number' from config")
		}
		ctx, cancel := context.WithTimeout(context.Background(), c.Cfg.Network.TxnTimeout.D)
		defer cancel()

		gasPrice, err := c.GetSuggestedLegacyFees(ctx, c.Cfg.Network.GasPriceEstimationTxPriority)
		if err != nil {
			gasPrice = big.NewInt(c.Cfg.Network.GasPrice)
		}

		bd, err := c.CalculateSubKeyFunding(*cfg.EphemeralAddrs, gasPrice.Int64(), *cfg.RootKeyFundsBuffer)
		if err != nil {
			return nil, err
		}
		L.Warn().Msg("Ephemeral mode, all funds will be lost!")

		eg, egCtx := errgroup.WithContext(ctx)
		// root key is element 0 in ephemeral
		for _, addr := range c.Addresses[1:] {
			eg.Go(func() error {
				return c.TransferETHFromKey(egCtx, 0, addr.Hex(), bd.AddrFunding, gasPrice)
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}
	}

	// we cannot use the tracer with simulated backend, because it doesn't expose a method to get rpcClient (even though it has one)
	// and Tracer needs rpcClient to call debug_traceTransaction
	if shouldInitializeTracer(c.Client, cfg) && c.Cfg.TracingLevel != TracingLevel_None && c.Tracer == nil {
		if c.ContractStore == nil {
			cs, err := NewContractStore(filepath.Join(cfg.ConfigDir, cfg.ABIDir), filepath.Join(cfg.ConfigDir, cfg.BINDir), cfg.GethWrappersDirs)
			if err != nil {
				return nil, fmt.Errorf("failed to create contract store for tracing: %w\n"+
					"Tracing requires ABI files. Check 'abi_dir' configuration or disable tracing with tracing_level = 'NONE'",
					err)
			}
			c.ContractStore = cs
		}
		if c.ABIFinder == nil {
			abiFinder := NewABIFinder(c.ContractAddressToNameMap, c.ContractStore)
			c.ABIFinder = &abiFinder
		}
		tr, err := NewTracer(c.ContractStore, c.ABIFinder, cfg, c.ContractAddressToNameMap, addrs)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction tracer: %w\n"+
				"Possible causes:\n"+
				"  1. RPC endpoint is not accessible\n"+
				"  2. RPC node doesn't support debug_traceTransaction API\n"+
				"Solutions:\n"+
				"  1. Verify RPC endpoint URL is correct and accessible\n"+
				"  2. Use an RPC node with debug API enabled (archive node recommended)\n"+
				"  3. Set tracing_level = 'NONE' in config to disable tracing",
				err)
		}

		c.Tracer = tr
	}

	now := time.Now().Format("2006-01-02-15-04-05")
	c.Cfg.revertedTransactionsFile = filepath.Join(c.Cfg.ArtifactsDir, fmt.Sprintf(RevertedTransactionsFilePattern, c.Cfg.Network.Name, now))

	if c.Cfg.Network.GasPriceEstimationEnabled {
		L.Debug().Msg("Gas estimation is enabled")
		L.Debug().Msg("Initializing LFU block header cache")
		c.HeaderCache = NewLFUBlockCache(c.Cfg.Network.GasPriceEstimationBlocks)

		if c.Cfg.Network.EIP1559DynamicFees {
			L.Debug().Msg("Checking if EIP-1559 is supported by the network")
			c.CalculateGasEstimations(GasEstimationRequest{
				GasEstimationEnabled: true,
				FallbackGasPrice:     c.Cfg.Network.GasPrice,
				FallbackGasFeeCap:    c.Cfg.Network.GasFeeCap,
				FallbackGasTipCap:    c.Cfg.Network.GasTipCap,
				Priority:             Priority_Standard,
			})
		}
	}

	if c.Cfg.GasBump != nil && c.Cfg.GasBump.Retries != 0 && c.Cfg.ReadOnly {
		return nil, fmt.Errorf("gas bumping is not supported in read-only mode because it requires sending replacement transactions. " +
			"Either:\n" +
			"  1. Set 'read_only = false' to enable transaction capabilities\n" +
			"  2. Set 'gas_bump.retries = 0' to disable gas bumping")
	}

	// if gas bumping is enabled, but no strategy is set, we set the default one; otherwise we set the no-op strategy (defensive programming to avoid NPE)
	if c.Cfg.GasBump != nil && c.Cfg.GasBump.StrategyFn == nil {
		if c.Cfg.GasBumpRetries() != 0 {
			c.Cfg.GasBump.StrategyFn = PriorityBasedGasBumpingStrategyFn(c.Cfg.Network.GasPriceEstimationTxPriority)
		} else {
			c.Cfg.GasBump.StrategyFn = NoOpGasBumpStrategyFn
		}
	}

	return c, nil
}

func (m *Client) checkRPCHealth() error {
	L.Info().Str("RPC node", m.URL).Msg("---------------- !!!!! ----------------> Checking RPC health")
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Network.TxnTimeout.Duration())
	defer cancel()

	gasPrice, err := m.GetSuggestedLegacyFees(context.Background(), m.Cfg.Network.GasPriceEstimationTxPriority)
	if err != nil {
		gasPrice = big.NewInt(m.Cfg.Network.GasPrice)
	}

	if err := m.validateAddressesKeyNum(0); err != nil {
		return err
	}

	err = m.TransferETHFromKey(ctx, 0, m.Addresses[0].Hex(), big.NewInt(10_000), gasPrice)
	if err != nil {
		return fmt.Errorf("RPC health check failed: %w\n"+
			"The health check sends a small self-transfer transaction to verify RPC functionality.\n"+
			"Possible issues:\n"+
			"  1. RPC node is not accepting transactions\n"+
			"  2. Root key has insufficient balance\n"+
			"  3. Network connectivity problems\n"+
			"  4. Gas price estimation failed\n"+
			"You can disable this check with 'check_rpc_health_on_start = false' in config (seth.toml or ClientBuilder)",
			err)
	}

	L.Info().Msg("RPC health check passed <---------------- !!!!! ----------------")
	return nil
}

// TransferETHFromKey initiates a transfer of Ether from a specified key to a recipient address.
// It validates the private key index, estimates gas limit if not provided, and sends the transaction.
// The function signs and sends an Ethereum transaction using the specified fromKeyNum and recipient address,
// with the specified amount and gas price.
func (m *Client) TransferETHFromKey(ctx context.Context, fromKeyNum int, to string, value *big.Int, gasPrice *big.Int) error {
	if err := m.validatePrivateKeysKeyNum(fromKeyNum); err != nil {
		return err
	}
	toAddr := common.HexToAddress(to)
	ctx, chainCancel := context.WithTimeout(ctx, m.Cfg.Network.TxnTimeout.Duration())
	defer chainCancel()

	chainID, err := m.Client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID from RPC: %w\n"+
			"Ensure the RPC endpoint is accessible. Chain ID is required for EIP-155 transaction signing",
			err)
	}

	var gasLimit int64
	//nolint
	gasLimitRaw, err := m.EstimateGasLimitForFundTransfer(m.Addresses[fromKeyNum], common.HexToAddress(to), value)
	if err != nil {
		gasLimit = m.Cfg.Network.TransferGasFee
	} else {
		gasLimit = mustSafeInt64(gasLimitRaw)
	}

	if gasPrice == nil {
		gasPrice = big.NewInt(m.Cfg.Network.GasPrice)
	}

	rawTx := &types.LegacyTx{
		Nonce:    m.NonceManager.NextNonce(m.Addresses[fromKeyNum]).Uint64(),
		To:       &toAddr,
		Value:    value,
		Gas:      mustSafeUint64(gasLimit),
		GasPrice: gasPrice,
	}
	L.Debug().Interface("TransferTx", rawTx).Send()
	signedTx, err := types.SignNewTx(m.PrivateKeys[fromKeyNum], types.NewEIP155Signer(chainID), rawTx)
	if err != nil {
		return fmt.Errorf("failed to sign transaction with key #%d (address: %s): %w\n"+
			"Verify the private key is valid and corresponds to the expected address",
			fromKeyNum, m.Addresses[fromKeyNum].Hex(), err)
	}

	err = m.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction to network: %w\n"+
			"Common causes:\n"+
			"  1. RPC node rejected the transaction\n"+
			"  2. Gas price too low (try increasing gas_price or gas_fee_cap)\n"+
			"  3. Nonce conflict (transaction with same nonce already pending)\n"+
			"  4. Insufficient funds for gas (check account balance)\n"+
			"  5. Transaction timeout (check transaction_timeout in config)",
			err)
	}
	l := L.With().Str("Transaction", signedTx.Hash().Hex()).Logger()
	l.Info().
		Int("FromKeyNum", fromKeyNum).
		Str("To", to).
		Interface("Value", value).
		Msg("Send ETH")
	_, err = m.WaitMined(ctx, l, m.Client, signedTx)
	if err != nil {
		return err
	}
	return err
}

// WaitMined the same as bind.WaitMined, awaits transaction receipt until timeout
func (m *Client) WaitMined(ctx context.Context, l zerolog.Logger, b bind.DeployBackend, tx *types.Transaction) (*types.Receipt, error) {
	l.Info().
		Msg("Waiting for transaction to be mined")
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	ctx, cancel := context.WithTimeout(ctx, m.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	for {
		receipt, err := b.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			l.Info().
				Int64("BlockNumber", receipt.BlockNumber.Int64()).
				Msg("Transaction receipt found")
			return receipt, nil
		} else if errors.Is(err, ethereum.NotFound) {
			l.Debug().
				Str("Timeout", m.Cfg.Network.TxnTimeout.String()).
				Msg("Awaiting transaction")
		} else {
			l.Debug().
				Msgf("Failed to get receipt due to: %s", err)
		}
		select {
		case <-ctx.Done():
			l.Error().Err(err).Str("Tx hash", tx.Hash().Hex()).Msg("Timed out, while waiting for transaction to be mined")
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

/* ClientOpts client functional options */

// ClientOpt is a client functional option
type ClientOpt func(c *Client)

// WithContractStore ContractStore functional option
func WithContractStore(as *ContractStore) ClientOpt {
	return func(c *Client) {
		c.ContractStore = as
	}
}

// WithContractMap contractAddressToNameMap functional option
func WithContractMap(contractAddressToNameMap ContractMap) ClientOpt {
	return func(c *Client) {
		c.ContractAddressToNameMap = contractAddressToNameMap
	}
}

// WithABIFinder ABIFinder functional option
func WithABIFinder(abiFinder *ABIFinder) ClientOpt {
	return func(c *Client) {
		c.ABIFinder = abiFinder
	}
}

// WithNonceManager NonceManager functional option
func WithNonceManager(nm *NonceManager) ClientOpt {
	return func(c *Client) {
		c.NonceManager = nm
	}
}

// WithTracer Tracer functional option
func WithTracer(t *Tracer) ClientOpt {
	return func(c *Client) {
		c.Tracer = t
	}
}

/* CallOpts function options */

// CallOpt is a functional option for bind.CallOpts
type CallOpt func(o *bind.CallOpts)

// WithPending sets pending option for bind.CallOpts
func WithPending(pending bool) CallOpt {
	return func(o *bind.CallOpts) {
		o.Pending = pending
	}
}

// WithBlockNumber sets blockNumber option for bind.CallOpts
func WithBlockNumber(bn uint64) CallOpt {
	return func(o *bind.CallOpts) {
		o.BlockNumber = big.NewInt(mustSafeInt64(bn))
	}
}

// NewCallOpts returns a new sequential call options wrapper
func (m *Client) NewCallOpts(o ...CallOpt) *bind.CallOpts {
	if errCallOpts := m.errCallOptsIfAddressCountTooLow(0); errCallOpts != nil {
		return errCallOpts
	}
	co := &bind.CallOpts{
		Pending: false,
		From:    m.Addresses[0],
	}
	for _, f := range o {
		f(co)
	}
	return co
}

// NewCallKeyOpts returns a new sequential call options wrapper from the key N
func (m *Client) NewCallKeyOpts(keyNum int, o ...CallOpt) *bind.CallOpts {
	if errCallOpts := m.errCallOptsIfAddressCountTooLow(keyNum); errCallOpts != nil {
		return errCallOpts
	}

	co := &bind.CallOpts{
		Pending: false,
		From:    m.Addresses[keyNum],
	}
	for _, f := range o {
		f(co)
	}
	return co
}

// errCallOptsIfAddressCountTooLow returns non-nil CallOpts with error in Context if keyNum is out of range
func (m *Client) errCallOptsIfAddressCountTooLow(keyNum int) *bind.CallOpts {
	if err := m.validateAddressesKeyNum(keyNum); err != nil {
		errText := err.Error()
		if keyNum == TimeoutKeyNum {
			errText += " (this is a probably because we didn't manage to find any synced key before timeout)"
		}

		err := errors.New(errText)
		m.Errors = append(m.Errors, err)
		opts := &bind.CallOpts{}

		// can't return nil, otherwise RPC wrapper will panic and we might lose funds on testnets/mainnets, that's why
		// error is passed in Context here to avoid panic, whoever is using Seth should make sure that there is no error
		// present in Context before using *bind.TransactOpts
		opts.Context = context.WithValue(context.Background(), ContextErrorKey{}, err)

		return opts
	}

	return nil
}

// errTxOptsIfPrivateKeysCountTooLow returns non-nil TransactOpts with error in Context if keyNum is out of range
func (m *Client) errTxOptsIfPrivateKeysCountTooLow(keyNum int) *bind.TransactOpts {
	if err := m.validatePrivateKeysKeyNum(keyNum); err != nil {
		errText := err.Error()
		if keyNum == TimeoutKeyNum {
			errText += " (this is a probably because we didn't manage to find any synced key before timeout)"
		}

		err := errors.New(errText)
		m.Errors = append(m.Errors, err)
		opts := &bind.TransactOpts{}

		// can't return nil, otherwise RPC wrapper will panic and we might lose funds on testnets/mainnets, that's why
		// error is passed in Context here to avoid panic, whoever is using Seth should make sure that there is no error
		// present in Context before using *bind.TransactOpts
		opts.Context = context.WithValue(context.Background(), ContextErrorKey{}, err)

		return opts
	}

	return nil
}

// TransactOpt is a wrapper for bind.TransactOpts
type TransactOpt func(o *bind.TransactOpts)

// WithValue sets value option for bind.TransactOpts
func WithValue(value *big.Int) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.Value = value
	}
}

// WithGasPrice sets gasPrice option for bind.TransactOpts
func WithGasPrice(gasPrice *big.Int) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.GasPrice = gasPrice
	}
}

// WithGasLimit sets gasLimit option for bind.TransactOpts
func WithGasLimit(gasLimit uint64) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.GasLimit = gasLimit
	}
}

// WithNoSend sets noSend option for bind.TransactOpts
func WithNoSend(noSend bool) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.NoSend = noSend
	}
}

// WithNonce sets nonce option for bind.TransactOpts
func WithNonce(nonce *big.Int) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.Nonce = nonce
	}
}

// WithGasFeeCap sets gasFeeCap option for bind.TransactOpts
func WithGasFeeCap(gasFeeCap *big.Int) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.GasFeeCap = gasFeeCap
	}
}

// WithGasTipCap sets gasTipCap option for bind.TransactOpts
func WithGasTipCap(gasTipCap *big.Int) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.GasTipCap = gasTipCap
	}
}

// WithSignerFn sets signerFn option for bind.TransactOpts
func WithSignerFn(signerFn bind.SignerFn) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.Signer = signerFn
	}
}

// WithFrom sets from option for bind.TransactOpts
func WithFrom(fromAddress common.Address) TransactOpt {
	return func(o *bind.TransactOpts) {
		o.From = fromAddress
	}
}

type ContextErrorKey struct{}

// NewTXOpts returns a new transaction options wrapper,
// Sets gas price/fee tip/cap and gas limit either based on TOML config or estimations.
func (m *Client) NewTXOpts(o ...TransactOpt) *bind.TransactOpts {
	opts, nonce, estimations := m.getProposedTransactionOptions(0)
	m.configureTransactionOpts(opts, nonce.PendingNonce, estimations, o...)
	L.Debug().
		Interface("Nonce", opts.Nonce).
		Interface("Value", opts.Value).
		Interface("GasPrice", opts.GasPrice).
		Interface("GasFeeCap", opts.GasFeeCap).
		Interface("GasTipCap", opts.GasTipCap).
		Uint64("GasLimit", opts.GasLimit).
		Msg("New transaction options")
	return opts
}

// NewTXKeyOpts returns a new transaction options wrapper,
// sets opts.GasPrice and opts.GasLimit from seth.toml or override with options
func (m *Client) NewTXKeyOpts(keyNum int, o ...TransactOpt) *bind.TransactOpts {
	if errTxOpts := m.errTxOptsIfPrivateKeysCountTooLow(keyNum); errTxOpts != nil {
		return errTxOpts
	}

	L.Debug().
		Interface("KeyNum", keyNum).
		Interface("Address", m.Addresses[keyNum]).
		Msg("Estimating transaction")
	opts, nonceStatus, estimations := m.getProposedTransactionOptions(keyNum)

	m.configureTransactionOpts(opts, nonceStatus.PendingNonce, estimations, o...)
	L.Debug().
		Interface("KeyNum", keyNum).
		Interface("Nonce", opts.Nonce).
		Interface("Value", opts.Value).
		Interface("GasPrice", opts.GasPrice).
		Interface("GasFeeCap", opts.GasFeeCap).
		Interface("GasTipCap", opts.GasTipCap).
		Uint64("GasLimit", opts.GasLimit).
		Msg("New transaction options")
	return opts
}

// AnySyncedKey returns the first synced key
func (m *Client) AnySyncedKey() int {
	return m.NonceManager.anySyncedKey()
}

type GasEstimations struct {
	// GasPrice for legacy transactions. If nil, RPC node will auto-estimate.
	GasPrice *big.Int
	// GasTipCap for EIP-1559 transactions. If nil, RPC node will auto-estimate.
	GasTipCap *big.Int
	// GasFeeCap for EIP-1559 transactions. If nil, RPC node will auto-estimate.
	GasFeeCap *big.Int
}

type NonceStatus struct {
	LastNonce    uint64
	PendingNonce uint64
}

func (m *Client) getNonceStatus(address common.Address) (NonceStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	pendingNonce, err := m.Client.PendingNonceAt(ctx, address)
	if err != nil {
		L.Error().Err(err).Msg("Failed to get pending nonce from RPC node")
		return NonceStatus{}, err
	}

	lastNonce, err := m.Client.NonceAt(ctx, address, nil)
	if err != nil {
		return NonceStatus{}, err
	}

	return NonceStatus{
		LastNonce:    lastNonce,
		PendingNonce: pendingNonce,
	}, nil
}

// getProposedTransactionOptions gets all the tx info that network proposed
func (m *Client) getProposedTransactionOptions(keyNum int) (*bind.TransactOpts, NonceStatus, GasEstimations) {
	if errTxOpts := m.errTxOptsIfPrivateKeysCountTooLow(keyNum); errTxOpts != nil {
		return errTxOpts, NonceStatus{}, GasEstimations{}
	}

	nonceStatus, err := m.getNonceStatus(m.Addresses[keyNum])
	if err != nil {
		m.Errors = append(m.Errors, err)
		// can't return nil, otherwise RPC wrapper will panic
		ctx := context.WithValue(context.Background(), ContextErrorKey{}, err)

		return &bind.TransactOpts{Context: ctx}, NonceStatus{}, GasEstimations{}
	}

	var ctx context.Context

	if m.Cfg.PendingNonceProtectionEnabled {
		if pendingErr := m.WaitUntilNoPendingTxForKeyNum(keyNum, m.Cfg.PendingNonceProtectionTimeout.Duration()); pendingErr != nil {
			errMsg := `
pending nonce for key %d is higher than last nonce, there are %d pending transactions.

This issue is caused by one of two things:
1. You are using the same keyNum in multiple goroutines, which is not supported. Each goroutine should use an unique keyNum
2. You have stuck transaction(s). Speed them up by sending replacement transactions with higher gas price before continuing, otherwise future transactions most probably will also get stuck`
			err := fmt.Errorf(errMsg, keyNum, nonceStatus.PendingNonce-nonceStatus.LastNonce)
			m.Errors = append(m.Errors, err)
			// can't return nil, otherwise RPC wrapper will panic, and we might lose funds on testnets/mainnets, that's why
			// error is passed in Context here to avoid panic, whoever is using Seth should make sure that there is no error
			// present in Context before using *bind.TransactOpts
			ctx = context.WithValue(context.Background(), ContextErrorKey{}, err)
		}
		L.Debug().
			Msg("Pending nonce protection is enabled. Nonce status is OK")
	}

	estimations := m.CalculateGasEstimations(m.NewDefaultGasEstimationRequest())

	L.Debug().
		Interface("KeyNum", keyNum).
		Uint64("Nonce", nonceStatus.PendingNonce).
		Interface("GasEstimations", estimations).
		Msg("Proposed transaction options")

	opts, err := bind.NewKeyedTransactorWithChainID(m.PrivateKeys[keyNum], big.NewInt(m.ChainID))
	if err != nil {
		err = fmt.Errorf("failed to create transactor for key #%d (address: %s, chain ID: %d): %w\n"+
			"This usually indicates:\n"+
			"  1. Invalid private key format\n"+
			"  2. Chain ID mismatch\n"+
			"  3. Key not properly loaded\n"+
			"Verify the private key is valid and chain_id is correct in config",
			keyNum, m.Addresses[keyNum].Hex(), m.ChainID, err)
		m.Errors = append(m.Errors, err)
		// can't return nil, otherwise RPC wrapper will panic and we might lose funds on testnets/mainnets, that's why
		// error is passed in Context here to avoid panic, whoever is using Seth should make sure that there is no error
		// present in Context before using *bind.TransactOpts
		ctx := context.WithValue(context.Background(), ContextErrorKey{}, err)

		return &bind.TransactOpts{Context: ctx}, NonceStatus{}, GasEstimations{}
	}

	if ctx != nil {
		opts.Context = ctx
	}

	return opts, nonceStatus, estimations
}

type GasEstimationRequest struct {
	GasEstimationEnabled bool
	FallbackGasPrice     int64
	FallbackGasFeeCap    int64
	FallbackGasTipCap    int64
	Priority             string
}

// NewDefaultGasEstimationRequest creates a new default gas estimation request based on current network configuration
func (m *Client) NewDefaultGasEstimationRequest() GasEstimationRequest {
	return GasEstimationRequest{
		GasEstimationEnabled: m.Cfg.Network.GasPriceEstimationEnabled,
		FallbackGasPrice:     m.Cfg.Network.GasPrice,
		FallbackGasFeeCap:    m.Cfg.Network.GasFeeCap,
		FallbackGasTipCap:    m.Cfg.Network.GasTipCap,
		Priority:             m.Cfg.Network.GasPriceEstimationTxPriority,
	}
}

// CalculateGasEstimations calculates gas estimations (price, tip/cap) or uses hardcoded values if estimation is disabled,
// estimation errors or network is a simulated one.
func (m *Client) CalculateGasEstimations(request GasEstimationRequest) GasEstimations {
	estimations := GasEstimations{}

	if request.Priority == Priority_Auto {
		// Return empty estimations with nil gas prices and fees.
		// This signals to the RPC node to auto-estimate gas prices.
		// The nil values will be passed to bind.TransactOpts, which treats
		// nil gas fields as a request for automatic estimation.
		L.Debug().Msg("Auto priority selected, skipping gas estimations")
		return estimations
	}

	if m.Cfg.IsSimulatedNetwork() || !request.GasEstimationEnabled {
		estimations.GasPrice = big.NewInt(request.FallbackGasPrice)
		estimations.GasFeeCap = big.NewInt(request.FallbackGasFeeCap)
		estimations.GasTipCap = big.NewInt(request.FallbackGasTipCap)

		return estimations
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Network.TxnTimeout.Duration())
	defer cancel()

	var disableEstimationsIfNeeded = func(err error) {
		if errors.Is(err, GasEstimationErr) {
			L.Warn().Msg("Received incorrect gas estimations. Disabling them and reverting to hardcoded values. Remember to update your config!")
			m.Cfg.Network.GasPriceEstimationEnabled = false
		}
	}

	var calculateLegacyFees = func() {
		gasPrice, err := m.GetSuggestedLegacyFees(ctx, request.Priority)
		if err != nil {
			disableEstimationsIfNeeded(err)
			L.Debug().Err(err).Msg("Failed to get suggested Legacy fees. Using hardcoded values")
			estimations.GasPrice = big.NewInt(request.FallbackGasPrice)
		} else {
			estimations.GasPrice = gasPrice
		}
	}

	if m.Cfg.Network.EIP1559DynamicFees {
		maxFee, priorityFee, err := m.GetSuggestedEIP1559Fees(ctx, request.Priority)
		if err != nil {
			L.Debug().Err(err).Msg("Failed to get suggested EIP1559 fees. Using hardcoded values")
			estimations.GasFeeCap = big.NewInt(request.FallbackGasFeeCap)
			estimations.GasTipCap = big.NewInt(request.FallbackGasTipCap)

			disableEstimationsIfNeeded(err)

			if strings.Contains(err.Error(), "method eth_maxPriorityFeePerGas") || strings.Contains(err.Error(), "method eth_maxFeePerGas") || strings.Contains(err.Error(), "method eth_feeHistory") || strings.Contains(err.Error(), "expected input list for types.txdata") {
				L.Warn().Msg("EIP1559 fees are not supported by the network. Switching to Legacy fees. Remember to update your config!")
				if m.Cfg.Network.GasPrice == 0 {
					L.Warn().Msg("Gas price is 0. If Legacy estimations fail, there will no fallback price and transactions will start fail. Set gas price in config and disable EIP1559DynamicFees")
				}
				m.Cfg.Network.EIP1559DynamicFees = false
				calculateLegacyFees()
			}
		} else {
			estimations.GasFeeCap = maxFee
			estimations.GasTipCap = priorityFee
		}
	} else {
		calculateLegacyFees()
	}

	return estimations
}

// EstimateGasLimitForFundTransfer estimates gas limit for fund transfer
func (m *Client) EstimateGasLimitForFundTransfer(from, to common.Address, amount *big.Int) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	gasLimit, err := m.Client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: amount,
	})
	if err != nil {
		L.Debug().Msgf("Failed to estimate gas for fund transfer due to: %s", err.Error())
		return 0, fmt.Errorf("failed to estimate gas for fund transfer from %s to %s (amount: %s wei): %w\n"+
			"Possible causes:\n"+
			"  1. Insufficient balance in sender account\n"+
			"  2. Invalid recipient address\n"+
			"  3. RPC node doesn't support gas estimation\n"+
			"  4. Network congestion or RPC issues\n"+
			"Try setting an explicit gas_limit in config if estimation continues to fail",
			from.Hex(), to.Hex(), amount.String(), err)
	}
	return gasLimit, nil
}

// configureTransactionOpts configures transaction for legacy or type-2
func (m *Client) configureTransactionOpts(
	opts *bind.TransactOpts,
	nonce uint64,
	estimations GasEstimations,
	o ...TransactOpt,
) *bind.TransactOpts {
	opts.Nonce = big.NewInt(mustSafeInt64(nonce))
	opts.GasPrice = estimations.GasPrice
	opts.GasLimit = m.Cfg.Network.GasLimit

	if m.Cfg.Network.EIP1559DynamicFees {
		opts.GasPrice = nil
		opts.GasTipCap = estimations.GasTipCap
		opts.GasFeeCap = estimations.GasFeeCap

		// Log when using auto-estimated gas (all nil values)
		if opts.GasTipCap == nil && opts.GasFeeCap == nil {
			L.Debug().Msg("Using RPC node's automatic gas estimation (EIP-1559)")
		}
	} else if opts.GasPrice == nil {
		L.Debug().Msg("Using RPC node's automatic gas estimation (Legacy)")
	}

	for _, f := range o {
		f(opts)
	}
	return opts
}

// ContractLoader is a helper struct for loading contracts
type ContractLoader[T any] struct {
	Client *Client
}

// NewContractLoader creates a new contract loader
func NewContractLoader[T any](client *Client) *ContractLoader[T] {
	return &ContractLoader[T]{
		Client: client,
	}
}

// LoadContract loads contract by name, address, ABI loader and wrapper init function, it adds contract ABI to Seth Contract Store and address to Contract Map. Thanks to that we can easily
// trace and debug interactions with the contract. Signatures of functions passed to this method were chosen to conform to Geth wrappers' GetAbi() and NewXXXContract() functions.
func (cl *ContractLoader[T]) LoadContract(name string, address common.Address, abiLoadFn func() (*abi.ABI, error), wrapperInitFn func(common.Address, bind.ContractBackend) (*T, error)) (*T, error) {
	abiData, err := abiLoadFn()
	if err != nil {
		return new(T), err
	}
	cl.Client.ContractStore.AddABI(name, *abiData)
	cl.Client.ContractAddressToNameMap.AddContract(address.Hex(), name)

	return wrapperInitFn(address, cl.Client.Client)
}

// DeployContract deploys contract using ABI and bytecode passed to it, waits for transaction to be minted and contract really
// available at the address, so that when the method returns it's safe to interact with it. It also saves the contract address and ABI name
// to the contract map, so that we can use that, when tracing transactions. It is suggested to use name identical to the name of the contract Solidity file.
func (m *Client) DeployContract(auth *bind.TransactOpts, name string, abi abi.ABI, bytecode []byte, params ...interface{}) (DeploymentData, error) {
	L.Info().
		Msgf("Started deploying %s contract", name)

	if auth.Context != nil {
		if err, ok := auth.Context.Value(ContextErrorKey{}).(error); ok {
			return DeploymentData{}, fmt.Errorf("aborted contract deployment for '%s': context error: %w\n"+
				"This usually means there was an error creating transaction options. "+
				"Check the error above for details about what went wrong",
				name, err)
		}
	}

	if m.Cfg.Hooks != nil && m.Cfg.Hooks.ContractDeployment.Pre != nil {
		if err := m.Cfg.Hooks.ContractDeployment.Pre(auth, name, abi, bytecode, params...); err != nil {
			return DeploymentData{}, fmt.Errorf("contract deployment pre-hook failed for '%s': %w\n"+
				"The pre-deployment hook returned an error. "+
				"Check your hook implementation for issues",
				name, err)
		}
	} else {
		L.Trace().Msg("No pre-contract deployment hook defined. Skipping")
	}

	address, tx, contract, err := bind.DeployContract(auth, abi, bytecode, m.Client, params...)
	if err != nil {
		return DeploymentData{}, wrapErrInMessageWithASuggestion(err)
	}

	L.Info().
		Str("Address", address.Hex()).
		Str("TXHash", tx.Hash().Hex()).
		Msgf("Waiting for %s contract deployment to finish", name)

	m.ContractAddressToNameMap.AddContract(address.Hex(), name)

	if _, ok := m.ContractStore.GetABI(name); !ok {
		m.ContractStore.AddABI(name, abi)
	}

	if m.Cfg.Hooks != nil && m.Cfg.Hooks.ContractDeployment.Post != nil {
		if err := m.Cfg.Hooks.ContractDeployment.Post(m, tx); err != nil {
			return DeploymentData{}, fmt.Errorf("contract deployment post-hook failed for transaction %s: %w\n"+
				"The post-deployment hook returned an error. "+
				"Check your hook implementation for issues",
				tx.Hash().Hex(), err)
		}
	} else {
		L.Trace().Msg("No post-contract deployment hook defined. Skipping")
	}

	if err := retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Network.TxnTimeout.Duration())
			_, err := bind.WaitDeployed(ctx, m.Client, tx)
			cancel()

			// let's make sure that deployment transaction was successful, before retrying
			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Network.TxnTimeout.Duration())
				receipt, mineErr := bind.WaitMined(ctx, m.Client, tx)
				if mineErr != nil {
					cancel()
					return mineErr
				}
				cancel()

				if receipt.Status == 0 {
					return fmt.Errorf("contract '%s' deployment transaction was reverted. "+
						"Transaction hash: %s\n"+
						"Common causes:\n"+
						"  1. Constructor parameters are incorrect\n"+
						"  2. Insufficient gas limit\n"+
						"  3. Constructor validation/require failed\n"+
						"  4. Contract bytecode is invalid\n"+
						"Check transaction trace with decode for the specific revert reason",
						name, tx.Hash().Hex())
				}
			}

			return err
		}, retry.OnRetry(func(i uint, retryErr error) {
			switch {
			case errors.Is(retryErr, context.DeadlineExceeded):
				replacementTx, replacementErr := prepareReplacementTransaction(m, tx)
				if replacementErr != nil {
					L.Debug().Str("Current error", retryErr.Error()).Str("Replacement error", replacementErr.Error()).Uint("Attempt", i+1).Msg("Failed to prepare replacement transaction for contract deployment. Retrying with the original one")
					return
				}
				tx = replacementTx
			default:
				// do nothing, just wait again until it's mined
			}
			L.Debug().Str("Current error", retryErr.Error()).Uint("Attempt", i+1).Msg("Waiting for contract to be deployed")
		}),
		retry.DelayType(retry.FixedDelay),
		// if gas bump retries are set to 0, we still want to retry 10 times, because what we will be retrying will be other errors (no code at address, etc.)
		// downside is that if retries are enabled and their number is low other retry errors will be retried only that number of times
		// (we could have custom logic for different retry count per error, but that seemed like an overkill, so it wasn't implemented)
		retry.Attempts(func() uint {
			if m.Cfg.GasBumpRetries() != 0 {
				return m.Cfg.GasBumpRetries()
			}
			return 10
		}()),
		retry.RetryIf(func(err error) bool {
			return strings.Contains(strings.ToLower(err.Error()), "no contract code at given address") ||
				strings.Contains(strings.ToLower(err.Error()), "no contract code after deployment") ||
				(m.Cfg.GasBumpRetries() != 0 && errors.Is(err, context.DeadlineExceeded))
		}),
	); err != nil {
		// pass this specific error, so that Decode knows that it's not the actual revert reason
		_, _ = m.Decode(tx, errors.New(ErrContractDeploymentFailed))

		return DeploymentData{}, wrapErrInMessageWithASuggestion(m.rewriteDeploymentError(err))
	}

	L.Info().
		Str("Address", address.Hex()).
		Str("TXHash", tx.Hash().Hex()).
		Msgf("Deployed %s contract", name)

	if !m.Cfg.ShouldSaveDeployedContractMap() {
		return DeploymentData{Address: address, Transaction: tx, BoundContract: contract}, nil
	}

	if err := SaveDeployedContract(m.Cfg.ContractMapFile, name, address.Hex()); err != nil {
		L.Warn().
			Err(err).
			Msg("Failed to save deployed contract address to file")
	}

	return DeploymentData{Address: address, Transaction: tx, BoundContract: contract}, nil
}

// rewriteDeploymentError makes some known errors more human friendly
func (m *Client) rewriteDeploymentError(err error) error {
	var maybeRetryErr retry.Error
	switch {
	case errors.As(err, &maybeRetryErr):
		areAllTimeouts := false
		for _, e := range maybeRetryErr.WrappedErrors() {
			if !errors.Is(e, context.DeadlineExceeded) {
				break
			}
			areAllTimeouts = true
		}

		if areAllTimeouts {
			newErr := retry.Error{}
			for range maybeRetryErr.WrappedErrors() {
				newErr = append(newErr, fmt.Errorf("deployment transaction was not mined within %s", m.Cfg.Network.TxnTimeout.Duration().String()))
			}
			err = newErr
		}
	case errors.Is(err, context.DeadlineExceeded):
		{
			err = fmt.Errorf("deployment transaction was not mined within %s", m.Cfg.Network.TxnTimeout.Duration().String())
		}
	}

	return err
}

type DeploymentData struct {
	Address       common.Address
	Transaction   *types.Transaction
	BoundContract *bind.BoundContract
}

// DeployContractFromContractStore deploys contract from Seth's Contract Store, waits for transaction to be minted and contract really
// available at the address, so that when the method returns it's safe to interact with it. It also saves the contract address and ABI name
// to the contract map, so that we can use that, when tracing transactions. Name by which you refer the contract should be the same as the
// name of ABI file (you can omit the .abi suffix).
func (m *Client) DeployContractFromContractStore(auth *bind.TransactOpts, name string, params ...interface{}) (DeploymentData, error) {
	if m.ContractStore == nil {
		return DeploymentData{}, fmt.Errorf("contract store is nil. Cannot deploy contract from store.\n" +
			"This usually means:\n" +
			"  1. Seth client wasn't properly initialized\n" +
			"  2. ABI directory path is incorrect in config\n" +
			"Solutions:\n" +
			"  1. Set 'abi_dir' and 'bin_dir' in seth.toml config file\n" +
			"  2. Use ClientBuilder: builder.WithABIDir(path).WithBINDir(path)\n" +
			"  3. Use DeployContract() method with explicit ABI/bytecode instead")
	}

	name = strings.TrimSuffix(name, ".abi")
	name = strings.TrimSuffix(name, ".bin")

	contractAbi, ok := m.ContractStore.ABIs[name+".abi"]
	if !ok {
		return DeploymentData{}, fmt.Errorf("ABI for contract '%s' not found in contract store.\n"+
			"Total ABIs loaded: %d\n"+
			"Ensure the ABI file '%s.abi' exists in the directory specified by 'abi_dir' in your config",
			name, len(m.ContractStore.ABIs), name)
	}

	bytecode, ok := m.ContractStore.BINs[name+".bin"]
	if !ok {
		return DeploymentData{}, fmt.Errorf("bytecode (BIN) for contract '%s' not found in contract store.\n"+
			"Total BINs loaded: %d\n"+
			"Ensure the BIN file '%s.bin' exists in the directory specified by 'bin_dir' in your config",
			name, len(m.ContractStore.BINs), name)
	}

	data, err := m.DeployContract(auth, name, contractAbi, bytecode, params...)
	if err != nil {
		return DeploymentData{}, err
	}

	return data, nil
}

func (m *Client) SaveDecodedCallsAsJson(dirname string) error {
	return m.Tracer.SaveDecodedCallsAsJson(dirname)
}

type TransactionLog struct {
	Topics []common.Hash
	Data   []byte
}

func (t TransactionLog) GetTopics() []common.Hash {
	return t.Topics
}

func (t TransactionLog) GetData() []byte {
	return t.Data
}

func (m *Client) decodeContractLogs(l zerolog.Logger, logs []types.Log, allABIs []*abi.ABI) ([]DecodedTransactionLog, []EventDecodingError, error) {
	l.Trace().
		Msg("Decoding events")
	sigMap := buildEventSignatureMap(allABIs)

	var eventsParsed []DecodedTransactionLog
	var decodeErrors []EventDecodingError
	for _, lo := range logs {
		if len(lo.Topics) == 0 {
			l.Debug().
				Msg("Log has no topics; skipping")
			continue
		}

		eventSig := lo.Topics[0].Hex()
		possibleEvents, exists := sigMap[eventSig]
		if !exists {
			l.Trace().
				Str("Event signature", eventSig).
				Msg("No matching events found for signature")
			continue
		}

		// Check if we know what contract is this log from and if we do, get its ABI to skip unnecessary iterations
		var knownContractABI *abi.ABI
		if contractName := m.ContractAddressToNameMap.GetContractName(lo.Address.Hex()); contractName != "" {
			maybeABI, ok := m.ContractStore.GetABI(contractName)
			if !ok {
				l.Trace().
					Str("Event signature", eventSig).
					Str("Contract name", contractName).
					Str("Contract address", lo.Address.Hex()).
					Msg("No ABI found for known contract; this is unexpected. Continuing with step-by-step ABI search")
			} else {
				knownContractABI = maybeABI
			}
		}

		// Iterate over possible events with the same signature
		matched := false
		var decodeAttempts []ABIDecodingError
		for _, evWithABI := range possibleEvents {
			evSpec := evWithABI.EventSpec
			contractABI := evWithABI.ContractABI

			// Check if known contract ABI matches candidate ABI and if not, skip this ABI and try the next one
			if knownContractABI != nil && !reflect.DeepEqual(knownContractABI, contractABI) {
				l.Trace().
					Str("Event signature", eventSig).
					Str("Contract address", lo.Address.Hex()).
					Msg("ABI doesn't match known ABI for this address; trying next ABI")
				continue
			}

			// Validate indexed parameters count
			// Non-indexed parameters are stored in the Data field,
			// and much harder to validate due to dynamic types,
			// so we skip them for now
			var indexedParams abi.Arguments
			for _, input := range evSpec.Inputs {
				if input.Indexed {
					indexedParams = append(indexedParams, input)
				}
			}

			expectedIndexed := len(indexedParams)
			actualIndexed := len(lo.Topics) - 1 // First topic is the event signature

			if expectedIndexed != actualIndexed {
				l.Trace().
					Str("Event", evSpec.Name).
					Int("Expected indexed param count", expectedIndexed).
					Int("Actual indexed param count", actualIndexed).
					Msg("Mismatch in indexed parameters; skipping event")
				continue
			}

			// Proceed to decode the event
			// Find ABI name for this contract ABI
			abiName := m.findABIName(contractABI)
			d := TransactionLog{lo.Topics, lo.Data}
			l.Trace().
				Str("Name", evSpec.RawName).
				Str("Signature", evSpec.Sig).
				Str("ABI", abiName).
				Msg("Unpacking event")

			eventsMap, topicsMap, err := decodeEventFromLog(l, *contractABI, *evSpec, d)
			if err != nil {
				l.Debug().
					Err(err).
					Str("Event", evSpec.Name).
					Str("ABI", abiName).
					Msg("Failed to decode event; trying next ABI")
				decodeAttempts = append(decodeAttempts, ABIDecodingError{
					ABIName:   abiName,
					EventName: evSpec.Name,
					Error:     err.Error(),
				})
				continue // Try next ABI instead of giving up
			}

			parsedEvent := decodedLogFromMaps(&DecodedTransactionLog{}, eventsMap, topicsMap)
			decodedTransactionLog, ok := parsedEvent.(*DecodedTransactionLog)
			if ok {
				decodedTransactionLog.Signature = evSpec.Sig
				m.mergeLogMeta(decodedTransactionLog, lo)
				eventsParsed = append(eventsParsed, *decodedTransactionLog)
				l.Trace().
					Interface("Log", parsedEvent).
					Msg("Transaction log decoded successfully")
				matched = true
				break // Move to the next log after successful decoding
			}

			l.Trace().
				Str("Actual type", fmt.Sprintf("%T", decodedTransactionLog)).
				Msg("Failed to cast decoded event to DecodedTransactionLog")
		}

		if !matched {
			// Record the decode failure
			topics := make([]string, len(lo.Topics))
			for i, topic := range lo.Topics {
				topics[i] = topic.Hex()
			}

			decodeError := EventDecodingError{
				Signature: eventSig,
				LogIndex:  lo.Index,
				Address:   lo.Address.Hex(),
				Topics:    topics,
				Errors:    decodeAttempts,
			}
			decodeErrors = append(decodeErrors, decodeError)

			abiNames := make([]string, len(decodeAttempts))
			for i, attempt := range decodeAttempts {
				abiNames[i] = attempt.ABIName
			}

			l.Warn().
				Str("Signature", eventSig).
				Uint("LogIndex", lo.Index).
				Str("Address", lo.Address.Hex()).
				Strs("AttemptedABIs", abiNames).
				Int("FailedAttempts", len(decodeAttempts)).
				Msg("Failed to decode event log")
		}
	}
	return eventsParsed, decodeErrors, nil
}

type eventWithABI struct {
	ContractABI *abi.ABI
	EventSpec   *abi.Event
}

// buildEventSignatureMap precomputes a mapping from event signature to events with their ABIs
func buildEventSignatureMap(allABIs []*abi.ABI) map[string][]*eventWithABI {
	sigMap := make(map[string][]*eventWithABI)
	for _, a := range allABIs {
		for _, ev := range a.Events {
			event := ev //nolint:copyloopvar // Explicitly keeping the copy for clarity
			sigMap[ev.ID.Hex()] = append(sigMap[ev.ID.Hex()], &eventWithABI{
				ContractABI: a,
				EventSpec:   &event,
			})
		}
	}

	return sigMap
}

// findABIName finds the name of the ABI in the ContractStore, returns "unknown" if not found
func (m *Client) findABIName(targetABI *abi.ABI) string {
	if m.ContractStore == nil {
		return "unknown"
	}

	for name, storedABI := range m.ContractStore.ABIs {
		if reflect.DeepEqual(storedABI, *targetABI) {
			return strings.TrimSuffix(name, ".abi")
		}
	}

	return "unknown"
}

// WaitUntilNoPendingTxForRootKey waits until there's no pending transaction for root key. If after timeout there are still pending transactions, it returns error.
func (m *Client) WaitUntilNoPendingTxForRootKey(timeout time.Duration) error {
	return m.WaitUntilNoPendingTx(m.MustGetRootKeyAddress(), timeout)
}

// WaitUntilNoPendingTxForKeyNum waits until there's no pending transaction for key at index `keyNum`. If index is out of range or
// if after timeout there are still pending transactions, it returns error.
func (m *Client) WaitUntilNoPendingTxForKeyNum(keyNum int, timeout time.Duration) error {
	if err := m.validateAddressesKeyNum(keyNum); err != nil {
		return err
	}
	return m.WaitUntilNoPendingTx(m.Addresses[keyNum], timeout)
}

// WaitUntilNoPendingTx waits until there's no pending transaction for address. If after timeout there are still pending transactions, it returns error.
func (m *Client) WaitUntilNoPendingTx(address common.Address, timeout time.Duration) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	waitTimeout := time.NewTimer(timeout)
	defer waitTimeout.Stop()

	for {
		select {
		case <-waitTimeout.C:
			return fmt.Errorf("timeout after %s: address %s still has pending transactions. "+
				"This means transactions haven't been mined within the timeout period.\n"+
				"Troubleshooting:\n"+
				"  1. Check if the network is processing transactions (check block explorer)\n"+
				"  2. Gas price might be too low (increase gas_price or gas_fee_cap)\n"+
				"  3. Network congestion (wait longer or increase timeout)\n"+
				"  4. Enable gas bumping: set gas_bump.retries > 0 in config",
				timeout, address.Hex())
		case <-ticker.C:
			nonceStatus, err := m.getNonceStatus(address)
			// if there is an error, we can't be sure if there are pending transactions or not, let's retry on next tick
			if err != nil {
				L.Debug().Err(err).Msg("Failed to get nonce status")
				continue
			}
			L.Debug().Msgf("Nonce status for address %s: %v", address.Hex(), nonceStatus)

			if nonceStatus.PendingNonce > nonceStatus.LastNonce {
				L.Debug().Uint64("Pending transactions", nonceStatus.PendingNonce-nonceStatus.LastNonce).Msgf("There are still pending transactions for %s", address.Hex())
				continue
			}

			return nil
		}
	}
}

func (m *Client) validatePrivateKeysKeyNum(keyNum int) error {
	if keyNum >= len(m.PrivateKeys) || keyNum < 0 {
		if len(m.PrivateKeys) == 0 {
			return fmt.Errorf("no private keys loaded, but tried to use key #%d.\n"+
				"Load private keys by:\n"+
				"  1. Setting SETH_ROOT_PRIVATE_KEY environment variable\n"+
				"  2. Adding 'private_keys_secret' to your seth.toml network config\n"+
				"  3. Using WithPrivateKeys() with ClientBuilder",
				keyNum)
		}
		return fmt.Errorf("keyNum %d is out of range. Available keys: 0-%d (total: %d keys loaded).\n"+
			"Common causes:\n"+
			"  1. Using keyNum from another test/context\n"+
			"  2. Not enough keys configured for parallel test execution\n"+
			"  3. Consider enabling ephemeral_addresses_number in your config for more keys",
			keyNum, len(m.PrivateKeys)-1, len(m.PrivateKeys))
	}

	return nil
}

func (m *Client) validateAddressesKeyNum(keyNum int) error {
	if keyNum >= len(m.Addresses) || keyNum < 0 {
		if len(m.Addresses) == 0 {
			return fmt.Errorf("no addresses loaded, but tried to use key #%d.\n"+
				"This should not happen if private keys were loaded correctly. "+
				"Please report this issue at https://github.com/smartcontractkit/chainlink-testing-framework/issues with the stack trace",
				keyNum)
		}
		return fmt.Errorf("keyNum %d is out of range. Available addresses: 0-%d (total: %d addresses loaded).\n"+
			"This indicates the keyNum is invalid for the current client configuration",
			keyNum, len(m.Addresses)-1, len(m.Addresses))
	}

	return nil
}

// mergeLogMeta add metadata from log
func (m *Client) mergeLogMeta(pe *DecodedTransactionLog, l types.Log) {
	pe.Address = l.Address
	pe.Topics = make([]string, 0)
	for _, t := range l.Topics {
		pe.Topics = append(pe.Topics, t.String())
	}
	pe.BlockNumber = l.BlockNumber
	pe.Index = l.Index
	pe.TXHash = l.TxHash.Hex()
	pe.TXIndex = l.TxIndex
	pe.Removed = l.Removed
}

func shouldInitializeTracer(client simulated.Client, cfg *Config) bool {
	return len(cfg.Network.URLs) > 0 && supportsTracing(client)
}

func supportsTracing(client simulated.Client) bool {
	return strings.Contains(reflect.TypeOf(client).String(), "ethclient.Client")
}
