package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// PolygonEdgeMultinodeClient represents a multi-node, EVM compatible client for the Klaytn network
type PolygonEdgeMultinodeClient struct {
	*EthereumMultinodeClient
}

// PolygonEdgeClient represents a single node, EVM compatible client for the Polygon edge network
type PolygonEdgeClient struct {
	*EthereumClient
}

func (p *PolygonEdgeClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	privateKey, err := crypto.HexToECDSA(p.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return err
	}
	nonce, err := p.GetNonce(context.Background(), p.DefaultWallet.address)
	if err != nil {
		return err
	}
	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(p.GetChainID()), &types.DynamicFeeTx{
		ChainID: p.GetChainID(),
		Nonce:   nonce,
		To:      &to,
		Value:   utils.EtherToWei(amount),
		Gas:     100000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("From", p.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := p.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	return p.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (p *PolygonEdgeClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := p.TransactionOpts(p.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasTipCap = nil
	opts.GasPrice = nil

	contractAddress, transaction, contractInstance, err := deployer(opts, p.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = p.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", p.DefaultWallet.Address()).
		Str("Total Gas Cost", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", p.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}
