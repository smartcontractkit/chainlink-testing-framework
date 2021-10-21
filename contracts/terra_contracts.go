package contracts

import (
	context "context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts/terra/actypes"
	"github.com/smartcontractkit/integrations-framework/contracts/terra/cw20types"
	"github.com/smartcontractkit/integrations-framework/contracts/terra/ocr2types"
	terraClient "github.com/smartcontractkit/terra.go/client"
	"github.com/smartcontractkit/terra.go/msg"
	"math/big"
)

// TerraWASMLinkToken represents a LinkToken deployed on terra as WASM
type TerraWASMLinkToken struct {
	client       *client.TerraLCDClient
	callerWallet client.BlockchainWallet
	address      string
}

func (t *TerraWASMLinkToken) Address() string {
	return t.address
}

func (t *TerraWASMLinkToken) Approve(fromWallet client.BlockchainWallet, to string, amount *big.Int) error {
	panic("implement me")
}

func (t *TerraWASMLinkToken) Transfer(fromWallet client.BlockchainWallet, to string, amount *big.Int) error {
	ownerAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	linkAddrBech32, _ := msg.AccAddressFromBech32(t.address)
	toAddr, _ := msg.AccAddressFromBech32(to)
	executeMsg := cw20types.ExecuteTransferMsg{
		Transfer: cw20types.ExecuteTransferTypeMsg{
			Amount:    amount.String(),
			Recipient: toAddr,
		}}
	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		return err
	}
	txBlockResp, err := t.client.SendTX(terraClient.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgExecuteContract(
				ownerAddr,
				linkAddrBech32,
				executeMsgBytes,
				msg.NewCoins(),
			),
		},
	})
	if err != nil {
		return err
	}
	log.Info().
		Str("ContractAddress", linkAddrBech32.String()).
		Str("From", ownerAddr.String()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return nil
}

func (t *TerraWASMLinkToken) BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	panic("implement me")
}

func (t *TerraWASMLinkToken) TransferAndCall(fromWallet client.BlockchainWallet, to string, amount *big.Int, data []byte) error {
	panic("implement me")
}

func (t *TerraWASMLinkToken) Fund(fromWallet client.BlockchainWallet, ethAmount *big.Float) error {
	if err := t.client.Fund(fromWallet, t.address, ethAmount, nil); err != nil {
		return err
	}
	return nil
}

func (t *TerraWASMLinkToken) Name(ctx context.Context) (string, error) {
	panic("implement me")
}

// TerraWASMAccessController represents a AccessController deployed on terra as WASM
type TerraWASMAccessController struct {
	client       *client.TerraLCDClient
	callerWallet client.BlockchainWallet
	address      string
}

func (t *TerraWASMAccessController) HasAccess(to string) (bool, error) {
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	toAddr, _ := msg.AccAddressFromHex(to)
	resp := &actypes.QueryHasAccessResponse{}
	err := t.client.QuerySmart(
		context.Background(),
		myContractAddr,
		actypes.QueryHasAccessMsg{HasAccess: actypes.QueryHasAccessTypeMsg{Address: toAddr}},
		resp,
	)
	if err != nil {
		return false, err
	}
	log.Debug().
		Interface("Response", resp).
		Msg("Query response")
	return false, nil
}

func (t *TerraWASMAccessController) RemoveAccess(fromWallet client.BlockchainWallet, to string) error {
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	toAddr, _ := msg.AccAddressFromHex(to)
	executeMsg := actypes.ExecuteRemoveAccessMsg{
		RemoveAccess: actypes.ExecuteRemoveAccessTypeMsg{
			Address: toAddr,
		}}
	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		return err
	}
	txBlockResp, err := t.client.SendTX(terraClient.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgExecuteContract(
				fromAddr,
				myContractAddr,
				executeMsgBytes,
				msg.NewCoins(),
			),
		},
	})
	if err != nil {
		return err
	}
	log.Info().
		Str("ContractAddress", myContractAddr.String()).
		Str("From", fromAddr.String()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return nil
}

func (t *TerraWASMAccessController) Address() string {
	return t.address
}

func (t *TerraWASMAccessController) AddAccess(fromWallet client.BlockchainWallet, to string) error {
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	toAddr, _ := msg.AccAddressFromHex(to)
	executeMsg := actypes.ExecuteAddAccessMsg{
		AddAccess: actypes.ExecuteAddAccessTypeMsg{
			Address: toAddr,
		}}
	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		return err
	}
	txBlockResp, err := t.client.SendTX(terraClient.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgExecuteContract(
				fromAddr,
				myContractAddr,
				executeMsgBytes,
				msg.NewCoins(),
			),
		},
	})
	if err != nil {
		return err
	}
	log.Info().
		Str("ContractAddress", myContractAddr.String()).
		Str("From", fromAddr.String()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return nil
}

