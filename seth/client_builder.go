package seth

import (
	"time"
)

type ClientBuilder struct {
	config *Config
}

// NewClientBuilder creates a new ClientBuilder with reasonable default values. You only need to pass private key(s) and RPC URL to build a usable config.
func NewClientBuilder() *ClientBuilder {
	network := &Network{
		Name:                         DefaultNetworkName,
		EIP1559DynamicFees:           true,
		TxnTimeout:                   MustMakeDuration(5 * time.Minute),
		DialTimeout:                  MustMakeDuration(DefaultDialTimeout),
		TransferGasFee:               DefaultTransferGasFee,
		GasPriceEstimationEnabled:    true,
		GasPriceEstimationBlocks:     200,
		GasPriceEstimationTxPriority: Priority_Standard,
		GasPrice:                     DefaultGasPrice,
		GasFeeCap:                    DefaultGasFeeCap,
		GasTipCap:                    DefaultGasTipCap,
	}

	return &ClientBuilder{
		config: &Config{
			ArtifactsDir:          "seth_artifacts",
			EphemeralAddrs:        &ZeroInt64,
			RootKeyFundsBuffer:    &ZeroInt64,
			Network:               network,
			Networks:              []*Network{network},
			TracingLevel:          TracingLevel_Reverted,
			TraceOutputs:          []string{TraceOutput_Console, TraceOutput_DOT},
			ExperimentsEnabled:    []string{},
			CheckRpcHealthOnStart: true,
			BlockStatsConfig:      &BlockStatsConfig{RPCRateLimit: 10},
			NonceManager:          &NonceManagerCfg{KeySyncRateLimitSec: 10, KeySyncRetries: 3, KeySyncTimeout: MustMakeDuration(60 * time.Second), KeySyncRetryDelay: MustMakeDuration(5 * time.Second)},
			GasBump: &GasBumpConfig{
				Retries: 10,
			},
		},
	}
}

// WithRpcUrl sets the RPC URL for the config.
// Default value is an empty string (which is an incorrect value).
func (c *ClientBuilder) WithRpcUrl(url string) *ClientBuilder {
	c.config.Network.URLs = []string{url}
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].URLs = []string{url}
	}
	return c
}

// WithPrivateKeys sets the private keys for the config. At least one is required to build a valid config.
// Default value is an empty slice (which is an incorrect value).
func (c *ClientBuilder) WithPrivateKeys(pks []string) *ClientBuilder {
	c.config.Network.PrivateKeys = pks
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].PrivateKeys = pks
	}
	return c
}

// WithNetworkName sets the network name, useful mostly for debugging and logging.
// Default value is "default".
func (c *ClientBuilder) WithNetworkName(name string) *ClientBuilder {
	c.config.Network.Name = name
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].Name = name
	}
	return c
}

// WithGasPriceEstimations enables or disables gas price estimations, sets the number of blocks to use for estimation or transaction priority.
// Even with estimations enabled you should still either set legacy gas price with `WithLegacyGasPrice()` or EIP-1559 dynamic fees with `WithDynamicGasPrices()`
// ss they will be used as fallback values, if the estimations fail.
// Following priorities are supported: "slow", "standard" and "fast"
// Default values are true for enabled, 200 blocks for estimation and "standard" for priority.
func (c *ClientBuilder) WithGasPriceEstimations(enabled bool, estimationBlocks uint64, txPriority string) *ClientBuilder {
	c.config.Network.GasPriceEstimationEnabled = enabled
	c.config.Network.GasPriceEstimationBlocks = estimationBlocks
	c.config.Network.GasPriceEstimationTxPriority = txPriority
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].GasPriceEstimationEnabled = enabled
		c.config.Networks[0].GasPriceEstimationBlocks = estimationBlocks
		c.config.Networks[0].GasPriceEstimationTxPriority = txPriority
	}
	return c
}

// WithEIP1559DynamicFees enables or disables EIP-1559 dynamic fees. If enabled, you should set gas fee cap and gas tip cap with `WithDynamicGasPrices()`
// Default value is true.
func (c *ClientBuilder) WithEIP1559DynamicFees(enabled bool) *ClientBuilder {
	c.config.Network.EIP1559DynamicFees = enabled
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].EIP1559DynamicFees = enabled
	}
	return c
}

// WithLegacyGasPrice sets the gas price for legacy transactions that will be used only if EIP-1559 dynamic fees are disabled.
// Default value is 1 gwei.
func (c *ClientBuilder) WithLegacyGasPrice(gasPrice int64) *ClientBuilder {
	c.config.Network.GasPrice = gasPrice
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].GasPrice = gasPrice
	}
	return c
}

// WithDynamicGasPrices sets the gas fee cap and gas tip cap for EIP-1559 dynamic fees. These values will be used only if EIP-1559 dynamic fees are enabled.
// Default values are 150 gwei for gas fee cap and 50 gwei for gas tip cap.
func (c *ClientBuilder) WithDynamicGasPrices(gasFeeCap, gasTipCap int64) *ClientBuilder {
	c.config.Network.GasFeeCap = gasFeeCap
	c.config.Network.GasTipCap = gasTipCap
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].GasFeeCap = gasFeeCap
		c.config.Networks[0].GasTipCap = gasTipCap
	}
	return c
}

