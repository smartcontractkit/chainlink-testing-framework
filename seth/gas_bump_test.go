package seth_test

import (
	"github.com/ethereum/go-ethereum/common"
	link_token "github.com/smartcontractkit/seth/contracts/bind/link"
	"github.com/smartcontractkit/seth/contracts/bind/link_token_interface"
	"github.com/smartcontractkit/seth/test_utils"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/seth"
)

var oneEth = big.NewInt(1000000000000000000)
var zero int64 = 0

func TestGasBumping_Contract_Deployment_Legacy_SufficientBumping(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, oneEth)

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	gasBumps := 0

	// Set a low gas price and a short timeout
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasPrice = 1
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		Retries:     10,
		MaxGasPrice: 100000000,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return new(big.Int).Mul(gasPrice, big.NewInt(100))
		},
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	t.Cleanup(func() {
		configCopy.Network.GasPrice = 1_000_000_000
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	// Send a transaction with low gas price
	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one gas bump")
	require.Greater(t, data.Transaction.GasPrice().Int64(), int64(1), "expected gas price to be bumped")
}

func TestGasBumping_Contract_Deployment_Legacy_InsufficientBumping(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, oneEth)

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	gasBumps := 0

	// Set a low gas price and a short timeout
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasPrice = 1
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		Retries: 2,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return new(big.Int).Add(gasPrice, big.NewInt(1))
		},
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	t.Cleanup(func() {
		configCopy.Network.GasPrice = 1_000_000_000
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	// Send a transaction with a low gas price
	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))

	require.Error(t, err, "contract was deployed, but gas bumping shouldn't be sufficient to deploy it")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one gas bump")
}

func TestGasBumping_Contract_Deployment_Legacy_FailedBumping(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, oneEth)

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	gasBumps := 0

	// Set a low gas price and a short timeout
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasPrice = 1
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		Retries: 2,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return new(big.Int).Mul(gasPrice, big.NewInt(1000000000000))
		},
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	t.Cleanup(func() {
		configCopy.Network.GasPrice = 1_000_000_000
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	// Send a transaction with a low gas price and then bump it too high to be accepted
	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.Error(t, err, "contract was deployed, but gas bumping should be failing")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one gas bump")
}

func TestGasBumping_Contract_Deployment_Legacy_BumpingDisabled(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, oneEth)

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	gasBumps := 0

	// Set a low gas price and a short timeout, but disable gas bumping
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasPrice = 1
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return gasPrice
		},
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	t.Cleanup(func() {
		configCopy.Network.GasPrice = 1_000_000_000
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	// Send a transaction with a low gas price
	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.Error(t, err, "contract was deployed, but gas bumping is disabled")
	require.GreaterOrEqual(t, gasBumps, 0, "expected no gas bumps")
}

func TestGasBumping_Contract_Deployment_Legacy_CustomBumpingFunction(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, oneEth)

	customGasBumps := 0

	client, err := seth.NewClientBuilder().
		WithRpcUrl(c.Cfg.Network.URLs[0]).
		WithPrivateKeys([]string{newPk}).
		WithGasPriceEstimations(false, 0, "").
		WithEIP1559DynamicFees(false).
		WithLegacyGasPrice(1).
		WithTransactionTimeout(10*time.Second).
		WithProtections(false, false).
		WithGasBumping(5, 0, func(gasPrice *big.Int) *big.Int {
			customGasBumps++
			return new(big.Int).Mul(gasPrice, big.NewInt(512))
		}).
		Build()
	require.NoError(t, err)

	t.Cleanup(func() {
		client.Cfg.Network.GasPrice = 1_000_000_000
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
	})

	// we don't want to expose it in builder, but setting it to 0 (automatic gas limit estimation) doesn't work well with gas price of 1 wei
	client.Cfg.Network.GasLimit = 8_000_000

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract was not deployed")
	require.GreaterOrEqual(t, customGasBumps, 1, "expected at least one custom gas bump")
}

