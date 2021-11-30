package contracts

import (
	"context"
)

type EthereumOCRv2 struct{}

func (e *EthereumOCRv2) ProgramAddress() string {
	panic("implement me")
}

func (e *EthereumOCRv2) SetOffChainConfig(_ OffChainAggregatorV2Config) error {
	panic("implement me")
}

func (e *EthereumOCRv2) DumpState() error {
	panic("implement me")
}

func (e *EthereumOCRv2) AuthorityAddr(s string) (string, error) {
	panic("implement me")
}

func (e *EthereumOCRv2) SetValidatorConfig(flaggingThreshold uint32, validatorAddr string) error {
	panic("implement me")
}

func (e *EthereumOCRv2) SetBilling(op uint32, tp uint32, controllerAddr string) error {
	panic("implement me")
}

func (e *EthereumOCRv2) GetContractData(ctx context.Context) (*OffchainAggregatorData, error) {
	panic("implement me")
}

func (e *EthereumOCRv2) SetOracles(cocParams OffChainAggregatorV2Config) error {
	panic("implement me")
}

func (e *EthereumOCRv2) RequestNewRound() error {
	panic("implement me")
}

func (e *EthereumOCRv2) Address() string {
	panic("implement me")
}

func (e *EthereumOCRv2) TransmissionsAddr() string {
	panic("implement me")
}

func (e *EthereumOCRv2) TransferOwnership(to string) error {
	panic("implement me")
}

func (e *EthereumOCRv2) GetLatestConfigDetails() (map[string]interface{}, error) {
	panic("implement me")
}

func (e *EthereumOCRv2) GetRoundData(roundID uint32) (map[string]interface{}, error) {
	panic("implement me")
}

func (e *EthereumOCRv2) GetOwedPayment(transmitterAddr string) (map[string]interface{}, error) {
	panic("implement me")
}