// WithTransferGasFee sets the gas fee for transfer transactions. This value is used, when sending funds to ephemeral keys or returning funds to root private key.
// Default value is 21_000 wei.
func (c *ClientBuilder) WithTransferGasFee(gasFee int64) *ClientBuilder {
	c.config.Network.TransferGasFee = gasFee
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].TransferGasFee = gasFee
	}
	return c
}

// WithGasBumping sets the number of retries for gas bumping and max gas price. You can also provide a custom bumping strategy. If the transaction is not mined within this number of retries, it will be considered failed.
// If the gas price is bumped to a value higher than max gas price, no more gas bumping will be attempted and previous gas price will be used by all subsequent attempts. If set to 0 max price is not checked.
// Default value is 10 retries, no max gas price and a default bumping strategy (with gas increase % based on gas_price_estimation_tx_priority)
func (c *ClientBuilder) WithGasBumping(retries uint, maxGasPrice int64, customBumpingStrategy GasBumpStrategyFn) *ClientBuilder {
	c.config.GasBump = &GasBumpConfig{
		Retries:     retries,
		MaxGasPrice: maxGasPrice,
		StrategyFn:  customBumpingStrategy,
	}
	return c
}

// WithTransactionTimeout sets the timeout for transactions. If the transaction is not mined within this time, it will be considered failed.
// Default value is 5 minutes.
func (c *ClientBuilder) WithTransactionTimeout(timeout time.Duration) *ClientBuilder {
	c.config.Network.TxnTimeout = MustMakeDuration(timeout)
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].TxnTimeout = MustMakeDuration(timeout)
	}
	return c
}

// WithRpcDialTimeout sets the timeout for dialing the RPC server. If the connection is not established within this time, it will be considered failed.
// Default value is 1 minute.
func (c *ClientBuilder) WithRpcDialTimeout(timeout time.Duration) *ClientBuilder {
	c.config.Network.DialTimeout = MustMakeDuration(timeout)
	// defensive programming
	if len(c.config.Networks) == 0 {
		c.config.Networks = append(c.config.Networks, c.config.Network)
	} else {
		c.config.Networks[0].DialTimeout = MustMakeDuration(timeout)
	}
	return c
}

// WithEphemeralAddresses sets the number of ephemeral addresses to generate and the amount of funds to keep in the root private key.
// Default values are 0 for ephemeral addresses and 0 for root key funds buffer.
func (c *ClientBuilder) WithEphemeralAddresses(ephemeralAddressCount, rootKeyBufferAmount int64) *ClientBuilder {
	c.config.EphemeralAddrs = &ephemeralAddressCount
	c.config.RootKeyFundsBuffer = &rootKeyBufferAmount

	return c
}

// WithTracing sets the tracing level and outputs. Tracing level can be one of: "all", "reverted", "none". Outputs can be one or more of: "console", "dot" or "json".
// Default values are "reverted" and ["console", "dot"].
func (c *ClientBuilder) WithTracing(level string, outputs []string) *ClientBuilder {
	c.config.TracingLevel = level
	c.config.TraceOutputs = outputs
	return c
}

// WithProtections enables or disables nonce protection (fails, when key has a pending transaction and you try to submit another one) and node health check on startup.
// Default values are false for nonce protection and true for node health check.
func (c *ClientBuilder) WithProtections(pendingNonceProtectionEnabled, nodeHealthStartupCheck bool) *ClientBuilder {
	c.config.PendingNonceProtectionEnabled = pendingNonceProtectionEnabled
	c.config.CheckRpcHealthOnStart = nodeHealthStartupCheck
	return c
}

// WithArtifactsFolder sets the folder where the Seth artifacts such as DOT graphs or JSON will be saved.
// Default value is "seth_artifacts".
func (c *ClientBuilder) WithArtifactsFolder(folder string) *ClientBuilder {
	c.config.ArtifactsDir = folder
	return c
}

// WithNonceManager sets the rate limit for key sync, number of retries, timeout and retry delay.
// Default values are 10 calls per second, 3 retires, 60s timeout and 5s retry delay.
func (c *ClientBuilder) WithNonceManager(rateLimitSec int, retries uint, timeout, retryDelay time.Duration) *ClientBuilder {
	c.config.NonceManager = &NonceManagerCfg{
		KeySyncRateLimitSec: rateLimitSec,
		KeySyncRetries:      retries,
		KeySyncTimeout:      MustMakeDuration(timeout),
		KeySyncRetryDelay:   MustMakeDuration(retryDelay),
	}

	return c
}

// Build creates a new Client from the builder.
func (c *ClientBuilder) Build() (*Client, error) {
	return NewClientWithConfig(c.config)
}
