package ocr2types

import (
	"fmt"
	"github.com/smartcontractkit/terra.go/msg"
	"strings"
)

const (
	QueryLatestConfigDetails        = "latest_config_details"
	QueryTransmitters               = "transmitters"
	QueryLatestTransmissionDetails  = "latest_transmission_details"
	QueryLatestConfigDigestAndEpoch = "latest_config_digest_and_epoch"
	QueryDescription                = "description"
	QueryDecimals                   = "decimals"
	QueryLatestRoundData            = "latest_round_data"
	QueryLinkToken                  = "link_token"
	QueryBilling                    = "billing"
	QueryBillingAccessController    = "billing_access_controller"
	QueryRequesterAccessController  = "requester_access_controller"
	QueryLinkAvailableForPayment    = "link_available_for_payment"
)

type QueryOwedPaymentMsg struct {
	OwedPayment QueryOwedPaymentTypeMsg `json:"owed_payment"`
}

type QueryOwedPaymentTypeMsg struct {
	Transmitter msg.AccAddress `json:"transmitter"`
}

type QueryRoundDataMsg struct {
	RoundData QueryRoundDataTypeMsg `json:"round_data"`
}

type QueryRoundDataTypeMsg struct {
	RoundID uint32 `json:"round_id"`
}

type OCRv2InstantiateMsg struct {
	BillingAccessController   msg.AccAddress `json:"billing_access_controller"`
	RequesterAccessController msg.AccAddress `json:"requester_access_controller"`
	LinkToken                 msg.AccAddress `json:"link_token"`
	Decimals                  uint8          `json:"decimals"`
	Description               string         `json:"description"`
	MinAnswer                 string         `json:"min_answer"`
	MaxAnswer                 string         `json:"max_answer"`
}

type ExecuteSetBillingMsg struct {
	SetBilling ExecuteSetBillingMsgType `json:"set_billing"`
}

type ExecuteSetBillingMsgType struct {
	Config ExecuteSetBillingConfigMsgType `json:"config"`
}

type ExecuteSetBillingConfigMsgType struct {
	ObservationPayment  uint32 `json:"observation_payment"`
	RecommendedGasPrice uint32 `json:"recommended_gas_price"`
}

type ExecuteTransferOwnershipMsg struct {
	TransferOwnership ExecuteTransferOwnershipMsgType `json:"transfer_ownership"`
}

type ExecuteTransferOwnershipMsgType struct {
	To msg.AccAddress `json:"to"`
}

type ExecuteSetConfigMsg struct {
	SetConfig ExecuteSetConfigMsgType `json:"set_config"`
}

type ExecuteSetConfigMsgType struct {
	Signers               ByteArrayArray `json:"signers"`
	Transmitters          []string       `json:"transmitters"`
	F                     uint8          `json:"f"`
	OnchainConfig         ByteArray      `json:"onchain_config"`
	OffchainConfigVersion uint64         `json:"offchain_config_version"`
	OffchainConfig        ByteArray      `json:"offchain_config"`
}

type ByteArray []byte

func (b ByteArray) MarshalJSON() ([]byte, error) {
	var result string
	if b == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
	}
	return []byte(result), nil
}

type ByteArrayArray [][]byte

func (b ByteArrayArray) MarshalJSON() ([]byte, error) {
	var result string
	if b == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
	}
	return []byte(result), nil
}