func TestGasBumping_Contract_Interaction_Legacy_MaxGas(t *testing.T) {
	spammer := test_utils.NewClientWithAddresses(t, 5, oneEth)

	configCopy, err := test_utils.CopyConfig(spammer.Cfg)
	require.NoError(t, err, "failed to copy config")

	newPk := test_utils.NewPrivateKeyWithFunds(t, spammer, oneEth)
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.EphemeralAddrs = &zero

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err, "failed to create client")

	contractAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")

	linkContract, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to instantiate contract")

	var gasPrices []*big.Int

	// Update config and set a low gas price and a short timeout
	client.Cfg.Network.GasPrice = 1
	client.Cfg.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	client.Cfg.GasBump = &seth.GasBumpConfig{
		Retries:     5,
		MaxGasPrice: 5, //after 2 retries gas price will be 5
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasPrices = append(gasPrices, gasPrice)
			return new(big.Int).Add(gasPrice, big.NewInt(2))
		},
	}

	// introduce some traffic, so that bumping is necessary to mine the transaction
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = spammer.DeployContract(spammer.NewTXKeyOpts(spammer.AnySyncedKey()), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
		}
	}()

	// Send a transaction with a low gas price
	_, _ = client.Decode(linkContract.Transfer(client.NewTXOpts(), client.Addresses[0], big.NewInt(1000000000000000000)))
	require.NoError(t, err, "failed to transfer tokens")
	require.GreaterOrEqual(t, len(gasPrices), 3, "expected 2 gas bumps")
	require.True(t, func() bool {
		for _, gasPrice := range gasPrices {
			if gasPrice.Cmp(big.NewInt(client.Cfg.GasBump.MaxGasPrice)) > 0 {
				return false
			}
		}
		return true
	}(), "at least one gas bump was too high")
}

func TestGasBumping_Contract_Interaction_EIP1559_MaxGas(t *testing.T) {
	spammer := test_utils.NewClientWithAddresses(t, 5, oneEth)

	configCopy, err := test_utils.CopyConfig(spammer.Cfg)
	require.NoError(t, err, "failed to copy config")

	newPk := test_utils.NewPrivateKeyWithFunds(t, spammer, oneEth)
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.EphemeralAddrs = &zero

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err, "failed to create client")

	contractAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")

	linkContract, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to instantiate contract")

	var gasPrices []*big.Int

	// Update config and set a low gas price and a short timeout
	client.Cfg.Network.GasFeeCap = 1
	client.Cfg.Network.GasTipCap = 1
	client.Cfg.Network.EIP1559DynamicFees = true
	client.Cfg.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	client.Cfg.GasBump = &seth.GasBumpConfig{
		Retries:     4,
		MaxGasPrice: 5, // for both fee and tip, which means that after a single bump, the gas price will be 5 and no more bumps should ever occur
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasPrices = append(gasPrices, gasPrice)
			return new(big.Int).Add(gasPrice, big.NewInt(2))
		},
	}

	// introduce some traffic, so that bumping is necessary to mine the transaction
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = spammer.DeployContract(spammer.NewTXKeyOpts(spammer.AnySyncedKey()), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
		}
	}()

	// Send a transaction with a low gas price
	_, _ = client.Decode(linkContract.Transfer(client.NewTXOpts(), client.Addresses[0], big.NewInt(1000000000000000000)))
	require.NoError(t, err, "failed to transfer tokens")
	require.GreaterOrEqual(t, len(gasPrices), 3, "expected at least 3 gas bumps")
	require.True(t, func() bool {
		for _, gasPrice := range gasPrices {
			// any other price higher than 2 would result in cumulated gas price (fee + cap) > 5
			if gasPrice.Cmp(big.NewInt(2)) > 0 {
				return false
			}
		}
		return true
	}(), "at least one gas bump was too high")
}

