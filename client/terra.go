package client

import (
	"context"
	"encoding/json"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"
	"github.com/smartcontractkit/terra.go/client"
	"github.com/smartcontractkit/terra.go/key"
	"github.com/smartcontractkit/terra.go/msg"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	// DefaultTerraTXTimeout is default http client timeout
	DefaultTerraTXTimeout = 20 * time.Second
	// DefaultBroadcastMode is set to MODE_BLOCK it means when call returns, tx is mined and accepted in the next block
	DefaultBroadcastMode = tx.BroadcastMode_BROADCAST_MODE_BLOCK
	// EventAttrKeyCodeID code id
	EventAttrKeyCodeID = "code_id"
	// EventAttrKeyContractAddress contract address as bech32
	EventAttrKeyContractAddress = "contract_address"
)

// TerraLCDClient is terra lite chain client allowing to upload and interact with the contracts
type TerraLCDClient struct {
	*client.LCDClient
	BroadcastMode tx.BroadcastMode
	Sequence      uint64
	ID            int
	Config        *config.NetworkConfig
	RootPrivKey   key.PrivKey
	RootAddr      []byte
}

// Get gets default client as an interface{}
func (c *TerraLCDClient) Get() interface{} {
	return c
}

// GetNetworkName gets the ID of the chain that the clients are connected to
func (c *TerraLCDClient) GetNetworkName() string {
	return c.Config.ChainName
}

// GetID gets client ID, node number it's connected to
func (c *TerraLCDClient) GetID() int {
	return c.ID
}

// SetID sets client ID (node)
func (c *TerraLCDClient) SetID(id int) {
	c.ID = id
}

// SetDefaultClient sets default client to perform calls to the network
func (c *TerraLCDClient) SetDefaultClient(clientID int) error {
	// We are using SetDefaultClient and GetClients only for multinode networks to check reorgs,
	// but Terra uses Tendermint PBFT with an absolute finality
	return nil
}

// GetClients gets clients for all nodes connected
func (c *TerraLCDClient) GetClients() []BlockchainClient {
	return []BlockchainClient{c}
}

// SuggestGasPrice gets suggested gas price
func (c *TerraLCDClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	// client already have simulation for gas estimation by default turned on
	panic("implement me")
}

// HeaderHashByNumber gets header hash by block number
func (c *TerraLCDClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	panic("implement me")
}

// BlockNumber gets block number
func (c *TerraLCDClient) BlockNumber(ctx context.Context) (uint64, error) {
	panic("implement me")
}

// HeaderTimestampByNumber gets header timestamp by number
func (c *TerraLCDClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	panic("implement me")
}

// CalculateTxGas calculates tx gas cost accordingly gas used plus buffer, converts it to big.Float for funding
func (c *TerraLCDClient) CalculateTxGas(gasUsedValue *big.Int) (*big.Float, error) {
	panic("implement me")
}

