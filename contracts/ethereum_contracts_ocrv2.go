package contracts

import (
	"context"
)

type EthereumOCRv2Store struct{}

func (e *EthereumOCRv2Store) TransmissionsAddress() string {
	panic("implement me")
}

func (e *EthereumOCRv2Store) GetLatestRoundData() (uint64, uint64, uint64, error) {
	panic("implement me")
}

func (e *EthereumOCRv2Store) SetWriter(writerAuthority string) error {
	panic("implement me")
}

func (e *EthereumOCRv2Store) CreateFeed(desc string, decimals uint8, granylarity int, liveLength int) error {
	panic("implement me")
}

func (e *EthereumOCRv2Store) SetValidatorConfig(flaggingThreshold uint32) error {
	panic("implement me")
}

func (e *EthereumOCRv2Store) Address() string {
	panic("implement me")
}

func (e *EthereumOCRv2Store) ProgramAddress() string {
	panic("implement me")
}

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

func (e *EthereumOCRv2) TransferOwnership(to string) error {
	panic("implement me")
}

func (e *EthereumOCRv2) GetLatestConfigDetails() (map[string]interface{}, error) {
	panic("implement me")
}

func (e *EthereumOCRv2) GetOwedPayment(transmitterAddr string) (map[string]interface{}, error) {
	panic("implement me")
}