// TerraWASMOCRv2 represents a OVR v2 contract deployed on terra as WASM
type TerraWASMOCRv2 struct {
	client       *client.TerraLCDClient
	callerWallet client.BlockchainWallet
	address      string
}

func (t *TerraWASMOCRv2) GetOwedPayment(transmitter string) (map[string]interface{}, error) {
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	transmitterAddr, _ := msg.AccAddressFromBech32(transmitter)
	resp := make(map[string]interface{})
	if err := t.client.QuerySmart(
		context.Background(),
		myContractAddr,
		ocr2types.QueryOwedPaymentMsg{
			OwedPayment: ocr2types.QueryOwedPaymentTypeMsg{
				Transmitter: transmitterAddr,
			},
		},
		&resp,
	); err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *TerraWASMOCRv2) GetRoundData(roundID uint32) (map[string]interface{}, error) {
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	resp := make(map[string]interface{})
	if err := t.client.QuerySmart(
		context.Background(),
		myContractAddr,
		ocr2types.QueryRoundDataMsg{
			RoundData: ocr2types.QueryRoundDataTypeMsg{
				RoundID: roundID,
			},
		},
		&resp,
	); err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *TerraWASMOCRv2) GetLatestConfigDetails() (map[string]interface{}, error) {
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	resp := make(map[string]interface{})
	if err := t.client.QuerySmart(context.Background(), myContractAddr, ocr2types.QueryLatestConfigDetails, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *TerraWASMOCRv2) SetBilling(fromWallet client.BlockchainWallet, observationPayment uint32, recommendedGasPrice uint32) error {
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	executeMsg := ocr2types.ExecuteSetBillingMsg{
		SetBilling: ocr2types.ExecuteSetBillingMsgType{
			Config: ocr2types.ExecuteSetBillingConfigMsgType{
				ObservationPayment:  observationPayment,
				RecommendedGasPrice: recommendedGasPrice,
			},
		}}
	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		return err
	}
	txBlockResp, err := t.client.SendTX(terraClient.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgExecuteContract(
				fromAddr,
				myContractAddr,
				executeMsgBytes,
				msg.NewCoins(),
			),
		},
	})
	if err != nil {
		return err
	}
	log.Info().
		Str("ContractAddress", myContractAddr.String()).
		Str("From", fromAddr.String()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return nil
}

func (t *TerraWASMOCRv2) TransferOwnership(fromWallet client.BlockchainWallet, to string) error {
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	toAddr, _ := msg.AccAddressFromHex(to)
	executeMsg := ocr2types.ExecuteTransferOwnershipMsg{
		TransferOwnership: ocr2types.ExecuteTransferOwnershipMsgType{
			To: toAddr,
		}}
	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		return err
	}
	txBlockResp, err := t.client.SendTX(terraClient.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgExecuteContract(
				fromAddr,
				myContractAddr,
				executeMsgBytes,
				msg.NewCoins(),
			),
		},
	})
	if err != nil {
		return err
	}
	log.Info().
		Str("ContractAddress", myContractAddr.String()).
		Str("From", fromAddr.String()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return nil
}

func (t *TerraWASMOCRv2) Address() string {
	return t.address
}

func (t *TerraWASMOCRv2) SetConfig(fromWallet client.BlockchainWallet) error {
	fromAddr, _ := msg.AccAddressFromHex(fromWallet.Address())
	myContractAddr, _ := msg.AccAddressFromBech32(t.address)
	executeMsg := ocr2types.ExecuteSetConfigMsg{
		SetConfig: ocr2types.ExecuteSetConfigMsgType{
			Signers:               [][]byte{{0, 0}},
			Transmitters:          []string{},
			F:                     1,
			OffchainConfig:        []byte{1, 2, 3},
			OffchainConfigVersion: 1,
			OnchainConfig:         []byte{4, 5, 6},
		}}
	executeMsgBytes, err := json.Marshal(executeMsg)
	if err != nil {
		return err
	}
	txBlockResp, err := t.client.SendTX(terraClient.CreateTxOptions{
		Msgs: []msg.Msg{
			msg.NewMsgExecuteContract(
				fromAddr,
				myContractAddr,
				executeMsgBytes,
				msg.NewCoins(),
			),
		},
	})
	if err != nil {
		return err
	}
	log.Info().
		Str("ContractAddress", myContractAddr.String()).
		Str("From", fromAddr.String()).
		Interface("TX", txBlockResp).
		Msg("Result")
	return nil
}
