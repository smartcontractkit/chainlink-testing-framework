package simple_node_set

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	er "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

func SendETH(client *ethclient.Client, privateKeyHex string, toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return er.Wrap(err, "failed to parse private key")
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
		return er.Wrap(err, "failed to fetch nonce")
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return er.Wrap(err, "failed to fetch gas price")
	}
	gasLimit := uint64(21000) // Standard gas limit for ETH transfer

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), wei, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return er.Wrap(err, "failed to fetch chain ID")
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return er.Wrap(err, "failed to sign transaction")
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return er.Wrap(err, "failed to send transaction")
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
	chainID, err := c.ChainID(context.Background())
	if err != nil {
		return er.Wrap(err, "failed to fetch chain ID")
	}
	for _, cl := range nodes {
		ek, err := cl.ReadPrimaryETHKey(chainID.String())
		if err != nil {
			return err
		}
		if err := SendETH(c, pkey, ek.Attributes.Address, big.NewFloat(ethAmount)); err != nil {
			return er.Wrapf(err, "failed to fund CL node %s", ek.Attributes.Address)
		}
	}
	return nil
}