func TestGasBumping_Contract_Deployment_EIP_1559_SufficientBumping(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, oneEth)

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	gasBumps := 0

	// Set a low gas fee and tip cap and a short timeout
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasTipCap = 1
	configCopy.Network.GasFeeCap = 1
	configCopy.Network.EIP1559DynamicFees = true
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		Retries:     10,
		MaxGasPrice: 10000000,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return new(big.Int).Mul(gasPrice, big.NewInt(100))
		},
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	t.Cleanup(func() {
		client.Cfg.Network.GasTipCap = 50_000_000_000
		client.Cfg.Network.GasFeeCap = 100_000_000_000
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	// Send a transaction with a low gas price
	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one gas bump")
	require.Greater(t, data.Transaction.GasTipCap().Int64(), int64(1), "expected gas tip cap to be bumped")
	require.Greater(t, data.Transaction.GasFeeCap().Int64(), int64(1), "expected gas fee cap to be bumped")
}

func TestGasBumping_Contract_Deployment_EIP_1559_NonRootKey(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, big.NewInt(0).Mul(oneEth, big.NewInt(10)))

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	gasBumps := 0
	var one int64 = 1

	// Set a low gas fee and tip cap and a short timeout
	configCopy.EphemeralAddrs = &one
	configCopy.RootKeyFundsBuffer = &one
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasTipCap = 1
	configCopy.Network.GasFeeCap = 1
	configCopy.Network.EIP1559DynamicFees = true
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		Retries:     10,
		MaxGasPrice: 10000000,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return new(big.Int).Mul(gasPrice, big.NewInt(100))
		},
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := client.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(client, client.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	// Send a transaction with a low gas price
	data, err := client.DeployContract(client.NewTXKeyOpts(1), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed from key 1")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one gas bump")
	require.Greater(t, data.Transaction.GasTipCap().Int64(), int64(1), "expected gas tip cap to be bumped")
	require.Greater(t, data.Transaction.GasFeeCap().Int64(), int64(1), "expected gas fee cap to be bumped")
}

func TestGasBumping_Contract_Deployment_EIP_1559_UnknownKey(t *testing.T) {
	c := newClient(t)
	newPk := test_utils.NewPrivateKeyWithFunds(t, c, big.NewInt(0).Mul(oneEth, big.NewInt(10)))

	configCopy, err := test_utils.CopyConfig(c.Cfg)
	require.NoError(t, err, "failed to copy config")

	var one int64 = 1

	// Set a low gas fee and tip cap and a short timeout
	configCopy.EphemeralAddrs = &one
	configCopy.RootKeyFundsBuffer = &one
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.Network.GasTipCap = 1
	configCopy.Network.GasFeeCap = 1
	configCopy.Network.EIP1559DynamicFees = true
	configCopy.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	configCopy.GasBump = &seth.GasBumpConfig{
		Retries: 2,
	}

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err)

	removedAddress := client.Addresses[1]

	gasBumps := 0

	client.Cfg.GasBump.StrategyFn = func(gasPrice *big.Int) *big.Int {
		// remove address from client to simulate an unlikely situation, where we try to bump a transaction with having sender's private key
		client.Addresses = client.Addresses[:1]
		gasBumps++
		return gasPrice
	}

	t.Cleanup(func() {
		client.Addresses = append(client.Addresses, removedAddress)
		err := client.NonceManager.UpdateNonces()
		require.NoError(t, err, "failed to update nonces")
		err = seth.ReturnFunds(client, client.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
		err = test_utils.TransferAllFundsBetweenKeyAndAddress(client, 0, c.Addresses[0])
		require.NoError(t, err, "failed to transfer funds back to original root key")
	})

	contractAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	_, err = client.DeployContract(client.NewTXKeyOpts(1), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.Error(t, err, "contract was deployed from unknown key")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one gas bump attempt")
}

func TestGasBumping_Contract_Interaction_Legacy_SufficientBumping(t *testing.T) {
	spammer := test_utils.NewClientWithAddresses(t, 5, oneEth)

	configCopy, err := test_utils.CopyConfig(spammer.Cfg)
	require.NoError(t, err, "failed to copy config")

	newPk := test_utils.NewPrivateKeyWithFunds(t, spammer, oneEth)
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.EphemeralAddrs = &zero

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err, "failed to create client")

	contractAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")

	linkContract, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to instantiate contract")

	gasBumps := 0

	// Update config and set a low gas price and a short timeout
	client.Cfg.Network.GasPrice = 1
	client.Cfg.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	client.Cfg.GasBump = &seth.GasBumpConfig{
		Retries:     10,
		MaxGasPrice: 10000000,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			return new(big.Int).Mul(gasPrice, big.NewInt(100))
		},
	}

	// introduce some traffic, so that bumping is necessary to mine the transaction
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = spammer.DeployContract(spammer.NewTXKeyOpts(spammer.AnySyncedKey()), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
		}
	}()

	// Send a transaction with a low gas price
	_, err = client.Decode(linkContract.Transfer(client.NewTXOpts(), client.Addresses[0], big.NewInt(1000000000000000000)))
	require.NoError(t, err, "failed to mint tokens")
	require.GreaterOrEqual(t, gasBumps, 1, "expected at least one transaction gas bump")
}

