package config

import (
	"fmt"
	"strings"

	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/crypto"

	"github.com/smartcontractkit/integrations-framework/libocr/gethwrappers/exposedoffchainaggregator"

	"github.com/smartcontractkit/integrations-framework/libocr/offchainreporting/types"
)

func makeConfigDigestArgs() abi.Arguments {
	ABi, err := abi.JSON(strings.NewReader(
		exposedoffchainaggregator.ExposedOffchainAggregatorABI))
	if err != nil {
		// assertion
		panic(fmt.Sprintf("could not parse aggregator ABI: %s", err.Error()))
	}
	return ABi.Methods["exposedConfigDigestFromConfigData"].Inputs
}

var configDigestArgs = makeConfigDigestArgs()

func ConfigDigest(
	contractAddress common.Address,
	configCount uint64,
	oracles []common.Address,
	transmitters []common.Address,
	threshold uint8,
	encodedConfigVersion uint64,
	config []byte,
) types.ConfigDigest {
	msg, err := configDigestArgs.Pack(
		contractAddress,
		configCount,
		oracles,
		transmitters,
		threshold,
		encodedConfigVersion,
		config,
	)
	if err != nil {
		// assertion
		panic(err)
	}
	rawHash := crypto.Keccak256(msg)
	configDigest := types.ConfigDigest{}
	if n := copy(configDigest[:], rawHash); n != len(configDigest) {
		// assertion
		panic("copy too little data")
	}
	return configDigest
}