// GasStats gets gas stats instance
func (c *TerraLCDClient) GasStats() *GasStats {
	panic("implement me")
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (c *TerraLCDClient) ParallelTransactions(enabled bool) {
	// need to check if it can be done after ws support through tendermint API
}

// Close tears down the current open Terra client
func (c *TerraLCDClient) Close() error {
	// no shutdown in underlying client for now,
	// implement shutdown for ws if it needed in the future
	return nil
}

// Currently only BroadcastMode_BROADCAST_MODE_BLOCK is working, we can get events directly from tx response after block is mined
// so we skip all subscription methods

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (c *TerraLCDClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (c *TerraLCDClient) DeleteHeaderEventSubscription(key string) {}

// WaitForEvents is a blocking function that waits for all event subscriptions for all clients
func (c *TerraLCDClient) WaitForEvents() error { return nil }

// NewTerraClient derives deployer key and creates new LCD client for Terra
func NewTerraClient(network BlockchainNetwork) (*TerraLCDClient, error) {
	cfg := network.Config()
	// Derive and set the key from first mnemonic, later keys can be changed by calling other methods with particular wallet
	privKeyBz, err := key.DerivePrivKeyBz(cfg.PrivateKeys[0], key.CreateHDPath(0, 0))
	if err != nil {
		return nil, err
	}
	privKey, err := key.PrivKeyGen(privKeyBz)
	if err != nil {
		return nil, err
	}
	return &TerraLCDClient{
		LCDClient: client.NewLCDClient(
			network.LocalURL(),
			cfg.ChainName,
			msg.NewDecCoinFromDec(cfg.Currency, msg.NewDecFromIntWithPrec(msg.NewInt(15), 2)),
			msg.NewDecFromIntWithPrec(msg.NewInt(15), 1),
			privKey,
			DefaultTerraTXTimeout,
		),
		Config:        cfg,
		RootPrivKey:   privKey,
		RootAddr:      privKey.PubKey().Address(),
		BroadcastMode: DefaultBroadcastMode,
	}, nil
}

// Instantiate instantiates particular uploaded WASM code using codeID and some instantiate message
func (c *TerraLCDClient) Instantiate(fromWallet BlockchainWallet, codeID uint64, instMsg interface{}) (string, error) {
	c.PrivKey = fromWallet.RawPrivateKey().(key.PrivKey)
	dat, err := json.Marshal(instMsg)
	if err != nil {
		return "", err
	}
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	log.Info().
		Uint64("CodeID", codeID).
		Msg("Instantiating contract")
	txBlockResp, err := c.SendTX(client.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgInstantiateContract(
				fromAddr,
				fromAddr,
				codeID,
				dat,
				msg.NewCoins(msg.NewInt64Coin(c.Config.Currency, 1000)),
			),
		},
	})
	if err != nil {
		return "", err
	}
	contractAddr, err := c.GetEventAttrValue(txBlockResp, EventAttrKeyContractAddress)
	if err != nil {
		return "", err
	}
	log.Info().
		Str("ContractAddress", contractAddr).
		Str("From", fromWallet.Address()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return contractAddr, nil
}

// SendTX signs and broadcast tx using default broadcast mode
func (c *TerraLCDClient) SendTX(txOpts client.CreateTxOptions) (*sdkTypes.TxResponse, error) {
	txn, err := c.CreateAndSignTx(context.Background(), txOpts)
	if err != nil {
		return nil, err
	}
	txBlockResp, err := c.Broadcast(context.Background(), txn, c.BroadcastMode)
	if err != nil {
		return nil, err
	}
	return txBlockResp, nil
}

// GetEventAttrValue gets attr value by key from sdkTypes.TxResponse
func (c *TerraLCDClient) GetEventAttrValue(tx *sdkTypes.TxResponse, attrKey string) (string, error) {
	for _, eventLog := range tx.Logs {
		for _, event := range eventLog.Events {
			for _, eventAttr := range event.Attributes {
				if eventAttr.Key == attrKey {
					return eventAttr.Value, nil
				}
			}
		}
	}
	return "", errors.New("No code_id found")
}

// DeployWASMCode deploys .wasm code file without instantiation from selected wallet and returns "code_id"
func (c *TerraLCDClient) DeployWASMCode(fromWallet BlockchainWallet, path string) (uint64, error) {
	c.PrivKey = fromWallet.RawPrivateKey().(key.PrivKey)
	dat, err := os.ReadFile(filepath.Join(tools.ProjectRoot, path))
	if err != nil {
		return 0, err
	}
	log.Info().
		Str("File", path).
		Str("From", fromWallet.Address()).
		Uint64("Sequence", c.Sequence).
		Msg("Deploying .wasm code")

	addr, _ := msg.AccAddressFromHex(fromWallet.Address())
	txBlockResp, err := c.SendTX(client.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgStoreCode(addr, dat),
		},
	})
	if err != nil {
		return 0, err
	}
	codeID, err := c.GetEventAttrValue(txBlockResp, EventAttrKeyCodeID)
	if err != nil {
		return 0, err
	}
	log.Info().
		Str("From", fromWallet.Address()).
		Str("CodeID", codeID).
		Interface("TX", txBlockResp).
		Msg("Result")
	codeUint, err := strconv.ParseUint(codeID, 10, 64)
	if err != nil {
		return 0, err
	}
	return codeUint, nil
}

// Fund funds a contracts with both native currency and LINK token
func (c *TerraLCDClient) Fund(fromWallet BlockchainWallet, toAddress string, nativeAmount, linkAmount *big.Float) error {
	c.PrivKey = fromWallet.RawPrivateKey().(key.PrivKey)
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	toAddrBech32, _ := msg.AccAddressFromBech32(toAddress)
	if big.NewFloat(0).Cmp(nativeAmount) != 0 {
		amount, _ := nativeAmount.Int64()
		txBlockResp, err := c.SendTX(client.CreateTxOptions{
			Msgs: []msg.Msg{
				msg.NewMsgSend(
					fromAddr,
					toAddrBech32,
					msg.NewCoins(msg.NewInt64Coin(c.Config.Currency, amount))),
			},
		})
		if err != nil {
			return err
		}
		log.Info().Str("From", fromWallet.Address()).Interface("TX", txBlockResp).Msg("Result")
	}
	return nil
}