func TestGasBumping_Contract_Interaction_Legacy_BumpingDisabled(t *testing.T) {
	spammer := test_utils.NewClientWithAddresses(t, 5, oneEth)

	configCopy, err := test_utils.CopyConfig(spammer.Cfg)
	require.NoError(t, err, "failed to copy config")

	newPk := test_utils.NewPrivateKeyWithFunds(t, spammer, oneEth)
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.EphemeralAddrs = &zero

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err, "failed to create client")

	contractAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")

	linkContract, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to instantiate contract")

	gasBumps := 0

	// Update config and set a low gas price and a short timeout
	client.Cfg.Network.GasPrice = 1
	client.Cfg.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	client.Cfg.GasBump = &seth.GasBumpConfig{
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			// do not bump anything
			return gasPrice
		},
	}

	// introduce some traffic, so that bumping is necessary to mine the transaction
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = spammer.DeployContract(spammer.NewTXKeyOpts(spammer.AnySyncedKey()), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
		}
	}()

	// Send a transaction with a low gas price
	_, err = client.Decode(linkContract.Transfer(client.NewTXOpts(), client.Addresses[0], big.NewInt(1000000000000000000)))
	require.Error(t, err, "did not fail to transfer tokens, even though gas bumping is disabled")
	require.Equal(t, gasBumps, 0, "expected no gas bumps")
}

func TestGasBumping_Contract_Interaction_Legacy_FailedBumping(t *testing.T) {
	spammer := test_utils.NewClientWithAddresses(t, 5, oneEth)

	configCopy, err := test_utils.CopyConfig(spammer.Cfg)
	require.NoError(t, err, "failed to copy config")

	newPk := test_utils.NewPrivateKeyWithFunds(t, spammer, oneEth)
	configCopy.Network.PrivateKeys = []string{newPk}
	configCopy.EphemeralAddrs = &zero

	client, err := seth.NewClientWithConfig(configCopy)
	require.NoError(t, err, "failed to create client")

	contractAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *contractAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	require.NoError(t, err, "contract wasn't deployed")

	linkContract, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to instantiate contract")

	gasBumps := 0

	// Update config and set a low gas price and a short timeout
	client.Cfg.Network.GasPrice = 1
	client.Cfg.Network.TxnTimeout = seth.MustMakeDuration(10 * time.Second)
	client.Cfg.GasBump = &seth.GasBumpConfig{
		Retries: 3,
		StrategyFn: func(gasPrice *big.Int) *big.Int {
			gasBumps++
			// this results in a gas bump that is too high to be accepted
			return new(big.Int).Mul(gasPrice, big.NewInt(1000000000000))
		},
	}

	// introduce some traffic, so that bumping is necessary to mine the transaction
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = spammer.DeployContract(spammer.NewTXKeyOpts(spammer.AnySyncedKey()), "LinkToken", *contractAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
		}
	}()

	// Send a transaction with a low gas price
	_, err = client.Decode(linkContract.Transfer(client.NewTXOpts(), client.Addresses[0], big.NewInt(1000000000000000000)))
	require.Error(t, err, "did not fail to transfer tokens, even though gas bumping is disabled")
	require.Equal(t, 3, gasBumps, "expected 2 gas bumps")
}
