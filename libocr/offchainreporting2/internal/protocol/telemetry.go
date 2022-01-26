package protocol

import (
	"github.com/smartcontractkit/integrations-framework/libocr/commontypes"
	"github.com/smartcontractkit/integrations-framework/libocr/offchainreporting2/types"
)

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		leader commontypes.OracleID,
	)
}
