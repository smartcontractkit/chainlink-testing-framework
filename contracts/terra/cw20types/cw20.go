package cw20types

import "github.com/smartcontractkit/terra.go/msg"

type InstantiateMsg struct {
	Name            string              `json:"name"`
	Symbol          string              `json:"symbol"`
	Decimals        int                 `json:"decimals"`
	InitialBalances []InitialBalanceMsg `json:"initial_balances"`
}

type InitialBalanceMsg struct {
	Address msg.AccAddress `json:"address"`
	Amount  string         `json:"amount"`
}

type ExecuteSendMsg struct {
	Send ExecuteSendTypeMsg `json:"send"`
}

type ExecuteSendTypeMsg struct {
	Contract msg.AccAddress `json:"contract"`
	Amount   string         `json:"amount"`
	Msg      []byte         `json:"msg"`
}

type ExecuteTransferMsg struct {
	Transfer ExecuteTransferTypeMsg `json:"transfer"`
}

type ExecuteTransferTypeMsg struct {
	Amount    string         `json:"amount"`
	Recipient msg.AccAddress `json:"recipient"`
}
