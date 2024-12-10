package simple_node_set

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"math/big"
)

func SendETH(client *ethclient.Client, privateKeyHex string, toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	wei := new(big.Int)
	amountWei := new(big.Float).Mul(amount, big.NewFloat(1e18))
	amountWei.Int(wei)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to fetch nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch gas price: %v", err)
	}
	gasLimit := uint64(21000) // Standard gas limit for ETH transfer

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), wei, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch chain ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}
	framework.L.Info().Msgf("Transaction sent: %s", signedTx.Hash().Hex())
	_, err = bind.WaitMined(context.Background(), client, signedTx)
	return err
}

// FundNodes funds Chainlink nodes with N ETH each
func FundNodes(c *ethclient.Client, nodes []*clclient.ChainlinkClient, pkey string, ethAmount float64) error {
	if ethAmount == 0 {
		return errors.New("funds_eth is 0, set some value in config, ex.: funds_eth = 30.0")
	}
	for _, cl := range nodes {
		ek, err := cl.ReadPrimaryETHKey()
		if err != nil {
			return err
		}
		if err := SendETH(c, pkey, ek.Attributes.Address, big.NewFloat(ethAmount)); err != nil {
			return fmt.Errorf("failed to fund CL node %s: %w", ek.Attributes.Address, err)
		}
	}
	return nil
}
